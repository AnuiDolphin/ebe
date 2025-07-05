package utils

import (
	"bytes"
	"ebe/types"
	"fmt"
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

// convertToMatchingType converts value b to match the type of value a
func ConvertToMatchingType(a, b interface{}) (interface{}, error) {
	if a == nil {
		return nil, fmt.Errorf("cannot determine target type from nil value")
	}

	// Get the target type from a
	targetType := reflect.TypeOf(a)

	// If b is nil, return nil
	if b == nil {
		return nil, fmt.Errorf("cannot convert nil value")
	}

	// If types are already the same, return b as-is
	if reflect.TypeOf(b) == targetType {
		return b, nil
	}

	// Handle conversions based on target type
	switch targetType.Kind() {
	case reflect.Uint64:
		val, err := ToUint64(b)
		if err != nil {
			return nil, err
		}
		return val, nil

	case reflect.Uint32:
		val, err := ToUint64(b)
		if err != nil {
			return nil, err
		}
		if val > math.MaxUint32 {
			return nil, fmt.Errorf("value %d exceeds uint32 range", val)
		}
		return uint32(val), nil

	case reflect.Uint16:
		val, err := ToUint64(b)
		if err != nil {
			return nil, err
		}
		if val > math.MaxUint16 {
			return nil, fmt.Errorf("value %d exceeds uint16 range", val)
		}
		return uint16(val), nil

	case reflect.Uint8:
		val, err := ToUint64(b)
		if err != nil {
			return nil, err
		}
		if val > math.MaxUint8 {
			return nil, fmt.Errorf("value %d exceeds uint8 range", val)
		}
		return uint8(val), nil

	case reflect.Uint:
		val, err := ToUint64(b)
		if err != nil {
			return nil, err
		}
		return uint(val), nil

	case reflect.Int64:
		val, err := ToSint64(b)
		if err != nil {
			return nil, err
		}
		return val, nil

	case reflect.Int32:
		val, err := ToSint64(b)
		if err != nil {
			return nil, err
		}
		if val < math.MinInt32 || val > math.MaxInt32 {
			return nil, fmt.Errorf("value %d exceeds int32 range", val)
		}
		return int32(val), nil

	case reflect.Int16:
		val, err := ToSint64(b)
		if err != nil {
			return nil, err
		}
		if val < math.MinInt16 || val > math.MaxInt16 {
			return nil, fmt.Errorf("value %d exceeds int16 range", val)
		}
		return int16(val), nil

	case reflect.Int8:
		val, err := ToSint64(b)
		if err != nil {
			return nil, err
		}
		if val < math.MinInt8 || val > math.MaxInt8 {
			return nil, fmt.Errorf("value %d exceeds int8 range", val)
		}
		return int8(val), nil

	case reflect.Int:
		val, err := ToSint64(b)
		if err != nil {
			return nil, err
		}
		return int(val), nil

	case reflect.Float64:
		switch v := b.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case uint64, uint32, uint16, uint8, uint:
			uval, err := ToUint64(b)
			if err != nil {
				return nil, err
			}
			return float64(uval), nil
		case int64, int32, int16, int8, int:
			ival, err := ToSint64(b)
			if err != nil {
				return nil, err
			}
			return float64(ival), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float64", b)
		}

	case reflect.Float32:
		switch v := b.(type) {
		case float32:
			return v, nil
		case float64:
			if v > math.MaxFloat32 || v < -math.MaxFloat32 {
				return nil, fmt.Errorf("value %f exceeds float32 range", v)
			}
			return float32(v), nil
		case uint64, uint32, uint16, uint8, uint:
			uval, err := ToUint64(b)
			if err != nil {
				return nil, err
			}
			return float32(uval), nil
		case int64, int32, int16, int8, int:
			ival, err := ToSint64(b)
			if err != nil {
				return nil, err
			}
			return float32(ival), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float32", b)
		}

	case reflect.String:
		if str, ok := b.(string); ok {
			return str, nil
		}
		return nil, fmt.Errorf("cannot convert %T to string", b)

	case reflect.Bool:
		if boolean, ok := b.(bool); ok {
			return boolean, nil
		}
		return nil, fmt.Errorf("cannot convert %T to bool", b)

	case reflect.Slice:
		if targetType.Elem().Kind() == reflect.Uint8 { // []byte
			if slice, ok := b.([]byte); ok {
				return slice, nil
			}
			if buffer, ok := b.(*bytes.Buffer); ok {
				return buffer.Bytes(), nil
			}
		}
		return nil, fmt.Errorf("cannot convert %T to %v", b, targetType)

	case reflect.Ptr:
		if targetType == reflect.TypeOf((*bytes.Buffer)(nil)) {
			if buffer, ok := b.(*bytes.Buffer); ok {
				return buffer, nil
			}
			if slice, ok := b.([]byte); ok {
				buf := bytes.NewBuffer(slice)
				return buf, nil
			}
		}
		return nil, fmt.Errorf("cannot convert %T to %v", b, targetType)

	default:
		return nil, fmt.Errorf("unsupported target type: %v", targetType)
	}
}

// Helper function to convert interface{} to uint64
func ToUint64(v interface{}) (uint64, error) {
	switch val := v.(type) {
	case uint64:
		return val, nil
	case uint32:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint8:
		return uint64(val), nil
	case uint:
		return uint64(val), nil
	case int64:
		if val < 0 {
			return 0, fmt.Errorf("cannot convert negative value %d to uint64", val)
		}
		return uint64(val), nil
	case int32:
		if val < 0 {
			return 0, fmt.Errorf("cannot convert negative value %d to uint64", val)
		}
		return uint64(val), nil
	case int16:
		if val < 0 {
			return 0, fmt.Errorf("cannot convert negative value %d to uint64", val)
		}
		return uint64(val), nil
	case int8:
		if val < 0 {
			return 0, fmt.Errorf("cannot convert negative value %d to uint64", val)
		}
		return uint64(val), nil
	case int:
		if val < 0 {
			return 0, fmt.Errorf("cannot convert negative value %d to uint64", val)
		}
		return uint64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

// Helper function to convert interface{} to int64
func ToSint64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int64:
		return val, nil
	case int32:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int:
		return int64(val), nil
	case uint64:
		if val > math.MaxInt64 {
			return 0, fmt.Errorf("value %d exceeds int64 range", val)
		}
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint:
		if uint64(val) > math.MaxInt64 {
			return 0, fmt.Errorf("value %d exceeds int64 range", val)
		}
		return int64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

// Helper function to check if two values are equivalent, handling type conversions
func CompareValue(a, b interface{}) bool {

	// If either value is nil, they are not equivalent
	if a == nil || b == nil {
		return false
	}

	// Try to convert b to match the type of a
	converted, err := ConvertToMatchingType(a, b)
	if err != nil {
		// If conversion fails, fall back to reflect.DeepEqual
		return reflect.DeepEqual(a, b)
	}

	// Compare the original value with the converted value
	return reflect.DeepEqual(a, converted)
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

			if headerValue & 0x08 != 0 {

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
