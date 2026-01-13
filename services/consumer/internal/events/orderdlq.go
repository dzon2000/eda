package events

import "time"

type OrderDLQEvent struct {
	EventID       string
	OriginalTopic string
	Partition     int
	Offset        int64
	ErrorType     string
	ErrorMessage  string
	Payload       []byte
	FailedAt      string
}

func NewOrderDLQEvent(
	eventID string,
	originalTopic string,
	partition int,
	offset int64,
	errorType string,
	errorMessage string,
	payload []byte,
) *OrderDLQEvent {
	return &OrderDLQEvent{
		EventID:       eventID,
		OriginalTopic: originalTopic,
		Partition:     partition,
		Offset:        offset,
		ErrorType:     errorType,
		ErrorMessage:  errorMessage,
		Payload:       payload,
		FailedAt:      time.Now().UTC().Format(time.RFC3339),
	}
}

// ToMap converts to format expected by Avro encoder
func (e *OrderDLQEvent) ToMap() map[string]interface{} {
	return map[string]interface{}{
		// eventId is a nullable union, so we must wrap it
		"eventId":       map[string]interface{}{"string": e.EventID},
		"originalTopic": e.OriginalTopic,
		"partition":     e.Partition,
		"offset":        e.Offset,
		"errorType":     e.ErrorType,
		"errorMessage":  e.ErrorMessage,
		"payload":       e.Payload,
		"failedAt":      e.FailedAt,
	}
}
