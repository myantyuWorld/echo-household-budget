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
