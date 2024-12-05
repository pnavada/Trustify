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
)

// Use the most preferred cryptographic library for generating key pairs and signing data in blockchain applications.
// TO-DO: Use a secure cryptographic library to generate an elliptic curve (or RSA) key pair.
// TO-DO: Ensure private keys are securely stored and not exposed.
// TO-DO: Gracefully handle invalid inputs, such as malformed keys or signatures.

func Sign(data []byte, privateKey []byte) ([]byte, error) {
	// Use the most preferred cryptographic and hashing algorithms for blockchain applications.
	// Create a digital signature for the provided data using the private key.
	// Hash the data to create a digest.
	// Use the private key to sign the hashed data.
	// Return the signature in a format suitable for verification (e.g., DER-encoded).
	// Decode the private key

	// Decode the private key from PEM format
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	var privKey *ecdsa.PrivateKey
	var err error

	if block.Type == "EC PRIVATE KEY" {
		privKey, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else if block.Type == "PRIVATE KEY" {
		keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		var ok bool
		privKey, ok = keyInterface.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("not ECDSA private key")
		}
	} else {
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
	type ecdsaSignature struct {
		R, S *big.Int
	}
	signature, err := asn1.Marshal(ecdsaSignature{r, s})
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func VerifySignature(hash []byte, signature []byte, publicKey []byte) bool {
	//  Verify the digital signature of the data using the public key.
	// Hash the input data to create a digest (same method used in Sign).
	// Use the public key to verify the signature against the hash.
	// Return true if the signature is valid, otherwise false.
	// Ensure the public key is correctly formatted and corresponds to the private key used for signing.

	// Decode the public key
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false
	}

	var pubKeyInterface interface{}
	var err error

	if block.Type == "PUBLIC KEY" {
		pubKeyInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return false
		}
	} else if block.Type == "EC PUBLIC KEY" {
		pubKeyInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return false
		}
	} else {
		return false
	}

	pubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return false
	}

	// Decode the signature from ASN.1 DER
	type ecdsaSignature struct {
		R, S *big.Int
	}
	var sig ecdsaSignature
	_, err = asn1.Unmarshal(signature, &sig)
	if err != nil {
		return false
	}

	// Verify the signature using ECDSA
	return ecdsa.Verify(pubKey, hash[:], sig.R, sig.S)
}
