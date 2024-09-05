// Code generated by mockery v2.45.0. DO NOT EDIT.

package newsletters

import (
	context "context"
	model "mailchump/pkg/model"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MockNewsletterStore is an autogenerated mock type for the NewsletterStore type
type MockNewsletterStore struct {
	mock.Mock
}

type MockNewsletterStore_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNewsletterStore) EXPECT() *MockNewsletterStore_Expecter {
	return &MockNewsletterStore_Expecter{mock: &_m.Mock}
}

// DeleteNewsletter provides a mock function with given fields: ctx, id
func (_m *MockNewsletterStore) DeleteNewsletter(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteNewsletter")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockNewsletterStore_DeleteNewsletter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteNewsletter'
type MockNewsletterStore_DeleteNewsletter_Call struct {
	*mock.Call
}

// DeleteNewsletter is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockNewsletterStore_Expecter) DeleteNewsletter(ctx interface{}, id interface{}) *MockNewsletterStore_DeleteNewsletter_Call {
	return &MockNewsletterStore_DeleteNewsletter_Call{Call: _e.mock.On("DeleteNewsletter", ctx, id)}
}

func (_c *MockNewsletterStore_DeleteNewsletter_Call) Run(run func(ctx context.Context, id string)) *MockNewsletterStore_DeleteNewsletter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockNewsletterStore_DeleteNewsletter_Call) Return(_a0 error) *MockNewsletterStore_DeleteNewsletter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNewsletterStore_DeleteNewsletter_Call) RunAndReturn(run func(context.Context, string) error) *MockNewsletterStore_DeleteNewsletter_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllNewsletters provides a mock function with given fields: ctx
func (_m *MockNewsletterStore) GetAllNewsletters(ctx context.Context) (model.Newsletters, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllNewsletters")
	}

	var r0 model.Newsletters
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (model.Newsletters, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) model.Newsletters); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Newsletters)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNewsletterStore_GetAllNewsletters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllNewsletters'
type MockNewsletterStore_GetAllNewsletters_Call struct {
	*mock.Call
}

// GetAllNewsletters is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockNewsletterStore_Expecter) GetAllNewsletters(ctx interface{}) *MockNewsletterStore_GetAllNewsletters_Call {
	return &MockNewsletterStore_GetAllNewsletters_Call{Call: _e.mock.On("GetAllNewsletters", ctx)}
}

func (_c *MockNewsletterStore_GetAllNewsletters_Call) Run(run func(ctx context.Context)) *MockNewsletterStore_GetAllNewsletters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockNewsletterStore_GetAllNewsletters_Call) Return(_a0 model.Newsletters, _a1 error) *MockNewsletterStore_GetAllNewsletters_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNewsletterStore_GetAllNewsletters_Call) RunAndReturn(run func(context.Context) (model.Newsletters, error)) *MockNewsletterStore_GetAllNewsletters_Call {
	_c.Call.Return(run)
	return _c
}

// GetNewsletterById provides a mock function with given fields: ctx, id
func (_m *MockNewsletterStore) GetNewsletterById(ctx context.Context, id string) (model.Newsletter, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetNewsletterById")
	}

	var r0 model.Newsletter
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.Newsletter, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.Newsletter); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Newsletter)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNewsletterStore_GetNewsletterById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNewsletterById'
type MockNewsletterStore_GetNewsletterById_Call struct {
	*mock.Call
}

// GetNewsletterById is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockNewsletterStore_Expecter) GetNewsletterById(ctx interface{}, id interface{}) *MockNewsletterStore_GetNewsletterById_Call {
	return &MockNewsletterStore_GetNewsletterById_Call{Call: _e.mock.On("GetNewsletterById", ctx, id)}
}

func (_c *MockNewsletterStore_GetNewsletterById_Call) Run(run func(ctx context.Context, id string)) *MockNewsletterStore_GetNewsletterById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockNewsletterStore_GetNewsletterById_Call) Return(_a0 model.Newsletter, _a1 error) *MockNewsletterStore_GetNewsletterById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNewsletterStore_GetNewsletterById_Call) RunAndReturn(run func(context.Context, string) (model.Newsletter, error)) *MockNewsletterStore_GetNewsletterById_Call {
	_c.Call.Return(run)
	return _c
}

// GetNewsletterOwnerID provides a mock function with given fields: ctx, id
func (_m *MockNewsletterStore) GetNewsletterOwnerID(ctx context.Context, id string) (uuid.UUID, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetNewsletterOwnerID")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (uuid.UUID, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) uuid.UUID); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNewsletterStore_GetNewsletterOwnerID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNewsletterOwnerID'
type MockNewsletterStore_GetNewsletterOwnerID_Call struct {
	*mock.Call
}

// GetNewsletterOwnerID is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockNewsletterStore_Expecter) GetNewsletterOwnerID(ctx interface{}, id interface{}) *MockNewsletterStore_GetNewsletterOwnerID_Call {
	return &MockNewsletterStore_GetNewsletterOwnerID_Call{Call: _e.mock.On("GetNewsletterOwnerID", ctx, id)}
}

func (_c *MockNewsletterStore_GetNewsletterOwnerID_Call) Run(run func(ctx context.Context, id string)) *MockNewsletterStore_GetNewsletterOwnerID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockNewsletterStore_GetNewsletterOwnerID_Call) Return(_a0 uuid.UUID, _a1 error) *MockNewsletterStore_GetNewsletterOwnerID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNewsletterStore_GetNewsletterOwnerID_Call) RunAndReturn(run func(context.Context, string) (uuid.UUID, error)) *MockNewsletterStore_GetNewsletterOwnerID_Call {
	_c.Call.Return(run)
	return _c
}

// HideNewsletter provides a mock function with given fields: ctx, id, owner
func (_m *MockNewsletterStore) HideNewsletter(ctx context.Context, id string, owner string) (bool, error) {
	ret := _m.Called(ctx, id, owner)

	if len(ret) == 0 {
		panic("no return value specified for HideNewsletter")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, id, owner)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, id, owner)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNewsletterStore_HideNewsletter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HideNewsletter'
type MockNewsletterStore_HideNewsletter_Call struct {
	*mock.Call
}

// HideNewsletter is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - owner string
func (_e *MockNewsletterStore_Expecter) HideNewsletter(ctx interface{}, id interface{}, owner interface{}) *MockNewsletterStore_HideNewsletter_Call {
	return &MockNewsletterStore_HideNewsletter_Call{Call: _e.mock.On("HideNewsletter", ctx, id, owner)}
}

func (_c *MockNewsletterStore_HideNewsletter_Call) Run(run func(ctx context.Context, id string, owner string)) *MockNewsletterStore_HideNewsletter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockNewsletterStore_HideNewsletter_Call) Return(isHidden bool, err error) *MockNewsletterStore_HideNewsletter_Call {
	_c.Call.Return(isHidden, err)
	return _c
}

func (_c *MockNewsletterStore_HideNewsletter_Call) RunAndReturn(run func(context.Context, string, string) (bool, error)) *MockNewsletterStore_HideNewsletter_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNewsletterStore creates a new instance of MockNewsletterStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNewsletterStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNewsletterStore {
	mock := &MockNewsletterStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
