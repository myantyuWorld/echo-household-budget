package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"template-echo-notion-integration/internal/domain/household"
	mockDomainService "template-echo-notion-integration/internal/domain/mock/service"
	mock "template-echo-notion-integration/internal/infrastructure/persistence/mock/repository"
	"template-echo-notion-integration/internal/infrastructure/persistence/repository"
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
		mockSetup     func(*mock.MockLineRepository, *mockDomainService.MockUserAccountService)
		expectedError error
	}{
		{
			name:  "コールバック処理に成功",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(household.LINEUserID("user123")).Return(false, nil)
				us.EXPECT().CreateUserAccount(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "stateが一致しない場合",
			code:  "valid_code",
			state: "invalid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("invalid_state").Return(false)
			},
			expectedError: errors.New("state does not match"),
		},
		{
			name:  "ユーザー情報の取得に失敗",
			code:  "invalid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("invalid_code").Return(nil, errors.New("failed to get user info"))
			},
			expectedError: errors.New("failed to get user info: failed to get user info"),
		},
		{
			name:  "ユーザーアカウントの重複チェックに失敗",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(household.LINEUserID("user123")).Return(false, errors.New("failed to check duplicate"))
			},
			expectedError: errors.New("failed to check if user account exists: failed to check duplicate"),
		},
		{
			name:  "ユーザーアカウントの作成に失敗",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(household.LINEUserID("user123")).Return(false, nil)
				us.EXPECT().CreateUserAccount(gomock.Any()).Return(errors.New("failed to create user account"))
			},
			expectedError: errors.New("failed to create user account: failed to create user account"),
		},
		{
			name:  "既存のユーザーアカウントの場合",
			code:  "valid_code",
			state: "valid_state",
			mockSetup: func(lr *mock.MockLineRepository, us *mockDomainService.MockUserAccountService) {
				lr.EXPECT().MatchState("valid_state").Return(true)
				lr.EXPECT().GetUserInfo("valid_code").Return(&repository.UserInfo{
					UserID:      "user123",
					DisplayName: "テストユーザー",
					PictureURL:  "https://example.com/photo.jpg",
				}, nil)
				us.EXPECT().IsDuplicateUserAccount(household.LINEUserID("user123")).Return(true, nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockLineRepository(ctrl)
			mockUserService := mockDomainService.NewMockUserAccountService(ctrl)
			tt.mockSetup(mockRepo, mockUserService)

			service := NewLineAuthService(mockRepo)
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
