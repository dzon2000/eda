package producer

import (
	"context"
	"strconv"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/dzon2000/eda/producer/internal/events"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	config config.KafkaConfig
}

func New(cfg config.KafkaConfig) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.Brokers,
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
		MaxAttempts:  cfg.MaxRetries,
		Async:        false,
	})
	return &Producer{
		writer: writer,
		config: cfg,
	}
}

func (p *Producer) Send(ctx context.Context, event events.OutboxEvent, avroBytes []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID.String()),
		Value: avroBytes,
		Headers: []kafka.Header{
			{Key: "event_id", Value: []byte(event.ID.String())},
			{Key: "event_type", Value: []byte(event.EventType)},
			{Key: "schema_version", Value: []byte(strconv.Itoa(event.SchemaVersion))},
		},
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
