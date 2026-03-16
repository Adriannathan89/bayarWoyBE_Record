package botmodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscordBotOtp struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	OTP       string    `gorm:"column:otp" json:"otp"`
	ExpiredAt time.Time `gorm:"column:expired_at" json:"expired_at"`
}

func (d *DiscordBotOtp) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
