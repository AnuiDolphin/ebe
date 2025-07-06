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
		return SerializeJson(v, w)

	case uint64:
		return SerializeUint64(v, w)

	case uint32:
		return SerializeUint64(uint64(v), w)

	case uint16:
		return SerializeUint64(uint64(v), w)

	case uint8:
		return SerializeUint64(uint64(v), w)

	case uint:
		return SerializeUint64(uint64(v), w)

	case int64:
		return SerializeSint64(v, w)

	case int32:
		return SerializeSint64(int64(v), w)

	case int16:
		return SerializeSint64(int64(v), w)

	case int8:
		return SerializeSint64(int64(v), w)

	case int:
		return SerializeSint64(int64(v), w)

	case float64:
		return SerializeFloat64(v, w)

	case float32:
		return SerializeFloat32(v, w)

	case bool:
		return SerializeBoolean(v, w)

	case string:
		return SerializeString(v, w)

	case []byte:
		return SerializeBuffer(v, w)

	case *bytes.Buffer:
		return SerializeBuffer(v.Bytes(), w)

	default:
		// Check if it's a struct
		if rv.Kind() == reflect.Struct {
			return SerializeStruct(value, w)
		}

		// Check if it's an array or slice
		if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
			return SerializeArray(value, w)
		}

		return fmt.Errorf("unsupported type for serialization: %T", value)
	}
}
