//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainmodel

import (
	"echo-household-budget/internal/infrastructure/persistence/models"
)

type ShoppingMemo struct {
	ID          ShoppingID  `json:"id"`
	HouseholdID HouseHoldID `json:"householdID"`
	CategoryID  CategoryID  `json:"categoryID"`
	Title       string      `json:"title"`
	Memo        string      `json:"memo"`
	IsCompleted IsCompleted `json:"isCompleted"`
	Category    Category    `json:"category"`
}

type ShoppingAmount struct {
	ID          ShoppingID  `json:"id"`
	HouseholdID HouseHoldID `json:"household_id"`
	CategoryID  CategoryID  `json:"category_id"`
	Amount      int         `json:"amount"`
	Date        string      `json:"date"`
	Memo        string      `json:"memo"`
	Category    Category    `json:"category"`
	AnalyzeID   int         `json:"analyze_id"`
}

type CategoryAmount struct {
	Category Category `json:"category"`
	Amount   int      `json:"amount"`
}

type ShoppingAmounts []*ShoppingAmount
type CategoryAmounts []*CategoryAmount

type SummarizeShoppingAmounts struct {
	ShoppingAmounts ShoppingAmounts `json:"shoppingAmounts"`
	TotalAmount     int             `json:"totalAmount"`
	CategoryAmounts CategoryAmounts `json:"categoryAmounts"`
}

func (s *ShoppingAmounts) SummarizeMonthlyGroupByCategory() CategoryAmounts {
	amounts := CategoryAmounts{}
	categoryMap := make(map[CategoryID]*CategoryAmount)
	for _, amount := range *s {
		if existing, ok := categoryMap[amount.CategoryID]; ok {
			existing.Amount += amount.Amount
		} else {
			categoryMap[amount.CategoryID] = &CategoryAmount{
				Category: amount.Category,
				Amount:   amount.Amount,
			}
		}
	}
	for _, amount := range categoryMap {
		amounts = append(amounts, amount)
	}
	return amounts
}

func (s *ShoppingAmounts) SummarizeMonthly() int {
	amounts := 0
	for _, amount := range *s {
		amounts += amount.Amount
	}
	return amounts
}

func NewSummarizeShoppingAmounts(shoppingAmounts ShoppingAmounts) *SummarizeShoppingAmounts {
	arr := ShoppingAmounts{}
	arr = append(arr, shoppingAmounts...)

	return &SummarizeShoppingAmounts{
		ShoppingAmounts: arr,
		TotalAmount:     arr.SummarizeMonthly(),
		CategoryAmounts: arr.SummarizeMonthlyGroupByCategory(),
	}
}

func ConvertShoppingAmountsToShoppingAmount(shoppingAmount *models.ShoppingAmount) *ShoppingAmount {
	return &ShoppingAmount{
		ID:          ShoppingID(shoppingAmount.ID),
		HouseholdID: HouseHoldID(shoppingAmount.HouseholdBookID),
		CategoryID:  CategoryID(shoppingAmount.CategoryID),
		Amount:      shoppingAmount.Amount,
		Date:        shoppingAmount.Date.Format("2006-01-02"),
		Memo:        shoppingAmount.Memo,
		// TODO : カテゴリについて、家計簿ごとに、上限金額を設定できるようにした上で、上限金額を取得するようにする
		// HouseHoldCategory、のようなモデルが必要か
		Category: Category{ID: CategoryID(shoppingAmount.CategoryID), Name: shoppingAmount.Category.Name, Color: shoppingAmount.Category.Color},
	}
}

func NewShoppingMemo(householdID HouseHoldID, categoryID CategoryID, title string, memo string) *ShoppingMemo {
	return &ShoppingMemo{
		HouseholdID: householdID,
		CategoryID:  categoryID,
		Title:       title,
		Memo:        memo,
		IsCompleted: NotDone,
	}
}

func NewShoppingAmount(householdID HouseHoldID, categoryID CategoryID, amount int, date string, memo string, analyzeID int) *ShoppingAmount {
	return &ShoppingAmount{
		HouseholdID: householdID,
		CategoryID:  categoryID,
		Amount:      amount,
		Date:        date,
		Memo:        memo,
		AnalyzeID:   analyzeID,
	}
}

type ShoppingID uint
type IsCompleted bool

const (
	Done    IsCompleted = true
	NotDone IsCompleted = false
)

type ShoppingRepository interface {
	RegisterShoppingMemo(shopping *ShoppingMemo) error
	FetchShoppingMemoItem(householdID HouseHoldID) ([]*ShoppingMemo, error)
	DeleteShoppingMemo(id ShoppingID) error
	RegisterShoppingAmount(shopping *models.ShoppingAmount) error
	FetchShoppingAmountItemByHouseholdID(householdID HouseHoldID, date string) ([]*models.ShoppingAmount, error)
	DeleteShoppingAmount(id ShoppingID) error
}
