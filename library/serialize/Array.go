package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
	"reflect"
)

// This function appends the serialized array to the existing writer
func serializeArray(rv reflect.Value, w io.Writer) error {

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
		if err := utils.WriteByte(w, types.CreateHeader(types.Array, byte(length))); err != nil {
			return err
		}
	} else {
		if err := utils.WriteByte(w, types.CreateHeader(types.Array, 0x08)); err != nil {
			return err
		}
		if err := serializeUint(uint64(length), w); err != nil {
			return err
		}
	}

	// Write the element type
	if err := utils.WriteByte(w, byte(elementType)); err != nil {
		return err
	}

	// Serialize each element with their normal headers
	for i := range length {
		element := rv.Index(i).Interface()
		if err := Serialize(element, w); err != nil {
			return fmt.Errorf("failed to serialize array element %d: %w", i, err)
		}
	}

	return nil
}

// Fast path serialization for integer arrays - avoids reflection overhead
func serializeIntArray(arr interface{}, w io.Writer) error {
	var length int
	var elementType types.Types

	// Handle different integer slice types
	switch v := arr.(type) {
	case []int:
		length = len(v)
		elementType = types.SInt
	case []int32:
		length = len(v)
		elementType = types.SInt
	case []int64:
		length = len(v)
		elementType = types.SInt
	case []int8:
		length = len(v)
		elementType = types.SInt
	case []int16:
		length = len(v)
		elementType = types.SInt
	default:
		return fmt.Errorf("unsupported integer array type: %T", arr)
	}

	// Write the array header with length
	if length <= 0x07 {
		if err := utils.WriteByte(w, types.CreateHeader(types.Array, byte(length))); err != nil {
			return err
		}
	} else {
		if err := utils.WriteByte(w, types.CreateHeader(types.Array, 0x08)); err != nil {
			return err
		}
		if err := serializeUint(uint64(length), w); err != nil {
			return err
		}
	}

	// Write the element type
	if err := utils.WriteByte(w, byte(elementType)); err != nil {
		return err
	}

	// Serialize each element directly without reflection
	switch v := arr.(type) {
	case []int:
		for _, elem := range v {
			if err := serializeSint(int64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize int element: %w", err)
			}
		}
	case []int32:
		for _, elem := range v {
			if err := serializeSint(int64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize int32 element: %w", err)
			}
		}
	case []int64:
		for _, elem := range v {
			if err := serializeSint(elem, w); err != nil {
				return fmt.Errorf("failed to serialize int64 element: %w", err)
			}
		}
	case []int8:
		for _, elem := range v {
			if err := serializeSint(int64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize int8 element: %w", err)
			}
		}
	case []int16:
		for _, elem := range v {
			if err := serializeSint(int64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize int16 element: %w", err)
			}
		}
	}

	return nil
}

func deserializeArray(r io.Reader, out interface{}) error {

	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return fmt.Errorf("failed to read array header: %w", err)
	}

	if headerType != types.Array {
		return fmt.Errorf("expected Array type, got %v", types.TypeName(headerType))
	}

	length := uint64(headerValue)

	// If the 4th bit of the nibble is not set then this is a special cased array whose length fits in the length nibble
	// Otherwise get the array length from integer in the next data type
	if length&0x08 != 0 {

		// Parse the length using deserializeUint64 directly from the reader
		arrayLength, err := deserializeUint(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize array length: %w", err)
		}
		length = arrayLength
	}

	// Read the element type
	elementTypeByte, err := utils.ReadByte(r)
	if err != nil {
		return fmt.Errorf("failed to read element type: %w", err)
	}
	elementType := types.Types(elementTypeByte)

	// Validate element type (optional - could be used for type checking)
	_ = elementType // Currently unused, but reserved for future type validation

	// Get the output value and validate it's a pointer to slice or array
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return fmt.Errorf("output must be a pointer")
	}

	outElem := outValue.Elem()
	if outElem.Kind() != reflect.Slice && outElem.Kind() != reflect.Array {
		return fmt.Errorf("output must be a pointer to slice or array")
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
			return fmt.Errorf("array index %d out of bounds (length %d)", i, outElem.Len())
		}

		elemPtr := outElem.Index(i).Addr().Interface()

		// Deserialize the element using the generic deserializer
		err = Deserialize(r, elemPtr)
		if err != nil {
			return fmt.Errorf("failed to deserialize array element %d: %w", i, err)
		}
	}

	return nil
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
