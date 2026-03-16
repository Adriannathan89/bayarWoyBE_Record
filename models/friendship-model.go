package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendship struct {
	ID       string `gorm:"primaryKey" json:"id"`
	UserID   string `gorm:"column:user_id" json:"userId"`
	FriendID string `gorm:"column:friend_id" json:"friendId"`

	User   User `gorm:"foreignKey:UserID; reference:ID" json:"user"`
	Friend User `gorm:"foreignKey:FriendID; reference:ID" json:"friend"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
