package payload

import "errors"

type CreateOrderRequest struct {
	OrderID    string   `json:"order_id"`
	CustomerID string   `json:"customer_id"`
	Amount     float64  `json:"amount"`
	Discount   *float64 `json:"discount,omitempty"`
}

func (r CreateOrderRequest) Validate() error {
	if r.OrderID == "" {
		return errors.New("order_id is required")
	}
	if r.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	return nil
}
