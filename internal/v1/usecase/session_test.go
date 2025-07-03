package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/v1/usecase"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewSessionUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) *usecase.SessionUCImpl {
	return usecase.NewSessionUC(
		builders.CacheBuilder.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestDeleteExpiredBefore(t *testing.T) {
	t.Run("delete got db error", func(t *testing.T) {
		builders := mockbuilder.NewBuilderContainer(t)
		builders.SessionRepoBuilder.DeleteExpiredBeforeFailed()
		sut := NewSessionUseCaseUT(t, builders)
		cutoff := time.Now().AddDate(0, 0, -30)
		ctx := context.Background()

		err := sut.DeleteExpiredBefore(ctx, cutoff)

		require.ErrorIs(t, err, mockbuilder.ErrDeleteExpiredBefore)
	})

	t.Run("delete got db success", func(t *testing.T) {
		builders := mockbuilder.NewBuilderContainer(t)
		builders.SessionRepoBuilder.DeleteExpiredBeforeSuccess()
		sut := NewSessionUseCaseUT(t, builders)
		cutoff := time.Now().AddDate(0, 0, -30)
		ctx := context.Background()

		err := sut.DeleteExpiredBefore(ctx, cutoff)

		require.NoError(t, err)
	})
}

func TestMarkExpiredSessions(t *testing.T) {
	type testcase struct {
		name        string
		setup       func(t *testing.T) *usecase.SessionUCImpl
		expectedErr error
	}

	testTable := []testcase{
		{
			name: "failed to get list active session from db",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrFindActiveSession,
		},
		{
			name: "failed when try to find sessionID key expired",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMissingKeys,
		},
		{
			name: "failed when trying to mark session expires time",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysSuccess()
				builders.SessionRepoBuilder.MarkSessionsExpiredFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMarkSessionsExpired,
		},
		{
			name: "no session are active",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindNoActiveSession()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: nil,
		},
		{
			name: "no session expire yet",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.NoMissingKeysFound()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: nil,
		},
		{
			name: "mark session expires success",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysSuccess()
				builders.SessionRepoBuilder.MarkSessionsExpiredSuccess()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			ctx := context.Background()

			err := sut.MarkExpiredSessions(ctx)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestValidateSession(t *testing.T) {
	type testcase struct {
		name              string
		setup             func(t *testing.T) *usecase.SessionUCImpl
		sessionID         string
		expectedAccountID uuid.UUID
		expectedErr       error
	}

	testTable := []testcase{
		{
			name: "invalid session id",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
				builders := mockbuilder.NewBuilderContainer(t)
				return NewSessionUseCaseUT(t, builders)
			},
			sessionID:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			expectedErr: model.ErrInvalidSessionID,
		},
		{
			name: "session found in cache",
			setup: func(t *testing.T) *usecase.SessionUCImpl {
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
			setup: func(t *testing.T) *usecase.SessionUCImpl {
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
			setup: func(t *testing.T) *usecase.SessionUCImpl {
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
			setup: func(t *testing.T) *usecase.SessionUCImpl {
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
			setup: func(t *testing.T) *usecase.SessionUCImpl {
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
			sut := tc.setup(t)
			ctx := context.Background()

			session, err := sut.ValidateSession(ctx, tc.sessionID)

			if err != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.Equal(t, tc.expectedAccountID, session.AccountID)
			}
		})
	}
}
