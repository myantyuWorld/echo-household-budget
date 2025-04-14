package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 環境変数をセットするヘルパー関数
func setEnv(key, value string) {
	os.Setenv(key, value)
}

// テスト後に環境変数をクリアする関数
func unsetEnv(keys ...string) {
	for _, key := range keys {
		os.Unsetenv(key)
	}
}

func TestLoadConfig_Success(t *testing.T) {
	// テスト用の環境変数をセット
	setEnv("NOTION_API_KEY", "test-api-key")
	setEnv("NOTION_KAIMEMO_DB_INPUT_ID", "test-database-input-id")
	setEnv("NOTION_KAIMEMO_DB_SUMMARY_ID", "test-database-summary-id")
	setEnv("ALLOW_ORIGINS", "https://example.com")
	setEnv("LINE_CHANNEL_ID", "test-client-id")
	setEnv("LINE_CHANNEL_SECRET", "test-client-secret")
	setEnv("LINE_REDIRECT_URI", "https://example.com/callback")

	// テスト終了後に環境変数をリセット
	defer unsetEnv(
		"NOTION_API_KEY",
		"NOTION_KAIMEMO_DB_INPUT_ID",
		"NOTION_KAIMEMO_DB_SUMMARY_ID",
		"ALLOW_ORIGINS",
		"LINE_CHANNEL_ID",
		"LINE_CHANNEL_SECRET",
		"LINE_REDIRECT_URI",
	)

	// Configロード
	config := LoadConfig()

	// 検証
	assert.NotNil(t, config)
	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, "test-api-key", config.NotionAPIKey)
	assert.Equal(t, "test-database-input-id", config.NotionKaimemoDatabaseInputID)
	assert.Equal(t, "test-database-summary-id", config.NotionKaimemoDatabaseSummaryRecordID)
	assert.Contains(t, config.AllowOrigins, "https://example.com")

	// LINE設定の検証
	assert.NotNil(t, config.LINEConfig)
	assert.Equal(t, "test-client-id", config.LINEConfig.ClientID)
	assert.Equal(t, "test-client-secret", config.LINEConfig.ClientSecret)
	assert.Equal(t, "https://example.com/callback", config.LINEConfig.RedirectURL)
}

func TestLoadConfig(t *testing.T) {
	// テスト用の環境変数を設定
	os.Setenv("PORT", "3000")
	os.Setenv("ALLOW_ORIGINS", "http://localhost:5173")
	os.Setenv("NOTION_API_KEY", "test_notion_key")
	os.Setenv("NOTION_KAIMEMO_DB_INPUT_ID", "test_input_id")
	os.Setenv("NOTION_KAIMEMO_DB_SUMMARY_ID", "test_summary_id")
	os.Setenv("LINE_CHANNEL_ID", "test_channel_id")
	os.Setenv("LINE_CHANNEL_SECRET", "test_channel_secret")
	os.Setenv("LINE_REDIRECT_URI", "http://localhost:5173/line/callback")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "kakeibo")

	// テスト終了後に環境変数をクリア
	defer unsetEnv(
		"PORT",
		"ALLOW_ORIGINS",
		"NOTION_API_KEY",
		"NOTION_KAIMEMO_DB_INPUT_ID",
		"NOTION_KAIMEMO_DB_SUMMARY_ID",
		"LINE_CHANNEL_ID",
		"LINE_CHANNEL_SECRET",
		"LINE_REDIRECT_URI",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	)

	config := LoadConfig()

	// 設定値の検証
	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, []string{"http://localhost:5173", "https://access.line.me/oauth2/v2.1/authorize"}, config.AllowOrigins)
	assert.Equal(t, "test_notion_key", config.NotionAPIKey)
	assert.Equal(t, "test_input_id", config.NotionKaimemoDatabaseInputID)
	assert.Equal(t, "test_summary_id", config.NotionKaimemoDatabaseSummaryRecordID)
	assert.Equal(t, "test_channel_id", config.LINEConfig.ClientID)
	assert.Equal(t, "test_channel_secret", config.LINEConfig.ClientSecret)
	assert.Equal(t, "http://localhost:5173/line/callback", config.LINEConfig.RedirectURL)
	assert.Equal(t, "localhost", config.DatabaseConfig.Host)
	assert.Equal(t, "5432", config.DatabaseConfig.Port)
	assert.Equal(t, "postgres", config.DatabaseConfig.User)
	assert.Equal(t, "postgres", config.DatabaseConfig.Password)
	assert.Equal(t, "kakeibo", config.DatabaseConfig.DBName)
}

func TestNewDBConnection(t *testing.T) {
	tests := []struct {
		name        string
		config      *DatabaseConfig
		expectError bool
	}{
		{
			name: "正常系：データベース接続が成功",
			config: &DatabaseConfig{
				Host:     "db", // Dockerコンテナ名
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				DBName:   "kakeibo",
			},
			expectError: false,
		},
		{
			name: "異常系：無効なホスト",
			config: &DatabaseConfig{
				Host:     "invalid-host",
				Port:     "5432",
				User:     "postgres",
				Password: "postgres",
				DBName:   "kakeibo",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDBConnection(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				if err != nil {
					t.Logf("データベース接続エラー: %v", err)
					t.Skip("データベースが利用できません")
				}
				assert.NoError(t, err)
				assert.NotNil(t, db)

				// 接続を閉じる
				sqlDB, err := db.DB()
				assert.NoError(t, err)
				sqlDB.Close()
			}
		})
	}
}
