package middleware

import (
	"echo-household-budget/internal/shared/errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler はグローバルエラーハンドリングミドルウェア
func ErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			fmt.Println("===============")
			fmt.Println("ErrorHandler")
			fmt.Println("===============")
			spew.Dump(err)

			// AppErrorの場合は適切なステータスコードとメッセージを返す
			if appErr, ok := errors.GetAppError(err); ok {
				status := getHTTPStatus(appErr.Code)
				return c.JSON(status, ErrorResponse{
					Code:    string(appErr.Code),
					Message: appErr.Message,
				})
			}

			// その他のエラーは500エラーとして扱う
			return c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    string(errors.ErrorCodeInternalError),
				Message: "Internal Server Error",
			})
		}
	}
}

// getHTTPStatus はエラーコードに対応するHTTPステータスコードを返す
func getHTTPStatus(code errors.ErrorCode) int {
	switch code {
	case errors.ErrorCodeInvalidInput:
		return http.StatusBadRequest
	case errors.ErrorCodeNotFound:
		return http.StatusNotFound
	case errors.ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case errors.ErrorCodeDatabaseError:
		return http.StatusInternalServerError
	case errors.ErrorCodeExternalService:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
