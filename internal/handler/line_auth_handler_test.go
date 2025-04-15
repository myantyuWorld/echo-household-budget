package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"echo-household-budget/config"
	domainmodel "echo-household-budget/internal/domain/model"
	serviceMock "echo-household-budget/internal/mock/service"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	appConfig := &config.AppConfig{
		LINELoginFrontendCallbackURL: "http://localhost:5173/line/callback",
	}
	defer ctrl.Finish()

	mockLineAuthService := serviceMock.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{
		lineAuthService: mockLineAuthService,
		appConfig:       appConfig,
	}

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
			req := httptest.NewRequest(http.MethodGet, "/line/login", nil)
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
	ctrl := gomock.NewController(t)
	appConfig := &config.AppConfig{
		LINELoginFrontendCallbackURL: "http://localhost:5173/line/callback",
	}
	defer ctrl.Finish()

	mockLineAuthService := serviceMock.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{
		lineAuthService: mockLineAuthService,
		appConfig:       appConfig,
	}

	tests := []struct {
		name           string
		code           string
		setupMock      func()
		expectedStatus int
		expectedURL    string
	}{
		{
			name: "successful callback",
			code: "valid_code",
			setupMock: func() {
				mockLineAuthService.EXPECT().Callback(gomock.Any(), "valid_code").Return(nil)
			},
			expectedStatus: http.StatusFound,
			expectedURL:    "http://localhost:5173/line/callback",
		},
		{
			name: "missing code",
			code: "",
			setupMock: func() {
				// コードが空の場合はサービスが呼ばれない
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "callback failed",
			code: "valid_code",
			setupMock: func() {
				mockLineAuthService.EXPECT().Callback(gomock.Any(), "valid_code").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/line/callback?code="+tt.code, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupMock()

			err := handler.Callback(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedURL != "" {
				assert.Equal(t, tt.expectedURL, rec.Header().Get("Location"))
			}
		})
	}
}

func TestAuthHandler_FetchMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	appConfig := &config.AppConfig{
		LINELoginFrontendCallbackURL: "http://localhost:5173/line/callback",
	}
	defer ctrl.Finish()

	tests := []struct {
		name           string
		setupMock      func(*serviceMock.MockLineAuthService)
		setupContext   func(*echo.Context)
		expectedStatus int
		expectedBody   *domainmodel.UserAccount
	}{
		{
			name: "認証済みユーザーの情報取得に成功",
			setupMock: func(m *serviceMock.MockLineAuthService) {
				m.EXPECT().CheckAuth(gomock.Any()).Return(&domainmodel.UserAccount{
					ID:         domainmodel.UserID(1),
					UserID:     "user123",
					Name:       "テストユーザー",
					PictureURL: "https://example.com/picture.jpg",
				}, nil)
			},
			setupContext: func(c *echo.Context) {
				cookie := &http.Cookie{
					Name:  "session",
					Value: "valid_session",
				}
				(*c).Request().AddCookie(cookie)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &domainmodel.UserAccount{
				ID:         domainmodel.UserID(1),
				UserID:     "user123",
				Name:       "テストユーザー",
				PictureURL: "https://example.com/picture.jpg",
			},
		},
		{
			name: "未認証の場合",
			setupMock: func(m *serviceMock.MockLineAuthService) {
				m.EXPECT().CheckAuth(gomock.Any()).Return(nil, errors.New("not logged in"))
			},
			setupContext: func(c *echo.Context) {
				// セッションクッキーなし
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "認証チェックに失敗",
			setupMock: func(m *serviceMock.MockLineAuthService) {
				m.EXPECT().CheckAuth(gomock.Any()).Return(nil, errors.New("failed to check auth"))
			},
			setupContext: func(c *echo.Context) {
				cookie := &http.Cookie{
					Name:  "session",
					Value: "invalid_session",
				}
				(*c).Request().AddCookie(cookie)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := serviceMock.NewMockLineAuthService(ctrl)
			tt.setupMock(mockService)

			handler := NewLineAuthHandler(mockService, appConfig)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/line/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.setupContext != nil {
				tt.setupContext(&c)
			}

			err := handler.FetchMe(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedBody != nil {
				var response domainmodel.UserAccount
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, *tt.expectedBody, response)
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	appConfig := &config.AppConfig{
		LINELoginFrontendCallbackURL: "http://localhost:5173/line/callback",
	}
	defer ctrl.Finish()

	mockLineAuthService := serviceMock.NewMockLineAuthService(ctrl)
	handler := &lineAuthHandler{
		lineAuthService: mockLineAuthService,
		appConfig:       appConfig,
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/line/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockLineAuthService.EXPECT().Logout(gomock.Any()).Return(nil)

	err := handler.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
