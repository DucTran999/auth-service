package shared_test

import (
	"context"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase/shared"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/assert"
)

func NewAccountVerifierUT(t *testing.T, builders *mockbuilder.BuilderContainer) shared.AccountVerifier {
	return shared.NewAccountVerifier(
		builders.HasherBuilder.GetInstance(),
		builders.AccountRepoBuilder.GetInstance(),
	)
}

type LoginInput struct {
	Email    string
	Password string
}

func TestVerify(t *testing.T) {
	type testCase struct {
		name        string
		setup       func(t *testing.T) shared.AccountVerifier
		loginInput  LoginInput
		expected    *model.Account
		expectedErr error
	}

	loginInput := LoginInput{
		Email:    mockbuilder.FakeEmail,
		Password: mockbuilder.FakeOldPass,
	}

	testTable := []testCase{
		{
			name: "find email in db failed",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailError()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: mockbuilder.ErrFindAccountByEmail,
			expected:    nil,
		},
		{
			name: "email not existed",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailNoResult()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: model.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "account inactive",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailAccountInactive()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: model.ErrAccountDisabled,
			expected:    nil,
		},
		{
			name: "password not match",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordNotMatch()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: model.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "compare password unexpected error",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.CompareHashPasswordGotError()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: mockbuilder.ErrCompareHashPassword,
			expected:    nil,
		},
		{
			name: "verify success",
			setup: func(t *testing.T) shared.AccountVerifier {
				builders := mockbuilder.NewBuilderContainer(t)
				builders.AccountRepoBuilder.FindByEmailHasResult()
				builders.HasherBuilder.HashPasswordMatch()
				return NewAccountVerifierUT(t, builders)
			},
			loginInput:  loginInput,
			expectedErr: nil,
			expected: &model.Account{
				ID:    mockbuilder.FakeAccountID,
				Email: mockbuilder.FakeEmail,
			},
		},
	}

	ctx := context.Background()

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			account, err := sut.Verify(ctx, tc.loginInput.Email, tc.loginInput.Password)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.ID, account.ID)
			} else {
				assert.Nil(t, account)
			}
		})
	}
}
