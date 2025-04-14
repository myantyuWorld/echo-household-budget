//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package handler

import (
	"fmt"
	"net/http"
	"template-echo-notion-integration/config"
	"template-echo-notion-integration/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	Login(c echo.Context) error
	Callback(c echo.Context) error
	FetchMe(c echo.Context) error
	Logout(c echo.Context) error
}

type lineAuthHandler struct {
	lineAuthService service.LineAuthService
	appConfig       *config.AppConfig
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
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Code is required"})
	}

	err := a.lineAuthService.Callback(c, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("Callback Failed: %v", err))
	}

	return c.Redirect(http.StatusFound, a.appConfig.LINELoginFrontendCallbackURL)
}

func (a *lineAuthHandler) FetchMe(c echo.Context) error {
	userInfo, err := a.lineAuthService.CheckAuth(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, fmt.Errorf("CheckAuth Failed: %v", err))
	}

	return c.JSON(http.StatusOK, userInfo)
}

func (a *lineAuthHandler) Logout(c echo.Context) error {
	a.lineAuthService.Logout(c)

	return c.JSON(http.StatusOK, echo.Map{"message": "Logged out"})
}

func NewLineAuthHandler(lineAuthService service.LineAuthService, appConfig *config.AppConfig) AuthHandler {
	return &lineAuthHandler{
		lineAuthService: lineAuthService,
		appConfig:       appConfig,
	}
}
