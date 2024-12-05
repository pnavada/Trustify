package blockchain

import (
	"bytes"
	"trustify/logger"
)

type Miner struct {
	Blockchain *Blockchain
	Mempool    *Mempool
	Wallet     *Wallet
}

func NewMiner(bc *Blockchain, mp *Mempool, wallet *Wallet) *Miner {
	return &Miner{Blockchain: bc, Mempool: mp, Wallet: wallet}
}

func (m *Miner) MineBlock(blockSize int) (*Block, error) {
	// Collect transactions, create block, perform proof of work
	// This method should be called when there are enough transactions in the mempool
	// form a block.
	// The required block size is defined in the configuation object
	// Then pick the required number of transactions from the mempool
	// Add coinbase transactions to reward the miners and reviewers (if applicable).
	// Use your best decision to position this logic
	// Create the block using the NewBlock method
	// The target hash is part of the configuration. Dont receive the confuguration object as a parameter
	// USe intelligent ways to pass this data between objects
	// Define the parameters where it makes the most sense to be a part of
	// Perform proof of work to find the valid nonce for the block
	// Add the block to the ledger and broadcast over the network
	// Ensure there are enough transactions

	// Get transactions from mempool
	transactions := m.Mempool.GetTransactions(blockSize)

	// Add coinbase transaction
	coinbaseTx := m.createCoinbaseTransaction(len(m.Blockchain.Ledger))

	// Add transaction fee and review rewards
	for _, tx := range transactions {
		switch data := tx.Data.(type) {
		case *PurchaseTransactionData:
			// Handle ProductTransactionData
			logger.InfoLogger.Println("Processing PurchaseTransactionData for transaction:", tx.ID)
			coinbaseTx.Outputs = append(coinbaseTx.Outputs, &UTXOTransaction{
				ID:      coinbaseTx.Outputs[0].ID,
				Address: m.Wallet.BitcoinAddress,
				Amount:  tx.GetTransactionFee(),
			})
			// Add your logic for ProductTransactionData here
		case *ReviewTransactionData:
			// Handle ReviewTransactionData
			logger.InfoLogger.Println("Processing ReviewTransactionData for transaction:", tx.ID)
			coinbaseTx.Outputs = append(coinbaseTx.Outputs, &UTXOTransaction{
				ID:      coinbaseTx.Outputs[0].ID,
				Address: data.ReviewerAddress,
				Amount:  m.Blockchain.ReviewReward,
			})
			// Add your logic for ReviewTransactionData here
		default:
			logger.InfoLogger.Println("Unknown transaction data type for transaction:", tx.ID)
		}
	}

	transactions = append([]*Transaction{coinbaseTx}, transactions...)

	// Create new block
	previousHash := m.Blockchain.LatestBlock().Header.BlockHash
	targetHash := m.Blockchain.LatestBlock().Header.TargetHash
	block, err := NewBlock(transactions, previousHash, targetHash)
	if err != nil {
		logger.ErrorLogger.Println("Failed to create new block:", err)
		return nil, err
	}

	// Perform Proof of Work
	m.ProofOfWork(block)

	// Add block to blockchain
	err = m.Blockchain.AddBlock(block)
	if err != nil {
		logger.ErrorLogger.Println("Failed to add block to blockchain:", err)
		return nil, err
	}

	// Broadcast block
	// Assuming n.BroadcastBlock(block) exists
	logger.InfoLogger.Println("New block mined and added to blockchain:", block.Header.BlockHash)
	return block, nil

}

func (m *Miner) ProofOfWork(b *Block) {
	// Perform POW to find valid nonce
	// Here  the block is the block for which the nonce is to be found
	// The nonce is intiialized to zero and incremented until the hash of the block is less than the target hash
	// Basic idea is to find a nonce such that the hash of the block is less than the target hash
	var hash []byte
	nonce := int64(0)
	for {
		b.Header.Nonce = nonce
		hash = HashObject(Serialize(b))
		if bytes.Compare(hash, b.Header.TargetHash) < 0 {
			b.Header.BlockHash = hash
			logger.InfoLogger.Println("Proof of Work successful with nonce:", nonce)
			break
		} else {
			nonce++
		}
	}

}

func (m *Miner) createCoinbaseTransaction(numBlocks int) *Transaction {
	// Create a coinbase transaction rewarding the miner
	coinbaseTx := &Transaction{
		Outputs: make([]*UTXOTransaction, 1),
		Data: CoinbaseTransactionData{
			BlockHeight: numBlocks,
		},
	}
	coinbaseTx.Outputs[0] = &UTXOTransaction{
		Address: m.Wallet.BitcoinAddress,
		Amount:  m.Blockchain.MiningReward,
	}
	coinbaseTx.ID = HashObject(Serialize(coinbaseTx))
	coinbaseTx.Outputs[0].ID = UTXOTransactionID{
		TxHash:  coinbaseTx.ID,
		TxIndex: 0,
	}
	logger.InfoLogger.Println("Coinbase transaction created for miner reward")
	return coinbaseTx
}
