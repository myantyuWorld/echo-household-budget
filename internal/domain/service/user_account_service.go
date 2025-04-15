//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainservice

import (
	domainmodel "echo-household-budget/internal/domain/model"

	"gorm.io/gorm"
)

type UserAccountService interface {
	CreateUserAccount(lineUserInfo *domainmodel.LINEUserInfo) error
	IsDuplicateUserAccount(lineUserID domainmodel.LINEUserID) (bool, error)
}

type userAccountService struct {
	repository domainmodel.UserAccountRepository
}

// IsDuplicateUserAccount implements UserAccountService.
func (s *userAccountService) IsDuplicateUserAccount(lineUserID domainmodel.LINEUserID) (bool, error) {
	account, err := s.repository.FindByLINEUserID(lineUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return account != nil, nil
}

func (s *userAccountService) CreateUserAccount(lineUserInfo *domainmodel.LINEUserInfo) error {
	userAccount := domainmodel.NewUserAccount(lineUserInfo)
	return s.repository.Create(userAccount)
}

func NewUserAccountService(repository domainmodel.UserAccountRepository) UserAccountService {
	return &userAccountService{
		repository: repository,
	}
}
