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

func NewAIChatReplyMessage(householdID int) *ChatMessage {
	return &ChatMessage{
		HouseholdID: householdID,
		UserID:      0,
		MessageType: ChatMessageTypeAI,
		Content:     "AIサービスで受け付けました、解析中です🤖(現在、実装中です)",
		CreatedAt:   time.Now(),
	}
}
