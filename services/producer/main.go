package main

import (
	"context"
	"log"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/dzon2000/eda/producer/internal/events"
	"github.com/dzon2000/eda/producer/internal/producer"
	"github.com/dzon2000/eda/producer/internal/schema"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Producer Service")
	_ = godotenv.Load(".env.development")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	schemaContent, err := schema.NewRegistry(cfg.Schema).LoadSchemaFromFile()
	if err != nil {
		log.Fatal(err)
	}
	encoder, err := schema.NewEncoder(schemaContent, cfg.Schema.SchemaID)
	if err != nil {
		log.Fatal(err)
	}
	producer := producer.New(cfg.Kafka)
	defer producer.Close()

	discount := 10.0
	event, _ := events.NewOrderCreatedEvent("456", "cost-999", 199.00, &discount)

	value, err := encoder.Encode(event.ToMap())
	if err != nil {
		log.Fatal(err)
	}

	err = producer.Send(context.Background(), []byte(event.OrderID), value)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Message sent successfully")
}
