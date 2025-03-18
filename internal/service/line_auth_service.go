//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package service

import (
	"errors"
	"fmt"
	"template-echo-notion-integration/internal/repository"

	"github.com/labstack/echo/v4"
)

type LineAuthService interface {
	Login(c echo.Context) (string, error)
	Logout(c echo.Context) error
	Callback(c echo.Context, code string) error
	CheckAuth(c echo.Context) error
}

// セッション管理用の仮データ
var sessionStore = make(map[string]string)

type lineAuthService struct {
	repository     repository.LineRepository
	sessionManager SessionManager
	cookieManager  CookieManager
}

func NewLineAuthService(repository repository.LineRepository) LineAuthService {
	return &lineAuthService{
		repository:     repository,
		sessionManager: NewSessionManager(),
		cookieManager:  NewCookieManager(),
	}
}

// Callback implements LineAuthService.
func (l *lineAuthService) Callback(c echo.Context, code string) error {
	if !l.repository.MatchState(c.QueryParam("state")) {
		return errors.New("State does not match")
	}

	userInfo, err := l.repository.GetUserInfo(c.FormValue("code"))
	fmt.Println(userInfo)

	//
	// TODO : システムに登録されていなければ、ユーザー情報をDBに保存する
	//

	sessionID, err := l.sessionManager.CreateSession(userInfo.UserID)
	if err != nil {
		return errors.New("Failed to create session")
	}

	// HACK : リファクタリング後、クッキーに保存できていない。ので、checkAuthで取得に失敗しているよう
	if err := l.cookieManager.SetSessionCookie(c, sessionID); err != nil {
		return errors.New("Failed to set session cookie")
	}

	// TODO : Redirectでフロントエンドに戻す
	return nil
}

func (l *lineAuthService) CheckAuth(c echo.Context) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		return errors.New("Not logged in")
	}
	userID, exists := sessionStore[cookie.Value]
	if !exists {
		return errors.New("Session invalid")
	}

	// TODO : userIDをもとに、ユーザー情報を取得して返す
	fmt.Println(userID)

	return nil
}

func (l *lineAuthService) Login(c echo.Context) (string, error) {
	return l.repository.GetAuthCodeUrl(), nil
}

func (l *lineAuthService) Logout(c echo.Context) error {
	if err := l.cookieManager.ClearSessionCookie(c); err != nil {
		return errors.New("Failed to clear session cookie")
	}

	return nil
}
