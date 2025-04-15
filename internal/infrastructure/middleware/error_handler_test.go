package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "echo-household-budget/internal/shared/errors"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   ErrorResponse
	}{
		{
			name:           "AppError - InvalidInput",
			err:            apperrors.NewAppError(apperrors.ErrorCodeInvalidInput, "invalid input", nil),
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeInvalidInput),
				Message: "invalid input",
			},
		},
		{
			name:           "AppError - NotFound",
			err:            apperrors.NewAppError(apperrors.ErrorCodeNotFound, "not found", nil),
			expectedStatus: http.StatusNotFound,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeNotFound),
				Message: "not found",
			},
		},
		{
			name:           "AppError - Unauthorized",
			err:            apperrors.NewAppError(apperrors.ErrorCodeUnauthorized, "unauthorized", nil),
			expectedStatus: http.StatusUnauthorized,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeUnauthorized),
				Message: "unauthorized",
			},
		},
		{
			name:           "AppError - DatabaseError",
			err:            apperrors.NewAppError(apperrors.ErrorCodeDatabaseError, "database error", nil),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeDatabaseError),
				Message: "database error",
			},
		},
		{
			name:           "AppError - ExternalService",
			err:            apperrors.NewAppError(apperrors.ErrorCodeExternalService, "external service error", nil),
			expectedStatus: http.StatusBadGateway,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeExternalService),
				Message: "external service error",
			},
		},
		{
			name:           "標準エラー",
			err:            errors.New("standard error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: ErrorResponse{
				Code:    string(apperrors.ErrorCodeInternalError),
				Message: "Internal Server Error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoインスタンスの作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// エラーハンドラーミドルウェアの作成
			h := ErrorHandler()

			// エラーを返すハンドラー
			handler := func(c echo.Context) error {
				return tt.err
			}

			// ミドルウェアの実行
			err := h(handler)(c)
			assert.NoError(t, err)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスボディの検証
			var response ErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     apperrors.ErrorCode
		expected int
	}{
		{
			name:     "InvalidInput",
			code:     apperrors.ErrorCodeInvalidInput,
			expected: http.StatusBadRequest,
		},
		{
			name:     "NotFound",
			code:     apperrors.ErrorCodeNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "Unauthorized",
			code:     apperrors.ErrorCodeUnauthorized,
			expected: http.StatusUnauthorized,
		},
		{
			name:     "DatabaseError",
			code:     apperrors.ErrorCodeDatabaseError,
			expected: http.StatusInternalServerError,
		},
		{
			name:     "ExternalService",
			code:     apperrors.ErrorCodeExternalService,
			expected: http.StatusBadGateway,
		},
		{
			name:     "Unknown",
			code:     "UNKNOWN",
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := getHTTPStatus(tt.code)
			assert.Equal(t, tt.expected, status)
		})
	}
}
