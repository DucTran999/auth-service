package auth_test

import (
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAuthUseCaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) port.AuthSessionUsecase {
	return auth.NewAuthSessionUsecase(
		builders.CacheBuilder.GetInstance(),
		builders.AccountVerifier.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
		builders.SessionRepoBuilder.GetInstance(),
	)
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) port.AuthSessionUsecase
		loginInput  dto.LoginInput
		expectedErr error
		expected    *model.Account
	}

	loginInput := dto.LoginInput{
		Email:    mockbuilder.FakeEmail,
		Password: mockbuilder.FakeOldPass,
	}
	expectedAccount := &model.Account{
		ID:       mockbuilder.FakeAccountID,
		Email:    mockbuilder.FakeEmail,
		IsActive: true,
	}

	testTable := []testCase{
		{
			name: "no session id provided",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountVerifier.VerifySuccess()
				builders.SessionRepoBuilder.CreateSessionSuccess()
				builders.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: uuid.Nil.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "verify account failed",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountVerifier.VerifyFailed(model.ErrInvalidCredentials)
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: uuid.Nil.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: model.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "find session in db got error",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.GetCacheErr()
				builders.SessionRepoBuilder.FindByIdFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: mockbuilder.ErrFindSessionByID,
			expected:    nil,
		},
		{
			name: "session miss cache still valid in db",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSuccess()
				builders.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "extend expires time in cache failed",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSuccess()
				builders.CacheBuilder.SetCacheSessionFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "expired session create new one",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSessionExpired()
				builders.AccountVerifier.VerifySuccess()
				builders.SessionRepoBuilder.CreateSessionSuccess()
				builders.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
		{
			name: "expired session create new one failed",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.SessionMissCache()
				builders.SessionRepoBuilder.FindByIDSessionExpired()
				builders.AccountVerifier.VerifySuccess()
				builders.SessionRepoBuilder.CreateSessionFailed()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: mockbuilder.ErrCreateSession,
			expected:    nil,
		},
		{
			name: "reuse session in cache",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.CacheBuilder.ValidSessionCached()
				builders.CacheBuilder.GetTTLSuccess()
				builders.CacheBuilder.SetExpireSuccess()
				return NewAuthUseCaseUT(t, builders)
			},
			loginInput: dto.LoginInput{
				CurrentSessionID: mockbuilder.FakeSessionID.String(),
				Email:            loginInput.Email,
				Password:         loginInput.Password,
			},
			expectedErr: nil,
			expected:    expectedAccount,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			ctx := t.Context()

			// Act
			session, err := sut.Login(ctx, tc.loginInput)

			// Assert
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
		setup     func(t *testing.T) port.AuthSessionUsecase
		expectErr error
	}

	testTable := []testcase{
		{
			name:      "invalid session id",
			sessionID: "98f0fc0ec6d13b7b9c6b04d62e3de8bd4acdc2e5e7e017fc6fa3a1c8a36c9f4a",
			setup: func(t *testing.T) port.AuthSessionUsecase {
				builders := mockbuilder.NewBuilderContainer(t)
				return NewAuthUseCaseUT(t, builders)
			},
			expectErr: model.ErrInvalidSessionID,
		},
		{
			name:      "failed to update session expires at",
			sessionID: mockbuilder.FakeSessionID.String(),
			setup: func(t *testing.T) port.AuthSessionUsecase {
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
			setup: func(t *testing.T) port.AuthSessionUsecase {
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
			ctx := t.Context()

			err := sut.Logout(ctx, tc.sessionID)

			require.ErrorIs(t, err, tc.expectErr)
		})
	}
}
