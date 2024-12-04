package blockchain

type UTXOTransaction struct {
	ID      UTXOTransactionID
	Address []byte
	Amount  int
	Fee     int
}

type UTXOTransactionID struct {
	BlockHash []byte
	TxIndex   int
}

type UTXOSet struct {
	UTXOs map[string]*UTXOTransaction
}

func NewUTXOSet() *UTXOSet {
	// Initialize UTXO set
	return nil
}

func (u *UTXOSet) Add(utxo UTXOTransaction) {
	// Implementation here
}

func (u *UTXOSet) Remove(id string) {
	// Implementation here
}

func (u *UTXOSet) Get(id string) (UTXOTransaction, bool) {
	// Implementation here
	return UTXOTransaction{}, false
}

func (u *UTXOSet) GetAllForAddress(address []byte) []UTXOTransaction {
	// Implementation here
	return nil
}
