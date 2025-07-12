package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/repository"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSessionRepo_MarkSessionsExpired_DBError(t *testing.T) {
	t.Parallel()
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	sessionIDs := []string{
		"9f21c863-8b63-4b91-9eb3-649d0e6a8d1e",
		"b3d9c7fa-93cb-4d3e-85b6-63a4c742a420",
	}
	expiresAt := time.Now()

	mock.ExpectBegin()
	// expireAt, updated_at, IN(...)
	mock.ExpectExec(`UPDATE "sessions"`).
		WithArgs(expiresAt, sqlmock.AnyArg(), sessionIDs[0], sessionIDs[1]).
		WillReturnError(errors.New("simulated db failure"))

	mock.ExpectRollback()

	sut := repository.NewSessionRepository(db)
	err := sut.MarkSessionsExpired(t.Context(), sessionIDs, expiresAt)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated db failure")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_MarkSessionsExpired_EmptyIDs(t *testing.T) {
	t.Parallel()
	db, _, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	repo := repository.NewSessionRepository(db)
	err := repo.MarkSessionsExpired(t.Context(), []string{}, time.Now())
	require.NoError(t, err)
}

func TestSessionRepo_FindAllActiveSession_DBError(t *testing.T) {
	t.Parallel()
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT "id" FROM "sessions" WHERE expires_at IS NULL`).
		WillReturnError(errors.New("simulated db error"))

	repo := repository.NewSessionRepository(db)
	sessions, err := repo.FindAllActiveSession(t.Context())

	require.Error(t, err)
	require.Nil(t, sessions)
	require.Contains(t, err.Error(), "simulated db error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_DeleteExpiredBefore_DBError(t *testing.T) {
	t.Parallel()
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	cutoff := time.Now()

	mock.ExpectExec(`DELETE FROM sessions WHERE expires_at < ?`).
		WithArgs(cutoff).
		WillReturnError(errors.New("simulated delete error"))

	repo := repository.NewSessionRepository(db)
	err := repo.DeleteExpiredBefore(t.Context(), cutoff)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated delete error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_UpdateExpiresAt_DBError(t *testing.T) {
	t.Parallel()
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	sessionID := "11111111-1111-1111-1111-111111111111"
	expiresAt := time.Now()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "sessions"`).
		WithArgs(expiresAt, sqlmock.AnyArg(), sessionID).
		WillReturnError(errors.New("simulated update failure"))
	mock.ExpectRollback()

	repo := repository.NewSessionRepository(db)
	err := repo.UpdateExpiresAt(t.Context(), sessionID, expiresAt)

	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated update failure")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_FindByID_DBError(t *testing.T) {
	db, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	sessionID := "broken-id"

	mock.ExpectQuery(`FROM "sessions" WHERE id = \$1`).
		WithArgs(sessionID, sqlmock.AnyArg()).
		WillReturnError(errors.New("simulated db error"))

	repo := repository.NewSessionRepository(db)
	session, err := repo.FindByID(t.Context(), sessionID)

	require.Error(t, err)
	require.Nil(t, session)
	require.Contains(t, err.Error(), "simulated db error")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSessionRepo_Create_DBError(t *testing.T) {
	gormDB, mock, cleanup := mockbuilder.NewMockGormDB(t)
	defer cleanup()

	repo := repository.NewSessionRepository(gormDB)

	now := time.Now().Add(24 * time.Hour)
	session := &model.Session{
		AccountID: uuid.New(),
		ExpiresAt: &now,
	}

	// Setup mock to simulate DB error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "sessions"`).
		WithArgs(
			session.AccountID.String(),
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			session.ExpiresAt,
		).
		WillReturnError(errors.New("simulated db error"))
	mock.ExpectRollback()

	// Act
	err := repo.Create(t.Context(), session)

	// Assert
	require.Error(t, err)
	require.Contains(t, err.Error(), "simulated db error")

	// Check mock expectations
	require.NoError(t, mock.ExpectationsWereMet())
}
