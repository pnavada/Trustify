package blockchain

import (
	"crypto/sha256"
	"errors"
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

	// Hash each transaction to create leaf nodes
	var leaves []*MerkleNode
	for i, tx := range transactions {
		hash := sha256.Sum256(SerializeTransaction(tx))
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
	logger.InfoLogger.Println("Merkle tree built successfully with root hash:", root.Hash)
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
		var left = nodes[i]
		var right *MerkleNode

		// Ensure the left node has a valid hash
		if len(left.Hash) == 0 {
			logger.ErrorLogger.Printf("Invalid hash for left node at index %d\n", i)
			return nil, errors.New("invalid hash in Merkle tree node")
		}

		// Handle odd number of nodes by duplicating the last node
		if i+1 < len(nodes) {
			right = nodes[i+1]
		} else {
			right = nodes[i]
		}

		// Ensure the right node has a valid hash
		if len(right.Hash) == 0 {
			logger.ErrorLogger.Printf("Invalid hash for right node at index %d\n", i+1)
			return nil, errors.New("invalid hash in Merkle tree node")
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
	root, err := buildMerkleTree(parentLevel)
	if err != nil {
		logger.ErrorLogger.Println("Failed to build next level of Merkle tree:", err)
		return nil, err
	}

	return root, nil
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

func (mt *MerkleTree) VerifyTransaction(tx UTXOTransaction, proof [][]byte) bool {
	// Verify that a transaction exists in the Merkle Tree using a proof.
	// Hash the provided transaction using the same algorithm used for tree construction.
	// Iterate through the proof, hashing the current hash with each proof node’s hash.
	// If the final computed hash matches the Merkle Root, the transaction is verified.
	// Return true if the transaction is valid; otherwise, return false.

	// if mt.Root == nil {
	//     logger.ErrorLogger.Println("Cannot verify transaction in an empty Merkle tree")
	//     return false
	// }

	// currentHash := tx.Hash()

	// for _, siblingHash := range proof {
	//     combined := append(currentHash, siblingHash...)
	//     newHash := sha256.Sum256(combined)
	//     currentHash = newHash[:]
	// }

	// isValid := bytes.Equal(currentHash, mt.Root.Hash)
	// if isValid {
	//     logger.InfoLogger.Println("Transaction verified in Merkle tree")
	// } else {
	//     logger.ErrorLogger.Println("Transaction verification failed in Merkle tree")
	// }
	// return isValid

	return false
}
