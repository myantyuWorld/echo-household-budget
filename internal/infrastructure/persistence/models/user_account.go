package models

// UserAccount はユーザーアカウントモデル
type UserAccount struct {
	Base
	UserID string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name   string `gorm:"type:varchar(255);not null"`
}

func (UserAccount) TableName() string { return "user_accounts" }
