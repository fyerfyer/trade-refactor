package grpc

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/order"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
)

func (a *Adapter) ProcessItems(ctx context.Context, req *pb.ProcessItemsRequest) (*pb.ProcessItemsResponse, error) {
	// convert grpc object into domain object
	var items []*order.OrderItemDTO
	for _, item := range req.OrderItems {
		items = append(items, &order.OrderItemDTO{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	err := a.service.ProcessItems(ctx, &order.ProcessItemsRequest{
		Customer: order.CustomerEntity{
			CustomerID: req.GetCustomer().GetCustomerId(),
			Balance:    req.GetCustomer().GetBalance(),
		},
		OrderItems: items,
	})
	if err != nil {
		return &pb.ProcessItemsResponse{Message: err.Error()}, err
	}

	return &pb.ProcessItemsResponse{Message: "successfully process items"}, nil
}

func (a *Adapter) ProcessOrder(ctx context.Context, req *pb.ProcessOrderRequest) (*pb.ProcessOrderResponse, error) {
	err := a.service.ProcessOrder(ctx, &order.ProcessOrderRequest{
		Customer: order.CustomerEntity{
			CustomerID: req.GetCustomer().GetCustomerId(),
			Balance:    req.GetCustomer().GetBalance(),
		},
		OrderID: req.GetOrderId(),
	})

	if err != nil {
		return &pb.ProcessOrderResponse{Message: err.Error()}, err
	}

	return &pb.ProcessOrderResponse{Message: "successfully process order"}, nil
}
