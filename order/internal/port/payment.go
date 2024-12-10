package port

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/payment"
)

type PaymentPort interface {
	Charge(ctx context.Context, req *payment.ChargeRequest) error
}
