package usecase_test

import (
	"context"
	"log"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/internal/usecase"
	mockbuilder "github.com/DucTran999/auth-service/test/mock-builder"
	"github.com/stretchr/testify/assert"
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
			// 	name: "failed to find email in db",
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
				builders.HasherBuilder.HashPasswordGotError()
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
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			sut := tc.setup(t)
			session, err := sut.Login(context.Background(), tc.loginInput)

			assert.Equal(t, tc.expectedErr, err)
			if tc.expected != nil {
				log.Println(tc.expected.ID)
				assert.Equal(t, tc.expected.ID, session.AccountID)
			} else {
				assert.Nil(t, session)
			}
		})
	}
}

// func TestLogin(t *testing.T) {

// 	userSample := usecase.LoginInput{
// 		Email:    "daniel@example.com",
// 		Password: "abc1234!",
// 	}

// 	testTable := []testCase{

// 		{
// 			name: "session not found allow create new one",
// 			sut: func() *authUseCaseUT {
// 				sut := NewAuthUseCaseUT()
// 				sut.mockSessionNotFound()
// 				sut.mockFindByEmailHasResult()
// 				sut.mockHashPasswordMatch()
// 				sut.mockCreateSessionSuccess()
// 				return sut
// 			}(),
// 			loginInput: usecase.LoginInput{
// 				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
// 				Email:            userSample.Email,
// 				Password:         userSample.Password,
// 			},
// 			expectedErr: nil,
// 			expected: &model.Account{
// 				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
// 				Email:    "daniel@example.com",
// 				IsActive: true,
// 			},
// 		},
// 		{
// 			name: "failed when update expires",
// 			sut: func() *authUseCaseUT {
// 				sut := NewAuthUseCaseUT()
// 				sut.mockSessionCanReuse()
// 				sut.mockSessionUpdateExpiresAtErr()
// 				return sut
// 			}(),
// 			loginInput: usecase.LoginInput{
// 				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
// 				Email:            userSample.Email,
// 				Password:         userSample.Password,
// 			},
// 			expectedErr: errors.New("update expires error"),
// 			expected:    nil,
// 		},
// 		{
// 			name: "reuse session",
// 			sut: func() *authUseCaseUT {
// 				sut := NewAuthUseCaseUT()
// 				sut.mockSessionCanReuse()
// 				sut.mockSessionUpdateExpiresAt()
// 				return sut
// 			}(),
// 			loginInput: usecase.LoginInput{
// 				CurrentSessionID: "123e4567-e89b-12d3-a456-426614174000",
// 				Email:            userSample.Email,
// 				Password:         userSample.Password,
// 			},
// 			expectedErr: nil,
// 			expected: &model.Account{
// 				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
// 				Email:    "daniel@example.com",
// 				IsActive: true,
// 			},
// 		},
// 		{
// 			name: "failed to create session",
// 			sut: func() *authUseCaseUT {
// 				sut := NewAuthUseCaseUT()
// 				sut.mockFindByEmailHasResult()
// 				sut.mockHashPasswordMatch()
// 				sut.mockCreateSessionFailed()
// 				return sut
// 			}(),
// 			loginInput:  userSample,
// 			expectedErr: errors.New("create session failed"),
// 			expected:    nil,
// 		},
// 		{
// 			name: "login success",
// 			sut: func() *authUseCaseUT {
// 				sut := NewAuthUseCaseUT()
// 				sut.mockFindByEmailHasResult()
// 				sut.mockHashPasswordMatch()
// 				sut.mockCreateSessionSuccess()
// 				return sut
// 			}(),
// 			loginInput:  userSample,
// 			expectedErr: nil,
// 			expected: &model.Account{
// 				ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
// 				Email:    "daniel@example.com",
// 				IsActive: true,
// 			},
// 		},
// 	}
// }
