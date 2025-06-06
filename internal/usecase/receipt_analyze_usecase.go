//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/domain/repository"
	domainservice "echo-household-budget/internal/domain/service"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type receiptAnalyzeUsecase struct {
	repo             domainmodel.ReceiptAnalyzeRepository
	fileStorage      repository.FileStorageRepository
	houseHoldService domainservice.HouseHoldService
}

// CreateReceiptAnalyzeReception implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error {
	b64data := receipt.ImageData[strings.IndexByte(receipt.ImageData, ',')+1:]

	imageData, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return err
	}
	// uuid-household_id-category_id-yyyyMMddHHmmss.jpg
	fileName := fmt.Sprintf("%s-%d-%d-%s.jpg", uuid.New().String(), receipt.HouseholdBookID, receipt.CategoryID, time.Now().Format("20060102150405"))
	_, err = r.fileStorage.UploadFile(imageData, fileName)
	if err != nil {
		return err
	}

	receipt.ImageURL = fileName
	return r.repo.CreateReceiptAnalyzeReception(receipt)
}

// CreateReceiptAnalyzeResult implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error {
	receiptAnalyze, err := r.repo.FindReceiptAnalyzeByS3FilePath(receipt.S3FilePath)
	if err != nil {
		return err
	}

	receiptAnalyze.TotalPrice = receipt.TotalPrice
	receiptAnalyze.Items = receipt.Items

	if err := r.repo.CreateReceiptAnalyzeResult(receiptAnalyze); err != nil {
		return err
	}

	shoppingAmount := domainmodel.NewShoppingAmount(receiptAnalyze.HouseholdBookID, receipt.CategoryID, int(receiptAnalyze.TotalPrice), time.Now().Format("2006-01-02"), "aiによるレシート分析", int(receiptAnalyze.ID))
	if err := r.houseHoldService.CreateShoppingAmount(shoppingAmount); err != nil {
		return err
	}

	return nil
}

// FindByID implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error) {
	return r.repo.FindByID(id)
}

type ReceiptAnalyzeUsecase interface {
	CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error
	CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error
	FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error)
}

func NewReceiptAnalyzeUsecase(repo domainmodel.ReceiptAnalyzeRepository, fileStorage repository.FileStorageRepository, houseHoldService domainservice.HouseHoldService) ReceiptAnalyzeUsecase {
	return &receiptAnalyzeUsecase{repo: repo, fileStorage: fileStorage, houseHoldService: houseHoldService}
}
