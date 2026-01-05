package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("Consumer Service")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "orders.v1",
		GroupID:  "orders-consumer",
		MinBytes: 1e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		err = handleMessage(msg.Value)
		if err != nil {
			log.Printf("failed to handle message: %v", err)
		}
	}

}

func handleMessage(value []byte) error {
	if len(value) < 5 {
		return fmt.Errorf("invalid message")
	}

	// 1. Magic byte
	if value[0] != 0 {
		return fmt.Errorf("unknown magic byte")
	}

	// 2. Schema ID
	schemaID := int(binary.BigEndian.Uint32(value[1:5]))

	// 3. Avro payload
	avroPayload := value[5:]

	return deserialize(schemaID, avroPayload)
}

func deserialize(schemaID int, payload []byte) error {
	codec, err := getCodec(schemaID)
	if err != nil {
		return err
	}

	native, _, err := codec.NativeFromBinary(payload)
	if err != nil {
		return err
	}

	record := native.(map[string]interface{})
	log.Printf("event: %+v", record)

	return nil
}

var codecCache = map[int]*goavro.Codec{}

func getCodec(schemaID int) (*goavro.Codec, error) {
	if codec, ok := codecCache[schemaID]; ok {
		return codec, nil
	}

	url := fmt.Sprintf(
		"http://schema-registry:8081/schemas/ids/%d",
		schemaID,
	)

	resp, err := http.Get(url)
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

	codecCache[schemaID] = codec
	return codec, nil
}
