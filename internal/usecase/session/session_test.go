package session_test

import (
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	"github.com/DucTran999/auth-service/internal/usecase/session"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewSessionUseCaseBackgroundUT(t *testing.T, builders *mockbuilder.BuilderContainer) port.SessionMaintenanceUsecase {
	return session.NewSessionUC(
		builders.CacheBuilder.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestDeleteExpiredBefore(t *testing.T) {
	t.Run("delete got db error", func(t *testing.T) {
		t.Parallel()
		builders := mockbuilder.NewBuilderContainer(t)
		builders.SessionRepoBuilder.DeleteExpiredBeforeFailed()
		sut := NewSessionUseCaseBackgroundUT(t, builders)
		cutoff := time.Now().AddDate(0, 0, -30)
		ctx := t.Context()

		err := sut.DeleteExpiredBefore(ctx, cutoff)

		require.ErrorIs(t, err, mockbuilder.ErrDeleteExpiredBefore)
	})

	t.Run("delete got db success", func(t *testing.T) {
		t.Parallel()
		builders := mockbuilder.NewBuilderContainer(t)
		builders.SessionRepoBuilder.DeleteExpiredBeforeSuccess()
		sut := NewSessionUseCaseBackgroundUT(t, builders)
		cutoff := time.Now().AddDate(0, 0, -30)
		ctx := t.Context()

		err := sut.DeleteExpiredBefore(ctx, cutoff)

		require.NoError(t, err)
	})
}

func TestMarkExpiredSessions(t *testing.T) {
	type testcase struct {
		name        string
		setup       func(t *testing.T) port.SessionMaintenanceUsecase
		expectedErr error
	}

	testTable := []testcase{
		{
			name: "failed to get list active session from db",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionFailed()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: mockbuilder.ErrFindActiveSession,
		},
		{
			name: "failed when try to find sessionID key expired",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysFailed()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMissingKeys,
		},
		{
			name: "failed when trying to mark session expires time",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysSuccess()
				builders.SessionRepoBuilder.MarkSessionsExpiredFailed()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMarkSessionsExpired,
		},
		{
			name: "no session are active",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindNoActiveSession()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: nil,
		},
		{
			name: "no session expire yet",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.NoMissingKeysFound()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: nil,
		},
		{
			name: "mark session expires success",
			setup: func(t *testing.T) port.SessionMaintenanceUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysSuccess()
				builders.SessionRepoBuilder.MarkSessionsExpiredSuccess()
				return NewSessionUseCaseBackgroundUT(t, builders)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)
			ctx := t.Context()

			err := sut.MarkExpiredSessions(ctx)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func NewSessionUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) port.SessionUsecase {
	return session.NewSessionUC(
		builders.CacheBuilder.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestValidateSession(t *testing.T) {
	type testcase struct {
		name              string
		setup             func(t *testing.T) port.SessionUsecase
		sessionID         string
		expectedAccountID uuid.UUID
		expectedErr       error
	}

	testTable := []testcase{
		{
			name: "invalid session id",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			expectedErr: model.ErrInvalidSessionID,
		},
		{
			name: "session found in cache",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.ValidSessionCached()
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:         mockbuilder.FakeSessionID.String(),
			expectedErr:       nil,
			expectedAccountID: mockbuilder.FakeAccountID,
		},
		{
			name: "miss cached but query db failed",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIdFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:   mockbuilder.FakeSessionID.String(),
			expectedErr: mockbuilder.ErrFindSessionByID,
		},
		{
			name: "session not found",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDNotFound()
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:   mockbuilder.FakeSessionID.String(),
			expectedErr: model.ErrSessionNotFound,
		},
		{
			name: "session already expired",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSessionExpired()
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:   mockbuilder.FakeSessionID.String(),
			expectedErr: model.ErrSessionNotFound,
		},
		{
			name: "failed when set cache",
			setup: func(t *testing.T) port.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSuccess()
				builders.CacheBuilder.SetCacheSessionFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:         mockbuilder.FakeSessionID.String(),
			expectedErr:       nil,
			expectedAccountID: mockbuilder.FakeAccountID,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)
			session, err := sut.Validate(t.Context(), tc.sessionID)

			if err != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.Equal(t, tc.expectedAccountID, session.AccountID)
			}
		})
	}
}
