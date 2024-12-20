// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/BariVakhidov/rssaggregator/internal/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockUserStorage is an autogenerated mock type for the UserStorage type
type MockUserStorage struct {
	mock.Mock
}

type MockUserStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserStorage) EXPECT() *MockUserStorage_Expecter {
	return &MockUserStorage_Expecter{mock: &_m.Mock}
}

// ChangePendingUserStatus provides a mock function with given fields: ctx, userId, status
func (_m *MockUserStorage) ChangePendingUserStatus(ctx context.Context, userId uuid.UUID, status string) (model.PendingUser, error) {
	ret := _m.Called(ctx, userId, status)

	if len(ret) == 0 {
		panic("no return value specified for ChangePendingUserStatus")
	}

	var r0 model.PendingUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) (model.PendingUser, error)); ok {
		return rf(ctx, userId, status)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) model.PendingUser); ok {
		r0 = rf(ctx, userId, status)
	} else {
		r0 = ret.Get(0).(model.PendingUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, userId, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserStorage_ChangePendingUserStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChangePendingUserStatus'
type MockUserStorage_ChangePendingUserStatus_Call struct {
	*mock.Call
}

// ChangePendingUserStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - userId uuid.UUID
//   - status string
func (_e *MockUserStorage_Expecter) ChangePendingUserStatus(ctx interface{}, userId interface{}, status interface{}) *MockUserStorage_ChangePendingUserStatus_Call {
	return &MockUserStorage_ChangePendingUserStatus_Call{Call: _e.mock.On("ChangePendingUserStatus", ctx, userId, status)}
}

func (_c *MockUserStorage_ChangePendingUserStatus_Call) Run(run func(ctx context.Context, userId uuid.UUID, status string)) *MockUserStorage_ChangePendingUserStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *MockUserStorage_ChangePendingUserStatus_Call) Return(_a0 model.PendingUser, _a1 error) *MockUserStorage_ChangePendingUserStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserStorage_ChangePendingUserStatus_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) (model.PendingUser, error)) *MockUserStorage_ChangePendingUserStatus_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: ctx, userID
func (_m *MockUserStorage) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*model.User, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *model.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserStorage_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type MockUserStorage_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
func (_e *MockUserStorage_Expecter) GetUser(ctx interface{}, userID interface{}) *MockUserStorage_GetUser_Call {
	return &MockUserStorage_GetUser_Call{Call: _e.mock.On("GetUser", ctx, userID)}
}

func (_c *MockUserStorage_GetUser_Call) Run(run func(ctx context.Context, userID uuid.UUID)) *MockUserStorage_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockUserStorage_GetUser_Call) Return(_a0 *model.User, _a1 error) *MockUserStorage_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserStorage_GetUser_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*model.User, error)) *MockUserStorage_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// PendingUserByEmail provides a mock function with given fields: ctx, userEmail
func (_m *MockUserStorage) PendingUserByEmail(ctx context.Context, userEmail string) (model.PendingUser, error) {
	ret := _m.Called(ctx, userEmail)

	if len(ret) == 0 {
		panic("no return value specified for PendingUserByEmail")
	}

	var r0 model.PendingUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.PendingUser, error)); ok {
		return rf(ctx, userEmail)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.PendingUser); ok {
		r0 = rf(ctx, userEmail)
	} else {
		r0 = ret.Get(0).(model.PendingUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userEmail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserStorage_PendingUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PendingUserByEmail'
type MockUserStorage_PendingUserByEmail_Call struct {
	*mock.Call
}

// PendingUserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - userEmail string
func (_e *MockUserStorage_Expecter) PendingUserByEmail(ctx interface{}, userEmail interface{}) *MockUserStorage_PendingUserByEmail_Call {
	return &MockUserStorage_PendingUserByEmail_Call{Call: _e.mock.On("PendingUserByEmail", ctx, userEmail)}
}

func (_c *MockUserStorage_PendingUserByEmail_Call) Run(run func(ctx context.Context, userEmail string)) *MockUserStorage_PendingUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockUserStorage_PendingUserByEmail_Call) Return(_a0 model.PendingUser, _a1 error) *MockUserStorage_PendingUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserStorage_PendingUserByEmail_Call) RunAndReturn(run func(context.Context, string) (model.PendingUser, error)) *MockUserStorage_PendingUserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// SavePendingUser provides a mock function with given fields: ctx, pendingUserInfo
func (_m *MockUserStorage) SavePendingUser(ctx context.Context, pendingUserInfo model.PendingUserInfo) (model.PendingUser, error) {
	ret := _m.Called(ctx, pendingUserInfo)

	if len(ret) == 0 {
		panic("no return value specified for SavePendingUser")
	}

	var r0 model.PendingUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.PendingUserInfo) (model.PendingUser, error)); ok {
		return rf(ctx, pendingUserInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.PendingUserInfo) model.PendingUser); ok {
		r0 = rf(ctx, pendingUserInfo)
	} else {
		r0 = ret.Get(0).(model.PendingUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.PendingUserInfo) error); ok {
		r1 = rf(ctx, pendingUserInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserStorage_SavePendingUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SavePendingUser'
type MockUserStorage_SavePendingUser_Call struct {
	*mock.Call
}

// SavePendingUser is a helper method to define mock.On call
//   - ctx context.Context
//   - pendingUserInfo model.PendingUserInfo
func (_e *MockUserStorage_Expecter) SavePendingUser(ctx interface{}, pendingUserInfo interface{}) *MockUserStorage_SavePendingUser_Call {
	return &MockUserStorage_SavePendingUser_Call{Call: _e.mock.On("SavePendingUser", ctx, pendingUserInfo)}
}

func (_c *MockUserStorage_SavePendingUser_Call) Run(run func(ctx context.Context, pendingUserInfo model.PendingUserInfo)) *MockUserStorage_SavePendingUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.PendingUserInfo))
	})
	return _c
}

func (_c *MockUserStorage_SavePendingUser_Call) Return(_a0 model.PendingUser, _a1 error) *MockUserStorage_SavePendingUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserStorage_SavePendingUser_Call) RunAndReturn(run func(context.Context, model.PendingUserInfo) (model.PendingUser, error)) *MockUserStorage_SavePendingUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserStorage creates a new instance of MockUserStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserStorage {
	mock := &MockUserStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
