package blockchain

import "trustify/logger"

type Wallet struct {
	BitcoinAddress []byte
	PublicKey      []byte
	PrivateKey     []byte
	UTXOs          []*UTXOTransaction
}

func NewWallet(privateKey []byte, publicKey []byte, bitcoinAddress []byte) *Wallet {
	// The wallet is initialized with an empty list of UTXOs
	// Other parameters are part of the configuration object
	return &Wallet{
		BitcoinAddress: bitcoinAddress,
		PublicKey:      publicKey,
		PrivateKey:     privateKey,
		UTXOs:          make([]*UTXOTransaction, 0),
	}
}

func (w *Wallet) GetBalance() int {
	// Calculate balance from UTXOs
	// Calculate balance from a list of UTXO transactions
	// balance := 0
	// for _, utxo := range w.UTXOs {
	//     balance += utxo.Amount
	// }
	// logger.InfoLogger.Println("Wallet balance calculated:", balance)
	// return balance
	return 0
}

func (w *Wallet) SignTransaction(tx *UTXOTransaction) {
	// TODO
	// return crypto.Sign(tx, w.PrivateKey)
}

func (w *Wallet) CreateInputs(amount int) ([]*UTXOTransaction, int, error) {
	var inputs []*UTXOTransaction
	total := 0
	for _, utxo := range w.UTXOs {
		inputs = append(inputs, utxo)
		total += utxo.Amount
		if total >= amount {
			break
		}
	}
	if total < amount {
		logger.ErrorLogger.Println("Insufficient funds")
		return nil, 0, ErrInsufficientFunds
	}
	change := total - amount
	return inputs, change, nil
}
