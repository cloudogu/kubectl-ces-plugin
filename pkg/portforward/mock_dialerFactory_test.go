// Code generated by mockery v2.20.0. DO NOT EDIT.

package portforward

import mock "github.com/stretchr/testify/mock"

// mockDialerFactory is an autogenerated mock type for the dialerFactory type
type mockDialerFactory struct {
	mock.Mock
}

type mockDialerFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDialerFactory) EXPECT() *mockDialerFactory_Expecter {
	return &mockDialerFactory_Expecter{mock: &_m.Mock}
}

// create provides a mock function with given fields:
func (_m *mockDialerFactory) create() (dialer, error) {
	ret := _m.Called()

	var r0 dialer
	var r1 error
	if rf, ok := ret.Get(0).(func() (dialer, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() dialer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dialer)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDialerFactory_create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'create'
type mockDialerFactory_create_Call struct {
	*mock.Call
}

// create is a helper method to define mock.On call
func (_e *mockDialerFactory_Expecter) create() *mockDialerFactory_create_Call {
	return &mockDialerFactory_create_Call{Call: _e.mock.On("create")}
}

func (_c *mockDialerFactory_create_Call) Run(run func()) *mockDialerFactory_create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDialerFactory_create_Call) Return(_a0 dialer, _a1 error) *mockDialerFactory_create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDialerFactory_create_Call) RunAndReturn(run func() (dialer, error)) *mockDialerFactory_create_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDialerFactory interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDialerFactory creates a new instance of mockDialerFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDialerFactory(t mockConstructorTestingTnewMockDialerFactory) *mockDialerFactory {
	mock := &mockDialerFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
