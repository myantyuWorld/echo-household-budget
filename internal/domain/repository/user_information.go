package repository

import domainmodel "echo-household-budget/internal/domain/model"

type UserInformationRepository interface {
	Create(userInformation *domainmodel.UserInformation) error
	UpdateRead(informationIDs []int, userID int) error
	FindAllIsPublished(userID int) ([]*domainmodel.UserInformation, error)
}
