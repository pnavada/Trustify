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
	return nil, nil
}

func (m *Miner) ProofOfWork(b *Block) {
	// Perform POW to find valid nonce
}
