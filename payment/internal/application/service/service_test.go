package service

import (
	"context"
	"testing"

	"github.com/fyerfyer/trade-dependency/dto/payment"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mock test need repairing
// haven't included cache
func TestChargeSuccess(t *testing.T) {
	mockRepo := new(Repository)
	mockCache := new(Cache)
	srv := NewService(mockRepo, mockCache)

	req := &payment.ChargeRequest{
		Customer: payment.CustomerEntity{
			CustomerID: 1,
			Balance:    200.0,
		},
		Order: payment.OrderEntity{
			OrderID:    1,
			TotalPrice: 100.0,
		},
	}

	// set mock behaviour
	mockRepo.On("Save", mock.Anything,
		mock.AnythingOfType("*domain.Payment")).
		Return(nil)

	err := srv.Charge(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestChargeInsufficientBalance(t *testing.T) {
	mockRepo := new(Repository)
	mockCache := new(Cache)
	srv := NewService(mockRepo, mockCache)

	req := &payment.ChargeRequest{
		Customer: payment.CustomerEntity{
			CustomerID: 1,
			Balance:    50.0,
		},
		Order: payment.OrderEntity{
			OrderID:    1,
			TotalPrice: 100.0,
		},
	}

	err := srv.Charge(context.Background(), req)
	assert.ErrorIs(t, err, e.BALANCE_INSUFFICIENT)
}
