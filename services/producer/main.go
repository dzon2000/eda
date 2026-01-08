package main

import (
	"context"
	"flag"
	"log"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/dzon2000/eda/producer/internal/events"
	"github.com/dzon2000/eda/producer/internal/producer"
	"github.com/dzon2000/eda/producer/internal/schema"
	"github.com/joho/godotenv"
)

var orderID string
var customerID string
var amount float64
var discount float64

func init() {
	flag.StringVar(&orderID, "order-id", "", "Order ID for the event")
	flag.StringVar(&customerID, "customer-id", "", "Customer ID for the event")
	flag.Float64Var(&amount, "amount", 0.0, "Amount for the event")
	flag.Float64Var(&discount, "discount", 0.0, "Discount for the event")
	flag.Parse()
	_ = godotenv.Load(".env.development")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Producer Service")
	log.Printf("Creating event for order ID: %s, customer ID: %s, amount: %.2f, discount: %.2f", orderID, customerID, amount, discount)
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

	event, _ := events.NewOrderCreatedEvent(orderID, customerID, amount, &discount)

	value, err := encoder.Encode(event.ToMap())
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < cfg.ProducerConfig.MaxRetries; i++ {
		err = producer.Send(context.Background(), []byte(event.OrderID), value)
		if err == nil {
			break // Can be improved
		}
	}

	log.Println("Message sent successfully")
}
