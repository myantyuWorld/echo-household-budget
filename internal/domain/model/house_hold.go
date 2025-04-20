//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

type HouseHold struct {
	ID            HouseHoldID      `json:"id"`
	UserID        UserID           `json:"userID"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	CategoryLimit []*CategoryLimit `json:"categoryLimit"`
}

type HouseHoldID uint

func NewDefaultHouseHold(userAccount *UserAccount) *HouseHold {
	householdBook := &HouseHold{
		UserID:      userAccount.ID,
		Title:       "default",
		Description: "default",
		CategoryLimit: []*CategoryLimit{
			{
				Category: Category{
					ID:    CategoryIDFood,
					Name:  "食費",
					Color: "#FF0000",
				},
				LimitAmount: 40000,
			},
			{
				Category: Category{
					ID:    CategoryIDNecessary,
					Name:  "日用品",
					Color: "#00FF00",
				},
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
