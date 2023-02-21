// Code generated by mockery v2.20.0. DO NOT EDIT.

package keys

import (
	cesapp_libkeys "github.com/cloudogu/cesapp-lib/keys"
	mock "github.com/stretchr/testify/mock"
)

// mockKeyProvider is an autogenerated mock type for the keyProvider type
type mockKeyProvider struct {
	mock.Mock
}

type mockKeyProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *mockKeyProvider) EXPECT() *mockKeyProvider_Expecter {
	return &mockKeyProvider_Expecter{mock: &_m.Mock}
}

// ReadPublicKeyFromString provides a mock function with given fields: publicKeyString
func (_m *mockKeyProvider) ReadPublicKeyFromString(publicKeyString string) (*cesapp_libkeys.PublicKey, error) {
	ret := _m.Called(publicKeyString)

	var r0 *cesapp_libkeys.PublicKey
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*cesapp_libkeys.PublicKey, error)); ok {
		return rf(publicKeyString)
	}
	if rf, ok := ret.Get(0).(func(string) *cesapp_libkeys.PublicKey); ok {
		r0 = rf(publicKeyString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cesapp_libkeys.PublicKey)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(publicKeyString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockKeyProvider_ReadPublicKeyFromString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadPublicKeyFromString'
type mockKeyProvider_ReadPublicKeyFromString_Call struct {
	*mock.Call
}

// ReadPublicKeyFromString is a helper method to define mock.On call
//  - publicKeyString string
func (_e *mockKeyProvider_Expecter) ReadPublicKeyFromString(publicKeyString interface{}) *mockKeyProvider_ReadPublicKeyFromString_Call {
	return &mockKeyProvider_ReadPublicKeyFromString_Call{Call: _e.mock.On("ReadPublicKeyFromString", publicKeyString)}
}

func (_c *mockKeyProvider_ReadPublicKeyFromString_Call) Run(run func(publicKeyString string)) *mockKeyProvider_ReadPublicKeyFromString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockKeyProvider_ReadPublicKeyFromString_Call) Return(_a0 *cesapp_libkeys.PublicKey, _a1 error) *mockKeyProvider_ReadPublicKeyFromString_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockKeyProvider_ReadPublicKeyFromString_Call) RunAndReturn(run func(string) (*cesapp_libkeys.PublicKey, error)) *mockKeyProvider_ReadPublicKeyFromString_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockKeyProvider interface {
	mock.TestingT
	Cleanup(func())
}

// newMockKeyProvider creates a new instance of mockKeyProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockKeyProvider(t mockConstructorTestingTnewMockKeyProvider) *mockKeyProvider {
	mock := &mockKeyProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}