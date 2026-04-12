package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Record struct {
	ID          string `gorm:"column:id; primaryKey" json:"id"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	Amount      int    `gorm:"column:amount" json:"amount"`
	OwnerID     string `gorm:"column:owner_id" json:"ownerId"`
	Owner       User   `gorm:"foreignKey:OwnerID" json:"owner"`
	Type        string `gorm:"column:type" json:"type"`
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
