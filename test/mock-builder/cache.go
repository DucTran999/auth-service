package mockbuilder

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/DucTran999/auth-service/internal/domain"
	"github.com/DucTran999/auth-service/pkg/cache"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrGetCache        = errors.New("failed to get cache")
	ErrSetCacheSession = errors.New("failed to set cache session")
	ErrMissingKeys     = errors.New("failed to get list missing keys")
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
	sessionCached := domain.Session{
		ID:        FakeSessionID,
		AccountID: FakeAccountID,
		Account: domain.Account{
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
			if ptr, ok := dest.(*domain.Session); ok {
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
		Del(mock.Anything, mock.MatchedBy(func(key string) bool {
			return strings.HasPrefix(key, cache.SessionKeyPrefix)
		})).
		Return(nil)
}
