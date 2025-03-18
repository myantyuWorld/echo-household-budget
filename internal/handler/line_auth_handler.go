//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package handler

import (
	"fmt"
	"net/http"
	"template-echo-notion-integration/internal/service"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type AuthHandler interface {
	Login(c echo.Context) error
	Callback(c echo.Context) error
	FetchMe(c echo.Context) error
	Logout(c echo.Context) error
}

type lineAuthHandler struct {
	lineAuthService service.LineAuthService
	lineConfig      *oauth2.Config
}

// [Go言語]LINE ログイン連携方法 メモ | https://qiita.com/KWS_0901/items/8c4accdda43bc9f26a57
// Login implements AuthHandler.
func (a *lineAuthHandler) Login(c echo.Context) error {
	url, err := a.lineAuthService.Login(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Login Failed"})
	}
	return c.Redirect(http.StatusFound, url)
}

// Callback implements AuthHandler.
func (a *lineAuthHandler) Callback(c echo.Context) error {
	err := a.lineAuthService.Callback(c, c.QueryParam("code"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("Callback Failed: %v", err))
	}

	// TODO : 下記はフロントエンド未実装のため、仮実装。本来はRedirectでフロントエンドのホーム画面にルーティングする
	return c.JSON(http.StatusOK, map[string]string{"message": "Callback Success"})
}

func (a *lineAuthHandler) FetchMe(c echo.Context) error {
	// TODO : userInfo, errを返すように修正する
	err := a.lineAuthService.CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Not logged in"})
	}

	// TODO : Service層でユーザー情報を取得して返す
	return c.JSON(http.StatusOK, map[string]string{"user": "userID"})
}

func (a *lineAuthHandler) Logout(c echo.Context) error {
	a.lineAuthService.Logout(c)

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out"})
}

func NewLineAuthHandler(lineAuthService service.LineAuthService, lineConfig *oauth2.Config) AuthHandler {
	return &lineAuthHandler{
		lineAuthService: lineAuthService,
		lineConfig:      lineConfig,
	}
}
