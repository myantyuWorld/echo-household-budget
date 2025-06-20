package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
)

type (
	FetchChatMessageInput struct {
		HouseholdID int
		Limit       int
		Offset      int
	}

	FetchChatMessageOutput struct {
		ChatMessages []*domainmodel.ChatMessage
	}

	FetchChatMessageUsecase interface {
		Execute(request FetchChatMessageInput) (*FetchChatMessageOutput, error)
	}

	fetchChatMessageUsecase struct {
		chatMessageRepository repository.ChatMessageRepository
	}
)

func NewFetchChatMessageUsecase(chatMessageRepository repository.ChatMessageRepository) FetchChatMessageUsecase {
	return &fetchChatMessageUsecase{
		chatMessageRepository: chatMessageRepository,
	}
}

func (u *fetchChatMessageUsecase) Execute(request FetchChatMessageInput) (*FetchChatMessageOutput, error) {
	chatMessages, err := u.chatMessageRepository.FindByHouseholdID(request.HouseholdID, request.Limit, request.Offset)
	if err != nil {
		return nil, err
	}
	return &FetchChatMessageOutput{
		ChatMessages: chatMessages,
	}, nil
}
