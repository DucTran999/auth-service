package mockbuilder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DucTran999/auth-service/internal/model"
	"github.com/DucTran999/auth-service/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var (
	FakeSessionID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	FakeAccountID = uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	ErrCreateSession        = errors.New("unexpected create session error")
	ErrFindSessionByID      = errors.New("unexpected find session error")
	ErrUpdateSessionExpires = errors.New("unexpected update expires error")
)

type mockSessionRepoBuilder struct {
	inst *mocks.SessionRepository
}

func newMockSessionRepoBuilder(t *testing.T) *mockSessionRepoBuilder {
	return &mockSessionRepoBuilder{
		inst: mocks.NewSessionRepository(t),
	}
}

func (b *mockSessionRepoBuilder) GetInstance() *mocks.SessionRepository {
	return b.inst
}

func (blr *mockSessionRepoBuilder) FindSessionError() {
	blr.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(nil, ErrFindSessionByID)
}

func (blr *mockSessionRepoBuilder) FindSessionSuccess() {
	mockSession := model.Session{
		ID:        FakeSessionID,
		AccountID: FakeAccountID,
		Account: model.Account{
			ID:       FakeAccountID,
			Email:    FakeEmail,
			IsActive: true,
		},
	}

	blr.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(&mockSession, nil)
}

func (blr *mockSessionRepoBuilder) FindExpiredSession() {
	expiredAt := time.Now().Add(-1 * time.Hour)
	mockSession := model.Session{
		ID:        FakeSessionID,
		AccountID: FakeAccountID,
		Account: model.Account{
			ID:       FakeAccountID,
			Email:    FakeEmail,
			IsActive: true,
		},
		ExpiresAt: &expiredAt,
	}

	blr.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(&mockSession, nil)
}

func (blr *mockSessionRepoBuilder) NotFoundSession() {
	blr.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(nil, nil)
}

func (blr *mockSessionRepoBuilder) FindSessionReuse() {
	mockExpires := time.Now().Add(time.Minute)
	blr.inst.EXPECT().
		FindByID(mock.Anything, mock.Anything).
		Return(&model.Session{
			ID:        FakeSessionID,
			AccountID: FakeAccountID,
			ExpiresAt: &mockExpires,
		}, nil)
}

func (blr *mockSessionRepoBuilder) UpdateExpiresAtError() {
	blr.inst.EXPECT().
		UpdateExpiresAt(mock.Anything, mock.Anything, mock.Anything).
		Return(ErrUpdateSessionExpires)
}

func (blr *mockSessionRepoBuilder) UpdateExpiresAtSuccess() {
	blr.inst.EXPECT().
		UpdateExpiresAt(mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func (blr *mockSessionRepoBuilder) CreateSessionSuccess() {
	blr.inst.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*model.Session")).
		Run(func(ctx context.Context, s *model.Session) {
			s.ID = FakeSessionID
		}).
		Return(nil)
}

func (blr *mockSessionRepoBuilder) CreateSessionFailed() {
	blr.inst.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(ErrCreateSession)
}
