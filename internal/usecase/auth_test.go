package usecase_test

import (
	"context"
	"log"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAuthUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) usecase.AuthUseCase {

	return usecase.NewAuthUseCase(
		builders.HasherBuilder.GetInstance(),
		builders.CacheBuilder.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) usecase.AuthUseCase
		loginInput  usecase.LoginInput
		expectedErr error
		expected    *model.Account
	}

	loginInput := usecase.LoginInput{
		Email:    "daniel@example.com",
		Password: "abc1234!",
	}
	expectedAccount := &model.Account{
		ID:       mockbuilder.FakeAccountID,
		Email:    mockbuilder.FakeEmail,
		IsActive: true,
	}

	testTable := []testCase{
		{
			name: "failed to find email in db",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailError()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: mockbuilder.ErrFindAccountByEmail,
			expected:    nil,
		},
		{
			name: "failed email not existed",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailNoResult()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: usecase.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "account inactive",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailAccountInactive()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: usecase.ErrAccountDisabled,
			expected:    nil,
		},
		{
			name: "password not match",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordNotMatch()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: usecase.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "compare password unexpected error",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.CompareHashPasswordGotError()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: mockbuilder.ErrCompareHashPassword,
			expected:    nil,
		},
		{
			name: "find session in db got error",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.GetCacheErr()
				builders.SessionRepoBuilder.FindSessionError()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: mockbuilder.ErrFindSessionByID,
			expected:    nil,
		},
		{
			name: "session miss cache still valid in db",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindSessionSuccess()
				builders.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "extend expires time in cache failed",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindSessionSuccess()
				builders.CacheBuilder.SetCacheSessionFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "expired session create new one",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindExpiredSession()
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				builders.SessionRepoBuilder.CreateSessionSuccess()
				builders.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "expired session create new one failed",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindExpiredSession()
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				builders.SessionRepoBuilder.CreateSessionFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: mockbuilder.ErrCreateSession,
			expected:    nil,
		},
		{
			name: "reuse session in cache",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.ValidSessionCached()
				builders.CacheBuilder.GetTTLSuccess()
				builders.CacheBuilder.SetExpireSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: usecase.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
	}

	ctx := context.Background()
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			session, err := sut.Login(ctx, tc.loginInput)
			log.Println(session)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.ID, session.AccountID)
			} else {
				assert.Nil(t, session)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	type testcase struct {
		name      string
		sessionID string
		setup     func(t *testing.T) usecase.AuthUseCase
		expectErr error
	}

	testTable := []testcase{
		{
			name:      "invalid session id",
			sessionID: "98f0fc0ec6d13b7b9c6b04d62e3de8bd4acdc2e5e7e017fc6fa3a1c8a36c9f4a",
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: usecase.ErrInvalidSessionID,
		},
		{
			name:      "failed to update session expires at",
			sessionID: mockbuilder.FakeSessionID.String(),
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.DelKeySuccess()
				builders.SessionRepoBuilder.UpdateExpiresAtFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: mockbuilder.ErrUpdateSessionExpires,
		},
		{
			name:      "logout success",
			sessionID: mockbuilder.FakeSessionID.String(),
			setup: func(t *testing.T) usecase.AuthUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.DelKeySuccess()
				builders.SessionRepoBuilder.UpdateExpiresAtSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			ctx := context.Background()

			err := sut.Logout(ctx, tc.sessionID)

			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}
