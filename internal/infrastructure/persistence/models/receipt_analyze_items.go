package models

type ReceiptAnalyzeItems struct {
	ID    int    `gorm:"primary_key"`
	Name  string `gorm:"not null"`
	Price int    `gorm:"not null"`
}

func (ReceiptAnalyzeItems) TableName() string {
	return "receipt_analyze_items"
}
