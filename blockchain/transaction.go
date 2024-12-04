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
	// For a purchase transaction, the user's wallet contains the list of unspent tranactions used for spending
	// The amount is the amount to be spent
	// The fee is the transaction fee
	// Note that the amount does not include the transaction fee
	return nil
}

func NewReviewTransaction(w *Wallet, productID string, rating int) *Transaction {
	// Create a new review transaction
	return nil
}

func (tx *Transaction) Hash() []byte {
	// Generate the hash for the transaction
	return nil
}

func (tx *Transaction) Sign(privKey []byte) []byte {
	// Digitally sign the transaction using the private key
	return nil
}

func (tx *Transaction) Verify() bool {
	// Verify the transaction signatures
	return false
}
