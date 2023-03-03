// Code generated by mockery v2.20.0. DO NOT EDIT.

package logger

import mock "github.com/stretchr/testify/mock"

// mockLogger is an autogenerated mock type for the logger type
type mockLogger struct {
	mock.Mock
}

type mockLogger_Expecter struct {
	mock *mock.Mock
}

func (_m *mockLogger) EXPECT() *mockLogger_Expecter {
	return &mockLogger_Expecter{mock: &_m.Mock}
}

// Debug provides a mock function with given fields: args
func (_m *mockLogger) Debug(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Debug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debug'
type mockLogger_Debug_Call struct {
	*mock.Call
}

// Debug is a helper method to define mock.On call
//  - args ...interface{}
func (_e *mockLogger_Expecter) Debug(args ...interface{}) *mockLogger_Debug_Call {
	return &mockLogger_Debug_Call{Call: _e.mock.On("Debug",
		append([]interface{}{}, args...)...)}
}

func (_c *mockLogger_Debug_Call) Run(run func(args ...interface{})) *mockLogger_Debug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Debug_Call) Return() *mockLogger_Debug_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Debug_Call) RunAndReturn(run func(...interface{})) *mockLogger_Debug_Call {
	_c.Call.Return(run)
	return _c
}

// Debugf provides a mock function with given fields: format, args
func (_m *mockLogger) Debugf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Debugf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debugf'
type mockLogger_Debugf_Call struct {
	*mock.Call
}

// Debugf is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *mockLogger_Expecter) Debugf(format interface{}, args ...interface{}) *mockLogger_Debugf_Call {
	return &mockLogger_Debugf_Call{Call: _e.mock.On("Debugf",
		append([]interface{}{format}, args...)...)}
}

func (_c *mockLogger_Debugf_Call) Run(run func(format string, args ...interface{})) *mockLogger_Debugf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Debugf_Call) Return() *mockLogger_Debugf_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Debugf_Call) RunAndReturn(run func(string, ...interface{})) *mockLogger_Debugf_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: args
func (_m *mockLogger) Error(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type mockLogger_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//  - args ...interface{}
func (_e *mockLogger_Expecter) Error(args ...interface{}) *mockLogger_Error_Call {
	return &mockLogger_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{}, args...)...)}
}

func (_c *mockLogger_Error_Call) Run(run func(args ...interface{})) *mockLogger_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Error_Call) Return() *mockLogger_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Error_Call) RunAndReturn(run func(...interface{})) *mockLogger_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Errorf provides a mock function with given fields: format, args
func (_m *mockLogger) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Errorf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Errorf'
type mockLogger_Errorf_Call struct {
	*mock.Call
}

// Errorf is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *mockLogger_Expecter) Errorf(format interface{}, args ...interface{}) *mockLogger_Errorf_Call {
	return &mockLogger_Errorf_Call{Call: _e.mock.On("Errorf",
		append([]interface{}{format}, args...)...)}
}

func (_c *mockLogger_Errorf_Call) Run(run func(format string, args ...interface{})) *mockLogger_Errorf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Errorf_Call) Return() *mockLogger_Errorf_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Errorf_Call) RunAndReturn(run func(string, ...interface{})) *mockLogger_Errorf_Call {
	_c.Call.Return(run)
	return _c
}

// Info provides a mock function with given fields: args
func (_m *mockLogger) Info(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type mockLogger_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//  - args ...interface{}
func (_e *mockLogger_Expecter) Info(args ...interface{}) *mockLogger_Info_Call {
	return &mockLogger_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{}, args...)...)}
}

func (_c *mockLogger_Info_Call) Run(run func(args ...interface{})) *mockLogger_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Info_Call) Return() *mockLogger_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Info_Call) RunAndReturn(run func(...interface{})) *mockLogger_Info_Call {
	_c.Call.Return(run)
	return _c
}

// Infof provides a mock function with given fields: format, args
func (_m *mockLogger) Infof(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Infof_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Infof'
type mockLogger_Infof_Call struct {
	*mock.Call
}

// Infof is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *mockLogger_Expecter) Infof(format interface{}, args ...interface{}) *mockLogger_Infof_Call {
	return &mockLogger_Infof_Call{Call: _e.mock.On("Infof",
		append([]interface{}{format}, args...)...)}
}

func (_c *mockLogger_Infof_Call) Run(run func(format string, args ...interface{})) *mockLogger_Infof_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Infof_Call) Return() *mockLogger_Infof_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Infof_Call) RunAndReturn(run func(string, ...interface{})) *mockLogger_Infof_Call {
	_c.Call.Return(run)
	return _c
}

// Warning provides a mock function with given fields: args
func (_m *mockLogger) Warning(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Warning_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warning'
type mockLogger_Warning_Call struct {
	*mock.Call
}

// Warning is a helper method to define mock.On call
//  - args ...interface{}
func (_e *mockLogger_Expecter) Warning(args ...interface{}) *mockLogger_Warning_Call {
	return &mockLogger_Warning_Call{Call: _e.mock.On("Warning",
		append([]interface{}{}, args...)...)}
}

func (_c *mockLogger_Warning_Call) Run(run func(args ...interface{})) *mockLogger_Warning_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Warning_Call) Return() *mockLogger_Warning_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Warning_Call) RunAndReturn(run func(...interface{})) *mockLogger_Warning_Call {
	_c.Call.Return(run)
	return _c
}

// Warningf provides a mock function with given fields: format, args
func (_m *mockLogger) Warningf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// mockLogger_Warningf_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Warningf'
type mockLogger_Warningf_Call struct {
	*mock.Call
}

// Warningf is a helper method to define mock.On call
//  - format string
//  - args ...interface{}
func (_e *mockLogger_Expecter) Warningf(format interface{}, args ...interface{}) *mockLogger_Warningf_Call {
	return &mockLogger_Warningf_Call{Call: _e.mock.On("Warningf",
		append([]interface{}{format}, args...)...)}
}

func (_c *mockLogger_Warningf_Call) Run(run func(format string, args ...interface{})) *mockLogger_Warningf_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogger_Warningf_Call) Return() *mockLogger_Warningf_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogger_Warningf_Call) RunAndReturn(run func(string, ...interface{})) *mockLogger_Warningf_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockLogger interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLogger creates a new instance of mockLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLogger(t mockConstructorTestingTnewMockLogger) *mockLogger {
	mock := &mockLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
