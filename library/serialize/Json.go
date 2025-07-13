package serialize

import (
	"ebe-library/types"
	"ebe-library/utils"
	"encoding/json"
	"fmt"
	"io"
)

// SerializeJson serializes a json.RawMessage
// The jsonMessage parameter should come from json.Marshal() wrapped as json.RawMessage
func serializeJson(jsonMessage json.RawMessage, w io.Writer) error {
	// Write header with JSON type
	utils.WriteByte(w, types.CreateHeader(types.Json, 0x00))
	if err := serializeUint(uint64(len(jsonMessage)), w); err != nil {
		return err
	}

	// Write the JSON bytes
	_, err := w.Write(jsonMessage)
	return err
}

// DeserializeJson deserializes JSON data from a stream and unmarshals it into the provided output
func deserializeJson(r io.Reader, out interface{}) error {
	// Read the header using utils.ReadHeader
	headerType, _, err := utils.ReadHeader(r)
	if err != nil {
		return fmt.Errorf("failed to read JSON header: %w", err)
	}

	// Verify the header type
	if headerType != types.Json {
		return fmt.Errorf("expected Json type, got %s", types.TypeName(headerType))
	}

	// Length always follows as a UInt (no nibble optimization)
	length, err := deserializeUint(r)
	if err != nil {
		return fmt.Errorf("failed to deserialize JSON length: %w", err)
	}

	// Read the JSON bytes
	jsonBytes := make([]byte, length)
	_, err = io.ReadFull(r, jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to read JSON data: %w", err)
	}

	// Unmarshal the JSON into the output
	if err := json.Unmarshal(jsonBytes, out); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
