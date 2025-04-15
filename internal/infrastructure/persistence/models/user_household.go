package models

type UserHouseHold struct {
	Base
	UserID      uint `gorm:"not null"`
	HouseholdID uint `gorm:"not null"`
}

func (UserHouseHold) TableName() string { return "user_households" }
