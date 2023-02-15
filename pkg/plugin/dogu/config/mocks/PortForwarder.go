// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PortForwarder is an autogenerated mock type for the PortForwarder type
type PortForwarder struct {
	mock.Mock
}

// ExecuteWithPortForward provides a mock function with given fields: fn
func (_m *PortForwarder) ExecuteWithPortForward(fn func() error) error {
	ret := _m.Called(fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(func() error) error); ok {
		r0 = rf(fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPortForwarder interface {
	mock.TestingT
	Cleanup(func())
}

// NewPortForwarder creates a new instance of PortForwarder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPortForwarder(t mockConstructorTestingTNewPortForwarder) *PortForwarder {
	mock := &PortForwarder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}