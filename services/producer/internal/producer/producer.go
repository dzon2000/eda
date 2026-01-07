package producer

import (
	"context"

	"github.com/dzon2000/eda/producer/internal/config"
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
	})
	return &Producer{
		writer: writer,
		config: cfg,
	}
}

func (p *Producer) Send(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
