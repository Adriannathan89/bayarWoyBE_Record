package dto

type CreateRecordDto struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=income expense"`
}
