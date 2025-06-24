package models

import "time"

type ChatMessage struct {
	ID          int
	HouseholdID int
	UserID      int
	MessageType string
	Content     string
	User        UserAccount
	CreatedAt   time.Time
}

func (m *ChatMessage) TableName() string {
	return "chat_messages"
}
