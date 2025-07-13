package main_test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

// BenchmarkMapStringInt tests the optimized fast path
func BenchmarkMapStringInt(b *testing.B) {
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

// BenchmarkMapStringString tests the optimized fast path
func BenchmarkMapStringString(b *testing.B) {
	testMap := map[string]string{
		"hello": "world",
		"foo":   "bar",
		"test":  "data",
		"key":   "value",
		"alpha": "beta",
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

// BenchmarkMapStringInterface tests the semi-optimized path
func BenchmarkMapStringInterface(b *testing.B) {
	testMap := map[string]interface{}{
		"int":    42,
		"string": "hello",
		"bool":   true,
		"float":  3.14,
		"data":   []byte{1, 2, 3},
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

// BenchmarkMapUnoptimized tests the generic reflection path
func BenchmarkMapUnoptimized(b *testing.B) {
	// Use a map type that won't hit the fast paths
	testMap := map[int64]float64{
		1: 1.1,
		2: 2.2,
		3: 3.3,
		4: 4.4,
		5: 5.5,
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

// BenchmarkMapLargeStringInt tests performance with larger maps
func BenchmarkMapLargeStringInt(b *testing.B) {
	testMap := make(map[string]int)
	for i := 0; i < 50; i++ {
		testMap[string(rune('A'+i%26))+string(rune('0'+i%10))] = i * i
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

// Deserialization benchmarks

// BenchmarkMapStringIntDeserialize tests optimized deserialization
func BenchmarkMapStringIntDeserialize(b *testing.B) {
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

// BenchmarkMapStringStringDeserialize tests optimized deserialization
func BenchmarkMapStringStringDeserialize(b *testing.B) {
	testMap := map[string]string{
		"hello": "world",
		"foo":   "bar",
		"test":  "data",
		"key":   "value",
		"alpha": "beta",
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
		var result map[string]string
		err := serialize.Deserialize(bytes.NewReader(data), &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMapUnoptimizedDeserialize tests generic reflection deserialization
func BenchmarkMapUnoptimizedDeserialize(b *testing.B) {
	// Use a map type that won't hit the fast paths
	testMap := map[int64]float64{
		1: 1.1,
		2: 2.2,
		3: 3.3,
		4: 4.4,
		5: 5.5,
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
		var result map[int64]float64
		err := serialize.Deserialize(bytes.NewReader(data), &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}