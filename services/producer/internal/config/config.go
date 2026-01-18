package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Kafka          KafkaConfig
	Schema         SchemaConfig
	Environment    string
	ProducerConfig ProducerConfig
	DB             DBConfig
}

type KafkaConfig struct {
	Brokers    []string
	Topic      string
	MaxRetries int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type SchemaConfig struct {
	FilePath string // Path to .avsc file
	SchemaID int
	URL      string
	Timeout  string
}

type ProducerConfig struct {
	MaxRetries int
}

func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Kafka: KafkaConfig{
			Brokers:    getBrokersFromEnv(),
			Topic:      getEnv("KAFKA_TOPIC", "orders.v1"),
			MaxRetries: getEnvAsInt("KAFKA_MAX_RETRIES", 10),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "producer_db"),
		},
		Schema: SchemaConfig{
			FilePath: getEnv("SCHEMA_REGISTRY_FILE_PATH", "order_created.avsc"),
			SchemaID: getEnvAsInt("SCHEMA_REGISTRY_ID", 3),
			URL:      getEnv("SCHEMA_REGISTRY_URL", "http://schema-registry:8081"),
			Timeout:  getEnv("SCHEMA_REGISTRY_TIMEOUT", "10s"),
		},
		ProducerConfig: ProducerConfig{
			MaxRetries: getEnvAsInt("PRODUCER_MAX_RETRIES", 5),
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
