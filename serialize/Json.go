package serialize

import (
	"ebe/types"
	"ebe/utils"
	"encoding/json"
	"fmt"
	"io"
)

// SerializeJson serializes a json.RawMessage
// The jsonMessage parameter should come from json.Marshal() wrapped as json.RawMessage
func SerializeJson(jsonMessage json.RawMessage, w io.Writer) error {
	// Write header with JSON type
	utils.WriteByte(w, types.CreateHeader(types.Json, 0x00))
	if err := SerializeUint64(uint64(len(jsonMessage)), w); err != nil {
		return err
	}

	// Write the JSON bytes
	_, err := w.Write(jsonMessage)
	return err
}

// DeserializeJson deserializes JSON data and unmarshals it into the provided output
func DeserializeJson(data []byte, out interface{}) ([]byte, error) {

	if len(data) == 0 {
		return data, fmt.Errorf("empty data")
	}

	// Verify the header type
	header := data[0]
	headerType := types.TypeFromHeader(header)
	if headerType != types.Json {
		return data, fmt.Errorf("expected Json type, got %s", types.TypeName(headerType))
	}
	remaining := data[1:]

	// Length always follows as a UInt (no nibble optimization)
	length, remaining, err := DeserializeUint64(remaining)
	if err != nil {
		return remaining, fmt.Errorf("failed to deserialize JSON length: %w", err)
	}

	// Check if we have enough data
	if uint64(len(remaining)) < length {
		return remaining, fmt.Errorf("insufficient data for JSON: expected %d bytes, got %d", length, len(remaining))
	}

	// Extract the JSON bytes
	jsonBytes := remaining[:length]
	remaining = remaining[length:]

	// Unmarshal the JSON into the output
	if err := json.Unmarshal(jsonBytes, out); err != nil {
		return remaining, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return remaining, nil
}
