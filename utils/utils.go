package utils

import (
	"bytes"
	"ebe/types"
	"fmt"
	"io"
	"math"
	"reflect"
)

func Abs(value int64) uint64 {
	if value < 0 {
		return uint64(-value)
	} else {
		return uint64(value)
	}
}

// Helper function to check if two values are equivalent, handling type conversions
func CompareValue(a, b interface{}) bool {

	// If either value is nil, they are not equivalent
	if a == nil || b == nil {
		return false
	}

	// Special case: compare []byte with *bytes.Buffer
	if aBytes, ok := a.([]byte); ok {
		if bBuffer, ok := b.(*bytes.Buffer); ok {
			bBytes := bBuffer.Bytes()
			// Handle nil vs empty slice comparison
			if len(aBytes) == 0 && len(bBytes) == 0 {
				return true
			}
			return bytes.Equal(aBytes, bBytes)
		}
	}
	if aBuffer, ok := a.(*bytes.Buffer); ok {
		if bBytes, ok := b.([]byte); ok {
			aBytes := aBuffer.Bytes()
			// Handle nil vs empty slice comparison
			if len(aBytes) == 0 && len(bBytes) == 0 {
				return true
			}
			return bytes.Equal(aBytes, bBytes)
		}
	}

	// Try to convert b to the type of a
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)
	aType := aValue.Type()
	bType := bValue.Type()

	// If types are already the same, use special handling for floats or DeepEqual
	if aType == bType {
		// Special case: handle nil vs empty slice for []byte
		if aType == reflect.TypeOf([]byte{}) {
			aBytes := a.([]byte)
			bBytes := b.([]byte)
			if len(aBytes) == 0 && len(bBytes) == 0 {
				return true
			}
		}
		// Special handling for floating-point numbers
		if aType.Kind() == reflect.Float32 || aType.Kind() == reflect.Float64 {
			return compareFloats(a, b)
		}
		return reflect.DeepEqual(a, b)
	}

	// Try to convert b to the type of a
	if bType.ConvertibleTo(aType) {
		convertedB := bValue.Convert(aType).Interface()
		// Special handling for floating-point numbers after conversion
		if aType.Kind() == reflect.Float32 || aType.Kind() == reflect.Float64 {
			return compareFloats(a, convertedB)
		}
		return reflect.DeepEqual(a, convertedB)
	}

	// If conversion is not possible, fall back to reflect.DeepEqual
	return reflect.DeepEqual(a, b)
}

// Helper function to compare floating-point numbers with tolerance and special value handling
func compareFloats(a, b interface{}) bool {
	const tolerance = 1e-6 // More generous tolerance for float32 precision

	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	var aFloat, bFloat float64

	switch aVal.Kind() {
	case reflect.Float32:
		aFloat = float64(aVal.Float())
	case reflect.Float64:
		aFloat = aVal.Float()
	default:
		return false
	}

	switch bVal.Kind() {
	case reflect.Float32:
		bFloat = float64(bVal.Float())
	case reflect.Float64:
		bFloat = bVal.Float()
	default:
		return false
	}

	// Handle special float values first
	// NaN comparison - both must be NaN
	if math.IsNaN(aFloat) && math.IsNaN(bFloat) {
		return true
	}
	if math.IsNaN(aFloat) || math.IsNaN(bFloat) {
		return false // Only one is NaN
	}

	// Positive infinity comparison
	if math.IsInf(aFloat, 1) && math.IsInf(bFloat, 1) {
		return true
	}
	if math.IsInf(aFloat, 1) || math.IsInf(bFloat, 1) {
		return false // Only one is positive infinity
	}

	// Negative infinity comparison
	if math.IsInf(aFloat, -1) && math.IsInf(bFloat, -1) {
		return true
	}
	if math.IsInf(aFloat, -1) || math.IsInf(bFloat, -1) {
		return false // Only one is negative infinity
	}

	// Regular float comparison with tolerance
	diff := aFloat - bFloat
	if diff < 0 {
		diff = -diff
	}

	return diff < tolerance
}

// PrintSerializedData takes a byte array and prints out each serialized type and value
func PrintSerializedData(data []byte) {
	fmt.Printf("Parsing %d bytes of serialized data: [% x]\n", len(data), data)

	remaining := data
	valueIndex := 0

	for len(remaining) > 0 {
		fmt.Printf("[%d] ", valueIndex)

		// Parse the header
		header := remaining[0]
		headerType := types.TypeFromHeader(header)
		headerValue := types.ValueFromHeader(header)

		fmt.Printf("Type: %s, Value: %d", types.TypeName(headerType), headerValue)

		// Move past the header
		remaining = remaining[1:]

		// Estimate how many bytes this value uses based on type and header value
		var bytesToSkip int
		switch headerType {

		case types.UNibble, types.SNibble, types.Boolean:
			// These store their value in the header nibble
			bytesToSkip = 0

		case types.UInt, types.SInt:
			// These use the header value as the length
			bytesToSkip = int(headerValue)

		case types.Float:
			// Float uses header value as length (4 or 8 bytes typically)
			bytesToSkip = int(headerValue)

		case types.String, types.Buffer:

			if headerValue&0x08 != 0 {

				// Length stored as separate UInt - need to parse it
				if len(remaining) > 0 {

					lengthHeader := remaining[0]
					lengthBytes := int(types.ValueFromHeader(lengthHeader))
					if len(remaining) >= 1+lengthBytes {
						// Parse the actual length
						var actualLength uint64 = 0
						for i := 1; i <= lengthBytes; i++ {
							actualLength |= uint64(remaining[i]) << (8 * (lengthBytes - (i - 1) - 1))
						}
						bytesToSkip = 1 + lengthBytes + int(actualLength)
						fmt.Printf(", Length bytes: %d, Actual length: %d", lengthBytes, actualLength)
					} else {
						bytesToSkip = len(remaining) // Skip rest if malformed
					}
				} else {
					bytesToSkip = 0
				}
			} else {
				// Length stored in header
				bytesToSkip = int(headerValue)
			}

		case types.Array:
			// Array format: header, [length if > 7], element_type, then elements
			var arrayLength uint64 = uint64(headerValue)
			elementTypeBytes := 1 // Always 1 byte for element type

			if headerValue&0x08 != 0 {
				// Length stored as separate UInt - need to parse it
				if len(remaining) > 0 {
					lengthHeader := remaining[0]
					lengthBytes := int(types.ValueFromHeader(lengthHeader))
					if len(remaining) >= 1+lengthBytes {
						// Parse the actual length
						arrayLength = 0
						for i := 1; i <= lengthBytes; i++ {
							arrayLength |= uint64(remaining[i]) << (8 * (lengthBytes - (i - 1) - 1))
						}
						elementTypeBytes += 1 + lengthBytes
						fmt.Printf(", Length bytes: %d, Array length: %d", lengthBytes, arrayLength)
					}
				}
			}

			// For arrays, we just skip the element type byte and let the elements be parsed individually
			bytesToSkip = elementTypeBytes
			if len(remaining) >= elementTypeBytes {
				elementType := types.Types(remaining[elementTypeBytes-1])
				fmt.Printf(", Element type: %s", types.TypeName(elementType))
			}
		default:
			// Unknown type, skip 1 byte and continue
			bytesToSkip = 1
		}

		// Print the data bytes if any
		if bytesToSkip > 0 {
			if len(remaining) >= bytesToSkip {
				fmt.Printf(", Data (%d bytes): [% x]", bytesToSkip, remaining[:bytesToSkip])
				remaining = remaining[bytesToSkip:]
			} else {
				fmt.Printf(", Data (remaining %d bytes): [% x]", len(remaining), remaining)
				remaining = nil
			}
		}

		fmt.Println()
		valueIndex++

		// Safety check to prevent infinite loops
		if valueIndex > 100 {
			fmt.Println("Stopping after 100 values to prevent infinite loop")
			break
		}
	}

	fmt.Printf("Parsed %d values total\n", valueIndex)
}

// SetValueWithConversion sets a reflect.Value with type conversion support
func SetValueWithConversion(rhs reflect.Value, lhs interface{}) error {
	valueReflect := reflect.ValueOf(lhs)
	outType := rhs.Type()
	valueType := valueReflect.Type()

	// Special case: convert *bytes.Buffer to []byte
	if outType == reflect.TypeOf([]byte{}) && valueType == reflect.TypeOf(&bytes.Buffer{}) {
		buffer := lhs.(*bytes.Buffer)
		bufBytes := buffer.Bytes()
		// Ensure we always return a non-nil slice, even if empty
		if bufBytes == nil {
			bufBytes = []byte{}
		}
		rhs.Set(reflect.ValueOf(bufBytes))
		return nil
	}

	// Special case: convert []byte to *bytes.Buffer
	if outType == reflect.TypeOf(&bytes.Buffer{}) && valueType == reflect.TypeOf([]byte{}) {
		byteSlice := lhs.([]byte)
		buffer := bytes.NewBuffer(byteSlice)
		rhs.Set(reflect.ValueOf(buffer))
		return nil
	}

	// Direct assignment if types match
	if valueType.AssignableTo(outType) {
		rhs.Set(valueReflect)
		return nil
	}

	// Type conversion if possible
	if valueType.ConvertibleTo(outType) {
		rhs.Set(valueReflect.Convert(outType))
		return nil
	}

	return fmt.Errorf("cannot convert %v to %v", valueType, outType)
}

// WriteByte writes a single byte to the provided io.Writer.
func WriteByte(w io.Writer, b byte) error {
	buf := []byte{b}
	_, err := w.Write(buf)
	return err
}
