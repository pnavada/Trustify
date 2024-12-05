package blockchain

import (
	"bytes"
	"fmt"
	"sync"
	"trustify/logger"
)

type UTXOTransaction struct {
	ID      *UTXOTransactionID
	Address []byte
	Amount  int
}

type UTXOTransactionID struct {
	TxHash  []byte
	TxIndex int
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
	Mutex sync.Mutex
}

func NewUTXOSet() *UTXOSet {
	// Initialize UTXO set
	return &UTXOSet{UTXOs: make(map[string]*UTXOTransaction)}
}

// Helper method to convert UTXOTransactionID to string
func (id UTXOTransactionID) String() string {
	return fmt.Sprintf("%x:%d", id.TxHash, id.TxIndex)
}

func (u *UTXOSet) Add(utxo *UTXOTransaction) bool {
	// Add a UXTO transaction to the set
	// Make sure the transaction is unique
	// There cannot be duplocates in a set!
	// Return a boolean indicating success or failure
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	key := utxo.ID.String()
	if _, exists := u.UTXOs[key]; exists {
		logger.ErrorLogger.Println("UTXO already exists:", key)
		return false
	}
	u.UTXOs[key] = utxo
	logger.InfoLogger.Println("UTXO added:", key)
	return true
}

func (u *UTXOSet) Remove(id *UTXOTransactionID) bool {
	// Remove the transaction from the set
	// Return a boolean indicating success or failure
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	key := id.String()
	if _, exists := u.UTXOs[key]; !exists {
		logger.ErrorLogger.Println("UTXO not found:", key)
		return false
	}
	delete(u.UTXOs, key)
	logger.InfoLogger.Println("UTXO removed:", key)
	return true
}

func (u *UTXOSet) Get(id *UTXOTransactionID) (*UTXOTransaction, bool) {
	// Get the transaction
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	key := id.String()
	utxo, exists := u.UTXOs[key]
	if !exists {
		logger.ErrorLogger.Println("UTXO not found:", key)
		return nil, false
	}
	return utxo, true
}

func (u *UTXOSet) GetAllForAddress(address []byte) []*UTXOTransaction {
	// // Get all transcations for the specified address
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	var utxos []*UTXOTransaction
	for _, utxo := range u.UTXOs {
		if bytes.Equal(utxo.Address, address) {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}
