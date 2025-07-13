package auth_test

import (
	"testing"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/usecase/auth"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/internal/usecase/port"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAuthJWTUsecaseUT(t *testing.T, builders *mockbuilder.BuilderContainer) port.AuthJWTUsecase {
	return auth.NewAuthJWTUsecase(
		builders.TokenSigner.GetInstance(),
		builders.CacheBuilder.GetInstance(),
		builders.AccountVerifier.GetInstance(),
	)
}

func TestLoginWithJWT(t *testing.T) {
	type testcase struct {
		name        string
		setup       func(t *testing.T) port.AuthJWTUsecase
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "authenticate failed",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountVerifier.VerifyFailed(errs.ErrInvalidCredentials)
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name: "sign access token failed",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountVerifier.VerifySuccess()
				builder.TokenSigner.SignFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSigningToken,
		},
		{
			name: "sign refresh token failed",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountVerifier.VerifySuccess()
				builder.TokenSigner.SignAccessSuccessAndSignRefreshFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSigningToken,
		},
		{
			name: "set cache session failed",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountVerifier.VerifySuccess()
				builder.TokenSigner.SignSuccess()
				builder.CacheBuilder.SetCacheSessionFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSetCacheSession,
		},
		{
			name: "login success",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.AccountVerifier.VerifySuccess()
				builder.TokenSigner.SignSuccess()
				builder.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := dto.LoginJWTInput{
				Email:    "daniel@gmail.com",
				Password: "test1234",
			}
			sut := tc.setup(t)

			// Act
			tokens, err := sut.Login(t.Context(), input)

			// Assert
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	type testcase struct {
		name        string
		token       string
		setup       func(t *testing.T) port.AuthJWTUsecase
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "empty",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name:  "invalid token",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name:  "delete token cache failed",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.DelKeyErr()
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
		{
			name:  "revoke success",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.DelKeySuccess()
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)

			// Act
			err := sut.RevokeRefreshToken(t.Context(), tc.token)

			// Assert
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type testcase struct {
		name        string
		token       string
		setup       func(t *testing.T) port.AuthJWTUsecase
		expectedErr error
	}

	testcases := []testcase{
		{
			name: "empty",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name:  "invalid token",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrParsingToken,
		},
		{
			name:  "get refresh token in cache error",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.CheckRefreshTokenFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrHasCache,
		},
		{
			name:  "refresh token invalid cache",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenInvalidCache()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: errs.ErrInvalidCredentials,
		},
		{
			name:  "revoke old token failed",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenValidCache()
				builder.CacheBuilder.DelKeyErr()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrDelCacheDelete,
		},
		{
			name:  "sign tokens failed",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenValidCache()
				builder.CacheBuilder.DelKeySuccess()
				builder.TokenSigner.SignFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSigningToken,
		},
		{
			name:  "sign refresh token failed",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenValidCache()
				builder.CacheBuilder.DelKeySuccess()
				builder.TokenSigner.SignAccessSuccessAndSignRefreshFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSigningToken,
		},
		{
			name:  "cache new refresh token failed",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenValidCache()
				builder.CacheBuilder.DelKeySuccess()
				builder.TokenSigner.SignSuccess()
				builder.CacheBuilder.SetRefreshTokenFailed()
				return NewAuthJWTUsecaseUT(t, builder)
			},
			expectedErr: mockbuilder.ErrSetCache,
		},
		{
			name:  "refresh success",
			token: "dummy token",
			setup: func(t *testing.T) port.AuthJWTUsecase {
				t.Helper()
				builder := mockbuilder.NewBuilderContainer(t)
				builder.TokenSigner.ParseIntoSuccess()
				builder.CacheBuilder.RefreshTokenValidCache()
				builder.CacheBuilder.DelKeySuccess()
				builder.TokenSigner.SignSuccess()
				builder.CacheBuilder.SetCacheSessionSuccess()
				return NewAuthJWTUsecaseUT(t, builder)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sut := tc.setup(t)

			// Act
			tokens, err := sut.RefreshToken(t.Context(), tc.token)

			// Assert
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				require.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
		})
	}
}
