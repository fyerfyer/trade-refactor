package payment

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/payment"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentAdapter struct {
	payment pb.PaymentClient
}

func NewPaymentAdapter(paymentserviceURL string) (*PaymentAdapter, error) {
	conn, err := grpc.Dial(paymentserviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &PaymentAdapter{
		payment: pb.NewPaymentClient(conn),
	}, nil
}

func (a *PaymentAdapter) Charge(ctx context.Context, req *payment.ChargeRequest) error {
	_, err := a.payment.Charge(ctx, &pb.ChargeRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: req.Customer.CustomerID,
			Balance:    req.Customer.Balance,
		},

		Order: &pb.OrderEntity{
			OrderId:    req.Order.OrderID,
			TotalPrice: req.Order.TotalPrice,
		},
	})

	return err
}
