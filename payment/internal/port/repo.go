package port

import (
	"context"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/domain"
)

type Repository interface {
	Save(ctx context.Context, payment *domain.Payment) error
}
