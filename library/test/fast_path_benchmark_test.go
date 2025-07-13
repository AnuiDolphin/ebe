package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

// Benchmark our fast paths vs the generic reflection-based approach
func BenchmarkIntArraySerialization(b *testing.B) {
	data := make([]int32, 1000)
	for i := range data {
		data[i] = int32(i)
	}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}

func BenchmarkStringArraySerialization(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = "test string " + string(rune(i))
	}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}

func BenchmarkUintArraySerialization(b *testing.B) {
	data := make([]uint64, 1000)
	for i := range data {
		data[i] = uint64(i * 1000)
	}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}

// Benchmark small vs large arrays to see scaling
func BenchmarkSmallIntArray(b *testing.B) {
	data := []int32{1, 2, 3, 4, 5}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}

func BenchmarkLargeIntArray(b *testing.B) {
	data := make([]int32, 10000)
	for i := range data {
		data[i] = int32(i)
	}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}
func BenchmarkFloatArraySerialization(b *testing.B) {
	data := make([]float64, 1000)
	for i := range data {
		data[i] = float64(i) * 3.14159
	}
	
	var buf bytes.Buffer
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		buf.Reset()
		err := serialize.Serialize(data, &buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	b.ReportAllocs()
}
