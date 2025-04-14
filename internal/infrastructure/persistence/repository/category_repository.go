package repository

import (
	"errors"
	"template-echo-notion-integration/internal/domain/household"
	"template-echo-notion-integration/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

// CategoryRepository はカテゴリリポジトリの実装
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository は新しいCategoryRepositoryを作成します
func NewCategoryRepository(db *gorm.DB) household.CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// Create は新しいカテゴリを作成します
func (r *CategoryRepository) Create(category *household.Category) error {
	model := &models.Category{
		Name:  category.Name,
		Color: category.Color,
	}

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	category.ID = model.ID
	return nil
}

// FindByID は指定されたIDのカテゴリを取得します
func (r *CategoryRepository) FindByID(id uint) (*household.Category, error) {
	var model models.Category
	if err := r.db.First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &household.Category{
		ID:    model.ID,
		Name:  model.Name,
		Color: model.Color,
	}, nil
}

// Update は既存のカテゴリを更新します
func (r *CategoryRepository) Update(category *household.Category) error {
	model := &models.Category{
		Base: models.Base{
			ID: category.ID,
		},
		Name:  category.Name,
		Color: category.Color,
	}

	result := r.db.Model(model).Updates(map[string]interface{}{
		"name":  category.Name,
		"color": category.Color,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete は指定されたIDのカテゴリを削除します
func (r *CategoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
