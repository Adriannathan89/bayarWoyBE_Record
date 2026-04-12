package dto

type DebtDTO struct {
	Amount      int    `json:"amount" binding:"required,gt=0"`
	Description string `json:"description" binding:"required"`
	DebtorID    string `json:"debtorId" binding:"omitempty"`
}