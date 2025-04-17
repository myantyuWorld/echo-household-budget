package models

// HouseholdBook は家計簿モデル
type HouseholdBook struct {
	Base
	Title          string          `gorm:"type:varchar(255);not null"`
	Description    string          `gorm:"type:text"`
	CategoryLimits []CategoryLimit `gorm:"foreignKey:HouseholdBookID"`
	Users          []UserAccount   `gorm:"many2many:user_households;foreignKey:ID;joinForeignKey:HouseholdID;References:ID;joinReferences:UserID"`
}

func (HouseholdBook) TableName() string { return "household_books" }
