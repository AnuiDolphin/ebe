package serialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// Serialize takes any supported value and serializes it to the writer
func Serialize(value interface{}, w io.Writer) error {

	// If the value is a pointer, dereference it
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("cannot serialize nil pointer")
		}
		rv = rv.Elem()
		value = rv.Interface() // Update value to the dereferenced value
	}

	switch v := value.(type) {

	case json.RawMessage:
		return serializeJson(v, w)

	case uint64:
		return serializeUint(v, w)

	case uint32:
		return serializeUint(uint64(v), w)

	case uint16:
		return serializeUint(uint64(v), w)

	case uint8:
		return serializeUint(uint64(v), w)

	case uint:
		return serializeUint(uint64(v), w)

	case int64:
		return serializeSint(v, w)

	case int32:
		return serializeSint(int64(v), w)

	case int16:
		return serializeSint(int64(v), w)

	case int8:
		return serializeSint(int64(v), w)

	case int:
		return serializeSint(int64(v), w)

	case float64:
		return serializeFloat(v, w)

	case float32:
		return serializeFloat(float64(v), w)

	case bool:
		return serializeBoolean(v, w)

	case string:
		return serializeString(v, w)

	case []byte:
		return serializeBuffer(v, w)

	case *bytes.Buffer:
		return serializeBuffer(v.Bytes(), w)

	// Fast paths for common array types
	case []int, []int32, []int64, []int8, []int16:
		return serializeIntArray(v, w)

	case []uint, []uint32, []uint64, []uint16:
		return serializeUintArray(v, w)

	case []float32, []float64:
		return serializeFloatArray(v, w)

	case []string:
		return serializeStringArray(v, w)

	default:

		// Handle structs by serializing each exported field in order
		if rv.Kind() == reflect.Struct {
			return serializeStruct(value, w)
		}

		// Check if it's an array or slice
		if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
			return serializeArray(rv, w)
		}

		// Check if it's a map
		if rv.Kind() == reflect.Map {
			return serializeMap(value, w)
		}

		return fmt.Errorf("unsupported type for serialization: %T", value)
	}
}
