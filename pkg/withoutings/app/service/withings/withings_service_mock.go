// Code generated by mockery v2.20.0. DO NOT EDIT.

package withings

import (
	context "context"

	account "github.com/roessland/withoutings/pkg/withoutings/domain/account"

	domainwithings "github.com/roessland/withoutings/pkg/withoutings/domain/withings"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

type MockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockService) EXPECT() *MockService_Expecter {
	return &MockService_Expecter{mock: &_m.Mock}
}

// NotifyList provides a mock function with given fields: ctx, acc, params
func (_m *MockService) NotifyList(ctx context.Context, acc *account.Account, params domainwithings.NotifyListParams) (*domainwithings.NotifyListResponse, error) {
	ret := _m.Called(ctx, acc, params)

	var r0 *domainwithings.NotifyListResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.NotifyListParams) (*domainwithings.NotifyListResponse, error)); ok {
		return rf(ctx, acc, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.NotifyListParams) *domainwithings.NotifyListResponse); ok {
		r0 = rf(ctx, acc, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domainwithings.NotifyListResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *account.Account, domainwithings.NotifyListParams) error); ok {
		r1 = rf(ctx, acc, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_NotifyList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NotifyList'
type MockService_NotifyList_Call struct {
	*mock.Call
}

// NotifyList is a helper method to define mock.On call
//   - ctx context.Context
//   - acc *account.Account
//   - params domainwithings.NotifyListParams
func (_e *MockService_Expecter) NotifyList(ctx interface{}, acc interface{}, params interface{}) *MockService_NotifyList_Call {
	return &MockService_NotifyList_Call{Call: _e.mock.On("NotifyList", ctx, acc, params)}
}

func (_c *MockService_NotifyList_Call) Run(run func(ctx context.Context, acc *account.Account, params domainwithings.NotifyListParams)) *MockService_NotifyList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*account.Account), args[2].(domainwithings.NotifyListParams))
	})
	return _c
}

func (_c *MockService_NotifyList_Call) Return(_a0 *domainwithings.NotifyListResponse, _a1 error) *MockService_NotifyList_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_NotifyList_Call) RunAndReturn(run func(context.Context, *account.Account, domainwithings.NotifyListParams) (*domainwithings.NotifyListResponse, error)) *MockService_NotifyList_Call {
	_c.Call.Return(run)
	return _c
}

// NotifySubscribe provides a mock function with given fields: ctx, acc, params
func (_m *MockService) NotifySubscribe(ctx context.Context, acc *account.Account, params domainwithings.NotifySubscribeParams) (*domainwithings.NotifySubscribeResponse, error) {
	ret := _m.Called(ctx, acc, params)

	var r0 *domainwithings.NotifySubscribeResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.NotifySubscribeParams) (*domainwithings.NotifySubscribeResponse, error)); ok {
		return rf(ctx, acc, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.NotifySubscribeParams) *domainwithings.NotifySubscribeResponse); ok {
		r0 = rf(ctx, acc, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domainwithings.NotifySubscribeResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *account.Account, domainwithings.NotifySubscribeParams) error); ok {
		r1 = rf(ctx, acc, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_NotifySubscribe_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NotifySubscribe'
type MockService_NotifySubscribe_Call struct {
	*mock.Call
}

// NotifySubscribe is a helper method to define mock.On call
//   - ctx context.Context
//   - acc *account.Account
//   - params domainwithings.NotifySubscribeParams
func (_e *MockService_Expecter) NotifySubscribe(ctx interface{}, acc interface{}, params interface{}) *MockService_NotifySubscribe_Call {
	return &MockService_NotifySubscribe_Call{Call: _e.mock.On("NotifySubscribe", ctx, acc, params)}
}

func (_c *MockService_NotifySubscribe_Call) Run(run func(ctx context.Context, acc *account.Account, params domainwithings.NotifySubscribeParams)) *MockService_NotifySubscribe_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*account.Account), args[2].(domainwithings.NotifySubscribeParams))
	})
	return _c
}

func (_c *MockService_NotifySubscribe_Call) Return(_a0 *domainwithings.NotifySubscribeResponse, _a1 error) *MockService_NotifySubscribe_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_NotifySubscribe_Call) RunAndReturn(run func(context.Context, *account.Account, domainwithings.NotifySubscribeParams) (*domainwithings.NotifySubscribeResponse, error)) *MockService_NotifySubscribe_Call {
	_c.Call.Return(run)
	return _c
}

// SleepGetsummary provides a mock function with given fields: ctx, acc, params
func (_m *MockService) SleepGetsummary(ctx context.Context, acc *account.Account, params domainwithings.SleepGetsummaryParams) (*domainwithings.SleepGetsummaryResponse, error) {
	ret := _m.Called(ctx, acc, params)

	var r0 *domainwithings.SleepGetsummaryResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.SleepGetsummaryParams) (*domainwithings.SleepGetsummaryResponse, error)); ok {
		return rf(ctx, acc, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *account.Account, domainwithings.SleepGetsummaryParams) *domainwithings.SleepGetsummaryResponse); ok {
		r0 = rf(ctx, acc, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domainwithings.SleepGetsummaryResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *account.Account, domainwithings.SleepGetsummaryParams) error); ok {
		r1 = rf(ctx, acc, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_SleepGetsummary_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SleepGetsummary'
type MockService_SleepGetsummary_Call struct {
	*mock.Call
}

// SleepGetsummary is a helper method to define mock.On call
//   - ctx context.Context
//   - acc *account.Account
//   - params domainwithings.SleepGetsummaryParams
func (_e *MockService_Expecter) SleepGetsummary(ctx interface{}, acc interface{}, params interface{}) *MockService_SleepGetsummary_Call {
	return &MockService_SleepGetsummary_Call{Call: _e.mock.On("SleepGetsummary", ctx, acc, params)}
}

func (_c *MockService_SleepGetsummary_Call) Run(run func(ctx context.Context, acc *account.Account, params domainwithings.SleepGetsummaryParams)) *MockService_SleepGetsummary_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*account.Account), args[2].(domainwithings.SleepGetsummaryParams))
	})
	return _c
}

func (_c *MockService_SleepGetsummary_Call) Return(_a0 *domainwithings.SleepGetsummaryResponse, _a1 error) *MockService_SleepGetsummary_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_SleepGetsummary_Call) RunAndReturn(run func(context.Context, *account.Account, domainwithings.SleepGetsummaryParams) (*domainwithings.SleepGetsummaryResponse, error)) *MockService_SleepGetsummary_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockService interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockService(t mockConstructorTestingTNewMockService) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}