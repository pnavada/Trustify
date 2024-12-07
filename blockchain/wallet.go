package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"trustify/cryptography"
	"trustify/logger"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	BitcoinAddress []byte
	PublicKey      *ecdsa.PublicKey
	PrivateKey     *ecdsa.PrivateKey
	UTXOs          []*UTXOTransaction
}

func NewWallet(privateKeyPEM string) (*Wallet, error) {
	// Parse the private key
	privateKey, err := cryptography.ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	// Serialize the public key
	publicKey, err := cryptography.SerializePublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	bitcoinAddress, err := calculateBitcoinAddress(publicKey)
	if err != nil {
		return nil, err
	}

	// PRINT the bitcoin address
	logger.InfoLogger.Println("Bitcoin address:", string(bitcoinAddress))

	// Return the wallet
	return &Wallet{
		BitcoinAddress: bitcoinAddress,
		PublicKey:      &privateKey.PublicKey,
		PrivateKey:     privateKey,
		UTXOs:          make([]*UTXOTransaction, 0),
	}, nil
}

func calculateBitcoinAddress(publicKey []byte) ([]byte, error) {
	// Step 1: Perform SHA-256 hashing on the public key
	sha256Hash := sha256.Sum256(publicKey)

	// Step 2: Perform RIPEMD-160 hashing on the SHA-256 hash
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		return nil, errors.New("failed to hash public key with RIPEMD-160")
	}
	publicKeyHash := ripemd160Hasher.Sum(nil)

	// Step 3: Add Network Prefix (0x00 for Bitcoin mainnet)
	versionedPayload := append([]byte{0x00}, publicKeyHash...)

	// Step 4: Calculate Checksum (first 4 bytes of double SHA-256)
	checksum := sha256.Sum256(versionedPayload)
	checksum = sha256.Sum256(checksum[:])
	finalPayload := append(versionedPayload, checksum[:4]...)

	// Return the Bitcoin address as bytes
	return finalPayload, nil
}

func (w *Wallet) GetBalance() int {
	// Calculate balance from UTXOs
	// Calculate balance from a list of UTXO transactions
	balance := 0
	for _, utxo := range w.UTXOs {
		balance += utxo.Amount
	}
	logger.InfoLogger.Println("Wallet balance calculated:", balance)
	return balance
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

	// log total amount
	logger.InfoLogger.Println("Total amount:", total)

	if total < amount {
		logger.ErrorLogger.Println("Insufficient funds")
		return nil, 0, ErrInsufficientFunds
	}
	change := total - amount
	return inputs, change, nil
}
