package model

import (
	"reflect"
	"testing"
)

func TestKaimemoAmountRecords_GroupByWeek(t *testing.T) {
	tests := []struct {
		name     string
		records  KaimemoAmountRecords
		expected []WeeklySummary
	}{
		{
			name: "empty records",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{},
			},
			expected: []WeeklySummary{},
		},
		{
			name: "single week records",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-15", Amount: 1000},
					{Date: "2023-05-16", Amount: 2000},
					{Date: "2023-05-17", Amount: 3000},
				},
			},
			expected: []WeeklySummary{
				{
					WeekStart:   "2023-05-14",
					WeekEnd:     "2023-05-20",
					TotalAmount: 6000,
					Items: []KaimemoAmount{
						{Date: "2023-05-15", Amount: 1000},
						{Date: "2023-05-16", Amount: 2000},
						{Date: "2023-05-17", Amount: 3000},
					},
				},
			},
		},
		{
			name: "multiple weeks across month boundary",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-30", Amount: 1000},
					{Date: "2023-06-01", Amount: 2000},
					{Date: "2023-06-05", Amount: 3000},
				},
			},
			expected: []WeeklySummary{
				{
					WeekStart:   "2023-05-28",
					WeekEnd:     "2023-06-03",
					TotalAmount: 3000,
					Items: []KaimemoAmount{
						{Date: "2023-05-30", Amount: 1000},
						{Date: "2023-06-01", Amount: 2000},
					},
				},
				{
					WeekStart:   "2023-06-04",
					WeekEnd:     "2023-06-10",
					TotalAmount: 3000,
					Items: []KaimemoAmount{
						{Date: "2023-06-05", Amount: 3000},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.records.GroupByWeek()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\nGroupByWeek() =\n%+v\nwant\n%+v", result, tt.expected)
			}
		})
	}
}
func TestKaimemoAmountRecords_GroupByMonth(t *testing.T) {
	tests := []struct {
		name     string
		records  KaimemoAmountRecords
		expected []MonthlySummary
	}{
		{
			name: "empty records",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{},
			},
			expected: []MonthlySummary{},
		},
		{
			name: "single month with multiple tags",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-15", Amount: 1000, Tag: "food"},
					{Date: "2023-05-16", Amount: 2000, Tag: "transport"},
					{Date: "2023-05-17", Amount: 3000, Tag: "food"},
				},
			},
			expected: []MonthlySummary{
				{
					Month:       "2023-05",
					TotalAmount: 6000,
					TagSummary: map[string]int{
						"food":      4000,
						"transport": 2000,
					},
				},
			},
		},
		{
			name: "multiple months with sorting",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-06-01", Amount: 2000, Tag: "food"},
					{Date: "2023-05-15", Amount: 1000, Tag: "transport"},
					{Date: "2023-07-01", Amount: 3000, Tag: "entertainment"},
				},
			},
			expected: []MonthlySummary{
				{
					Month:       "2023-05",
					TotalAmount: 1000,
					TagSummary: map[string]int{
						"transport": 1000,
					},
				},
				{
					Month:       "2023-06",
					TotalAmount: 2000,
					TagSummary: map[string]int{
						"food": 2000,
					},
				},
				{
					Month:       "2023-07",
					TotalAmount: 3000,
					TagSummary: map[string]int{
						"entertainment": 3000,
					},
				},
			},
		},
		{
			name: "same tag across different months",
			records: KaimemoAmountRecords{
				Records: []KaimemoAmount{
					{Date: "2023-05-15", Amount: 1000, Tag: "food"},
					{Date: "2023-06-16", Amount: 2000, Tag: "food"},
					{Date: "2023-06-17", Amount: 3000, Tag: "food"},
				},
			},
			expected: []MonthlySummary{
				{
					Month:       "2023-05",
					TotalAmount: 1000,
					TagSummary: map[string]int{
						"food": 1000,
					},
				},
				{
					Month:       "2023-06",
					TotalAmount: 5000,
					TagSummary: map[string]int{
						"food": 5000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.records.GroupByMonth()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("\nGroupByMonth() =\n%+v\nwant\n%+v", result, tt.expected)
			}
		})
	}
}
