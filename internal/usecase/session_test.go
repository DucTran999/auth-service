package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/usecase"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/require"
)

func NewSessionUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) usecase.SessionUsecase {
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
		setup       func(t *testing.T) usecase.SessionUsecase
		expectedErr error
	}

	testTable := []testcase{
		{
			name: "failed to get list active session from db",
			setup: func(t *testing.T) usecase.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrFindActiveSession,
		},
		{
			name: "failed when try to find sessionID key expired",
			setup: func(t *testing.T) usecase.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMissingKeys,
		},
		{
			name: "failed when trying to mark session expires time",
			setup: func(t *testing.T) usecase.SessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.SessionRepoBuilder.FindAllActiveSessionSuccess()
				builders.CacheBuilder.CallMissingKeysSuccess()
				builders.SessionRepoBuilder.MarkSessionsExpiredFailed()
				return NewSessionUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrMarkSessionsExpired,
		},
		{
			name: "mark session expires success",
			setup: func(t *testing.T) usecase.SessionUsecase {
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
