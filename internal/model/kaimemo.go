package model

import (
	"fmt"
	"sort"
	"time"
)

type KaimemoResponse struct {
	ID   string `json:"id"`
	Tag  string `json:"tag"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateKaimemoRequest struct {
	TempUserID string `json:"tempUserID"`
	Tag        string `json:"tag"`
	Name       string `json:"name"`
}

type TelegraphRequest struct {
	MethodType string  `json:"methodType" validate:"required"` // 1 : 登録 2 : 削除	TempUserID *string `json:"tempUserID"`
	ID         *string `json:"id"`
	Tag        *string `json:"tag"`
	Name       *string `json:"name"`
}

type RemoveKaimemoRequest struct {
	TempUserID string `json:"tempUserID"`
}
type KaimemoAmountResponse struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

type CreateKaimemoAmountRequest struct {
	TempUserID string `json:"tempUserID"`
	Date       string `json:"date"`
	Tag        string `json:"tag"`
	Amount     int    `json:"amount"`
}

type RemoveKaimemoAmountRequest struct {
	TempUserID string `json:"tempUserID"`
}

type KaimemoAmount struct {
	ID     string `json:"id"`
	Date   string `json:"date"`
	Tag    string `json:"tag"`
	Amount int    `json:"amount"`
}

type KaimemoAmountRecords struct {
	Records []KaimemoAmount
}

type MonthlySummary struct {
	Month       string         `json:"month"`
	TotalAmount int            `json:"totalAmount"`
	TagSummary  map[string]int `json:"tagSummary"`
}

type KaimemoSummaryResponse struct {
	MonthlySummaries []MonthlySummary `json:"monthlySummaries"`
	WeeklySummaries  []WeeklySummary  `json:"weeklySummaries"`
}

func (k KaimemoAmountRecords) GroupByMonth() []MonthlySummary {
	summaries := make(map[string]*MonthlySummary)

	for _, amount := range k.Records {
		date, _ := time.Parse("2006-01-02", amount.Date)
		monthKey := date.Format("2006-01")

		if _, exists := summaries[monthKey]; !exists {
			summaries[monthKey] = &MonthlySummary{
				Month:       monthKey,
				TotalAmount: 0,
				TagSummary:  make(map[string]int),
			}
		}

		summary := summaries[monthKey]
		summary.TotalAmount += amount.Amount
		summary.TagSummary[amount.Tag] += amount.Amount
	}

	result := make([]MonthlySummary, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, *summary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Month < result[j].Month
	})

	return result
}

type WeeklySummary struct {
	WeekStart   string          `json:"weekStart"`
	WeekEnd     string          `json:"weekEnd"`
	TotalAmount int             `json:"totalAmount"`
	Items       []KaimemoAmount `json:"items"`
}

func (k KaimemoAmountRecords) GroupByWeek() []WeeklySummary {
	summaries := make(map[string]*WeeklySummary)

	for _, amount := range k.Records {
		date, _ := time.Parse("2006-01-02", amount.Date)
		year, week := date.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%d", year, week)

		if _, exists := summaries[weekKey]; !exists {
			// 週の開始日を計算
			weekStart := date.AddDate(0, 0, -int(date.Weekday()))
			weekEnd := weekStart.AddDate(0, 0, 6)

			summaries[weekKey] = &WeeklySummary{
				WeekStart:   weekStart.Format("2006-01-02"),
				WeekEnd:     weekEnd.Format("2006-01-02"),
				TotalAmount: 0,
				Items:       []KaimemoAmount{},
			}
		}

		summary := summaries[weekKey]
		summary.TotalAmount += amount.Amount
		summary.Items = append(summary.Items, amount)
	}

	// マップをスライスに変換
	result := make([]WeeklySummary, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, *summary)
	}

	// 日付順にソート
	sort.Slice(result, func(i, j int) bool {
		return result[i].WeekStart < result[j].WeekStart
	})

	return result
}
