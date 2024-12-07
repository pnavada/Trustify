package blockchain

import (
	"bytes"
	"container/heap"
	"sync"
)

// Feel free to correct any mistakes or define new methods if needed

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
}

func (mp *Mempool) RemoveTransaction(tx *Transaction) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	for i, t := range *mp.Transactions {
		if t.Equals(tx) {
			heap.Remove(mp.Transactions, i)
			break
		}
	}
}

func (tx *Transaction) Equals(other *Transaction) bool {
	SerializedTx := SerializeTransaction(tx)
	hashedTx := HashObject(SerializedTx)

	SerializedOther := SerializeTransaction(other)
	hashedOther := HashObject(SerializedOther)

	// Compare the hashes of the two transactions
	return bytes.Equal(hashedTx, hashedOther)
}

func (mp *Mempool) HasTransaction(tx *Transaction) bool {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	for _, t := range *mp.Transactions {
		if t.Equals(tx) {
			return true
		}
	}
	return false
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
