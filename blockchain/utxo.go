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

// The reason to use UTXOSet is for faster lookups
// Instaed of scanning the entire blockchain
// The node refers to this set to validate if a transaction is unspent
// For instance, the input transactions of a purchase can be validated quickly
// You can easily say if a user is doing double spending
// Note this set is updated only after a block is tagged as committed
// The block is assumed to be committed when its x block away from the tip
// The confirmation depth is part of the configuration object
// The UTXO set should be updated as and when needed
// When a block is committed, where the node earns rewards for mining
// Or when the node's review transaction is committed to the blockchain
// Also the reward transactions in the genesis block

type UTXOSet struct {
	UTXOs map[string]*UTXOTransaction
}

func NewUTXOSet() *UTXOSet {
	// Initialize UTXO set
	return nil
}

func (u *UTXOSet) Add(utxo UTXOTransaction) {
	// Add a UXTO transaction to the set
	// Make sure the transaction is unique
	// There cannot be duplocates in a set!
	// Return a boolean indicating success or failure
}

func (u *UTXOSet) Remove(id UTXOTransactionID) {
	// Remove the transaction from the set
	// Return a boolean indicating success or failure
}

func (u *UTXOSet) Get(id UTXOTransactionID) (UTXOTransaction, bool) {
	// Get the transaction
	return UTXOTransaction{}, false
}

func (u *UTXOSet) GetAllForAddress(address []byte) []UTXOTransaction {
	// Get all transcations for the specified address
	return nil
}
