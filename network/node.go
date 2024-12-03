package network

import (
	"trustify/blockchain"
	"trustify/config"
)

type Node struct {
	Config     *config.Config
	Wallet     *blockchain.Wallet `yaml:"wallet"`
	Blockchain *blockchain.Blockchain
	Mempool    *blockchain.Mempool
	// Miner      *mining.Miner
	// Network    *network.Network
}

func NewNode(cfg *config.Config) *Node {
	// Initialize node with wallet, blockchain, mempool, miner, and network
	return nil
}

func (n *Node) Start() {
	// Start node operations: networking, transaction processing, mining - concurrent
}
