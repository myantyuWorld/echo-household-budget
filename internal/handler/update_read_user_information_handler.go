package handler

import (
	"echo-household-budget/internal/domain/repository"
	"echo-household-budget/internal/infrastructure/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	UpdateReadUserInformationRequest struct {
		InformationIDs []int `json:"informationIDs"`
	}

	updateReadUserInformationHandler struct {
		userInformationRepository repository.UserInformationRepository
	}
)

func NewUpdateReadUserInformationHandler(userInformationRepository repository.UserInformationRepository) *updateReadUserInformationHandler {
	return &updateReadUserInformationHandler{
		userInformationRepository: userInformationRepository,
	}
}

func (h *updateReadUserInformationHandler) Handle(c echo.Context) error {
	request := UpdateReadUserInformationRequest{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	user, ok := middleware.GetUserFromContext(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	err := h.userInformationRepository.UpdateRead(request.InformationIDs, int(user.ID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
}
