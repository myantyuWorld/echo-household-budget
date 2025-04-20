package repository

import (
	"testing"

	domainmodel "echo-household-budget/internal/domain/model"
	"echo-household-budget/internal/infrastructure/persistence/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserAccountRepository_Create(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	userAccount := &domainmodel.UserAccount{
		UserID:     "user123",
		Name:       "テストユーザー",
		PictureURL: "https://example.com/picture.jpg",
	}
	householdBook := &models.HouseholdBook{
		Title: "初期家計簿",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "user_accounts"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userAccount.UserID, userAccount.Name, userAccount.PictureURL).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`INSERT INTO "household_books"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), householdBook.Title, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(userAccount)
	assert.NoError(t, err)
	assert.Equal(t, domainmodel.UserID(1), userAccount.ID)
}

func TestUserAccountRepository_Delete(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "user_accounts" WHERE "user_accounts"."id" = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(domainmodel.UserID(1))
	assert.NoError(t, err)
}

func TestUserAccountRepository_NotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	// Delete - Not Found
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "user_accounts" WHERE "user_accounts"."id" = \$1`).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := repo.Delete(domainmodel.UserID(999))
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserAccountRepository_FindByLINEUserID(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	userAccount := &domainmodel.UserAccount{
		ID:         1,
		UserID:     "user123",
		Name:       "テストユーザー",
		PictureURL: "https://example.com/picture.jpg",
	}

	// SQLクエリのモックを修正
	mock.ExpectQuery(`SELECT \* FROM "user_accounts" WHERE user_id = \$1 ORDER BY "user_accounts"."id" LIMIT \$2`).
		WithArgs(string(userAccount.UserID), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "picture_url"}).
			AddRow(1, string(userAccount.UserID), userAccount.Name, userAccount.PictureURL))

	// HouseholdBooksのプリロード用のモック（中間テーブルを経由）
	mock.ExpectQuery(`SELECT \* FROM "user_households" WHERE "user_households"."user_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "household_id"}).
			AddRow(1, 1))

	mock.ExpectQuery(`SELECT \* FROM "household_books" WHERE "household_books"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).
			AddRow(1, "テスト家計簿", "テスト用"))

	// CategoryLimitsのプリロード用のモック
	mock.ExpectQuery(`SELECT \* FROM "category_limits" WHERE "category_limits"."household_book_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "household_book_id", "category_id", "limit_amount"}).
			AddRow(1, 1, 1, 10000))

	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "食費"))

	foundAccount, err := repo.FindByLINEUserID(userAccount.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, foundAccount)
	spew.Dump(foundAccount)
	assert.Equal(t, userAccount.Name, foundAccount.Name)
	assert.Equal(t, userAccount.PictureURL, foundAccount.PictureURL)

	// 関連テーブルの値を確認
	assert.Len(t, foundAccount.HouseholdBooks, 1)
	assert.Equal(t, "テスト家計簿", foundAccount.HouseholdBooks[0].Title)
	assert.Equal(t, "テスト用", foundAccount.HouseholdBooks[0].Description)

	assert.Len(t, foundAccount.HouseholdBooks[0].CategoryLimit, 1)
	assert.Equal(t, domainmodel.HouseHoldID(1), foundAccount.HouseholdBooks[0].CategoryLimit[0].HouseholdBookID)
	assert.Equal(t, domainmodel.CategoryID(1), foundAccount.HouseholdBooks[0].CategoryLimit[0].Category.ID)
	assert.Equal(t, 10000, foundAccount.HouseholdBooks[0].CategoryLimit[0].LimitAmount)
}

// 存在しないユーザーのテストケースを追加
func TestUserAccountRepository_FindByLINEUserID_NotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	// 存在しないユーザーIDで検索
	mock.ExpectQuery(`SELECT \* FROM "user_accounts" WHERE user_id = \$1 ORDER BY "user_accounts"."id" LIMIT \$2`).
		WithArgs("non_existent_user", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	foundAccount, err := repo.FindByLINEUserID(domainmodel.LINEUserID("non_existent_user"))
	assert.Error(t, err)
	assert.Nil(t, foundAccount)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserAccountRepository_FetchMe(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	userAccount := &domainmodel.UserAccount{
		ID:         1,
		UserID:     "user123",
		Name:       "テストユーザー",
		PictureURL: "https://example.com/picture.jpg",
	}

	// SQLクエリのモック
	mock.ExpectQuery(`SELECT \* FROM "user_accounts" WHERE user_accounts.id = \$1 ORDER BY "user_accounts"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "picture_url"}).
			AddRow(1, userAccount.UserID, userAccount.Name, userAccount.PictureURL))

	// HouseholdBooksのプリロード用のモック（中間テーブルを経由）
	mock.ExpectQuery(`SELECT \* FROM "user_households" WHERE "user_households"."user_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "household_id"}).
			AddRow(1, 1))

	mock.ExpectQuery(`SELECT \* FROM "household_books" WHERE "household_books"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).
			AddRow(1, "テスト家計簿", "テスト用"))

	// CategoryLimitsのプリロード用のモック
	mock.ExpectQuery(`SELECT \* FROM "category_limits" WHERE "category_limits"."household_book_id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "household_book_id", "category_id", "limit_amount"}).
			AddRow(1, 1, 1, 10000))

	mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."id" = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "食費"))

	foundAccount, err := repo.FetchMe(domainmodel.UserID(1))
	assert.NoError(t, err)
	assert.NotNil(t, foundAccount)
	assert.Equal(t, userAccount.Name, foundAccount.Name)
	assert.Equal(t, userAccount.PictureURL, foundAccount.PictureURL)

	// 関連テーブルの値を確認
	assert.Len(t, foundAccount.HouseholdBooks, 1)
	assert.Equal(t, "テスト家計簿", foundAccount.HouseholdBooks[0].Title)
	assert.Equal(t, "テスト用", foundAccount.HouseholdBooks[0].Description)

	assert.Len(t, foundAccount.HouseholdBooks[0].CategoryLimit, 1)
	assert.Equal(t, domainmodel.HouseHoldID(1), foundAccount.HouseholdBooks[0].CategoryLimit[0].HouseholdBookID)
	assert.Equal(t, domainmodel.CategoryID(1), foundAccount.HouseholdBooks[0].CategoryLimit[0].Category.ID)
	assert.Equal(t, 10000, foundAccount.HouseholdBooks[0].CategoryLimit[0].LimitAmount)
}
