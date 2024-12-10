package order

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/order"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderAdapter struct {
	order pb.OrderClient
}

func NewOrderAdapter(orderserviceURL string) (*OrderAdapter, error) {
	conn, err := grpc.Dial(orderserviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &OrderAdapter{
		order: pb.NewOrderClient(conn),
	}, nil
}

func (a *OrderAdapter) ProcessItems(ctx context.Context, req *order.ProcessItemsRequest) error {
	var items []*pb.OrderItem
	for _, item := range req.OrderItems {
		items = append(items, &pb.OrderItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	_, err := a.order.ProcessItems(ctx, &pb.ProcessItemsRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: req.Customer.CustomerID,
			Balance:    req.Customer.Balance,
		},

		OrderItems: items,
	})

	return err
}

func (a *OrderAdapter) ProcessOrder(ctx context.Context, req *order.ProcessOrderRequest) error {
	_, err := a.order.ProcessOrder(ctx, &pb.ProcessOrderRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: req.Customer.CustomerID,
			Balance:    req.Customer.Balance,
		},
		OrderId: req.OrderID,
	})

	return err
}

func (a *OrderAdapter) GetUnpaidOrder(ctx context.Context, req *order.GetUnpaidOrdersRequest) (*order.GetUnpaidOrdersResponse, error) {
	res, err := a.order.GetUnpaidOrders(ctx, &pb.GetUnpaidOrdersRequest{
		CustomerId: req.CustomerID,
	})

	if err != nil {
		return nil, err
	}

	var orders []order.OrderEntity
	for _, o := range res.Orders {
		var items []*order.OrderItemDTO
		for _, item := range o.OrderItems {
			items = append(items, &order.OrderItemDTO{
				ProductCode: item.ProductCode,
				UnitPrice:   item.UnitPrice,
				Quantity:    item.Quantity,
			})
		}
		orders = append(orders, order.OrderEntity{
			OrderID: o.OrderId,
			Items:   items,
			Status:  o.Status,
		})
	}

	return &order.GetUnpaidOrdersResponse{
		Orders: orders,
	}, nil
}
