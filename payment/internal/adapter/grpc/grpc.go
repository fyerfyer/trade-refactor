package grpc

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/payment"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/payment"
)

func (a *Adapter) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	customerEntity := payment.CustomerEntity{
		CustomerID: req.GetCustomer().GetCustomerId(),
		Balance:    req.GetCustomer().GetBalance(),
	}
	orderEntity := payment.OrderEntity{
		OrderID:    req.GetOrder().GetOrderId(),
		TotalPrice: req.GetOrder().GetTotalPrice(),
	}
	err := a.service.Charge(ctx, &payment.ChargeRequest{
		Customer: customerEntity,
		Order:    orderEntity,
	})

	if err != nil {
		return &pb.ChargeResponse{
			Status:  "failure",
			Message: err.Error(),
		}, err
	}

	return &pb.ChargeResponse{
		Status:  "success",
		Message: "successfully purchase the payment",
	}, nil
}
