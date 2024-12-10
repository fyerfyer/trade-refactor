package domain

import "time"

type Order struct {
	ID         uint64      `json:"id"`
	CustomerID uint64      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

func (o *Order) TotalPrice() float32 {
	var total float32
	for _, item := range o.Items {
		total += item.UnitPrice * float32(item.Quantity)
	}
	return total
}
