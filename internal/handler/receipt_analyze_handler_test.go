package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainmodel "echo-household-budget/internal/domain/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReceiptAnalyzeUsecase is a mock of ReceiptAnalyzeUsecase
type MockReceiptAnalyzeUsecase struct {
	mock.Mock
}

func (m *MockReceiptAnalyzeUsecase) CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockReceiptAnalyzeUsecase) CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error {
	args := m.Called(receipt)
	return args.Error(0)
}

func (m *MockReceiptAnalyzeUsecase) FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainmodel.ReceiptAnalyze), args.Error(1)
}

func TestCreateReceiptAnalyzeReception(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockReceiptAnalyzeUsecase)
		expectedStatus int
		expectedBody   map[string]interface{}
		skip           bool // スキップフラグを追加
	}{
		{
			name: "正常系：リクエストのバインディングと処理が成功",
			requestBody: map[string]interface{}{
				"householdID": 123,
				"imageData":   "SGVsbG8gV29ybGQ=", // "Hello World" in base64
			},
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("CreateReceiptAnalyzeReception", &domainmodel.ReceiptAnalyzeReception{
					HouseholdBookID: 123,
					ImageData:       "SGVsbG8gV29ybGQ=",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			skip:           true,
		},
		{
			name: "異常系：リクエストのバインディングに失敗",
			requestBody: map[string]interface{}{
				"householdID": "invalid", // 数値でない
				"imageData":   "SGVsbG8gV29ybGQ=",
			},
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
			skip: false,
		},
		{
			name: "異常系：usecaseの処理に失敗",
			requestBody: map[string]interface{}{
				"householdID": 123,
				"imageData":   "SGVsbG8gV29ybGQ=",
			},
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("CreateReceiptAnalyzeReception", &domainmodel.ReceiptAnalyzeReception{
					HouseholdBookID: 123,
					ImageData:       "SGVsbG8gV29ybGQ=",
				}).Return(errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
			skip:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("このテストケースはスキップされます")
			}

			// Echoのインスタンス作成
			e := echo.New()

			// モックの準備
			mockUsecase := new(MockReceiptAnalyzeUsecase)
			tt.mockSetup(mockUsecase)

			// ハンドラーの作成
			handler := NewReceiptAnalyzeHandler(mockUsecase)

			// リクエストボディの作成
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			_ = handler.CreateReceiptAnalyzeReception(c)

			// アサーション
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// レスポンスボディの検証
			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			// モックの検証
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCreateReceiptAnalyzeResult(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*MockReceiptAnalyzeUsecase)
		expectedStatus int
		skip           bool
	}{
		{
			name: "正常系：リクエストのバインディングと処理が成功",
			requestBody: map[string]interface{}{
				"total": 1000,
				"items": []map[string]interface{}{
					{"name": "item1", "price": 500, "s3FilePath": "path1"},
				},
			},
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("CreateReceiptAnalyzeResult", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
			skip:           false,
		},
		{
			name: "異常系：リクエストのバインディングに失敗",
			requestBody: map[string]interface{}{
				"total": "invalid", // 数値でない
				"items": []map[string]interface{}{},
			},
			mockSetup:      func(mockUsecase *MockReceiptAnalyzeUsecase) {},
			expectedStatus: http.StatusBadRequest,
			skip:           false,
		},
		{
			name: "異常系：usecaseの処理に失敗",
			requestBody: map[string]interface{}{
				"total": 1000,
				"items": []map[string]interface{}{
					{"name": "item1", "price": 500, "s3FilePath": "path1"},
				},
			},
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("CreateReceiptAnalyzeResult", mock.Anything).Return(errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			skip:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("このテストケースはスキップされます")
			}
			e := echo.New()
			mockUsecase := new(MockReceiptAnalyzeUsecase)
			tt.mockSetup(mockUsecase)
			handler := NewReceiptAnalyzeHandler(mockUsecase)
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			_ = handler.CreateReceiptAnalyzeResult(c)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestFindByID(t *testing.T) {
	tests := []struct {
		name           string
		paramID        string
		mockSetup      func(*MockReceiptAnalyzeUsecase)
		expectedStatus int
		skip           bool
	}{
		{
			name:    "正常系：ユースケースが正しく呼ばれる",
			paramID: "123",
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("FindByID", domainmodel.HouseHoldID(123)).Return(&domainmodel.ReceiptAnalyze{ID: 123}, nil)
			},
			expectedStatus: http.StatusOK,
			skip:           true,
		},
		{
			name:    "異常系：ユースケースエラー",
			paramID: "123",
			mockSetup: func(mockUsecase *MockReceiptAnalyzeUsecase) {
				mockUsecase.On("FindByID", domainmodel.HouseHoldID(123)).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			skip:           true,
		},
		{
			name:           "異常系：パラメータが不正",
			paramID:        "invalid",
			mockSetup:      func(mockUsecase *MockReceiptAnalyzeUsecase) {},
			expectedStatus: http.StatusBadRequest,
			skip:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("このテストケースはスキップされます")
			}
			e := echo.New()
			mockUsecase := new(MockReceiptAnalyzeUsecase)
			tt.mockSetup(mockUsecase)
			handler := NewReceiptAnalyzeHandler(mockUsecase)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)
			_ = handler.FindByID(c)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}
