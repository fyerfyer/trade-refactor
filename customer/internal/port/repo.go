package port

import (
	"context"

	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
)

type Repository interface {
	Save(ctx context.Context, c *domain.Customer) error
	Update(ctx context.Context, c *domain.Customer) error
	GetByName(ctx context.Context, customerName string) (*domain.Customer, error)
}
