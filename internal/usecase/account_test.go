package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type accountUseCaseUT struct {
	ut          usecase.AccountUseCase
	hasher      *mocks.Hasher
	accountRepo *mocks.AccountRepo
}

func NewAccountUseCaseUT() *accountUseCaseUT {
	hasher := new(mocks.Hasher)
	accountRepo := new(mocks.AccountRepo)

	return &accountUseCaseUT{
		ut:          usecase.NewAccountUseCase(hasher, accountRepo),
		hasher:      hasher,
		accountRepo: accountRepo,
	}
}

func (sut *accountUseCaseUT) mockFindByEmailFailed() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, errors.New("find user by email: unexpected error"))
}

func (sut *accountUseCaseUT) mockFindByEmailHasResult() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(&model.Account{Email: "daniel@example.com"}, nil)
}

func (sut *accountUseCaseUT) mockFindByEmailNoResult() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (sut *accountUseCaseUT) mockHashPasswordErr() {
	sut.hasher.EXPECT().
		HashPassword(mock.AnythingOfType("string")).
		Return("", errors.New("hash got err"))
}

func (sut *accountUseCaseUT) mockHashPasswordSuccess() {
	sut.hasher.EXPECT().
		HashPassword(mock.AnythingOfType("string")).
		Return("hashedPassword", nil)
}

func (sut *accountUseCaseUT) mockCreateError() {
	sut.accountRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil, errors.New("create user: unexpected error"))
}

func (sut *accountUseCaseUT) mockCreateSuccess() {
	sut.accountRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(&model.Account{Email: "daniel@example.com"}, nil)
}

func TestRegisterAccount(t *testing.T) {
	type testCase struct {
		name        string
		sut         *accountUseCaseUT
		accountInfo usecase.RegisterInput
		expectedErr error
		expected    *model.Account
	}

	userSample := usecase.RegisterInput{
		Email:    "daniel@example.com",
		Password: "abc1234!",
	}

	testTable := []testCase{
		{
			name: "failed to find email in db",
			sut: func() *accountUseCaseUT {
				sut := NewAccountUseCaseUT()
				sut.mockFindByEmailFailed()
				return sut
			}(),
			accountInfo: userSample,
			expectedErr: errors.New("find user by email: unexpected error"),
			expected:    nil,
		},
		{
			name: "failed caused email already taken",
			sut: func() *accountUseCaseUT {
				sut := NewAccountUseCaseUT()
				sut.mockFindByEmailHasResult()
				return sut
			}(),
			accountInfo: userSample,
			expectedErr: usecase.ErrEmailExisted,
			expected:    nil,
		},
		{
			name: "failed when hash password",
			sut: func() *accountUseCaseUT {
				sut := NewAccountUseCaseUT()
				sut.mockFindByEmailNoResult()
				sut.mockHashPasswordErr()
				return sut
			}(),
			accountInfo: userSample,
			expectedErr: errors.New("hash got err"),
			expected:    nil,
		},
		{
			name: "failed when persist to db",
			sut: func() *accountUseCaseUT {
				sut := NewAccountUseCaseUT()
				sut.mockFindByEmailNoResult()
				sut.mockHashPasswordSuccess()
				sut.mockCreateError()
				return sut
			}(),
			accountInfo: userSample,
			expectedErr: errors.New("create user: unexpected error"),
			expected:    nil,
		},
		{
			name: "register success",
			sut: func() *accountUseCaseUT {
				sut := NewAccountUseCaseUT()
				sut.mockFindByEmailNoResult()
				sut.mockHashPasswordSuccess()
				sut.mockCreateSuccess()
				return sut
			}(),
			accountInfo: userSample,
			expectedErr: nil,
			expected: &model.Account{
				Email: "daniel@example.com",
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			user, err := tc.sut.ut.Register(context.Background(), tc.accountInfo)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Email, user.Email)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}
