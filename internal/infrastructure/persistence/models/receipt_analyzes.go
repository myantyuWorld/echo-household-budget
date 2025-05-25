package models

type ReceiptAnalyzes struct {
	ID              int                   `gorm:"primary_key"`
	ImageURL        string                `gorm:"not null"`
	AnalyzeStatus   string                `gorm:"not null"`
	TotalPrice      int                   `gorm:"not null"`
	HouseholdBookID int                   `gorm:"not null"`
	HouseholdBook   HouseholdBook         `gorm:"foreignKey:HouseholdBookID"`
	Items           []ReceiptAnalyzeItems `gorm:"foreignKey:ReceiptAnalyzeID"`
}

func (ReceiptAnalyzes) TableName() string {
	return "receipt_analyzes"
}
