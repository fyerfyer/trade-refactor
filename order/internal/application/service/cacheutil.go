package service

import (
	"encoding/json"
	"fmt"

	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/order/internal/application/domain"
)

func LookUpOrderInCache(cache cachePort.Cache, id uint64, status string) (bool, *domain.Order) {
	var order *domain.Order
	if cache.Exists(GetOrderKey(id, status)) {
		data, err := cache.Get(GetOrderKey(id, status))
		if err != nil {
			return false, nil
		}
		json.Unmarshal(data, &order)
	} else {
		return false, nil
	}

	return true, order
}

func GetOrderKey(id uint64, status string) string {
	return fmt.Sprintf("Order_%v_%v", id, status)
}
