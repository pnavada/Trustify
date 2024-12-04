package crypto

import (
	"trustify/blockchain"
)

// Refer https://pkg.go.dev/github.com/wealdtech/go-merkletree#readme-maintainers and use the same for implementation under the hood

// TO-DO: Return meaningful errors for invalid inputs, such as empty transactions or malformed proofs.

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  []byte
}

func BuildTree(transactions []blockchain.UTXOTransaction) (*MerkleTree, error) {
	// Construct a Merkle Tree from a list of transactions.
	// Compute the Merkle Root, representing the cryptographic hash of all transactions.
	// Hash each transaction in the provided transactions list to generate the leaf nodes.
	// Pair up the leaf nodes and hash their concatenated values to create parent nodes.
	// Repeat this process until only one root node remains.
	// If the number of nodes in a level is odd, duplicate the last node to form a pair.
	// Return the constructed MerkleTree object with the root node.

	return nil, nil
}

func (mt *MerkleTree) GetRoot() []byte {
	// Retrieve the Merkle Root of the tree.
	// Return the Hash value of the root node (mt.Root).
	// If the tree is empty (mt.Root == nil), return a nil value.
	// Ensure the tree has been constructed before accessing the root.
	return nil
}

func (mt *MerkleTree) VerifyTransaction(tx blockchain.UTXOTransaction, proof [][]byte) bool {
	// Verify that a transaction exists in the Merkle Tree using a proof.
	// Hash the provided transaction using the same algorithm used for tree construction.
	// Iterate through the proof, hashing the current hash with each proof node’s hash.
	// If the final computed hash matches the Merkle Root, the transaction is verified.
	// Return true if the transaction is valid; otherwise, return false.
	return false
}
