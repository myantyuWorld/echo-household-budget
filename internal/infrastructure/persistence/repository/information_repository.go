package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type informationRepository struct {
	db *gorm.DB
}

// FindAll implements repository.InformationRepository.
func (i *informationRepository) FindAll() ([]*domainmodel.Information, error) {
	models := []*models.Information{}
	if err := i.db.Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}

	output := make([]*domainmodel.Information, len(models))
	for i, model := range models {
		output[i] = &domainmodel.Information{
			ID:          model.ID,
			Title:       model.Title,
			Content:     model.Content,
			Category:    model.Category,
			IsPublished: model.IsPublished,
		}
	}

	return output, nil
}

// Create implements repository.InformationRepository.
func (i *informationRepository) Create(information *domainmodel.Information) error {
	model := &models.Information{
		Title:       information.Title,
		Content:     information.Content,
		Category:    information.Category,
		IsPublished: information.IsPublished,
	}

	if err := i.db.Create(model).Error; err != nil {
		return err
	}

	information.ID = int(model.ID)

	return nil
}

// UpdateStatusPublished implements repository.InformationRepository.
func (i *informationRepository) UpdateStatusPublished(informationID int) error {
	if err := i.db.Model(&models.Information{}).Where("id = ?", informationID).Update("is_published", true).Error; err != nil {
		return err
	}

	return nil
}

func NewInformationRepository(db *gorm.DB) repository.InformationRepository {
	return &informationRepository{
		db: db,
	}
}
