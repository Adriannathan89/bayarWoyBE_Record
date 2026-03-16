package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendRequest struct {
	ID               string `gorm:"primaryKey" json:"id"`
	SenderUsername   string `gorm:"column:user_username" json:"senderUsername"`
	ReceiverUsername string `gorm:"column:friend_username" json:"receiverUsername"`

	Sender   User   `gorm:"foreignKey:SenderUsername;references:Username" json:"sender"`
	Receiver User   `gorm:"foreignKey:ReceiverUsername;references:Username" json:"receiver"`
}

func (f *FriendRequest) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
