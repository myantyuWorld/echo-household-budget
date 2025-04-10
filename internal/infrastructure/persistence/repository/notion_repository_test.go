package repository

import (
	"template-echo-notion-integration/internal/model"
	"testing"
	"time"

	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewNotionRepository(t *testing.T) {
	testCases := []struct {
		name                           string
		apiKey                         string
		databaseKaimemoInputID         string
		databaseKaimemoSummaryRecordID string
	}{
		{
			name:                           "should create new repository with valid credentials",
			apiKey:                         "test-api-key",
			databaseKaimemoInputID:         "test-database-id",
			databaseKaimemoSummaryRecordID: "test-database-id",
		},
		{
			name:                           "should create new repository with empty credentials",
			apiKey:                         "",
			databaseKaimemoInputID:         "",
			databaseKaimemoSummaryRecordID: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewNotionRepository(tc.apiKey, tc.databaseKaimemoInputID, tc.databaseKaimemoSummaryRecordID)

			assert.NotNil(t, repo)

			kaimemoRepo, ok := repo.(*kaimemoRepository)
			assert.True(t, ok)
			assert.NotNil(t, kaimemoRepo.client)
			assert.Equal(t, tc.databaseKaimemoInputID, kaimemoRepo.databaseKaimemoInputID)
			assert.Equal(t, tc.databaseKaimemoSummaryRecordID, kaimemoRepo.databaseKaimemoSummaryRecordID)
			assert.Equal(t, notionapi.Token(tc.apiKey), kaimemoRepo.client.Token)
		})
	}
}

type MockKaimemoRepository struct {
	mock.Mock
}

func (m *MockKaimemoRepository) FetchKaimemoAmountRecords() (model.KaimemoAmountRecords, error) {
	args := m.Called()
	return args.Get(0).(model.KaimemoAmountRecords), args.Error(1)
}

func (m *MockKaimemoRepository) RemoveKaimemoAmount(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestMockKaimemoRepository(t *testing.T) {
	mockRepo := new(MockKaimemoRepository)

	testData := model.KaimemoAmountRecords{
		Records: []model.KaimemoAmount{
			{
				ID:     "1",
				Date:   time.Now().Format("2006-01-02"),
				Amount: 1000,
			},
			{
				ID:     "2",
				Date:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
				Amount: 2000,
			},
		},
	}

	// モックの振る舞いを設定
	mockRepo.On("FetchKaimemoAmountRecords").Return(testData, nil)
	mockRepo.On("RemoveKaimemoAmount", "1").Return(nil)

	res, err := mockRepo.FetchKaimemoAmountRecords()
	assert.NoError(t, err)
	assert.NotEqual(t, res, nil)
	assert.Equal(t, 2, len(res.Records))
	assert.Equal(t, 1000, res.Records[0].Amount)

	err = mockRepo.RemoveKaimemoAmount("1")
	assert.NoError(t, err)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
}
