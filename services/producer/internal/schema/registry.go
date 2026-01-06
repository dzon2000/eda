package schema

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dzon2000/eda/producer/internal/config"
)

type Registry struct {
	config config.SchemaConfig
	client *http.Client
}

func NewRegistry(cfg config.SchemaConfig) *Registry {
	return &Registry{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterSchema registers a schema and returns its ID
func (r *Registry) RegisterSchema(subject string, schema string) (int, error) {
	return 0, nil // Placeholder for actual implementation
}

// LoadSchemaFromFile reads .avsc file
func (r *Registry) LoadSchemaFromFile() (string, error) {
	schemaBytes, err := os.ReadFile(r.config.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %w", err)
	}
	return string(schemaBytes), nil
}
