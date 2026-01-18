package schema

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dzon2000/eda/producer/internal/config"
	"github.com/linkedin/goavro/v2"
)

type Registry struct {
	config config.SchemaConfig
	client *http.Client
	cache  map[int]*goavro.Codec
}

func NewRegistry(cfg config.SchemaConfig) *Registry {
	return &Registry{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[int]*goavro.Codec),
	}
}

func (r *Registry) GetCodec(schemaID int) (*goavro.Codec, error) {
	if codec, ok := r.cache[schemaID]; ok {
		return codec, nil
	}
	url := fmt.Sprintf(
		"%s/schemas/ids/%d",
		r.config.URL,
		schemaID,
	)

	resp, err := r.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Schema string `json:"schema"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(res.Schema)
	if err != nil {
		return nil, err
	}

	r.cache[schemaID] = codec
	return codec, nil
}

func (r *Registry) LoadSchemaFromFile() (string, error) {
	schemaBytes, err := os.ReadFile(r.config.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read schema file: %w", err)
	}
	return string(schemaBytes), nil
}
