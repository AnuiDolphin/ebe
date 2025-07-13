package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

// Benchmark deserialization performance for arrays
func BenchmarkIntArrayDeserialization(b *testing.B) {
	// Create test data
	data := make([]int32, 1000)
	for i := range data {
		data[i] = int32(i)
	}
	
	// Serialize once to get the binary data
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []int32
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringArrayDeserialization(b *testing.B) {
	// Create test data
	data := make([]string, 100)
	for i := range data {
		data[i] = "test string " + string(rune(i))
	}
	
	// Serialize once to get the binary data
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []string
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUintArrayDeserialization(b *testing.B) {
	// Create test data
	data := make([]uint64, 1000)
	for i := range data {
		data[i] = uint64(i * 1000)
	}
	
	// Serialize once to get the binary data
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []uint64
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFloatArrayDeserialization(b *testing.B) {
	// Create test data
	data := make([]float64, 1000)
	for i := range data {
		data[i] = float64(i) * 3.14159
	}
	
	// Serialize once to get the binary data
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []float64
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark small vs large arrays to see scaling impact
func BenchmarkSmallIntArrayDeserialization(b *testing.B) {
	data := []int32{1, 2, 3, 4, 5}
	
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []int32
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeIntArrayDeserialization(b *testing.B) {
	data := make([]int32, 10000)
	for i := range data {
		data[i] = int32(i)
	}
	
	var buf bytes.Buffer
	err := serialize.Serialize(data, &buf)
	if err != nil {
		b.Fatal(err)
	}
	serializedData := buf.Bytes()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(serializedData)
		var result []int32
		err := serialize.Deserialize(reader, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}