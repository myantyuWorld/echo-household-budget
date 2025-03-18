package shared

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GenerateJWT は、指定されたユーザーIDを持つJWTを生成する
func GenerateJWT(userID string, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(secretKey)
}

// SetAuthCookie はJWTをクッキーに設定する
func SetAuthCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	}
	c.SetCookie(cookie)
}

// ParseJWT は、クッキーからJWTを取得し検証する
func ParseJWT(c echo.Context, secretKey []byte) (*jwt.Token, error) {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return nil, err
	}
	return jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
}
