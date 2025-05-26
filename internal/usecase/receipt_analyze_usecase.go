//go:generate mockgen -source=$GOFILE -destination=../mock/$GOPACKAGE/mock_$GOFILE -package=mock
package usecase

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/domain/repository"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type receiptAnalyzeUsecase struct {
	repo        domainmodel.ReceiptAnalyzeRepository
	fileStorage repository.FileStorageRepository
}

// CreateReceiptAnalyzeReception implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error {
	var imageURL string
	imageData, err := base64.StdEncoding.DecodeString(receipt.ImageData)
	if err != nil {
		return err
	}
	// uuid-household_id-yyyyMMddHHmmss.jpg
	fileName := fmt.Sprintf("%s-%d-%s.jpg", uuid.New().String(), receipt.HouseholdBookID, time.Now().Format("20060102150405"))
	imageURL, err = r.fileStorage.UploadFile(imageData, fileName)
	if err != nil {
		return err
	}

	receipt.ImageURL = imageURL
	return r.repo.CreateReceiptAnalyzeReception(receipt)
}

// CreateReceiptAnalyzeResult implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error {
	return r.repo.CreateReceiptAnalyzeResult(receipt)
}

// FindByID implements ReceiptAnalyzeUsecase.
func (r *receiptAnalyzeUsecase) FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error) {
	panic("unimplemented")
}

type ReceiptAnalyzeUsecase interface {
	CreateReceiptAnalyzeReception(receipt *domainmodel.ReceiptAnalyzeReception) error
	CreateReceiptAnalyzeResult(receipt *domainmodel.ReceiptAnalyze) error
	FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error)
}

func NewReceiptAnalyzeUsecase(repo domainmodel.ReceiptAnalyzeRepository) ReceiptAnalyzeUsecase {
	return &receiptAnalyzeUsecase{repo: repo}
}
