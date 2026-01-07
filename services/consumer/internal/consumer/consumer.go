package consumer

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/deduplicator"
	"github.com/dzon2000/eda/consumer/internal/events"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	kafkaConfig config.KafkaConfig
	reader      *kafka.Reader
	dedup       *deduplicator.Deduplicator
}

func New(kafkaConfig config.KafkaConfig) *Consumer {
	return &Consumer{
		kafkaConfig: kafkaConfig,
		dedup:       deduplicator.New(),
	}
}

func (c *Consumer) Start(registry *schema.Registry) error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.kafkaConfig.Brokers,
		Topic:    c.kafkaConfig.Topic,
		GroupID:  c.kafkaConfig.GroupID,
		MinBytes: c.kafkaConfig.MinBytes,
		MaxBytes: c.kafkaConfig.MaxBytes,
	})
	for {
		ctx := context.Background()
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue // Don't fatal, keep running
		}

		orderEvent, err := c.handleMessage(registry, msg.Value)
		if err != nil {
			log.Printf("Failed to handle message: %v", err)
			continue
		}
		if c.dedup.Seen(orderEvent.EventID) {
			log.Printf("Duplicate event detected: %s", orderEvent.EventID)
			continue
		} else {
			// Do something with the event
			log.Printf("Processed OrderCreated event: %+v", orderEvent)
		}
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
		}

	}
}

func (c *Consumer) Stop() error {
	return c.reader.Close()
}

func (c *Consumer) handleMessage(registry *schema.Registry, value []byte) (*events.OrderCreatedEvent, error) {
	if len(value) < 5 {
		return nil, fmt.Errorf("invalid message")
	}

	// 1. Magic byte
	if value[0] != 0 {
		return nil, fmt.Errorf("unknown magic byte")
	}

	// 2. Schema ID
	schemaID := int(binary.BigEndian.Uint32(value[1:5]))

	// 3. Avro payload
	avroPayload := value[5:]
	payload, err := registry.Deserialize(schemaID, avroPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %w", err)
	}
	log.Printf("Received message: %v", payload)
	orderEvent := events.ParseOrderCreated(payload.(map[string]interface{}))
	return orderEvent, nil
}
