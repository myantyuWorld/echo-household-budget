package functioncalling

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/llm"
	"fmt"
)

type LimitRetrievalTool struct {
	repository domainmodel.HouseHoldRepository
}

func NewLimitRetrievalTool(repository domainmodel.HouseHoldRepository) llm.Tool {
	return &LimitRetrievalTool{
		repository: repository,
	}
}

func (t *LimitRetrievalTool) Name() string {
	return "get_monthly_limits"
}

func (t *LimitRetrievalTool) Description() string {
	return "指定された家計簿の月間支出制限を取得する"
}

func (t *LimitRetrievalTool) Execute(params map[string]interface{}) (interface{}, error) {
	householdID, ok := params["household_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("household_id is required")
	}

	household, err := t.repository.FindByHouseHoldID(domainmodel.HouseHoldID(householdID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch household: %w", err)
	}

	if household == nil {
		return nil, fmt.Errorf("household not found")
	}

	return t.processHouseholdLimits(household), nil
}

func (t *LimitRetrievalTool) processHouseholdLimits(household *domainmodel.HouseHold) map[string]interface{} {
	totalLimit := 0
	categoryLimits := make([]map[string]interface{}, 0)

	for _, categoryLimit := range household.CategoryLimit {
		totalLimit += categoryLimit.LimitAmount
		categoryLimits = append(categoryLimits, map[string]interface{}{
			"category_id":    int(categoryLimit.Category.ID),
			"category_name":  categoryLimit.Category.Name,
			"limit_amount":   categoryLimit.LimitAmount,
			"category_color": categoryLimit.Category.Color,
		})
	}

	return map[string]interface{}{
		"household_id":     int(household.ID),
		"household_title":  household.Title,
		"total_limit":      totalLimit,
		"category_limits":  categoryLimits,
		"limit_count":      len(household.CategoryLimit),
	}
}