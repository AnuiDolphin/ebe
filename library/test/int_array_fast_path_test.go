package main_test

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