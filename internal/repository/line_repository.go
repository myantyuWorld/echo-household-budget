//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/oauth2"
)

type LineRepository interface {
	GetUserInfo(code string) (*UserInfo, error)
	MatchState(state string) bool
	GetAuthCodeUrl() string
}

type lineRepository struct {
	lineConfig *oauth2.Config
	state      string
}

func NewLineRepository(lineConfig *oauth2.Config) LineRepository {
	return &lineRepository{
		lineConfig: lineConfig,
		state:      generateState(),
	}
}

// GetUserInfo implements LineRepository.
func (l *lineRepository) GetUserInfo(code string) (*UserInfo, error) {
	// Call OAuth2.0 Token Endpoint
	token, err := l.lineConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, errors.New("Token Exchange Failed")
	}

	// Call User Profile Endpoint
	client := l.lineConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.line.me/v2/profile")
	if err != nil {
		return nil, errors.New("Failed to get user info")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read user info")
	}

	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, errors.New("Failed to parse user info")
	}

	return &userInfo, nil
}

// MatchState implements LineRepository.
func (l *lineRepository) MatchState(state string) bool {
	if state != l.state {
		return false
	}
	return true
}

// GetAuthCode implements LineRepository.
func (l *lineRepository) GetAuthCodeUrl() string {
	url := l.lineConfig.AuthCodeURL(l.state)
	return url
}

// UserInfo From LINE
type UserInfo struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
}

func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	// base64 Encoding
	return base64.URLEncoding.EncodeToString(b)
}
