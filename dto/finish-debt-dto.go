package dto

type FinishDebtDTO struct {
	DebtID string `json:"debtId" binding:"required"`
}