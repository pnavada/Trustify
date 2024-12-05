package blockchain

import "errors"

var (
	ErrEmptyTransactions      = errors.New("block must contain at least one transaction")
	ErrInvalidPreviousHash    = errors.New("invalid previous hash")
	ErrInvalidTargetHash      = errors.New("invalid target hash")
	ErrInvalidMerkleRoot      = errors.New("invalid Merkle root")
	ErrInvalidBlockHash       = errors.New("invalid block hash")
	ErrInvalidTimestamp       = errors.New("invalid timestamp")
	ErrInvalidNonce           = errors.New("invalid nonce")
	ErrBlockNotFound          = errors.New("block not found")
	ErrTransactionInvalid     = errors.New("transaction invalid")
	ErrReviewNotPurchased     = errors.New("reviewer has not purchased the product")
	ErrReviewDuplicate        = errors.New("duplicate review submission")
	ErrInsufficientFunds      = errors.New("insufficient funds for transaction")
	ErrUTXONotFound           = errors.New("input UTXO not found in UTXO set")
	ErrDoubleSpending         = errors.New("double spending detected")
	ErrInvalidSignature       = errors.New("invalid digital signature")
	ErrProductNotPurchased    = errors.New("product not purchased, cannot submit review")
	ErrDuplicateReview        = errors.New("duplicate review for the product by the same reviewer")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInputOutputMismatch    = errors.New("input and output sum mismatch in transaction")
	ErrInvalidProofOfWork     = errors.New("invalid proof of work")
)
