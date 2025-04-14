//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package service

import "template-echo-notion-integration/internal/domain/household"

type UserAccountService interface {
	CreateUserAccount(lineUserInfo *household.LINEUserInfo) error
	IsDuplicateUserAccount(lineUserID household.LINEUserID) (bool, error)
}

type userAccountService struct {
	repository household.UserAccountRepository
}

// IsDuplicateUserAccount implements UserAccountService.
func (s *userAccountService) IsDuplicateUserAccount(lineUserID household.LINEUserID) (bool, error) {
	account, err := s.repository.FindByLINEUserID(lineUserID)
	if err != nil {
		return false, err
	}

	return account != nil, nil
}

func (s *userAccountService) CreateUserAccount(lineUserInfo *household.LINEUserInfo) error {
	userAccount := household.NewUserAccount(lineUserInfo.UserID, lineUserInfo.DisplayName)
	return s.repository.Create(userAccount)
}

func NewUserAccountService(repository household.UserAccountRepository) UserAccountService {
	return &userAccountService{
		repository: repository,
	}
}
