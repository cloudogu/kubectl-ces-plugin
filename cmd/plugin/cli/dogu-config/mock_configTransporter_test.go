// Code generated by mockery v2.20.0. DO NOT EDIT.

package dogu_config

import mock "github.com/stretchr/testify/mock"

// mockConfigTransporter is an autogenerated mock type for the configTransporter type
type mockConfigTransporter struct {
	mock.Mock
}

type mockConfigTransporter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockConfigTransporter) EXPECT() *mockConfigTransporter_Expecter {
	return &mockConfigTransporter_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: key
func (_m *mockConfigTransporter) Get(key string) interface{} {
	ret := _m.Called(key)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// mockConfigTransporter_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockConfigTransporter_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//  - key string
func (_e *mockConfigTransporter_Expecter) Get(key interface{}) *mockConfigTransporter_Get_Call {
	return &mockConfigTransporter_Get_Call{Call: _e.mock.On("Get", key)}
}

func (_c *mockConfigTransporter_Get_Call) Run(run func(key string)) *mockConfigTransporter_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockConfigTransporter_Get_Call) Return(_a0 interface{}) *mockConfigTransporter_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockConfigTransporter_Get_Call) RunAndReturn(run func(string) interface{}) *mockConfigTransporter_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: key, value
func (_m *mockConfigTransporter) Set(key string, value interface{}) {
	_m.Called(key, value)
}

// mockConfigTransporter_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type mockConfigTransporter_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//  - key string
//  - value interface{}
func (_e *mockConfigTransporter_Expecter) Set(key interface{}, value interface{}) *mockConfigTransporter_Set_Call {
	return &mockConfigTransporter_Set_Call{Call: _e.mock.On("Set", key, value)}
}

func (_c *mockConfigTransporter_Set_Call) Run(run func(key string, value interface{})) *mockConfigTransporter_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(interface{}))
	})
	return _c
}

func (_c *mockConfigTransporter_Set_Call) Return() *mockConfigTransporter_Set_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockConfigTransporter_Set_Call) RunAndReturn(run func(string, interface{})) *mockConfigTransporter_Set_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockConfigTransporter interface {
	mock.TestingT
	Cleanup(func())
}

// newMockConfigTransporter creates a new instance of mockConfigTransporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockConfigTransporter(t mockConstructorTestingTnewMockConfigTransporter) *mockConfigTransporter {
	mock := &mockConfigTransporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}