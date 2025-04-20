package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type shoppingRepository struct {
	db *gorm.DB
}

func NewShoppingRepository(db *gorm.DB) domainmodel.ShoppingRepository {
	return &shoppingRepository{db: db}
}

// DeleteShoppingAmount implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) DeleteShoppingAmount(id string) error {
	panic("unimplemented")
}

// DeleteShoppingMemo implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) DeleteShoppingMemo(id domainmodel.ShoppingID) error {
	model := models.ShoppingMemo{
		Base: models.Base{
			ID: uint(id),
		},
	}

	if err := s.db.Delete(&model).Error; err != nil {
		return err
	}
	return nil
}

// FetchShoppingAmountItem implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) FetchShoppingAmountItem(id string) (*domainmodel.ShoppingAmount, error) {
	panic("unimplemented")
}

// FetchShoppingMemoItem implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) FetchShoppingMemoItem(householdID domainmodel.HouseHoldID) ([]*domainmodel.ShoppingMemo, error) {
	model := []models.ShoppingMemo{}
	if err := s.db.Where("household_book_id = ?", householdID).Find(&model).Error; err != nil {
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
		})
	}
	return shoppingMemo, nil
}

// RegisterShoppingAmount implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) RegisterShoppingAmount(shopping *domainmodel.ShoppingAmount) error {
	panic("unimplemented")
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

// UpdateShoppingAmount implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) UpdateShoppingAmount(shopping *domainmodel.ShoppingAmount) error {
	panic("unimplemented")
}

// UpdateShoppingMemo implements domainmodel.ShoppingRepository.
func (s *shoppingRepository) UpdateShoppingMemo(shopping *domainmodel.ShoppingMemo) error {
	panic("unimplemented")
}
