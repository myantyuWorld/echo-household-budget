//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package service

import (
	"errors"
	"fmt"
	"template-echo-notion-integration/internal/domain/household"
	"template-echo-notion-integration/internal/domain/service"
	"template-echo-notion-integration/internal/infrastructure/persistence/repository"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

type LineAuthService interface {
	Login(c echo.Context) (string, error)
	Logout(c echo.Context) error
	Callback(c echo.Context, code string) error
	CheckAuth(c echo.Context) (*household.UserAccount, error)
}

type lineAuthService struct {
	repository            repository.LineRepository
	userAccountRepository household.UserAccountRepository
	sessionManager        SessionManager
	cookieManager         CookieManager
	userAccountService    service.UserAccountService
}

func NewLineAuthService(repository repository.LineRepository, userAccountRepository household.UserAccountRepository, userAccountService service.UserAccountService) LineAuthService {
	return &lineAuthService{
		repository:            repository,
		userAccountRepository: userAccountRepository,
		sessionManager:        NewSessionManager(),
		cookieManager:         NewCookieManager(),
		userAccountService:    userAccountService,
	}
}

// Callback implements LineAuthService.
func (l *lineAuthService) Callback(c echo.Context, code string) error {
	if !l.repository.MatchState(c.QueryParam("state")) {
		return errors.New("state does not match")
	}

	userInfo, err := l.repository.GetUserInfo(c.FormValue("code"))
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	spew.Dump(userInfo)
	lineUserInfo := household.NewLINEUserInfo(household.LINEUserID(userInfo.UserID), userInfo.DisplayName, userInfo.PictureURL)
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
	spew.Dump(sessionID)

	if err := l.cookieManager.SetSessionCookie(c, sessionID); err != nil {
		return errors.New("failed to set session cookie")
	}

	return nil
}

func (l *lineAuthService) CheckAuth(c echo.Context) (*household.UserAccount, error) {
	cookie, err := c.Cookie("session")
	if err != nil {
		return nil, errors.New("not logged in")
	}
	userID, err := l.sessionManager.GetSession(cookie.Value)
	if err != nil {
		return nil, errors.New("session invalid")
	}

	userAaccount, err := l.userAccountRepository.FindByLINEUserID(household.LINEUserID(userID))
	fmt.Println("===============")
	spew.Dump(userAaccount)
	fmt.Println("===============")
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
