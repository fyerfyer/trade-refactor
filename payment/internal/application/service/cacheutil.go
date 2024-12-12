package service

import (
	"encoding/json"
	"fmt"

	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/payment/internal/application/domain"
)

func LookUpPaymentInCache(cache cachePort.Cache, id uint64) (bool, *domain.Payment) {
	var payment *domain.Payment
	if cache.Exists(GetPaymentKey(id)) {
		data, err := cache.Get(GetPaymentKey(id))
		if err != nil {
			return false, nil
		}
		json.Unmarshal(data, &payment)
		return true, payment
	} else {
		return false, nil
	}
}

func GetPaymentKey(id uint64) string {
	return fmt.Sprintf("Payment_%v", id)
}
