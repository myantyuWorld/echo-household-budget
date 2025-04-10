package models

// Category はカテゴリモデル
type Category struct {
	Base
	HouseholdBookID uint   `gorm:"not null;index"`
	Name            string `gorm:"type:varchar(255);not null;index"`
	Color           string `gorm:"type:varchar(7)"`
	HouseholdBook   HouseholdBook
	ShoppingMemos   []ShoppingMemo
	CategoryLimit   *CategoryLimit
}

func (Category) TableName() string { return "categories" }
