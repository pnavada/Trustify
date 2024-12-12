package blockchain

import (
	"sync"
	"time"
	"trustify/logger"
)

type Transaction struct {
	ID        []byte
	Inputs    []*UTXOTransaction
	Outputs   []*UTXOTransaction
	Data      interface{}
	Timestamp int64  // Add timestamp field
	Sequence  uint64 // Add sequence field for tracking order
}

type PurchaseTransactionData struct {
	BuyerAddress  []byte
	SellerAddress []byte
	ProductID     string
	Amount        int
}

type ReviewTransactionData struct {
	ReviewerAddress []byte
	Rating          int
	ProductID       string
}

type CoinbaseTransactionData struct {
	BlockHeight int
}

func (tx *Transaction) GetTransactionFee() int {
	fee := 0
	for _, input := range tx.Inputs {
		fee += input.Amount
	}
	for _, output := range tx.Outputs {
		fee -= output.Amount
	}
	return fee
}

func NewPurchaseTransaction(w *Wallet, to string, amount int, fee int, productID string) *Transaction {
	// Create a new purchase transaction
	// For a purchase transaction, the user's wallet contains the list of unspent tranactions used for spending
	// The amount is the amount to be spent
	// The fee is the transaction fee
	// Note that the amount does not include the transaction fee

	// Create inputs from UTXOs
	inputs, change, err := w.CreateInputs(amount + fee)
	if err != nil {
		logger.ErrorLogger.Println("Failed to create inputs for purchase transaction:", err)
		return nil
	}

	// Create outputs
	outputs := []*UTXOTransaction{
		{Address: []byte(to), Amount: amount},
	}
	if change > 0 {
		outputs = append(outputs, &UTXOTransaction{Address: w.BitcoinAddress, Amount: change})
	}

	txData := &PurchaseTransactionData{
		BuyerAddress:  w.BitcoinAddress,
		SellerAddress: []byte(to),
		ProductID:     productID,
		Amount:        amount,
	}

	tx := &Transaction{
		Inputs:    inputs,
		Outputs:   outputs,
		Data:      txData,
		Timestamp: time.Now().Unix(),        // Add timestamp
		Sequence:  generateSequenceNumber(), // Implement this method
	}

	// Serialize the transaction and generate a hash
	Serialized := SerializeTransaction(tx)
	hashed := HashObject(Serialized)
	tx.ID = hashed

	// Assign the IDS for the outputs
	for i, output := range tx.Outputs {
		output.ID = &UTXOTransactionID{
			TxHash:  tx.ID,
			TxIndex: i,
		}
	}

	// List down the inputs and outputs for the transaction
	logger.InfoLogger.Println("Transaction:", tx.ID)
	for _, input := range tx.Inputs {
		logger.InfoLogger.Printf("Input - Address: %x, Amount: %d\n", input.Address, input.Amount)
	}

	logger.InfoLogger.Println("Transaction Outputs:")
	for _, output := range tx.Outputs {
		logger.InfoLogger.Printf("Output - Address: %x, Amount: %d\n", output.Address, output.Amount)
	}

	logger.InfoLogger.Println("New purchase transaction created:", tx.ID, " with data:", txData)
	return tx
}

func NewReviewTransaction(w *Wallet, productID string, rating int) *Transaction {
	// Create a new review transaction
	txData := &ReviewTransactionData{
		ReviewerAddress: w.BitcoinAddress,
		Rating:          rating,
		ProductID:       productID,
	}

	tx := &Transaction{
		Data:      txData,
		Timestamp: time.Now().Unix(),        // Add timestamp
		Sequence:  generateSequenceNumber(), // Implement this method
	}

	// Serialize the transaction and generate a hash
	Serialized := SerializeTransaction(tx)
	hashed := HashObject(Serialized)
	tx.ID = hashed
	logger.InfoLogger.Println("New review transaction created:", tx.ID)
	return tx
}

// Helper method to generate unique sequence number
var (
	sequenceMutex   sync.Mutex
	currentSequence uint64 = 0
)

func generateSequenceNumber() uint64 {
	sequenceMutex.Lock()
	defer sequenceMutex.Unlock()
	currentSequence++
	return currentSequence
}
