package models

type ReceiptAnalyzeItems struct {
	ID               int             `gorm:"primary_key"`
	ReceiptAnalyzeID int             `gorm:"not null"`
	ReceiptAnalyze   ReceiptAnalyzes `gorm:"foreignKey:ReceiptAnalyzeID"`
	Name             string          `gorm:"not null"`
	Price            int             `gorm:"not null"`
}

func (ReceiptAnalyzeItems) TableName() string {
	return "receipt_analyze_items"
}
