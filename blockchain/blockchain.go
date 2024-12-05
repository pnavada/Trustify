package blockchain

import (
	"bytes"
	"trustify/config"
	"trustify/logger"
)

type Blockchain struct {
	Ledger            []*Block
	MiningReward      int
	ReviewReward      int
	ConfirmationDepth int
	TargetHash        []byte
	UTXOSet           *UTXOSet
}

// We are getting the geneisis block from config file and through an ConfigGenesisBlock object.
// So here leverage that and create an actual blockchain block and include it for further processing
// If the genesis block is invalid or initialization fails, return an appropriate error.
// Create a new Blockchain instance and initialize its fields, including the Ledger with the given genesisBlock.
// Add the provided genesisBlock to the Ledger as the first block in the chain.
// Initialize the blockchain’s settings using input parameters (i.e., MiningReward, ReviewReward, and ConfirmationDepth).
// Initialize any auxiliary structures required for managing transactions, such as UTXO sets or review tracking.
// Return the newly created Blockchain instance ready for use.
func NewBlockchain(genesisBlock *config.ConfigGenesisBlock, blockchainSettings *config.ConfigBlockchainSettings, utxoSet *UTXOSet, targetHash []byte) (*Blockchain, error) {
	// Convert ConfigGenesisBlock to Block
	block, err := convertConfigGenesisBlockToBlock(genesisBlock)
	if err != nil {
		logger.ErrorLogger.Println("Failed to convert genesis block:", err)
		return nil, err
	}

	logger.InfoLogger.Printf("Genesis Block: %+v\n", block)

	bc := &Blockchain{
		Ledger:            []*Block{block},
		MiningReward:      blockchainSettings.MiningReward,
		ReviewReward:      blockchainSettings.ReviewReward,
		ConfirmationDepth: blockchainSettings.BlockConfirmationDepth,
		UTXOSet:           utxoSet,
		TargetHash:        targetHash,
	}

	logger.InfoLogger.Printf("Blockchain initialized with genesis block:  %+v\n", bc)
	return bc, nil
}

func (bc *Blockchain) AddBlock(block *Block) error {

	err := bc.ValidateBlock(block)
	if err != nil {
		logger.ErrorLogger.Println("Invalid block:", err)
		return err
	}

	// Validate previous hash
	lastBlock := bc.LatestBlock()
	
	if !bytes.Equal(block.Header.PreviousHash, lastBlock.Header.BlockHash) {
		logger.ErrorLogger.Println("Block's previous hash does not match the latest block's hash")
		// TODO:  Start getblocks protocol
	} else {
		// Add block to blockchain
		bc.Ledger = append(bc.Ledger, block)
		logger.InfoLogger.Println("Block added to blockchain:", block.Header.BlockHash)
		// Commit block if it reaches the confirmation depth
		bc.CommitBlock()
	}

	return nil

}

func (bc *Blockchain) ValidateProofOfWork(block *Block, target []byte) bool {

	hash := bc.ComputeHash(block)
	return bytes.Compare(hash, target) == -1

}

func (bc *Blockchain) ValidateMerkeRoot(block *Block) bool {

	merkleTree, err := BuildTree(block.Transactions)
	if err != nil {
		logger.ErrorLogger.Println("Failed to build Merkle tree:", err)
		return false
	}

	return bytes.Equal(merkleTree.Root.Hash, block.Header.MerkleRoot)

}

func (bc *Blockchain) ValidateTimeStamp(block *Block) bool {

	lastBlock := bc.LatestBlock()
	return block.Header.Timestamp > lastBlock.Header.Timestamp

}

func (bc *Blockchain) ComputeHash(block *Block) []byte {

	serializedBlockHeader := Serialize(block.Header)
	hash := HashObject(serializedBlockHeader)
	return hash

}

func (bc *Blockchain) ValidateBlock(block *Block) error {

	// Verify proof of work
	if !bc.ValidateProofOfWork(block, bc.TargetHash) {
		logger.ErrorLogger.Println("Invalid proof of work")
		return ErrInvalidProofOfWork
	}

	// Verify merkle root
	if !bc.ValidateMerkeRoot(block) {
		logger.ErrorLogger.Println("Invalid Merkle root")
		return ErrInvalidMerkleRoot
	}

	// Verify timestamp
	if !bc.ValidateTimeStamp(block) {
		logger.ErrorLogger.Println("Invalid timestamp")
		return ErrInvalidTimestamp
	}

	// Validate transactions
	for _, tx := range block.Transactions {
		err := bc.ValidateTransaction(tx, bc.UTXOSet)
		if err != nil {
			logger.ErrorLogger.Println("Invalid transaction:", err)
			return err
		}
	}

	return nil

}

func (bc *Blockchain) GetBlockByHash(hash []byte) (*Block, error) {
	// Ensure the provided hash is in the correct format and non-empty.
	// Iterate through the blockchain’s list of blocks to locate the block that matches the provided hash.
	// If a block with the matching hash is found, return it.
	// If no block is found with the given hash, return a meaningful error indicating that the block does not exist.
	// Ensure that the retrieved block is valid within the context of the current chain state (e.g., hasn’t been replaced by a fork).
	for _, block := range bc.Ledger {
		if bytes.Equal(block.Header.BlockHash, hash) {
			logger.InfoLogger.Println("Block found for hash:", hash)
			return block, nil
		}
	}
	logger.ErrorLogger.Println("Block not found for hash:", hash)
	return nil, ErrBlockNotFound

}

func (bc *Blockchain) LatestBlock() *Block {

	if len(bc.Ledger) == 0 {
		logger.ErrorLogger.Println("Blockchain is empty")
		return nil
	}
	return bc.Ledger[len(bc.Ledger)-1]

}

// Add a method to identify committed blocks and transactions based on the confirmation depth available from the configuration
// This method should check for committed blocks and transactions
// Update the UTXO set with committed transactions
func (bc *Blockchain) CommitBlock() {
	// Implement block and transaction confirmation logic
	numBlocks := len(bc.Ledger)
	blockToCommitIndex := numBlocks - bc.ConfirmationDepth - 1
	blockToCommit := bc.Ledger[blockToCommitIndex]
	for _, tx := range blockToCommit.Transactions {
		for _, utxo := range tx.Inputs {
			_, hasUTXO := bc.UTXOSet.Get(utxo.ID)
			if hasUTXO {
				bc.UTXOSet.Remove(utxo.ID)
			} else {
				// TODO: Handle error
			}
		}

		for _, utxo := range tx.Outputs {
			_, hasUTXO := bc.UTXOSet.Get(utxo.ID)
			if hasUTXO {
				// Handle error
			} else {
				bc.UTXOSet.Add(utxo)
			}
		}
	}
}

func (bc *Blockchain) GetUTXOSet() *UTXOSet {
	return bc.UTXOSet
}

// validateTransaction validates a single transaction for correctness.
func (bc *Blockchain) ValidateTransaction(tx *Transaction, utxoSet *UTXOSet) error {
	// Implement validation logic for transactions
	// Check UTXOSet for inputs
	// Verify signatures, double-spending, etc.

	// Separate handling for Purchase and Review transactions
	switch data := tx.Data.(type) {
	case *PurchaseTransactionData:
		return bc.validatePurchaseTransaction(tx, data, utxoSet)
	case *ReviewTransactionData:
		return bc.validateReviewTransaction(data)
	default:
		logger.ErrorLogger.Println("Invalid transaction type detected")
		return ErrInvalidTransactionType
	}
}

// validatePurchaseTransaction validates purchase transactions.
func (bc *Blockchain) validatePurchaseTransaction(tx *Transaction, data *PurchaseTransactionData, utxoSet *UTXOSet) error {
	inputSum := 0
	outputSum := 0
	usedUTXOs := make(map[string]bool)

	// Validate Inputs
	for _, input := range tx.Inputs {
		utxo, exists := utxoSet.Get(input.ID)
		if !exists {
			logger.ErrorLogger.Printf("Input UTXO not found: %v\n", input.ID)
			return ErrUTXONotFound
		}

		// Check for double-spending within the same transaction
		utxoKey := input.ID.String()
		if usedUTXOs[utxoKey] {
			logger.ErrorLogger.Println("Double-spending detected within the transaction")
			return ErrDoubleSpending
		}
		usedUTXOs[utxoKey] = true

		// Accumulate the input sum
		inputSum += utxo.Amount
	}

	// Validate Outputs
	for _, output := range tx.Outputs {
		outputSum += output.Amount
	}

	// Check if input sum covers output sum and fee
	if inputSum < outputSum {
		logger.ErrorLogger.Printf("Input sum (%d) is less than output sum (%d)\n", inputSum, outputSum)
		return ErrInsufficientFunds
	}

	logger.InfoLogger.Printf("Purchase transaction validated successfully for Buyer: %x, Product: %s\n", data.BuyerAddress, data.ProductID)
	return nil
}

// validateReviewTransaction validates review transactions.
func (bc *Blockchain) validateReviewTransaction(data *ReviewTransactionData) error {
	reviewer := data.ReviewerAddress
	product := data.ProductID

	// Check if the product was purchased by the reviewer
	if !bc.hasPurchasedProduct(reviewer, product) {
		logger.ErrorLogger.Printf("Reviewer %x has not purchased product %s\n", reviewer, product)
		return ErrProductNotPurchased
	}

	// Check if a review already exists for this product by the reviewer
	if bc.hasDuplicateReview(reviewer, product) {
		logger.ErrorLogger.Printf("Duplicate review detected by %x for product %s\n", reviewer, product)
		return ErrDuplicateReview
	}

	logger.InfoLogger.Printf("Review transaction validated successfully for Reviewer: %x, Product: %s\n", reviewer, product)
	return nil
}

// hasPurchasedProduct checks the blockchain ledger for a purchase transaction by the reviewer for the product.
func (bc *Blockchain) hasPurchasedProduct(reviewer []byte, product string) bool {
	for _, block := range bc.Ledger {
		for _, tx := range block.Transactions {
			if purchaseData, ok := tx.Data.(*PurchaseTransactionData); ok {
				if bytes.Equal(purchaseData.BuyerAddress, reviewer) && purchaseData.ProductID == product {
					return true
				}
			}
		}
	}
	return false
}

// hasDuplicateReview checks the blockchain ledger for duplicate reviews.
func (bc *Blockchain) hasDuplicateReview(reviewer []byte, product string) bool {
	for _, block := range bc.Ledger {
		for _, tx := range block.Transactions {
			if reviewData, ok := tx.Data.(*ReviewTransactionData); ok {
				if bytes.Equal(reviewData.ReviewerAddress, reviewer) && reviewData.ProductID == product {
					return true
				}
			}
		}
	}
	return false
}

func (bc *Blockchain) updateUTXOSet(b *Block) {
	// Remove spent UTXOs and add new UTXOs from the block's transactions
}
