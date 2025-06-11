package handler

import "github.com/labstack/echo/v4"

type (
	FetchInformationDetailRequest struct {
		ID int `json:"id"`
	}

	FetchInformationDetailResponse struct {
		Information FetchInformationItemResponse `json:"information"`
	}

	fetchInformationDetailHandler struct {
	}

	FetchInformationDetailHandler interface {
		Handle(c echo.Context) error
	}
)

func NewFetchInformationDetailHandler() FetchInformationDetailHandler {
	return &fetchInformationDetailHandler{}
}

// Handle implements FetchInformationDetailHandler.
func (f *fetchInformationDetailHandler) Handle(c echo.Context) error {
	panic("unimplemented")
}
