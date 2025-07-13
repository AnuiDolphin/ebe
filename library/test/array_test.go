package test

import (
	"bytes"
	"ebe/serialize"
	"ebe/types"
	"ebe/utils"
	"fmt"
	"math"
	"reflect"
	"testing"
)

// TestArraySerialization tests basic array serialization/deserialization
func TestArraySerialization(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "int slice",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "uint32 slice",
			input:    []uint32{100, 200, 300},
			expected: []uint32{100, 200, 300},
		},
		{
			name:     "string slice",
			input:    []string{"hello", "world", "test"},
			expected: []string{"hello", "world", "test"},
		},
		{
			name:     "bool slice",
			input:    []bool{true, false, true},
			expected: []bool{true, false, true},
		},
		{
			name:     "float64 slice",
			input:    []float64{1.5, 2.7, 3.14159},
			expected: []float64{1.5, 2.7, 3.14159},
		},
		{
			name:     "int array",
			input:    [3]int{10, 20, 30},
			expected: [3]int{10, 20, 30},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			// Serialize
			if err := serialize.Serialize(tc.input, &buf); err != nil {
				t.Fatalf("Error serializing %s: %v", tc.name, err)
			}

			t.Logf("Serialized %s: %d bytes", tc.name, buf.Len())
			if testing.Verbose() {
				utils.PrintSerializedData(buf.Bytes())
			}

			// Create output variable of the same type as expected
			outputType := reflect.TypeOf(tc.expected)
			output := reflect.New(outputType).Interface()

			// Deserialize
			err := serialize.Deserialize(bytes.NewReader(buf.Bytes()), output)
			if err != nil {
				t.Fatalf("Error deserializing %s: %v", tc.name, err)
			}

			// Get the actual value from the pointer
			actualValue := reflect.ValueOf(output).Elem().Interface()

			// Compare using custom comparison for floats
			if !compareArrayValues(tc.expected, actualValue) {
				t.Errorf("Round trip mismatch for %s: expected %v, got %v", tc.name, tc.expected, actualValue)
			}

			t.Logf("%s: PASS", tc.name)
		})
	}
}

// TestEmptyArrays tests serialization of empty arrays and slices
func TestEmptyArrays(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "empty int slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "empty string slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "empty bool slice",
			input:    []bool{},
			expected: []bool{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := serialize.Serialize(tc.input, &buf); err != nil {
				t.Fatalf("Error serializing %s: %v", tc.name, err)
			}

			outputType := reflect.TypeOf(tc.expected)
			output := reflect.New(outputType).Interface()

			err := serialize.Deserialize(bytes.NewReader(buf.Bytes()), output)
			if err != nil {
				t.Fatalf("Error deserializing %s: %v", tc.name, err)
			}

			actualValue := reflect.ValueOf(output).Elem().Interface()
			if !reflect.DeepEqual(tc.expected, actualValue) {
				t.Errorf("Round trip mismatch for %s: expected %v, got %v", tc.name, tc.expected, actualValue)
			}
		})
	}
}

// TestLargeArrays tests arrays with more than 7 elements (requiring length encoding)
func TestLargeArrays(t *testing.T) {
	// Create arrays larger than 7 elements to test length encoding
	largeIntSlice := make([]int, 100)
	for i := range largeIntSlice {
		largeIntSlice[i] = i * 2
	}

	largeStringSlice := make([]string, 20)
	for i := range largeStringSlice {
		largeStringSlice[i] = fmt.Sprintf("item_%d", i)
	}

	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "large int slice (100 elements)",
			input:    largeIntSlice,
			expected: largeIntSlice,
		},
		{
			name:     "large string slice (20 elements)",
			input:    largeStringSlice,
			expected: largeStringSlice,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := serialize.Serialize(tc.input, &buf); err != nil {
				t.Fatalf("Error serializing %s: %v", tc.name, err)
			}

			t.Logf("Serialized %s: %d bytes", tc.name, buf.Len())

			outputType := reflect.TypeOf(tc.expected)
			output := reflect.New(outputType).Interface()

			err := serialize.Deserialize(bytes.NewReader(buf.Bytes()), output)
			if err != nil {
				t.Fatalf("Error deserializing %s: %v", tc.name, err)
			}

			actualValue := reflect.ValueOf(output).Elem().Interface()
			if !reflect.DeepEqual(tc.expected, actualValue) {
				t.Errorf("Round trip mismatch for %s", tc.name)
			}
		})
	}
}

// TestArrayWithSpecialValues tests arrays containing edge case values
func TestArrayWithSpecialValues(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "int slice with zero and negative values",
			input:    []int{-100, 0, 100, -1, 1},
			expected: []int{-100, 0, 100, -1, 1},
		},
		{
			name:     "uint slice with boundary values",
			input:    []uint64{0, 255, 65535, 4294967295, 18446744073709551615},
			expected: []uint64{0, 255, 65535, 4294967295, 18446744073709551615},
		},
		{
			name:     "string slice with special strings",
			input:    []string{"", "hello", "Hello 世界! 🌍🚀", string(make([]byte, 300))},
			expected: []string{"", "hello", "Hello 世界! 🌍🚀", string(make([]byte, 300))},
		},
		{
			name:     "float slice with special values",
			input:    []float64{0.0, math.Copysign(0, -1), 1.5, -1.5, 3.141592653589793},
			expected: []float64{0.0, math.Copysign(0, -1), 1.5, -1.5, 3.141592653589793},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			if err := serialize.Serialize(tc.input, &buf); err != nil {
				t.Fatalf("Error serializing %s: %v", tc.name, err)
			}

			outputType := reflect.TypeOf(tc.expected)
			output := reflect.New(outputType).Interface()

			err := serialize.Deserialize(bytes.NewReader(buf.Bytes()), output)
			if err != nil {
				t.Fatalf("Error deserializing %s: %v", tc.name, err)
			}

			actualValue := reflect.ValueOf(output).Elem().Interface()
			if !compareArrayValues(tc.expected, actualValue) {
				t.Errorf("Round trip mismatch for %s: expected %v, got %v", tc.name, tc.expected, actualValue)
			}
		})
	}
}

// TestArraySerializationErrors tests error conditions during serialization
func TestArraySerializationErrors(t *testing.T) {
	t.Run("unsupported element type", func(t *testing.T) {
		// Create a slice with unsupported element type (map)
		unsupportedSlice := []map[string]int{
			{"key": 1},
		}

		var buf bytes.Buffer
		err := serialize.Serialize(unsupportedSlice, &buf)
		if err == nil {
			t.Error("Expected error when serializing slice with unsupported element type, got nil")
		}
	})

}

// TestArrayDeserializationErrors tests error conditions during deserialization
func TestArrayDeserializationErrors(t *testing.T) {
	t.Run("corrupted array header", func(t *testing.T) {
		// Create corrupted data with wrong type header
		data := []byte{0x60, 0x02, 0x03} // Wrong type (6 = Boolean, should be 9 = Array)
		var out []int

		err := serialize.Deserialize(bytes.NewReader(data), &out)
		if err == nil {
			t.Error("Expected error when deserializing corrupted array header, got nil")
		}
	})

	t.Run("incomplete data", func(t *testing.T) {
		// Create incomplete array data (missing elements)
		data := []byte{0x92, 0x02} // Array header (9 = Array, 2 = length), element type (2 = SInt), but no elements
		var out []int

		err := serialize.Deserialize(bytes.NewReader(data), &out)
		if err == nil {
			t.Error("Expected error when deserializing incomplete array data, got nil")
		}
	})

	t.Run("non-pointer output", func(t *testing.T) {
		data := []byte{0x91, 0x02, 0x01} // Valid array data
		var out []int

		err := serialize.Deserialize(bytes.NewReader(data), out) // Not a pointer
		if err == nil {
			t.Error("Expected error when deserializing to non-pointer, got nil")
		}
	})

	t.Run("output not array or slice", func(t *testing.T) {
		data := []byte{0x91, 0x02, 0x01} // Valid array data
		var out int

		err := serialize.Deserialize(bytes.NewReader(data), &out) // Pointer to int, not array/slice
		if err == nil {
			t.Error("Expected error when deserializing to non-array/slice type, got nil")
		}
	})
}

// TestArrayWireFormat tests the specific wire format of arrays
func TestArrayWireFormat(t *testing.T) {
	t.Run("small array format", func(t *testing.T) {
		// Test array with ≤7 elements (length in header nibble)
		input := []int{1, 2, 3}
		var buf bytes.Buffer

		if err := serialize.Serialize(input, &buf); err != nil {
			t.Fatalf("Error serializing small array: %v", err)
		}

		data := buf.Bytes()
		if len(data) < 2 {
			t.Fatalf("Expected at least 2 bytes for small array, got %d", len(data))
		}

		// Check header: type=9 (Array), value=3 (length)
		expectedHeader := types.CreateHeader(types.Array, 3)
		if data[0] != expectedHeader {
			t.Errorf("Expected header %02x, got %02x", expectedHeader, data[0])
		}

		// Check element type: should be SInt (2)
		if data[1] != byte(types.SInt) {
			t.Errorf("Expected element type %d (SInt), got %d", types.SInt, data[1])
		}
	})

	t.Run("large array format", func(t *testing.T) {
		// Test array with >7 elements (length as separate UInt)
		input := make([]int, 10)
		for i := range input {
			input[i] = i
		}

		var buf bytes.Buffer
		if err := serialize.Serialize(input, &buf); err != nil {
			t.Fatalf("Error serializing large array: %v", err)
		}

		data := buf.Bytes()
		if len(data) < 3 {
			t.Fatalf("Expected at least 3 bytes for large array, got %d", len(data))
		}

		// Check header: type=9 (Array), value=8 (indicates length follows)
		expectedHeader := types.CreateHeader(types.Array, 0x08)
		if data[0] != expectedHeader {
			t.Errorf("Expected header %02x, got %02x", expectedHeader, data[0])
		}

		// Length should follow as UNibble (10 = 0x0a, since it's ≤ 15)
		// UNibble header: type=0, value=10
		expectedLengthHeader := types.CreateHeader(types.UNibble, 10)
		if data[1] != expectedLengthHeader {
			t.Errorf("Expected length header %02x, got %02x", expectedLengthHeader, data[1])
		}

		// Element type should follow immediately (no separate length value byte for UNibble)
		if data[2] != byte(types.SInt) {
			t.Errorf("Expected element type %d (SInt), got %d", types.SInt, data[2])
		}
	})
}

// Helper function to compare array values with float tolerance
func compareArrayValues(expected, actual interface{}) bool {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	if expectedValue.Type() != actualValue.Type() {
		return false
	}

	if expectedValue.Len() != actualValue.Len() {
		return false
	}

	// Compare each element
	for i := 0; i < expectedValue.Len(); i++ {
		expectedElem := expectedValue.Index(i).Interface()
		actualElem := actualValue.Index(i).Interface()

		if !utils.CompareValue(expectedElem, actualElem) {
			return false
		}
	}

	return true
}
