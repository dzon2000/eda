package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Kafka       KafkaConfig
	Schema      SchemaConfig
	Environment string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type SchemaConfig struct {
	FilePath string // Path to .avsc file
	SchemaID int
}

func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Kafka: KafkaConfig{
			Brokers: getBrokersFromEnv(),
			Topic:   getEnv("KAFKA_TOPIC", "orders.v1"),
		},
		Schema: SchemaConfig{
			FilePath: getEnv("SCHEMA_FILE_PATH", "order_created.avsc"),
			SchemaID: getEnvAsInt("SCHEMA_ID", 2),
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

func getBrokersFromEnv() []string {
	brokers := getEnv("KAFKA_BROKERS", "kafka:9092")
	return strings.Split(brokers, ",")
}
