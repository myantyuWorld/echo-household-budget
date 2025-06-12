package models

type UserInformation struct {
	ID            int  `gorm:"primaryKey;autoIncrement"`
	InformationID int  `gorm:"not null"`
	UserID        int  `gorm:"not null"`
	IsRead        bool `gorm:"not null"`
	Information   Information
}
