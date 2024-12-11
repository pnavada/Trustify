package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"trustify/blockchain"
	"trustify/config"
	"trustify/cryptography"
	"trustify/logger"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
)

type Node struct {
	Config       *config.Config
	Wallet       *blockchain.Wallet
	Blockchain   *blockchain.Blockchain
	Mempool      *blockchain.Mempool
	UTXOSet      *blockchain.UTXOSet
	Miner        *blockchain.Miner
	Peers        []string
	TCPEgress    *ConnectionPool
	ReadChannel  chan InboundMessage
	WriteChannel chan OutboundMessage
	hostName     string
}

// Context - blockchain package files

func NewNode(cfg *config.Config) *Node {
	// Initialize node with wallet, blockchain, mempool, UTXOSet
	// The node's wallet is initialized using the public key, private key and bitcoin address
	// from the configuration
	// The node's blockchain should be initialized with the genesis block
	// from the configuration file
	// The rewards should be added to wallet as unspent transaction outputs
	// The transaction id is the hash of the genesis block combined with the index
	// The UTXOSet is initialized with the genesis block's transactions
	// The miner is initialized with the node's blockchain and mempool
	// The peers are the list of nodes except the host under the nodes section of the configuration

	me, err := os.Hostname()
	if err != nil {
		logger.ErrorLogger.Println("Failed to get hostname:", err)
		return nil
	}

	utxoSet := blockchain.NewUTXOSet()

	cfgNode := cfg.Nodes[me]
	wallet, err := blockchain.NewWallet(cfgNode.Wallet.PrivateKey) // Need to get self private key

	if err != nil {
		logger.ErrorLogger.Println("Failed to initialize wallet:", err)
		return nil
	}

	mempool := blockchain.NewMempool()
	host, err := libp2p.New() // TODO: Verify if this is the correct way to initialize host

	getBlocksProtocol := blockchain.NewGetBlocksProtocol(
		host,
		cfg.BlockchainSettings.Protocols.GetBlocks.Timeout,
	)

	// BestBlocksChannel := make(chan *GetBlocksResponse)
	chain, err := blockchain.NewBlockchain(&cfg.GenesisBlock, &cfg.BlockchainSettings, utxoSet, getBlocksProtocol, mempool) // BestBlocksChannel

	if err != nil {
		logger.ErrorLogger.Println("Failed to initialize blockchain:", err)
		return nil
	}

	// Initialize UTXOSet with genesis block's transactions
	for _, tx := range chain.Ledger[0].Transactions {
		for _, output := range tx.Outputs {
			utxoSet.Add(output) // TODO: create a copy of output?
			// print output address and wallet bitcoin address
			logger.InfoLogger.Printf("Output address: %x, Wallet Bitcoin address: %x\n", output.Address, wallet.BitcoinAddress)

			if bytes.Equal(output.Address, wallet.BitcoinAddress) {
				wallet.UTXOs = append(wallet.UTXOs, output)
			}
		}
	}

	miner := blockchain.NewMiner(chain, mempool, wallet)

	// Initialize peers
	var peers []string
	for nodeName := range cfg.Nodes {
		if nodeName != me {
			peers = append(peers, nodeName)
		}
	}

	node := &Node{
		Config:       cfg,
		Wallet:       wallet,
		Blockchain:   chain,
		Mempool:      mempool,
		UTXOSet:      utxoSet,
		Miner:        miner,
		Peers:        peers,
		TCPEgress:    NewTCPConnectionPool(8080, Outgoing),
		ReadChannel:  make(chan InboundMessage),
		WriteChannel: make(chan OutboundMessage),
		hostName:     me,
	}

	logger.InfoLogger.Printf("Node initialized: %+v\n", node)

	// logger.InfoLogger.Println("Node initialized with address:", wallet.BitcoinAddress)
	return node
}

// HandleGetBlocksRequest handles incoming GetBlocks requests from peers
func (n *Node) HandleGetBlocksRequest(s network.Stream) {
	defer s.Close()

	var request blockchain.GetBlocksRequest
	if err := blockchain.Receive(s, &request); err != nil {
		logger.ErrorLogger.Printf("Failed to receive request: %v", err)
		return
	}

	blocks, err := n.Blockchain.GetBlocksSinceHash(request.LastKnownHash)
	if err != nil {
		logger.ErrorLogger.Printf("Error retrieving blocks since hash: %v", err)
		blockchain.Send(s, blockchain.GetBlocksResponse{Success: false})
		return
	}

	response := blockchain.GetBlocksResponse{
		Blocks:  blocks,
		Success: true,
	}

	if err := blockchain.Send(s, response); err != nil {
		logger.ErrorLogger.Printf("Failed to send response: %v", err)
	}
}

func (n *Node) StartMining() {
	for {
		blockSize := n.Config.BlockchainSettings.BlockSize
		// logger.InfoLogger.Println("Number of transaction in mempool:", n.Mempool.Transactions.Len())
		if n.Mempool.Transactions.Len() >= blockSize {
			block, err := n.Miner.MineBlock(blockSize)
			if err != nil {
				logger.ErrorLogger.Println("Failed to mine block:", err)
				continue
			}
			if block != nil {
				n.BroadcastBlock(block)
			} else {
				logger.ErrorLogger.Println("Block is nil")
			}
		}
	}
}

func (n *Node) Start() {
	// Start node operations: networking, transaction processing, mining - concurrent
	// A node should start listening for incoming transactions and blocks on a specified port
	// The node should create an outgoing connection to broadcast data over the network
	// The nodes should start mining to add new blocks to blockchain
	// Add additional methods or files as needed maintaining separation of concerns

	// Start networking, transaction processing, mining
	go n.ListenForTCPConnections()

	// Listen UDP
	go ReceiveMessages(n.ReadChannel)

	time.Sleep(5 * time.Second)
	go n.StartMining()

	// TEST
	for _, peer := range n.Peers {
		go n.SendMessageToHost(peer, []byte("hello from "+n.hostName))
	}

	// Handle transactions from configuration file
	for _, tx := range n.Config.Nodes[n.hostName].Transactions {
		go n.handleConfigTransaction(tx)
	}

	n.HandleMessages()
}

func (n *Node) handleConfigTransaction(tx config.ConfigTransaction) {
	// Wait for the specified delay
	logger.InfoLogger.Printf("Waiting for %d seconds before processing transaction", tx.Delay)
	time.Sleep(time.Duration(tx.Delay) * time.Second)

	// Construct the transaction based on the configuration
	var transaction *blockchain.Transaction

	switch tx.Type {
	case "purchase":
		logger.InfoLogger.Println("Processing purchase transaction")
		// Create a purchase transaction
		transaction = blockchain.NewPurchaseTransaction(
			n.Wallet,
			tx.SellerAddress,
			tx.Amount,
			tx.Fee,
			tx.ProductID,
		)
		if transaction == nil {
			logger.ErrorLogger.Println("Failed to create purchase transaction")
			return
		}
	case "review":
		logger.InfoLogger.Println("Processing review transaction")
		// Create a review transaction
		transaction = blockchain.NewReviewTransaction(
			n.Wallet,
			tx.ProductID,
			tx.Rating,
		)
		if transaction == nil {
			logger.ErrorLogger.Println("Failed to create review transaction")
			return
		}

	default:
		logger.ErrorLogger.Printf("Unknown transaction type: %s", tx.Type)
		return
	}

	// Print the constructed transaction
	logger.InfoLogger.Printf("Constructed transaction: %+v", transaction)

	// Validate the transaction
	if err := n.Blockchain.ValidateTransaction(transaction, n.UTXOSet); err != nil {
		logger.ErrorLogger.Println("Transaction validation failed:", err)
		return
	}

	// Add the transaction to the mempool
	n.Mempool.AddTransaction(transaction)

	// Broadcast the transaction to the network
	go n.BroadcastTransaction(transaction)

}

// Network communication
func (n *Node) ListenForTCPConnections() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go n.HandleTCPConnection(conn)
	}
}

func (node *Node) HandleTCPConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from TCP connection:", err)
			}
			break
		}
		node.ReadChannel <- InboundMessage{
			Data:   buffer[:n],
			Sender: conn.RemoteAddr(),
		}
	}
}

func (node *Node) SendMessageToHost(host string, data []byte) {
	addr, err := GetAddrFromHostname(host)
	if err != nil {
		fmt.Printf("error resolving address for host %s: %v\n", host, err)
		return
	}

	node.WriteChannel <- OutboundMessage{
		Data:      data,
		Recipient: addr,
	}
}

func GetAddrFromHostname(hostname string) (net.Addr, error) {
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("no addresses found for hostname: %s", hostname)
	}
	return &net.TCPAddr{IP: addrs[0]}, nil
}

func (n *Node) HandleMessages() {
	for {
		select {
		case inboundMessage := <-n.ReadChannel:
			go n.handleIncomingMessage(inboundMessage)
		case outboundMessage := <-n.WriteChannel:
			conn, _ := n.TCPEgress.Get(outboundMessage.Recipient)
			tcpConn := conn.(net.Conn)
			tcpConn.Write(outboundMessage.Data)
		}
	}
}

func (n *Node) handleIncomingMessage(message InboundMessage) {
	if len(message.Data) < 5 { // 1 byte for type + 4 bytes for length
		logger.ErrorLogger.Println("Received block data too short")
		return
	}

	messageType := message.Data[0]
	payload := message.Data[1:]

	switch messageType {
	case MessageTypeTransaction:
		logger.InfoLogger.Printf("Received transaction from %s\n", message.Sender)
		tx, signature, publicKey, err := deserializeTransactionMessage(payload)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to deserialize transaction from %s\n", message.Sender)
			return
		}
		n.HandleIncomingTransaction(tx, signature, publicKey)
	case MessageTypeBlock:
		logger.InfoLogger.Printf("Received block from sender %v with payload %v\n", message.Sender, payload)

		// Read the length of the serialized block
		length := binary.BigEndian.Uint32(message.Data[1:5])
		if int(length) > len(message.Data[5:]) {
			logger.ErrorLogger.Printf("Declared block length %d exceeds received data %d\n", length, len(message.Data[5:]))
			return
		}

		serializedBlock := message.Data[5 : 5+length]
		block := blockchain.DeserializeBlock(serializedBlock)
		if block == nil {
			logger.ErrorLogger.Println("Failed to deserialize block from UDP data")
		}
		n.HandleIncomingBlock(block)
	default:
		logger.ErrorLogger.Printf("Unknown message type %v , msg: %v from %v\n", messageType, string(message.Data), message.Sender)
	}
}

func (n *Node) BroadcastBlock(block *blockchain.Block) {
	// Broadcast block to the network
	// Broadcast the block data over the network to all peers
	// do not use peer to peer multicasting instead use broadcasting

	// Print block data
	logger.InfoLogger.Printf("Broadcasting block: %+v\n", block)

	// Log the block and block header
	logger.InfoLogger.Printf("Block: %+v\n", block)
	logger.InfoLogger.Printf("Block Header: %+v\n", block.Header)

	// Network broadcasting
	err := SendBlock(block)
	if err != nil {
		logger.ErrorLogger.Println("Failed to broadcast block:", err)
	}

}

func (n *Node) BroadcastTransaction(tx *blockchain.Transaction) {
	// Broadcast the transaction data over the network to all the peers.
	// Do not use peer to peer multicasting instead use broadcasting

	// Sign the transaction before broadcasting
	// Will be sending transcation, signature and public key

	Serialized := blockchain.SerializeTransaction(tx)
	hashed := blockchain.HashObject(Serialized)
	signature, err := cryptography.SignMessage(n.Wallet.PrivateKey, hashed)

	if err != nil {
		logger.ErrorLogger.Println("Failed to sign transaction:", err)
		return
	}

	// Serialize public key for broadcasting
	publicKey, err := cryptography.SerializePublicKey(n.Wallet.PublicKey)
	if err != nil {
		logger.ErrorLogger.Println("Failed to serialize public key:", err)
		return
	}

	// Print tx, signature and public key
	logger.InfoLogger.Printf("Transaction: %+v, Signature : %+v and Public Key: %+v are ready for broadcasting\n", tx.Data, signature, publicKey)

	// Parse into purchase transaction and log the contents
	if purchaseTx, ok := tx.Data.(*blockchain.PurchaseTransactionData); ok {
		logger.InfoLogger.Printf("Purchase Transaction Details: Seller Address: %x, Amount: %d, ProductID: %s\n", purchaseTx.SellerAddress, purchaseTx.Amount, purchaseTx.ProductID)
	} else {
		logger.InfoLogger.Println("Transaction is not a purchase transaction")
	}

	// Network broadcasting
	err = SendTransaction(tx, signature, publicKey)
	if err != nil {
		logger.ErrorLogger.Println("Failed to broadcast transaction:", err)
	}
}

func (n *Node) HandleIncomingTransaction(tx *blockchain.Transaction, signature []byte, publicKey []byte) {
	// Handle incoming transaction
	// The peers are responsible for validating these transactions.
	// They verify if the sender bitcoin address is valid and if the transaction is signed by the sender.
	// For purchase transactions, they validate if the input transactions are in UXTO set
	// (if they can be spent and there is no double spending).
	// They validate the seller and buyer bitcoin address for purchase transaction.
	// For review transaction, they check if
	// the corresponding purchase transaction is in the ledger.
	// A user can only add 1 review
	// for a product. This verification is also done.
	// Also, the reviewer bitcoin address is verified.
	// If all the checks pass, the transaction is added to the memory pool.
	// Add additional methods or files as needed maintaining separation of concerns

	// Verify the signature
	serialized := blockchain.SerializeTransaction(tx)
	hashed := blockchain.HashObject(serialized)

	// Deserialize the public key
	pubKey, err := cryptography.DeserializePublicKey(publicKey)
	if err != nil {
		logger.ErrorLogger.Println("Failed to deserialize public key:", err)
		return
	}

	if !cryptography.VerifySignature(pubKey, hashed, signature) {
		logger.ErrorLogger.Println("Transaction signature verification failed")
		return
	}

	// Validate the transaction
	if err := n.Blockchain.ValidateTransaction(tx, n.UTXOSet); err != nil {
		logger.ErrorLogger.Println("Transaction validation failed:", err)
		return
	}

	tran, exists := n.Mempool.HasTransaction(tx)
	if exists {
		logger.InfoLogger.Println("Transaction already in mempool")
		// log the details of both the transactions
		logger.InfoLogger.Printf("Existing transaction: %+v\n", tran)
		logger.InfoLogger.Printf("New transaction: %+v\n", tx)
		return
	} else {
		// Add the transaction to the mempool
		n.Mempool.AddTransaction(tx)

		// Broadcast the transaction to the network
		go n.BroadcastTransaction(tx)
	}
}

func (n *Node) HandleIncomingBlock(block *blockchain.Block) error {
	// Handle incoming block
	// Verify the block coming, verifying the transactions in it and if its the succeeding block
	// Based on the validation, carrry out the next operation -
	// If it is validated - add it to the ledger and make sure the mining is done for the next block
	// If not, then initiate the getBlocks protocol to figure out the missing blocks and act accordingly or weather to drop this block
	// Add additional methods or files as needed maintaining separation of concerns

	// Try to add the block. validation happens inside it
	if err := n.Blockchain.AddBlock(block); err != nil {
		logger.ErrorLogger.Println("Failed to add incoming block:", err)
		return err
	}

	logger.InfoLogger.Println("Incoming block added to blockchain:", block.Header.BlockHash)
	return nil
}
