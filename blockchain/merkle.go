package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"trustify/logger"
)

// import ()

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

func BuildTree(transactions []*Transaction) (*MerkleTree, error) {
	// Construct a Merkle Tree from a list of transactions.
	// Compute the Merkle Root, representing the cryptographic hash of all transactions.
	// Hash each transaction in the provided transactions list to generate the leaf nodes.
	// Pair up the leaf nodes and hash their concatenated values to create parent nodes.
	// Repeat this process until only one root node remains.
	// If the number of nodes in a level is odd, duplicate the last node to form a pair.
	// Return the constructed MerkleTree object with the root node.

	// Validate the input
	if len(transactions) == 0 {
		logger.ErrorLogger.Println("No transactions provided to build the Merkle tree")
		return nil, errors.New("cannot build Merkle tree with zero transactions")
	}

	// Sort transactions deterministically by their IDs (hashes)
	sort.Slice(transactions, func(i, j int) bool {
		return bytes.Compare(transactions[i].ID, transactions[j].ID) < 0
	})

	// Hash each transaction to create leaf nodes
	var leaves []*MerkleNode
	for i, tx := range transactions {
		serializedTx := SerializeTransaction(tx)
		hash := sha256.Sum256(serializedTx)
		if len(hash) == 0 {
			logger.ErrorLogger.Printf("Failed to hash transaction at index %d\n", i)
			return nil, errors.New("failed to hash transaction")
		}

		node := &MerkleNode{
			Left:  nil,
			Right: nil,
			Hash:  hash[:],
		}
		leaves = append(leaves, node)
	}

	// Build the Merkle Tree
	root, err := buildMerkleTree(leaves)
	if err != nil {
		logger.ErrorLogger.Println("Failed to build Merkle tree:", err)
		return nil, err
	}

	// Construct the Merkle Tree object
	tree := &MerkleTree{Root: root}
	logger.InfoLogger.Println("Merkle tree built successfully with root hash:", hex.EncodeToString(root.Hash))
	return tree, nil
}

func buildMerkleTree(nodes []*MerkleNode) (*MerkleNode, error) {
	// Ensure there are nodes to process
	if len(nodes) == 0 {
		logger.ErrorLogger.Println("Cannot build Merkle tree with no nodes")
		return nil, errors.New("no nodes provided to build the Merkle tree")
	}

	// If there's only one node, it's the root of the Merkle Tree
	if len(nodes) == 1 {
		return nodes[0], nil
	}

	var parentLevel []*MerkleNode

	// Iterate through the current level of nodes
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		var right *MerkleNode

		// Handle odd number of nodes by duplicating the last node
		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = nodes[i]
		}

		// Compute the parent hash
		parentHash := sha256.Sum256(append(left.Hash, right.Hash...))

		// Create the parent node
		parentNode := &MerkleNode{
			Left:  left,
			Right: right,
			Hash:  parentHash[:],
		}
		parentLevel = append(parentLevel, parentNode)
	}

	// Recursively build the next level
	return buildMerkleTree(parentLevel)
}

func (mt *MerkleTree) GetRoot() []byte {
	// Retrieve the Merkle Root of the tree.
	// Return the Hash value of the root node (mt.Root).
	// If the tree is empty (mt.Root == nil), return a nil value.
	// Ensure the tree has been constructed before accessing the root.
	if mt.Root == nil {
		logger.ErrorLogger.Println("Merkle tree root is nil")
		return nil
	}
	return mt.Root.Hash
}
