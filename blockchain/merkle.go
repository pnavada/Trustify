package blockchain

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

func BuildTree(transactions []UTXOTransaction) (*MerkleTree, error) {
	// Construct a Merkle Tree from a list of transactions.
	// Compute the Merkle Root, representing the cryptographic hash of all transactions.
	// Hash each transaction in the provided transactions list to generate the leaf nodes.
	// Pair up the leaf nodes and hash their concatenated values to create parent nodes.
	// Repeat this process until only one root node remains.
	// If the number of nodes in a level is odd, duplicate the last node to form a pair.
	// Return the constructed MerkleTree object with the root node.

    if len(transactions) == 0 {
        logger.ErrorLogger.Println("No transactions provided to build the Merkle tree")
        return nil, errors.New("cannot build Merkle tree with zero transactions")
    }

    // Hash each transaction to create leaf nodes
    var leaves []*MerkleNode
    for _, tx := range transactions {
        hash := tx.Hash()
        node := &MerkleNode{
            Left:  nil,
            Right: nil,
            Hash:  hash,
        }
        leaves = append(leaves, node)
    }

    // Build the tree
    root := buildMerkleTree(leaves)
    tree := &MerkleTree{Root: root}
    logger.InfoLogger.Println("Merkle tree built successfully")
    return tree, nil
}

// Recursive function to build the Merkle tree
func buildMerkleTree(nodes []*MerkleNode) *MerkleNode {
    if len(nodes) == 1 {
        return nodes[0]
    }

    var parentLevel []*MerkleNode

    for i := 0; i < len(nodes); i += 2 {
        var left = nodes[i]
        var right *MerkleNode

        if i+1 < len(nodes) {
            right = nodes[i+1]
        } else {
            // Duplicate the last node if the number is odd
            right = nodes[i]
        }

        parentHash := sha256.Sum256(append(left.Hash, right.Hash...))
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

func (mt *MerkleTree) VerifyTransaction(tx UTXOTransaction, proof [][]byte) bool {
	// Verify that a transaction exists in the Merkle Tree using a proof.
	// Hash the provided transaction using the same algorithm used for tree construction.
	// Iterate through the proof, hashing the current hash with each proof node’s hash.
	// If the final computed hash matches the Merkle Root, the transaction is verified.
	// Return true if the transaction is valid; otherwise, return false.
	
	if mt.Root == nil {
        logger.ErrorLogger.Println("Cannot verify transaction in an empty Merkle tree")
        return false
    }

    currentHash := tx.Hash()

    for _, siblingHash := range proof {
        combined := append(currentHash, siblingHash...)
        newHash := sha256.Sum256(combined)
        currentHash = newHash[:]
    }

    isValid := bytes.Equal(currentHash, mt.Root.Hash)
    if isValid {
        logger.InfoLogger.Println("Transaction verified in Merkle tree")
    } else {
        logger.ErrorLogger.Println("Transaction verification failed in Merkle tree")
    }
    return isValid
}
