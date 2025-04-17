package models

// CategoryLimit はカテゴリ予算モデル
type CategoryLimit struct {
	Base
	HouseholdBookID uint          `gorm:"not null"`
	CategoryID      uint          `gorm:"not null"`
	LimitAmount     int           `gorm:"not null"`
	HouseholdBook   HouseholdBook `gorm:"foreignKey:HouseholdBookID"`
	Category        Category      `gorm:"foreignKey:CategoryID"`
}

func (CategoryLimit) TableName() string { return "category_limits" }
