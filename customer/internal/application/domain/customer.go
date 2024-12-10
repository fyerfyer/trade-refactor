package domain

import "time"

type Customer struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Customer) CanPlaceOrder() bool {
	return c.Status == "active"
}

func (c *Customer) DeductBalance(amount float32) bool {
	if c.Balance >= amount {
		c.Balance -= amount
		return true
	}
	return false
}

func (c *Customer) AddBalance(amount float32) {
	c.Balance += amount
}
