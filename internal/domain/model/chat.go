package domainmodel

import "time"

const (
	ChatMessageTypeUser ChatMessageType = "user"
	ChatMessageTypeAI   ChatMessageType = "ai"
)

type ChatMessage struct {
	ID          int
	HouseholdID int
	UserID      int
	MessageType ChatMessageType
	Content     string
	CreatedAt   time.Time
	User        *UserAccount
}

type ChatMessageType string

func NewChatMessage(householdID int, userID int, content string) *ChatMessage {
	return &ChatMessage{
		HouseholdID: householdID,
		UserID:      userID,
		MessageType: ChatMessageTypeUser,
		Content:     content,
		CreatedAt:   time.Now(),
	}
}

func NewAIChatReplyMessage(householdID int, content string) *ChatMessage {
	return &ChatMessage{
		HouseholdID: householdID,
		UserID:      0,
		MessageType: ChatMessageTypeAI,
		Content:     content,
		CreatedAt:   time.Now(),
	}
}
