package grpc

import (
	"context"
	"testing"

	"github.com/fyerfyer/trade-dependency/pkg/e"
	pb "github.com/fyerfyer/trade-dependency/proto/grpc/payment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCharge_Success(t *testing.T) {
	mockSrv := new(PaymentPort)
	adapter := NewAdapter(mockSrv, 50051)

	req := &pb.ChargeRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: 1,
			Balance:    200.0,
		},

		Order: &pb.OrderEntity{
			OrderId:    1,
			TotalPrice: 100.0,
		},
	}

	mockSrv.On("Charge", mock.Anything,
		mock.AnythingOfType("*payment.ChargeRequest")).
		Return(nil)

	res, err := adapter.Charge(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "success", res.Status)
	assert.Equal(t, "successfully purchase the payment", res.Message)
}

func TestCharge_Failure(t *testing.T) {
	mockSrv := new(PaymentPort)
	adapter := NewAdapter(mockSrv, 50051)

	req := &pb.ChargeRequest{
		Customer: &pb.CustomerEntity{
			CustomerId: 1,
			Balance:    50.0,
		},

		Order: &pb.OrderEntity{
			OrderId:    1,
			TotalPrice: 100.0,
		},
	}

	mockSrv.On("Charge", mock.Anything,
		mock.AnythingOfType("*payment.ChargeRequest")).
		Return(e.BALANCE_INSUFFICIENT)

	res, err := adapter.Charge(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, "failure", res.Status)
	assert.Equal(t, e.BALANCE_INSUFFICIENT.Error(), res.Message)
}
