package grpc

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/customer"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/customer"
)

func (a *Adapter) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CreateCustomerResponse, error) {
	res, err := a.service.CreateCustomer(ctx, &customer.CreateCustomerRequest{
		CustomerName: req.GetCustomerName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateCustomerResponse{
		CustomerId: res.CustomerID,
		Success:    true,
	}, nil
}

func (a *Adapter) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerResponse, error) {
	res, err := a.service.GetCustomer(ctx, &customer.GetCustomerRequest{
		CustomerName: req.GetCustomerName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetCustomerResponse{Customer: &pb.CustomerEntity{
		CustomerId: res.Customer.ID,
		Name:       res.Customer.Name,
		Status:     res.Customer.Status,
		Balance:    res.Customer.Balance,
		CreateAt:   uint64(res.Customer.CreatedAt.Unix()),
	}}, nil
}

func (a *Adapter) DeactiveCustomer(ctx context.Context, req *pb.DeactivateCustomerRequest) (*pb.DeactivateCustomerResponse, error) {
	err := a.service.DeactivateCustomer(ctx, &customer.DeactivateCustomerRequest{
		CustomerName: req.GetCustomerName(),
	})

	if err != nil {
		return &pb.DeactivateCustomerResponse{Success: false}, err
	}

	return &pb.DeactivateCustomerResponse{Success: true}, nil
}

func (a *Adapter) ActivateCustomer(ctx context.Context, req *pb.ActivateCustomerRequest) (*pb.ActivateCustomerResponse, error) {
	err := a.service.ActivateCustomer(ctx, &customer.ActivateCustomerRequest{
		CustomerName: req.GetCustomerName(),
	})
	if err != nil {
		return &pb.ActivateCustomerResponse{Success: false}, err
	}

	return &pb.ActivateCustomerResponse{Success: true}, nil
}

func (a *Adapter) SubmitOrder(ctx context.Context, req *pb.SubmitOrderRequest) (*pb.SubmitOrderResponse, error) {
	var items []*customer.OrderItem
	for _, item := range req.GetOrderItems() {
		items = append(items, &customer.OrderItem{
			ProductCode: item.GetProductCode(),
			UnitPrice:   item.GetUnitPrice(),
			Quantity:    item.GetQuantity(),
		})
	}

	err := a.service.SubmitOrder(ctx, &customer.SubmitOrderRequest{
		CustomerName: req.CustomerName,
		OrderItems:   items,
	})

	if err != nil {
		return &pb.SubmitOrderResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &pb.SubmitOrderResponse{
		Success: true,
		Message: "items processed successfully",
	}, nil
}

func (a *Adapter) PayOrder(ctx context.Context, req *pb.PayOrderRequest) (*pb.PayOrderResponse, error) {
	err := a.service.PayOrder(ctx, &customer.PayOrderRequest{
		CustomerName: req.GetCustomerName(),
		OrderID:      req.GetOrderId(),
	})
	if err != nil {
		return &pb.PayOrderResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &pb.PayOrderResponse{
		Success: true,
		Message: "order processed successfully",
	}, nil
}

func (a *Adapter) GetUnpaidOrders(ctx context.Context, req *pb.GetUnpaidOrdersRequest) (*pb.GetUnpaidOrdersResponse, error) {
	res, err := a.service.GetUnpaidOrders(ctx, &customer.GetUnpaidOrdersRequest{
		CustomerID: req.GetCustomerId(),
	})
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, o := range res.UnpaidOrders {
		var items []*pb.OrderItem
		for _, item := range o.Items {
			items = append(items, &pb.OrderItem{
				ProductCode: item.ProductCode,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			})
		}

		pbOrders = append(pbOrders, &pb.Order{
			OrderId: o.OrderID,
			Items:   items,
			Status:  o.Status,
		})
	}
	return &pb.GetUnpaidOrdersResponse{
		UnpaidOrders: pbOrders,
	}, nil
}
