package port

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/customer"
)

type CustomerPort interface {
	CreateCustomer(ctx context.Context, req *customer.CreateCustomerRequest) (*customer.CreateCustomerResponse, error)
	GetCustomer(ctx context.Context, req *customer.GetCustomerRequest) (*customer.GetCustomerResponse, error)
	DeactivateCustomer(ctx context.Context, req *customer.DeactivateCustomerRequest) error
	ActicateCustomer(ctx context.Context, req *customer.ActivateCustomerRequest) error
	PayOrder(ctx context.Context, req *customer.PayOrderRequest) error
	SubmitOrder(ctx context.Context, req *customer.SubmitOrderRequest) error
	GetUnpaidOrders(ctx context.Context, req *customer.GetUnpaidOrdersRequest) (*customer.GetUnpaidOrdersResponse, error)
	StoreBalance(ctx context.Context, req *customer.StoreBalanceRequest) error
}
