package blockchain

type Wallet struct {
	BitcoinAddress []byte
	PublicKey      []byte
	PrivateKey     []byte
	UTXOs          []UTXOTransaction
}

func NewWallet(privateKey string) *Wallet {
	// Generate wallet from private key
	return nil
}

func (w *Wallet) GetBalance() int {
	// Calculate balance from UTXOs
	return 0
}

func (w *Wallet) SignTransaction(tx *UTXOTransaction) {
	// Sign transaction using wallet's private key
}
