//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package handler

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/usecase"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

type receiptAnalyzeHandler struct {
	usecase usecase.ReceiptAnalyzeUsecase
}

type CreateReceiptRequest struct {
	HouseholdID uint   `json:"householdID" param:"houseHoldID"`
	ImageData   string `json:"imageData"`
}

type CreateReceiptAnalyzeResultRequest struct {
	Total      uint                 `json:"total"`
	S3FilePath string               `json:"s3FilePath"`
	Items      []ReceiptAnalyzeItem `json:"items"`
}

type ReceiptAnalyzeItem struct {
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

// CreateReceiptAnalyzeReception implements ReceiptAnalyzeHandler.
func (r *receiptAnalyzeHandler) CreateReceiptAnalyzeReception(c echo.Context) error {
	req := CreateReceiptRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if len(req.ImageData) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "image_data is required",
		})
	}

	receipt := &domainmodel.ReceiptAnalyzeReception{
		HouseholdBookID: domainmodel.HouseHoldID(req.HouseholdID),
		ImageData:       req.ImageData,
	}

	if err := r.usecase.CreateReceiptAnalyzeReception(receipt); err != nil {
		spew.Dump(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

// CreateReceiptAnalyzeResult implements ReceiptAnalyzeHandler.
func (r *receiptAnalyzeHandler) CreateReceiptAnalyzeResult(c echo.Context) error {
	req := CreateReceiptAnalyzeResultRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	spew.Dump(req)

	// TODO：ここ、わざわざハンドラーでやらない方がいい｜具体的には、ドメインモデルで、変換処理をしたらいい？
	items := make([]domainmodel.ReceiptAnalyzeItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domainmodel.ReceiptAnalyzeItem{
			Name:  item.Name,
			Price: item.Price,
		}
	}

	result := &domainmodel.ReceiptAnalyze{
		TotalPrice: req.Total,
		S3FilePath: req.S3FilePath,
		Items:      items,
	}

	if err := r.usecase.CreateReceiptAnalyzeResult(result); err != nil {
		log.Println(err)
		spew.Dump(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)

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
