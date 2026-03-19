package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          string `gorm:"column:id primaryKey" json:"id"`
	Description string `gorm:"column:description" json:"description"`
	Amount      int    `gorm:"column:amount" json:"amount"`
	OwnerID     string `gorm:"column:owner_id" json:"ownerId"`
	Owner       User   `gorm:"foreignKey:OwnerID" json:"owner"`
}

func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	return nil
}