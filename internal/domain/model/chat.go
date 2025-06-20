package domainmodel

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
	User        *UserAccount
}

type ChatMessageType string

func NewChatMessage(householdID int, userID int, content string) *ChatMessage {
	return &ChatMessage{
		HouseholdID: householdID,
		UserID:      userID,
		MessageType: ChatMessageTypeUser,
		Content:     content,
	}
}

func NewAIChatReplyMessage(householdID int) *ChatMessage {
	return &ChatMessage{
		HouseholdID: householdID,
		UserID:      0,
		MessageType: ChatMessageTypeAI,
		Content:     "受け付けました、解析中です🤖",
	}
}
