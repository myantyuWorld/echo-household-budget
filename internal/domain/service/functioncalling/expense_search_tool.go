package functioncalling

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/llm"
	"echo-household-budget/internal/infrastructure/persistence/models"
	"fmt"
	"strconv"
	"time"
)

type ExpenseSearchTool struct {
	repository domainmodel.ShoppingRepository
}

func NewExpenseSearchTool(repository domainmodel.ShoppingRepository) llm.Tool {
	return &ExpenseSearchTool{
		repository: repository,
	}
}

func (t *ExpenseSearchTool) Name() string {
	return "search_monthly_expenses"
}

func (t *ExpenseSearchTool) Description() string {
	return "指定された月の支出データを検索し、カテゴリ別に集計する"
}

func (t *ExpenseSearchTool) Execute(params map[string]interface{}) (interface{}, error) {
	householdID, ok := params["household_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("household_id is required")
	}

	year, ok := params["year"].(float64)
	if !ok {
		return nil, fmt.Errorf("year is required")
	}

	month, ok := params["month"].(float64)
	if !ok {
		return nil, fmt.Errorf("month is required")
	}

	date := fmt.Sprintf("%04d-%02d-01", int(year), int(month))
	
	shoppingAmounts, err := t.repository.FetchShoppingAmountItemByHouseholdID(
		domainmodel.HouseHoldID(householdID),
		date,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shopping amounts: %w", err)
	}

	return t.processShoppingAmounts(shoppingAmounts), nil
}

func (t *ExpenseSearchTool) processShoppingAmounts(shoppingAmounts []*models.ShoppingAmount) map[string]interface{} {
	totalAmount := 0
	categoryAmounts := make(map[string]int)
	categoryNames := make(map[string]string)

	for _, shopping := range shoppingAmounts {
		totalAmount += shopping.Amount
		categoryIDStr := strconv.Itoa(int(shopping.CategoryID))
		categoryAmounts[categoryIDStr] += shopping.Amount
		categoryNames[categoryIDStr] = shopping.Category.Name
	}

	categoryResults := make([]map[string]interface{}, 0)
	for categoryID, amount := range categoryAmounts {
		categoryResults = append(categoryResults, map[string]interface{}{
			"category_id":   categoryID,
			"category_name": categoryNames[categoryID],
			"amount":        amount,
		})
	}

	return map[string]interface{}{
		"total_amount":      totalAmount,
		"category_amounts":  categoryResults,
		"shopping_amounts":  len(shoppingAmounts),
		"analysis_date":     time.Now().Format("2006-01-02"),
	}
}