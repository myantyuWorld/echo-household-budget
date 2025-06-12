package domainmodel

type UserInformation struct {
	ID            int
	InformationID int
	UserID        UserID
	IsRead        bool
	Information   *Information
}

func NewUserInformation(informationID int, userID UserID) *UserInformation {
	return &UserInformation{
		InformationID: informationID,
		UserID:        userID,
		IsRead:        false,
	}
}
