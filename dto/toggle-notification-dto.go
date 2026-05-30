package dto

type ToggleNotificationDto struct {
	Type    string `json:"type" binding:"required,oneof=commit weekly"`
	Enabled bool   `json:"enabled"`
}
