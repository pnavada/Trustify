package blockchain

import "errors"

var (
	ErrEmptyTransactions   = errors.New("block must contain at least one transaction")
	ErrInvalidPreviousHash = errors.New("invalid previous hash")
	ErrInvalidTargetHash   = errors.New("invalid target hash")
	ErrInvalidMerkleRoot   = errors.New("invalid Merkle root")
	ErrInvalidBlockHash    = errors.New("invalid block hash")
	ErrInvalidTimestamp    = errors.New("invalid timestamp")
	ErrInvalidNonce        = errors.New("invalid nonce")
	ErrBlockNotFound       = errors.New("block not found")
	ErrTransactionInvalid  = errors.New("transaction invalid")
	ErrDoubleSpending      = errors.New("double spending detected")
	ErrReviewNotPurchased  = errors.New("reviewer has not purchased the product")
	ErrReviewDuplicate     = errors.New("duplicate review submission")
	ErrInvalidSignature    = errors.New("invalid digital signature")
	ErrUTXONotFound        = errors.New("UTXO not found")
)
