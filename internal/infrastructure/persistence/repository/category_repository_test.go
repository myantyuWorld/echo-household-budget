package repository

import (
	"template-echo-notion-integration/internal/domain/household"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCategoryRepository_Create(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	category := &household.Category{
		Name:  "食費",
		Color: "#FF0000",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "categories"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), category.Name, category.Color).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(category)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), category.ID)
}

func TestCategoryRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1 ORDER BY "categories"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "color", "created_at", "updated_at"}).
			AddRow(1, "食費", "#FF0000", nil, nil))

	category, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, uint(1), category.ID)
	assert.Equal(t, "食費", category.Name)
	assert.Equal(t, "#FF0000", category.Color)
}

func TestCategoryRepository_Update(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	category := &household.Category{
		ID:    1,
		Name:  "食費（更新）",
		Color: "#FF0000",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "categories" SET "color"=\$1,"name"=\$2,"updated_at"=\$3 WHERE "id" = \$4`).
		WithArgs(category.Color, category.Name, sqlmock.AnyArg(), category.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(category)
	assert.NoError(t, err)
}

func TestCategoryRepository_Delete(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "categories" WHERE "categories"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(1)
	assert.NoError(t, err)
}

func TestCategoryRepository_NotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// FindByID - Not Found
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1 ORDER BY "categories"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "color"}))

	category, err := repo.FindByID(999)
	assert.NoError(t, err)
	assert.Nil(t, category)

	// Update - Not Found
	category = &household.Category{
		ID:    999,
		Name:  "存在しないカテゴリ",
		Color: "#FF0000",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "categories" SET "color"=\$1,"name"=\$2,"updated_at"=\$3 WHERE "id" = \$4`).
		WithArgs(category.Color, category.Name, sqlmock.AnyArg(), category.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.Update(category)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Delete - Not Found
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "categories" WHERE "categories"."id" = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.Delete(999)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
