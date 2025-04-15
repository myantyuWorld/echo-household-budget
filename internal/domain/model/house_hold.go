//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

type HouseHold struct {
	ID            HouseHoldID
	UserID        UserID
	Title         string
	Description   string
	CategoryLimit []CategoryLimit
}

type HouseHoldID uint

func NewDefaultHouseHold(userAccount *UserAccount) *HouseHold {
	householdBook := &HouseHold{
		UserID:      userAccount.ID,
		Title:       "default",
		Description: "default",
		CategoryLimit: []CategoryLimit{
			{
				CategoryID:  CategoryIDFood,
				LimitAmount: 40000,
			},
			{
				CategoryID:  CategoryIDNecessary,
				LimitAmount: 10000,
			},
		},
	}
	return householdBook
}

type UserHouseHold struct {
	UserID      UserID
	HouseHoldID HouseHoldID
}

type HouseHoldRepository interface {
	Create(houseHold *HouseHold) error
	CreateUserHouseHold(userHouseHold *UserHouseHold) error
	FindByUserID(userID UserID) (*HouseHold, error)
	FindByHouseHoldID(houseHoldID HouseHoldID) (*HouseHold, error)
	Update(houseHold *HouseHold) error
	Delete(houseHoldID HouseHoldID) error
}
