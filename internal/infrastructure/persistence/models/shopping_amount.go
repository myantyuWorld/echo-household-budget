package models

import "time"

// ShoppingAmount は買い物金額モデル
type ShoppingAmount struct {
	Base
	HouseholdBookID uint      `gorm:"not null;index"`
	CategoryID      uint      `gorm:"not null;index"`
	Amount          int       `gorm:"not null;default:0"`
	Date            time.Time `gorm:"not null"`
	Memo            string    `gorm:"type:text"`
	HouseholdBook   HouseholdBook
	Category        Category
}

func (ShoppingAmount) TableName() string { return "shopping_amounts" }
