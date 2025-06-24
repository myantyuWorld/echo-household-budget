package main

import (
	"echo-household-budget/config"
	"echo-household-budget/internal/infrastructure/middleware"
	"echo-household-budget/internal/setup"
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
	// spew.Dump(appConfig)

	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	setupMiddleware(e, appConfig)

	// 依存関係の初期化
	dependencies := setup.NewDependencies(appConfig)

	// ルーティングの設定
	setupRoutes(e, dependencies)

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appConfig.Port)))
}

func setupMiddleware(e *echo.Echo, appConfig *config.AppConfig) {
	e.Use(middleware.RequestLoggerMiddleware())
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     appConfig.AllowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))
	e.Use(middleware.ErrorHandler())
}

func setupRoutes(e *echo.Echo, deps *setup.Dependencies) {
	// 買い物メモ関連のエンドポイント
	kaimemo := e.Group("/kaimemo", middleware.AuthMiddleware(deps.SessionManager, deps.UserAccountRepository))
	kaimemo.GET("", deps.KaimemoHandler.FetchKaimemo)
	kaimemo.POST("", deps.KaimemoHandler.CreateKaimemo)
	kaimemo.DELETE("/:id", deps.KaimemoHandler.RemoveKaimemo)
	kaimemo.GET("/ws", deps.KaimemoHandler.WebsocketTelegraph)
	kaimemo.GET("/summary", deps.KaimemoHandler.FetchKaimemoSummaryRecord)
	kaimemo.POST("/summary", deps.KaimemoHandler.CreateKaimemoAmount)
	kaimemo.DELETE("/summary/:id", deps.KaimemoHandler.RemoveKaimemoAmount)

	// 家計簿関連のエンドポイント
	houseHold := e.Group("/household", middleware.AuthMiddleware(deps.SessionManager, deps.UserAccountRepository))
	houseHold.GET("/:id", deps.HouseHoldHandler.FetchHouseHold)
	houseHold.GET("/user/:id", deps.HouseHoldHandler.FetchHouseHoldUser)
	houseHold.POST("/user/:id", deps.HouseHoldHandler.AddHouseHold)
	houseHold.POST("/:householdID/share/:inviteUserID", deps.HouseHoldHandler.ShareHouseHold)
	houseHold.POST("/:householdID/category", deps.HouseHoldHandler.AddHouseHoldCategory)
	houseHold.GET("/:householdID/shopping/record", deps.HouseHoldHandler.FetchShoppingRecord)
	houseHold.POST("/:householdID/shopping/record", deps.HouseHoldHandler.CreateShoppingRecord)
	houseHold.DELETE("/:householdID/shopping/record/:shoppingID", deps.HouseHoldHandler.RemoveShoppingRecord)

	// LINE認証関連のエンドポイント
	lineAuth := e.Group("/line")
	lineAuth.GET("/login", deps.LineAuthHandler.Login)
	lineAuth.GET("/callback", deps.LineAuthHandler.Callback)
	lineAuth.POST("/logout", deps.LineAuthHandler.Logout)
	lineAuth.GET("/me", deps.LineAuthHandler.FetchMe)

	// OpenAI関連のエンドポイント
	openAI := e.Group("/openai/analyze")
	openAI.POST("/:householdID/receipt/reception", deps.ReceiptAnalyzeHandler.CreateReceiptAnalyzeReception)
	openAI.POST("/:householdID/receipt/result", deps.ReceiptAnalyzeHandler.CreateReceiptAnalyzeResult)

	// 管理系のエンドポイント
	admin := e.Group("/admin", middleware.AuthMiddleware(deps.SessionManager, deps.UserAccountRepository))
	admin.GET("/informations", deps.FetchInformationHandler.Handle)
	admin.POST("/informations", deps.CreateInformationHandler.Handle)
	admin.DELETE("/informations/:id", deps.DeleteInformationHandler.Handle)
	admin.GET("/informations/:id", deps.FetchInformationDetailHandler.Handle)
	admin.PUT("/informations/:id", deps.PutInformationHandler.Handle)
	admin.POST("/informations/:id/publish", deps.PublishInformationHandler.Handle)

	// ユーザー関連のエンドポイント
	user := e.Group("/user", middleware.AuthMiddleware(deps.SessionManager, deps.UserAccountRepository))
	user.GET("/informations", deps.FetchUserInformationHandler.Handle)
	user.POST("/informations", deps.UpdateReadUserInformationHandler.Handle)

	// チャット関連のエンドポイント
	chat := e.Group("/chat", middleware.AuthMiddleware(deps.SessionManager, deps.UserAccountRepository))
	chat.GET("/messages/ws", deps.ChatMessageTelegraphHandler.WebSocketChat)
}
