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
	// Return a boolean indicating success or failure
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	key := utxo.ID.String()
	if _, exists := u.UTXOs[key]; exists {
		logger.InfoLogger.Println("UTXO already exists:", key)
		return false
	}
	u.UTXOs[key] = utxo
	logger.InfoLogger.Println("UTXO added:", key)
	return true
}

func (u *UTXOSet) Remove(id *UTXOTransactionID) bool {
	// Remove the transaction from the set
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
		logger.InfoLogger.Println("UTXO not found:", key)
		logger.InfoLogger.Println("Expected UTXO:", id)
		logger.InfoLogger.Println("Available UTXOs:")
		for k := range u.UTXOs {
			logger.InfoLogger.Println(k)
		}
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
