package domainmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummarizeMonthlyGroupByCategory(t *testing.T) {
	tests := []struct {
		name     string
		amounts  ShoppingAmounts
		expected CategoryAmounts
	}{
		{
			name: "カテゴリーごとの合計金額が正しく計算される",
			amounts: ShoppingAmounts{
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     1000,
					Date:       "2024-01-01",
				},
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     2000,
					Date:       "2024-01-02",
				},
				{
					CategoryID: CategoryID(2),
					Category:   Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"},
					Amount:     1500,
					Date:       "2024-01-03",
				},
			},
			expected: CategoryAmounts{
				{Category: Category{ID: CategoryID(1), Name: "食費", Color: "#000000"}, Amount: 3000},
				{Category: Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"}, Amount: 1500},
			},
		},
		{
			name:     "空のデータの場合、空のマップが返される",
			amounts:  ShoppingAmounts{},
			expected: CategoryAmounts{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.amounts.SummarizeMonthlyGroupByCategory()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSummarizeMonthly(t *testing.T) {
	tests := []struct {
		name     string
		amounts  ShoppingAmounts
		expected int
	}{
		{
			name: "日付ごとの合計金額が正しく計算される",
			amounts: ShoppingAmounts{
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     1000,
					Date:       "2024-01-01",
				},
				{
					CategoryID: CategoryID(2),
					Category:   Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"},
					Amount:     2000,
					Date:       "2024-01-01",
				},
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     1500,
					Date:       "2024-01-02",
				},
			},
			expected: 4500,
		},
		{
			name:     "空のデータの場合、0が返される",
			amounts:  ShoppingAmounts{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.amounts.SummarizeMonthly()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewSummarizeShoppingAmounts(t *testing.T) {
	tests := []struct {
		name            string
		shoppingAmounts ShoppingAmounts
		expected        *SummarizeShoppingAmounts
	}{
		{
			name: "買い物金額の集計が正しく行われる",
			shoppingAmounts: ShoppingAmounts{
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     1000,
					Date:       "2024-01-01",
				},
				{
					CategoryID: CategoryID(1),
					Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
					Amount:     2000,
					Date:       "2024-01-02",
				},
				{
					CategoryID: CategoryID(2),
					Category:   Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"},
					Amount:     1500,
					Date:       "2024-01-03",
				},
			},
			expected: &SummarizeShoppingAmounts{
				ShoppingAmounts: ShoppingAmounts{
					{
						CategoryID: CategoryID(1),
						Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
						Amount:     1000,
						Date:       "2024-01-01",
					},
					{
						CategoryID: CategoryID(1),
						Category:   Category{ID: CategoryID(1), Name: "食費", Color: "#000000"},
						Amount:     2000,
						Date:       "2024-01-02",
					},
					{
						CategoryID: CategoryID(2),
						Category:   Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"},
						Amount:     1500,
						Date:       "2024-01-03",
					},
				},
				TotalAmount: 4500,
				CategoryAmounts: CategoryAmounts{
					{Category: Category{ID: CategoryID(1), Name: "食費", Color: "#000000"}, Amount: 3000},
					{Category: Category{ID: CategoryID(2), Name: "日用品", Color: "#000000"}, Amount: 1500},
				},
			},
		},
		{
			name:            "空のデータの場合、空の集計結果が返される",
			shoppingAmounts: ShoppingAmounts{},
			expected: &SummarizeShoppingAmounts{
				ShoppingAmounts: ShoppingAmounts{},
				TotalAmount:     0,
				CategoryAmounts: CategoryAmounts{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSummarizeShoppingAmounts(tt.shoppingAmounts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
