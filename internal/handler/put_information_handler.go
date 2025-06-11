package handler

import "github.com/labstack/echo/v4"

type (
	PutInformationRequest struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	PutInformationResponse struct {
		ID int `json:"id"`
	}

	putInformationHandler struct {
	}

	PutInformationHandler interface {
		Handle(c echo.Context) error
	}
)

func NewPutInformationHandler() PutInformationHandler {
	return &putInformationHandler{}
}

// Handle implements PutInformationHandler.
func (p *putInformationHandler) Handle(c echo.Context) error {
	panic("unimplemented")
}
