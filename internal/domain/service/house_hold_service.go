package domainservice

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"
	"errors"
	"time"
)

type HouseHoldService interface {
	FetchHouseHold(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error)
	ShareHouseHold(houseHoldID domainmodel.HouseHoldID, inviteUserID domainmodel.UserID) error
	FetchShoppingAmount(input FetchShoppingRecordInput) ([]*domainmodel.ShoppingAmount, error)
	AddUserHouseHold(houseHold *domainmodel.HouseHold) error
	AddHouseHoldCategory(houseHoldID domainmodel.HouseHoldID, categoryName string, categoryLimitAmount int) error
	CreateShoppingAmount(shoppingAmount *domainmodel.ShoppingAmount) error
	RemoveShoppingAmount(shoppingAmountID domainmodel.ShoppingID) error
	SummarizeShoppingAmount(input FetchShoppingRecordInput) (*domainmodel.SummarizeShoppingAmounts, error)
}

type houseHoldService struct {
	houseHoldRepository domainmodel.HouseHoldRepository
	shoppingRepository  domainmodel.ShoppingRepository
	categoryRepository  domainmodel.CategoryRepository
}

// AddHouseHoldCategory implements HouseHoldService.
func (h *houseHoldService) AddHouseHoldCategory(houseHoldID domainmodel.HouseHoldID, categoryName string, categoryLimitAmount int) error {
	category := &domainmodel.Category{
		Name:  categoryName,
		Color: "#000000",
	}

	if err := h.categoryRepository.CreateMasterCategory(category); err != nil {
		return err
	}

	categoryLimit := &domainmodel.CategoryLimit{
		HouseholdBookID: houseHoldID,
		Category:        *category,
		LimitAmount:     categoryLimitAmount,
	}

	if err := h.categoryRepository.CreateHouseHoldCategory(categoryLimit); err != nil {
		return err
	}

	return nil
}

// AddUserHouseHold implements HouseHoldService.
func (h *houseHoldService) AddUserHouseHold(houseHold *domainmodel.HouseHold) error {
	if err := h.houseHoldRepository.Create(houseHold); err != nil {
		return err
	}

	userHouseHold := &domainmodel.UserHouseHold{
		HouseHoldID: houseHold.ID,
		UserID:      houseHold.UserID,
	}

	if err := h.houseHoldRepository.CreateUserHouseHold(userHouseHold); err != nil {
		return err
	}

	return nil
}

type FetchShoppingRecordInput struct {
	HouseholdID domainmodel.HouseHoldID
	Date        string
}

// SummarizeShoppingAmount implements HouseHoldService.
func (h *houseHoldService) SummarizeShoppingAmount(input FetchShoppingRecordInput) (*domainmodel.SummarizeShoppingAmounts, error) {
	results, err := h.FetchShoppingAmount(input)
	if err != nil {
		return nil, err
	}

	shoppingAmounts := domainmodel.ShoppingAmounts{}
	for _, v := range results {
		shoppingAmounts = append(shoppingAmounts, v)
	}

	return domainmodel.NewSummarizeShoppingAmounts(shoppingAmounts), nil
}

// CreateShoppingAmount implements HouseHoldService.
func (h *houseHoldService) CreateShoppingAmount(shoppingAmount *domainmodel.ShoppingAmount) error {
	date, err := time.Parse("2006-01-02", shoppingAmount.Date)
	if err != nil {
		return errors.New("domainservice::CreateShoppingAmount failed to parse date")
	}
	model := &models.ShoppingAmount{
		HouseholdBookID: uint(shoppingAmount.HouseholdID),
		CategoryID:      uint(shoppingAmount.CategoryID),
		Amount:          shoppingAmount.Amount,
		Date:            date,
		Memo:            shoppingAmount.Memo,
	}

	return h.shoppingRepository.RegisterShoppingAmount(model)
}

// FetchShoppingAmount implements HouseHoldService.
func (h *houseHoldService) FetchShoppingAmount(input FetchShoppingRecordInput) ([]*domainmodel.ShoppingAmount, error) {
	shoppingAmount, err := h.shoppingRepository.FetchShoppingAmountItemByHouseholdID(input.HouseholdID, input.Date)
	if err != nil {
		return nil, err
	}

	shoppingAmounts := []*domainmodel.ShoppingAmount{}
	for _, v := range shoppingAmount {
		shoppingAmounts = append(shoppingAmounts, domainmodel.ConvertShoppingAmountsToShoppingAmount(v))
	}
	return shoppingAmounts, nil
}

// RemoveShoppingAmount implements HouseHoldService.
func (h *houseHoldService) RemoveShoppingAmount(shoppingAmountID domainmodel.ShoppingID) error {
	return h.shoppingRepository.DeleteShoppingAmount(shoppingAmountID)
}

// FetchHouseHold implements HouseHoldService.
func (h *houseHoldService) FetchHouseHold(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error) {
	houseHold, err := h.houseHoldRepository.FindByHouseHoldID(houseHoldID)
	if err != nil {
		return nil, err
	}
	return houseHold, nil
}

// ShareHouseHold implements HouseHoldService.
func (h *houseHoldService) ShareHouseHold(houseHoldID domainmodel.HouseHoldID, inviteUserID domainmodel.UserID) error {
	userHouseHold := &domainmodel.UserHouseHold{
		HouseHoldID: houseHoldID,
		UserID:      inviteUserID,
	}

	return h.houseHoldRepository.CreateUserHouseHold(userHouseHold)
}

func NewHouseHoldService(houseHoldRepository domainmodel.HouseHoldRepository, shoppingRepository domainmodel.ShoppingRepository, categoryRepository domainmodel.CategoryRepository) HouseHoldService {
	return &houseHoldService{
		houseHoldRepository: houseHoldRepository,
		shoppingRepository:  shoppingRepository,
		categoryRepository:  categoryRepository,
	}
}
