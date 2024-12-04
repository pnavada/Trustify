package crypto

import (
	"trustify/blockchain"
)

// Refer https://pkg.go.dev/github.com/wealdtech/go-merkletree#readme-maintainers

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

func BuildTree(transactions []blockchain.UTXOTransaction) (*MerkleTree, error) {
	// Method implementation goes here
	return nil, nil
}

func (mt *MerkleTree) GetRoot() []byte {
	// Method implementation goes here
	return nil
}

func (mt *MerkleTree) VerifyTransaction(tx blockchain.UTXOTransaction, proof [][]byte) bool {
	// Method implementation goes here
	return false
}
