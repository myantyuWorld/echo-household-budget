package repository

import domainmodel "echo-household-budget/internal/domain/model"

type InformationRepository interface {
	Create(information *domainmodel.Information) error
	UpdateStatusPublished(informationID int) error
	FindAll() ([]*domainmodel.Information, error)
}
