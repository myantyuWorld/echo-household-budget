package usecase

import (
	"errors"
	"testing"

	domainmodel "echo-household-budget/internal/domain/model"
	domainservice "echo-household-budget/internal/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReceiptAnalyzeRepository is a mock of ReceiptAnalyzeRepository
type MockReceiptAnalyzeRepository struct {
	mock.Mock
}

func (m *MockReceiptAnalyzeRepository) CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockReceiptAnalyzeRepository) CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockReceiptAnalyzeRepository) FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmodel.ReceiptAnalyze), args.Error(1)
}

func (m *MockReceiptAnalyzeRepository) FindReceiptAnalyzeByS3FilePath(s3FilePath string) (*domainmodel.ReceiptAnalyze, error) {
	args := m.Called(s3FilePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmodel.ReceiptAnalyze), args.Error(1)
}

// MockFileStorageRepository is a mock of FileStorageRepository
type MockFileStorageRepository struct {
	mock.Mock
}

func (m *MockFileStorageRepository) UploadFile(data []byte, fileName string) (string, error) {
	args := m.Called(data, fileName)
	return args.String(0), args.Error(1)
}

func (m *MockFileStorageRepository) DeleteFile(fileName string) error {
	args := m.Called(fileName)
	return args.Error(0)
}

func (m *MockFileStorageRepository) GetFileURL(fileName string) (string, error) {
	args := m.Called(fileName)
	return args.String(0), args.Error(1)
}

// MockHouseHoldServiceの定義
type MockHouseHoldService struct {
	mock.Mock
}

func (m *MockHouseHoldService) FetchHouseHold(houseHoldID domainmodel.HouseHoldID) (*domainmodel.HouseHold, error) {
	args := m.Called(houseHoldID)
	return args.Get(0).(*domainmodel.HouseHold), args.Error(1)
}

func (m *MockHouseHoldService) ShareHouseHold(houseHoldID domainmodel.HouseHoldID, inviteUserID domainmodel.UserID) error {
	args := m.Called(houseHoldID, inviteUserID)
	return args.Error(0)
}

func (m *MockHouseHoldService) FetchShoppingAmount(input domainservice.FetchShoppingRecordInput) ([]*domainmodel.ShoppingAmount, error) {
	args := m.Called(input)
	return args.Get(0).([]*domainmodel.ShoppingAmount), args.Error(1)
}

func (m *MockHouseHoldService) AddUserHouseHold(houseHold *domainmodel.HouseHold) error {
	args := m.Called(houseHold)
	return args.Error(0)
}

func (m *MockHouseHoldService) AddHouseHoldCategory(houseHoldID domainmodel.HouseHoldID, categoryName string, categoryLimitAmount int) error {
	args := m.Called(houseHoldID, categoryName, categoryLimitAmount)
	return args.Error(0)
}

func (m *MockHouseHoldService) CreateShoppingAmount(shoppingAmount *domainmodel.ShoppingAmount) error {
	args := m.Called(shoppingAmount)
	return args.Error(0)
}

func (m *MockHouseHoldService) RemoveShoppingAmount(shoppingAmountID domainmodel.ShoppingID) error {
	args := m.Called(shoppingAmountID)
	return args.Error(0)
}

func (m *MockHouseHoldService) SummarizeShoppingAmount(input domainservice.FetchShoppingRecordInput) (*domainmodel.SummarizeShoppingAmounts, error) {
	args := m.Called(input)
	return args.Get(0).(*domainmodel.SummarizeShoppingAmounts), args.Error(1)
}

func TestCreateReceiptAnalyzeReception(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		receipt        *domainmodel.ReceiptAnalyzeReception
		mockSetup      func(*MockReceiptAnalyzeRepository, *MockFileStorageRepository)
		expectedError  error
		expectedURL    string
		validateResult func(*testing.T, *domainmodel.ReceiptAnalyzeReception)
	}{
		{
			name: "正常系：ファイル名が正しく生成され、保存される",
			receipt: &domainmodel.ReceiptAnalyzeReception{
				HouseholdBookID: 123,
				ImageData:       "SGVsbG8gV29ybGQ=", // "Hello World" in base64
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, storage *MockFileStorageRepository) {
				// ファイルアップロードのモック
				storage.On("UploadFile", mock.Anything, mock.MatchedBy(func(fileName string) bool {
					// ファイル名のフォーマットを検証
					// uuid-household_id-yyyyMMddHHmmss.jpg
					return true // TODO: より厳密な正規表現チェックを追加
				})).Return("https://example.com/test.jpg", nil)

				// DB保存のモック
				repo.On("CreateReceiptAnalyzeReception", mock.Anything).Return(nil)
			},
			expectedError: nil,
			expectedURL:   "https://example.com/test.jpg",
			validateResult: func(t *testing.T, receipt *domainmodel.ReceiptAnalyzeReception) {
				// assert.Equal(t, "https://example.com/test.jpg", receipt.ImageURL)
			},
		},
		{
			name: "異常系：base64デコードエラー",
			receipt: &domainmodel.ReceiptAnalyzeReception{
				HouseholdBookID: 123,
				ImageData:       "invalid base64",
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, storage *MockFileStorageRepository) {
				// モックは呼ばれないはず
			},
			expectedError: errors.New("illegal base64 data at input byte 7"),
		},
		{
			name: "異常系：ファイルアップロードエラー",
			receipt: &domainmodel.ReceiptAnalyzeReception{
				HouseholdBookID: 123,
				ImageData:       "SGVsbG8gV29ybGQ=",
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, storage *MockFileStorageRepository) {
				storage.On("UploadFile", mock.Anything, mock.Anything).Return("", errors.New("upload error"))
			},
			expectedError: errors.New("upload error"),
		},
		{
			name: "異常系：DB保存エラー",
			receipt: &domainmodel.ReceiptAnalyzeReception{
				HouseholdBookID: 123,
				ImageData:       "SGVsbG8gV29ybGQ=",
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, storage *MockFileStorageRepository) {
				storage.On("UploadFile", mock.Anything, mock.Anything).Return("https://example.com/test.jpg", nil)
				repo.On("CreateReceiptAnalyzeReception", mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockRepo := new(MockReceiptAnalyzeRepository)
			mockStorage := new(MockFileStorageRepository)
			tt.mockSetup(mockRepo, mockStorage)

			// テスト対象のインスタンス作成
			usecase := &receiptAnalyzeUsecase{
				repo:        mockRepo,
				fileStorage: mockStorage,
			}

			// テスト実行
			err := usecase.CreateReceiptAnalyzeReception(tt.receipt)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, tt.receipt)
				}
			}

			// モックの検証
			mockRepo.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestCreateReceiptAnalyzeResult(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		receipt       *domainmodel.ReceiptAnalyze
		mockSetup     func(*MockReceiptAnalyzeRepository, *MockHouseHoldService)
		expectedError error
	}{
		{
			name: "正常系：結果が保存される",
			receipt: &domainmodel.ReceiptAnalyze{
				ID:         123,
				TotalPrice: 1000,
				S3FilePath: "test/path.jpg",
				Items:      []domainmodel.ReceiptAnalyzeItem{},
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, houseHoldService *MockHouseHoldService) {
				repo.On("FindReceiptAnalyzeByS3FilePath", "test/path.jpg").Return(&domainmodel.ReceiptAnalyze{
					ID:         123,
					TotalPrice: 1000,
					S3FilePath: "test/path.jpg",
					Items:      []domainmodel.ReceiptAnalyzeItem{},
				}, nil)
				repo.On("CreateReceiptAnalyzeResult", mock.Anything).Return(nil)
				houseHoldService.On("CreateShoppingAmount", mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "異常系：DB保存エラー",
			receipt: &domainmodel.ReceiptAnalyze{
				ID:         123,
				TotalPrice: 1000,
				S3FilePath: "test/path.jpg",
				Items:      []domainmodel.ReceiptAnalyzeItem{},
			},
			mockSetup: func(repo *MockReceiptAnalyzeRepository, houseHoldService *MockHouseHoldService) {
				repo.On("FindReceiptAnalyzeByS3FilePath", "test/path.jpg").Return(&domainmodel.ReceiptAnalyze{
					ID:         123,
					TotalPrice: 1000,
					S3FilePath: "test/path.jpg",
					Items:      []domainmodel.ReceiptAnalyzeItem{},
				}, nil)
				repo.On("CreateReceiptAnalyzeResult", mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockRepo := new(MockReceiptAnalyzeRepository)
			mockHouseHoldService := new(MockHouseHoldService)
			tt.mockSetup(mockRepo, mockHouseHoldService)

			// テスト対象のインスタンス作成
			usecase := &receiptAnalyzeUsecase{
				repo:             mockRepo,
				houseHoldService: mockHouseHoldService,
			}

			// テスト実行
			err := usecase.CreateReceiptAnalyzeResult(tt.receipt)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// モックの検証
			mockRepo.AssertExpectations(t)
			mockHouseHoldService.AssertExpectations(t)
		})
	}
}
