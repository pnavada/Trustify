package blockchain

// import (
//     "trustify/crypto"
// )

type Transaction struct {
	ID      []byte
	Inputs  []*UTXOTransaction
	Outputs []*UTXOTransaction
	Data    TransactionData
	Fee     int
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
	// inputs, change, err := w.CreateInputs(amount + fee)
	// if err != nil {
	//     logger.ErrorLogger.Println("Failed to create inputs for purchase transaction:", err)
	//     return nil, err
	// }

	// // Create outputs
	// outputs := []UTXOTransaction{
	//     {Address: to, Amount: amount},
	// }
	// if change > 0 {
	//     outputs = append(outputs, UTXOTransaction{Address: w.BitcoinAddress, Amount: change})
	// }

	// txData := &PurchaseTransactionData{
	//     BuyerAddress:  w.BitcoinAddress,
	//     SellerAddress: to,
	//     ProductID:     productID,
	//     Amount:        amount,
	// }

	// tx := &Transaction{
	//     Inputs:  inputs,
	//     Outputs: outputs,
	//     Data:    txData,
	// }

	// tx.ID = tx.Hash()
	// crypto.Sign(w.PrivateKey)
	// logger.InfoLogger.Println("New purchase transaction created:", tx.ID)
	// return tx, nil
	return nil
}

func NewReviewTransaction(w *Wallet, productID string, rating int) *Transaction {
	// Create a new review transaction
	// txData := &ReviewTransactionData{
	//     ReviewerAddress: w.BitcoinAddress,
	//     Rating:          rating,
	//     ProductID:       productID,
	// }

	// tx := &Transaction{
	//     Data: txData,
	// }

	// tx.ID = tx.Hash()
	// crypto.Sign(w.PrivateKey)
	// logger.InfoLogger.Println("New review transaction created:", tx.ID)
	// return tx, nil
	return nil
}

func (tx *Transaction) Hash() []byte {
	// // Generate the hash for the transaction
	// data := crypto.Serialize(tx)
	// hash := sha256.Sum256(data)
	// return hash[:]
	return nil
}

func (tx *Transaction) Sign(privKey []byte) []byte {
	// Digitally sign the transaction using the private key
	// signature, err := crypto.Sign(tx.Hash(), privKey)
	// if err != nil {
	//     logger.ErrorLogger.Println("Failed to sign transaction:", err)
	//     return err
	// }
	// crypto.Signature = signature
	return nil
}

func (tx *Transaction) Verify() bool {
	// Extract public key from address
	// var pubKey []byte
	// switch data := tx.Data.(type) {
	// case *PurchaseTransactionData:
	//     pubKey = data.BuyerAddress
	// case *ReviewTransactionData:
	//     pubKey = data.ReviewerAddress
	// default:
	//     logger.ErrorLogger.Println("Unknown transaction data type")
	//     return false
	// }

	// valid := crypto.Verify(tx.Hash(), tx.Signature, pubKey)
	// if !valid {
	//     logger.ErrorLogger.Println("Transaction signature verification failed")
	// }
	// return valid
	return false
}
