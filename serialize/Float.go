package serialize

import (
	"bytes"
	"ebe/types"
	"encoding/binary"
	"fmt"
	"math"
)

func SerializeFloat64(value float64, data *bytes.Buffer) {
	// This function appends the serialized float64 to the existing buffer
	// Figure out what size of float is needed for the data
	switch {

	case value >= -math.SmallestNonzeroFloat32 && value <= math.MaxFloat32:
		// Set the header for the type in the data buffer
		data.WriteByte(types.CreateHeader(types.Float, 4))

		// Write the value
		_ = binary.Write(data, binary.LittleEndian, float32(value))

	default:
		// Set the header for the type in the data buffer
		data.WriteByte(types.CreateHeader(types.Float, 8))

		// Write the value
		_ = binary.Write(data, binary.LittleEndian, float64(value))
	}
}

func SerializeFloat32(value float32, data *bytes.Buffer) {
	// This function appends the serialized float32 to the existing buffer
	// float32 always uses 4 bytes
	data.WriteByte(types.CreateHeader(types.Float, 4))
	_ = binary.Write(data, binary.LittleEndian, value)
}

func DeserializeFloat(data []byte) (float64, []byte, error) {

	if len(data) == 0 {
		return 0, data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)

	// Make sure the data is a valid integer value
	if headerType != types.Float {
		return 0, data, fmt.Errorf("expected Float type, got %v", headerType)
	}

	if length != 4 && length != 8 {
		return 0, data, fmt.Errorf("invalid float length: expected 4 or 8, got %d", length)
	}

	if len(data) < int(length) {
		return 0, data, fmt.Errorf("insufficient data: need %d bytes, got %d", length, len(data))
	}

	// Convert byte slice into a reader
	buf := bytes.NewReader(data)

	// Read the value
	// If the value is a float32 then read into float32 then copy to float64
	var value float64 = 0
	if length == 4 {
		var float32Value float32 = 0
		err := binary.Read(buf, binary.LittleEndian, &float32Value)
		if err != nil {
			return 0, data, fmt.Errorf("failed to read float32: %w", err)
		}
		data = data[4:]
		value = float64(float32Value)
	} else {
		err := binary.Read(buf, binary.LittleEndian, &value)
		if err != nil {
			return 0, data, fmt.Errorf("failed to read float64: %w", err)
		}
		data = data[8:]
	}

	return value, data, nil
}

func DeserializeFloat32(data []byte) (float32, []byte, error) {
	if len(data) == 0 {
		return 0, data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)

	// Make sure the data is a valid float value
	if headerType != types.Float {
		return 0, data, fmt.Errorf("expected Float type, got %v", headerType)
	}

	if length != 4 {
		return 0, data, fmt.Errorf("invalid float32 length: expected 4, got %d", length)
	}

	if len(data) < 4 {
		return 0, data, fmt.Errorf("insufficient data: need 4 bytes, got %d", len(data))
	}

	// Convert byte slice into a reader
	buf := bytes.NewReader(data)

	// Read the float32 value
	var value float32 = 0
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		return 0, data, fmt.Errorf("failed to read float32: %w", err)
	}
	data = data[4:]

	return value, data, nil
}

func DeserializeFloat64(data []byte) (float64, []byte, error) {
	return DeserializeFloat(data)
}
