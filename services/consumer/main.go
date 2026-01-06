package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/dzon2000/eda/consumer/internal/consumer"
	"github.com/dzon2000/eda/consumer/internal/schema"
	"github.com/joho/godotenv"
)

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

	consumer := consumer.New(cfg.Kafka)
	defer consumer.Stop()
	registry := schema.New(cfg.SchemaRegistry)

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := consumer.Start(registry); err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down gracefully...")
	consumer.Stop()

}
