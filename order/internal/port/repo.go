package port

import (
	"context"

	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
)

type Repository interface {
	Save(ctx context.Context, order *domain.Order) error
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, orderID uint64) error
	GetUnpaidOrders(ctx context.Context, customerID uint64) ([]domain.Order, error)
	GetUnpaidOrder(ctx context.Context, orderID uint64) (*domain.Order, error)
}
