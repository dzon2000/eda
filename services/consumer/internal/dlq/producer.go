package dlq

import (
	"context"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

type DLQProducer interface {
	Send(ctx context.Context, msg kafka.Message, err error) error
}

type KafkaDLQProducer struct {
	writer  *kafka.Writer
	encoder *goavro.Encoder
}

func (p *KafkaDLQProducer) Send(
	ctx context.Context,
	msg kafka.Message,
	cause error,
) error {

	dlqEvent := map[string]interface{}{
		"eventId":       extractEventID(msg.Value),
		"originalTopic": msg.Topic,
		"partition":     msg.Partition,
		"offset":        msg.Offset,
		"errorType":     classifyError(cause),
		"errorMessage":  cause.Error(),
		"payload":       msg.Value,
		"failedAt":      time.Now().UTC().Format(time.RFC3339),
	}

	value, err := p.encoder.Encode(dlqEvent)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: value,
	})
}

func extractEventID(value []byte) string {
	// Placeholder implementation
	return "unknown-event-id"
}

func classifyError(err error) string {
	// Placeholder implementation
	return "processing_error"
}
