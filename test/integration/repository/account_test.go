package repository_test

import (
	"testing"

	"github.com/DucTran999/auth-service/config"
	"github.com/DucTran999/auth-service/internal/container"
	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/DucTran999/dbkit"
	"github.com/stretchr/testify/require"
)

func SetupTestDB(t *testing.T) dbkit.Connection {
	cfg, err := config.LoadConfig(".test.env")
	require.NoError(t, err)

	c, err := container.NewContainer(cfg)
	require.NoError(t, err)

	db := c.AuthDB.DB()

	err = db.Exec("TRUNCATE sessions, accounts CASCADE").Error
	require.NoError(t, err)

	return c.AuthDB
}

func TestAccountRepo(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ctx := t.Context()
	repo := repository.NewAccountRepo(db.DB())
	hasher := hasher.NewHasher()

	hashedPassword, err := hasher.HashPassword("sTr0ngP@ssg0rk")
	require.NoError(t, err)

	account := model.Account{
		Email:        "daniel@example.go",
		PasswordHash: hashedPassword,
	}

	err = repo.Create(ctx, &account)
	require.NoError(t, err)
	require.NotEmpty(t, account.ID)

	t.Run("find by ID found", func(t *testing.T) {
		found, err := repo.FindByID(ctx, account.ID.String())
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, account.Email, found.Email)
	})

	t.Run("find by email found", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, account.Email)
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, account.ID, found.ID)
	})

	t.Run("find by email not found", func(t *testing.T) {
		found, err := repo.FindByEmail(ctx, "notfound@example.com")
		require.ErrorIs(t, err, errs.ErrAccountNotFound)
		require.Nil(t, found)
	})

	t.Run("find by ID not found", func(t *testing.T) {
		found, err := repo.FindByID(ctx, "8f5c6b1e-dc99-4e33-a8c0-3e58fba86a65")
		require.ErrorIs(t, err, errs.ErrAccountNotFound)
		require.Nil(t, found)
	})

	t.Run("update password ", func(t *testing.T) {
		newPassword := "n3wP@ssW0rd"
		hashedNewPassword, err := hasher.HashPassword(newPassword)
		require.NoError(t, err)

		err = repo.UpdatePasswordHash(ctx, account.ID.String(), hashedNewPassword)
		require.NoError(t, err)

		updated, err := repo.FindByID(ctx, account.ID.String())
		require.NoError(t, err)

		ok, err := hasher.ComparePasswordAndHash(newPassword, updated.PasswordHash)
		require.NoError(t, err)
		require.True(t, ok)
	})
}
