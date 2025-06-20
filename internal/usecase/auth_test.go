package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type authUseCaseUT struct {
	ut usecase.AuthUseCase

	hasher      *mocks.Hasher
	accountRepo *mocks.AccountRepo
	sessionRepo *mocks.SessionRepository
}

func NewAuthUseCaseUT() *authUseCaseUT {
	hasher := new(mocks.Hasher)
	accountRepo := new(mocks.AccountRepo)
	sessionRepo := new(mocks.SessionRepository)

	return &authUseCaseUT{
		ut:          usecase.NewAuthUseCase(hasher, accountRepo, sessionRepo),
		hasher:      hasher,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
	}
}

func (sut *authUseCaseUT) mockFindByEmailFailed() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, errors.New("find user by email: unexpected error"))
}

func (sut *authUseCaseUT) mockFindByEmailHasResult() {
	activeAccount := &model.Account{
		ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Email:    "daniel@example.com",
		IsActive: true,
	}

	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(activeAccount, nil)
}

func (sut *authUseCaseUT) mockFindByEmailAccountInactive() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(&model.Account{Email: "daniel@example.com", IsActive: false}, nil)
}

func (sut *authUseCaseUT) mockFindByEmailNoResult() {
	sut.accountRepo.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (sut *authUseCaseUT) mockHashPasswordMatch() {
	sut.hasher.EXPECT().
		ComparePasswordAndHash(mock.Anything, mock.AnythingOfType("string")).
		Return(true, nil)
}

func (sut *authUseCaseUT) mockHashPasswordNotMatch() {
	sut.hasher.EXPECT().
		ComparePasswordAndHash(mock.AnythingOfType("string"), mock.Anything).
		Return(false, nil)
}

func (sut *authUseCaseUT) mockHashPasswordGotError() {
	sut.hasher.EXPECT().
		ComparePasswordAndHash(mock.AnythingOfType("string"), mock.Anything).
		Return(false, errors.New("compare password unexpected error"))
}

func (sut *authUseCaseUT) mockCreateSessionFailed() {
	sut.sessionRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(errors.New("create session failed"))
}

func (sut *authUseCaseUT) mockFindSessionError() {
	sut.sessionRepo.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(nil, errors.New("failed to find session"))
}

func (sut *authUseCaseUT) mockSessionNotFound() {
	sut.sessionRepo.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (sut *authUseCaseUT) mockSessionCanReuse() {
	mockExpires := time.Now().Add(time.Minute)
	sut.sessionRepo.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(&model.Session{
			ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			AccountID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			ExpiresAt: &mockExpires,
		}, nil)
}

func (sut *authUseCaseUT) mockSessionUpdateExpiresAtErr() {
	sut.sessionRepo.EXPECT().
		UpdateExpiresAt(mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("update expires error"))
}

func (sut *authUseCaseUT) mockCreateSessionSuccess() {
	sut.sessionRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*model.Session")).
		Run(func(ctx context.Context, s *model.Session) {
			s.ID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		}).
		Return(nil)
}

func (sut *authUseCaseUT) mockSessionUpdateExpiresAt() {
	sut.sessionRepo.EXPECT().
		UpdateExpiresAt(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func TestLogin(t *testing.T) {
	type testCase struct {
		name        string
		sut         *authUseCaseUT
		loginInput  usecase.LoginInput
		expectedErr error
		expected    *model.Account
	}

	userSample := usecase.LoginInput{
		Email:    "daniel@example.com",
		Password: "abc1234!",
	}

	testTable := []testCase{
		{
			name: "failed to find email in db",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailFailed()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: errors.New("find user by email: unexpected error"),
			expected:    nil,
		},
		{
			name: "failed email not existed",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailNoResult()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: usecase.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "account inactive",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailAccountInactive()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: usecase.ErrAccountDisabled,
			expected:    nil,
		},
		{
			name: "password not match",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailHasResult()
				sut.mockHashPasswordNotMatch()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: usecase.ErrInvalidCredentials,
			expected:    nil,
		},
		{
			name: "compare password unexpected error",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailHasResult()
				sut.mockHashPasswordGotError()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: errors.New("compare password unexpected error"),
			expected:    nil,
		},
		{
			name: "failed to find session",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindSessionError()
				return sut
			}(),
			loginInput: usecase.LoginInput{
				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
				Email:            userSample.Email,
				Password:         userSample.Password,
			},
			expectedErr: errors.New("failed to find session"),
			expected:    nil,
		},
		{
			name: "session not found allow create new one",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockSessionNotFound()
				sut.mockFindByEmailHasResult()
				sut.mockHashPasswordMatch()
				sut.mockCreateSessionSuccess()
				return sut
			}(),
			loginInput: usecase.LoginInput{
				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
				Email:            userSample.Email,
				Password:         userSample.Password,
			},
			expectedErr: nil,
			expected: &model.Account{
				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Email:    "daniel@example.com",
				IsActive: true,
			},
		},
		{
			name: "failed when update expires",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockSessionCanReuse()
				sut.mockSessionUpdateExpiresAtErr()
				return sut
			}(),
			loginInput: usecase.LoginInput{
				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
				Email:            userSample.Email,
				Password:         userSample.Password,
			},
			expectedErr: errors.New("update expires error"),
			expected:    nil,
		},
		{
			name: "reuse session",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockSessionCanReuse()
				sut.mockSessionUpdateExpiresAt()
				return sut
			}(),
			loginInput: usecase.LoginInput{
				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
				Email:            userSample.Email,
				Password:         userSample.Password,
			},
			expectedErr: nil,
			expected: &model.Account{
				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Email:    "daniel@example.com",
				IsActive: true,
			},
		},
		{
			name: "failed to create session",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailHasResult()
				sut.mockHashPasswordMatch()
				sut.mockCreateSessionFailed()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: errors.New("create session failed"),
			expected:    nil,
		},
		{
			name: "login success",
			sut: func() *authUseCaseUT {
				sut := NewAuthUseCaseUT()
				sut.mockFindByEmailHasResult()
				sut.mockHashPasswordMatch()
				sut.mockCreateSessionSuccess()
				return sut
			}(),
			loginInput:  userSample,
			expectedErr: nil,
			expected: &model.Account{
				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Email:    "daniel@example.com",
				IsActive: true,
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			account, err := tc.sut.ut.Login(context.Background(), tc.loginInput)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.ID, account.ID)
			} else {
				assert.Nil(t, account)
			}
		})
	}
}
