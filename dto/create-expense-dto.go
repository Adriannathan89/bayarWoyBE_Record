package dto

type CreateExpenseDto struct {
	Description string `json:"description" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
}