package middleware

import (
	"context"
	"fmt"
	"net/http"

	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/usecase"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

// UserContextKey はコンテキストにユーザー情報を格納する際のキー
type UserContextKey string

const (
	// UserKey はコンテキストにユーザー情報を格納する際のキー
	UserKey UserContextKey = "user"
)

// AuthMiddleware はセッションからユーザー情報を取得してコンテキストに渡すミドルウェア
func AuthMiddleware(sessionManager usecase.SessionManager, userAccountRepository domainmodel.UserAccountRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("===============")
			fmt.Println("AuthMiddleware")
			fmt.Println("===============")
			// セッションからLINEユーザーIDを取得
			cookie, err := c.Cookie("session")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "middleware not logged in",
				})
			}
			spew.Dump(cookie)

			lineUserID, err := sessionManager.GetSession(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "middleware session invalid",
				})
			}

			// ユーザーアカウントを取得
			userAccount, err := userAccountRepository.FindByLINEUserID(domainmodel.LINEUserID(lineUserID))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "middleware user not found",
				})
			}

			// コンテキストにユーザー情報を格納
			ctx := context.WithValue(c.Request().Context(), UserKey, userAccount)
			c.SetRequest(c.Request().WithContext(ctx))
			fmt.Println("===============")
			return next(c)
		}
	}
}

// GetUserFromContext はコンテキストからユーザー情報を取得するヘルパー関数
func GetUserFromContext(ctx context.Context) (*domainmodel.UserAccount, bool) {
	user, ok := ctx.Value(UserKey).(*domainmodel.UserAccount)
	return user, ok
}
