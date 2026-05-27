package dto

type CommitRecordDto struct {
	RecordID string `json:"recordId" binding:"required"`
	Category string `json:"category"` // optional — only if user corrects
}
