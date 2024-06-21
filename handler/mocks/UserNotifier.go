// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserNotifier is an autogenerated mock type for the UserNotifier type
type UserNotifier struct {
	mock.Mock
}

// Notify provides a mock function with given fields: _a0, _a1, _a2
func (_m *UserNotifier) Notify(_a0 context.Context, _a1 string, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewUserNotifier creates a new instance of UserNotifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserNotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserNotifier {
	mock := &UserNotifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
