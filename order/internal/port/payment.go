package port

import (
	"context"

	"github.com/fyerfyer/trade-dependency/dto/payment"
)

type PaymentPort interface {
	Charge(ctx context.Context, req *payment.ChargeRequest) error
	GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error)
}
