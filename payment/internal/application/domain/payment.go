package domain

type Payment struct {
	CustomerID uint64  `json:"customer_id"`
	OrderID    uint64  `json:"order_id"`
	TotalPrice float32 `json:"total_price"`
	Status     string  `json:"status"`
	Message    string  `json:"message"`
}
