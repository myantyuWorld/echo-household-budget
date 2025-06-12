package handler

import (
	"echo-household-budget/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	FetchInformationItemResponse struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Content     string `json:"content"`
		IsPublished bool   `json:"isPublished"`
		Category    string `json:"category"`
	}

	fetchInformationsHandler struct {
		fetchInformationUsecase usecase.FetchInformationUsecase
	}

	FetchInformationsHandler interface {
		Handle(c echo.Context) error
	}
)

// Handle implements FetchInformationsHandler.
func (f *fetchInformationsHandler) Handle(c echo.Context) error {
	output, err := f.fetchInformationUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := f.makeOutput(output)
	return c.JSON(http.StatusOK, response)
}

// TODO: 後ほど、ユースケースのoutputを受け取るようにする
func (f *fetchInformationsHandler) makeOutput(output []usecase.FetchInformationOutput) []FetchInformationItemResponse {
	response := make([]FetchInformationItemResponse, len(output))
	for i, item := range output {
		response[i] = FetchInformationItemResponse{
			ID:          item.ID,
			Title:       item.Title,
			Content:     item.Content,
			IsPublished: item.IsPublished,
			Category:    item.Category,
		}
	}

	return response
}

func NewFetchInformationsHandler(fetchInformationUsecase usecase.FetchInformationUsecase) FetchInformationsHandler {
	return &fetchInformationsHandler{
		fetchInformationUsecase: fetchInformationUsecase,
	}
}
