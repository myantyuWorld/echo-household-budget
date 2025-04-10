package main

import (
	"fmt"
	"net/http"
	"template-echo-notion-integration/config"
	"template-echo-notion-integration/internal/handler"
	"template-echo-notion-integration/internal/infrastructure/persistence/repository"
	"template-echo-notion-integration/internal/middleware"
	"template-echo-notion-integration/internal/service"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// export XXは、開いているターミナルのみ有効
// export PATH=$PATH:$(go env GOPATH)/bin && air -c .air.toml でホットリロードを有効化
func main() {
	// 設定の読み込み
	appConfig := config.LoadConfig()

	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     appConfig.AllowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))
	e.Use(middleware.ErrorHandler())

	// リポジトリの初期化
	kaimemoRepository := repository.NewNotionRepository(
		appConfig.NotionAPIKey,
		appConfig.NotionKaimemoDatabaseInputID,
		appConfig.NotionKaimemoDatabaseSummaryRecordID,
	)
	lineRepository := repository.NewLineRepository(appConfig.LINEConfig)

	// サービスの初期化
	kaimemoService := service.NewKaimemoService(kaimemoRepository)
	lineAuthService := service.NewLineAuthService(lineRepository)

	// ハンドラーの初期化
	kaimemoHandler := handler.NewKaimemoHandler(kaimemoService)
	lineAuthHandler := handler.NewLineAuthHandler(lineAuthService, appConfig.LINEConfig)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 買い物メモ関連のエンドポイント
	kaimemo := e.Group("/kaimemo")
	kaimemo.GET("", kaimemoHandler.FetchKaimemo)
	kaimemo.POST("", kaimemoHandler.CreateKaimemo)
	kaimemo.DELETE("/:id", kaimemoHandler.RemoveKaimemo)
	kaimemo.GET("/ws", kaimemoHandler.WebsocketTelegraph)
	kaimemo.GET("/summary", kaimemoHandler.FetchKaimemoSummaryRecord)
	kaimemo.POST("/summary", kaimemoHandler.CreateKaimemoAmount)
	kaimemo.DELETE("/summary/:id", kaimemoHandler.RemoveKaimemoAmount)

	// LINE認証関連のエンドポイント
	lineAuth := e.Group("/line")
	lineAuth.GET("/login", lineAuthHandler.Login)
	lineAuth.GET("/callback", lineAuthHandler.Callback)
	lineAuth.POST("/logout", lineAuthHandler.Logout)
	lineAuth.GET("/me", lineAuthHandler.FetchMe)

	// サーバーの起動
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appConfig.Port)))
}
