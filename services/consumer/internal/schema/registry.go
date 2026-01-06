// internal/schema/registry.go
package schema

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzon2000/eda/consumer/internal/config"
	"github.com/linkedin/goavro/v2"
)

type Registry struct {
	config config.SchemaRegistryConfig
	client *http.Client
	cache  map[int]*goavro.Codec
}

func New(cfg config.SchemaRegistryConfig) *Registry {
	return &Registry{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
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

func (r *Registry) Deserialize(schemaID int, payload []byte) (interface{}, error) {
	codec, err := r.GetCodec(schemaID)
	if err != nil {
		return nil, err
	}

	native, _, err := codec.NativeFromBinary(payload)
	if err != nil {
		return nil, err
	}

	record := native.(map[string]interface{})

	return record, nil
}
