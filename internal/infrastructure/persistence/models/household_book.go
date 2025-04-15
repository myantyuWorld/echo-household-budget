package models

// HouseholdBook は家計簿モデル
type HouseholdBook struct {
	Base
	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
}

func (HouseholdBook) TableName() string { return "household_books" }
