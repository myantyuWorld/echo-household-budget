package usecase

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	mockUserAccountRepository "echo-household-budget/internal/domain/mock/domainmodel"
	mockDomainService "echo-household-budget/internal/domain/mock/domainservice"
	domainmodel "echo-household-budget/internal/domain/model"
	mock "echo-household-budget/internal/infrastructure/persistence/mock/repository"
	"echo-household-budget/internal/infrastructure/persistence/repository"
)

func setupEchoContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestLineAuthService_Callback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		code          string
		state         string
		mockSetup     func(*mock.MockLineRepository, *mockDomainService.MockUserAccountService, *mockUserAccountRepository.MockUserAccountRepository)
		expectedError error
	}{
		{
			name:  "コールバック処理に成功",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(domainmodel.LINEUserID("user123")).Return(false, nil)
				us.EXPECT().CreateUserAccount(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "stateが一致しない場合",
			code:  "valid_code",
			state: "invalid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("invalid_state").Return(false)
			},
			expectedError: errors.New("state does not match"),
		},
		{
			name:  "ユーザー情報の取得に失敗",
			code:  "invalid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("invalid_code").Return(nil, errors.New("failed to get user info"))
			},
			expectedError: errors.New("failed to get user info: failed to get user info"),
		},
		{
			name:  "ユーザーアカウントの重複チェックに失敗",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(domainmodel.LINEUserID("user123")).Return(false, errors.New("failed to check duplicate"))
			},
			expectedError: errors.New("failed to check if user account exists: failed to check duplicate"),
		},
		{
			name:  "ユーザーアカウントの作成に失敗",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(domainmodel.LINEUserID("user123")).Return(false, nil)
				us.EXPECT().CreateUserAccount(gomock.Any()).Return(errors.New("failed to create user account"))
			},
			expectedError: errors.New("failed to create user account: failed to create user account"),
		},
		{
			name:  "既存のユーザーアカウントの場合",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService, ua *mockUserAccountRepository.MockUserAccountRepository) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(domainmodel.LINEUserID("user123")).Return(true, nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockLineRepository(ctrl)
			mockUserAccountRepository := mockUserAccountRepository.NewMockUserAccountRepository(ctrl)
			mockUserService := mockDomainService.NewMockUserAccountService(ctrl)
			tt.mockSetup(mockRepo, mockUserService, mockUserAccountRepository)

			service := NewLineAuthService(mockRepo, mockUserAccountRepository, mockUserService)
			service.(*lineAuthService).userAccountService = mockUserService

			c := setupEchoContext()
			c.QueryParams().Set("state", tt.state)
			c.Request().Form = map[string][]string{
				"code": {tt.code},
			}

			err := service.Callback(c, tt.code)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLineAuthService_CheckAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		setupMock     func(*mockUserAccountRepository.MockUserAccountRepository)
		setupContext  func(*echo.Context)
		expectedUser  *domainmodel.UserAccount
		expectedError error
	}{
		{
			name: "認証済みユーザーの情報取得に成功",
			setupMock: func(ua *mockUserAccountRepository.MockUserAccountRepository) {
				ua.EXPECT().FindByLINEUserID(domainmodel.LINEUserID("user123")).Return(&domainmodel.UserAccount{
					ID:         1,
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
			expectedUser: &domainmodel.UserAccount{
				ID:         1,
				UserID:     "user123",
				Name:       "テストユーザー",
				PictureURL: "https://example.com/picture.jpg",
			},
			expectedError: nil,
		},
		{
			name:      "セッションクッキーが存在しない場合",
			setupMock: func(ua *mockUserAccountRepository.MockUserAccountRepository) {},
			setupContext: func(c *echo.Context) {
				// セッションクッキーなし
			},
			expectedUser:  nil,
			expectedError: errors.New("not logged in"),
		},
		{
			name: "ユーザーアカウントが見つからない場合",
			setupMock: func(ua *mockUserAccountRepository.MockUserAccountRepository) {
				ua.EXPECT().FindByLINEUserID(domainmodel.LINEUserID("user123")).Return(nil, gorm.ErrRecordNotFound)
			},
			setupContext: func(c *echo.Context) {
				cookie := &http.Cookie{
					Name:  "session",
					Value: "valid_session",
				}
				(*c).Request().AddCookie(cookie)
			},
			expectedUser:  nil,
			expectedError: errors.New("failed to find user account"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockLineRepository(ctrl)
			mockUserAccountRepository := mockUserAccountRepository.NewMockUserAccountRepository(ctrl)
			mockUserService := mockDomainService.NewMockUserAccountService(ctrl)

			service := NewLineAuthService(mockRepo, mockUserAccountRepository, mockUserService)
			lineService := service.(*lineAuthService)

			// セッションマネージャーのモック設定
			lineService.sessionManager = &mockSessionManager{
				getSessionFunc: func(sessionID string) (string, error) {
					if sessionID == "valid_session" {
						return "user123", nil
					}
					return "", errors.New("session invalid")
				},
			}

			tt.setupMock(mockUserAccountRepository)

			c := setupEchoContext()
			if tt.setupContext != nil {
				tt.setupContext(&c)
			}

			user, err := service.CheckAuth(c)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedUser != nil {
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

type mockSessionManager struct {
	getSessionFunc func(string) (string, error)
}

func (m *mockSessionManager) GetSession(sessionID string) (string, error) {
	return m.getSessionFunc(sessionID)
}

func (m *mockSessionManager) CreateSession(userID string) (string, error) {
	return "", nil
}

func (m *mockSessionManager) DestroySession(sessionID string) error {
	return nil
}
