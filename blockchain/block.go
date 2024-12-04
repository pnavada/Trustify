package blockchain

type BlockHeader struct {
	BlockHash    []byte
	PreviousHash []byte
	MerkleRoot   []byte
	Timestamp    int64
	TargetHash   []byte
	Nonce        int64
}

type Block struct {
	Header           BlockHeader
	TransactionCount int
	Transactions     []UTXOTransaction
}

func NewBlock(transactions []*UTXOTransaction, previousHash string, targetHash string) *Block {
	// Create a new block with transactions
	return nil
}

// TODO: Verify if this method should be moved to mining.go
func (b *Block) ComputeHash() []byte {
	// Compute the block's hash
	return nil
}

// TODO: Add a method to identify committed blocks and transactions
