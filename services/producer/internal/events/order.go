package events

import (
	"fmt"
	"time"
)

type OrderCreated struct {
	OrderID    string
	CustomerID string
	Amount     float64
	CreatedAt  string
	Discount   *float64 // nullable
}

// ToMap converts to format expected by Avro encoder
func (o *OrderCreated) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"orderId":    o.OrderID,
		"customerId": o.CustomerID,
		"amount":     o.Amount,
		"createdAt":  o.CreatedAt,
	}
	if o.Discount != nil {
		data["discount"] = map[string]interface{}{
			"double": *o.Discount,
		}
	} else {
		data["discount"] = nil
	}
	return data
}

func NewOrderCreatedEvent(orderID, customerID string, amount float64, discount *float64) (*OrderCreated, error) {
	if orderID == "" {
		return nil, fmt.Errorf("orderID is required")
	}
	if amount < 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	return &OrderCreated{
		OrderID:    orderID,
		CustomerID: customerID,
		Amount:     amount,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		Discount:   discount,
	}, nil
}
