package dto

type CreateRecordDto struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Amount      float32 `json:"amount" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
}
