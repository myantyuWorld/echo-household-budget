package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type HouseHoldRepository struct {
	db *gorm.DB
}

// Create implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) Create(houseHold *domainmodel.HouseHold) error {
	model := &models.HouseholdBook{
		Title:       houseHold.Title,
		Description: houseHold.Description,
	}

	if err := h.db.Create(model).Error; err != nil {
		return err
	}

	houseHold.ID = domainmodel.HouseHoldID(model.ID)

	return nil
}

// CreateUserHouseHold implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) CreateUserHouseHold(userHouseHold *domainmodel.UserHouseHold) error {
	model := &models.UserHouseHold{
		UserID:      uint(userHouseHold.UserID),
		HouseholdID: uint(userHouseHold.HouseHoldID),
	}

	if err := h.db.Create(model).Error; err != nil {
		return err
	}

	userHouseHold.HouseHoldID = domainmodel.HouseHoldID(model.HouseholdID)

	return nil
}

// Delete implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) Delete(houseHoldID domainmodel.HouseHoldID) error {
	panic("unimplemented")
}

// FindByHouseHoldID implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) FindByHouseHoldID(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error) {
	model := &models.HouseholdBook{}
	if err := h.db.Where("id = ?", houseHoldID).
		Preload("CategoryLimits").
		Preload("CategoryLimits.Category").
		First(model).Error; err != nil {
		return nil, err
	}

	categoryLimits := make([]*domainmodel.CategoryLimit, len(model.CategoryLimits))
	for i, categoryLimit := range model.CategoryLimits {
		categoryLimits[i] = &domainmodel.CategoryLimit{
			Category: domainmodel.Category{
				ID:    domainmodel.CategoryID(categoryLimit.Category.ID),
				Name:  categoryLimit.Category.Name,
				Color: categoryLimit.Category.Color,
			},
			LimitAmount: categoryLimit.LimitAmount,
		}
	}

	return &domainmodel.HouseHold{
		ID:            domainmodel.HouseHoldID(model.ID),
		Title:         model.Title,
		Description:   model.Description,
		CategoryLimit: categoryLimits,
	}, nil
}

// FindByUserID implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) FindByUserID(userID domainmodel.UserID) (*domainmodel.HouseHold, error) {
	panic("unimplemented")
}

// Update implements domainmodel.HouseHoldRepository.
func (h *HouseHoldRepository) Update(houseHold *domainmodel.HouseHold) error {
	panic("unimplemented")
}

func NewHouseHoldRepository(db *gorm.DB) domainmodel.HouseHoldRepository {
	return &HouseHoldRepository{
		db: db,
	}
}
