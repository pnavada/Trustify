package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

// NodeKeyPair stores the private key, public key, and Bitcoin address
type NodeKeyPair struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

// GenerateKeyPair generates an ECDSA key pair and Bitcoin address
func GenerateUniqueKeyPair() (*NodeKeyPair, error) {
	// Step 1: Generate private key using secp256k1
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Use secp256k1 for Bitcoin
	if err != nil {
		return nil, err
	}

	// Step 2: Extract the public key
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...) // Uncompressed

	// Step 3: Generate Bitcoin address
	address, err := GenerateBitcoinAddress(pubKey)
	if err != nil {
		return nil, err
	}

	// Encode private and public keys as hex strings
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := hex.EncodeToString(pubKey)

	return &NodeKeyPair{
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
		Address:    address,
	}, nil
}

// GenerateBitcoinAddress creates a Bitcoin address from the public key
func GenerateBitcoinAddress(pubKey []byte) (string, error) {
	// Step 1: Perform SHA-256 hash on the public key
	shaHash := sha256.Sum256(pubKey)

	// Step 2: Perform RIPEMD-160 hash on the result of SHA-256
	ripemd := ripemd160.New()
	_, err := ripemd.Write(shaHash[:])
	if err != nil {
		return "", err
	}
	publicKeyHash := ripemd.Sum(nil)

	// Step 3: Add version byte (0x00 for Mainnet)
	versionedHash := append([]byte{0x00}, publicKeyHash...)

	// Step 4: Perform double SHA-256 to calculate the checksum
	checksum := sha256.Sum256(versionedHash)
	checksum = sha256.Sum256(checksum[:])

	// Step 5: Append the first 4 bytes of the checksum to the versioned hash
	fullHash := append(versionedHash, checksum[:4]...)

	// Step 6: Encode the result using Base58 encoding
	address := Base58Encode(fullHash)
	return address, nil
}

// Base58Encode encodes a byte slice to a Base58 string
func Base58Encode(input []byte) string {
	const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	result := ""
	intVal := new(big.Int).SetBytes(input)

	for intVal.Sign() > 0 {
		mod := new(big.Int)
		intVal.DivMod(intVal, big.NewInt(58), mod)
		result = string(base58Alphabet[mod.Int64()]) + result
	}

	// Handle leading zero bytes
	for _, b := range input {
		if b == 0x00 {
			result = string(base58Alphabet[0]) + result
		} else {
			break
		}
	}

	return result
}

// GenerateKeyPairs generates key pairs for n nodes
func GenerateKeyPairs(n int) ([]NodeKeyPair, error) {
	var keyPairs []NodeKeyPair
	for i := 0; i < n; i++ {
		keyPair, err := GenerateUniqueKeyPair()
		if err != nil {
			return nil, err
		}
		keyPairs = append(keyPairs, *keyPair)
	}
	return keyPairs, nil
}

func main2() {
	n := 5 // Number of nodes (replace with desired number)
	keyPairs, err := GenerateKeyPairs(n)
	if err != nil {
		fmt.Printf("Error generating key pairs: %v\n", err)
		return
	}

	fmt.Println("Generated Key Pairs:")
	for i, kp := range keyPairs {
		fmt.Printf("Node %d:\n", i+1)
		fmt.Printf("  Private Key: %s\n", kp.PrivateKey)
		fmt.Printf("  Public Key: %s\n", kp.PublicKey)
		fmt.Printf("  Address: %s\n", kp.Address)
	}
}
