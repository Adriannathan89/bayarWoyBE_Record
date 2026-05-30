package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscordVerification struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id; index" json:"userId"`
	Code      string    `gorm:"column:code; uniqueIndex" json:"code"`
	ExpiresAt time.Time `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (d *DiscordVerification) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
