// Code generated by mockery v2.20.0. DO NOT EDIT.

package config

import mock "github.com/stretchr/testify/mock"

// mockDoguConfigService is an autogenerated mock type for the doguConfigService type
type mockDoguConfigService struct {
	mock.Mock
}

type mockDoguConfigService_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguConfigService) EXPECT() *mockDoguConfigService_Expecter {
	return &mockDoguConfigService_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: registryKey
func (_m *mockDoguConfigService) Delete(registryKey string) error {
	ret := _m.Called(registryKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(registryKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDoguConfigService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//  - registryKey string
func (_e *mockDoguConfigService_Expecter) Delete(registryKey interface{}) *mockDoguConfigService_Delete_Call {
	return &mockDoguConfigService_Delete_Call{Call: _e.mock.On("Delete", registryKey)}
}

func (_c *mockDoguConfigService_Delete_Call) Run(run func(registryKey string)) *mockDoguConfigService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguConfigService_Delete_Call) Return(_a0 error) *mockDoguConfigService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigService_Delete_Call) RunAndReturn(run func(string) error) *mockDoguConfigService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Edit provides a mock function with given fields: registryKey, registryValue, deleteOnEmpty
func (_m *mockDoguConfigService) Edit(registryKey string, registryValue string, deleteOnEmpty bool) error {
	ret := _m.Called(registryKey, registryValue, deleteOnEmpty)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool) error); ok {
		r0 = rf(registryKey, registryValue, deleteOnEmpty)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigService_Edit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Edit'
type mockDoguConfigService_Edit_Call struct {
	*mock.Call
}

// Edit is a helper method to define mock.On call
//  - registryKey string
//  - registryValue string
//  - deleteOnEmpty bool
func (_e *mockDoguConfigService_Expecter) Edit(registryKey interface{}, registryValue interface{}, deleteOnEmpty interface{}) *mockDoguConfigService_Edit_Call {
	return &mockDoguConfigService_Edit_Call{Call: _e.mock.On("Edit", registryKey, registryValue, deleteOnEmpty)}
}

func (_c *mockDoguConfigService_Edit_Call) Run(run func(registryKey string, registryValue string, deleteOnEmpty bool)) *mockDoguConfigService_Edit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *mockDoguConfigService_Edit_Call) Return(_a0 error) *mockDoguConfigService_Edit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigService_Edit_Call) RunAndReturn(run func(string, string, bool) error) *mockDoguConfigService_Edit_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllForDogu provides a mock function with given fields:
func (_m *mockDoguConfigService) GetAllForDogu() (map[string]string, error) {
	ret := _m.Called()

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigService_GetAllForDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllForDogu'
type mockDoguConfigService_GetAllForDogu_Call struct {
	*mock.Call
}

// GetAllForDogu is a helper method to define mock.On call
func (_e *mockDoguConfigService_Expecter) GetAllForDogu() *mockDoguConfigService_GetAllForDogu_Call {
	return &mockDoguConfigService_GetAllForDogu_Call{Call: _e.mock.On("GetAllForDogu")}
}

func (_c *mockDoguConfigService_GetAllForDogu_Call) Run(run func()) *mockDoguConfigService_GetAllForDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDoguConfigService_GetAllForDogu_Call) Return(_a0 map[string]string, _a1 error) *mockDoguConfigService_GetAllForDogu_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigService_GetAllForDogu_Call) RunAndReturn(run func() (map[string]string, error)) *mockDoguConfigService_GetAllForDogu_Call {
	_c.Call.Return(run)
	return _c
}

// GetValue provides a mock function with given fields: registryKey
func (_m *mockDoguConfigService) GetValue(registryKey string) (string, error) {
	ret := _m.Called(registryKey)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(registryKey)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(registryKey)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(registryKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigService_GetValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValue'
type mockDoguConfigService_GetValue_Call struct {
	*mock.Call
}

// GetValue is a helper method to define mock.On call
//  - registryKey string
func (_e *mockDoguConfigService_Expecter) GetValue(registryKey interface{}) *mockDoguConfigService_GetValue_Call {
	return &mockDoguConfigService_GetValue_Call{Call: _e.mock.On("GetValue", registryKey)}
}

func (_c *mockDoguConfigService_GetValue_Call) Run(run func(registryKey string)) *mockDoguConfigService_GetValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockDoguConfigService_GetValue_Call) Return(_a0 string, _a1 error) *mockDoguConfigService_GetValue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigService_GetValue_Call) RunAndReturn(run func(string) (string, error)) *mockDoguConfigService_GetValue_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguConfigService interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguConfigService creates a new instance of mockDoguConfigService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguConfigService(t mockConstructorTestingTnewMockDoguConfigService) *mockDoguConfigService {
	mock := &mockDoguConfigService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
