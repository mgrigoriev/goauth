// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	io "io"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// HTTPClient is an autogenerated mock type for the HTTPClient type
type HTTPClient struct {
	mock.Mock
}

// Post provides a mock function with given fields: url, contentType, body
func (_m *HTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	ret := _m.Called(url, contentType, body)

	if len(ret) == 0 {
		panic("no return value specified for Post")
	}

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string, io.Reader) (*http.Response, error)); ok {
		return rf(url, contentType, body)
	}
	if rf, ok := ret.Get(0).(func(string, string, io.Reader) *http.Response); ok {
		r0 = rf(url, contentType, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, io.Reader) error); ok {
		r1 = rf(url, contentType, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHTTPClient creates a new instance of HTTPClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHTTPClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *HTTPClient {
	mock := &HTTPClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
