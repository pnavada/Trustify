package network

import (
	"bytes"
	"os"
	"trustify/blockchain"
	"trustify/config"
	"trustify/logger"

	"context"

	"log"

	"github.com/libp2p/go-libp2p"

)

type Node struct {
	Config     *config.Config
	Wallet     *blockchain.Wallet
	Blockchain *blockchain.Blockchain
	Mempool    *blockchain.Mempool
	UTXOSet    *blockchain.UTXOSet
	Miner      *blockchain.Miner
	Peers      []string
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
		utxoSet.Add(tx) // TODO: create a copy of tx
		if bytes.Equal(tx.Address, wallet.BitcoinAddress) {
			wallet.UTXOs = append(wallet.UTXOs, tx)
		}
	}

	miner := blockchain.NewMiner(chain, mempool)

	// Initialize peers
	var peers []string
	for nodeName := range cfg.Nodes {
		if nodeName != me {
			peers = append(peers, nodeName)
		}
	}

	node := &Node{
		Config:     cfg,
		Wallet:     wallet,
		Blockchain: chain,
		Mempool:    mempool,
		UTXOSet:    utxoSet,
		Miner:      miner,
		Peers:      peers,
	}

	logger.InfoLogger.Printf("Node initialized: %+v\n", node)

	// logger.InfoLogger.Println("Node initialized with address:", wallet.BitcoinAddress)
	return node
}

func (n *Node) Start() {
	// Start node operations: networking, transaction processing, mining - concurrent
	// A node should start listening for incoming transactions and blocks on a specified port
	// The node should create an outgoing connection to broadcast data over the network
	// The nodes should start mining to add new blocks to blockchain
	// Add additional methods or files as needed maintaining separation of concerns

	// Start networking, transaction processing, mining
	go n.listenForIncomingData()
	go n.mineBlocks()
	logger.InfoLogger.Println("Node started operations")
}

func (n *Node) BroadcastTransaction(tx blockchain.UTXOTransaction) {
	// Broadcast transaction to the network
	// Broadcast the transaction data over the network to all the peers
	// do not use peer to peer multicasting instead use broadcasting

	// Serialize and broadcast the transaction to peers
    data := utils.SerializeTransaction(tx)
    for _, peer := range n.Peers {
        go n.sendDataToPeer(peer, data)
    }
    logger.InfoLogger.Println("Transaction broadcasted:", tx.ID)
}

func (n *Node) BroadcastBlock(block blockchain.Block) {
	// Broadcast block to the network
	// Broadcast the block data over the network to all peers
	// do not use peer to peer multicasting instead use broadcasting

	// Serialize and broadcast the block to peers
    data := utils.SerializeBlock(block)
    for _, peer := range n.Peers {
        go n.sendDataToPeer(peer, data)
    }
    logger.InfoLogger.Println("Block broadcasted:", block.Header.BlockHash)
}

func (n *Node) HandleIncomingTransaction(tx blockchain.UTXOTransaction) error {
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

	// Validate and add to mempool
    if !tx.Verify() {
        logger.ErrorLogger.Println("Invalid transaction signature:", tx.ID)
        return ErrInvalidSignature
    }

    // Additional validation
    if err := n.validateTransaction(tx); err != nil {
        logger.ErrorLogger.Println("Transaction validation failed:", err)
        return err
    }

    n.Mempool.AddTransaction(tx)
    logger.InfoLogger.Println("Transaction added to mempool:", tx.ID)
    return nil
}

func (n *Node) HandleIncomingBlock(block blockchain.Block) error {
	// Handle incoming block
	// Verify the block coming, verifying the transactions in it and if its the succeeding block
	// Based on the validation, carrry out the next operation -
	// If it is validated - add it to the ledger and make sure the mining is done for the next block
	// If not, then initiate the getBlocks protocol to figure out the missing blocks and act accordingly or weather to drop this block
	// Add additional methods or files as needed maintaining separation of concerns

	n.Mutex.Lock()
    defer n.Mutex.Unlock()

    // Validate block
    if err := n.Blockchain.AddBlock(block); err != nil {
        logger.ErrorLogger.Println("Failed to add incoming block:", err)
        // Initiate GetBlocks protocol if necessary
        return err
    }

    logger.InfoLogger.Println("Incoming block added to blockchain:", block.Header.BlockHash)
    return nil
}

func (n *Node) mineBlocks() {
    // Continuously attempt to mine new blocks
    for {
        block, err := n.Miner.MineBlock()
        if err != nil {
            logger.ErrorLogger.Println("Mining failed:", err)
        } else if block != nil {
            n.BroadcastBlock(block)
        }
        // Wait or check for new transactions before attempting next block
    }
}