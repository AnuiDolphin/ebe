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

	// Fast path type assertions - handle ALL common types before any reflection
	switch v := value.(type) {

	// Primitive types (non-pointer)
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

	// Array types (non-pointer)
	case []int, []int32, []int64, []int8, []int16:
		return serializeIntArray(v, w)
	case []uint, []uint32, []uint64, []uint16:
		return serializeUintArray(v, w)
	case []float32, []float64:
		return serializeFloatArray(v, w)
	case []string:
		return serializeStringArray(v, w)
	case []bool:
		return serializeArray(reflect.ValueOf(v), w)

	// Map types (non-pointer)
	case map[string]int:
		return serializeMap(v, w)
	case map[string]string:
		return serializeMap(v, w)
	case map[string]interface{}:
		return serializeMap(v, w)
	case map[int]string:
		return serializeMap(v, w)
	case map[int]int:
		return serializeMap(v, w)

	// Pointer types - fast paths
	case *int:
		if v != nil {
			return serializeSint(int64(*v), w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *int32:
		if v != nil {
			return serializeSint(int64(*v), w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *int64:
		if v != nil {
			return serializeSint(*v, w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *uint:
		if v != nil {
			return serializeUint(uint64(*v), w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *uint32:
		if v != nil {
			return serializeUint(uint64(*v), w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *uint64:
		if v != nil {
			return serializeUint(*v, w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *float32:
		if v != nil {
			return serializeFloat(float64(*v), w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *float64:
		if v != nil {
			return serializeFloat(*v, w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *bool:
		if v != nil {
			return serializeBoolean(*v, w)
		}
		return fmt.Errorf("cannot serialize nil pointer")
	case *string:
		if v != nil {
			return serializeString(*v, w)
		}
		return fmt.Errorf("cannot serialize nil pointer")

	default:
		// Only use reflection as last resort for unknown types
		return serializeWithReflection(value, w)
	}
}

// serializeWithReflection handles types that couldn't be handled by fast path type assertions
func serializeWithReflection(value interface{}, w io.Writer) error {
	rv := reflect.ValueOf(value)
	
	// Handle other pointer types with reflection
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return fmt.Errorf("cannot serialize nil pointer")
		}
		rv = rv.Elem()
		value = rv.Interface() // Update value to the dereferenced value
		// Try fast path again after dereferencing
		return Serialize(value, w)
	}

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