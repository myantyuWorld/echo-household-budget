package main

import (
	"net/http"
	"template-echo-notion-integration/config"
	"template-echo-notion-integration/internal/handler"
	"template-echo-notion-integration/internal/repository"
	"template-echo-notion-integration/internal/service"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// export XXは、開いているターミナルのみ有効
// export PATH=$PATH:$(go env GOPATH)/bin && air -c .air.toml でホットリロードを有効化
func main() {
	appConfig := config.LoadConfig()
	spew.Dump(appConfig)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConfig.AllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPost, http.MethodDelete},
	}))

	kaimemoRepository := repository.NewNotionRepository(
		appConfig.NotionAPIKey,
		appConfig.NotionKaimemoDatabaseInputID,
		appConfig.NotionKaimemoDatabaseSummaryRecordID,
	)
	kaimemoService := service.NewKaimemoService(kaimemoRepository)
	kaimemoHandler := handler.NewKaimemoHandler(kaimemoService)

	lineRepository := repository.NewLineRepository(appConfig.LINEConfig)
	lineAuthService := service.NewLineAuthService(lineRepository)
	lineAuthHandler := handler.NewLineAuthHandler(lineAuthService, appConfig.LINEConfig)

	kaimemo := e.Group("/kaimemo")
	kaimemo.GET("", kaimemoHandler.FetchKaimemo)
	kaimemo.POST("", kaimemoHandler.CreateKaimemo)
	kaimemo.DELETE("/:id", kaimemoHandler.RemoveKaimemo)

	kaimemo.GET("/ws", kaimemoHandler.WebsocketTelegraph)

	kaimemo.GET("/summary", kaimemoHandler.FetchKaimemoSummaryRecord)
	kaimemo.POST("/summary", kaimemoHandler.CreateKaimemoAmount)
	kaimemo.DELETE("/summary/:id", kaimemoHandler.RemoveKaimemoAmount)

	lineAuth := e.Group("/line")
	lineAuth.GET("/login", lineAuthHandler.Login)
	lineAuth.GET("/callback", lineAuthHandler.Callback)
	lineAuth.GET("/logout", lineAuthHandler.Logout)
	lineAuth.GET("/me", lineAuthHandler.FetchMe)

	port := "3000"
	e.Logger.Fatal(e.Start(":" + port))
}
