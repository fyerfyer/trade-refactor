package service

import (
	"encoding/json"

	cachePort "github.com/fyerfyer/trade-dependency/pkg/cache"
	"github.com/fyerfyer/trade-refactor/customer/internal/application/domain"
)

func LookUpCustomerInCache(cache cachePort.Cache, name string) (bool, *domain.Customer) {
	var customer *domain.Customer
	if cache.Exists(GetCustomerKey(name)) {
		data, err := cache.Get(GetCustomerKey(name))
		if err != nil {
			return false, nil
		}
		json.Unmarshal(data, &customer)
	} else {
		return false, nil
	}

	return true, customer
}

func GetCustomerKey(name string) string {
	return "Customer_" + name
}
