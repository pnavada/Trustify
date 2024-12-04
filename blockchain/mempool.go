package blockchain

import (
	"container/heap"
	"sync"
)

type TransactionHeap []*UTXOTransaction

func (th TransactionHeap) Len() int { return len(th) }
func (th TransactionHeap) Less(i, j int) bool {
	return th[i].Fee > th[j].Fee // Higher fee, higher priority
}
func (th TransactionHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

func (th *TransactionHeap) Push(x interface{}) {
	*th = append(*th, x.(*UTXOTransaction))
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

func (mp *Mempool) AddTransaction(tx *UTXOTransaction) {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	heap.Push(mp.Transactions, tx)
}

func (mp *Mempool) GetTransactions(count int) []*UTXOTransaction {
	mp.Mutex.Lock()
	defer mp.Mutex.Unlock()
	var txs []*UTXOTransaction
	for i := 0; i < count && mp.Transactions.Len() > 0; i++ {
		tx := heap.Pop(mp.Transactions).(*UTXOTransaction)
		txs = append(txs, tx)
	}
	return txs
}
