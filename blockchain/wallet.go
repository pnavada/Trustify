package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"trustify/config"
	"trustify/logger"

	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	BitcoinAddress []byte
	PublicKey      []byte
	PrivateKey     []byte
	UTXOs          []*UTXOTransaction
}

func NewWallet(privateKeyPEM []byte) (*Wallet, error) {
	// Decode the PEM-encoded private key
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("invalid private key format")
	}

	var privKey *ecdsa.PrivateKey
	var err error

	switch block.Type {
	case "EC PRIVATE KEY":
		privKey, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	case "PRIVATE KEY":
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		privKey, ok = keyInterface.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("not ECDSA private key")
		}
	default:
		return nil, errors.New("unknown private key type")
	}

	// Serialize the public key in compressed format (33 bytes)
	// pubKeyBytes := append(
	// 	[]byte{},
	// 	privKey.PublicKey.X.Bytes()...,
	// )
	// You can choose to serialize in compressed or uncompressed format based on your needs
	// Here's how to serialize in compressed format:
	pubKeyCompressed := config.CompressPublicKey(&privKey.PublicKey)

	// Generate the Bitcoin address as bytes
	bitcoinAddress, err := calculateBitcoinAddress(pubKeyCompressed)
	if err != nil {
		return nil, errors.New("failed to calculate Bitcoin address")
	}

	// Log Bitcoin address
	logger.InfoLogger.Printf("Bitcoin address: %x and %s\n", bitcoinAddress, string(bitcoinAddress))

	// Return the wallet
	return &Wallet{
		BitcoinAddress: bitcoinAddress,
		PublicKey:      pubKeyCompressed,
		PrivateKey:     privateKeyPEM,
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
