package blockchain

import (
	"trustify/logger"
)

type Transaction struct {
	ID      []byte
	Inputs  []*UTXOTransaction
	Outputs []*UTXOTransaction
	Data    TransactionData
}

type TransactionData interface{}

type PurchaseTransactionData struct {
	TransactionData
	BuyerAddress  []byte
	SellerAddress []byte
	ProductID     string
	Amount        int
}

type ReviewTransactionData struct {
	TransactionData
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
		Inputs:  inputs,
		Outputs: outputs,
		Data:    txData,
	}

	// Serialize the transaction and generate a hash
	Serialized := SerializeTransaction(tx)
	hashed := HashObject(Serialized)
	tx.ID = hashed
	logger.InfoLogger.Println("New purchase transaction created:", tx.ID)
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
		Data: txData,
	}

	// Serialize the transaction and generate a hash
	Serialized := SerializeTransaction(tx)
	hashed := HashObject(Serialized)
	tx.ID = hashed
	logger.InfoLogger.Println("New review transaction created:", tx.ID)
	return tx
}
