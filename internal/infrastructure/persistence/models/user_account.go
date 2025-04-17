package models

// UserAccount はユーザーアカウントモデル
type UserAccount struct {
	Base
	UserID         string          `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name           string          `gorm:"type:varchar(255);not null"`
	PictureURL     string          `gorm:"type:varchar(255);not null"`
	HouseholdBooks []HouseholdBook `gorm:"many2many:user_households;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:HouseholdID"`
}

func (UserAccount) TableName() string { return "user_accounts" }
