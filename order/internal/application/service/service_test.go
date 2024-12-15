package service

import (
	"context"
	"testing"
	"time"

	"github.com/fyerfyer/trade-dependency/dto/order"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessItemsSuccess(t *testing.T) {
	mockRepo := new(Repository)
	mockCache := new(Cache)
	mockPayment := new(PaymentPort)
	srv := NewService(mockRepo, mockCache, mockPayment)

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

	mockRepo.On("Save", mock.Anything,
		mock.AnythingOfType("*domain.Order")).
		Return(nil).
		Once()

	mockRepo.On("Update", mock.Anything,
		mock.AnythingOfType("*domain.Order")).
		Return(nil).
		Once()

	// mock the payment behaviour
	mockPayment.On("Charge", mock.Anything,
		mock.AnythingOfType("*payment.ChargeRequest")).
		Return(nil).
		Once()

	// mock the cache behaviour in handlePayment
	mockCache.On("Set", mock.AnythingOfType("string"),
		mock.AnythingOfType("*domain.Order"),
		mock.AnythingOfType("int")).
		Return(nil)

	err := srv.ProcessItems(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertNumberOfCalls(t, "Set", 2)
	mockPayment.AssertExpectations(t)
}

func TestProcessItemsFailure(t *testing.T) {
	mockRepo := new(Repository)
	mockCache := new(Cache)
	mockPayment := new(PaymentPort)
	srv := NewService(mockRepo, mockCache, mockPayment)

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

	mockRepo.On("Save", mock.Anything,
		mock.AnythingOfType("*domain.Order")).
		Return(nil).
		Once()

	// mock the payment behaviour
	mockPayment.On("Charge", mock.Anything,
		mock.AnythingOfType("*payment.ChargeRequest")).
		Return(e.BALANCE_INSUFFICIENT).
		Once()

	// mock the cache behaviour in handlePayment
	mockCache.On("Set", mock.AnythingOfType("string"),
		mock.AnythingOfType("*domain.Order"),
		mock.AnythingOfType("int")).
		Return(nil).
		Once()

	err := srv.ProcessItems(context.Background(), req)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockPayment.AssertExpectations(t)
}

func TestProcessOrderSuccess(t *testing.T) {
	mockRepo := new(Repository)
	mockCache := new(Cache)
	mockPayment := new(PaymentPort)
	srv := NewService(mockRepo, mockCache, mockPayment)

	req := &order.ProcessOrderRequest{
		Customer: order.CustomerEntity{
			CustomerID: 1,
			Balance:    200.0,
		},

		OrderID: 1,
	}

	mockRepo.On("GetUnpaidOrder", mock.Anything,
		mock.AnythingOfType("uint64")).
		Return(&domain.Order{
			ID:         1,
			CustomerID: 1,
			Items: []domain.OrderItem{
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
			Status:    "unpaid",
			CreatedAt: time.Now(),
		}, nil)

	mockRepo.On("Update", mock.Anything,
		mock.AnythingOfType("*domain.Order")).
		Return(nil)

	mockCache.On("Exists", mock.AnythingOfType("string")).
		Return(false)

	mockCache.On("Set", mock.Anything,
		mock.AnythingOfType("*domain.Order"),
		mock.AnythingOfType("int")).
		Return(nil)

	mockPayment.On("Charge", mock.Anything,
		mock.AnythingOfType("*payment.ChargeRequest")).
		Return(nil)

	err := srv.ProcessOrder(context.Background(), req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockPayment.AssertExpectations(t)
}
