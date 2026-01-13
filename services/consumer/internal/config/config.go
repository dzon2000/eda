package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Kafka          KafkaConfig
	SchemaRegistry SchemaRegistryConfig
	Environment    string
}

type KafkaConfig struct {
	Brokers    []string
	Topic      string
	GroupID    string
	MinBytes   int
	MaxBytes   int
	DLQTopic   string
	MaxRetries int
}

type SchemaRegistryConfig struct {
	URL         string
	Timeout     time.Duration
	DLQSchemaID int
}

func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Kafka: KafkaConfig{
			Brokers:    getBrokersFromEnv(),
			Topic:      getEnv("KAFKA_TOPIC", "orders.v1"),
			GroupID:    getEnv("KAFKA_GROUP_ID", "orders-consumer"),
			MinBytes:   getEnvAsInt("KAFKA_MIN_BYTES", 1000),
			MaxBytes:   getEnvAsInt("KAFKA_MAX_BYTES", 10000000),
			DLQTopic:   getEnv("KAFKA_DLQ_TOPIC", "orders.dlq"),
			MaxRetries: getEnvAsInt("KAFKA_MAX_RETRIES", 5),
		},
		SchemaRegistry: SchemaRegistryConfig{
			URL:         getEnv("SCHEMA_REGISTRY_URL", "http://schema-registry:8081"),
			Timeout:     getEnvAsDuration("SCHEMA_REGISTRY_TIMEOUT", 10*time.Second),
			DLQSchemaID: getEnvAsInt("SCHEMA_REGISTRY_DLQ_SCHEMA_ID", 67),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker is required")
	}
	if c.Kafka.Topic == "" {
		return fmt.Errorf("Kafka topic is required")
	}
	if c.SchemaRegistry.URL == "" {
		return fmt.Errorf("Schema Registry URL is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getBrokersFromEnv() []string {
	brokers := getEnv("KAFKA_BROKERS", "kafka:9092")
	return strings.Split(brokers, ",")
}
