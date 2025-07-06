package main_test

import (
	"bytes"
	"ebe/serialize"
	"ebe/utils"
	"testing"
)

func TestMultipleValueSerialization(t *testing.T) {
	var data bytes.Buffer

	// Serialize multiple values into buffer
	values := []interface{}{
		uint64(0xffffffffffffffff),
		"The quick brown fox jumps over the lazy dog",
		int64(-0x7fffffffffffffff),
		"Hello, World!",
		true,
		3.141592653589793,
	}

	for i, value := range values {
		if err := serialize.Serialize(value, &data); err != nil {
			t.Fatalf("Error serializing value %d (%T): %v", i, value, err)
		}
	}

	t.Logf("Serialized %d values into %d bytes", len(values), data.Len())

	// Test manual deserialization of multiple values
	remaining := data.Bytes()
	var deserializedValues []interface{}

	// Deserialize each value manually
	var val1 uint64
	remaining, err := serialize.Deserialize(bytes.NewReader(remaining), &val1)
	if err != nil {
		t.Fatalf("Error deserializing value 1: %v", err)
	}
	deserializedValues = append(deserializedValues, val1)

	var val2 string
	remaining, err = serialize.Deserialize(bytes.NewReader(remaining), &val2)
	if err != nil {
		t.Fatalf("Error deserializing value 2: %v", err)
	}
	deserializedValues = append(deserializedValues, val2)

	var val3 int64
	remaining, err = serialize.Deserialize(bytes.NewReader(remaining), &val3)
	if err != nil {
		t.Fatalf("Error deserializing value 3: %v", err)
	}
	deserializedValues = append(deserializedValues, val3)

	var val4 string
	remaining, err = serialize.Deserialize(bytes.NewReader(remaining), &val4)
	if err != nil {
		t.Fatalf("Error deserializing value 4: %v", err)
	}
	deserializedValues = append(deserializedValues, val4)

	var val5 bool
	remaining, err = serialize.Deserialize(bytes.NewReader(remaining), &val5)
	if err != nil {
		t.Fatalf("Error deserializing value 5: %v", err)
	}
	deserializedValues = append(deserializedValues, val5)

	var val6 float64
	remaining, err = serialize.Deserialize(bytes.NewReader(remaining), &val6)
	if err != nil {
		t.Fatalf("Error deserializing value 6: %v", err)
	}
	deserializedValues = append(deserializedValues, val6)

	// Verify all remaining data was consumed
	if len(remaining) != 0 {
		t.Errorf("Expected no remaining bytes after deserialization, got %d", len(remaining))
	}

	// Verify the values
	if len(values) != len(deserializedValues) {
		t.Fatalf("Value count mismatch: expected %d, got %d", len(values), len(deserializedValues))
	}

	for i, expected := range values {
		if i >= len(deserializedValues) {
			t.Errorf("Value %d: missing deserialized value", i)
			continue
		}

		actual := deserializedValues[i]
		if !utils.CompareValue(expected, actual) {
			t.Errorf("Value %d: expected %v (%T), got %v (%T)", i, expected, expected, actual, actual)
		} else {
			t.Logf("Value %d: PASS - %T: %v", i, actual, actual)
		}
	}
}
