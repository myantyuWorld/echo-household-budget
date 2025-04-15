package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	domainmodel "echo-household-budget/internal/domain/model"
)

func TestCategoryRepository_Create(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	category := &domainmodel.Category{
		Name:  "食費",
		Color: "#FF0000",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "categories"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), category.Name, category.Color).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateMasterCategory(category)
	assert.NoError(t, err)
	assert.Equal(t, domainmodel.CategoryID(1), category.ID)
}

func TestCategoryRepository_FindByID(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1 ORDER BY "categories"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "color", "created_at", "updated_at"}).
			AddRow(1, "食費", "#FF0000", nil, nil))

	category, err := repo.FindMasterCategoryByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, domainmodel.CategoryID(1), category.ID)
	assert.Equal(t, "食費", category.Name)
	assert.Equal(t, "#FF0000", category.Color)
}

func TestCategoryRepository_Update(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	category := &domainmodel.Category{
		ID:    domainmodel.CategoryID(1),
		Name:  "食費（更新）",
		Color: "#FF0000",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "categories" SET "color"=\$1,"name"=\$2,"updated_at"=\$3 WHERE "id" = \$4`).
		WithArgs(category.Color, category.Name, sqlmock.AnyArg(), category.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateMasterCategory(category)
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

	err := repo.DeleteMasterCategory(1)
	assert.NoError(t, err)
}

func TestCategoryRepository_NotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// FindByID - Not Found
	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1 ORDER BY "categories"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "color"}))

	category, err := repo.FindMasterCategoryByID(999)
	assert.NoError(t, err)
	assert.Nil(t, category)

	// Update - Not Found
	category = &domainmodel.Category{
		ID:    domainmodel.CategoryID(999),
		Name:  "存在しないカテゴリ",
		Color: "#FF0000",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "categories" SET "color"=\$1,"name"=\$2,"updated_at"=\$3 WHERE "id" = \$4`).
		WithArgs(category.Color, category.Name, sqlmock.AnyArg(), category.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.UpdateMasterCategory(category)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Delete - Not Found
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "categories" WHERE "categories"."id" = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.DeleteMasterCategory(999)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCategoryRepository_CreateHouseHoldCategory(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	categoryLimit := &domainmodel.CategoryLimit{
		HouseholdBookID: domainmodel.HouseHoldID(1),
		CategoryID:      domainmodel.CategoryID(1),
		LimitAmount:     10000,
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "category_limits"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), categoryLimit.HouseholdBookID, categoryLimit.CategoryID, categoryLimit.LimitAmount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateHouseHoldCategory(categoryLimit)
	assert.NoError(t, err)
}

func TestCategoryRepository_FindHouseHoldCategoryByHouseHoldID(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectQuery(`SELECT \* FROM "category_limits" WHERE "category_limits"."id" = \$1 ORDER BY "category_limits"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "household_book_id", "category_id", "limit_amount", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 10000, nil, nil))

	categoryLimit, err := repo.FindHouseHoldCategoryByHouseHoldID(1)
	assert.NoError(t, err)
	assert.NotNil(t, categoryLimit)
	assert.Equal(t, domainmodel.CategoryLimitID(1), categoryLimit.ID)
	assert.Equal(t, domainmodel.HouseHoldID(1), categoryLimit.HouseholdBookID)
	assert.Equal(t, domainmodel.CategoryID(1), categoryLimit.CategoryID)
	assert.Equal(t, 10000, categoryLimit.LimitAmount)
}

func TestCategoryRepository_UpdateHouseHoldCategory(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	categoryLimit := &domainmodel.CategoryLimit{
		ID:              domainmodel.CategoryLimitID(1),
		HouseholdBookID: domainmodel.HouseHoldID(1),
		CategoryID:      domainmodel.CategoryID(1),
		LimitAmount:     15000,
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "category_limits" SET "category_id"=\$1,"household_book_id"=\$2,"limit_amount"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
		WithArgs(categoryLimit.CategoryID, categoryLimit.HouseholdBookID, categoryLimit.LimitAmount, sqlmock.AnyArg(), categoryLimit.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateHouseHoldCategory(categoryLimit)
	assert.NoError(t, err)
}

func TestCategoryRepository_DeleteHouseHoldCategory(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "category_limits" WHERE "category_limits"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteHouseHoldCategory(1)
	assert.NoError(t, err)
}

func TestCategoryRepository_HouseHoldCategoryNotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewCategoryRepository(gormDB)

	// FindHouseHoldCategoryByHouseHoldID - Not Found
	mock.ExpectQuery(`SELECT \* FROM "category_limits" WHERE "category_limits"."id" = \$1 ORDER BY "category_limits"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "household_book_id", "category_id", "limit_amount"}))

	categoryLimit, err := repo.FindHouseHoldCategoryByHouseHoldID(999)
	assert.NoError(t, err)
	assert.Nil(t, categoryLimit)

	// UpdateHouseHoldCategory - Not Found
	categoryLimit = &domainmodel.CategoryLimit{
		ID:              domainmodel.CategoryLimitID(999),
		HouseholdBookID: domainmodel.HouseHoldID(1),
		CategoryID:      domainmodel.CategoryID(1),
		LimitAmount:     10000,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "category_limits" SET "category_id"=\$1,"household_book_id"=\$2,"limit_amount"=\$3,"updated_at"=\$4 WHERE "id" = \$5`).
		WithArgs(categoryLimit.CategoryID, categoryLimit.HouseholdBookID, categoryLimit.LimitAmount, sqlmock.AnyArg(), categoryLimit.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.UpdateHouseHoldCategory(categoryLimit)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// DeleteHouseHoldCategory - Not Found
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "category_limits" WHERE "category_limits"."id" = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err = repo.DeleteHouseHoldCategory(999)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
