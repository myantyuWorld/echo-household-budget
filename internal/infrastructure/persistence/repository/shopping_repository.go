package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"
	"time"

	"gorm.io/gorm"
)

type shoppingRepository struct {
	db *gorm.DB
}

func NewShoppingRepository(db *gorm.DB) domainmodel.ShoppingRepository {
	return &shoppingRepository{db: db}
}

// DeleteShoppingAmount implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) DeleteShoppingAmount(id domainmodel.ShoppingID) error {
	model := models.ShoppingAmount{
		Base: models.Base{
			ID: uint(id),
		},
	}

	if err := s.db.Delete(&model).Error; err != nil {
		return err
	}
	return nil
}

// FetchShoppingAmountItemByHouseholdID implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) FetchShoppingAmountItemByHouseholdID(householdID domainmodel.HouseHoldID, date string) ([]*models.ShoppingAmount, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	startDateMonth := time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, dateTime.Location())
	endDateMonth := startDateMonth.AddDate(0, 1, 0)

	// TODO
	// カテゴリについて、家計簿ごとに、上限金額を設定できるようにした上で、上限金額を取得するようにする
	model := []*models.ShoppingAmount{}
	if err := s.db.Debug().Where("household_book_id = ? AND date BETWEEN ? AND ?", householdID, startDateMonth, endDateMonth).Preload("Category").
		Preload("Analyze.Items").Find(&model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

// DeleteShoppingMemo implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) DeleteShoppingMemo(id domainmodel.ShoppingID) error {
	model := models.ShoppingMemo{
		Base: models.Base{
			ID: uint(id),
		},
	}

	if err := s.db.Model(&model).Update("is_completed", true).Error; err != nil {
		return err
	}
	return nil
}

// FetchShoppingMemoItem implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) FetchShoppingMemoItem(householdID domainmodel.HouseHoldID) ([]*domainmodel.ShoppingMemo, error) {
	model := []models.ShoppingMemo{}
	if err := s.db.Debug().Where("household_book_id = ? AND is_completed = ?", householdID, false).Preload("Category").Find(&model).Error; err != nil {
		return nil, err
	}

	shoppingMemo := []*domainmodel.ShoppingMemo{}
	for _, v := range model {
		shoppingMemo = append(shoppingMemo, &domainmodel.ShoppingMemo{
			ID:          domainmodel.ShoppingID(v.ID),
			HouseholdID: domainmodel.HouseHoldID(v.HouseholdBookID),
			CategoryID:  domainmodel.CategoryID(v.CategoryID),
			Title:       v.Title,
			Memo:        v.Memo,
			IsCompleted: domainmodel.IsCompleted(v.IsCompleted),
			Category: domainmodel.Category{
				ID:   domainmodel.CategoryID(v.Category.ID),
				Name: v.Category.Name,
			},
		})
	}
	return shoppingMemo, nil
}

// RegisterShoppingAmount implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) RegisterShoppingAmount(shopping *models.ShoppingAmount) error {
	if err := s.db.Create(shopping).Error; err != nil {
		return err
	}
	return nil
}

// RegisterShoppingMemo implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) RegisterShoppingMemo(shopping *domainmodel.ShoppingMemo) error {
	model := models.ShoppingMemo{
		HouseholdBookID: uint(shopping.HouseholdID),
		CategoryID:      uint(shopping.CategoryID),
		Title:           shopping.Title,
		Memo:            shopping.Memo,
		IsCompleted:     bool(shopping.IsCompleted),
	}

	if err := s.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}
