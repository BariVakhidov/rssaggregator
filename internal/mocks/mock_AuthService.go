// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	ssov1 "github.com/BariVakhidov/ssoprotos/gen/go/sso"
	mock "github.com/stretchr/testify/mock"
)

// MockAuthService is an autogenerated mock type for the AuthService type
type MockAuthService struct {
	mock.Mock
}

type MockAuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthService) EXPECT() *MockAuthService_Expecter {
	return &MockAuthService_Expecter{mock: &_m.Mock}
}

// Login provides a mock function with given fields: ctx, req
func (_m *MockAuthService) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 *ssov1.LoginResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssov1.LoginRequest) (*ssov1.LoginResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssov1.LoginRequest) *ssov1.LoginResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssov1.LoginResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssov1.LoginRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type MockAuthService_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - req *ssov1.LoginRequest
func (_e *MockAuthService_Expecter) Login(ctx interface{}, req interface{}) *MockAuthService_Login_Call {
	return &MockAuthService_Login_Call{Call: _e.mock.On("Login", ctx, req)}
}

func (_c *MockAuthService_Login_Call) Run(run func(ctx context.Context, req *ssov1.LoginRequest)) *MockAuthService_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ssov1.LoginRequest))
	})
	return _c
}

func (_c *MockAuthService_Login_Call) Return(_a0 *ssov1.LoginResponse, _a1 error) *MockAuthService_Login_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_Login_Call) RunAndReturn(run func(context.Context, *ssov1.LoginRequest) (*ssov1.LoginResponse, error)) *MockAuthService_Login_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, req
func (_m *MockAuthService) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 *ssov1.RegisterResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssov1.RegisterRequest) *ssov1.RegisterResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssov1.RegisterResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssov1.RegisterRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAuthService_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type MockAuthService_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - req *ssov1.RegisterRequest
func (_e *MockAuthService_Expecter) Register(ctx interface{}, req interface{}) *MockAuthService_Register_Call {
	return &MockAuthService_Register_Call{Call: _e.mock.On("Register", ctx, req)}
}

func (_c *MockAuthService_Register_Call) Run(run func(ctx context.Context, req *ssov1.RegisterRequest)) *MockAuthService_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ssov1.RegisterRequest))
	})
	return _c
}

func (_c *MockAuthService_Register_Call) Return(_a0 *ssov1.RegisterResponse, _a1 error) *MockAuthService_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAuthService_Register_Call) RunAndReturn(run func(context.Context, *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error)) *MockAuthService_Register_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuthService creates a new instance of MockAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthService {
	mock := &MockAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
