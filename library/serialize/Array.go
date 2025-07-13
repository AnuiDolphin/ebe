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

	// Write the array header
	if err := writeArrayHeader(w, length, elementType); err != nil {
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

	// Write the array header
	if err := writeArrayHeader(w, length, elementType); err != nil {
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

// Fast path serialization for unsigned integer arrays - avoids reflection overhead
func serializeUintArray(arr interface{}, w io.Writer) error {
	var length int
	var elementType types.Types

	// Handle different unsigned integer slice types
	switch v := arr.(type) {
	case []uint:
		length = len(v)
		elementType = types.UInt
	case []uint32:
		length = len(v)
		elementType = types.UInt
	case []uint64:
		length = len(v)
		elementType = types.UInt
	case []uint16:
		length = len(v)
		elementType = types.UInt
	default:
		return fmt.Errorf("unsupported unsigned integer array type: %T", arr)
	}

	// Write the array header
	if err := writeArrayHeader(w, length, elementType); err != nil {
		return err
	}

	// Serialize each element directly without reflection
	switch v := arr.(type) {
	case []uint:
		for _, elem := range v {
			if err := serializeUint(uint64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize uint element: %w", err)
			}
		}
	case []uint32:
		for _, elem := range v {
			if err := serializeUint(uint64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize uint32 element: %w", err)
			}
		}
	case []uint64:
		for _, elem := range v {
			if err := serializeUint(elem, w); err != nil {
				return fmt.Errorf("failed to serialize uint64 element: %w", err)
			}
		}
	case []uint16:
		for _, elem := range v {
			if err := serializeUint(uint64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize uint16 element: %w", err)
			}
		}
	}

	return nil
}

// Fast path serialization for float arrays - avoids reflection overhead
func serializeFloatArray(arr interface{}, w io.Writer) error {
	var length int
	var elementType types.Types

	// Handle different float slice types
	switch v := arr.(type) {
	case []float32:
		length = len(v)
		elementType = types.Float
	case []float64:
		length = len(v)
		elementType = types.Float
	default:
		return fmt.Errorf("unsupported float array type: %T", arr)
	}

	// Write the array header
	if err := writeArrayHeader(w, length, elementType); err != nil {
		return err
	}

	// Serialize each element directly without reflection
	switch v := arr.(type) {
	case []float32:
		for _, elem := range v {
			if err := serializeFloat(float64(elem), w); err != nil {
				return fmt.Errorf("failed to serialize float32 element: %w", err)
			}
		}
	case []float64:
		for _, elem := range v {
			if err := serializeFloat(elem, w); err != nil {
				return fmt.Errorf("failed to serialize float64 element: %w", err)
			}
		}
	}

	return nil
}

// Fast path serialization for string arrays - avoids reflection overhead
func serializeStringArray(arr []string, w io.Writer) error {
	length := len(arr)
	elementType := types.String

	// Write the array header
	if err := writeArrayHeader(w, length, elementType); err != nil {
		return err
	}

	// Serialize each string element directly without reflection
	for _, elem := range arr {
		if err := serializeString(elem, w); err != nil {
			return fmt.Errorf("failed to serialize string element: %w", err)
		}
	}

	return nil
}

func deserializeArray(r io.Reader, out interface{}) error {
	// Direct deserialization based on output type
	switch out.(type) {
	case *[]int, *[]int32, *[]int64, *[]int8, *[]int16:
		return deserializeIntArray(r, out)
	case *[]uint, *[]uint32, *[]uint64, *[]uint16:
		return deserializeUintArray(r, out)
	case *[]float32, *[]float64:
		return deserializeFloatArray(r, out)
	case *[]string:
		return deserializeStringArray(r, out)
	default:
		// Use generic reflection-based deserialization for unsupported types
		return deserializeArrayGeneric(r, out)
	}
}

// deserializeIntArray performs deserialization for integer arrays
func deserializeIntArray(r io.Reader, out interface{}) error {
	// Read array header and length
	length, elementType, err := readArrayHeader(r)
	if err != nil {
		return err
	}

	// Verify element type is SInt
	if elementType != types.SInt {
		return fmt.Errorf("expected SInt element type, got %v", types.TypeName(elementType))
	}

	// Type switch to handle different integer slice types
	switch ptr := out.(type) {
	case *[]int:
		*ptr = make([]int, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeInt64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize int element %d: %w", i, err)
			}
			(*ptr)[i] = int(elem)
		}
	case *[]int32:
		*ptr = make([]int32, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeInt64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize int32 element %d: %w", i, err)
			}
			(*ptr)[i] = int32(elem)
		}
	case *[]int64:
		*ptr = make([]int64, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeInt64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize int64 element %d: %w", i, err)
			}
			(*ptr)[i] = elem
		}
	case *[]int8:
		*ptr = make([]int8, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeInt64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize int8 element %d: %w", i, err)
			}
			(*ptr)[i] = int8(elem)
		}
	case *[]int16:
		*ptr = make([]int16, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeInt64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize int16 element %d: %w", i, err)
			}
			(*ptr)[i] = int16(elem)
		}
	default:
		return fmt.Errorf("unsupported integer array type: %T", out)
	}

	return nil
}

// deserializeUintArray performs deserialization for unsigned integer arrays
func deserializeUintArray(r io.Reader, out interface{}) error {
	// Read array header and length
	length, elementType, err := readArrayHeader(r)
	if err != nil {
		return err
	}

	// Verify element type is UInt
	if elementType != types.UInt {
		return fmt.Errorf("expected UInt element type, got %v", types.TypeName(elementType))
	}

	// Type switch to handle different unsigned integer slice types
	switch ptr := out.(type) {
	case *[]uint:
		*ptr = make([]uint, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeUint64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize uint element %d: %w", i, err)
			}
			(*ptr)[i] = uint(elem)
		}
	case *[]uint32:
		*ptr = make([]uint32, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeUint64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize uint32 element %d: %w", i, err)
			}
			(*ptr)[i] = uint32(elem)
		}
	case *[]uint64:
		*ptr = make([]uint64, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeUint64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize uint64 element %d: %w", i, err)
			}
			(*ptr)[i] = elem
		}
	case *[]uint16:
		*ptr = make([]uint16, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeUint64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize uint16 element %d: %w", i, err)
			}
			(*ptr)[i] = uint16(elem)
		}
	default:
		return fmt.Errorf("unsupported unsigned integer array type: %T", out)
	}

	return nil
}

// deserializeFloatArray performs deserialization for float arrays
func deserializeFloatArray(r io.Reader, out interface{}) error {
	// Read array header and length
	length, elementType, err := readArrayHeader(r)
	if err != nil {
		return err
	}

	// Verify element type is Float
	if elementType != types.Float {
		return fmt.Errorf("expected Float element type, got %v", types.TypeName(elementType))
	}

	// Type switch to handle different float slice types
	switch ptr := out.(type) {
	case *[]float32:
		*ptr = make([]float32, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeFloat64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize float32 element %d: %w", i, err)
			}
			(*ptr)[i] = float32(elem)
		}
	case *[]float64:
		*ptr = make([]float64, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeFloat64(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize float64 element %d: %w", i, err)
			}
			(*ptr)[i] = elem
		}
	default:
		return fmt.Errorf("unsupported float array type: %T", out)
	}

	return nil
}

// deserializeStringArray performs deserialization for string arrays
func deserializeStringArray(r io.Reader, out interface{}) error {
	// Read array header and length
	length, elementType, err := readArrayHeader(r)
	if err != nil {
		return err
	}

	// Verify element type is String
	if elementType != types.String {
		return fmt.Errorf("expected String element type, got %v", types.TypeName(elementType))
	}

	// Type switch to handle string slice
	switch ptr := out.(type) {
	case *[]string:
		*ptr = make([]string, length)
		for i := 0; i < int(length); i++ {
			elem, err := DeserializeString(r)
			if err != nil {
				return fmt.Errorf("failed to deserialize string element %d: %w", i, err)
			}
			(*ptr)[i] = elem
		}
	default:
		return fmt.Errorf("unsupported string array type: %T", out)
	}

	return nil
}

// deserializeArrayGeneric is the original reflection-based array deserialization
func deserializeArrayGeneric(r io.Reader, out interface{}) error {

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

// readArrayHeader reads and parses the array header, returning length and element type
func readArrayHeader(r io.Reader) (uint64, types.Types, error) {
	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read array header: %w", err)
	}

	if headerType != types.Array {
		return 0, 0, fmt.Errorf("expected Array type, got %v", types.TypeName(headerType))
	}

	length := uint64(headerValue)

	// If the 4th bit of the nibble is set, read the actual length from the next uint
	if length&0x08 != 0 {
		arrayLength, err := deserializeUint(r)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to deserialize array length: %w", err)
		}
		length = arrayLength
	}

	// Read the element type
	elementTypeByte, err := utils.ReadByte(r)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read element type: %w", err)
	}
	elementType := types.Types(elementTypeByte)

	return length, elementType, nil
}

// writeArrayHeader writes the array header with length and element type
func writeArrayHeader(w io.Writer, length int, elementType types.Types) error {
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

	return nil
}
