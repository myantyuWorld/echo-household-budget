package models

type ChatMessage struct {
	ID          int
	HouseholdID int
	UserID      int
	MessageType string
	Content     string
	User        UserAccount
}

func (m *ChatMessage) TableName() string {
	return "chat_messages"
}
