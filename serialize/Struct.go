package serialize

import (
	"fmt"
	"reflect"
	"io"
)

// SerializeStruct serializes a struct by serializing each exported field in order
// Unexported fields are skipped
func SerializeStruct(value interface{}, w io.Writer) error {
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
