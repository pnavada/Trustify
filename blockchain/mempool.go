package blockchain

import (
	"bytes"
	"container/heap"
	"sync"
	"trustify/logger"
)

type TransactionHeap []*Transaction

func (th TransactionHeap) Len() int { return len(th) }
func (th TransactionHeap) Less(i, j int) bool {
	return th[i].GetTransactionFee() > th[j].GetTransactionFee() // Higher fee, higher priority
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

type Mempool struct {
	Transactions *TransactionHeap
	Mutex        sync.Mutex
}

func NewMempool() *Mempool {
	th := &TransactionHeap{}
	heap.Init(th)
	return &Mempool{Transactions: th}
}

func (mp *Mempool) AddTransaction(tx *Transaction) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()

	heap.Push(mp.Transactions, tx)
	logger.InfoLogger.Println("Transaction added to mempool: ", tx.ID)
}

func (mp *Mempool) RemoveTransaction(tx *Transaction) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	for i, t := range *mp.Transactions {
		if bytes.Equal(t.ID, tx.ID) {
			heap.Remove(mp.Transactions, i)
			break
		}
	}
}

func (mp *Mempool) HasTransaction(tx *Transaction) (*Transaction, bool) { // THIS IS WITHOUT LOCKING and should be called only after locking
	for _, t := range *mp.Transactions {
		if bytes.Equal(t.ID, tx.ID) {
			return t, true
		}
	}
	return nil, false
}

func (mp *Mempool) GetTransactions(count int) []*Transaction {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	var txs []*Transaction
	for i := 0; i < count && mp.Transactions.Len() > 0; i++ {
		tx := heap.Pop(mp.Transactions).(*Transaction)
		txs = append(txs, tx)
	}
	return txs
}
