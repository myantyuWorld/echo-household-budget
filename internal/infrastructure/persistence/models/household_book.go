package models

// HouseholdBook は家計簿モデル
type HouseholdBook struct {
	Base
	UserID      string         `gorm:"type:varchar(255);not null;index"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Categories  []Category     `gorm:"foreignKey:HouseholdBookID"`
	Memos       []ShoppingMemo `gorm:"foreignKey:HouseholdBookID"`
}

func (HouseholdBook) TableName() string { return "household_books" }
