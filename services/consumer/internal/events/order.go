package events

type OrderCreatedEvent struct {
	EventID    string
	OrderID    string
	CustomerID string
	Amount     float64
	Discount   *float64
}

// Parse from Avro deserialized map
func ParseOrderCreated(data map[string]interface{}) *OrderCreatedEvent {
	event := &OrderCreatedEvent{
		EventID:    data["eventId"].(string),
		OrderID:    data["orderId"].(string),
		CustomerID: data["customerId"].(string),
		Amount:     data["amount"].(float64),
	}

	if d, ok := data["discount"].(map[string]interface{}); ok {
		val := d["double"].(float64)
		event.Discount = &val
	}

	return event
}
