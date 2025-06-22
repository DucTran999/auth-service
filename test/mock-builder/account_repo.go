package mockbuilder

import (
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	FakeEmail             = "daniel@example.com"
	ErrFindAccountByEmail = errors.New("find unexpected error")
)

type mockAccountRepoBuilder struct {
	inst *mocks.AccountRepo
}

func newMockAccountRepoBuilder(t *testing.T) *mockAccountRepoBuilder {
	return &mockAccountRepoBuilder{
		inst: mocks.NewAccountRepo(t),
	}
}

func (b *mockAccountRepoBuilder) GetInstance() *mocks.AccountRepo {
	return b.inst
}

func (b *mockAccountRepoBuilder) FindByEmailError() {
	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, ErrFindAccountByEmail)
}

func (b *mockAccountRepoBuilder) FindByEmailHasResult() {
	activeAccount := &model.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: true,
	}

	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(activeAccount, nil)
}

func (b *mockAccountRepoBuilder) FindByEmailAccountInactive() {
	mockAccount := &model.Account{
		ID:       FakeAccountID,
		Email:    FakeEmail,
		IsActive: false,
	}

	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(mockAccount, nil)
}

func (b *mockAccountRepoBuilder) FindByEmailNoResult() {
	b.inst.EXPECT().
		FindByEmail(mock.Anything, mock.Anything).
		Return(nil, nil)
}
