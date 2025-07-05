package serialize

import (
	"bytes"
	"fmt"
)

// Serialize takes any supported value and serializes it to the buffer
// This function appends the serialized value to the existing buffer
func Serialize(value interface{}, data *bytes.Buffer) error {

	switch v := value.(type) {

	case uint64, uint32, uint16, uint8, uint:
		switch uv := v.(type) {
		case uint64:
			SerializeUint64(uv, data)
		case uint32:
			SerializeUint64(uint64(uv), data)
		case uint16:
			SerializeUint64(uint64(uv), data)
		case uint8:
			SerializeUint64(uint64(uv), data)
		case uint:
			SerializeUint64(uint64(uv), data)
		}

	case int64, int32, int16, int8, int:
		switch iv := v.(type) {
		case int64:
			SerializeSint64(iv, data)
		case int32:
			SerializeSint64(int64(iv), data)
		case int16:
			SerializeSint64(int64(iv), data)
		case int8:
			SerializeSint64(int64(iv), data)
		case int:
			SerializeSint64(int64(iv), data)
		}

	case float64, float32:
		// Use type-specific serialize methods
		switch fv := v.(type) {
		case float32:
			SerializeFloat32(fv, data)
		case float64:
			SerializeFloat64(fv, data)
		}

	case bool:
		SerializeBoolean(v, data)

	case string:
		SerializeString(v, data)

	case []byte:
		SerializeBuffer(v, data)

	case *bytes.Buffer:
		SerializeBuffer(v.Bytes(), data)

	default:
		return fmt.Errorf("unsupported type for serialization: %T", value)
	}
	return nil
}

// SerializeAll serializes multiple values into a single buffer
// This function appends all serialized values to the existing buffer
func SerializeAll(values []interface{}, data *bytes.Buffer) error {
	for i, value := range values {
		if err := Serialize(value, data); err != nil {
			return fmt.Errorf("error serializing value at index %d: %w", i, err)
		}
	}
	return nil
}
