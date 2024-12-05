package blockchain

import (
	"encoding/hex"
	"errors"
	"trustify/config"
	"trustify/logger"
)

func convertConfigGenesisBlockToBlock(genesisConfig *config.ConfigGenesisBlock) (*Block, error) {
	if genesisConfig == nil {
		logger.ErrorLogger.Println("Genesis config is nil")
		return nil, errors.New("genesis config is nil")
	}

	// Parse and validate the target hash
	targetHash, err := hex.DecodeString(genesisConfig.TargetHash)
	if err != nil || len(targetHash) == 0 {
		logger.ErrorLogger.Printf("Invalid target hash: %s\n", genesisConfig.TargetHash)
		return nil, errors.New("invalid target hash in genesis block")
	}

	// Parse and validate the block hash
	blockHash, err := hex.DecodeString(genesisConfig.BlockHash)
	if err != nil || len(blockHash) == 0 {
		logger.ErrorLogger.Printf("Invalid block hash: %s\n", genesisConfig.BlockHash)
		return nil, errors.New("invalid block hash in genesis block")
	}

	// Parse and validate the previous hash
	previousHash, err := hex.DecodeString(genesisConfig.PreviousHash)
	if err != nil || len(previousHash) == 0 {
		logger.ErrorLogger.Printf("Invalid previous hash: %s\n", genesisConfig.PreviousHash)
		return nil, errors.New("invalid previous hash in genesis block")
	}

	// Parse and validate the Merkle root
	merkleRoot := []byte(genesisConfig.MerkleRoot)
	// if err != nil || len(merkleRoot) == 0 {
	//     logger.ErrorLogger.Printf("Invalid Merkle root: %s\n", genesisConfig.MerkleRoot)
	//     return nil, errors.New("invalid Merkle root in genesis block")
	// }

	// Convert the transactions from ConfigUTXOTransaction to UTXOTransaction
	var transactions []*UTXOTransaction
	for _, tx := range genesisConfig.Transactions.Outputs {
		transaction := &UTXOTransaction{
			ID: UTXOTransactionID{
				TxHash:  blockHash,
				TxIndex: len(transactions),
			},
			Address: []byte(tx.Address),
			Amount:  tx.Amount,
		}
		transactions = append(transactions, transaction)
	}

	var blockTransactions []*Transaction
	transaction := &Transaction{
		Outputs: transactions,
		Inputs:  nil,
		Data:    genesisConfig.Transactions.Data,
		ID:      genesisConfig.Transactions.ID,
	}

	blockTransactions = append(blockTransactions, transaction)

	if len(transactions) != genesisConfig.TransactionCount {
		logger.ErrorLogger.Println("Mismatch in transaction count in genesis block")
		return nil, errors.New("transaction count mismatch in genesis block")
	}

	// Create the BlockHeader
	header := BlockHeader{
		BlockHash:    blockHash,
		PreviousHash: previousHash,
		MerkleRoot:   merkleRoot,
		Timestamp:    int64(genesisConfig.Timestamp),
		TargetHash:   targetHash,
		Nonce:        int64(genesisConfig.Nonce),
	}

	// Create the Block
	block := &Block{
		Header:           header,
		TransactionCount: len(transactions),
		Transactions:     blockTransactions,
	}

	logger.InfoLogger.Printf("Genesis block converted successfully with hash: %x\n", blockHash)
	return block, nil
}
