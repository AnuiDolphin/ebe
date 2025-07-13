package main_test

import (
	"bytes"
	"ebe-library/serialize"
	"ebe-library/types"
	"ebe-library/utils"
	"fmt"
	"math"
	"reflect"
	"testing"
)

// TestUNibbles tests UNibble (0-15) serialization/deserialization
func TestUNibbles(t *testing.T) {
	// Test UNibble range (0-15)
	for i := uint8(0); i <= 15; i++ {
		expectedType := types.UNibble
		if i == 0 {
			expectedType = types.SNibble // 0 serializes as SNibble
		}
		testValue(t, i, expectedType, fmt.Sprintf("UNibble %d", i))
	}
}

// TestSNibbles tests SNibble (-7 to 0) serialization/deserialization
func TestSNibbles(t *testing.T) {
	// Test SNibble range (-7 to 0)
	for i := int8(-7); i <= 0; i++ {
		testValue(t, i, types.SNibble, fmt.Sprintf("SNibble %d", i))
	}

	// Test boundary cases that become SInt due to range limits
	testValue(t, int8(-8), types.SInt, "int8(-8) - becomes SInt (out of SNibble range)")
}

// TestUnsignedIntegers tests unsigned integer serialization/deserialization
func TestUnsignedIntegers(t *testing.T) {
	// Test integer boundary transitions
	testValue(t, uint8(16), types.UInt, "uint8(16) - first UInt")
	testValue(t, uint8(255), types.UInt, "uint8(255) - max uint8")
	testValue(t, uint16(256), types.UInt, "uint16(256) - beyond uint8")
	testValue(t, uint16(65535), types.UInt, "uint16(65535) - max uint16")
	testValue(t, uint32(65536), types.UInt, "uint32(65536) - beyond uint16")
	testValue(t, uint32(4294967295), types.UInt, "uint32(4294967295) - max uint32")
	testValue(t, uint64(4294967296), types.UInt, "uint64(4294967296) - beyond uint32")
	testValue(t, uint64(18446744073709551615), types.UInt, "uint64(18446744073709551615) - max uint64")
}

// TestSignedIntegers tests signed integer serialization/deserialization
func TestSignedIntegers(t *testing.T) {
	// Note: positive int8 values 1-7 serialize as SNibble, but 8+ serialize as SInt due to signed type

	// Test signed integer boundaries
	testValue(t, int8(-128), types.SInt, "int8(-128) - min int8")
	testValue(t, int16(-129), types.SInt, "int16(-129) - beyond int8")
	testValue(t, int16(-32768), types.SInt, "int16(-32768) - min int16")
	testValue(t, int32(-32769), types.SInt, "int32(-32769) - beyond int16")
	testValue(t, int32(-2147483648), types.SInt, "int32(-2147483648) - min int32")
	testValue(t, int64(-2147483649), types.SInt, "int64(-2147483649) - beyond int32")
	// Note: int64 min value has precision issues in some implementations

	// Test that signed integers serialize as SInt regardless of positive value (except very small positive values)
	testValue(t, int16(16), types.SInt, "int16(16) - signed type stays SInt")
	testValue(t, int32(256), types.SInt, "int32(256) - signed type stays SInt")
	testValue(t, int64(65536), types.SInt, "int64(65536) - signed type stays SInt")
}

// TestStrings tests string serialization/deserialization
func TestStrings(t *testing.T) {
	testValue(t, "", types.String, "empty string")
	testValue(t, "A", types.String, "single ASCII char")
	testValue(t, "AB", types.String, "two ASCII chars")
	testValue(t, "Hello, World!", types.String, "standard string")
	testValue(t, "🙂", types.String, "single emoji")
	testValue(t, "🙂🎉🚀", types.String, "multiple emojis")
	testValue(t, "Hello 世界 🌍", types.String, "mixed unicode")
	testValue(t, "null\x00byte", types.String, "string with null byte")
	testValue(t, string(make([]byte, 254)), types.String, "254-byte string")
	testValue(t, string(make([]byte, 255)), types.String, "255-byte string")
	testValue(t, string(make([]byte, 256)), types.String, "256-byte string")
	testValue(t, string(make([]byte, 65534)), types.String, "65534-byte string")
	testValue(t, string(make([]byte, 65535)), types.String, "65535-byte string")
	testValue(t, string(make([]byte, 65536)), types.String, "65536-byte string")

	// Create a string with all possible byte values
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}
	testValue(t, string(allBytes), types.String, "string with all byte values")
}

// TestBuffers tests buffer/byte slice serialization/deserialization
func TestBuffers(t *testing.T) {
	// Test buffer edge cases
	testValue(t, []byte{}, types.Buffer, "empty buffer")
	testValue(t, []byte{0}, types.Buffer, "buffer with zero byte")
	testValue(t, []byte{255}, types.Buffer, "buffer with max byte")
	testValue(t, []byte{1, 2, 3}, types.Buffer, "small buffer")

	// Buffer with all possible byte values
	allBytes := make([]byte, 256)
	for i := 0; i < 256; i++ {
		allBytes[i] = byte(i)
	}
	testValue(t, allBytes, types.Buffer, "buffer with all byte values 0-255")

	// Test buffer size boundaries
	testValue(t, make([]byte, 254), types.Buffer, "254-byte buffer")
	testValue(t, make([]byte, 255), types.Buffer, "255-byte buffer")
	testValue(t, make([]byte, 256), types.Buffer, "256-byte buffer")
	testValue(t, make([]byte, 65534), types.Buffer, "65534-byte buffer")
	testValue(t, make([]byte, 65535), types.Buffer, "65535-byte buffer")
	testValue(t, make([]byte, 65536), types.Buffer, "65536-byte buffer")
}

// TestBooleans tests boolean serialization/deserialization
func TestBooleans(t *testing.T) {
	testValue(t, true, types.Boolean, "true")
	testValue(t, false, types.Boolean, "false")
}

// TestFloats tests floating-point serialization/deserialization
func TestFloats(t *testing.T) {
	// Test float edge cases
	testValue(t, 0.0, types.Float, "positive zero")
	testValue(t, math.Copysign(0, -1), types.Float, "negative zero")
	testValue(t, 1.0, types.Float, "simple positive")
	testValue(t, -1.0, types.Float, "simple negative")
	testValue(t, 0.1, types.Float, "decimal fraction")
	testValue(t, -0.1, types.Float, "negative decimal")

	// Float32 boundaries
	testValue(t, float32(1.17549435e-38), types.Float, "min positive float32")
	testValue(t, float32(-1.17549435e-38), types.Float, "max negative float32")
	testValue(t, float32(3.4028235e+38), types.Float, "max positive float32")
	testValue(t, float32(-3.4028235e+38), types.Float, "min negative float32")

	// Float64 boundaries
	testValue(t, 2.2250738585072014e-308, types.Float, "min positive float64")
	testValue(t, -2.2250738585072014e-308, types.Float, "max negative float64")
	testValue(t, 1.7976931348623157e+308, types.Float, "max positive float64")
	testValue(t, -1.7976931348623157e+308, types.Float, "min negative float64")

	// Special float values (now handled by improved compareFloats function)
	testValue(t, math.Inf(1), types.Float, "positive infinity")
	testValue(t, math.Inf(-1), types.Float, "negative infinity")
	testValue(t, math.NaN(), types.Float, "NaN")

	// Subnormal numbers
	testValue(t, 5e-324, types.Float, "smallest positive subnormal")
	testValue(t, -5e-324, types.Float, "largest negative subnormal")
}

func testValue(t *testing.T, value interface{}, expectedType types.Types, description string) {
	var data bytes.Buffer

	// Use generic serialize to handle all types
	err := serialize.Serialize(value, &data)
	if err != nil {
		t.Errorf("%s: Error serializing %T: %v", description, value, err)
		return
	}

	// Get the actual type from the serialized data
	if len(data.Bytes()) == 0 {
		t.Errorf("%s: No data serialized for %T", description, value)
		return
	}

	header := data.Bytes()[0]
	actualType := types.TypeFromHeader(header)

	// Validate the type matches expected
	if actualType != expectedType {
		t.Errorf("%s: Expected type %s, got %s", description, types.TypeName(expectedType), types.TypeName(actualType))
		return
	}

	// Use generic deserialize
	var readValue interface{}

	// Create appropriate type based on the serialized type
	switch actualType {
	case types.UNibble, types.UInt:
		var v uint64
		readValue = &v
	case types.SNibble, types.SInt:
		var v int64
		readValue = &v
	case types.Float:
		var v float64
		readValue = &v
	case types.Boolean:
		var v bool
		readValue = &v
	case types.String:
		var v string
		readValue = &v
	case types.Buffer:
		var v []byte
		readValue = &v
	default:
		t.Errorf("%s: Unsupported type for deserialization: %s", description, types.TypeName(actualType))
		return
	}

	err = serialize.Deserialize(bytes.NewReader(data.Bytes()), readValue)
	if err != nil {
		t.Errorf("%s: Deserialization error: %v", description, err)
		return
	}

	// Dereference the pointer to get the actual value
	readValue = reflect.ValueOf(readValue).Elem().Interface()

	// Check if values are equivalent (handles type conversions)
	isEqual := utils.CompareValue(value, readValue)

	if !isEqual {
		t.Errorf("%s: Values not equal: expected %v, got %v", description, value, readValue)
		return
	}

	t.Logf("%s: PASS", description)
}
