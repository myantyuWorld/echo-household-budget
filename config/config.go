package config

import (
	"log"
	"os"

	"golang.org/x/oauth2"
)

type AppConfig struct {
	Port                                 string
	NotionAPIKey                         string
	NotionKaimemoDatabaseInputID         string
	NotionKaimemoDatabaseSummaryRecordID string
	AllowOrigins                         []string
	// LINEConfig                           *LINEConfig
	LINEConfig *oauth2.Config
}

type LINEConfig struct {
	ClientID     string
	ClientSecret string
	JwtSecret    string
	State        string
	RedirectURI  string
}

func LoadConfig() *AppConfig {
	port := "3000"
	// HACK : productionなら、.envを読み込まない設定にしたい
	// if err := dotenv.Load(); err != nil {
	// 	log.Fatalln(err)
	// }

	apiKey := os.Getenv("NOTION_API_KEY")
	if apiKey == "" {
		log.Fatal("NOTION_API_KEY is not set")
	}

	notionKaimemoDatabaseInputID := os.Getenv("NOTION_DATABASE_KAIMEMO_INPUT")
	if notionKaimemoDatabaseInputID == "" {
		log.Fatal("NOTION_DATABASE_KAIMEMO_INPUT is not set")
	}

	notionKaimemoDatabaseSummaryRecordID := os.Getenv("NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD")
	if notionKaimemoDatabaseSummaryRecordID == "" {
		log.Fatal("NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD is not set")
	}

	frontEndUrl := os.Getenv("FRONTEND_URL")
	if frontEndUrl == "" {
		log.Fatal("FRONTEND_URL is not set")
	}

	lineClientID := os.Getenv("LINE_CLIENT_ID")
	if lineClientID == "" {
		log.Fatal("LINE_CLIENT_ID is not set")
	}
	lineClientSecret := os.Getenv("LINE_CLIENT_SECRET")
	if lineClientSecret == "" {
		log.Fatal("LINE_CLIENT_SECRET is not set")
	}
	lineJwtSecret := os.Getenv("LINE_JWT_SECRET")
	if lineJwtSecret == "" {
		log.Fatal("LINE_JWT_SECRET is not set")
	}
	lineState := os.Getenv("LINE_STATE")
	if lineState == "" {
		log.Fatal("LINE_STATE is not set")
	}

	lineRedirectURI := os.Getenv("LINE_REDIRECT_URI")
	if lineRedirectURI == "" {
		log.Fatal("LINE_REDIRECT is not set")
	}
	// lineTokenURL := os.Getenv("LINE_TOKEN_URL")
	// if lineTokenURL == "" {
	// 	log.Fatal("LINE_TOKEN_URL is not set")
	// }

	return &AppConfig{
		Port:                                 port,
		NotionAPIKey:                         apiKey,
		NotionKaimemoDatabaseInputID:         notionKaimemoDatabaseInputID,
		NotionKaimemoDatabaseSummaryRecordID: notionKaimemoDatabaseSummaryRecordID,
		AllowOrigins: []string{
			"http://localhost:5173", "http://localhost:4173", frontEndUrl,
		},
		// LINEConfig: &LINEConfig{
		// 	ClientID:     lineClientID,
		// 	ClientSecret: lineClientSecret,
		// 	JwtSecret:    lineJwtSecret,
		// 	State:        lineState,
		// 	RedirectURI:  lineRedirectURI,
		// },
		LINEConfig: &oauth2.Config{
			ClientID:     lineClientID,
			ClientSecret: lineClientSecret,
			RedirectURL:  lineRedirectURI,
			Scopes:       []string{"profile", "openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
				TokenURL: "https://api.line.me/oauth2/v2.1/token",
			},
		},
	}
}
