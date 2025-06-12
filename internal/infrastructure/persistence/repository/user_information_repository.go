package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/domain/repository"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type UserInformationRepository struct {
	db *gorm.DB
}

// Create implements repository.UserInformationRepository.
func (u *UserInformationRepository) Create(userInformation *domainmodel.UserInformation) error {
	model := &models.UserInformation{
		InformationID: userInformation.InformationID,
		UserID:        int(userInformation.UserID),
		IsRead:        userInformation.IsRead,
	}

	if err := u.db.Create(model).Error; err != nil {
		return err
	}

	userInformation.ID = model.ID

	return nil
}

func NewUserInformationRepository(db *gorm.DB) repository.UserInformationRepository {
	return &UserInformationRepository{
		db: db,
	}
}
