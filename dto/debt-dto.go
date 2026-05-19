package dto

type DebtDTO struct {
	Amount      float32 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description" binding:"required"`
	DebtorID    string  `json:"debtorId" binding:"omitempty"`
}
