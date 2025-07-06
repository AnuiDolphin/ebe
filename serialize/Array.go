package serialize

import (
	"bytes"
	"ebe/types"
	"fmt"
	"io"
	"reflect"
)

// This function appends the serialized array to the existing writer
func SerializeArray(value interface{}, w io.Writer) error {

	// Get the reflect value to work with arrays/slices
	rv := reflect.ValueOf(value)

	// Handle both arrays and slices
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return fmt.Errorf("expected array or slice, got %v", rv.Kind())
	}

	var length = rv.Len()

	// Determine the element type from the array's declared type
	elemType := rv.Type().Elem()
	elementType, err := getTypeForReflectType(elemType)
	if err != nil {
		return fmt.Errorf("unsupported array element type: %w", err)
	}

	// Write the array header with length (similar to string format)
	// Special case arrays that are less than 8 elements in length
	if length <= 0x07 {
		w.Write([]byte{types.CreateHeader(types.Array, byte(length))})
	} else {
		w.Write([]byte{types.CreateHeader(types.Array, 0x08)})
		SerializeUint64(uint64(length), w)
	}

	// Write the element type
	w.Write([]byte{byte(elementType)})

	// Serialize each element with their normal headers
	for i := range length {
		element := rv.Index(i).Interface()
		if err := Serialize(element, w); err != nil {
			return fmt.Errorf("failed to serialize array element %d: %w", i, err)
		}
	}

	return nil
}

func DeserializeArray(data []byte, out interface{}) ([]byte, error) {
	if len(data) == 0 {
		return data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	remaining := data[1:]

	var headerType = types.TypeFromHeader(header)
	if headerType != types.Array {
		return data, fmt.Errorf("expected Array type, got %v", types.TypeName(headerType))
	}

	var length = uint64(types.ValueFromHeader(header))

	// If the 4th bit of the nibble is not set then this is a special cased array whose length fits in the length nibble
	// Otherwise get the array length from integer in the next data type
	var err error = nil
	if length&0x08 != 0 {
		length, remaining, err = DeserializeUint64(remaining)
		if err != nil {
			return remaining, fmt.Errorf("failed to deserialize array length: %w", err)
		}
	}

	// Read the element type
	if len(remaining) == 0 {
		return remaining, fmt.Errorf("no data for element type")
	}
	elementType := types.Types(remaining[0])
	remaining = remaining[1:]

	// Validate element type (optional - could be used for type checking)
	_ = elementType // Currently unused, but reserved for future type validation

	// Get the output value and validate it's a pointer to slice or array
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return remaining, fmt.Errorf("output must be a pointer")
	}

	outElem := outValue.Elem()
	if outElem.Kind() != reflect.Slice && outElem.Kind() != reflect.Array {
		return remaining, fmt.Errorf("output must be a pointer to slice or array")
	}

	// For slices, create a new slice of the appropriate length
	if outElem.Kind() == reflect.Slice {
		sliceType := outElem.Type()
		newSlice := reflect.MakeSlice(sliceType, int(length), int(length))
		outElem.Set(newSlice)
	}

	// Deserialize each element using the generic deserializer
	for i := 0; i < int(length); i++ {
		if i >= outElem.Len() {
			return remaining, fmt.Errorf("array index %d out of bounds (length %d)", i, outElem.Len())
		}

		elemPtr := outElem.Index(i).Addr().Interface()

		// Deserialize the element using the generic deserializer
		newRemaining, err := Deserialize(bytes.NewReader(remaining), elemPtr)
		if err != nil {
			return remaining, fmt.Errorf("failed to deserialize array element %d: %w", i, err)
		}
		remaining = newRemaining
	}

	return remaining, nil
}

// Helper function to determine the Types enum value for a reflect.Type
func getTypeForReflectType(t reflect.Type) (types.Types, error) {
	switch t.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return types.UInt, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return types.SInt, nil
	case reflect.Float32, reflect.Float64:
		return types.Float, nil
	case reflect.Bool:
		return types.Boolean, nil
	case reflect.String:
		return types.String, nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return types.Buffer, nil
		}
		return 0, fmt.Errorf("nested slices not yet supported")
	case reflect.Struct:
		return 0, fmt.Errorf("structs in arrays not yet supported")
	default:
		return 0, fmt.Errorf("unsupported type: %v", t)
	}
}

// deserializeArrayElement deserializes a single array element of the specified type
// func deserializeArrayElement(data []byte, elementType types.Types, outValue reflect.Value) ([]byte, error) {
// 	switch elementType {
// 	case types.UNibble:
// 		value, remaining, err := DeserializeUNibble(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set UNibble value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.SNibble:
// 		value, remaining, err := DeserializeSNibble(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set SNibble value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.UInt:
// 		value, remaining, err := DeserializeUint64(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set UInt value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.SInt:
// 		value, remaining, err := DeserializeSint64(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set SInt value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.Float:
// 		value, remaining, err := DeserializeFloat(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set Float value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.Boolean:
// 		value, remaining, err := DeserializeBoolean(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set Boolean value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.String:
// 		value, remaining, err := DeserializeString(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set String value: %w", err)
// 		}
// 		return remaining, nil

// 	case types.Buffer:
// 		value, remaining, err := DeserializeBuffer(data)
// 		if err != nil {
// 			return remaining, err
// 		}
// 		if err := utils.SetValueWithConversion(outValue, value); err != nil {
// 			return remaining, fmt.Errorf("failed to set Buffer value: %w", err)
// 		}
// 		return remaining, nil

// 	default:
// 		return data, fmt.Errorf("unsupported array element type: %s", types.TypeName(elementType))
// 	}
// }
