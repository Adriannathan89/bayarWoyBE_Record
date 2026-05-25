package responses

type RecordResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	CreatedAt   string  `json:"createdAt"`
}
