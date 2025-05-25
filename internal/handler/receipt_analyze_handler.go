package handler

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type receiptAnalyzeHandler struct {
	usecase usecase.ReceiptAnalyzeUsecase
}

type CreateReceiptRequest struct {
	HouseholdID domainmodel.HouseHoldID `json:"household_id" param:"household_id"`
	ImageData   string                  `json:"image_data"`
}

// CreateReceiptAnalyzeReception implements ReceiptAnalyzeHandler.
func (r *receiptAnalyzeHandler) CreateReceiptAnalyzeReception(c echo.Context) error {
	req := CreateReceiptRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	receipt := &domainmodel.ReceiptAnalyzeReception{
		HouseholdBookID: req.HouseholdID,
		ImageData:       req.ImageData,
	}

	return r.usecase.CreateReceiptAnalyzeReception(receipt)
}

// CreateReceiptAnalyzeResult implements ReceiptAnalyzeHandler.
func (r *receiptAnalyzeHandler) CreateReceiptAnalyzeResult(c echo.Context) error {
	panic("unimplemented")
}

// FindByID implements ReceiptAnalyzeHandler.
func (r *receiptAnalyzeHandler) FindByID(c echo.Context) error {
	panic("unimplemented")
}

type ReceiptAnalyzeHandler interface {
	CreateReceiptAnalyzeResult(c echo.Context) error
	CreateReceiptAnalyzeReception(c echo.Context) error
	FindByID(c echo.Context) error
}

func NewReceiptAnalyzeHandler(usecase usecase.ReceiptAnalyzeUsecase) ReceiptAnalyzeHandler {
	return &receiptAnalyzeHandler{usecase: usecase}
}
