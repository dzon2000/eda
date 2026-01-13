package dlq

import (
	"context"
	"strings"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/events"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/segmentio/kafka-go"
)

type DLQProducer interface {
	Send(ctx context.Context, msg kafka.Message, err error) error
	Close() error
}

type KafkaDLQProducer struct {
	writer  *kafka.Writer
	encoder *schema.Encoder
}

func NewKafkaDLQProducer(kafkaConfig config.KafkaConfig, encoder *schema.Encoder) *KafkaDLQProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      kafkaConfig.Brokers,
		Topic:        kafkaConfig.DLQTopic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
		MaxAttempts:  kafkaConfig.MaxRetries,
	})
	return &KafkaDLQProducer{
		writer:  writer,
		encoder: encoder,
	}
}

func (p *KafkaDLQProducer) Send(
	ctx context.Context,
	msg kafka.Message,
	cause error,
) error {

	event := events.NewOrderDLQEvent(
		extractEventID(msg.Value),
		msg.Topic,
		msg.Partition,
		msg.Offset,
		classifyError(cause),
		cause.Error(),
		msg.Value,
	)

	value, err := p.encoder.Encode(event.ToMap())
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   msg.Key,
		Value: value,
	})
}

func (p *KafkaDLQProducer) Close() error {
	return p.writer.Close()
}

func extractEventID(value []byte) string {
	if len(value) < 5 {
		return "malformed-message"
	}
	return "unknown-event-id"
}

func classifyError(err error) string {
	switch {
	case strings.Contains(err.Error(), "schema"):
		return "schema_error"
	case strings.Contains(err.Error(), "deserialize"):
		return "deserialization_error"
	case strings.Contains(err.Error(), "invalid message"):
		return "invalid_message"
	default:
		return "processing_error"
	}
}
