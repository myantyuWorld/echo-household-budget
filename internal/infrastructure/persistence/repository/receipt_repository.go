package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

// FindReceiptAnalyzeByS3FilePath implements domainmodel.ReceiptAnalyzeRepository.
func (r *ReceiptRepository) FindReceiptAnalyzeByS3FilePath(s3FilePath string) (*domainmodel.ReceiptAnalyze, error) {
	var models models.ReceiptAnalyzes
	if err := r.db.Where("image_url = ?", s3FilePath).
		First(&models).Error; err != nil {
		return nil, err
	}

	var items []domainmodel.ReceiptAnalyzeItem
	for _, item := range models.Items {
		items = append(items, domainmodel.ReceiptAnalyzeItem{
			Name:  item.Name,
			Price: uint(item.Price),
		})
	}

	return &domainmodel.ReceiptAnalyze{
		ID:              uint(models.ID),
		TotalPrice:      uint(models.TotalPrice),
		S3FilePath:      models.ImageURL,
		HouseholdBookID: domainmodel.HouseHoldID(models.HouseholdBookID),
		Items:           items,
	}, nil
}

// CreateReceiptAnalyzeReception implements domainmodel.ReceiptAnalyzeRepository.
func (r *ReceiptRepository) CreateReceiptAnalyzeReception(receiptAnalyze *domainmodel.ReceiptAnalyzeReception) error {
	model := models.ReceiptAnalyzes{
		ImageURL:        receiptAnalyze.ImageURL,
		HouseholdBookID: int(receiptAnalyze.HouseholdBookID),
		AnalyzeStatus:   "pending",
	}

	return r.db.Create(&model).Error
}

// CreateReceiptAnalyzeResult implements domainmodel.ReceiptAnalyzeRepository.
func (r *ReceiptRepository) CreateReceiptAnalyzeResult(receiptAnalyze *domainmodel.ReceiptAnalyze) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		items := make([]models.ReceiptAnalyzeItems, len(receiptAnalyze.Items))
		for i, item := range receiptAnalyze.Items {
			items[i] = models.ReceiptAnalyzeItems{
				ReceiptAnalyzeID: int(receiptAnalyze.ID),
				Name:             item.Name,
				Price:            int(item.Price),
			}
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		model := models.ReceiptAnalyzes{
			TotalPrice:    int(receiptAnalyze.TotalPrice),
			AnalyzeStatus: "finished",
			Items:         items,
		}

		if err := tx.Model(&models.ReceiptAnalyzes{}).Where("id = ?", receiptAnalyze.ID).Updates(&model).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

// FindByID implements domainmodel.ReceiptAnalyzeRepository.
func (r *ReceiptRepository) FindByID(id domainmodel.HouseHoldID) (*domainmodel.ReceiptAnalyze, error) {
	var models models.ReceiptAnalyzes
	if err := r.db.Where("id = ?", id).
		Preload("Items").
		First(&models).Error; err != nil {
		return nil, err
	}

	var items []domainmodel.ReceiptAnalyzeItem
	for _, item := range models.Items {
		items = append(items, domainmodel.ReceiptAnalyzeItem{
			Name:  item.Name,
			Price: uint(item.Price),
		})
	}

	return &domainmodel.ReceiptAnalyze{
		ID:         uint(models.ID),
		TotalPrice: uint(models.TotalPrice),
		S3FilePath: models.ImageURL,
		Items:      items,
	}, nil
}

func NewReceiptRepository(db *gorm.DB) domainmodel.ReceiptAnalyzeRepository {
	return &ReceiptRepository{db: db}
}
