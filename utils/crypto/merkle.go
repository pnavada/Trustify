package crypto

// Refer https://pkg.go.dev/github.com/wealdtech/go-merkletree#readme-maintainers

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

func BuildTree(transactions []UXTOTransaction) (*MerkleTree, error) {
	// Method implementation goes here
	return nil, nil
}

func (mt *MerkleTree) GetRoot() []byte {
	// Method implementation goes here
	return nil
}

func (mt *MerkleTree) VerifyTransaction(tx Transaction, proof [][]byte) bool {
	// Method implementation goes here
	return false
}
