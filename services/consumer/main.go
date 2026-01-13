package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/consumer"
	"github.com/dzon2000/eda/consumer/internal/dlq"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/joho/godotenv"
)

func initializeDLQProducer(cfg *config.Config, encoder *schema.Encoder) (dlq.DLQProducer, error) {
	dlqProducer := dlq.NewKafkaDLQProducer(cfg.Kafka, encoder)
	return dlqProducer, nil
}

func initializeEncoder(registry *schema.Registry, cfg *config.Config) *schema.Encoder {
	dlqCodec, err := registry.GetCodec(cfg.SchemaRegistry.DLQSchemaID)
	if err != nil {
		log.Fatalf("Failed to get codec from schema registry: %v", err)
	}
	encoder, _ := schema.NewEncoder(dlqCodec, cfg.SchemaRegistry.DLQSchemaID)
	return encoder
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Consumer Service")
	_ = godotenv.Load(".env.development")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting consumer in %s environment", cfg.Environment)
	log.Printf("Kafka brokers: %v", cfg.Kafka.Brokers)
	log.Printf("Topic: %s", cfg.Kafka.Topic)
	registry := schema.New(cfg.SchemaRegistry)
	dlqEncoder := initializeEncoder(registry, cfg)
	dlqProducer, err := initializeDLQProducer(cfg, dlqEncoder)
	if err != nil {
		log.Fatalf("Failed to initialize DLQ producer: %v", err)
	}
	consumer, err := consumer.New(cfg.Kafka, registry, dlqProducer)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Stop()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down gracefully...")
	consumer.Stop()
}
