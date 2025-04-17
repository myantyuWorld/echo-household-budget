package models

// Category はカテゴリモデル
type Category struct {
	Base
	Name  string `gorm:"type:varchar(255);not null"`
	Color string `gorm:"type:varchar(7);not null"`
}

func (Category) TableName() string { return "categories" }
