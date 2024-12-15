package service

import (
	context "context"
	"testing"

	"github.com/fyerfyer/trade-dependency/dto/customer"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestSubmitOrder(t *testing.T) {
	// Prepare mocks
	mockCache := new(Cache)
	mockRepo := new(Repository)
	mockOrder := new(OrderPort)

	// Create the service instance
	service := &CustomerService{
		cache: mockCache,
		repo:  mockRepo,
		order: mockOrder,
	}

	// Test data
	req := &customer.SubmitOrderRequest{
		CustomerName: "john_doe",
		OrderItems: []*customer.OrderItem{
			{ProductCode: "P123", Quantity: 1, UnitPrice: 100},
		},
	}

	customer := &domain.Customer{
		Name:    "john_doe",
		ID:      1,
		Balance: 200,
		Status:  "active",
	}

	mockCache.On("Exists", mock.AnythingOfType("string")).Return(false)
	mockCache.On("LookUpCustomerInCache", "john_doe").Return(false, nil)
	mockRepo.On("GetByName", mock.Anything, "john_doe").Return(customer, nil)
	mockCache.On("Set", "Customer_john_doe", customer, 0).Return(nil)
	mockOrder.On("ProcessItems", mock.Anything, mock.Anything).Return(nil)

	err := service.SubmitOrder(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockOrder.AssertExpectations(t)
}
