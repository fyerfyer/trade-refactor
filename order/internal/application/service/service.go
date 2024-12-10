package service

import (
	"time"

	"github.com/fyerfyer/trade-dependency/dto/order"
	"github.com/fyerfyer/trade-dependency/dto/payment"
	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
	"github.com/fyerfyer/trade-refactor/order/internal/port"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	repo    port.Repository
	cache   cachePort.Cache
	payment port.PaymentPort
}

func NewService(repo port.Repository, cache cachePort.Cache, paymentPort port.PaymentPort) *OrderService {
	return &OrderService{
		repo:    repo,
		cache:   cache,
		payment: paymentPort,
	}
}

func (s *OrderService) ProcessItems(ctx context.Context, req *order.ProcessItemsRequest) error {
	var items []domain.OrderItem
	for _, item := range req.OrderItems {
		items = append(items, domain.OrderItem{
			ProductCode: item.ProductCode,
			UnitPrice:   item.UnitPrice,
			Quantity:    item.Quantity,
		})
	}

	o := &domain.Order{
		CustomerID: req.Customer.CustomerID,
		Items:      items,
		Status:     "unpaid",
		CreatedAt:  time.Now(),
	}

	err := s.repo.Save(ctx, o)
	if err != nil {
		return e.FAILED_TO_STORE_DB_ERROR
	}

	return s.handlePayment(ctx, o, req.Customer)
}

func (s *OrderService) ProcessOrder(ctx context.Context, req *order.ProcessOrderRequest) error {
	// fix: I put the get order logic into order microservice

	// find in cache first
	var (
		o       *domain.Order
		success bool
		err     error
	)

	if success, o = LookUpOrderInCache(s.cache, req.OrderID, "unpaid"); !success {
		if o, err = s.repo.GetUnpaidOrder(ctx, req.OrderID); err != nil {
			return e.CUSTOMER_NOT_FOUND
		}

		// update the cache
		if err = s.cache.Set(GetOrderKey(req.OrderID, "unpaid"), o, 0); err != nil {
			return e.FAILED_TO_UPDATE_CACHE_ERROR
		}
	}

	return s.handlePayment(ctx, o, req.Customer)
}

func (s *OrderService) handlePayment(ctx context.Context, o *domain.Order, c order.CustomerEntity) error {
	err := s.payment.Charge(ctx, &payment.ChargeRequest{
		Order: payment.OrderEntity{
			OrderID:    o.ID,
			TotalPrice: o.TotalPrice(),
		},
		Customer: payment.CustomerEntity{
			CustomerID: o.CustomerID,
			Balance:    c.Balance,
		},
	})

	if err == nil {
		o.Status = "success"
		if err := s.cache.Set(GetOrderKey(o.ID, o.Status), o, 0); err != nil {
			return e.FAILED_TO_UPDATE_CACHE_ERROR
		}
		if err := s.repo.Update(ctx, o); err != nil {
			return e.FAILED_TO_UPDATE_DB_ERROR
		}

		return nil
	}

	o.Status = "unpaid"
	if err := s.cache.Set(GetOrderKey(o.ID, o.Status), o, 0); err != nil {
		return e.FAILED_TO_UPDATE_CACHE_ERROR
	}
	if err := s.repo.Update(ctx, o); err != nil {
		return e.FAILED_TO_UPDATE_DB_ERROR
	}

	return status.New(codes.Internal, err.Error()).Err()
}

func (s *OrderService) GetUnpaidOrders(ctx context.Context, req *order.GetUnpaidOrdersRequest) (*order.GetUnpaidOrdersResponse, error) {
	orders, err := s.repo.GetUnpaidOrders(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}

	var ordersEntity []order.OrderEntity

	for _, o := range orders {
		var items []*order.OrderItemDTO
		for _, item := range o.Items {
			items = append(items, &order.OrderItemDTO{
				ProductCode: item.ProductCode,
				UnitPrice:   item.UnitPrice,
				Quantity:    item.Quantity,
			})
		}

		ordersEntity = append(ordersEntity, order.OrderEntity{
			OrderID: o.ID,
			Items:   items,
			Status:  o.Status,
		})
	}

	return &order.GetUnpaidOrdersResponse{
		Orders: ordersEntity,
	}, nil
}
