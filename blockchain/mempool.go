package blockchain

import (
	"bytes"
	"container/heap"
	"sync"
	"time"
	"trustify/logger"
)

// TransactionHeap implements heap.Interface and holds Transactions.
type TransactionHeap []*Transaction

func (th TransactionHeap) Len() int { return len(th) }

func (th TransactionHeap) Less(i, j int) bool {
	// Prioritize by fee, then by sequence (for transactions with same fee)
	if th[i].GetTransactionFee() == th[j].GetTransactionFee() {
		return th[i].Sequence < th[j].Sequence
	}
	return th[i].GetTransactionFee() > th[j].GetTransactionFee()
}

func (th TransactionHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

func (th *TransactionHeap) Push(x interface{}) {
	*th = append(*th, x.(*Transaction))
}

func (th *TransactionHeap) Pop() interface{} {
	old := *th
	n := len(old)
	x := old[n-1]
	*th = old[0 : n-1]
	return x
}

// Mempool manages pending transactions.
type Mempool struct {
	Transactions *TransactionHeap
	TxMap        map[string]*Transaction // Map from TxID to Transaction
	Mutex        sync.RWMutex            // Read-Write mutex for better concurrency
}

// NewMempool initializes a new Mempool.
func NewMempool() *Mempool {
	th := &TransactionHeap{}
	heap.Init(th)
	return &Mempool{
		Transactions: th,
		TxMap:        make(map[string]*Transaction),
	}
}

// AddTransaction adds a new transaction to the mempool if it's not a duplicate.
// Returns true if the transaction was added, false otherwise.
func (mp *Mempool) AddTransaction(tx *Transaction) bool {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	txID := string(tx.ID)
	if _, exists := mp.TxMap[txID]; exists {
		logger.InfoLogger.Printf("Transaction already in mempool: %x\n", tx.ID)
		return false
	}

	// Add timestamp when transaction enters mempool
	tx.Timestamp = time.Now().Unix()

	heap.Push(mp.Transactions, tx)
	mp.TxMap[txID] = tx
	logger.InfoLogger.Printf("Transaction added to mempool: %x, Pool size: %d\n", tx.ID, mp.Transactions.Len())
	return true
}

// RemoveTransaction removes a transaction from the mempool.
func (mp *Mempool) RemoveTransaction(tx *Transaction) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	txID := string(tx.ID)
	if _, exists := mp.TxMap[txID]; exists {
		for i, t := range *mp.Transactions {
			if bytes.Equal(t.ID, tx.ID) {
				heap.Remove(mp.Transactions, i)
				delete(mp.TxMap, txID)
				logger.InfoLogger.Printf("Transaction removed from mempool: %x, Pool size: %d\n", tx.ID, mp.Transactions.Len())
				break
			}
		}
	}
}

// ContainsTransaction checks if a transaction exists in the mempool.
func (mp *Mempool) ContainsTransaction(tx *Transaction) bool {
	mp.Mutex.RLock()
	defer mp.Mutex.RUnlock()

	_, exists := mp.TxMap[string(tx.ID)]
	return exists
}

// GetTransactions retrieves and removes up to 'count' transactions from the mempool.
// Modified GetTransactions to include transaction tracking
func (mp *Mempool) GetTransactions(count int) []*Transaction {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	var txs []*Transaction
	usedTxIDs := make(map[string]bool)

	for i := 0; i < count && mp.Transactions.Len() > 0; i++ {
		tx := heap.Pop(mp.Transactions).(*Transaction)
		txID := string(tx.ID)

		// Additional check to prevent duplicate transactions
		if _, used := usedTxIDs[txID]; used {
			logger.InfoLogger.Printf("Skipping already used transaction: %x\n", tx.ID)
			continue
		}

		txs = append(txs, tx)
		usedTxIDs[txID] = true
		delete(mp.TxMap, txID)

		logger.InfoLogger.Printf("Transaction fetched for block: %x, Remaining pool size: %d\n", tx.ID, mp.Transactions.Len())
	}
	return txs
}

// HasTransaction remains for legacy support or internal usage, but it's recommended to use ContainsTransaction instead.
func (mp *Mempool) HasTransaction(tx *Transaction) (*Transaction, bool) {
	mp.Mutex.RLock()
	defer mp.Mutex.RUnlock()

	existingTx, exists := mp.TxMap[string(tx.ID)]
	return existingTx, exists
}
