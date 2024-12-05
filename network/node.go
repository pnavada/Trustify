package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"trustify/blockchain"
	"trustify/config"
	"trustify/logger"
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

	cfgNode := cfg.Nodes[me]
	wallet := blockchain.NewWallet([]byte(cfgNode.Wallet.PrivateKey), []byte(cfgNode.Wallet.PublicKey), []byte(cfgNode.Wallet.BitcoinAddress)) // Need to get self private key
	chain, err := blockchain.NewBlockchain(&cfg.GenesisBlock, &cfg.BlockchainSettings)
	if err != nil {
		logger.ErrorLogger.Println("Failed to initialize blockchain:", err)
		return nil
	}

	mempool := blockchain.NewMempool()
	utxoSet := blockchain.NewUTXOSet()

	// Initialize UTXOSet with genesis block's transactions
	for _, tx := range chain.Ledger[0].Transactions {

		for _, output := range tx.Outputs {
			utxoSet.Add(output) // TODO: create a copy of output
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

func (n *Node) StartMining() {
	for {
		blockSize := n.Config.BlockchainSettings.BlockSize
		if n.Mempool.Transactions.Len() >= blockSize {
			block, _ := n.Miner.MineBlock(blockSize)
			n.BroadcastBlock(block)
		}
	}
}

func (n *Node) CommitBlocks() {
	for {

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
	time.Sleep(5 * time.Second)
	go n.StartMining()

	// TEST
	for _, peer := range n.Peers {
		go n.SendMessageToHost(peer, []byte("hello"))
	}

	// Handle transactions from configuration file
	for _, tx := range n.Config.Nodes[n.hostName].Transactions {
		go n.handleConfigTransaction(tx)
	}
	go n.CommitBlocks()

	n.HandleMessages()
}

func (n *Node) handleConfigTransaction(tx config.ConfigTransaction) {
	// Wait for the specified delay
	logger.InfoLogger.Printf("Waiting for %d seconds before processing transaction", tx.Delay)
	time.Sleep(time.Duration(tx.Delay) * time.Second)

	// Construct the transaction based on the configuration
	var transaction *blockchain.Transaction
	var err error

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
		if err != nil {
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
		if err != nil {
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
	n.BroadcastTransaction(*transaction)
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
			// Handle incoming message
			fmt.Println("Received message:", string(inboundMessage.Data))
		case outboundMessage := <-n.WriteChannel:
			conn, _ := n.TCPEgress.Get(outboundMessage.Recipient)
			tcpConn := conn.(net.Conn)
			tcpConn.Write(outboundMessage.Data)
		}
	}
}

func (n *Node) BroadcastBlock(block *blockchain.Block) {
	// Broadcast block to the network
	// Broadcast the block data over the network to all peers
	// do not use peer to peer multicasting instead use broadcasting

	// Serialize and broadcast the block to peers
	// data := utils.SerializeBlock(block)
	// for _, peer := range n.Peers {
	// 	go n.sendDataToPeer(peer, data)
	// }
	// logger.InfoLogger.Println("Block broadcasted:", block.Header.BlockHash)
}

func (n *Node) BroadcastTransaction(tx blockchain.Transaction) {
	// Broadcast transaction to the network
	// Broadcast the transaction data over the network to all the peers
	// do not use peer to peer multicasting instead use broadcasting

	// Serialize and broadcast the transaction to peers
	// data := utils.SerializeTransaction(tx)
	// for _, peer := range n.Peers {
	//     go n.sendDataToPeer(peer, data)
	// }
	// logger.InfoLogger.Println("Transaction broadcasted:", tx.ID)
}

func (n *Node) HandleIncomingTransaction(tx *blockchain.Transaction) error {
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

	// // Validate and add to mempool
	// if !tx.Verify() {
	//     logger.ErrorLogger.Println("Invalid transaction signature:", tx.ID)
	//     return ErrInvalidSignature
	// }

	// // Additional validation
	// if err := n.validateTransaction(tx); err != nil {
	//     logger.ErrorLogger.Println("Transaction validation failed:", err)
	//     return err
	// }

	// n.Mempool.AddTransaction(tx)
	// logger.InfoLogger.Println("Transaction added to mempool:", tx.ID)
	return nil
}

func (n *Node) HandleIncomingBlock(block *blockchain.Block) error {
	// Handle incoming block
	// Verify the block coming, verifying the transactions in it and if its the succeeding block
	// Based on the validation, carrry out the next operation -
	// If it is validated - add it to the ledger and make sure the mining is done for the next block
	// If not, then initiate the getBlocks protocol to figure out the missing blocks and act accordingly or weather to drop this block
	// Add additional methods or files as needed maintaining separation of concerns

	// n.Mutex.Lock()
	// defer n.Mutex.Unlock()

	// // Validate block
	// if err := n.Blockchain.AddBlock(block); err != nil {
	//     logger.ErrorLogger.Println("Failed to add incoming block:", err)
	//     // Initiate GetBlocks protocol if necessary
	//     return err
	// }

	// logger.InfoLogger.Println("Incoming block added to blockchain:", block.Header.BlockHash)
	return nil
}
