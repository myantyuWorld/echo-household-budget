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

// FindAllIsPublished implements repository.UserInformationRepository.
func (u *UserInformationRepository) FindAllIsPublished(userID int) ([]*domainmodel.UserInformation, error) {
	userInformationModels := make([]*models.UserInformation, 0)
	if err := u.db.Where("user_id = ?", userID).Preload("Information").Order("created_at DESC").Find(&userInformationModels).Error; err != nil {
		return nil, err
	}

	userInformation := make([]*domainmodel.UserInformation, 0, len(userInformationModels))
	for _, userInformationModel := range userInformationModels {
		userInformation = append(userInformation, &domainmodel.UserInformation{
			InformationID: userInformationModel.InformationID,
			UserID:        domainmodel.UserID(userInformationModel.UserID),
			IsRead:        userInformationModel.IsRead,
			Information: &domainmodel.Information{
				ID:       userInformationModel.InformationID,
				Title:    userInformationModel.Information.Title,
				Content:  userInformationModel.Information.Content,
				Category: userInformationModel.Information.Category,
			},
		})
	}

	return userInformation, nil
}

// UpdateRead implements repository.UserInformationRepository.
func (u *UserInformationRepository) UpdateRead(informationIDs []int, userID int) error {
	userInformationModels := make([]*models.UserInformation, 0, len(informationIDs))
	for _, informationID := range informationIDs {
		userInformationModels = append(userInformationModels, &models.UserInformation{
			InformationID: informationID,
			UserID:        userID,
			IsRead:        true,
		})
	}

	if err := u.db.Model(&models.UserInformation{}).Where("user_id = ? AND information_id IN ?", userID, informationIDs).Update("is_read", true).Error; err != nil {
		return err
	}

	return nil
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
