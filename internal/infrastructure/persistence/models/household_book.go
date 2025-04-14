package models

// HouseholdBook は家計簿モデル
type HouseholdBook struct {
	Base
	UserID      string `gorm:"type:varchar(255);not null;index"`
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
}

func (HouseholdBook) TableName() string { return "household_books" }
