package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"
	"errors"

	"gorm.io/gorm"
)

// CategoryRepository はカテゴリリポジトリの実装
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository は新しいCategoryRepositoryを作成します
func NewCategoryRepository(db *gorm.DB) domainmodel.CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// CreateHouseHoldCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) CreateHouseHoldCategory(categoryLimit *domainmodel.CategoryLimit) error {
	model := &models.CategoryLimit{
		HouseholdBookID: uint(categoryLimit.HouseholdBookID),
		CategoryID:      uint(categoryLimit.CategoryID),
		LimitAmount:     categoryLimit.LimitAmount,
	}

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	return nil
}

// CreateMasterCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) CreateMasterCategory(category *domainmodel.Category) error {
	model := &models.Category{
		Name:  category.Name,
		Color: category.Color,
	}

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	category.ID = domainmodel.CategoryID(model.ID)

	return nil
}

// DeleteHouseHoldCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) DeleteHouseHoldCategory(categoryLimitID domainmodel.CategoryLimitID) error {
	result := r.db.Delete(&models.CategoryLimit{}, categoryLimitID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// DeleteMasterCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) DeleteMasterCategory(id domainmodel.CategoryID) error {
	result := r.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// FindHouseHoldCategoryByHouseHoldID implements domainmodel.CategoryRepository.
func (r *CategoryRepository) FindHouseHoldCategoryByHouseHoldID(categoryLimitID domainmodel.CategoryLimitID) (*domainmodel.CategoryLimit, error) {
	var model models.CategoryLimit
	if err := r.db.First(&model, categoryLimitID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domainmodel.CategoryLimit{
		ID:              domainmodel.CategoryLimitID(model.ID),
		HouseholdBookID: domainmodel.HouseHoldID(model.HouseholdBookID),
		CategoryID:      domainmodel.CategoryID(model.CategoryID),
		LimitAmount:     model.LimitAmount,
	}, nil
}

// FindMasterCategoryByID implements domainmodel.CategoryRepository.
func (r *CategoryRepository) FindMasterCategoryByID(categoryID domainmodel.CategoryID) (*domainmodel.Category, error) {
	var model models.Category
	if err := r.db.First(&model, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &domainmodel.Category{
		ID:    domainmodel.CategoryID(model.ID),
		Name:  model.Name,
		Color: model.Color,
	}, nil
}

// UpdateHouseHoldCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) UpdateHouseHoldCategory(categoryLimit *domainmodel.CategoryLimit) error {
	model := &models.CategoryLimit{
		Base: models.Base{
			ID: uint(categoryLimit.ID),
		},
		HouseholdBookID: uint(categoryLimit.HouseholdBookID),
		CategoryID:      uint(categoryLimit.CategoryID),
		LimitAmount:     categoryLimit.LimitAmount,
	}

	result := r.db.Model(model).Updates(map[string]interface{}{
		"household_book_id": categoryLimit.HouseholdBookID,
		"category_id":       categoryLimit.CategoryID,
		"limit_amount":      categoryLimit.LimitAmount,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateMasterCategory implements domainmodel.CategoryRepository.
func (r *CategoryRepository) UpdateMasterCategory(category *domainmodel.Category) error {
	model := &models.Category{
		Base: models.Base{
			ID: uint(category.ID),
		},
		Name:  category.Name,
		Color: category.Color,
	}

	result := r.db.Model(model).Updates(map[string]interface{}{
		"name":  category.Name,
		"color": category.Color,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
