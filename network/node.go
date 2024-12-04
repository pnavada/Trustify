package network

import (
	"trustify/blockchain"
	"trustify/config"
)

type Node struct {
	Config     *config.Config
	Wallet     *blockchain.Wallet
	Blockchain *blockchain.Blockchain
	Mempool    *blockchain.Mempool
	UTXOSet    *blockchain.UTXOSet
	Miner      *blockchain.Miner
	peers      []string
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
	return nil
}

func (n *Node) Start() {
	// Start node operations: networking, transaction processing, mining - concurrent
	// A node should start listening for incoming transactions and blocks on a specified port
	// The node should create an outgoing connection to broadcast data over the network
	// The nodes should start mining to add new blocks to blockchain
	// Add additional methods or files as needed maintaining separation of concerns
}

func (n *Node) BroadcastTransaction(tx blockchain.UTXOTransaction) {
	// Broadcast transaction to the network
	// Broadcast the transaction data over the network to all the peers
	// do not use peer to peer multicasting instead use broadcasting
}

func (n *Node) BroadcastBlock(block blockchain.Block) {
	// Broadcast block to the network
	// Broadcast the block data over the network to all peers
	// do not use peer to peer multicasting instead use broadcasting
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

	return nil
}

func (n *Node) HandleIncomingBlock(block blockchain.Block) error {
	// Handle incoming block
	// Verify the block coming, verifying the transactions in it and if its the succeeding block
	// Based on the validation, carrry out the next operation -
	// If it is validated - add it to the ledger and make sure the mining is done for the next block
	// If not, then initiate the getBlocks protocol to figure out the missing blocks and act accordingly or weather to drop this block
	// Add additional methods or files as needed maintaining separation of concerns

	return nil
}
