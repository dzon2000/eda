package schema

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/linkedin/goavro/v2"
)

type Encoder struct {
	codec    *goavro.Codec
	schemaID int
}

func NewEncoder(codec *goavro.Codec, schemaID int) (*Encoder, error) {
	return &Encoder{
		codec:    codec,
		schemaID: schemaID,
	}, nil
}

func (e *Encoder) Decode(schemaID int, payload []byte) (interface{}, error) {
	log.Printf("Deserializing message with schema ID: %d", schemaID)

	native, _, err := e.codec.NativeFromBinary(payload)
	if err != nil {
		return nil, err
	}

	record := native.(map[string]interface{})

	return record, nil
}

func (e *Encoder) Encode(data map[string]interface{}) ([]byte, error) {
	avroBytes, err := e.codec.BinaryFromNative(nil, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode to avro: %w", err)
	}
	buf := new(bytes.Buffer)
	// Magic byte
	buf.WriteByte(0)
	// Schema ID (big-endian)
	binary.Write(buf, binary.BigEndian, int32(e.schemaID))
	// Avro payload
	buf.Write(avroBytes)
	value := buf.Bytes()
	return value, nil
}
