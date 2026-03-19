package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string        `gorm:"column:id primaryKey" json:"id"`
	Username     string        `gorm:"column:username; unique" json:"username"`
	Password     string        `gorm:"column:password" json:"password"`
	Debt         int           `gorm:"column:debt" json:"debt"`
	Receivable   int           `gorm:"column:receivable" json:"receivable"`
	Cash         int           `gorm:"column:cash" json:"cash"`
	DiscordID    *string       `gorm:"column:discord_id" json:"discordId"`
	IsValidated  bool          `gorm:"column:is_validated; default:false" json:"isValidated"`
	Transactions []Transaction `gorm:"foreignKey:OwnerID" json:"transactions"`
	Expenses     []Expense     `gorm:"foreignKey:OwnerID" json:"expenses"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}
