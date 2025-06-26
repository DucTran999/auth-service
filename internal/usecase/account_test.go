package usecase_test

import (
	"context"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewAccountUseCaseUT(
	t *testing.T,
	builders *mockbuilder.BuilderContainer,
) usecase.AccountUseCase {
	return usecase.NewAccountUseCase(
		builders.HasherBuilder.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
	)
}

func TestRegisterAccount(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) usecase.AccountUseCase
		accountInfo usecase.RegisterInput
		expectedErr error
		expected    *model.Account
	}

	userSample := usecase.RegisterInput{
		Email:    mockbuilder.FakeEmail,
		Password: "abc1234!",
	}

	testTable := []testCase{
		{
			name: "failed to find email in db",
			setup: func(t *testing.T) usecase.AccountUseCase {
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailError()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: mockbuilder.ErrFindAccountByEmail,
			expected:    nil,
		},
		{
			name: "failed caused email already taken",
			setup: func(t *testing.T) usecase.AccountUseCase {
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailHasResult()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: usecase.ErrEmailExisted,
			expected:    nil,
		},
		{
			name: "failed when hash password",
			setup: func(t *testing.T) usecase.AccountUseCase {
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordFailed()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: mockbuilder.ErrHashingPassword,
			expected:    nil,
		},
		{
			name: "failed when persist to db",
			setup: func(t *testing.T) usecase.AccountUseCase {
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordSuccess()
				b.AccountRepoBuilder.CreateAccountError()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: mockbuilder.ErrCreateAccount,
			expected:    nil,
		},
		{
			name: "register success",
			setup: func(t *testing.T) usecase.AccountUseCase {
				b := mockbuilder.NewBuilderContainer(t)
				b.AccountRepoBuilder.FindByEmailNoResult()
				b.HasherBuilder.HashingPasswordSuccess()
				b.AccountRepoBuilder.CreateAccountSuccess()
				return NewAccountUseCaseUT(t, b)
			},
			accountInfo: userSample,
			expectedErr: nil,
			expected: &model.Account{
				Email: "daniel@example.com",
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			user, err := sut.Register(context.Background(), tc.accountInfo)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Email, user.Email)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) usecase.AccountUseCase
		expectedErr error
	}

	input := usecase.ChangePasswordInput{
		AccountID:   mockbuilder.FakeAccountID.String(),
		OldPassword: mockbuilder.FakeOldPass,
		NewPassword: mockbuilder.FakeNewPass,
	}

	testTable := []testCase{
		{
			name: "failed find account by id",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrFindAccountByID,
		},
		{
			name: "failed to compare old password to hash",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdSuccess()
				builders.HasherBuilder.CompareHashPasswordGotError()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrCompareHashPassword,
		},
		{
			name: "old password not match",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdSuccess()
				builders.HasherBuilder.HashPasswordNotMatch()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: usecase.ErrInvalidCredentials,
		},
		{
			name: "failed to hashing password",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdSuccess()
				builders.HasherBuilder.HashPasswordMatch()
				builders.HasherBuilder.HashingPasswordFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrHashingPassword,
		},
		{
			name: "failed to update new password",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdSuccess()
				builders.HasherBuilder.HashPasswordMatch()
				builders.HasherBuilder.HashingPasswordSuccess()
				builders.AccountRepoBuilder.UpdatePasswordHashFailed()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: mockbuilder.ErrUpdateHashPassword,
		},
		{
			name: "change password success",
			setup: func(t *testing.T) usecase.AccountUseCase {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByIdSuccess()
				builders.HasherBuilder.HashPasswordMatch()
				builders.HasherBuilder.HashingPasswordSuccess()
				builders.AccountRepoBuilder.UpdatePasswordHashSuccess()
				return NewAccountUseCaseUT(t, builders)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			ctx := context.Background()

			err := sut.ChangePassword(ctx, input)

			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
