package handler

import (
	"echo-household-budget/internal/usecase"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

type (
	PublishInformationRequest struct {
		ID int `json:"id" param:"id"`
	}

	PublishInformationResponse struct {
		ID int `json:"id"`
	}

	publishInformationHandler struct {
		publishInformationUsecase usecase.PublishInformationUsecase
	}

	PublishInformationHandler interface {
		Handle(c echo.Context) error
	}
)

// Handle implements PublishInformationHandler.
func (p *publishInformationHandler) Handle(c echo.Context) error {
	request := &PublishInformationRequest{}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	spew.Dump(request)

	input := usecase.PublishInformationInput{
		ID: request.ID,
	}

	output, err := p.publishInformationUsecase.Execute(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, output)
}

func NewPublishInformationHandler(publishInformationUsecase usecase.PublishInformationUsecase) PublishInformationHandler {
	return &publishInformationHandler{
		publishInformationUsecase: publishInformationUsecase,
	}
}
