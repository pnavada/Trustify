package cryptography

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
)

// parsePrivateKey decodes the PEM-encoded private key from the configuration.
func ParsePrivateKey(pemKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	return privateKey, nil
}

// serializePublicKey serializes the ECDSA public key to a byte slice.
func SerializePublicKey(pubKey *ecdsa.PublicKey) ([]byte, error) {
	derPubKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	return derPubKey, nil
}

// deserializePublicKey deserializes the ECDSA public key from a byte slice.
func DeserializePublicKey(derPubKey []byte) (*ecdsa.PublicKey, error) {
	pubKeyInterface, err := x509.ParsePKIXPublicKey(derPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER public key: %w", err)
	}

	pubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not ECDSA public key")
	}

	return pubKey, nil
}

func SignMessage(privateKey *ecdsa.PrivateKey, hashedMessage []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashedMessage[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Serialize the signature as r || s
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// verifySignature verifies the ECDSA signature.
func VerifySignature(pubKey *ecdsa.PublicKey, hashedMessage []byte, signature []byte) bool {
	// Split the signature into r and s values
	keyLen := (pubKey.Curve.Params().BitSize + 7) / 8
	r := new(big.Int).SetBytes(signature[:keyLen])
	s := new(big.Int).SetBytes(signature[keyLen:])

	// Verify the signature
	return ecdsa.Verify(pubKey, hashedMessage[:], r, s)
}
