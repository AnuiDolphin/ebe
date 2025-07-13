package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

// BenchmarkMapSerialization tests map serialization performance
func BenchmarkMapSerialization(b *testing.B) {
	testMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := serialize.Serialize(testMap, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMapDeserialization tests map deserialization performance
func BenchmarkMapDeserialization(b *testing.B) {
	testMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
		"four":  4,
		"five":  5,
	}

	// Pre-serialize the map
	var buf bytes.Buffer
	err := serialize.Serialize(testMap, &buf)
	if err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]int
		err := serialize.Deserialize(bytes.NewReader(data), &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}