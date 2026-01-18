package events

import (
	"time"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID            uuid.UUID
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte // JSON from DB
	SchemaVersion int
	CreatedAt     time.Time
}

func NewOutboxEvent(
	ID uuid.UUID,
	aggregateType string,
	aggregateID string,
	eventType string,
	payload []byte,
	schemaVersion int,
	createdAt time.Time,
) (*OutboxEvent, error) {
	return &OutboxEvent{
		ID:            ID,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		SchemaVersion: schemaVersion,
		CreatedAt:     createdAt,
	}, nil
}
