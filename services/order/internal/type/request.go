package order

type CreateOrderRequest struct {
    OrderID    string   `json:"order_id"`
    CustomerID string   `json:"customer_id"`
    Amount     float64  `json:"amount"`
    Discount   *float64 `json:"discount,omitempty"`
}