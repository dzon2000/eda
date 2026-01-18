package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderCreated struct {
	EventID    string   `json:"event_id"`
	OrderID    string   `json:"order_id"`
	CustomerID string   `json:"customer_id"`
	Amount     float64  `json:"amount"`
	Discount   *float64 `json:"discount,omitempty"`
	CreatedAt  string   `json:"created_at"`
}

// ToMap converts to format expected by Avro encoder
func (o *OrderCreated) ToMap() map[string]interface{} {
	data := map[string]interface{}{
		"eventId":    o.EventID,
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
		EventID:    generateEventID(),
		OrderID:    orderID,
		CustomerID: customerID,
		Amount:     amount,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		Discount:   discount,
	}, nil
}

func generateEventID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
