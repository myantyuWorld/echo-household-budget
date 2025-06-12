package repository

import domainmodel "echo-household-budget/internal/domain/model"

type UserInformationRepository interface {
	Create(userInformation *domainmodel.UserInformation) error
}
