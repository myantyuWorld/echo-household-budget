package repository

import (
	domainmodel "echo-household-budget/internal/domain/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHouseHoldRepository_Create(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewHouseHoldRepository(gormDB)

	houseHold := &domainmodel.HouseHold{
		UserID: 1,
		Title:  "test",
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"household_books\"").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), houseHold.Title, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(houseHold)
	assert.NoError(t, err)
	assert.Equal(t, domainmodel.HouseHoldID(1), houseHold.ID)
}

func TestHouseHoldRepository_CreateUserHouseHold(t *testing.T) {
	gormDB, mock := setupTest(t)
	repo := NewHouseHoldRepository(gormDB)

	userHouseHold := &domainmodel.UserHouseHold{
		UserID:      1,
		HouseHoldID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"user_households\"").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userHouseHold.UserID, userHouseHold.HouseHoldID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateUserHouseHold(userHouseHold)
	assert.NoError(t, err)
	assert.Equal(t, domainmodel.HouseHoldID(1), userHouseHold.HouseHoldID)
}
