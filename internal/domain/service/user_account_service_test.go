package domainservice

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock "echo-household-budget/internal/domain/mock/domainmodel"
	domainmodel "echo-household-budget/internal/domain/model"
)

func TestUserAccountService_CreateUserAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		lineUserInfo  *domainmodel.LINEUserInfo
		mockSetup     func(*mock.MockUserAccountRepository)
		expectedError error
	}{
		{
			name: "新規ユーザーアカウントの作成に成功",
			lineUserInfo: &domainmodel.LINEUserInfo{
				UserID:      domainmodel.LINEUserID("user123"),
				DisplayName: "テストユーザー",
				PictureURL:  "https://example.com/photo.jpg",
			},
			mockSetup: func(m *mock.MockUserAccountRepository) {
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserAccountRepository(ctrl)
			tt.mockSetup(mockRepo)

			service := NewUserAccountService(mockRepo)
			err := service.CreateUserAccount(tt.lineUserInfo)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserAccountService_IsDuplicateUserAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		lineUserID     domainmodel.LINEUserID
		mockSetup      func(*mock.MockUserAccountRepository)
		expectedResult bool
		expectedError  error
	}{
		{
			name:       "ユーザーが存在する場合",
			lineUserID: domainmodel.LINEUserID("existing_user"),
			mockSetup: func(m *mock.MockUserAccountRepository) {
				m.EXPECT().FindByLINEUserID(domainmodel.LINEUserID("existing_user")).Return(&domainmodel.UserAccount{
					ID:     1,
					UserID: domainmodel.LINEUserID("existing_user"),
					Name:   "既存ユーザー",
				}, nil)
			},
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:       "ユーザーが存在しない場合",
			lineUserID: domainmodel.LINEUserID("non_existent_user"),
			mockSetup: func(m *mock.MockUserAccountRepository) {
				m.EXPECT().FindByLINEUserID(domainmodel.LINEUserID("non_existent_user")).Return(nil, nil)
			},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:       "エラーが発生した場合",
			lineUserID: domainmodel.LINEUserID("error_user"),
			mockSetup: func(m *mock.MockUserAccountRepository) {
				m.EXPECT().FindByLINEUserID(domainmodel.LINEUserID("error_user")).Return(nil, errors.New("database error"))
			},
			expectedResult: false,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserAccountRepository(ctrl)
			tt.mockSetup(mockRepo)

			service := NewUserAccountService(mockRepo)
			result, err := service.IsDuplicateUserAccount(tt.lineUserID)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
