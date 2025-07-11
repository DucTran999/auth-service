package mockbuilder

import (
	"context"
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrGetCache        = errors.New("failed to get cache")
	ErrSetCache        = errors.New("failed to set cache")
	ErrSetCacheSession = errors.New("failed to set cache session")
	ErrMissingKeys     = errors.New("failed to get list missing keys")
	ErrDelCacheDelete  = errors.New("failed to del cache key")
	ErrHasCache        = errors.New("failed to check key in cache")
)

type mockCacheBuilder struct {
	inst *mocks.Cache
}

func (b *mockCacheBuilder) GetInstance() *mocks.Cache {
	return b.inst
}

func newMockCacheBuilder(t *testing.T) *mockCacheBuilder {
	return &mockCacheBuilder{
		inst: mocks.NewCache(t),
	}
}

func (b *mockCacheBuilder) GetCacheErr() {
	b.inst.EXPECT().
		GetInto(mock.Anything, mock.Anything, mock.Anything).
		Return(ErrGetCache)
}

func (b *mockCacheBuilder) ValidSessionCached() {
	sessionCached := model.Session{
		ID:        FakeSessionID,
		AccountID: FakeAccountID,
		Account: model.Account{
			ID:       FakeAccountID,
			Email:    FakeEmail,
			IsActive: true,
		},
		ExpiresAt: nil,
	}

	b.inst.EXPECT().
		GetInto(mock.Anything, mock.Anything, mock.Anything).
		Run(func(ctx context.Context, key string, dest any) {
			// Type assert to pointer type
			if ptr, ok := dest.(*model.Session); ok {
				*ptr = sessionCached // Copy value into the pointed-to object
			}
		}).
		Return(nil)
}

func (b *mockCacheBuilder) SessionMissCache() {
	b.inst.EXPECT().
		GetInto(mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("missing key"))
}

func (b *mockCacheBuilder) SetCacheSessionSuccess() {
	b.inst.EXPECT().
		Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func (b *mockCacheBuilder) SetCacheSessionFailed() {
	b.inst.EXPECT().
		Set(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(ErrSetCacheSession)
}

func (b *mockCacheBuilder) SetExpireSuccess() {
	b.inst.EXPECT().
		Expire(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func (b *mockCacheBuilder) GetTTLSuccess() {
	b.inst.EXPECT().
		TTL(mock.Anything, mock.Anything).
		Return(555, nil)
}

func (b *mockCacheBuilder) CallMissingKeysFailed() {
	b.inst.EXPECT().
		MissingKeys(mock.Anything, mock.Anything).
		Return(nil, ErrMissingKeys)
}

func (b *mockCacheBuilder) CallMissingKeysSuccess() {
	cachedKey := cache.KeyFromSessionID(FakeSessionID.String())
	missingKeys := []string{cachedKey}

	b.inst.EXPECT().
		MissingKeys(mock.Anything, mock.Anything).
		Return(missingKeys, nil)
}

func (b *mockCacheBuilder) NoMissingKeysFound() {
	b.inst.EXPECT().
		MissingKeys(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (b *mockCacheBuilder) DelKeySuccess() {
	b.inst.EXPECT().
		Del(mock.Anything, mock.AnythingOfType("string")).
		Return(nil)
}

func (b *mockCacheBuilder) DelKeyErr() {
	b.inst.EXPECT().
		Del(mock.Anything, mock.AnythingOfType("string")).
		Return(ErrDelCacheDelete)
}

func (b *mockCacheBuilder) CheckRefreshTokenFailed() {
	b.inst.EXPECT().
		Has(mock.Anything, mock.AnythingOfType("string")).
		Return(false, ErrHasCache)
}

func (b *mockCacheBuilder) RefreshTokenInvalidCache() {
	b.inst.EXPECT().
		Has(mock.Anything, mock.AnythingOfType("string")).
		Return(false, nil)
}

func (b *mockCacheBuilder) RefreshTokenValidCache() {
	b.inst.EXPECT().
		Has(mock.Anything, mock.AnythingOfType("string")).
		Return(true, nil)
}

func (b *mockCacheBuilder) SetRefreshTokenFailed() {
	b.inst.EXPECT().
		Set(mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(ErrSetCache)
}

func (b *mockCacheBuilder) SetRefreshTokenSuccess() {
	b.inst.EXPECT().
		Set(mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything).
		Return(nil)
}
