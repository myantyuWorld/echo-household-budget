package models

type Information struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Content     string `gorm:"not null"`
	Category    string `gorm:"not null"`
	IsPublished bool   `gorm:"not null"`
}

func (Information) TableName() string {
	return "informations"
}
