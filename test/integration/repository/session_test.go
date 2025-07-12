package repository_test

import (
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	"github.com/DucTran999/auth-service/pkg/hasher"
	"github.com/stretchr/testify/require"
)

func TestSessionRepo(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ctx := t.Context()
	accountRepo := repository.NewAccountRepo(db.DB())
	sessionRepo := repository.NewSessionRepository(db.DB())
	hasher := hasher.NewHasher()

	hashedPassword, err := hasher.HashPassword("sTr0ngP@ssg0rk")
	require.NoError(t, err)

	account := model.Account{
		Email:        "daniel@example.go",
		PasswordHash: hashedPassword,
	}

	// Pre setup an account
	err = accountRepo.Create(ctx, &account)
	require.NoError(t, err)
	require.NotEmpty(t, account.ID)

	t.Run("create session then find by ID then mark to expire", func(t *testing.T) {
		session := &model.Session{
			AccountID: account.ID,
			IPAddress: "192.168.1.1",
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
			CreatedAt: time.Now(),
		}

		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Verify session created
		found, err := sessionRepo.FindByID(ctx, session.ID.String())
		require.NoError(t, err)
		require.NotNil(t, found)
		require.Equal(t, session.ID, found.ID)
		require.Equal(t, session.AccountID, found.AccountID)

		// Mark this session as failed
		err = sessionRepo.MarkSessionsExpired(ctx, []string{session.ID.String()}, time.Now().Add(-1*time.Hour))
		require.NoError(t, err)
	})

	t.Run("find all active sessions", func(t *testing.T) {
		session1 := &model.Session{
			AccountID: account.ID,
			IPAddress: "192.168.1.1",
			UserAgent: "Chrome",
			CreatedAt: time.Now(),
		}
		err := sessionRepo.Create(ctx, session1)
		require.NoError(t, err)

		sessions, err := sessionRepo.FindAllActiveSession(ctx)
		require.NoError(t, err)
		require.Len(t, sessions, 1)
	})

	t.Run("find session by ID but not found", func(t *testing.T) {
		session, err := sessionRepo.FindByID(ctx, "8f5c6b1e-dc99-4e33-a8c0-3e58fba86a65")
		require.ErrorIs(t, err, errs.ErrSessionNotFound)
		require.Nil(t, session)
	})
}
