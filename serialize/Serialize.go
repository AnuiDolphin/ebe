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

	case json.RawMessage:
		return serializeJson(v, w)

	case uint64:
		return serializeUint64(v, w)

	case uint32:
		return serializeUint64(uint64(v), w)

	case uint16:
		return serializeUint64(uint64(v), w)

	case uint8:
		return serializeUint64(uint64(v), w)

	case uint:
		return serializeUint64(uint64(v), w)

	case int64:
		return serializeSint64(v, w)

	case int32:
		return serializeSint64(int64(v), w)

	case int16:
		return serializeSint64(int64(v), w)

	case int8:
		return serializeSint64(int64(v), w)

	case int:
		return serializeSint64(int64(v), w)

	case float64:
		return serializeFloat64(v, w)

	case float32:
		return serializeFloat32(v, w)

	case bool:
		return serializeBoolean(v, w)

	case string:
		return serializeString(v, w)

	case []byte:
		return serializeBuffer(v, w)

	case *bytes.Buffer:
		return serializeBuffer(v.Bytes(), w)

	default:
		// Check if it's a struct
		if rv.Kind() == reflect.Struct {
			return serializeStruct(value, w)
		}

		// Check if it's an array or slice
		if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
			return serializeArray(value, w)
		}

		return fmt.Errorf("unsupported type for serialization: %T", value)
	}
}
