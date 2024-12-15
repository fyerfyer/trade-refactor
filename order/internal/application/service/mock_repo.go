// Code generated by mockery v2.25.0. DO NOT EDIT.

package service

import (
	context "context"

	domain "github.com/fyerfyer/trade-refactor/order/internal/application/domain"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, orderID
func (_m *Repository) Delete(ctx context.Context, orderID uint64) error {
	ret := _m.Called(ctx, orderID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, orderID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type Repository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - orderID uint64
func (_e *Repository_Expecter) Delete(ctx interface{}, orderID interface{}) *Repository_Delete_Call {
	return &Repository_Delete_Call{Call: _e.mock.On("Delete", ctx, orderID)}
}

func (_c *Repository_Delete_Call) Run(run func(ctx context.Context, orderID uint64)) *Repository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *Repository_Delete_Call) Return(_a0 error) *Repository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Delete_Call) RunAndReturn(run func(context.Context, uint64) error) *Repository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, orderID
func (_m *Repository) Get(ctx context.Context, orderID uint64) (*domain.Order, error) {
	ret := _m.Called(ctx, orderID)

	var r0 *domain.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*domain.Order, error)); ok {
		return rf(ctx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *domain.Order); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type Repository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - orderID uint64
func (_e *Repository_Expecter) Get(ctx interface{}, orderID interface{}) *Repository_Get_Call {
	return &Repository_Get_Call{Call: _e.mock.On("Get", ctx, orderID)}
}

func (_c *Repository_Get_Call) Run(run func(ctx context.Context, orderID uint64)) *Repository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *Repository_Get_Call) Return(_a0 *domain.Order, _a1 error) *Repository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Get_Call) RunAndReturn(run func(context.Context, uint64) (*domain.Order, error)) *Repository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnpaidOrder provides a mock function with given fields: ctx, orderID
func (_m *Repository) GetUnpaidOrder(ctx context.Context, orderID uint64) (*domain.Order, error) {
	ret := _m.Called(ctx, orderID)

	var r0 *domain.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*domain.Order, error)); ok {
		return rf(ctx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *domain.Order); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetUnpaidOrder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnpaidOrder'
type Repository_GetUnpaidOrder_Call struct {
	*mock.Call
}

// GetUnpaidOrder is a helper method to define mock.On call
//   - ctx context.Context
//   - orderID uint64
func (_e *Repository_Expecter) GetUnpaidOrder(ctx interface{}, orderID interface{}) *Repository_GetUnpaidOrder_Call {
	return &Repository_GetUnpaidOrder_Call{Call: _e.mock.On("GetUnpaidOrder", ctx, orderID)}
}

func (_c *Repository_GetUnpaidOrder_Call) Run(run func(ctx context.Context, orderID uint64)) *Repository_GetUnpaidOrder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *Repository_GetUnpaidOrder_Call) Return(_a0 *domain.Order, _a1 error) *Repository_GetUnpaidOrder_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetUnpaidOrder_Call) RunAndReturn(run func(context.Context, uint64) (*domain.Order, error)) *Repository_GetUnpaidOrder_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnpaidOrders provides a mock function with given fields: ctx, customerID
func (_m *Repository) GetUnpaidOrders(ctx context.Context, customerID uint64) ([]domain.Order, error) {
	ret := _m.Called(ctx, customerID)

	var r0 []domain.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) ([]domain.Order, error)); ok {
		return rf(ctx, customerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) []domain.Order); ok {
		r0 = rf(ctx, customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetUnpaidOrders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnpaidOrders'
type Repository_GetUnpaidOrders_Call struct {
	*mock.Call
}

// GetUnpaidOrders is a helper method to define mock.On call
//   - ctx context.Context
//   - customerID uint64
func (_e *Repository_Expecter) GetUnpaidOrders(ctx interface{}, customerID interface{}) *Repository_GetUnpaidOrders_Call {
	return &Repository_GetUnpaidOrders_Call{Call: _e.mock.On("GetUnpaidOrders", ctx, customerID)}
}

func (_c *Repository_GetUnpaidOrders_Call) Run(run func(ctx context.Context, customerID uint64)) *Repository_GetUnpaidOrders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *Repository_GetUnpaidOrders_Call) Return(_a0 []domain.Order, _a1 error) *Repository_GetUnpaidOrders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetUnpaidOrders_Call) RunAndReturn(run func(context.Context, uint64) ([]domain.Order, error)) *Repository_GetUnpaidOrders_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: ctx, order
func (_m *Repository) Save(ctx context.Context, order *domain.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type Repository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - order *domain.Order
func (_e *Repository_Expecter) Save(ctx interface{}, order interface{}) *Repository_Save_Call {
	return &Repository_Save_Call{Call: _e.mock.On("Save", ctx, order)}
}

func (_c *Repository_Save_Call) Run(run func(ctx context.Context, order *domain.Order)) *Repository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Order))
	})
	return _c
}

func (_c *Repository_Save_Call) Return(_a0 error) *Repository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Save_Call) RunAndReturn(run func(context.Context, *domain.Order) error) *Repository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, order
func (_m *Repository) Update(ctx context.Context, order *domain.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type Repository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - order *domain.Order
func (_e *Repository_Expecter) Update(ctx interface{}, order interface{}) *Repository_Update_Call {
	return &Repository_Update_Call{Call: _e.mock.On("Update", ctx, order)}
}

func (_c *Repository_Update_Call) Run(run func(ctx context.Context, order *domain.Order)) *Repository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Order))
	})
	return _c
}

func (_c *Repository_Update_Call) Return(_a0 error) *Repository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Update_Call) RunAndReturn(run func(context.Context, *domain.Order) error) *Repository_Update_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
