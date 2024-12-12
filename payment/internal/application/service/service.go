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

func (s *PaymentService) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	var (
		p       *domain.Payment
		success bool
		err     error
	)

	if success, p = LookUpPaymentInCache(s.cache, req.CustomerID); !success {
		if p, err = s.repo.Get(ctx, req.CustomerID); err != nil {
			return nil, e.PAYMENT_NOT_FOUND
		}

		if err := s.cache.Set(GetPaymentKey(req.CustomerID), p, 0); err != nil {
			return nil, e.FAILED_TO_UPDATE_CACHE_ERROR
		}
	}

	return &payment.GetPaymentResponse{
		Payment: payment.PaymentDTO{
			OrderID:    p.OrderID,
			CustomerID: p.CustomerID,
			TotalPrice: p.TotalPrice,
			Status:     p.Status,
			Message:    p.Message,
		},
	}, nil
}
