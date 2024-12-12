package port

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/order"
)

type OrderPort interface {
	ProcessItems(ctx context.Context, req *order.ProcessItemsRequest) error
	ProcessOrder(ctx context.Context, req *order.ProcessOrderRequest) error
	GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error)
	GetUnpaidOrders(ctx context.Context, req *order.GetUnpaidOrdersRequest) (*order.GetUnpaidOrdersResponse, error)
}
