package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
)

type (
	RegisterChatMessageInput struct {
		HouseholdID int
		UserID      int
		Message     string
	}

	RegisterChatMessageUsecase interface {
		Execute(request RegisterChatMessageInput) (*domainmodel.ChatMessage, error)
	}

	registerChatMessageUsecase struct {
		chatMessageRepository repository.ChatMessageRepository
	}
)

func NewRegisterChatMessageUsecase(chatMessageRepository repository.ChatMessageRepository) RegisterChatMessageUsecase {
	return &registerChatMessageUsecase{
		chatMessageRepository: chatMessageRepository,
	}
}

func (u *registerChatMessageUsecase) Execute(request RegisterChatMessageInput) (*domainmodel.ChatMessage, error) {
	chatMessage := domainmodel.NewChatMessage(request.HouseholdID, request.UserID, request.Message)

	if err := u.chatMessageRepository.Create(chatMessage); err != nil {
		return nil, err
	}

	// TODO : メッセージをAIサービスに送信して、返答を受け取る

	aiChatReplyMessage := domainmodel.NewAIChatReplyMessage(request.HouseholdID)
	if err := u.chatMessageRepository.Create(aiChatReplyMessage); err != nil {
		return nil, err
	}

	return aiChatReplyMessage, nil
}
