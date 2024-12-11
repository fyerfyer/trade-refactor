package service

import (
	"context"
	"time"

	"github.com/fyerfyer/trade-dependency/dto/customer"
	"github.com/fyerfyer/trade-dependency/dto/order"
	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-dependency/pkg/e"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
	"github.com/fyerfyer/trade-refactor/customer/internal/port"
)

type CustomerService struct {
	repo  port.Repository
	cache cachePort.Cache
	order port.OrderPort
}

func NewService(repo port.Repository, cache cachePort.Cache, order port.OrderPort) *CustomerService {
	return &CustomerService{
		repo:  repo,
		cache: cache,
		order: order,
	}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *customer.CreateCustomerRequest) (*customer.CreateCustomerResponse, error) {
	// find in the cache first
	if s.cache.Exists(req.CustomerName) {
		return nil, e.DUPLICATE_CUSTOMER_ERROR
	}

	// else, find in db
	// if we can find in db, then update our cache
	// we store id and name key separately
	if c, _ := s.repo.GetByName(ctx, req.CustomerName); c != nil {
		cKey := GetCustomerKey(req.CustomerName)
		s.cache.Set(req.CustomerName, 1, 0)
		s.cache.Set(cKey, c, 0)
		return nil, e.DUPLICATE_CUSTOMER_ERROR
	}

	c := &domain.Customer{
		Name:      req.CustomerName,
		Status:    "active",
		Balance:   0,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, c); err != nil {
		return nil, err
	}

	// store customer into cache
	if err := s.cache.Set(GetCustomerKey(c.Name), c, 0); err != nil {
		return &customer.CreateCustomerResponse{
			CustomerID: c.ID,
			Success:    false,
		}, e.FAILED_TO_STORE_CACHE_ERROR
	}

	return &customer.CreateCustomerResponse{
		CustomerID: c.ID,
		Success:    true,
	}, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, req *customer.GetCustomerRequest) (*customer.GetCustomerResponse, error) {
	// look up in cache first
	var (
		c       *domain.Customer
		success bool
		err     error
	)
	success, c = LookUpCustomerInCache(s.cache, req.CustomerName)
	if !success {
		c, err = s.repo.GetByName(ctx, req.CustomerName)
		if err != nil {
			return nil, e.CUSTOMER_NOT_FOUND
		}

		// store customer into cache
		if err := s.cache.Set(GetCustomerKey(c.Name), c, 0); err != nil {
			return nil, e.FAILED_TO_STORE_CACHE_ERROR
		}
	}

	return &customer.GetCustomerResponse{
		Customer: &customer.CustomerDTO{
			ID:        c.ID,
			Name:      c.Name,
			Status:    c.Status,
			Balance:   c.Balance,
			CreatedAt: c.CreatedAt,
		},
	}, nil
}

func (s *CustomerService) DeactivateCustomer(ctx context.Context, req *customer.DeactivateCustomerRequest) error {
	res, err := s.GetCustomer(ctx, &customer.GetCustomerRequest{
		CustomerName: req.CustomerName,
	})

	if err != nil {
		return e.CUSTOMER_NOT_FOUND
	}

	if res.Customer.Status != "inactive" {
		res.Customer.Status = "inactive"
		domainCustomer := &domain.Customer{
			ID:        res.Customer.ID,
			Name:      res.Customer.Name,
			Status:    res.Customer.Status,
			Balance:   res.Customer.Balance,
			CreatedAt: res.Customer.CreatedAt,
		}
		if err := s.cache.Set(GetCustomerKey(req.CustomerName), domainCustomer, 0); err != nil {
			return e.FAILED_TO_UPDATE_CACHE_ERROR
		}
		if err := s.repo.Update(ctx, domainCustomer); err != nil {
			return e.FAILED_TO_UPDATE_DB_ERROR
		}
	}

	return nil
}

func (s *CustomerService) ActivateCustomer(ctx context.Context, req *customer.ActivateCustomerRequest) error {
	res, err := s.GetCustomer(ctx, &customer.GetCustomerRequest{
		CustomerName: req.CustomerName,
	})

	if err != nil {
		return e.CUSTOMER_NOT_FOUND
	}

	if res.Customer.Status != "active" {
		res.Customer.Status = "active"
		domainCustomer := &domain.Customer{
			ID:        res.Customer.ID,
			Name:      res.Customer.Name,
			Status:    res.Customer.Status,
			Balance:   res.Customer.Balance,
			CreatedAt: res.Customer.CreatedAt,
		}
		if err := s.cache.Set(GetCustomerKey(req.CustomerName), domainCustomer, 0); err != nil {
			return e.FAILED_TO_UPDATE_CACHE_ERROR
		}
		if err := s.repo.Update(ctx, domainCustomer); err != nil {
			return e.FAILED_TO_UPDATE_DB_ERROR
		}
	}

	return nil
}

func (s *CustomerService) SubmitOrder(ctx context.Context, req *customer.SubmitOrderRequest) error {
	// look up in cache first
	var (
		c       *domain.Customer
		success bool
		err     error
	)

	success, c = LookUpCustomerInCache(s.cache, req.CustomerName)
	if !success {
		c, err = s.repo.GetByName(ctx, req.CustomerName)
		if err != nil {
			return e.CUSTOMER_NOT_FOUND
		}

		// store customer into cache
		if err := s.cache.Set(GetCustomerKey(c.Name), c, 0); err != nil {
			return e.FAILED_TO_STORE_CACHE_ERROR
		}
	}

	if !c.CanPlaceOrder() {
		return e.CUSTOMER_INACTIVE
	}

	var items []*order.OrderItemDTO
	for _, item := range req.OrderItems {
		items = append(items, &order.OrderItemDTO{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		})
	}

	// the CreateOrderMethod will process the order
	err = s.order.ProcessItems(ctx, &order.ProcessItemsRequest{
		Customer: order.CustomerEntity{
			CustomerID: c.ID,
			Balance:    c.Balance,
		},
		OrderItems: items,
	})

	// failed to pay
	if err != nil {
		return err
	}

	return nil
}

func (s *CustomerService) PayOrder(ctx context.Context, req *customer.PayOrderRequest) error {
	var (
		c       *domain.Customer
		success bool
		err     error
	)

	success, c = LookUpCustomerInCache(s.cache, req.CustomerName)
	if !success {
		c, err = s.repo.GetByName(ctx, req.CustomerName)
		if err != nil {
			return e.CUSTOMER_NOT_FOUND
		}

		// store customer into cache
		if err := s.cache.Set(GetCustomerKey(req.CustomerName), c, 0); err != nil {
			return e.FAILED_TO_STORE_CACHE_ERROR
		}
	}

	if !c.CanPlaceOrder() {
		return e.CUSTOMER_INACTIVE
	}

	return s.order.ProcessOrder(ctx, &order.ProcessOrderRequest{
		Customer: order.CustomerEntity{
			CustomerID: c.ID,
			Balance:    c.Balance,
		},
		OrderID: req.OrderID,
	})
}

func (s *CustomerService) GetUnpaidOrders(ctx context.Context, req *customer.GetUnpaidOrdersRequest) (*customer.GetUnpaidOrdersResponse, error) {
	res, err := s.order.GetUnpaidOrder(ctx, &order.GetUnpaidOrdersRequest{
		CustomerID: req.CustomerID,
	})

	if err != nil {
		return nil, err
	}

	var unpaidOrders []*customer.Order
	for _, o := range res.Orders {
		var items []*customer.OrderItem
		for _, item := range o.Items {
			items = append(items, &customer.OrderItem{
				ProductCode: item.ProductCode,
				UnitPrice:   item.UnitPrice,
				Quantity:    item.Quantity,
			})
		}
		unpaidOrders = append(unpaidOrders, &customer.Order{
			OrderID:   o.OrderID,
			Items:     items,
			Status:    o.Status,
			CreatedAt: time.Now(),
		})
	}
	return &customer.GetUnpaidOrdersResponse{
		UnpaidOrders: unpaidOrders,
	}, nil
}

func (s *CustomerService) StoreBalance(ctx context.Context, req *customer.StoreBalanceRequest) error {
	// look up in cache first
	var (
		c       *domain.Customer
		success bool
		err     error
	)
	success, c = LookUpCustomerInCache(s.cache, req.CustomerName)
	if !success {
		c, err = s.repo.GetByName(ctx, req.CustomerName)
		if err != nil {
			return e.CUSTOMER_NOT_FOUND
		}
	}

	c.AddBalance(req.Balance)

	// update cache and database
	if err := s.cache.Set(GetCustomerKey(req.CustomerName), c, 0); err != nil {
		return e.FAILED_TO_UPDATE_CACHE_ERROR
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return e.FAILED_TO_UPDATE_DB_ERROR
	}

	return nil
}
