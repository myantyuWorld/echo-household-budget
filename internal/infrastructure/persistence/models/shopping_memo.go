package models

// ShoppingMemo は買い物メモモデル
type ShoppingMemo struct {
	Base
	HouseholdBookID uint   `gorm:"not null;index"`
	CategoryID      *uint  `gorm:"index"`
	Title           string `gorm:"type:varchar(255);not null"`
	Memo            string `gorm:"type:text"`
	IsCompleted     bool   `gorm:"not null;default:false;index"`
	Category        *Category
}

func (ShoppingMemo) TableName() string { return "shopping_memos" }
