package db

import (
	"context"
	"database/sql"

	"github.com/dzon2000/eda/order/internal/payload"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) InsertOrder(ctx context.Context, tx *sql.Tx, req payload.CreateOrderRequest) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO orders (order_id, customer_id, amount, discount, status)
        VALUES ($1, $2, $3, $4, 'CREATED')
        ON CONFLICT (order_id) DO NOTHING
    `, req.OrderID, req.CustomerID, req.Amount, req.Discount)
	return err
}

func (r *OrderRepository) InsertOrderCreatedOutbox(
	ctx context.Context,
	tx *sql.Tx,
	payload []byte,
) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO outbox_events (
            id, aggregate_type, aggregate_id,
            event_type, payload, schema_version
        ) VALUES (
            gen_random_uuid(), 'order', $1,
            'OrderCreated', $2, 1
        )
    `, extractOrderID(payload), payload)
	return err
}

func extractOrderID(payload []byte) string {
	return "etc"
}
