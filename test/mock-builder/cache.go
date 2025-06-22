package mockbuilder

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

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
		Get(mock.Anything, mock.Anything).
		Return("", ErrGetCache)
}

func (b *mockCacheBuilder) ValidSessionCached() {
	notExpired := time.Now().Add(time.Hour)

	sessionCached := &model.Session{
		ID:        FakeSessionID,
		AccountID: FakeAccountID,
		Account: model.Account{
			ID:       FakeAccountID,
			Email:    FakeEmail,
			IsActive: true,
		},
		ExpiresAt: &notExpired,
	}
	bytes, _ := json.Marshal(sessionCached)

	b.inst.EXPECT().
		Get(mock.Anything, mock.Anything).
		Return(string(bytes), nil)
}

func (b *mockCacheBuilder) SessionMissCache() {
	b.inst.EXPECT().
		Get(mock.Anything, mock.Anything).
		Return("", nil)
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
