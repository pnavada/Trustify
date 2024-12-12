package types

type TransactionType string

const (
	TransactionTypeCoinbase TransactionType = "coinbase"
	TransactionTypePurchase TransactionType = "purchase"
	TransactionTypeReview   TransactionType = "review"
)
