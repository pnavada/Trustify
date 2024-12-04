package blockchain

type Miner struct {
	Blockchain *Blockchain
	Mempool    *Mempool
}

func NewMiner(bc *Blockchain, mp *Mempool) *Miner {
	return &Miner{Blockchain: bc, Mempool: mp}
}

func (m *Miner) MineBlock() (*Block, error) {
	// Collect transactions, create block, perform proof of work
	// This method should be called when there are enough transactions in the mempool
	// form a block.
	// The required block size is defined in the configuation object
	// Then pick the required number of transactions from the mempool
	// Add coinbase transactions to reward the miners and reviewers (if applicable).
	// Use your best decision to position this logic
	// Create the block using the NewBlock method
	// The target hash is part of the configuration. Dont receive the confuguration object as a parameter
	// USe intelligent ways to pass this data between objects
	// Define the parameters where it makes the most sense to be a part of
	// Perform proof of work to find the valid nonce for the block
	// Add the block to the ledger and broadcast over the network
	return nil, nil
}

func (m *Miner) ProofOfWork(b *Block) {
	// Perform POW to find valid nonce
	// Here  the block is the block for which the nonce is to be found
	// The nonce is intiialized to zero and incremented until the hash of the block is less than the target hash
	// Basic idea is to find a nonce such that the hash of the block is less than the target hash
}
