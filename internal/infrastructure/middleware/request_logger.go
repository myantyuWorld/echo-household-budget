package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"github.com/labstack/echo/v4"
)

// RequestLoggerMiddleware はリクエスト内容をログ出力するミドルウェア
func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			// ボディを読み込む
			var bodyBytes []byte
			if req.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(req.Body)
			}

			// ログ出力
			log.Printf(`[Request]
	Method: %s
	Path: %s
	Headers: %v
	Body: %s
	`,
				req.Method,
				req.URL.Path,
				req.Header,
				string(bodyBytes))

			// 読み込んだボディを元に戻す
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			return next(c)
		}
	}
}
