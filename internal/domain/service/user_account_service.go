//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package domainservice

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"fmt"

	"gorm.io/gorm"
)

type UserAccountService interface {
	CreateUserAccount(lineUserInfo *domainmodel.LINEUserInfo) error
	IsDuplicateUserAccount(lineUserID domainmodel.LINEUserID) (bool, error)
	FetchUserAccount(userID domainmodel.UserID) (*domainmodel.UserAccount, error)
	FetchAllUserAccount() ([]*domainmodel.UserAccount, error)
}

type userAccountService struct {
	userAccountRepository domainmodel.UserAccountRepository
	categoryRepository    domainmodel.CategoryRepository
	houseHoldRepository   domainmodel.HouseHoldRepository
}

// FetchUserAccount implements UserAccountService.
func (s *userAccountService) FetchUserAccount(userID domainmodel.UserID) (*domainmodel.UserAccount, error) {
	account, err := s.userAccountRepository.FetchMe(userID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// FetchAllUserAccount implements UserAccountService.
func (s *userAccountService) FetchAllUserAccount() ([]*domainmodel.UserAccount, error) {
	return s.userAccountRepository.FetchAll()
}

// IsDuplicateUserAccount implements UserAccountService.
func (s *userAccountService) IsDuplicateUserAccount(lineUserID domainmodel.LINEUserID) (bool, error) {
	account, err := s.userAccountRepository.FindByLINEUserID(lineUserID)
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

	err := s.userAccountRepository.Create(userAccount)
	if err != nil {
		return fmt.Errorf("failed to create user account: %w", err)
	}

	householdBook := domainmodel.NewDefaultHouseHold(userAccount)
	err = s.houseHoldRepository.Create(householdBook)
	if err != nil {
		return fmt.Errorf("failed to create household book: %w", err)
	}

	err = s.houseHoldRepository.CreateUserHouseHold(&domainmodel.UserHouseHold{
		UserID:      userAccount.ID,
		HouseHoldID: householdBook.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create user household book: %w", err)
	}

	for _, categoryLimit := range householdBook.CategoryLimit {
		err = s.categoryRepository.CreateHouseHoldCategory(&domainmodel.CategoryLimit{
			HouseholdBookID: householdBook.ID,
			Category:        categoryLimit.Category,
			LimitAmount:     categoryLimit.LimitAmount,
		})
		if err != nil {
			return fmt.Errorf("failed to create master category: %w", err)
		}
	}

	return nil
}

func NewUserAccountService(userAccountRepository domainmodel.UserAccountRepository, categoryRepository domainmodel.CategoryRepository, houseHoldRepository domainmodel.HouseHoldRepository) UserAccountService {
	return &userAccountService{
		userAccountRepository: userAccountRepository,
		categoryRepository:    categoryRepository,
		houseHoldRepository:   houseHoldRepository,
	}
}
