package config

import (
	"echo-household-budget/internal/shared/errors"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	"log"
)

// Config はアプリケーションの設定を表す構造体
type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	LINE     LINEConfig     `validate:"required"`
	Notion   NotionConfig   `validate:"required"`
	S3       S3Config       `validate:"required"`
}

// ServerConfig はサーバー関連の設定
type ServerConfig struct {
	Port         string   `validate:"required"`
	AllowOrigins []string `validate:"required"`
}

// DatabaseConfig はデータベース関連の設定
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DBName   string `validate:"required"`
	SSLMode  string `validate:"required"`
}

// LINEConfig はLINE関連の設定
type LINEConfig struct {
	ChannelID     string `validate:"required"`
	ChannelSecret string `validate:"required"`
	RedirectURI   string `validate:"required"`
}

// NotionConfig はNotion関連の設定
type NotionConfig struct {
	APIKey                   string `validate:"required"`
	KaimemoDatabaseInputID   string `validate:"required"`
	KaimemoDatabaseSummaryID string `validate:"required"`
}

type S3Config struct {
	BucketName      string `validate:"required"`
	Region          string `validate:"required"`
	AccessKeyID     string `validate:"required"`
	SecretAccessKey string `validate:"required"`
}

// validateRequiredEnvVars は必須環境変数の存在をチェックする
func validateRequiredEnvVars(config *Config) error {
	if config.LINE.ChannelID == "" {
		return errors.NewAppError(
			errors.ErrorCodeInvalidInput,
			"LINE_CHANNEL_ID is required",
			nil,
		)
	}
	if config.LINE.ChannelSecret == "" {
		return errors.NewAppError(
			errors.ErrorCodeInvalidInput,
			"LINE_CHANNEL_SECRET is required",
			nil,
		)
	}
	if config.Notion.APIKey == "" {
		return errors.NewAppError(
			errors.ErrorCodeInvalidInput,
			"NOTION_API_KEY is required",
			nil,
		)
	}
	return nil
}

type AppConfig struct {
	Port                                 string
	NotionAPIKey                         string
	NotionKaimemoDatabaseInputID         string
	NotionKaimemoDatabaseSummaryRecordID string
	AllowOrigins                         []string
	LINEConfig                           *oauth2.Config
	LINELoginFrontendCallbackURL         string
	DatabaseConfig                       *DatabaseConfig
	S3Config                             *S3Config
}

func LoadConfig() *AppConfig {
	// HACK : 本番デプロイ時には、コメントアウトすること
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// LINE OAuth2設定
	lineConfig := &oauth2.Config{
		ClientID:     os.Getenv("LINE_CHANNEL_ID"),
		ClientSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		RedirectURL:  os.Getenv("LINE_REDIRECT_URI"),
		Scopes:       []string{"profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
		},
	}

	// データベース設定
	dbConfig := &DatabaseConfig{
		Host:     getEnvWithDefault("DB_HOST", "localhost"),
		Port:     getEnvWithDefault("DB_PORT", "35432"),
		User:     getEnvWithDefault("DB_USER", "postgres"),
		Password: getEnvWithDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvWithDefault("DB_NAME", "echo-household-budget"),
		SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
	}

	// S3設定
	s3Config := &S3Config{
		BucketName:      getEnvWithDefault("S3_BUCKET_NAME", ""),
		Region:          getEnvWithDefault("S3_REGION", ""),
		AccessKeyID:     getEnvWithDefault("S3_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnvWithDefault("S3_SECRET_ACCESS_KEY", ""),
	}
	return &AppConfig{
		Port:                                 getEnvWithDefault("PORT", "3000"),
		NotionAPIKey:                         getEnvWithDefault("NOTION_API_KEY", ""),
		NotionKaimemoDatabaseInputID:         getEnvWithDefault("NOTION_KAIMEMO_DB_INPUT_ID", ""),
		NotionKaimemoDatabaseSummaryRecordID: getEnvWithDefault("NOTION_KAIMEMO_DB_SUMMARY_ID", ""),
		AllowOrigins:                         []string{os.Getenv("ALLOW_ORIGINS"), "https://access.line.me/oauth2/v2.1/authorize"},
		LINEConfig:                           lineConfig,
		LINELoginFrontendCallbackURL:         os.Getenv("LINE_LOGIN_FRONTEND_CALLBACK_URL"),
		DatabaseConfig:                       dbConfig,
		S3Config:                             s3Config,
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewDBConnection はデータベース接続を作成する
func NewDBConnection(config *DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}
