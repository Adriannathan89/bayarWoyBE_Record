package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Debt struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	OwnerID     string    `gorm:"column:owner_id" json:"ownerId"`
	Amount      float32   `gorm:"column:amount" json:"amount"`
	Description string    `gorm:"column:description" json:"description"`
	Category    string    `gorm:"column:category" json:"category"`
	Owner       User      `gorm:"foreignKey:OwnerID; reference:ID" json:"owner"`
	DebtorID    string    `gorm:"column:debtor_id" json:"debtorId"`
	Debtor      User      `gorm:"foreignKey:DebtorID; reference:ID" json:"debtor"`
	Status      string    `gorm:"column:status" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
}

func (d *Debt) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	return nil
}
