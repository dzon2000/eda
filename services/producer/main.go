package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/dzon2000/eda/producer/internal/db"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	// schemaRegistry := schema.NewRegistry(cfg.Schema)

	outboxRepo := db.NewOutboxRepository(dbPool)

	tx, err := dbPool.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	events, err := outboxRepo.FetchPending(context.Background(), tx, 10)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Fetched %d pending events from outbox", len(events))
	for _, event := range events {
		log.Printf("Outbox Event: %+v", event)
	}
	// for i := 0; i < cfg.ProducerConfig.MaxRetries; i++ {
	// 	err = producer.Send(context.Background(), []byte(event.OrderID), value)
	// 	if err == nil {
	// 		break // Can be improved
	// 	}
	// }

	log.Println("Message sent successfully")
}
