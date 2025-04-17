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

// fetchUserAccount はユーザーアカウントを取得する共通処理
func (r *UserAccountRepository) fetchUserAccount(condition string, args ...interface{}) (*domainmodel.UserAccount, error) {
	var userAccount models.UserAccount
	if err := r.db.
		Debug().
		Preload("HouseholdBooks.CategoryLimits").
		Preload("HouseholdBooks.CategoryLimits.Category").
		Where(condition, args...).
		First(&userAccount).Error; err != nil {
		return nil, err
	}

	// 関連テーブルの値をドメインモデルに変換
	householdBooks := make([]*domainmodel.HouseHold, len(userAccount.HouseholdBooks))
	var categoryLimits []*domainmodel.CategoryLimit

	for i, hb := range userAccount.HouseholdBooks {
		householdBooks[i] = &domainmodel.HouseHold{
			ID:          domainmodel.HouseHoldID(hb.ID),
			Title:       hb.Title,
			Description: hb.Description,
		}

		// HouseholdBookのCategoryLimitsを追加
		for _, cl := range hb.CategoryLimits {
			categoryLimits = append(categoryLimits, &domainmodel.CategoryLimit{
				ID:              domainmodel.CategoryLimitID(cl.ID),
				HouseholdBookID: domainmodel.HouseHoldID(cl.HouseholdBookID),
				CategoryID:      domainmodel.CategoryID(cl.CategoryID),
				LimitAmount:     cl.LimitAmount,
			})
		}

		householdBooks[i].CategoryLimit = categoryLimits
	}

	return &domainmodel.UserAccount{
		ID:             domainmodel.UserID(userAccount.ID),
		UserID:         domainmodel.LINEUserID(userAccount.UserID),
		Name:           userAccount.Name,
		PictureURL:     userAccount.PictureURL,
		HouseholdBooks: householdBooks,
	}, nil
}

// FetchMe implements domainmodel.UserAccountRepository.
func (r *UserAccountRepository) FetchMe(userID domainmodel.UserID) (*domainmodel.UserAccount, error) {
	return r.fetchUserAccount("user_accounts.id = ?", userID)
}

// FindByLINEUserID implements household.UserAccountRepository.
func (r *UserAccountRepository) FindByLINEUserID(userID domainmodel.LINEUserID) (*domainmodel.UserAccount, error) {
	return r.fetchUserAccount("user_id = ?", userID)
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
			Title: "初期家計簿",
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
	userAccount.ID = domainmodel.UserID(model.ID)
	return nil
}

// Delete は指定されたIDのユーザーアカウントを削除します
func (r *UserAccountRepository) Delete(id domainmodel.UserID) error {
	result := r.db.Delete(&models.UserAccount{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
