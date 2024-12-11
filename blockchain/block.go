package blockchain

import (
	"time"
	"trustify/logger"
)

// BlockHeader represents metadata of a block.
type BlockHeader struct {
	BlockHash    []byte // Hash of the block
	PreviousHash []byte // Hash of the previous block
	MerkleRoot   []byte // Root of the Merkle tree for transactions
	Timestamp    int64  // Unix timestamp when the block is created
	TargetHash   []byte // Target difficulty for proof of work
	Nonce        int64  // Nonce used for proof of work
}

// Block represents a complete blockchain block.
type Block struct {
	Header           *BlockHeader   // Metadata of the block
	TransactionCount int            // Total number of transactions in the block
	Transactions     []*Transaction // List of transactions in the block
}

// NewBlock creates and initializes a new block.
func NewBlock(transactions []*Transaction, previousHash []byte, targetHash []byte) (*Block, error) {
	// Validate inputs
	if len(transactions) == 0 {
		logger.ErrorLogger.Println("Cannot create block: No transactions provided")
		return nil, ErrEmptyTransactions
	}
	if len(previousHash) == 0 {
		logger.ErrorLogger.Println("Cannot create block: Invalid previous hash")
		return nil, ErrInvalidPreviousHash
	}
	if len(targetHash) == 0 {
		logger.ErrorLogger.Println("Cannot create block: Invalid target hash")
		return nil, ErrInvalidTargetHash
	}

	// Compute the Merkle root from the transaction list
	merkleRoot, err := BuildTree(transactions)
	if err != nil {
		logger.ErrorLogger.Println("Failed to compute Merkle root:", err)
		return nil, err
	}

	// Create the block header
	if merkleRoot == nil {
		return nil, ErrInvalidMerkleRoot
	}
	header := BlockHeader{
		PreviousHash: previousHash,
		MerkleRoot:   merkleRoot.GetRoot(),
		Timestamp:    time.Now().UnixNano(), // Using nanosecond precision
		TargetHash:   targetHash,
		Nonce:        0,        // Initial nonce (to be updated during mining)
		BlockHash:    []byte{}, // Placeholder for the block hash
	}

	// Create the block
	block := &Block{
		Header:           &header,
		TransactionCount: len(transactions),
		Transactions:     transactions,
	}

	// Compute the block hash and update the block header
	// blockHash :=
	// block.Header.BlockHash = blockHash

	logger.InfoLogger.Printf("New block created with hash: %x", computeBlockHash(block))
	return block, nil
}

// computeBlockHash calculates the hash of a block using its header and transaction data.
func computeBlockHash(block *Block) []byte {
	serializedBlock := SerializeBlockHeader(block.Header) // Serialize block data
	return HashObject(serializedBlock)                    // Compute the hash
}

// GetTransactionFee calculates the total transaction fees in a block.
func (b *Block) GetTransactionFee() int64 {
	var totalFees int64
	for _, transaction := range b.Transactions {
		totalFees += int64(transaction.GetTransactionFee())
	}
	return totalFees
}

// Print the complete contents of the block - similar to JAVA's toString() method in smart and consise way and also write all the transactions of the block
func (b *Block) PrintToString() {
	// Print returns a string representation of the block
	logger.InfoLogger.Printf("Block Hash: %x", b.Header.BlockHash)
	logger.InfoLogger.Printf("Previous Hash: %x", b.Header.PreviousHash)
	logger.InfoLogger.Printf("Merkle Root: %x", b.Header.MerkleRoot)
	logger.InfoLogger.Printf("Timestamp: %v", time.Unix(0, b.Header.Timestamp))
	logger.InfoLogger.Printf("Target Hash: %x", b.Header.TargetHash)
	logger.InfoLogger.Printf("Nonce: %d", b.Header.Nonce)
	logger.InfoLogger.Printf("Transaction Count: %d", b.TransactionCount)

	for i, tx := range b.Transactions {
		logger.InfoLogger.Printf("Transaction %d:", i+1)
		logger.InfoLogger.Printf("  Sender: %s", tx.Data)
	}
}
