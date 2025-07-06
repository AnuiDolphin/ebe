package serialize

import (
	"bytes"
	"fmt"
	"reflect"
)

// Serialize takes any supported value and serializes it to the buffer
// This function appends the serialized value to the existing buffer
func Serialize(value interface{}, data *bytes.Buffer) error {

	// Handle structs by serializing each exported field in order
	rv := reflect.ValueOf(value)

	// If the value is a pointer, dereference it
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("cannot serialize nil pointer")
		}
		rv = rv.Elem()
		value = rv.Interface() // Update value to the dereferenced value
	}

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
		// Check if it's a struct
		if rv.Kind() == reflect.Struct {
			return SerializeStruct(value, data)
		}
		
		// Check if it's an array or slice
		if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
			return SerializeArray(value, data)
		}

		return fmt.Errorf("unsupported type for serialization: %T", value)
	}
	return nil
}
