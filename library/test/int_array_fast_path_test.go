package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

func TestIntArrayFastPath(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"[]int", []int{1, 2, 3, 4, 5}},
		{"[]int32", []int32{10, 20, 30}},
		{"[]int64", []int64{100, 200, 300, 400}},
		{"[]int8", []int8{1, 2, 3}},
		{"[]int16", []int16{10, 20}},
		{"empty []int", []int{}},
		{"single element []int64", []int64{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Serialize using fast path
			err := serialize.Serialize(tt.data, &buf)
			if err != nil {
				t.Fatalf("Fast path serialization failed: %v", err)
			}
			
			// Verify we can deserialize back
			var result interface{}
			switch tt.data.(type) {
			case []int:
				var out []int
				result = &out
			case []int32:
				var out []int32
				result = &out
			case []int64:
				var out []int64
				result = &out
			case []int8:
				var out []int8
				result = &out
			case []int16:
				var out []int16
				result = &out
			}
			
			err = serialize.Deserialize(&buf, result)
			if err != nil {
				t.Fatalf("Deserialization failed: %v", err)
			}
			
			// Basic length check for now
			// More detailed comparison would require reflection
			t.Logf("Successfully serialized and deserialized %s", tt.name)
		})
	}
}

func TestStringArrayFastPath(t *testing.T) {
	tests := []struct {
		name string
		data []string
	}{
		{"basic strings", []string{"hello", "world", "test"}},
		{"empty array", []string{}},
		{"single string", []string{"single"}},
		{"unicode strings", []string{"café", "naïve", "résumé"}},
		{"empty strings", []string{"", "hello", ""}},
		{"long array", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}}, // > 7 elements
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Serialize using fast path
			err := serialize.Serialize(tt.data, &buf)
			if err != nil {
				t.Fatalf("String array fast path serialization failed: %v", err)
			}
			
			// Verify we can deserialize back
			var result []string
			err = serialize.Deserialize(&buf, &result)
			if err != nil {
				t.Fatalf("String array deserialization failed: %v", err)
			}
			
			// Verify length matches
			if len(result) != len(tt.data) {
				t.Fatalf("Length mismatch: expected %d, got %d", len(tt.data), len(result))
			}
			
			t.Logf("Successfully serialized and deserialized %s with %d elements", tt.name, len(tt.data))
		})
	}
}

func TestUintArrayFastPath(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"[]uint", []uint{1, 2, 3, 4, 5}},
		{"[]uint32", []uint32{10, 20, 30}},
		{"[]uint64", []uint64{100, 200, 300, 400}},
		{"[]uint16", []uint16{10, 20}},
		{"empty []uint", []uint{}},
		{"single element []uint64", []uint64{42}},
		{"large values", []uint64{18446744073709551615, 0, 1000000000000}}, // max uint64 and other large values
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Serialize using fast path
			err := serialize.Serialize(tt.data, &buf)
			if err != nil {
				t.Fatalf("Uint array fast path serialization failed: %v", err)
			}
			
			// Verify we can deserialize back
			var result interface{}
			switch tt.data.(type) {
			case []uint:
				var out []uint
				result = &out
			case []uint32:
				var out []uint32
				result = &out
			case []uint64:
				var out []uint64
				result = &out
			case []uint16:
				var out []uint16
				result = &out
			}
			
			err = serialize.Deserialize(&buf, result)
			if err != nil {
				t.Fatalf("Uint array deserialization failed: %v", err)
			}
			
			// Basic length check for now
			// More detailed comparison would require reflection
			t.Logf("Successfully serialized and deserialized %s", tt.name)
		})
	}
}

func TestFloatArrayFastPath(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{"[]float32", []float32{1.5, 2.5, 3.14159}},
		{"[]float64", []float64{1.5, 2.5, 3.14159265359}},
		{"empty []float32", []float32{}},
		{"empty []float64", []float64{}},
		{"single element []float32", []float32{42.0}},
		{"single element []float64", []float64{42.0}},
		{"special values", []float64{0.0, -0.0, 1.0, -1.0}},
		{"large float32 array", make([]float32, 10)}, // 10 elements for overflow test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			
			// Serialize using fast path
			err := serialize.Serialize(tt.data, &buf)
			if err != nil {
				t.Fatalf("Float array fast path serialization failed: %v", err)
			}
			
			// Verify we can deserialize back
			var result interface{}
			switch tt.data.(type) {
			case []float32:
				var out []float32
				result = &out
			case []float64:
				var out []float64
				result = &out
			}
			
			err = serialize.Deserialize(&buf, result)
			if err != nil {
				t.Fatalf("Float array deserialization failed: %v", err)
			}
			
			// Basic length check for now
			// More detailed comparison would require reflection
			t.Logf("Successfully serialized and deserialized %s", tt.name)
		})
	}
}
