//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

import "time"

type ShoppingMemo struct {
	ID          ShoppingID  `json:"id"`
	HouseholdID HouseHoldID `json:"household_id"`
	CategoryID  CategoryID  `json:"category_id"`
	Title       string      `json:"title"`
	Memo        string      `json:"memo"`
	IsCompleted IsCompleted `json:"is_completed"`
}

type ShoppingAmount struct {
	ID          ShoppingID  `json:"id"`
	HouseholdID HouseHoldID `json:"household_id"`
	CategoryID  CategoryID  `json:"category_id"`
	Amount      int         `json:"amount"`
	Date        time.Time   `json:"date"`
	Memo        string      `json:"memo"`
}

type ShoppingID string
type IsCompleted bool

const (
	Done    IsCompleted = true
	NotDone IsCompleted = false
)

type ShoppingRepository interface {
	RegisterShoppingMemo(shopping *ShoppingMemo) error
	FetchShoppingMemoItem(id string) (*ShoppingMemo, error)
	UpdateShoppingMemo(shopping *ShoppingMemo) error
	DeleteShoppingMemo(id string) error
	RegisterShoppingAmount(shopping *ShoppingAmount) error
	FetchShoppingAmountItem(id string) (*ShoppingAmount, error)
	UpdateShoppingAmount(shopping *ShoppingAmount) error
	DeleteShoppingAmount(id string) error
}
