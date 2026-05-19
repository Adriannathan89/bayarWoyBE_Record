package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string   `gorm:"column:id; primaryKey" json:"id"`
	Username    string   `gorm:"column:username; unique" json:"username"`
	Password    string   `gorm:"column:password" json:"password"`
	Debt        float32  `gorm:"column:debt" json:"debt"`
	Receivable  float32  `gorm:"column:receivable" json:"receivable"`
	Cash        float32  `gorm:"column:cash" json:"cash"`
	DiscordID   *string  `gorm:"column:discord_id" json:"discordId"`
	IsValidated bool     `gorm:"column:is_validated; default:false" json:"isValidated"`
	Debts       []Debt   `gorm:"foreignKey:OwnerID" json:"debtsx"`
	Records     []Record `gorm:"foreignKey:OwnerID" json:"records"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	return nil
}
