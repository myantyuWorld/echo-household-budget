package handler

import (
	"net/http"

	"echo-household-budget/internal/usecase"

	"github.com/labstack/echo/v4"
)

type (
	CreateInformationRequest struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}

	CreateInformationResponse struct {
		ID int `json:"id"`
	}

	createInformationHandler struct {
		usecase usecase.CreateInformationUsecase
	}

	CreateInformationHandler interface {
		Handle(c echo.Context) error
	}
)

// Handle implements CreateInformationHandler.
func (h *createInformationHandler) Handle(c echo.Context) error {
	request := CreateInformationRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	input := usecase.CreateInformationInput{
		Title:    request.Title,
		Content:  request.Content,
		Category: request.Category,
	}

	output, err := h.usecase.Execute(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	response := h.makeOutput(output)

	return c.JSON(http.StatusOK, response)

}

func (h *createInformationHandler) makeOutput(output usecase.CreateInformationOutput) CreateInformationResponse {
	return CreateInformationResponse{
		ID: output.ID,
	}
}

func NewCreateInformationHandler(usecase usecase.CreateInformationUsecase) CreateInformationHandler {
	return &createInformationHandler{
		usecase: usecase,
	}
}
