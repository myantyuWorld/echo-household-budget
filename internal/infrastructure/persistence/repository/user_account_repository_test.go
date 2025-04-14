package repository

import (
	"template-echo-notion-integration/internal/domain/household"
	"template-echo-notion-integration/internal/infrastructure/persistence/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserAccountRepository_Create(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	userAccount := &household.UserAccount{
		UserID: "user123",
		Name:   "テストユーザー",
	}
	householdBook := &models.HouseholdBook{
		UserID: string(userAccount.UserID),
		Title:  "初期家計簿",
	}

	// SQLクエリのモック
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "user_accounts"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userAccount.UserID, userAccount.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery(`INSERT INTO "household_books"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), householdBook.UserID, householdBook.Title, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(userAccount)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), userAccount.ID)
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

	err := repo.Delete(1)
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

	err := repo.Delete(999)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserAccountRepository_FindByLINEUserID(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	userAccount := &household.UserAccount{
		UserID: household.LINEUserID("user123"),
		Name:   "テストユーザー",
	}

	// SQLクエリのモックを修正
	mock.ExpectQuery(`SELECT \* FROM "user_accounts" WHERE user_id = \$1 ORDER BY "user_accounts"."id" LIMIT \$2`).
		WithArgs(string(userAccount.UserID), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).
			AddRow(1, string(userAccount.UserID), userAccount.Name))

	foundAccount, err := repo.FindByLINEUserID(userAccount.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, foundAccount)
	assert.Equal(t, userAccount.UserID, foundAccount.UserID)
	assert.Equal(t, userAccount.Name, foundAccount.Name)
}

// 存在しないユーザーのテストケースを追加
func TestUserAccountRepository_FindByLINEUserID_NotFound(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewUserAccountRepository(gormDB)

	// 存在しないユーザーIDで検索
	mock.ExpectQuery(`SELECT \* FROM "user_accounts" WHERE user_id = \$1 ORDER BY "user_accounts"."id" LIMIT \$2`).
		WithArgs("non_existent_user", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	foundAccount, err := repo.FindByLINEUserID(household.LINEUserID("non_existent_user"))
	assert.Error(t, err)
	assert.Nil(t, foundAccount)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
