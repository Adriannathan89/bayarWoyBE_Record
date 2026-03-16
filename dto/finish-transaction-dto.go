package dto

type FinishTransactionDTO struct {
	TransactionID string `json:"transactionId" binding:"required"`
}