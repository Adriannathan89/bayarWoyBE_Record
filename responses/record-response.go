package responses

type CategoryInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RecordResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Amount      float32        `json:"amount"`
	Categories  []CategoryInfo `json:"categories"`
	Type        string         `json:"type"`
	CreatedAt   string         `json:"createdAt"`
	IsCommitted bool           `json:"isCommitted"`
}
