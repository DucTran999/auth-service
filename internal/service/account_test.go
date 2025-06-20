package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/common"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/service"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type accountSvcUT struct {
	ut          service.AccountService
	accountRepo *mocks.AccountRepo
}

func NewAccountSvcUT() *accountSvcUT {
	accountRepo := new(mocks.AccountRepo)

	return &accountSvcUT{
		ut:          service.NewAccountService(accountRepo),
		accountRepo: accountRepo,
	}
}

func (sut *accountSvcUT) mockFindByEmailFailed() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, errors.New("find user by email: unexpected error"))
}

func (sut *accountSvcUT) mockFindByEmailHasResult() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(&model.Account{Email: "daniel@example.com"}, nil)
}

func (sut *accountSvcUT) mockFindByEmailNoResult() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (sut *accountSvcUT) mockCreateError() {
	sut.accountRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil, errors.New("create user: unexpected error"))
}

func (sut *accountSvcUT) mockCreateSuccess() {
	sut.accountRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(&model.Account{Email: "daniel@example.com"}, nil)
}

func TestRegisterAccount(t *testing.T) {
	type testCase struct {
		name        string
		sut         *accountSvcUT
		userInfo    model.Account
		expectedErr error
		expected    *model.Account
	}

	userSample := model.Account{
		Email: "daniel@example.com",
	}

	testTable := []testCase{
		{
			name: "WhenFindByEmailGotErr_ThenReturnErr",
			sut: func() *accountSvcUT {
				sut := NewAccountSvcUT()
				sut.mockFindByEmailFailed()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: errors.New("find user by email: unexpected error"),
			expected:    nil,
		},
		{
			name: "WhenEmailUsed_ThenReturnExistedErr",
			sut: func() *accountSvcUT {
				sut := NewAccountSvcUT()
				sut.mockFindByEmailHasResult()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: common.ErrEmailExisted,
			expected:    nil,
		},
		{
			name: "WhenCreateGotErr_ThenReturnErr",
			sut: func() *accountSvcUT {
				sut := NewAccountSvcUT()
				sut.mockFindByEmailNoResult()
				sut.mockCreateError()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: errors.New("create user: unexpected error"),
			expected:    nil,
		},
		{
			name: "RegisterSuccess",
			sut: func() *accountSvcUT {
				sut := NewAccountSvcUT()
				sut.mockFindByEmailNoResult()
				sut.mockCreateSuccess()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: nil,
			expected: &model.Account{
				Email: "daniel@example.com",
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			user, err := tc.sut.ut.Register(context.Background(), tc.userInfo)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Email, user.Email)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}
