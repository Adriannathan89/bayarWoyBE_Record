package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friendship struct {
	ID       string `gorm:"primaryKey" json:"id"`
	UserID   string `gorm:"column:user_id;uniqueIndex:unique_friendship,composite:1" json:"userId"`
	FriendID string `gorm:"column:friend_id;uniqueIndex:unique_friendship,composite:2" json:"friendId"`
	User     User   `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Friend   User   `gorm:"foreignKey:FriendID;references:ID" json:"friend"`
	Status   string `gorm:"column:status" json:"status"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	return nil
}
