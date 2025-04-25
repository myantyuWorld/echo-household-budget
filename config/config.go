package config

import (
	"echo-household-budget/internal/shared/errors"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config はアプリケーションの設定を表す構造体
type Config struct {
	Server   ServerConfig   `validate:"required"`
	Database DatabaseConfig `validate:"required"`
	LINE     LINEConfig     `validate:"required"`
	Notion   NotionConfig   `validate:"required"`
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

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	// .envファイルの読み込み
	// if err := godotenv.Load(); err != nil {
	// 	return nil, errors.NewAppError(
	// 		errors.ErrorCodeInternalError,
	// 		"Failed to load .env file",
	// 		err,
	// 	)
	// }

	config := &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("PORT", "3000"),
			AllowOrigins: []string{getEnvOrDefault("ALLOW_ORIGINS", "*")},
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			DBName:   getEnvOrDefault("DB_NAME", "kakeibo"),
		},
		LINE: LINEConfig{
			ChannelID:     getEnvOrDefault("LINE_CHANNEL_ID", ""),
			ChannelSecret: getEnvOrDefault("LINE_CHANNEL_SECRET", ""),
			RedirectURI:   getEnvOrDefault("LINE_REDIRECT_URI", ""),
		},
		Notion: NotionConfig{
			APIKey:                   getEnvOrDefault("NOTION_API_KEY", ""),
			KaimemoDatabaseInputID:   getEnvOrDefault("NOTION_KAIMEMO_DB_INPUT_ID", ""),
			KaimemoDatabaseSummaryID: getEnvOrDefault("NOTION_KAIMEMO_DB_SUMMARY_ID", ""),
		},
	}

	// 必須環境変数のチェック
	if err := validateRequiredEnvVars(config); err != nil {
		return nil, err
	}

	return config, nil
}

// getEnvOrDefault は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
}

func LoadConfig() *AppConfig {
	// if err := godotenv.Load(); err != nil {
	// 	log.Printf("Warning: .env file not found")
	// }

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

	return &AppConfig{
		Port:                                 getEnvWithDefault("PORT", "3000"),
		NotionAPIKey:                         getEnvWithDefault("NOTION_API_KEY", ""),
		NotionKaimemoDatabaseInputID:         getEnvWithDefault("NOTION_KAIMEMO_DB_INPUT_ID", ""),
		NotionKaimemoDatabaseSummaryRecordID: getEnvWithDefault("NOTION_KAIMEMO_DB_SUMMARY_ID", ""),
		AllowOrigins:                         []string{os.Getenv("ALLOW_ORIGINS"), "https://access.line.me/oauth2/v2.1/authorize"},
		LINEConfig:                           lineConfig,
		LINELoginFrontendCallbackURL:         os.Getenv("LINE_LOGIN_FRONTEND_CALLBACK_URL"),
		DatabaseConfig:                       dbConfig,
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
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}
