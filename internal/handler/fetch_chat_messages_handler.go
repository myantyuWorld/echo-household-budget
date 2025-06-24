package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	FetchChatMessagesRequest struct {
		HouseholdID int `json:"household_id" param:"household_id" validate:"required"`
		Limit       int `json:"limit" query:"limit"`
		Offset      int `json:"offset" query:"offset"`
	}

	ChatMessage struct {
		ID          int    `json:"id"`
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		Content     string `json:"content"`
		MessageType string `json:"messageType"`
		CreatedAt   string `json:"createdAt"`
	}

	fetchChatMessagesHandler struct {
		// usecase usecase.FetchChatMessagesUsecase
	}

	FetchChatMessagesHandler interface {
		Handle(c echo.Context) error
	}
)

func NewFetchChatMessagesHandler() FetchChatMessagesHandler {
	return &fetchChatMessagesHandler{}
}

func (h *fetchChatMessagesHandler) Handle(c echo.Context) error {
	request := FetchChatMessagesRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	output := h.makeOutput()

	return c.JSON(http.StatusOK, output)
}

func (h *fetchChatMessagesHandler) makeOutput() []ChatMessage {
	return []ChatMessage{
		{
			ID:          1,
			UserID:      1,
			UserName:    "end user",
			Content:     "test message",
			MessageType: "user",
			CreatedAt:   "2021-01-01T00:00:00Z",
		},
		{
			ID:          2,
			UserID:      2,
			UserName:    "ai assistant",
			Content:     "test message",
			MessageType: "ai",
			CreatedAt:   "2021-01-01T00:00:00Z",
		},
		{
			ID:          3,
			UserID:      3,
			UserName:    "end user",
			Content:     "今月のビール代っていくらかかった？",
			MessageType: "user",
			CreatedAt:   "2021-01-01T00:00:00Z",
		},
		{
			ID:          4,
			UserID:      4,
			UserName:    "ai assistant",
			Content:     "今月のビール代は、合計で「10000円」でした。結構買いましたね、疲れてる？ホゲホゲホゲホゲホゲホゲホゲホゲホゲホゲホゲホゲホゲホゲ",
			MessageType: "ai",
			CreatedAt:   "2021-01-01T00:00:00Z",
		},
	}
}
