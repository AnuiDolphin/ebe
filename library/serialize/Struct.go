package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
	"reflect"
)

// SerializeStruct serializes a struct by writing a struct header followed by each exported field
// Format: [Struct Header] [Optional Field Count] [Field Values...]
// Unexported fields are skipped
func serializeStruct(value interface{}, w io.Writer) error {
	rv := reflect.ValueOf(value)

	// Ensure we have a struct
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", rv.Kind())
	}

	// Use cached struct information for performance
	structInfo, err := typeCache.GetStructInfo(rv.Type())
	if err != nil {
		return fmt.Errorf("failed to get struct info: %w", err)
	}

	// Special case: empty structs serialize to 0 bytes (no header)
	if structInfo.Empty {
		return nil
	}

	// Count exported fields using cached information
	fieldCount := 0
	for _, fieldInfo := range structInfo.Fields {
		if fieldInfo.Exported {
			fieldCount++
		}
	}

	// Write struct header with field count optimization
	if err := writeStructHeader(fieldCount, w); err != nil {
		return err
	}

	// Serialize each exported field in order using cached field information
	for _, fieldInfo := range structInfo.Fields {
		if !fieldInfo.Exported {
			continue
		}

		// Get the field value
		fieldValue := rv.Field(fieldInfo.Index)

		// Recursively call Serialize to serialize the field value
		if err := Serialize(fieldValue.Interface(), w); err != nil {
			return fmt.Errorf("error serializing field %s: %w", fieldInfo.Name, err)
		}
	}

	return nil
}

// deserializeStruct deserializes data from a stream into a struct with a pre-read struct header
func deserializeStruct(r io.Reader, header byte, structValue reflect.Value) error {
	// Read and parse struct header
	expectedFieldCount, err := readStructHeader(r, header)
	if err != nil {
		return err
	}

	// Use cached struct information for performance
	structInfo, err := typeCache.GetStructInfo(structValue.Type())
	if err != nil {
		return fmt.Errorf("failed to get struct info: %w", err)
	}

	// Count actual exported fields using cached information
	actualFieldCount := 0
	for _, fieldInfo := range structInfo.Fields {
		if fieldInfo.Exported {
			actualFieldCount++
		}
	}

	// Validate field count matches
	if uint64(actualFieldCount) != expectedFieldCount {
		return fmt.Errorf("struct field count mismatch: expected %d, struct has %d exported fields", expectedFieldCount, actualFieldCount)
	}

	// Deserialize each exported field in order using cached field information
	for _, fieldInfo := range structInfo.Fields {
		if !fieldInfo.Exported {
			continue
		}

		// Get the field value
		field := structValue.Field(fieldInfo.Index)

		// Create a pointer to the field for deserialization
		fieldPtr := field.Addr().Interface()

		// Each field reads its own header
		err = Deserialize(r, fieldPtr)
		if err != nil {
			return fmt.Errorf("failed to deserialize field '%s': %w", fieldInfo.Name, err)
		}
	}

	return nil
}

// writeStructHeader writes the struct header with field count optimization
func writeStructHeader(fieldCount int, w io.Writer) error {
	if fieldCount <= 7 {
		// Small structs: store count directly in header nibble
		header := types.CreateHeader(types.Struct, byte(fieldCount))
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write struct header: %w", err)
		}
	} else {
		// Large structs: use overflow indicator (8) in header, then UInt for actual count
		header := types.CreateHeader(types.Struct, 8)
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write struct header: %w", err)
		}
		
		// Write actual count as standard EBE UInt
		if err := serializeUint(uint64(fieldCount), w); err != nil {
			return fmt.Errorf("failed to write struct field count: %w", err)
		}
	}
	return nil
}

// readStructHeader reads and parses the struct header, returning field count
func readStructHeader(r io.Reader, header byte) (uint64, error) {
	headerType := types.TypeFromHeader(header)
	headerValue := types.ValueFromHeader(header)

	if headerType != types.Struct {
		return 0, fmt.Errorf("expected Struct type, got %v", types.TypeName(headerType))
	}

	// Determine field count
	var fieldCount uint64
	var err error
	if headerValue <= 7 {
		// Small struct: count stored in header
		fieldCount = uint64(headerValue)
	} else if headerValue == 8 {
		// Large struct: read count as UInt
		fieldCount, err = deserializeUintWithHeader(r)
		if err != nil {
			return 0, fmt.Errorf("failed to read struct field count: %w", err)
		}
	} else {
		return 0, fmt.Errorf("invalid struct header value: %d", headerValue)
	}

	return fieldCount, nil
}
