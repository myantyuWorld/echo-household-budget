package main

import (
	"echo-household-budget/config"
	domainService "echo-household-budget/internal/domain/service"
	"echo-household-budget/internal/handler"
	"echo-household-budget/internal/infrastructure/middleware"
	"echo-household-budget/internal/infrastructure/persistence/repository"
	"echo-household-budget/internal/usecase"
	"fmt"
	"net/http"

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

	// データベース接続の設定
	db, err := config.NewDBConnection(appConfig.DatabaseConfig)
	if err != nil {
		e.Logger.Fatal(err)
	}

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
	userAccountRepository := repository.NewUserAccountRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	houseHoldRepository := repository.NewHouseHoldRepository(db)
	shoppingRepository := repository.NewShoppingRepository(db)
	userAccountService := domainService.NewUserAccountService(userAccountRepository, categoryRepository, houseHoldRepository)

	// サービスの初期化
	sessionManager := usecase.NewSessionManager()
	kaimemoService := usecase.NewKaimemoService(kaimemoRepository)
	shoppingUsecase := usecase.NewShoppingUsecase(shoppingRepository)
	lineAuthService := usecase.NewLineAuthService(lineRepository, userAccountRepository, userAccountService, sessionManager)

	// ハンドラーの初期化
	kaimemoHandler := handler.NewKaimemoHandler(kaimemoService, shoppingUsecase)
	lineAuthHandler := handler.NewLineAuthHandler(lineAuthService, appConfig)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// 買い物メモ関連のエンドポイント
	kaimemo := e.Group("/kaimemo", middleware.AuthMiddleware(sessionManager, userAccountRepository))
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
	// lineAuth.GET("/me", lineAuthHandler.FetchMe, middleware.AuthMiddleware(sessionManager, userAccountRepository))
	lineAuth.GET("/me", lineAuthHandler.FetchMe)

	// サーバーの起動
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appConfig.Port)))
}
