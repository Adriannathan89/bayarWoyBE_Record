package dto

type CreateRecordDto struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount" binding:"required"`
	Date        string  `json:"date" binding:"required,datetime=2006-01-02"`
}
