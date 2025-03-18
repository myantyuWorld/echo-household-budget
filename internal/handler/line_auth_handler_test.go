package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	service "template-echo-notion-integration/internal/mock/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLineAuthService := service.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{lineAuthService: mockLineAuthService}

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedURL    string
	}{
		{
			name: "successful login",
			setupMock: func() {
				mockLineAuthService.EXPECT().Login(gomock.Any()).Return("https://line.auth/redirect", nil)
			},
			expectedStatus: http.StatusFound,
			expectedURL:    "https://line.auth/redirect",
		},
		{
			name: "login service error",
			setupMock: func() {
				mockLineAuthService.EXPECT().Login(gomock.Any()).Return("", assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/login", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupMock()

			err := handler.Login(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedURL != "" {
				assert.Equal(t, tt.expectedURL, rec.Header().Get("Location"))
			}
		})
	}
}
func TestAuthHandler_Callback(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		expectedStatus int
	}{
		{
			name:           "successful callback",
			code:           "valid_code",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing code",
			code:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/callback?code="+tt.code, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := NewLineAuthHandler()
			err := handler.Callback(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestAuthHandler_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLineAuthService := service.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{lineAuthService: mockLineAuthService}

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "authenticated user",
			setupMock: func() {
				mockLineAuthService.EXPECT().CheckAuth(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthenticated user",
			setupMock: func() {
				mockLineAuthService.EXPECT().CheckAuth(gomock.Any()).Return(assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/check-auth", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupMock()

			err := handler.CheckAuth(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLineAuthService := service.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{lineAuthService: mockLineAuthService}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLineAuthService.EXPECT().Logout(gomock.Any()).Times(1)

	err := handler.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
