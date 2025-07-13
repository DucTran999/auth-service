package mockbuilder

import (
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/usecase/dto"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrLoginJwt     = errors.New("verify credential got db error")
	ErrRevokeToken  = errors.New("failed to revoke token")
	ErrRefreshToken = errors.New("failed to refresh token")
)

type mockAuthJWTUsecase struct {
	inst *mocks.AuthJWTUsecase
}

func newMockAuthJWTUsecase(t *testing.T) *mockAuthJWTUsecase {
	return &mockAuthJWTUsecase{
		inst: mocks.NewAuthJWTUsecase(t),
	}
}

func (m *mockAuthJWTUsecase) GetInstance() *mocks.AuthJWTUsecase {
	return m.inst
}

func (m *mockAuthJWTUsecase) LoginErrWrongCredentials() {
	m.inst.EXPECT().
		Login(mock.Anything, mock.Anything).
		Return(nil, errs.ErrInvalidCredentials)
}

func (m *mockAuthJWTUsecase) LoginErrDB() {
	m.inst.EXPECT().
		Login(mock.Anything, mock.Anything).
		Return(nil, ErrLoginJwt)
}

func (m *mockAuthJWTUsecase) LoginSuccess() {
	tokens := dto.TokenPairs{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	m.inst.EXPECT().
		Login(mock.Anything, mock.Anything).
		Return(&tokens, nil)
}

func (m *mockAuthJWTUsecase) RevokeRefreshTokenErr() {
	m.inst.EXPECT().
		RevokeRefreshToken(mock.Anything, mock.Anything).
		Return(ErrRevokeToken)
}

func (m *mockAuthJWTUsecase) RevokeRefreshTokenSuccess() {
	m.inst.EXPECT().
		RevokeRefreshToken(mock.Anything, mock.Anything).
		Return(nil)
}

func (m *mockAuthJWTUsecase) RefreshTokenError() {
	m.inst.EXPECT().
		RefreshToken(mock.Anything, mock.Anything).
		Return(nil, ErrRefreshToken)
}

func (m *mockAuthJWTUsecase) RefreshTokenSuccess() {
	tokens := dto.TokenPairs{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}
	m.inst.EXPECT().
		RefreshToken(mock.Anything, mock.Anything).
		Return(&tokens, nil)
}
