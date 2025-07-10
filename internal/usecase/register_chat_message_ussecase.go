package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	repository "echo-household-budget/internal/domain/repository"
	"echo-household-budget/internal/infrastructure/llm"
	"log"

	"github.com/davecgh/go-spew/spew"
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
		llmClient             *llm.LLMClient
	}
)

func NewRegisterChatMessageUsecase(chatMessageRepository repository.ChatMessageRepository, llmClient *llm.LLMClient) RegisterChatMessageUsecase {
	return &registerChatMessageUsecase{
		chatMessageRepository: chatMessageRepository,
		llmClient:             llmClient,
	}
}

func (u *registerChatMessageUsecase) Execute(request RegisterChatMessageInput) (*domainmodel.ChatMessage, error) {
	chatMessage := domainmodel.NewChatMessage(request.HouseholdID, request.UserID, request.Message)

	if err := u.chatMessageRepository.Create(chatMessage); err != nil {
		return nil, err
	}

	// LLMClientでFunction Calling処理
	response, err := u.llmClient.ProcessMessage(request.Message)
	log.Println("============ LLM response =============")
	spew.Dump(response)
	if err != nil {
		return nil, err
	}

	aiChatReplyMessage := domainmodel.NewAIChatReplyMessage(request.HouseholdID, response)
	if err := u.chatMessageRepository.Create(aiChatReplyMessage); err != nil {
		return nil, err
	}

	return aiChatReplyMessage, nil
}
