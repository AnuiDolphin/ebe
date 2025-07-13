package main_test

import (
	"bytes"
	"ebe-library/serialize"
	"ebe-library/utils"
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
	reader := bytes.NewReader(data.Bytes())
	var deserializedValues []interface{}

	// Deserialize each value manually
	var val1 uint64
	err := serialize.Deserialize(reader, &val1)
	if err != nil {
		t.Fatalf("Error deserializing value 1: %v", err)
	}
	deserializedValues = append(deserializedValues, val1)

	var val2 string
	err = serialize.Deserialize(reader, &val2)
	if err != nil {
		t.Fatalf("Error deserializing value 2: %v", err)
	}
	deserializedValues = append(deserializedValues, val2)

	var val3 int64
	err = serialize.Deserialize(reader, &val3)
	if err != nil {
		t.Fatalf("Error deserializing value 3: %v", err)
	}
	deserializedValues = append(deserializedValues, val3)

	var val4 string
	err = serialize.Deserialize(reader, &val4)
	if err != nil {
		t.Fatalf("Error deserializing value 4: %v", err)
	}
	deserializedValues = append(deserializedValues, val4)

	var val5 bool
	err = serialize.Deserialize(reader, &val5)
	if err != nil {
		t.Fatalf("Error deserializing value 5: %v", err)
	}
	deserializedValues = append(deserializedValues, val5)

	var val6 float64
	err = serialize.Deserialize(reader, &val6)
	if err != nil {
		t.Fatalf("Error deserializing value 6: %v", err)
	}
	deserializedValues = append(deserializedValues, val6)

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
