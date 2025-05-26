//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
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

	if req.HouseholdID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "household_id is required",
		})
	}

	if req.ImageData == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "image_data is required",
		})
	}

	receipt := &domainmodel.ReceiptAnalyzeReception{
		HouseholdBookID: req.HouseholdID,
		ImageData:       req.ImageData,
	}

	if err := r.usecase.CreateReceiptAnalyzeReception(receipt); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
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
