package blockchain

type Transaction struct {
	ID      string
	Inputs  []UTXOTransaction
	Outputs []UTXOTransaction
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

func NewPurchaseTransaction(w *Wallet, to string, amount int, fee int, productID string) *Transaction {
	// Create a new purchase transaction
	return nil
}

func NewReviewTransaction(w *Wallet, productID string, rating int) *Transaction {
	// Create a new review transaction
	return nil
}

func (tx *Transaction) Hash() []byte {
	// Generate the transaction hash
	return nil
}

func (tx *Transaction) Sign(privKey []byte) []byte {
	// Sign the transaction inputs
	return nil
}

func (tx *Transaction) Verify() bool {
	// Verify the transaction signatures
	return false
}
