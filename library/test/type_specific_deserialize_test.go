package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

func TestDeserializeInt64(t *testing.T) {
	tests := []struct {
		name  string
		value int64
	}{
		{"zero", 0},
		{"positive small", 7},
		{"negative small", -5},
		{"positive large", 1000},
		{"negative large", -1000},
		{"max int64", 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize the value
			var buf bytes.Buffer
			err := serialize.Serialize(tt.value, &buf)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Deserialize using type-specific function
			result, err := serialize.DeserializeInt64(&buf)
			if err != nil {
				t.Fatalf("Type-specific deserialization failed: %v", err)
			}

			if result != tt.value {
				t.Errorf("Expected %d, got %d", tt.value, result)
			}
		})
	}
}

func TestDeserializeUint64(t *testing.T) {
	tests := []struct {
		name  string
		value uint64
	}{
		{"zero", 0},
		{"small", 15},
		{"medium", 1000},
		{"large", 1000000},
		{"max uint64", 18446744073709551615},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize the value
			var buf bytes.Buffer
			err := serialize.Serialize(tt.value, &buf)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Deserialize using type-specific function
			result, err := serialize.DeserializeUint64(&buf)
			if err != nil {
				t.Fatalf("Type-specific deserialization failed: %v", err)
			}

			if result != tt.value {
				t.Errorf("Expected %d, got %d", tt.value, result)
			}
		})
	}
}

func TestDeserializeFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"positive", 3.5},
		{"negative", -2.5},
		{"large", 1000000.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize the value
			var buf bytes.Buffer
			err := serialize.Serialize(tt.value, &buf)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Deserialize using type-specific function
			result, err := serialize.DeserializeFloat64(&buf)
			if err != nil {
				t.Fatalf("Type-specific deserialization failed: %v", err)
			}

			// Use proper floating point comparison
			if result != tt.value && !(result != result && tt.value != tt.value) { // Handle NaN case
				t.Errorf("Expected %f, got %f", tt.value, result)
			}
		})
	}
}

func TestDeserializeString(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"unicode", "café naïve résumé"},
		{"long", "This is a longer string to test serialization and deserialization of strings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize the value
			var buf bytes.Buffer
			err := serialize.Serialize(tt.value, &buf)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Deserialize using type-specific function
			result, err := serialize.DeserializeString(&buf)
			if err != nil {
				t.Fatalf("Type-specific deserialization failed: %v", err)
			}

			if result != tt.value {
				t.Errorf("Expected %q, got %q", tt.value, result)
			}
		})
	}
}

func TestDeserializeBytes(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
	}{
		{"empty", []byte{}},
		{"simple", []byte{1, 2, 3, 4, 5}},
		{"large", make([]byte, 1000)}, // Initialize with zeros
	}

	// Initialize the large test case with some data
	for i := range tests[2].value {
		tests[2].value[i] = byte(i % 256)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize the value
			var buf bytes.Buffer
			err := serialize.Serialize(tt.value, &buf)
			if err != nil {
				t.Fatalf("Serialization failed: %v", err)
			}

			// Deserialize using type-specific function
			result, err := serialize.DeserializeBytes(&buf)
			if err != nil {
				t.Fatalf("Type-specific deserialization failed: %v", err)
			}

			if len(result) != len(tt.value) {
				t.Errorf("Length mismatch: expected %d, got %d", len(tt.value), len(result))
				return
			}

			for i, b := range tt.value {
				if result[i] != b {
					t.Errorf("Byte mismatch at index %d: expected %d, got %d", i, b, result[i])
					return
				}
			}
		})
	}
}

// Test type conversion edge cases
func TestDeserializeInt64FromUint(t *testing.T) {
	// Test converting a uint value to int64
	var buf bytes.Buffer
	err := serialize.Serialize(uint64(1000), &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	result, err := serialize.DeserializeInt64(&buf)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if result != 1000 {
		t.Errorf("Expected 1000, got %d", result)
	}
}

func TestDeserializeInt64OverflowError(t *testing.T) {
	// Test that very large uint64 values cause overflow error
	var buf bytes.Buffer
	err := serialize.Serialize(uint64(18446744073709551615), &buf) // max uint64
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	_, err = serialize.DeserializeInt64(&buf)
	if err == nil {
		t.Error("Expected overflow error, but got none")
	}
}

func TestDeserializeUint64FromNegativeError(t *testing.T) {
	// Test that negative values cause error when deserializing to uint64
	var buf bytes.Buffer
	err := serialize.Serialize(int64(-1), &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	_, err = serialize.DeserializeUint64(&buf)
	if err == nil {
		t.Error("Expected error for negative value, but got none")
	}
}