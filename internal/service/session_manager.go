package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

// セッション管理
type SessionManager interface {
	CreateSession(userID string) (string, error)
	GetSession(sessionID string) (string, error)
	DestroySession(sessionID string) error
}

type sessionManager struct {
	sessionStore map[string]string
}

func NewSessionManager() SessionManager {
	return &sessionManager{
		sessionStore: make(map[string]string),
	}
}

func (s *sessionManager) CreateSession(userID string) (string, error) {
	sessionID := fmt.Sprintf("session-%s", generateSessionID())
	s.sessionStore[sessionID] = userID
	return sessionID, nil
}

func (s *sessionManager) GetSession(sessionID string) (string, error) {
	userID, exists := s.sessionStore[sessionID]
	if !exists {
		return "", errors.New("Session invalid")
	}
	return userID, nil
}

func (s *sessionManager) DestroySession(sessionID string) error {
	delete(s.sessionStore, sessionID)
	return nil
}

// セッションIDをランダムに生成
func generateSessionID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // 生成に失敗した場合はエラーを返す
	}
	return hex.EncodeToString(bytes) // 16バイトのランダムな値を16進数文字列に変換
}
