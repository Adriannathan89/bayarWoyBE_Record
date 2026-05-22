package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendRequest struct {
	ID         string `gorm:"primaryKey" json:"id"`
	SenderID   string `gorm:"column:user_id" json:"senderId"`
	ReceiverID string `gorm:"column:friend_id" json:"receiverId"`

	Sender   User `gorm:"foreignKey:SenderID;references:ID" json:"sender"`
	Receiver User `gorm:"foreignKey:ReceiverID;references:ID" json:"receiver"`
}

func (f *FriendRequest) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
