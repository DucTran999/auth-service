package mockbuilder

import (
	"errors"
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockAuthSessionUsecase struct {
	inst *mocks.AuthSessionUsecase
}

func newMockAuthSessionUsecase(t *testing.T) *mockAuthSessionUsecase {
	return &mockAuthSessionUsecase{
		inst: mocks.NewAuthSessionUsecase(t),
	}
}

func (m *mockAuthSessionUsecase) GetInstance() *mocks.AuthSessionUsecase {
	return m.inst
}

func (m *mockAuthSessionUsecase) LoginInvalidCredentials() {
	m.inst.EXPECT().Login(mock.Anything, mock.Anything).
		Return(nil, errs.ErrInvalidCredentials)
}

func (m *mockAuthSessionUsecase) LoginInternalError() {
	m.inst.EXPECT().Login(mock.Anything, mock.Anything).
		Return(nil, errors.New("database error"))
}

func (m *mockAuthSessionUsecase) LoginSuccess() {
	m.inst.EXPECT().Login(mock.Anything, mock.Anything).
		Return(&model.Session{
			ID:        FakeSessionID,
			AccountID: uuid.New(),
			CreatedAt: time.Now(),
			UserAgent: "test-agent",
			IPAddress: "127.0.0.1",
		}, nil)
}

func (m *mockAuthSessionUsecase) LogoutSuccess() {
	m.inst.EXPECT().
		Logout(mock.Anything, mock.Anything).
		Return(nil)
}

func (m *mockAuthSessionUsecase) LogoutError() {
	m.inst.EXPECT().
		Logout(mock.Anything, mock.Anything).
		Return(errors.New("logout failed"))
}
