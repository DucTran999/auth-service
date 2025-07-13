package mockbuilder

import (
	"errors"
	"testing"

	"github.com/DucTran999/auth-service/internal/errs"
	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	ErrSessionValidate = errors.New("failed to validate session")
)

type mockSessionUsecase struct {
	inst *mocks.SessionUsecase
}

func newMockSessionUsecase(t *testing.T) *mockSessionUsecase {
	return &mockSessionUsecase{
		inst: mocks.NewSessionUsecase(t),
	}
}

func (m *mockSessionUsecase) GetInstance() *mocks.SessionUsecase {
	return m.inst
}

func (m *mockSessionUsecase) ValidateError() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(nil, ErrSessionValidate)
}

func (m *mockSessionUsecase) ValidateInvalidSession() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(nil, errs.ErrInvalidSessionID)
}

func (m *mockSessionUsecase) ValidateSessionSuccess() {
	m.inst.EXPECT().
		Validate(mock.Anything, mock.Anything).
		Return(&model.Session{
			ID:        FakeSessionID,
			AccountID: FakeAccountID,
		}, nil)
}
