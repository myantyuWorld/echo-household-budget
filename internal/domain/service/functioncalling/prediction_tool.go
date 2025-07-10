package functioncalling

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/llm"
	"echo-household-budget/internal/infrastructure/persistence/models"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type PredictionTool struct {
	shoppingRepository  domainmodel.ShoppingRepository
	householdRepository domainmodel.HouseHoldRepository
}

func NewPredictionTool(shoppingRepository domainmodel.ShoppingRepository, householdRepository domainmodel.HouseHoldRepository) llm.Tool {
	return &PredictionTool{
		shoppingRepository:  shoppingRepository,
		householdRepository: householdRepository,
	}
}

func (t *PredictionTool) Name() string {
	return "predict_monthly_expenses"
}

func (t *PredictionTool) Description() string {
	return "現在の支出ペースから月末の支出予測を生成する"
}

func (t *PredictionTool) Execute(params map[string]interface{}) (interface{}, error) {
	log.Println("============ PredictionTool Execute =============")
	spew.Dump(params)

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

	currentDateStr, ok := params["current_date"].(string)
	if !ok {
		currentDateStr = time.Now().Format("2006-01-02")
	}

	currentDate, err := time.Parse("2006-01-02", currentDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid current_date format: %w", err)
	}

	date := fmt.Sprintf("%04d-%02d-01", int(year), int(month))

	shoppingAmounts, err := t.shoppingRepository.FetchShoppingAmountItemByHouseholdID(
		domainmodel.HouseHoldID(householdID),
		date,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shopping amounts: %w", err)
	}

	household, err := t.householdRepository.FindByHouseHoldID(domainmodel.HouseHoldID(householdID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch household: %w", err)
	}

	return t.generatePrediction(shoppingAmounts, household, currentDate, int(year), int(month)), nil
}

func (t *PredictionTool) generatePrediction(shoppingAmounts []*models.ShoppingAmount, household *domainmodel.HouseHold, currentDate time.Time, year, month int) map[string]interface{} {
	// 月の日数を計算
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentDate.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	totalDaysInMonth := lastDay.Day()

	// 経過日数を計算
	daysPassed := currentDate.Day()

	// 残り日数を計算
	remainingDays := totalDaysInMonth - daysPassed

	// カテゴリ別の現在の支出を集計
	categoryAmounts := make(map[uint]int)
	for _, shopping := range shoppingAmounts {
		categoryAmounts[shopping.CategoryID] += shopping.Amount
	}

	// カテゴリ別の制限を取得
	categoryLimits := make(map[uint]int)
	for _, limit := range household.CategoryLimit {
		categoryLimits[uint(limit.Category.ID)] = limit.LimitAmount
	}

	// 予測結果を生成
	categoryPredictions := make([]map[string]interface{}, 0)
	totalCurrentAmount := 0
	totalPredictedAmount := 0
	totalLimitAmount := 0

	for categoryID, currentAmount := range categoryAmounts {
		totalCurrentAmount += currentAmount

		// 日割りで予測
		dailyAverage := float64(currentAmount) / float64(daysPassed)
		predictedAmount := int(dailyAverage * float64(totalDaysInMonth))
		totalPredictedAmount += predictedAmount

		// 制限との比較
		limitAmount := categoryLimits[categoryID]
		totalLimitAmount += limitAmount

		isOverBudget := predictedAmount > limitAmount
		remainingAmount := limitAmount - currentAmount

		categoryPredictions = append(categoryPredictions, map[string]interface{}{
			"category_id":      categoryID,
			"current_amount":   currentAmount,
			"predicted_amount": predictedAmount,
			"limit_amount":     limitAmount,
			"is_over_budget":   isOverBudget,
			"remaining_amount": remainingAmount,
			"daily_average":    dailyAverage,
		})
	}

	return map[string]interface{}{
		"household_id":           int(household.ID),
		"year":                   year,
		"month":                  month,
		"current_date":           currentDate.Format("2006-01-02"),
		"days_passed":            daysPassed,
		"remaining_days":         remainingDays,
		"total_days_in_month":    totalDaysInMonth,
		"total_current_amount":   totalCurrentAmount,
		"total_predicted_amount": totalPredictedAmount,
		"total_limit_amount":     totalLimitAmount,
		"category_predictions":   categoryPredictions,
		"overall_budget_status":  totalPredictedAmount <= totalLimitAmount,
		"prediction_accuracy":    fmt.Sprintf("%.1f%%", float64(daysPassed)/float64(totalDaysInMonth)*100),
	}
}
