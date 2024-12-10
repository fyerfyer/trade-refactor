package service

import (
	"context"
	"github.com/fyerfyer/trade-dependency/dto/payment"
	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/domain"
	"github.com/fyerfyer/trade-refactor/payment/internal/port"
)

type PaymentService struct {
	repo  port.Repository
	cache cachePort.Cache
}

func NewService(repo port.Repository, cache cachePort.Cache) *PaymentService {
	return &PaymentService{
		repo:  repo,
		cache: cache,
	}
}

// we pass order and customer entity to payment service
func (s *PaymentService) Charge(ctx context.Context, req *payment.ChargeRequest) error {
	// fix: I put the purchase logic into customer service
	// because the DeductBalance method belongs to customer domain
	// so the related logic should be in customer service, too

	if req.Customer.Balance < req.Order.TotalPrice {
		return e.BALANCE_INSUFFICIENT
	}

	pay := &domain.Payment{
		CustomerID: req.Customer.CustomerID,
		OrderID:    req.Order.OrderID,
		TotalPrice: req.Order.TotalPrice,
		Status:     "success",
		Message:    "successfully process payment",
	}

	if err := s.repo.Save(ctx, pay); err != nil {
		return e.FAILED_TO_STORE_DB_ERROR
	}

	return nil
}
