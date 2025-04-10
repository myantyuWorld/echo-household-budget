package models

// UserAccount はユーザーアカウントモデル
type UserAccount struct {
	Base
	UserID string `gorm:"type:varchar(255);not null;index"`
	Name   string `gorm:"type:varchar(255);not null"`
}

func (UserAccount) TableName() string { return "user_accounts" }
