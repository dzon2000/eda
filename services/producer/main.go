package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/segmentio/kafka-go"
)

type schemaRequest struct {
	Schema string `json:"schema"`
}

type schemaResponse struct {
	ID int `json:"id"`
}

const schemaID = 1

func main() {
	fmt.Println("Producer Service")
	// Load Avro schema from file
	schemaBytes := readAvroSchema()
	// registerAvroSchema(schemaBytes)
	log.Println("Schema loaded and registered")
	// Create Avro codec to validate schema
	codec, err := goavro.NewCodec(string(schemaBytes))
	if err != nil {
		log.Fatal(err)
	}

	nativeMessage := map[string]interface{}{
		"orderId":    "123",
		"customerId": "cust-42",
		"amount":     99.99,
		"createdAt":  time.Now().UTC().Format(time.RFC3339),
	}

	value := createAvroMessage(codec, nativeMessage)
	log.Println("Avro message created")
	writer := kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "orders.v1",
		Balancer: &kafka.Hash{},
	}
	log.Println("Sending message to Kafka topic...")
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("123"), // orderId
			Value: value,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Message sent successfully")
}

func createAvroMessage(codec *goavro.Codec, native map[string]interface{}) []byte {
	avroBytes, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	// Magic byte
	buf.WriteByte(0)
	// Schema ID (big-endian)
	binary.Write(buf, binary.BigEndian, int32(schemaID))
	// Avro payload
	buf.Write(avroBytes)
	value := buf.Bytes()
	return value
}

func registerAvroSchema(schemaBytes []byte) {
	payload, _ := json.Marshal(schemaRequest{
		Schema: string(schemaBytes),
	})

	resp, err := http.Post(
		"http://schema-registry:8081/subjects/orders.v1-value/versions",
		"application/vnd.schemaregistry.v1+json",
		bytes.NewReader(payload),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Failed to register schema: %s", resp.Status)
	}
}

func readAvroSchema() []byte {
	schemaBytes, err := os.ReadFile("../../schemas/order-created.avsc")
	if err != nil {
		log.Fatal(err)
	}
	return schemaBytes
}
