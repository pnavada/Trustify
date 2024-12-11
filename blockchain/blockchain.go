package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"sort"
	"trustify/config"
	"trustify/logger"
)

// Blockchain represents the blockchain ledger and associated settings.
type Blockchain struct {
	Ledger            []*Block
	MiningReward      int
	ReviewReward      int
	ConfirmationDepth int
	UTXOSet           *UTXOSet
	TargetHash        []byte
	GetBlocksProtocol *GetBlocksProtocol
	Mempool           *Mempool
}

// NewBlockchain initializes a new blockchain with the provided genesis block and settings.
func NewBlockchain(genesisBlock *config.ConfigGenesisBlock, blockchainSettings *config.ConfigBlockchainSettings, utxoSet *UTXOSet, getBlocksProtocol *GetBlocksProtocol, mempool *Mempool) (*Blockchain, error) {
	block, err := convertConfigGenesisBlockToBlock(genesisBlock)
	if err != nil {
		logger.ErrorLogger.Println("Failed to convert genesis block:", err)
		return nil, err
	}

	logger.InfoLogger.Printf("Genesis Block Initialized: %+v\n", block)
	targetHash, _ := hex.DecodeString(blockchainSettings.TargetHash)

	bc := &Blockchain{
		Ledger:            []*Block{block},
		MiningReward:      blockchainSettings.MiningReward,
		ReviewReward:      blockchainSettings.ReviewReward,
		ConfirmationDepth: blockchainSettings.BlockConfirmationDepth,
		UTXOSet:           utxoSet,
		TargetHash:        targetHash,
		GetBlocksProtocol: getBlocksProtocol,
		Mempool:           mempool,
	}

	logger.InfoLogger.Println("Blockchain successfully initialized with genesis block.")
	return bc, nil
}

// AddBlock adds a block to the blockchain after validating it.
func (bc *Blockchain) AddBlock(block *Block) error {
	if err := bc.ValidateBlock(block); err != nil {
		logger.ErrorLogger.Println("Block validation failed:", err)
		return err
	}

	lastBlock := bc.LatestBlock()
	lastBlockHash := bc.ComputeHash(lastBlock)

	if !bytes.Equal(block.Header.PreviousHash, lastBlockHash) {
		err := bc.resolveFork(lastBlock.Header.BlockHash)
		if err != nil {
			logger.ErrorLogger.Println("Failed to resolve fork:", err)
			return err
		}
	}

	bc.Ledger = append(bc.Ledger, block)
	logger.InfoLogger.Printf("Block added to blockchain: %x\n", block.Header.BlockHash)
	bc.CommitBlock()
	return nil
}

// resolveFork resolves forks in the blockchain by choosing the best chain.
func (bc *Blockchain) resolveFork(startHash []byte) error {
	for {
		// Start the getblocks protocol to fetch blocks from peers
		err := bc.GetBlocksProtocol.GetBlocks(startHash)
		if err != nil {
			logger.ErrorLogger.Println("Failed to start getblocks:", err)
			return err
		}

		// Wait for responses from BlocksChannel
		responses := bc.getResponses()
		if responses == nil {
			logger.ErrorLogger.Println("No responses received from peers")
			return errors.New("no responses received from peers")
		}

		// Process responses to extract blocks
		allBlocks := bc.processResponses(responses)
		if len(allBlocks) == 0 {
			logger.ErrorLogger.Println("No valid blocks received from peers")
			startBlock, err := bc.GetBlockByHash(startHash)
			if err != nil {
				return err
			}
			startHash = startBlock.Header.PreviousHash
			continue
		}

		// Validate all blocks in each chain
		if err := bc.validateChains(allBlocks); err != nil {
			return err
		}

		// Choose the best chain based on length and transaction fees
		bestChain := bc.chooseBestChain(allBlocks)
		if bestChain == nil {
			logger.ErrorLogger.Println("No valid chain found")
			return errors.New("no valid chain found")
		}

		// Find the index of the first node of the best chain in the ledger
		firstNodeIndex, err := bc.GetBlockIndexByHash(bestChain[0].Header.BlockHash)
		if err != nil {
			logger.ErrorLogger.Println("Failed to get block index by hash:", err)
		}

		if firstNodeIndex != -1 && firstNodeIndex+len(bestChain) > len(bc.Ledger) {
			// Add unmatched blocks' transactions to the mempool
			bc.addUnmatchedBlocksToMempool(startHash)

			// Add the best chain to the ledger
			bc.addBestChainToLedger(bestChain[firstNodeIndex:])
		}

		logger.InfoLogger.Println("Best chain added to ledger")
		return nil
	}
}

// GetBlockIndexByHash retrieves the index of a block in the ledger by its hash.
func (bc *Blockchain) GetBlockIndexByHash(hash []byte) (int, error) {
	for i, block := range bc.Ledger {
		if bytes.Equal(block.Header.BlockHash, hash) {
			logger.InfoLogger.Printf("Block index found for hash: %x\n", hash)
			return i, nil
		}
	}
	logger.ErrorLogger.Printf("Block index not found for hash: %x\n", hash)
	return -1, ErrBlockNotFound
}

// validateChains validates all blocks in each chain.
func (bc *Blockchain) validateChains(chains [][]*Block) error {
	var validChains [][]*Block
	for _, chain := range chains {
		isValidChain := bc.isValidChain(chain)
		if isValidChain {
			validChains = append(validChains, chain)
		}
	}
	if len(validChains) == 0 {
		logger.ErrorLogger.Println("No valid chains found")
		return errors.New("no valid chains found")
	}
	logger.InfoLogger.Printf("Found %d valid chains\n", len(validChains))
	return nil
}

func (bc *Blockchain) isValidChain(chain []*Block) bool {
	for i, blk := range chain {
		if err := bc.ValidateBlock(blk); err != nil {
			logger.ErrorLogger.Println("Block validation failed in chain:", err)
			return false
		}
		if i > 0 && !bytes.Equal(blk.Header.PreviousHash, chain[i-1].Header.BlockHash) {
			logger.ErrorLogger.Println("Previous hash does not match in chain")
			return false
		}
	}
	return true
}

// addBestChainToLedger adds the best chain to the ledger.
func (bc *Blockchain) addBestChainToLedger(bestChain []*Block) {
	bc.Ledger = append(bc.Ledger, bestChain...)
}

// addUnmatchedBlocksToMempool adds unmatched blocks' transactions to the mempool.
func (bc *Blockchain) addUnmatchedBlocksToMempool(startHash []byte) {

	for i := len(bc.Ledger) - 1; i >= 0; i-- {
		if bytes.Equal(bc.Ledger[i].Header.BlockHash, startHash) {
			bc.Ledger = bc.Ledger[:i]
			break
		}
		for _, tx := range bc.Ledger[i].Transactions {
			bc.Mempool.AddTransaction(tx)
		}
	}

}

// getResponses collects responses from the BlocksChannel.
func (bc *Blockchain) getResponses() []*GetBlocksResponse {
	var responses []*GetBlocksResponse
	for {
		select {
		case response := <-bc.GetBlocksProtocol.BlocksChannel:
			responses = append(responses, response)
		default:
			return responses
		}
	}
}

// chooseBestChain selects the best chain based on length and transaction fees.
func (bc *Blockchain) chooseBestChain(chains [][]*Block) []*Block {
	if len(chains) == 0 {
		return nil
	}

	sort.Slice(chains, func(i, j int) bool {
		if len(chains[i]) != len(chains[j]) {
			return len(chains[i]) > len(chains[j])
		}
		var feeI, feeJ int64
		for _, block := range chains[i] {
			feeI += block.GetTransactionFee()
		}
		for _, block := range chains[j] {
			feeJ += block.GetTransactionFee()
		}
		if feeI == feeJ {
			return chains[i][len(chains[i])-1].Header.Timestamp < chains[j][len(chains[j])-1].Header.Timestamp
		}
		return feeI > feeJ
	})

	return chains[0]
}

// ValidateBlock performs all validations on a block before adding it to the blockchain.
func (bc *Blockchain) ValidateBlock(block *Block) error {
	if !bc.ValidateProofOfWork(block, bc.TargetHash) {
		logger.ErrorLogger.Println("Invalid proof of work")
		return ErrInvalidProofOfWork
	} else {
		logger.InfoLogger.Println("Block's proof of work is valid.")
	}

	if !bc.ValidateMerkleRoot(block) {
		logger.ErrorLogger.Println("Invalid Merkle root")
		return ErrInvalidMerkleRoot
	} else {
		logger.InfoLogger.Println("Block's Merkle root is valid.")
	}

	if !bc.ValidateTimestamp(block) {
		logger.ErrorLogger.Println("Invalid timestamp")
		return ErrInvalidTimestamp
	} else {
		logger.InfoLogger.Println("Block's timestamp is valid.")
	}

	for _, tx := range block.Transactions {
		if err := bc.ValidateTransaction(tx, bc.UTXOSet); err != nil {
			logger.ErrorLogger.Printf("Invalid transaction %x: %v\n", tx.ID, err)
			return err
		}
	}

	logger.InfoLogger.Printf("Block %x passed all validations.\n", block.Header.BlockHash)
	return nil
}

// ValidateProofOfWork checks if the block's hash meets the target difficulty.
func (bc *Blockchain) ValidateProofOfWork(block *Block, target []byte) bool {
	hash := bc.ComputeHash(block)
	logger.InfoLogger.Printf("Block's nonce: %d\n", block.Header.Nonce)
	isValid := bytes.Compare(hash, target) < 0
	if !isValid {
		logger.InfoLogger.Printf("Proof of work failed for block %x\n", block.Header.BlockHash)
		// Print both the hash for debugging
		logger.InfoLogger.Printf("Block's hash: %x, Target hash: %x\n", hash, target)
	}
	return isValid
}

// ValidateMerkleRoot ensures the block's Merkle root matches the transactions.
func (bc *Blockchain) ValidateMerkleRoot(block *Block) bool {
	merkleTree, err := BuildTree(block.Transactions)
	// logger.InfoLogger.Printf("Received block's merkle root: %v\n", hex.EncodeToString(block.Header.MerkleRoot))
	// logger.InfoLogger.Printf("Calculated merkle root: %v\n", hex.EncodeToString(merkleTree.Root.Hash))
	if err != nil {
		logger.ErrorLogger.Println("Failed to build Merkle tree:", err)
		return false
	}
	return bytes.Equal(merkleTree.Root.Hash, block.Header.MerkleRoot)
}

// ValidateTimestamp checks if the block's timestamp is valid.
func (bc *Blockchain) ValidateTimestamp(block *Block) bool {
	lastBlock := bc.LatestBlock()
	isValid := block.Header.Timestamp > lastBlock.Header.Timestamp

	if !isValid {
		// Print both timestamps for debugging
		logger.ErrorLogger.Printf("Block timestamp is invalid. Last block timestamp: %d, Current block timestamp: %d\n", lastBlock.Header.Timestamp, block.Header.Timestamp)
	}

	return isValid
}

// ComputeHash calculates the hash of the block's header.
func (bc *Blockchain) ComputeHash(block *Block) []byte {
	// log the block contents
	logger.InfoLogger.Printf("Block: %v\n", block)
	// log the block header
	logger.InfoLogger.Printf("Block Header: %v\n", block.Header)

	if bytes.Equal(block.Header.BlockHash, bc.Ledger[0].Header.BlockHash) {
		return block.Header.BlockHash
	}
	return HashObject(SerializeBlockHeader(block.Header))

}

// CommitBlock commits blocks that have reached the confirmation depth.
func (bc *Blockchain) CommitBlock() {
	numBlocks := len(bc.Ledger)
	if numBlocks <= bc.ConfirmationDepth {
		logger.InfoLogger.Println("Not enough blocks to commit.")
		return
	}

	blockToCommitIndex := numBlocks - bc.ConfirmationDepth
	blockToCommit := bc.Ledger[blockToCommitIndex]

	// Create a list of UTXOTransaction
	var tempTransactions []*UTXOTransaction

	for _, tx := range blockToCommit.Transactions {
		// Remove UTXOs used in the transaction inputs
		for _, input := range tx.Inputs {
			if utxo, hasUTXO := bc.UTXOSet.Get(input.ID); hasUTXO {
				tempTransactions = append(tempTransactions, utxo)
				bc.UTXOSet.Remove(input.ID)
			} else {
				// Clear the ledger from this block onwards and add the transactions to the mempool
				for i := blockToCommitIndex; i < len(bc.Ledger); i++ {
					for _, tx := range bc.Ledger[i].Transactions {
						// neglect coinbase transactions
						if _, ok := tx.Data.(CoinbaseTransactionData); !ok {
							bc.Mempool.AddTransaction(tx)
						}
					}
				}
				bc.Ledger = bc.Ledger[:blockToCommitIndex]

				// Add all tempTransactions to utxoSet
				for _, tx := range tempTransactions {
					bc.UTXOSet.Add(tx)
				}
			}
		}

		// Add UTXOs created by the transaction outputs
		for _, output := range tx.Outputs {
			if _, hasUTXO := bc.UTXOSet.Get(output.ID); hasUTXO {
				logger.ErrorLogger.Printf("UTXO already exists for output ID: %v\n", output.ID)
				continue
			} else {
				bc.UTXOSet.Add(output)
			}
		}
	}

	logger.InfoLogger.Printf("Block at index %d committed successfully.\n", blockToCommitIndex)
}

// processResponses extracts blocks from received responses
func (bc *Blockchain) processResponses(responses []*GetBlocksResponse) [][]*Block {
	var allBlocks [][]*Block

	for _, response := range responses {
		if response.Success {
			allBlocks = append(allBlocks, response.Blocks)
		} else {
			logger.ErrorLogger.Println("Received unsuccessful response from a peer")
		}
	}

	logger.InfoLogger.Printf("Processed %d sets of blocks", len(allBlocks))
	return allBlocks
}

// GetBlockByHash retrieves a block from the ledger by its hash.
func (bc *Blockchain) GetBlockByHash(hash []byte) (*Block, error) {
	for _, block := range bc.Ledger {
		if bytes.Equal(block.Header.BlockHash, hash) {
			logger.InfoLogger.Printf("Block found for hash: %x\n", hash)
			return block, nil
		}
	}
	logger.ErrorLogger.Printf("Block not found for hash: %x\n", hash)
	return nil, ErrBlockNotFound
}

// LatestBlock returns the most recent block in the blockchain.
func (bc *Blockchain) LatestBlock() *Block {
	if len(bc.Ledger) == 0 {
		logger.ErrorLogger.Println("Blockchain is empty.")
		return nil
	}
	return bc.Ledger[len(bc.Ledger)-1]
}

// GetHeight returns the height of the blockchain.
func (bc *Blockchain) GetHeight() int {
	return len(bc.Ledger)
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
	case PurchaseTransactionData:
		return bc.validatePurchaseTransaction(tx, &data, utxoSet)
	case ReviewTransactionData:
		return bc.validateReviewTransaction(&data)
	case CoinbaseTransactionData:
		return nil
	case *CoinbaseTransactionData:
		return nil
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

// getBlocksSinceHash retrieves all blocks from the ledger after a given hash
func (bc *Blockchain) GetBlocksSinceHash(lastKnownHash []byte) ([]*Block, error) {
	blocks := bc.Ledger
	startIndex := -1

	// Find the index of the last known hash
	for i, block := range blocks {
		if bytes.Equal(block.Header.BlockHash, lastKnownHash) {
			startIndex = i + 1
			break
		}
	}

	if startIndex == -1 {
		return nil, errors.New("last known hash not found")
	}

	return blocks[startIndex:], nil
}
