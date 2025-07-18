// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"

	"github.com/DucTran999/auth-service/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// NewSessionUsecase creates a new instance of SessionUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSessionUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *SessionUsecase {
	mock := &SessionUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// SessionUsecase is an autogenerated mock type for the SessionUsecase type
type SessionUsecase struct {
	mock.Mock
}

type SessionUsecase_Expecter struct {
	mock *mock.Mock
}

func (_m *SessionUsecase) EXPECT() *SessionUsecase_Expecter {
	return &SessionUsecase_Expecter{mock: &_m.Mock}
}

// Validate provides a mock function for the type SessionUsecase
func (_mock *SessionUsecase) Validate(ctx context.Context, sessionID string) (*model.Session, error) {
	ret := _mock.Called(ctx, sessionID)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 *model.Session
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (*model.Session, error)); ok {
		return returnFunc(ctx, sessionID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) *model.Session); ok {
		r0 = returnFunc(ctx, sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Session)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, sessionID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// SessionUsecase_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type SessionUsecase_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - ctx context.Context
//   - sessionID string
func (_e *SessionUsecase_Expecter) Validate(ctx interface{}, sessionID interface{}) *SessionUsecase_Validate_Call {
	return &SessionUsecase_Validate_Call{Call: _e.mock.On("Validate", ctx, sessionID)}
}

func (_c *SessionUsecase_Validate_Call) Run(run func(ctx context.Context, sessionID string)) *SessionUsecase_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *SessionUsecase_Validate_Call) Return(session *model.Session, err error) *SessionUsecase_Validate_Call {
	_c.Call.Return(session, err)
	return _c
}

func (_c *SessionUsecase_Validate_Call) RunAndReturn(run func(ctx context.Context, sessionID string) (*model.Session, error)) *SessionUsecase_Validate_Call {
	_c.Call.Return(run)
	return _c
}
