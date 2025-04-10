package models

// Category はカテゴリモデル
type Category struct {
	Base
	Name  string `gorm:"type:varchar(255);not null;index"`
	Color string `gorm:"type:varchar(7)"`
}

func (Category) TableName() string { return "categories" }
