// Code generated by mockery v2.20.0. DO NOT EDIT.

package dogu_config

import (
	core "github.com/cloudogu/cesapp-lib/core"
	mock "github.com/stretchr/testify/mock"
)

// mockDoguConfigurationEditor is an autogenerated mock type for the doguConfigurationEditor type
type mockDoguConfigurationEditor struct {
	mock.Mock
}

type mockDoguConfigurationEditor_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguConfigurationEditor) EXPECT() *mockDoguConfigurationEditor_Expecter {
	return &mockDoguConfigurationEditor_Expecter{mock: &_m.Mock}
}

// DeleteField provides a mock function with given fields: field
func (_m *mockDoguConfigurationEditor) DeleteField(field core.ConfigurationField) error {
	ret := _m.Called(field)

	var r0 error
	if rf, ok := ret.Get(0).(func(core.ConfigurationField) error); ok {
		r0 = rf(field)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigurationEditor_DeleteField_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteField'
type mockDoguConfigurationEditor_DeleteField_Call struct {
	*mock.Call
}

// DeleteField is a helper method to define mock.On call
//  - field core.ConfigurationField
func (_e *mockDoguConfigurationEditor_Expecter) DeleteField(field interface{}) *mockDoguConfigurationEditor_DeleteField_Call {
	return &mockDoguConfigurationEditor_DeleteField_Call{Call: _e.mock.On("DeleteField", field)}
}

func (_c *mockDoguConfigurationEditor_DeleteField_Call) Run(run func(field core.ConfigurationField)) *mockDoguConfigurationEditor_DeleteField_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(core.ConfigurationField))
	})
	return _c
}

func (_c *mockDoguConfigurationEditor_DeleteField_Call) Return(_a0 error) *mockDoguConfigurationEditor_DeleteField_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigurationEditor_DeleteField_Call) RunAndReturn(run func(core.ConfigurationField) error) *mockDoguConfigurationEditor_DeleteField_Call {
	_c.Call.Return(run)
	return _c
}

// EditConfiguration provides a mock function with given fields: fields, deleteOnEmpty
func (_m *mockDoguConfigurationEditor) EditConfiguration(fields []core.ConfigurationField, deleteOnEmpty bool) error {
	ret := _m.Called(fields, deleteOnEmpty)

	var r0 error
	if rf, ok := ret.Get(0).(func([]core.ConfigurationField, bool) error); ok {
		r0 = rf(fields, deleteOnEmpty)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigurationEditor_EditConfiguration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EditConfiguration'
type mockDoguConfigurationEditor_EditConfiguration_Call struct {
	*mock.Call
}

// EditConfiguration is a helper method to define mock.On call
//  - fields []core.ConfigurationField
//  - deleteOnEmpty bool
func (_e *mockDoguConfigurationEditor_Expecter) EditConfiguration(fields interface{}, deleteOnEmpty interface{}) *mockDoguConfigurationEditor_EditConfiguration_Call {
	return &mockDoguConfigurationEditor_EditConfiguration_Call{Call: _e.mock.On("EditConfiguration", fields, deleteOnEmpty)}
}

func (_c *mockDoguConfigurationEditor_EditConfiguration_Call) Run(run func(fields []core.ConfigurationField, deleteOnEmpty bool)) *mockDoguConfigurationEditor_EditConfiguration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]core.ConfigurationField), args[1].(bool))
	})
	return _c
}

func (_c *mockDoguConfigurationEditor_EditConfiguration_Call) Return(_a0 error) *mockDoguConfigurationEditor_EditConfiguration_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigurationEditor_EditConfiguration_Call) RunAndReturn(run func([]core.ConfigurationField, bool) error) *mockDoguConfigurationEditor_EditConfiguration_Call {
	_c.Call.Return(run)
	return _c
}

// GetCurrentValue provides a mock function with given fields: field
func (_m *mockDoguConfigurationEditor) GetCurrentValue(field core.ConfigurationField) (string, error) {
	ret := _m.Called(field)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(core.ConfigurationField) (string, error)); ok {
		return rf(field)
	}
	if rf, ok := ret.Get(0).(func(core.ConfigurationField) string); ok {
		r0 = rf(field)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(core.ConfigurationField) error); ok {
		r1 = rf(field)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigurationEditor_GetCurrentValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCurrentValue'
type mockDoguConfigurationEditor_GetCurrentValue_Call struct {
	*mock.Call
}

// GetCurrentValue is a helper method to define mock.On call
//  - field core.ConfigurationField
func (_e *mockDoguConfigurationEditor_Expecter) GetCurrentValue(field interface{}) *mockDoguConfigurationEditor_GetCurrentValue_Call {
	return &mockDoguConfigurationEditor_GetCurrentValue_Call{Call: _e.mock.On("GetCurrentValue", field)}
}

func (_c *mockDoguConfigurationEditor_GetCurrentValue_Call) Run(run func(field core.ConfigurationField)) *mockDoguConfigurationEditor_GetCurrentValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(core.ConfigurationField))
	})
	return _c
}

func (_c *mockDoguConfigurationEditor_GetCurrentValue_Call) Return(_a0 string, _a1 error) *mockDoguConfigurationEditor_GetCurrentValue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigurationEditor_GetCurrentValue_Call) RunAndReturn(run func(core.ConfigurationField) (string, error)) *mockDoguConfigurationEditor_GetCurrentValue_Call {
	_c.Call.Return(run)
	return _c
}

// SetFieldToValue provides a mock function with given fields: field, value, deleteOnEmpty
func (_m *mockDoguConfigurationEditor) SetFieldToValue(field core.ConfigurationField, value string, deleteOnEmpty bool) error {
	ret := _m.Called(field, value, deleteOnEmpty)

	var r0 error
	if rf, ok := ret.Get(0).(func(core.ConfigurationField, string, bool) error); ok {
		r0 = rf(field, value, deleteOnEmpty)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigurationEditor_SetFieldToValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFieldToValue'
type mockDoguConfigurationEditor_SetFieldToValue_Call struct {
	*mock.Call
}

// SetFieldToValue is a helper method to define mock.On call
//  - field core.ConfigurationField
//  - value string
//  - deleteOnEmpty bool
func (_e *mockDoguConfigurationEditor_Expecter) SetFieldToValue(field interface{}, value interface{}, deleteOnEmpty interface{}) *mockDoguConfigurationEditor_SetFieldToValue_Call {
	return &mockDoguConfigurationEditor_SetFieldToValue_Call{Call: _e.mock.On("SetFieldToValue", field, value, deleteOnEmpty)}
}

func (_c *mockDoguConfigurationEditor_SetFieldToValue_Call) Run(run func(field core.ConfigurationField, value string, deleteOnEmpty bool)) *mockDoguConfigurationEditor_SetFieldToValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(core.ConfigurationField), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *mockDoguConfigurationEditor_SetFieldToValue_Call) Return(_a0 error) *mockDoguConfigurationEditor_SetFieldToValue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigurationEditor_SetFieldToValue_Call) RunAndReturn(run func(core.ConfigurationField, string, bool) error) *mockDoguConfigurationEditor_SetFieldToValue_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguConfigurationEditor interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguConfigurationEditor creates a new instance of mockDoguConfigurationEditor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguConfigurationEditor(t mockConstructorTestingTnewMockDoguConfigurationEditor) *mockDoguConfigurationEditor {
	mock := &mockDoguConfigurationEditor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
