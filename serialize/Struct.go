package serialize

import (
	"fmt"
	"io"
	"reflect"
)

// SerializeStruct serializes a struct by serializing each exported field in order
// Unexported fields are skipped
func serializeStruct(value interface{}, w io.Writer) error {
	rv := reflect.ValueOf(value)

	// Ensure we have a struct
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", rv.Kind())
	}

	// Serialize each exported field in order
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		if field.PkgPath != "" { // unexported field
			continue
		}

		// Recursively call Serialize to serialize the field value
		if err := Serialize(rv.Field(i).Interface(), w); err != nil {
			return fmt.Errorf("error serializing field %s: %w", field.Name, err)
		}
	}

	return nil
}

// deserializeStruct deserializes data from a stream into a struct by deserializing each field in order
func deserializeStruct(r io.Reader, structValue reflect.Value) error {
	structType := structValue.Type()

	// Iterate through each field in the struct
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Create a pointer to the field for deserialization
		fieldPtr := field.Addr().Interface()

		// Recursively call Deserialize to deserialize into this field
		err := Deserialize(r, fieldPtr)
		if err != nil {
			return fmt.Errorf("failed to deserialize field '%s': %w", fieldType.Name, err)
		}
	}

	return nil
}
