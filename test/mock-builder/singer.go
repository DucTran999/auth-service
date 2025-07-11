package mockbuilder

import (
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

var (
	ErrSigningToken = errors.New("unexpected error while signing token")
	ErrParsingToken = errors.New("unexpected error while parsing token")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type mockSignerBuilder struct {
	inst *mocks.TokenSigner
}

func NewMockSignerBuilder(t *testing.T) *mockSignerBuilder {
	return &mockSignerBuilder{
		inst: mocks.NewTokenSigner(t),
	}
}

func (b *mockSignerBuilder) GetInstance() *mocks.TokenSigner {
	return b.inst
}

// Sign method mocks
func (b *mockSignerBuilder) SignSuccess() {
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", nil)
}

func (b *mockSignerBuilder) SignAccessSuccessAndSignRefreshFailed() {
	// 1. Access token sign → success
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", nil).
		Times(1)

	// 2. Refresh token sign → fail
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return("", ErrSigningToken).
		Times(1)
}

func (b *mockSignerBuilder) SignFailed() {
	b.inst.EXPECT().
		Sign(mock.Anything).
		Return("", ErrSigningToken)
}

// ParseInto method mocks
func (b *mockSignerBuilder) ParseIntoSuccess() {
	fakeClaims := model.TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "abc",
			ID:      "token-abc",
		},
	}
	b.inst.EXPECT().
		ParseInto(mock.AnythingOfType("string"), mock.AnythingOfType("*model.TokenClaims")).
		Run(func(tokenStr string, dest jwt.Claims) {
			// Set the value inside the provided *model.TokenClaims
			if out, ok := dest.(*model.TokenClaims); ok {
				*out = fakeClaims
			}
		}).
		Return(nil)
}

func (b *mockSignerBuilder) ParseIntoInvalidToken() {
	b.inst.EXPECT().
		ParseInto(mock.AnythingOfType("string"), mock.Anything).
		Return(ErrInvalidToken)
}

func (b *mockSignerBuilder) ParseIntoExpiredToken() {
	b.inst.EXPECT().
		ParseInto(mock.AnythingOfType("string"), mock.Anything).
		Return(ErrExpiredToken)
}

func (b *mockSignerBuilder) ParseIntoFailed() {
	b.inst.EXPECT().
		ParseInto(mock.AnythingOfType("string"), mock.Anything).
		Return(ErrParsingToken)
}

// Specific token and claims combinations
func (b *mockSignerBuilder) ParseIntoSpecificToken(tokenStr string) {
	b.inst.EXPECT().
		ParseInto(tokenStr, mock.Anything).
		Return(nil)
}

func (b *mockSignerBuilder) SignSpecificClaims(claims jwt.Claims, expectedToken string) {
	b.inst.EXPECT().
		Sign(claims).
		Return(expectedToken, nil)
}

// Chain methods for common scenarios
func (b *mockSignerBuilder) SignAndParseSuccess() {
	b.SignSuccess()
	b.ParseIntoSuccess()
}

func (b *mockSignerBuilder) SignSuccessParseIntoFailed() {
	b.SignSuccess()
	b.ParseIntoInvalidToken()
}
