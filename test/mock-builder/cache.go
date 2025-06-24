package mockbuilder

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrGetCache        = errors.New("failed to get cache")
	ErrSetCacheSession = errors.New("failed to set cache session")
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
		Run(func(ctx context.Context, key string, dest any) {
			dest = nil
		}).
		Return(fmt.Errorf("missing key"))
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
