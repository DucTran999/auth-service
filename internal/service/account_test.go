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

type userBizUT struct {
	ut    service.IUserService
	uRepo *mocks.IUserRepo
}

func NewUserBizUT() *userBizUT {
	uRepo := new(mocks.IUserRepo)

	return &userBizUT{
		ut:    service.NewUserBiz(uRepo),
		uRepo: uRepo,
	}
}

func (sut *userBizUT) mockGetUserByEmailFailed() {
	sut.uRepo.EXPECT().
		GetUserByEmail(mock.Anything, mock.Anything).
		Return(nil, errors.New("get user by email unexpected error"))
}

func (sut *userBizUT) mockGetUserByEmailHasResult() {
	sut.uRepo.EXPECT().
		GetUserByEmail(mock.Anything, mock.Anything).
		Return(&model.User{Username: "daniel", Email: "daniel@example.com"}, nil)
}

func (sut *userBizUT) mockGetUserByEmailNoResult() {
	sut.uRepo.EXPECT().
		GetUserByEmail(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (sut *userBizUT) mockCreateUserErr() {
	sut.uRepo.EXPECT().
		CreateUser(mock.Anything, mock.Anything).
		Return(nil, errors.New("create user unexpected err"))
}

func (sut *userBizUT) mockCreateUserSuccess() {
	sut.uRepo.EXPECT().
		CreateUser(mock.Anything, mock.Anything).
		Return(&model.User{Username: "daniel", Email: "daniel@example.com"}, nil)
}

func TestRegisterUser(t *testing.T) {
	type testCase struct {
		name        string
		sut         *userBizUT
		userInfo    model.User
		expectedErr error
		expected    *model.User
	}

	userSample := model.User{
		Username: "daniel",
		Email:    "daniel@example.com",
	}

	testTable := []testCase{
		{
			name: "WhenGetUserByEmailGotErr_ThenReturnErr",
			sut: func() *userBizUT {
				sut := NewUserBizUT()
				sut.mockGetUserByEmailFailed()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: errors.New("get user by email unexpected error"),
			expected:    nil,
		},
		{
			name: "WhenEmailUsed_ThenReturnExistedErr",
			sut: func() *userBizUT {
				sut := NewUserBizUT()
				sut.mockGetUserByEmailHasResult()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: common.ErrEmailExisted,
			expected:    nil,
		},
		{
			name: "WhenCreateUserGotErr_ThenReturnErr",
			sut: func() *userBizUT {
				sut := NewUserBizUT()
				sut.mockGetUserByEmailNoResult()
				sut.mockCreateUserErr()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: errors.New("create user unexpected err"),
			expected:    nil,
		},
		{
			name: "RegisterSuccess",
			sut: func() *userBizUT {
				sut := NewUserBizUT()
				sut.mockGetUserByEmailNoResult()
				sut.mockCreateUserSuccess()
				return sut
			}(),
			userInfo:    userSample,
			expectedErr: nil,
			expected:    &userSample,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			user, err := tc.sut.ut.RegisterUser(context.Background(), tc.userInfo)

			assert.Equal(t, tc.expected, user)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
