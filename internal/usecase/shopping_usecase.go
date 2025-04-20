//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"fmt"
)

type shoppingUsecase struct {
	repo domainmodel.ShoppingRepository
}

// FetchShopping implements ShoppingUsecase.
func (s *shoppingUsecase) FetchShopping(householdID domainmodel.HouseHoldID) ([]*domainmodel.ShoppingMemo, error) {
	return s.repo.FetchShoppingMemoItem(householdID)
}

// CreateShopping implements ShoppingUsecase.
func (s *shoppingUsecase) CreateShopping(shopping *domainmodel.ShoppingMemo) error {
	fmt.Println("func (s *shoppingUsecase) CreateShopping(shopping *domainmodel.ShoppingMemo) error {")
	fmt.Println("shopping", shopping)
	return s.repo.RegisterShoppingMemo(shopping)
}

// DeleteShopping implements ShoppingUsecase.
func (s *shoppingUsecase) DeleteShopping(id string) error {
	panic("unimplemented")
}

type ShoppingUsecase interface {
	CreateShopping(shopping *domainmodel.ShoppingMemo) error
	FetchShopping(householdID domainmodel.HouseHoldID) ([]*domainmodel.ShoppingMemo, error)
	DeleteShopping(id string) error
}

func NewShoppingUsecase(repo domainmodel.ShoppingRepository) ShoppingUsecase {
	return &shoppingUsecase{repo: repo}
}
