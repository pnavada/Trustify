package blockchain

import (
	"crypto/sha256"
	"encoding/pem"
	"errors"
	"trustify/logger"

	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	BitcoinAddress []byte
	PublicKey      []byte
	PrivateKey     []byte
	UTXOs          []*UTXOTransaction
}

func NewWallet(privateKeyPEM []byte) *Wallet {
	// Decode the PEM-encoded private key
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		panic("Invalid private key format")
	}

	// Parse the private key to obtain both private and public keys
	_, publicKey := btcec.PrivKeyFromBytes(block.Bytes)

	// Serialize the public key in compressed format
	publicKeyBytes := publicKey.SerializeCompressed()

	// Generate the Bitcoin address as bytes
	bitcoinAddress, err := calculateBitcoinAddress(publicKeyBytes)
	if err != nil {
		panic("Failed to calculate Bitcoin address")
	}

	// print bitcoin address
	logger.InfoLogger.Printf("Bitcoin address: %v and %v\n", bitcoinAddress, string(bitcoinAddress))

	// Return the wallet
	return &Wallet{
		BitcoinAddress: bitcoinAddress,
		PublicKey:      publicKeyBytes,
		PrivateKey:     privateKeyPEM,
		UTXOs:          make([]*UTXOTransaction, 0),
	}
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
