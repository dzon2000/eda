package consumer

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/deduplicator"
	"github.com/dzon2000/eda/consumer/internal/dlq"
	"github.com/dzon2000/eda/consumer/internal/events"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	kafkaConfig config.KafkaConfig
	reader      *kafka.Reader
	dedup       *deduplicator.Deduplicator
	dlqProducer dlq.DLQProducer
	registry    *schema.Registry
}

func New(kafkaConfig config.KafkaConfig, registry *schema.Registry, dlqProducer dlq.DLQProducer) (*Consumer, error) {
	return &Consumer{
		kafkaConfig: kafkaConfig,
		dedup:       deduplicator.New(),
		dlqProducer: dlqProducer,
		registry:    registry,
	}, nil
}

func (c *Consumer) Start() error {
	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.kafkaConfig.Brokers,
		Topic:          c.kafkaConfig.Topic,
		GroupID:        c.kafkaConfig.GroupID,
		MinBytes:       c.kafkaConfig.MinBytes,
		MaxBytes:       c.kafkaConfig.MaxBytes,
		CommitInterval: 0, // manual commits
	})
	for {
		ctx := context.Background()
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue // Don't fatal, keep running
		}
		if err := c.processMessage(ctx, msg); err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	orderEvent, err := c.handleMessage(c.registry, msg.Value)
	if err != nil {
		return c.handleProcessingError(ctx, msg, err)
	}

	if c.dedup.Seen(orderEvent.EventID) {
		log.Printf("Duplicate event detected: %s", orderEvent.EventID)
		return c.commitMessage(ctx, msg)
	}

	log.Printf("Processed OrderCreated event: %+v", orderEvent)
	return c.commitMessage(ctx, msg)
}

func (c *Consumer) handleProcessingError(ctx context.Context, msg kafka.Message, err error) error {
	log.Printf("Processing error: %v", err)

	if dlqErr := c.dlqProducer.Send(ctx, msg, err); dlqErr != nil {
		log.Printf("Failed to send to DLQ: %v", dlqErr)
		return fmt.Errorf("both processing and DLQ failed: %w", err)
	}

	// Successfully sent to DLQ, commit to avoid reprocessing
	return c.commitMessage(ctx, msg)
}

func (c *Consumer) commitMessage(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.Printf("Failed to commit message: %v", err)
		return err
	}
	return nil
}

func (c *Consumer) Stop() error {
	c.dlqProducer.Close()
	return c.reader.Close()
}

func (c *Consumer) handleMessage(registry *schema.Registry, value []byte) (*events.OrderCreatedEvent, error) {
	if 1 == 1 {
		return nil, fmt.Errorf("simulated error for testing retries")
	}
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

	codec, err := registry.GetCodec(schemaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get codec for schema ID %d: %w", schemaID, err)
	}

	encoder, _ := schema.NewEncoder(codec, schemaID)

	payload, err := encoder.Decode(schemaID, avroPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %w", err)
	}
	log.Printf("Received message: %v", payload)
	orderEvent := events.ParseOrderCreated(payload.(map[string]interface{}))
	return orderEvent, nil
}
