package repository_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/require"
)

func TestAccountRepo_FindByEmail_DBError(t *testing.T) {
	gormDB, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	// Set up expectation — simulate query error
	mock.ExpectQuery(`SELECT .* FROM "accounts" WHERE email = \$1 ORDER BY "accounts"\."id" LIMIT \$2`).
		WillReturnError(errors.New("simulated db failure"))

	// Instantiate the repository and call the method
	sut := repository.NewAccountRepo(gormDB)

	account, err := sut.FindByEmail(t.Context(), "fail@example.com")

	// Assert expected outcomes
	require.Nil(t, account)
	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated db failure")

	// Ensure mock expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepo_FindByID_DBError(t *testing.T) {
	gormDB, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	// Set up expectation — simulate query error
	mock.ExpectQuery(`SELECT .* FROM "accounts" WHERE id = \$1 ORDER BY "accounts"\."id" LIMIT \$2`).
		WillReturnError(errors.New("simulated db failure"))

	// Instantiate the repository and call the method
	sut := repository.NewAccountRepo(gormDB)
	ctx := t.Context()

	account, err := sut.FindByID(ctx, "8f5c6b1e-dc99-4e33-a8c0-3e58fba86a65")

	// Assert expected outcomes
	require.Nil(t, account)
	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated db failure")

	// Ensure mock expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepo_Create_DBError(t *testing.T) {
	gormDB, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	account := model.Account{
		Email:        "daniel@example.go",
		PasswordHash: "hashedP4ssw0rd",
		IsVerified:   false,
		IsActive:     true,
		Role:         "user",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "accounts"`).
		WithArgs(
			account.Email,
			account.PasswordHash,
			account.IsVerified,
			account.IsActive,
			account.Role,
		).
		WillReturnError(errors.New("simulated insert failure"))
	mock.ExpectRollback()

	sut := repository.NewAccountRepo(gormDB)

	err := sut.Create(t.Context(), &account)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated insert failure")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepo_UpdatePasswordHash_DBError(t *testing.T) {
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	id := "8f5c6b1e-dc99-4e33-a8c0-3e58fba86a65"
	hash := "newhashedpassword"

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "accounts"`).
		WithArgs(hash, sqlmock.AnyArg(), id).
		WillReturnError(errors.New("simulated update failure"))
	mock.ExpectRollback()

	repo := repository.NewAccountRepo(db)
	err := repo.UpdatePasswordHash(t.Context(), id, hash)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated update failure")
	require.NoError(t, mock.ExpectationsWereMet())
}
