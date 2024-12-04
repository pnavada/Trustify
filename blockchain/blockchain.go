package blockchain

import (
	"trustify/config"
)

type Blockchain struct {
	Ledger            []*Block
	MiningReward      int
	ReviewReward      int
	ConfirmationDepth int
}

// Add a method to create the genesis block from the genesis config object

func NewBlockchain(genesisBlock *config.ConfigGenesisBlock) *Blockchain {
	// We are getting the geneisis block from config file and through an ConfigGenesisBlock object.
	// So here leverage that and create an actual blockchain block and include it for further processing
	// If the genesis block is invalid or initialization fails, return an appropriate error.
	// Create a new Blockchain instance and initialize its fields, including the Ledger with the given genesisBlock.
	// Add the provided genesisBlock to the Ledger as the first block in the chain.
	// Initialize the blockchain’s settings using input parameters (i.e., MiningReward, ReviewReward, and ConfirmationDepth).
	// Initialize any auxiliary structures required for managing transactions, such as UTXO sets or review tracking.
	// Return the newly created Blockchain instance ready for use.
	return nil
}

func (bc *Blockchain) AddBlock(b *Block) error {
	// Check the structure of the block, ensuring it contains all the required fields to create a block.
	// Perform all validations necessary
	// Validate the merkle root of the block’s transactions against the transactions themselves.
	// Validate the block’s Proof of Work (PoW) against the configured target hash.
	// Validate the block’s timestamp to ensure it is within a reasonable range.
	// Verify that the block references the correct hash of the previous block in the chain.
	// Iterate through all transactions in the block:
	// Perform the required checks for all transactions: verifying digital signatures, verifying input and output transactions, etc
	// For purchase transactions, validate the inputs against the unspent transaction outputs (UTXO) set.
	// Ensure there is no double-spending in the block.
	// For review transactions, ensure the reviwer has purchase the product
	// For review transactions, confirm that the reviewer has not already submitted a review for it.
	// Append the validated block to the chain if all checks pass.
	// Return meaningful error messages if the block fails any validation step.
	// Make sure the addition of the block is an atomic operation—either fully added or not at all, to maintain blockchain integrity.
	return nil
}

func (bc *Blockchain) GetBlockByHash(hash []byte) (*Block, error) {
	// Ensure the provided hash is in the correct format and non-empty.
	// Iterate through the blockchain’s list of blocks to locate the block that matches the provided hash.
	// If a block with the matching hash is found, return it.
	// If no block is found with the given hash, return a meaningful error indicating that the block does not exist.
	// Ensure that the retrieved block is valid within the context of the current chain state (e.g., hasn’t been replaced by a fork).
	return nil, nil
}

func (bc *Blockchain) LatestBlock() *Block {
	// Retrieve the last block added to the blockchain, which represents the current state of the ledger.
	// If the blockchain is empty (e.g., no blocks have been added), return nil.
	return nil
}

// Add a method to identify committed blocks and transactions based on the confirmation depth available from the configuration
// This method should check for committed blocks and transactions
// Update the UTXO set with committed transactions
