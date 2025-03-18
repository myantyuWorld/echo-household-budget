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
	setEnv("NOTION_DATABASE_KAIMEMO_INPUT", "test-database-input-id")
	setEnv("NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD", "test-database-summary-id")
	setEnv("FRONTEND_URL", "https://example.com")
	setEnv("LINE_CLIENT_ID", "test-client-id")
	setEnv("LINE_CLIENT_SECRET", "test-client-secret")
	setEnv("LINE_JWT_SECRET", "test-jwt-secret")
	setEnv("LINE_STATE", "test-state")
	setEnv("LINE_REDIRECT_URI", "https://example.com/callback")
	setEnv("LINE_TOKEN_URL", "https://example.com/token")

	// テスト終了後に環境変数をリセット
	defer unsetEnv("NOTION_API_KEY", "NOTION_DATABASE_KAIMEMO_INPUT", "NOTION_DATABASE_KAIMEMO_SUMMARY_RECORD", "FRONTEND_URL")

	// `log.Fatal` をキャッチするためにリカバリ
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected fatal error: %v", r)
		}
	}()

	// Configロード
	config := LoadConfig()

	// 検証
	assert.NotNil(t, config)
	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, "test-api-key", config.NotionAPIKey)
	assert.Equal(t, "test-database-input-id", config.NotionKaimemoDatabaseInputID)
	assert.Equal(t, "test-database-summary-id", config.NotionKaimemoDatabaseSummaryRecordID)
	assert.Contains(t, config.AllowOrigins, "https://example.com")

	assert.Equal(t, "test-client-id", config.LINEConfig.ClientID)
	assert.Equal(t, "test-client-secret", config.LINEConfig.ClientSecret)
	assert.Equal(t, "test-jwt-secret", config.LINEConfig.JwtSecret)
	assert.Equal(t, "test-state", config.LINEConfig.State)
	assert.Equal(t, "https://example.com/callback", config.LINEConfig.RedirectURI)
	assert.Equal(t, "https://example.com/token", config.LINEConfig.TokenURL)
}
