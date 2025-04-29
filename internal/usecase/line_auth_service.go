//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	domainservice "echo-household-budget/internal/domain/service"
	"echo-household-budget/internal/infrastructure/persistence/repository"
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type LineAuthService interface {
	Login(c echo.Context) (string, error)
	Logout(c echo.Context) error
	Callback(c echo.Context, code string) error
	CheckAuth(c echo.Context) (*domainmodel.UserAccount, error)
}

type lineAuthService struct {
	repository            repository.LineRepository
	userAccountRepository domainmodel.UserAccountRepository
	sessionManager        SessionManager
	cookieManager         CookieManager
	userAccountService    domainservice.UserAccountService
}

func NewLineAuthService(repository repository.LineRepository, userAccountRepository domainmodel.UserAccountRepository, userAccountService domainservice.UserAccountService, sessionManager SessionManager) LineAuthService {
	return &lineAuthService{
		repository:            repository,
		userAccountRepository: userAccountRepository,
		sessionManager:        sessionManager,
		cookieManager:         NewCookieManager(),
		userAccountService:    userAccountService,
	}
}

// Callback implements LineAuthService.
func (l *lineAuthService) Callback(c echo.Context, code string) error {
	fmt.Println("func (l *lineAuthService) Callback(c echo.Context, code string) error {")
	if !l.repository.MatchState(c.QueryParam("state")) {
		return errors.New("state does not match")
	}

	userInfo, err := l.repository.GetUserInfo(c.FormValue("code"))
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	lineUserInfo := domainmodel.NewLINEUserInfo(domainmodel.LINEUserID(userInfo.UserID), userInfo.DisplayName, userInfo.PictureURL)
	result, err := l.userAccountService.IsDuplicateUserAccount(lineUserInfo.UserID)
	if err != nil {
		return fmt.Errorf("failed to check if user account exists: %w", err)
	}

	if !result {
		err = l.userAccountService.CreateUserAccount(lineUserInfo)
		if err != nil {
			return fmt.Errorf("failed to create user account: %w", err)
		}
	}

	sessionID, err := l.sessionManager.CreateSession(userInfo.UserID)
	if err != nil {
		return errors.New("failed to create session")
	}

	if err := l.cookieManager.SetSessionCookie(c, sessionID); err != nil {
		return errors.New("failed to set session cookie")
	}

	return nil
}

func (l *lineAuthService) CheckAuth(c echo.Context) (*domainmodel.UserAccount, error) {
	cookie, _ := c.Cookie("session")
	if cookie == nil {
		return nil, errors.New("not logged in")
	}
	userID, err := l.sessionManager.GetSession(cookie.Value)
	if err != nil {
		return nil, errors.New("check auth session invalid")
	}

	userAaccount, err := l.userAccountRepository.FindByLINEUserID(domainmodel.LINEUserID(userID))
	if err != nil {
		return nil, errors.New("failed to find user account")
	}

	return userAaccount, nil
}

func (l *lineAuthService) Login(c echo.Context) (string, error) {
	return l.repository.GetAuthCodeUrl(), nil
}

func (l *lineAuthService) Logout(c echo.Context) error {
	if err := l.cookieManager.ClearSessionCookie(c); err != nil {
		return errors.New("failed to clear session cookie")
	}

	return nil
}
