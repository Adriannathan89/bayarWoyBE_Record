package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendship struct {
	ID       string `gorm:"primaryKey" json:"id"`
	UserID   string `gorm:"column:user_id" json:"userId"`
	FriendID string `gorm:"column:friend_id" json:"friendId"`

	User   User   `gorm:"foreignKey:UserID; reference:UserID" json:"user"`
	Friend User   `gorm:"foreignKey:FriendID; reference:FriendID" json:"friend"`
	Status string `gorm:"column:status" json:"status"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
