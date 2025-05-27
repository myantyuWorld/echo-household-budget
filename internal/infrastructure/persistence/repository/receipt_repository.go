package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
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
	panic("unimplemented")
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
			Price: item.Price,
		})
	}

	return &domainmodel.ReceiptAnalyze{
		ID:              models.ID,
		ImageURL:        models.ImageURL,
		AnalyzeStatus:   models.AnalyzeStatus,
		TotalPrice:      models.TotalPrice,
		HouseholdBookID: domainmodel.HouseHoldID(models.HouseholdBookID),
		Items:           items,
	}, nil
}

func NewReceiptRepository(db *gorm.DB) domainmodel.ReceiptAnalyzeRepository {
	return &ReceiptRepository{db: db}
}
