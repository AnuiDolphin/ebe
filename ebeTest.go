package main

import (
	"bytes"
	"ebe/serialize"
	"ebe/types"
	"ebe/utils"
	"fmt"
	"math"
)

func main() {

	// Test multiple serialization
	fmt.Println("\n--- Testing Multiple serialization ---")
	TestMultipleDeserialize()

	// Test different integer types (Go will infer appropriate types)
	TestValue(uint8(0), types.SNibble)
	TestValue(uint8(1), types.UNibble)
	TestValue(uint8(0xff), types.UInt)
	TestValue(uint16(0xffff), types.UInt)
	TestValue(uint32(0xffffffff), types.UInt)
	TestValue(uint64(0xffffffffffffffff), types.UInt) // Must cast to uint64 - too large for signed int

	// Test some boundary values
	TestValue(uint8(0xff), types.UInt)
	TestValue(uint16(0xff+1), types.UInt)
	TestValue(uint16(0xffff), types.UInt)
	TestValue(uint32(0xffff+1), types.UInt)
	TestValue(uint32(0xffffffff), types.UInt)
	TestValue(uint64(0xffffffff+1), types.UInt)

	// Test different signed integer values
	TestValue(-1, types.SNibble)
	TestValue(-7, types.SNibble)
	TestValue(-0x7f, types.SInt)
	TestValue(-127, types.SInt)
	TestValue(-0x7fff, types.SInt)
	TestValue(-32767, types.SInt)
	TestValue(-0x7fffffff, types.SInt)
	TestValue(-2147483647, types.SInt)
	TestValue(-0x7fffffffffffffff, types.SInt)
	TestValue(-9223372036854775807, types.SInt)

	// Test some additional signed values
	TestValue(-128, types.SInt)        // min int8 range
	TestValue(-32768, types.SInt)      // min int16 range
	TestValue(-2147483648, types.SInt) // min int32 range
	TestValue(-0x7fffffffffffffff, types.SInt)

	// Test strings
	TestValue("", types.String)
	TestValue("A", types.String)
	TestValue("AB", types.String)
	TestValue("Hello", types.String)
	TestValue("The quick brown fox jumped over the lazy dog", types.String)
	TestValue("🙂", types.String)

	// Test buffers
	TestValue([]byte{}, types.Buffer)
	TestValue([]byte{1}, types.Buffer)
	TestValue([]byte{1, 2, 3}, types.Buffer)
	TestValue([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, types.Buffer)

	// Test booleans
	TestValue(true, types.Boolean)
	TestValue(false, types.Boolean)

	// Test floats
	TestValue(0.0, types.Float)
	TestValue(0.1, types.Float)
	TestValue(-1.1, types.Float)
	TestValue(float32(1.17549435e-38), types.Float) // Min float32
	TestValue(1.17549435e-39, types.Float)
	TestValue(float32(math.MaxFloat32), types.Float)
	TestValue(math.MaxFloat64, types.Float)
}

func TestValue(value interface{}, expectedType types.Types) {
	var data bytes.Buffer

	// Use generic serialize to handle all types
	err := serialize.Serialize(value, &data)
	if err != nil {
		fmt.Printf("Error serializing %T: %v\n", value, err)
		return
	}

	// Get the actual type from the serialized data
	if len(data.Bytes()) == 0 {
		fmt.Printf("Error: no data serialized for %T\n", value)
		return
	}

	header := data.Bytes()[0]
	actualType := types.TypeFromHeader(header)

	// Format output based on type
	switch v := value.(type) {
	case string, []byte, bool:
		fmt.Print("value = '", v, "', ")
	default:
		fmt.Print("value = ", v, ", ")
	}

	fmt.Print(types.HeaderString(data.Bytes()))

	// Validate the type matches expected
	if actualType != expectedType {
		fmt.Printf(", Error: expected type: %s, got: %s\n", types.TypeName(expectedType), types.TypeName(actualType))
		return
	}

	// Use generic deserialize
	readValue, _, err := serialize.Deserialize(data.Bytes())
	if err != nil {
		fmt.Printf(", Error: %v\n", err)
		return
	}

	// Format read output based on original type
	switch value.(type) {
	case string:
		fmt.Print(", Read:'", readValue, "'")
	case []byte:
		if buffer, ok := readValue.(*bytes.Buffer); ok {
			fmt.Print(", Read:", buffer.Bytes())
		} else {
			fmt.Print(", Read:", readValue)
		}
	default:
		fmt.Print(", Read:", readValue)
	}

	// Check if values are equivalent (handles type conversions)
	var isEqual bool
	switch originalValue := value.(type) {
	case []byte:
		if buffer, ok := readValue.(*bytes.Buffer); ok {
			isEqual = bytes.Equal(originalValue, buffer.Bytes())
		}
	default:
		isEqual = utils.CompareValue(value, readValue)
	}

	if !isEqual {
		fmt.Printf(", Error: values not equal: %v != %v\n", value, readValue)
		return
	}

	fmt.Println(", Pass")
}

func TestMultipleDeserialize() {
	var data bytes.Buffer

	// Serialize multiple values into buffer
	serialize.Serialize(uint64(0xffffffffffffffff), &data)
	serialize.Serialize("The quick brown fox jumps over the lazy dog", &data)
	serialize.Serialize(-0x7fffffffffffffff, &data)
	serialize.Serialize("Hello, World!", &data)
	serialize.Serialize(true, &data)
	serialize.Serialize(3.141592653589793, &data)

	// Test DeserializeAll
	values, err := serialize.DeserializeAll(data.Bytes())
	if err != nil {
		fmt.Println("Error deserializing all:", err)
		return
	}

	// Verify the values
	expectedValues := []interface{}{
		uint64(0xffffffffffffffff),
		"The quick brown fox jumps over the lazy dog",
		int64(-0x7fffffffffffffff),
		"Hello, World!",
		true,
		3.141592653589793,
	}

	for i, expected := range expectedValues {
		if i < len(values) && values[i] == expected {
			fmt.Printf("  [%d] Pass: %T, %v\n", i, values[i], values[i])
		} else {
			fmt.Printf("  [%d] *** FAIL *** (expected %v, got %v)\n", i, expected, values[i])
		}
	}
}
