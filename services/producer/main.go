package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/dzon2000/eda/producer/internal/db"
	"github.com/dzon2000/eda/producer/internal/events"
	"github.com/dzon2000/eda/producer/internal/producer"
	"github.com/dzon2000/eda/producer/internal/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Publisher struct {
	outboxRepo     *db.OutboxRepository
	dbPool         *sql.DB
	schemaRegistry *schema.Registry
	kafkaProducer  *producer.Producer
}

func NewPublisher(outboxRepo *db.OutboxRepository, dbPool *sql.DB, schemaRegistry *schema.Registry, kafkaProducer *producer.Producer) *Publisher {
	return &Publisher{
		outboxRepo:     outboxRepo,
		dbPool:         dbPool,
		schemaRegistry: schemaRegistry,
		kafkaProducer:  kafkaProducer,
	}
}

func (p *Publisher) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.publishBatch(ctx); err != nil {
				log.Println("publish batch failed", err)
			}
		}
	}
}

func (p *Publisher) publishBatch(ctx context.Context) error {
	log.Println("Starting batch publishing.")
	tx, err := p.dbPool.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	events, err := p.outboxRepo.FetchPending(ctx, tx, 100)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		log.Println("No pending events to publish.")
		return tx.Commit()
	}

	for _, e := range events {
		if err := p.publishOne(ctx, e); err != nil {
			log.Println("Failed to publish event:", err)
			p.outboxRepo.MarkError(ctx, tx, e.ID, err)
			return err
		}
		log.Println("Published event ID:", e.ID)
		if err := p.outboxRepo.MarkSent(ctx, tx, e.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *Publisher) publishOne(ctx context.Context, event events.OutboxEvent) error {
	log.Printf("Publishing event ID: %s, Type: %s, Schema: %d", event.ID, event.EventType, event.SchemaVersion)
	codec, err := p.schemaRegistry.GetCodec(event.SchemaVersion)
	if err != nil {
		return fmt.Errorf("failed to get codec for schema ID %d: %w", event.SchemaVersion, err)
	}

	encoder, err := schema.NewEncoder(codec, event.SchemaVersion)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}

	orderEvent, err := p.DeserializeOrderCreated(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to deserialize event ID %s: %w", event.ID, err)
	}
	avroBytes, err := encoder.Encode(orderEvent.ToMap())
	if err != nil {
		return fmt.Errorf("failed to encode event ID %s: %w", event.ID, err)
	}

	if err := p.kafkaProducer.Send(ctx, event, avroBytes); err != nil {
		return fmt.Errorf("failed to send event ID %s to Kafka: %w", event.ID, err)
	}

	log.Printf("Successfully published event ID: %s", event.ID)
	return nil
}

func (p *Publisher) DeserializeOrderCreated(payload []byte) (*events.OrderCreated, error) {
	var orderCreated events.OrderCreated
	err := json.Unmarshal(payload, &orderCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize OrderCreated event: %w", err)
	}
	return &orderCreated, nil
}

func (p *Publisher) Close() {
	p.kafkaProducer.Close()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Producer Service")
	_ = godotenv.Load(".env.development")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DBName,
	)
	dbPool, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbPool.SetMaxOpenConns(20)
	dbPool.SetMaxIdleConns(5)
	dbPool.SetConnMaxLifetime(time.Hour)

	outboxRepo := db.NewOutboxRepository(dbPool)
	schemaRegistry := schema.NewRegistry(cfg.Schema)
	producer := producer.New(cfg.Kafka)
	defer producer.Close()
	publisher := NewPublisher(outboxRepo, dbPool, schemaRegistry, producer)

	ctx := context.Background()
	go publisher.Run(ctx)

	log.Println("Publisher started. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received, stopping...")
	publisher.Close()
	time.Sleep(1 * time.Second)
	log.Println("Shutdown complete")
}
