// Code generated by mockery v2.50.0. DO NOT EDIT.

package service

import (
	context "context"

	payment "github.com/fyerfyer/trade-dependency/dto/payment"
	mock "github.com/stretchr/testify/mock"
)

// PaymentPort is an autogenerated mock type for the PaymentPort type
type PaymentPort struct {
	mock.Mock
}

type PaymentPort_Expecter struct {
	mock *mock.Mock
}

func (_m *PaymentPort) EXPECT() *PaymentPort_Expecter {
	return &PaymentPort_Expecter{mock: &_m.Mock}
}

// Charge provides a mock function with given fields: ctx, req
func (_m *PaymentPort) Charge(ctx context.Context, req *payment.ChargeRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for Charge")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.ChargeRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PaymentPort_Charge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Charge'
type PaymentPort_Charge_Call struct {
	*mock.Call
}

// Charge is a helper method to define mock.On call
//   - ctx context.Context
//   - req *payment.ChargeRequest
func (_e *PaymentPort_Expecter) Charge(ctx interface{}, req interface{}) *PaymentPort_Charge_Call {
	return &PaymentPort_Charge_Call{Call: _e.mock.On("Charge", ctx, req)}
}

func (_c *PaymentPort_Charge_Call) Run(run func(ctx context.Context, req *payment.ChargeRequest)) *PaymentPort_Charge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*payment.ChargeRequest))
	})
	return _c
}

func (_c *PaymentPort_Charge_Call) Return(_a0 error) *PaymentPort_Charge_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PaymentPort_Charge_Call) RunAndReturn(run func(context.Context, *payment.ChargeRequest) error) *PaymentPort_Charge_Call {
	_c.Call.Return(run)
	return _c
}

// GetPayment provides a mock function with given fields: ctx, req
func (_m *PaymentPort) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetPayment")
	}

	var r0 *payment.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentRequest) *payment.GetPaymentResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.GetPaymentRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PaymentPort_GetPayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPayment'
type PaymentPort_GetPayment_Call struct {
	*mock.Call
}

// GetPayment is a helper method to define mock.On call
//   - ctx context.Context
//   - req *payment.GetPaymentRequest
func (_e *PaymentPort_Expecter) GetPayment(ctx interface{}, req interface{}) *PaymentPort_GetPayment_Call {
	return &PaymentPort_GetPayment_Call{Call: _e.mock.On("GetPayment", ctx, req)}
}

func (_c *PaymentPort_GetPayment_Call) Run(run func(ctx context.Context, req *payment.GetPaymentRequest)) *PaymentPort_GetPayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*payment.GetPaymentRequest))
	})
	return _c
}

func (_c *PaymentPort_GetPayment_Call) Return(_a0 *payment.GetPaymentResponse, _a1 error) *PaymentPort_GetPayment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PaymentPort_GetPayment_Call) RunAndReturn(run func(context.Context, *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error)) *PaymentPort_GetPayment_Call {
	_c.Call.Return(run)
	return _c
}

// NewPaymentPort creates a new instance of PaymentPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentPort {
	mock := &PaymentPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
