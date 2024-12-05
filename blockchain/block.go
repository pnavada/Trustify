package blockchain

// import (
//     "trustify/logger"
// )

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
	Transactions     []*UTXOTransaction
}

func NewBlock(transactions []*UTXOTransaction, previousHash []byte, targetHash []byte) (*Block, error) {
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
	// Ensure that the transactions list is not empty.
	// if len(transactions) == 0 {
    //     logger.ErrorLogger.Println("Attempted to create a block with no transactions")
    //     return nil, ErrEmptyTransactions
    // }

    // if len(previousHash) == 0 {
    //     logger.ErrorLogger.Println("Invalid previous hash provided")
    //     return nil, ErrInvalidPreviousHash
    // }

    // if len(targetHash) == 0 {
    //     logger.ErrorLogger.Println("Invalid target hash provided")
    //     return nil, ErrInvalidTargetHash
    // }

    // merkleRoot, err := crypto.ComputeMerkleRoot(transactions)
    // if err != nil {
    //     logger.ErrorLogger.Println("Failed to compute Merkle root:", err)
    //     return nil, err
    // }

    // header := BlockHeader{
    //     PreviousHash: previousHash,
    //     MerkleRoot:   merkleRoot,
    //     Timestamp:    time.Now().Unix(),
    //     TargetHash:   targetHash,
    //     Nonce:        0, // Will be updated during mining
    // }

    // block := &Block{
    //     Header:           header,
    //     TransactionCount: len(transactions),
    //     Transactions:     transactions,
    // }

    // // Compute the block hash (without Nonce for now)
    // blockHash := block.ComputeHash()
    // block.Header.BlockHash = blockHash

    // logger.InfoLogger.Println("New block created with hash:", blockHash)
    // return block, nil

    return nil, nil
}

// TODO: Verify if this method should be moved to mining.go
func (b *Block) ComputeHash() []byte {
	// Compute the block's hash
	// var buffer bytes.Buffer

    // buffer.Write(b.Header.PreviousHash)
    // buffer.Write(b.Header.MerkleRoot)
    // timestampBytes := make([]byte, 8)
    // binary.LittleEndian.PutUint64(timestampBytes, uint64(b.Header.Timestamp))
    // buffer.Write(timestampBytes)
    // buffer.Write(b.Header.TargetHash)
    // nonceBytes := make([]byte, 8)
    // binary.LittleEndian.PutUint64(nonceBytes, uint64(b.Header.Nonce))
    // buffer.Write(nonceBytes)

    // hash := sha256.Sum256(buffer.Bytes())
    // logger.InfoLogger.Println("Computed hash for block:", hash[:])
    // return hash[:]
    return nil
}