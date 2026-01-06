package consumer

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	kafkaConfig config.KafkaConfig
	reader      *kafka.Reader
}

func New(kafkaConfig config.KafkaConfig) *Consumer {
	return &Consumer{
		kafkaConfig: kafkaConfig,
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
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue // Don't fatal, keep running
		}

		if err := c.handleMessage(registry, msg.Value); err != nil {
			log.Printf("Failed to handle message: %v", err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	return c.reader.Close()
}

func (c *Consumer) handleMessage(registry *schema.Registry, value []byte) error {
	if len(value) < 5 {
		return fmt.Errorf("invalid message")
	}

	// 1. Magic byte
	if value[0] != 0 {
		return fmt.Errorf("unknown magic byte")
	}

	// 2. Schema ID
	schemaID := int(binary.BigEndian.Uint32(value[1:5]))

	// 3. Avro payload
	avroPayload := value[5:]
	payload, err := registry.Deserialize(schemaID, avroPayload)
	if err != nil {
		return fmt.Errorf("failed to deserialize message: %w", err)
	}
	log.Printf("Received message: %v", payload)
	return nil
}
