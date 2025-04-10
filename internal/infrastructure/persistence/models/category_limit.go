package models

// CategoryLimit はカテゴリ予算モデル
type CategoryLimit struct {
	Base
	HouseholdBookID uint `gorm:"not null;index"`
	CategoryID      uint `gorm:"not null;index;uniqueIndex:idx_category_limits_unique"`
	LimitAmount     int  `gorm:"not null;default:0"`
	HouseholdBook   HouseholdBook
	Category        Category
}

func (CategoryLimit) TableName() string { return "category_limits" }
