package db

import (
	"context"
	"database/sql"

	"github.com/dzon2000/eda/producer/internal/events"
	"github.com/google/uuid"
)

type OutboxRepository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) FetchPending(
	ctx context.Context,
	tx *sql.Tx,
	limit int,
) ([]events.OutboxEvent, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, schema_version, created_at
        FROM outbox_events
        WHERE status = 'PENDING'
        ORDER BY created_at
        LIMIT $1
        FOR UPDATE SKIP LOCKED
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventsList []events.OutboxEvent
	for rows.Next() {
		var event events.OutboxEvent
		if err := rows.Scan(
			&event.ID,
			&event.AggregateType,
			&event.AggregateID,
			&event.EventType,
			&event.Payload,
			&event.SchemaVersion,
			&event.CreatedAt,
		); err != nil {
			return nil, err
		}
		eventsList = append(eventsList, event)
	}
	return eventsList, rows.Err()
}

func (r *OutboxRepository) MarkSent(
	ctx context.Context,
	tx *sql.Tx,
	eventID uuid.UUID,
) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE outbox_events
		SET status = 'PUBLISHED', published_at = NOW()
		WHERE id = $1
	`, eventID)
	return err
}
func (r *OutboxRepository) MarkError(
	ctx context.Context,
	tx *sql.Tx,
	eventID uuid.UUID,
	cause error,
) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE outbox_events
		SET status = 'ERROR', last_error = $2, updated_at = NOW()
		WHERE id = $1
	`, eventID, cause.Error())
	return err
}
