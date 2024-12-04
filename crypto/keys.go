package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// Use the most preferred cryptographic library for generating key pairs and signing data in blockchain applications.
// TO-DO: Use a secure cryptographic library to generate an elliptic curve (or RSA) key pair.
// TO-DO: Ensure private keys are securely stored and not exposed.
// TO-DO: Gracefully handle invalid inputs, such as malformed keys or signatures.

// NOT IN USE
// func GenerateKeyPair() (KeyPair, error) {
// 	// Ensure the private key is securely generated and stored.
// 	// Derive the public key from the private key.
// 	// Return the key pair or an error if the generation fails.
// 	// Need to confirm - Use a standard curve like P256 for elliptic curve cryptography (or RSA-2048 for RSA).
// 	// Ensure randomness is sourced from a secure random number generator.
// 	return KeyPair{}, nil
// }

func Sign(data []byte, privateKey []byte) ([]byte, error) {
	// Use the most preferred cryptographic and hashing algorithms for blockchain applications.
	// Create a digital signature for the provided data using the private key.
	// Hash the data to create a digest.
	// Use the private key to sign the hashed data.
	// Return the signature in a format suitable for verification (e.g., DER-encoded).
	// Decode the private key
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid private key format")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Hash the data using SHA-256
	hash := sha256.Sum256(data)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Serialize the signature
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func Verify(data []byte, signature []byte, publicKey []byte) bool {
	//  Verify the digital signature of the data using the public key.
	// Hash the input data to create a digest (same method used in Sign).
	// Use the public key to verify the signature against the hash.
	// Return true if the signature is valid, otherwise false.
	// Ensure the public key is correctly formatted and corresponds to the private key used for signing.

	// Decode the public key
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return false
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return false
	}

	// Hash the data using SHA-256
	hash := sha256.Sum256(data)

	// Split the signature into r and s
	sigLen := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:sigLen])
	s := new(big.Int).SetBytes(signature[sigLen:])

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], r, s)
}

func HashData(data []byte) []byte {
	// Compute a secure hash for the input data.
	// Use a secure hashing algorithm like SHA-256.
	// Convert the hash output to a hexadecimal string for readability/storage.
	// Return the hash as a string.
	// Use a cryptographic hash function that ensures collision resistance and Do not use weak algorithms like MD5 or SHA-1.
	
	hash := sha256.Sum256(data)
	return hash[:]
}
