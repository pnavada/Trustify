package blockchain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"trustify/logger"
)

// Miner represents a miner in the blockchain system.
type Miner struct {
	Blockchain *Blockchain // Reference to the blockchain
	Mempool    *Mempool    // Reference to the transaction pool
	Wallet     *Wallet     // Wallet of the miner
}

// NewMiner creates a new Miner instance.
func NewMiner(blockchain *Blockchain, mempool *Mempool, wallet *Wallet) *Miner {
	return &Miner{
		Blockchain: blockchain,
		Mempool:    mempool,
		Wallet:     wallet,
	}
}

// MineBlock mines a new block by selecting transactions, creating a block, performing proof of work, and adding it to the blockchain.
func (m *Miner) MineBlock(blockSize int) (*Block, error) {
	logger.InfoLogger.Println("Starting block mining process")

	// Step 1: Retrieve transactions from the mempool.
	transactions := m.Mempool.GetTransactions(blockSize)
	if len(transactions) == 0 {
		logger.ErrorLogger.Println("Mempool is empty. Aborting mining process.")
		return nil, errors.New("not enough transactions to mine a block")
	}
	logger.InfoLogger.Printf("Retrieved %d transactions from mempool", len(transactions))

	// Step 2: Create a coinbase transaction to reward the miner.
	coinbaseTx := m.createCoinbaseTransaction(len(m.Blockchain.Ledger))

	// Step 3: Process transactions for additional rewards.
	m.processTransactionRewards(transactions, coinbaseTx)

	// Include the coinbase transaction at the beginning.
	transactions = append([]*Transaction{coinbaseTx}, transactions...)

	// Step 4: Create a new block.
	previousHash := m.Blockchain.LatestBlock().Header.BlockHash
	targetHash := m.Blockchain.LatestBlock().Header.TargetHash
	block, err := NewBlock(transactions, previousHash, targetHash)
	logger.InfoLogger.Println("Block successfully created for mining")
	logger.InfoLogger.Printf("Block's merkle root: %v\n", hex.EncodeToString(block.Header.MerkleRoot))
	if err != nil {
		logger.ErrorLogger.Printf("Failed to create new block: %v", err)
		return nil, err
	}
	logger.InfoLogger.Println("New block created successfully")

	logger.InfoLogger.Printf("Target hash for new block: %s", hex.EncodeToString(targetHash))

	// Step 5: Perform Proof of Work.
	m.performProofOfWork(block)

	// Step 6: Add the mined block to the blockchain.
	if err := m.Blockchain.AddBlock(block); err != nil {
		logger.ErrorLogger.Printf("Failed to add block to blockchain: %v", err)
		return nil, err
	}
	logger.InfoLogger.Printf("Block successfully added to blockchain with hash: %x", block.Header.BlockHash)

	return block, nil
}

// performProofOfWork calculates a valid nonce for the block by solving the proof-of-work puzzle.
func (m *Miner) performProofOfWork(block *Block) {
	logger.InfoLogger.Println("Starting proof of work")
	var hash []byte
	var nonce int64

	for {
		block.Header.Nonce = nonce
		hash = HashObject(SerializeBlockHeader(block.Header))
		logger.InfoLogger.Printf("Current nonce: %d, Hash: %x", nonce, hash)
		if bytes.Compare(hash, block.Header.TargetHash) < 0 {
			// block.Header.BlockHash = hash
			logger.InfoLogger.Printf("Proof of work completed. Nonce: %d, Hash: %x", nonce, hash)
			break
		}
		nonce++
	}
}

// createCoinbaseTransaction creates a coinbase transaction to reward the miner.
func (m *Miner) createCoinbaseTransaction(blockHeight int) *Transaction {
	logger.InfoLogger.Println("Creating coinbase transaction for miner reward")
	coinbaseTx := &Transaction{
		Outputs: make([]*UTXOTransaction, 1),
		Data: CoinbaseTransactionData{
			BlockHeight: blockHeight,
		},
	}
	coinbaseTx.Outputs[0] = &UTXOTransaction{
		Address: m.Wallet.BitcoinAddress,
		Amount:  m.Blockchain.MiningReward,
	}
	coinbaseTx.ID = HashObject(Serialize(coinbaseTx))
	coinbaseTx.Outputs[0].ID = &UTXOTransactionID{
		TxHash:  coinbaseTx.ID,
		TxIndex: 0,
	}
	logger.InfoLogger.Println("Coinbase transaction successfully created")
	return coinbaseTx
}

// processTransactionRewards adds transaction fees and review rewards to the coinbase transaction.
func (m *Miner) processTransactionRewards(transactions []*Transaction, coinbaseTx *Transaction) {
	logger.InfoLogger.Println("Processing transaction rewards")
	for _, tx := range transactions {
		switch data := tx.Data.(type) {
		case *PurchaseTransactionData:
			logger.InfoLogger.Printf("Adding transaction fee for purchase transaction: %x", tx.ID)
			coinbaseTx.Outputs = append(coinbaseTx.Outputs, &UTXOTransaction{
				ID:      coinbaseTx.Outputs[0].ID,
				Address: m.Wallet.BitcoinAddress,
				Amount:  tx.GetTransactionFee(),
			})
		case *ReviewTransactionData:
			logger.InfoLogger.Printf("Adding review reward for transaction: %x", tx.ID)
			coinbaseTx.Outputs = append(coinbaseTx.Outputs, &UTXOTransaction{
				ID:      coinbaseTx.Outputs[0].ID,
				Address: data.ReviewerAddress,
				Amount:  m.Blockchain.ReviewReward,
			})
		default:
			logger.ErrorLogger.Printf("Unknown transaction type for transaction: %x", tx.ID)
		}
	}
}
