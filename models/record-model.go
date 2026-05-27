package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Record struct {
	ID          string     `gorm:"column:id; primaryKey" json:"id"`
	Title       string     `gorm:"column:title" json:"title"`
	Description string     `gorm:"column:description; nullable" json:"description"`
	Amount      float32    `gorm:"column:amount" json:"amount"`
	OwnerID     string     `gorm:"column:owner_id" json:"ownerId"`
	Owner       User       `gorm:"foreignKey:OwnerID" json:"owner"`
	Type        string     `gorm:"column:type" json:"type"`
	Categories  []Category `gorm:"many2many:record_categories;" json:"categories"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	IsCommitted bool       `gorm:"column:is_committed;default:false" json:"isCommitted"`
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}
