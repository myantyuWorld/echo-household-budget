package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

// UserAccountRepository はユーザーアカウントリポジトリの実装
type UserAccountRepository struct {
	db *gorm.DB
}

// FindByLINEUserID implements household.UserAccountRepository.
func (r *UserAccountRepository) FindByLINEUserID(userID domainmodel.LINEUserID) (*domainmodel.UserAccount, error) {
	var userAccount models.UserAccount
	if err := r.db.Where("user_id = ?", userID).First(&userAccount).Error; err != nil {
		return nil, err
	}

	return &domainmodel.UserAccount{
		ID:         userAccount.ID,
		Name:       userAccount.Name,
		PictureURL: userAccount.PictureURL,
	}, nil
}

// NewUserAccountRepository は新しいUserAccountRepositoryを作成します
func NewUserAccountRepository(db *gorm.DB) domainmodel.UserAccountRepository {
	return &UserAccountRepository{
		db: db,
	}
}

// Create は新しいユーザーアカウントを作成します
func (r *UserAccountRepository) Create(userAccount *domainmodel.UserAccount) error {
	model := &models.UserAccount{
		UserID:     string(userAccount.UserID),
		Name:       userAccount.Name,
		PictureURL: userAccount.PictureURL,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// ユーザーアカウントの作成
		if err := tx.Create(model).Error; err != nil {
			return err
		}

		// 家計簿レコードの作成
		householdBook := &models.HouseholdBook{
			UserID: model.UserID,
			Title:  "初期家計簿",
		}

		if err := tx.Create(householdBook).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 作成したモデルのIDをドメインモデルに設定
	userAccount.ID = model.ID
	return nil
}

// Delete は指定されたIDのユーザーアカウントを削除します
func (r *UserAccountRepository) Delete(id uint) error {
	result := r.db.Delete(&models.UserAccount{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
