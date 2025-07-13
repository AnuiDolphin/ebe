package main_test

import (
	"bytes"
	"ebe/serialize"
	"fmt"
	"testing"
)

// These benchmarks measure the performance impact of using type-specific deserializers
// in array deserialization instead of the generic Deserialize function

func BenchmarkArrayDeserializationOptimized(b *testing.B) {
	// Create test data for different array types
	intData := make([]int32, 1000)
	for i := range intData {
		intData[i] = int32(i)
	}

	uintData := make([]uint64, 1000) 
	for i := range uintData {
		uintData[i] = uint64(i * 1000)
	}

	floatData := make([]float64, 1000)
	for i := range floatData {
		floatData[i] = float64(i) * 3.14159
	}

	stringData := make([]string, 100)
	for i := range stringData {
		stringData[i] = "test string " + string(rune(i))
	}

	// Serialize all test data once
	var intBuf, uintBuf, floatBuf, stringBuf bytes.Buffer
	
	serialize.Serialize(intData, &intBuf)
	serialize.Serialize(uintData, &uintBuf)  
	serialize.Serialize(floatData, &floatBuf)
	serialize.Serialize(stringData, &stringBuf)

	serializedInt := intBuf.Bytes()
	serializedUint := uintBuf.Bytes()
	serializedFloat := floatBuf.Bytes()
	serializedString := stringBuf.Bytes()

	b.Run("IntArray", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			reader := bytes.NewReader(serializedInt)
			var result []int32
			err := serialize.Deserialize(reader, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("UintArray", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			reader := bytes.NewReader(serializedUint)
			var result []uint64
			err := serialize.Deserialize(reader, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FloatArray", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			reader := bytes.NewReader(serializedFloat)
			var result []float64
			err := serialize.Deserialize(reader, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StringArray", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		
		for i := 0; i < b.N; i++ {
			reader := bytes.NewReader(serializedString)
			var result []string
			err := serialize.Deserialize(reader, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Test specific array sizes to see scaling impact
func BenchmarkOptimizedArrayScaling(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		// Create test data
		data := make([]int32, size)
		for i := range data {
			data[i] = int32(i)
		}

		// Serialize once
		var buf bytes.Buffer
		serialize.Serialize(data, &buf)
		serializedData := buf.Bytes()

		b.Run(fmt.Sprintf("IntArray_%d", size), func(b *testing.B) {
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
		})
	}
}