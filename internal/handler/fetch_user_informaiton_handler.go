package handler

import (
	"echo-household-budget/internal/infrastructure/middleware"
	"echo-household-budget/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	FetchUserInformationRequest struct {
	}

	FetchUserInformationItemResponse struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		IsRead   bool   `json:"isRead"`
		Category string `json:"category"`
	}

	fetchUserInformationHandler struct {
		fetchUserInformationUsecase usecase.FetchUserInformationUsecase
	}

	FetchUserInformationHandler interface {
		Handle(c echo.Context) error
	}
)

func NewFetchUserInformationHandler(fetchUserInformationUsecase usecase.FetchUserInformationUsecase) FetchUserInformationHandler {
	return &fetchUserInformationHandler{
		fetchUserInformationUsecase: fetchUserInformationUsecase,
	}
}

func (h *fetchUserInformationHandler) Handle(c echo.Context) error {
	user, ok := middleware.GetUserFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	informations, err := h.fetchUserInformationUsecase.Execute(usecase.FetchUserInformationInput{
		UserID: int(user.ID),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := h.makeOutput(informations)
	return c.JSON(http.StatusOK, response)
}

func (h *fetchUserInformationHandler) makeOutput(informations []usecase.FetchUserInformationOutput) []FetchUserInformationItemResponse {
	output := make([]FetchUserInformationItemResponse, len(informations))
	for i, information := range informations {
		output[i] = FetchUserInformationItemResponse{
			ID:       information.ID,
			Title:    information.Title,
			Content:  information.Content,
			IsRead:   information.IsRead,
			Category: information.Category,
		}
	}
	return output
}
