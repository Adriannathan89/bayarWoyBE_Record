package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"column:username" json:"username"`
	UserID       string    `gorm:"column:user_id" json:"userId"`
	IPAddress    string    `gorm:"column:ip_address" json:"ipAddress"`
	RefreshToken string    `gorm:"column:refresh_token" json:"refreshToken"`
	ExpiresAt    time.Time `gorm:"column:expires_at" json:"expiresAt"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	return nil
}