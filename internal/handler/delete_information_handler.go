package handler

import "github.com/labstack/echo/v4"

type (
	DeleteInformationRequest struct {
		ID int `json:"id"`
	}

	DeleteInformationResponse struct {
		ID int `json:"id"`
	}

	deleteInformationHandler struct {
	}

	DeleteInformationHandler interface {
		Handle(c echo.Context) error
	}
)

func NewDeleteInformationHandler() DeleteInformationHandler {
	return &deleteInformationHandler{}
}

// Handle implements DeleteInformationHandler.
func (d *deleteInformationHandler) Handle(c echo.Context) error {
	panic("unimplemented")
}
