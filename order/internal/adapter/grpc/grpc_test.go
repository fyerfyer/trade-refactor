package grpc

import (
	"context"
	"testing"

	"github.com/fyerfyer/trade-dependency/dto/order"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/fyerfyer/trade-dependency/proto/grpc/order"
)

func TestProcessItems_Success(t *testing.T) {
	mockOrder := new(OrderPort)
	adapter := NewAdapter(mockOrder, 50052)

	req := &order.ProcessItemsRequest{
		Customer: order.CustomerEntity{
			CustomerID: 1,
			Balance:    200.0,
		},
		OrderItems: []*order.OrderItemDTO{
			{
				ProductCode: "Juice",
				UnitPrice:   50.0,
				Quantity:    1,
			},
			{
				ProductCode: "Bread",
				UnitPrice:   50.0,
				Quantity:    2,
			},
		},
	}

	mockOrder.On("ProcessItems", mock.Anything,
		req).
		Return(nil).
		Once()

	res, err := adapter.ProcessItems(context.Background(), &pb.ProcessItemsRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: 1,
			Balance:    200.0,
		},

		OrderItems: []*pb.OrderItem{
			{
				ProductCode: "Juice",
				UnitPrice:   50.0,
				Quantity:    1,
			},
			{
				ProductCode: "Bread",
				UnitPrice:   50.0,
				Quantity:    2,
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "successfully process items", res.GetMessage())
	mockOrder.AssertExpectations(t)
}

func TestProcessItems_Error(t *testing.T) {
	mockOrder := new(OrderPort)
	adapter := NewAdapter(mockOrder, 50052)

	req := &order.ProcessItemsRequest{
		Customer: order.CustomerEntity{
			CustomerID: 1,
			Balance:    100.0,
		},
		OrderItems: []*order.OrderItemDTO{
			{
				ProductCode: "Juice",
				UnitPrice:   50.0,
				Quantity:    1,
			},
			{
				ProductCode: "Bread",
				UnitPrice:   50.0,
				Quantity:    2,
			},
		},
	}

	mockOrder.On("ProcessItems", mock.Anything,
		req).
		Return(e.BALANCE_INSUFFICIENT).
		Once()

	res, err := adapter.ProcessItems(context.Background(), &pb.ProcessItemsRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: 1,
			Balance:    100.0,
		},

		OrderItems: []*pb.OrderItem{
			{
				ProductCode: "Juice",
				UnitPrice:   50.0,
				Quantity:    1,
			},
			{
				ProductCode: "Bread",
				UnitPrice:   50.0,
				Quantity:    2,
			},
		},
	})

	assert.Error(t, err)
	assert.Equal(t, e.BALANCE_INSUFFICIENT.Error(),
		res.GetMessage())

	mockOrder.AssertExpectations(t)
}
