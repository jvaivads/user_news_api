// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	mail "gopkg.in/mail.v2"
)

// Dialer is an autogenerated mock type for the Dialer type
type Dialer struct {
	mock.Mock
}

// DialAndSend provides a mock function with given fields: m
func (_m *Dialer) DialAndSend(m ...*mail.Message) error {
	_va := make([]interface{}, len(m))
	for _i := range m {
		_va[_i] = m[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*mail.Message) error); ok {
		r0 = rf(m...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDialer creates a new instance of Dialer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDialer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Dialer {
	mock := &Dialer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
