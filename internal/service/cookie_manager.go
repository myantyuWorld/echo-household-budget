package service

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// クッキー操作
type CookieManager interface {
	SetSessionCookie(c echo.Context, sessionID string) error
	ClearSessionCookie(c echo.Context) error
}

type cookieManager struct{}

func NewCookieManager() CookieManager {
	return &cookieManager{}
}

func (cookieManager *cookieManager) SetSessionCookie(c echo.Context, sessionID string) error {
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func (cookieManager *cookieManager) ClearSessionCookie(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now(),
	})
	return nil
}
