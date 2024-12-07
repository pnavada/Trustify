package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"math/big"
	"trustify/logger"

	"github.com/btcsuite/btcd/btcec/v2"
)

// Use the most preferred cryptographic library for generating key pairs and signing data in blockchain applications.
// TO-DO: Use a secure cryptographic library to generate an elliptic curve (or RSA) key pair.
// TO-DO: Ensure private keys are securely stored and not exposed.
// TO-DO: Gracefully handle invalid inputs, such as malformed keys or signatures.

// ecdsaSignature represents the structure of an ECDSA signature in ASN.1 DER format.
type ecdsaSignature struct {
	R, S *big.Int
}

func Sign(data []byte, privateKeyPEM []byte) ([]byte, error) {
	// Log the private key and data
	logger.InfoLogger.Printf("Private Key: %x\n", privateKeyPEM)
	logger.InfoLogger.Printf("Data: %x\n", data)

	// Decode the private key from PEM format
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
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

	// Hash the data using SHA-256
	hash := sha256.Sum256(data)

	// Sign the hashed data using ECDSA
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}

	// Encode the signature as ASN.1 DER
	signature, err := asn1.Marshal(ecdsaSignature{R: r, S: s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func VerifySignature(hash []byte, signature []byte, publicKeyBytes []byte) bool {
	// Log the public key, hash, and signature in hexadecimal format
	logger.InfoLogger.Printf("Public Key: %x\n", publicKeyBytes)
	logger.InfoLogger.Printf("Hash: %x\n", hash)
	logger.InfoLogger.Printf("Signature: %x\n", signature)

	// Parse the compressed public key
	pubKey, err := parseCompressedPublicKey(publicKeyBytes)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to parse public key: %v\n", err)
		return false
	}

	// Unmarshal the signature from ASN.1 DER format to obtain R and S values
	var sig ecdsaSignature
	_, err = asn1.Unmarshal(signature, &sig)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal signature: %v\n", err)
		return false
	}

	// Verify the signature using the ecdsa.Verify function
	valid := ecdsa.Verify(pubKey, hash, sig.R, sig.S)
	if !valid {
		logger.ErrorLogger.Println("Signature verification failed")
	} else {
		logger.InfoLogger.Println("Signature verification succeeded")
	}

	return valid
}

// Helper function to parse a compressed public key into ecdsa.PublicKey
func parseCompressedPublicKey(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
	if len(pubKeyBytes) != 33 {
		return nil, errors.New("invalid compressed public key length")
	}

	// Determine the Y coordinate based on the prefix byte
	prefix := pubKeyBytes[0]
	if prefix != 0x02 && prefix != 0x03 {
		return nil, errors.New("invalid compressed public key prefix")
	}

	x := new(big.Int).SetBytes(pubKeyBytes[1:])
	curve := btcec.S256() // Replaceed with secp256k1

	// Compute Y coordinate
	y, err := decompressY(x, prefix)
	if err != nil {
		return nil, err
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// Helper function to decompress Y coordinate
func decompressY(x *big.Int, prefix byte) (*big.Int, error) {
	// This function needs to implement the Y coordinate recovery based on X and prefix
	// For secp256k1, you might use btcec or another library
	// Here's a simplified placeholder

	// Import btcec for Y coordinate decompression
	// Alternatively, use any other library that supports secp256k1
	privKey, err := btcec.ParsePubKey(append([]byte{prefix}, x.Bytes()...))
	if err != nil {
		return nil, err
	}
	return privKey.Y(), nil
}
