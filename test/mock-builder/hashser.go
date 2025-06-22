package mockbuilder

import (
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrCompareHashPassword = errors.New("compare password unexpected error")
)

type mockHasherBuilder struct {
	inst *mocks.Hasher
}

func newMockHasherBuilder(t *testing.T) *mockHasherBuilder {
	return &mockHasherBuilder{
		inst: mocks.NewHasher(t),
	}
}

func (b *mockHasherBuilder) GetInstance() *mocks.Hasher {
	return b.inst
}

func (b *mockHasherBuilder) HashPasswordMatch() {
	b.inst.EXPECT().
		ComparePasswordAndHash(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(true, nil)
}

func (b *mockHasherBuilder) HashPasswordNotMatch() {
	b.inst.EXPECT().
		ComparePasswordAndHash(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(false, nil)
}

func (b *mockHasherBuilder) HashPasswordGotError() {
	b.inst.EXPECT().
		ComparePasswordAndHash(mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(false, ErrCompareHashPassword)
}
