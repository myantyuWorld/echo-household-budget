package models

type UserHouseHold struct {
	Base
	UserID        uint          `gorm:"not null;index:idx_user_households_user_id"`
	HouseholdID   uint          `gorm:"not null;index:idx_user_households_household_id"`
	UserAccount   UserAccount   `gorm:"foreignKey:UserID;references:ID"`
	HouseholdBook HouseholdBook `gorm:"foreignKey:HouseholdID;references:ID"`
}

func (UserHouseHold) TableName() string { return "user_households" }
