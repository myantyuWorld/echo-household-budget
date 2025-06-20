package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
)

type chatMessageRepository struct {
	db *gorm.DB
}

// FindByHouseholdID implements repository.ChatMessageRepository.
func (r *chatMessageRepository) FindByHouseholdID(householdID int, limit int, offset int) ([]*domainmodel.ChatMessage, error) {
	chatMessages := []*models.ChatMessage{}

	if err := r.db.
		Where("household_id = ?", householdID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("User").
		Find(&chatMessages).Error; err != nil {
		return nil, err
	}

	domainChatMessages := []*domainmodel.ChatMessage{}
	for _, chatMessage := range chatMessages {
		domainChatMessages = append(domainChatMessages, &domainmodel.ChatMessage{
			ID:          chatMessage.ID,
			HouseholdID: chatMessage.HouseholdID,
			UserID:      chatMessage.UserID,
			MessageType: domainmodel.ChatMessageType(chatMessage.MessageType),
			Content:     chatMessage.Content,
			User: &domainmodel.UserAccount{
				ID:   domainmodel.UserID(chatMessage.User.ID),
				Name: chatMessage.User.Name,
			},
		})
	}

	return domainChatMessages, nil
}

func NewChatMessageRepository(db *gorm.DB) repository.ChatMessageRepository {
	return &chatMessageRepository{db: db}
}

func (r *chatMessageRepository) Create(input *domainmodel.ChatMessage) error {
	spew.Dump(input)

	chatMessage := &models.ChatMessage{
		ID:          input.ID,
		HouseholdID: input.HouseholdID,
		UserID:      input.UserID,
		MessageType: string(input.MessageType),
		Content:     input.Content,
	}

	return r.db.Create(chatMessage).Error
}
