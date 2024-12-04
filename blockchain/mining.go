package blockchain

type Miner struct {
	Blockchain *Blockchain
	Mempool    *Mempool
}

func NewMiner(bc *Blockchain, mp *Mempool) *Miner {
	return &Miner{Blockchain: bc, Mempool: mp}
}

func (m *Miner) MineBlock() (*Block, error) {
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
    txCount := m.Blockchain.Ledger[0].TransactionCount // Assuming block size from genesis block
    if m.Mempool.Transactions.Len() < txCount {
        logger.InfoLogger.Println("Not enough transactions to mine a block")
        return nil, nil
    }

    // Get transactions from mempool
    transactions := m.Mempool.GetTransactions(txCount)

    // Add coinbase transaction
    coinbaseTx := m.createCoinbaseTransaction()
    transactions = append([]*UTXOTransaction{coinbaseTx}, transactions...)

    // Create new block
    previousHash := m.Blockchain.LatestBlock().Header.BlockHash
    targetHash := m.Blockchain.LatestBlock().Header.TargetHash // Assuming same target
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
        hash = b.ComputeHash()
        if bytes.Compare(hash, b.Header.TargetHash) < 0 {
            b.Header.BlockHash = hash
            logger.InfoLogger.Println("Proof of Work successful with nonce:", nonce)
            break
        } else {
            nonce++
        }
    }
}

func (m *Miner) createCoinbaseTransaction() *UTXOTransaction {
    // Create a coinbase transaction rewarding the miner
    tx := &UTXOTransaction{
        ID: UTXOTransactionID{
            // Set appropriate BlockHash and TxIndex
        },
        Address: m.Wallet.BitcoinAddress,
        Amount:  m.Blockchain.MiningReward,
        Fee:     0,
    }
    logger.InfoLogger.Println("Coinbase transaction created for miner reward")
    return tx
}
