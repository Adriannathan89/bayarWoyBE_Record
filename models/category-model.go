package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID   string `gorm:"column:id; primaryKey" json:"id"`
	Name string `gorm:"column:name;uniqueIndex:uq_cat_name_type" json:"name"`
	Type string `gorm:"column:type;uniqueIndex:uq_cat_name_type" json:"type"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	return nil
}
