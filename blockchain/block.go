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

func NewBlock(transactions []*UTXOTransaction, previousHash []byte, targetHash []byte) *Block {
	// Ensure that the transactions list is not empty.
	// Check if previousHash and targetHash are valid
	// Use the transactions list to compute the Merkle Root: Hash each transaction and pair the hashes and iteratively hash them to compute the root.
	// The computed Merkle Root will represent the integrity of the transaction set.
	// Populate the BlockHeader structure with:
	// PreviousHash: Hash of the last block.
	// MerkleRoot: Root of the transaction tree.
	// Timestamp: Current timestamp.
	// TargetHash: The difficulty target for Proof of Work.
	// Leave Nonce empty; it will be updated during mining.
	// Assign the list of UTXOTransaction objects to the Transactions field.
	// Set the TransactionCount field to the length of the transactions list.
	// Package the BlockHeader and transaction data into a Block structure.
	// Return the new Block object for further processing.
	return nil
}

// TODO: Verify if this method should be moved to mining.go
func (b *Block) ComputeHash() []byte {
	// Compute the block's hash
	return nil
}
