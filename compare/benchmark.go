package main

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"ebe/serialize"
	"google.golang.org/protobuf/proto"
)

// BenchmarkResult represents the results of a serialization comparison
type BenchmarkResult struct {
	TestName        string
	EBESize         int
	ProtobufSize    int
	SizeDifference  int     // EBE - Protobuf (negative means EBE is smaller)
	SizeRatio       float64 // EBE / Protobuf
	EBETime         time.Duration
	ProtobufTime    time.Duration
	TimeDifference  time.Duration
	TimeRatio       float64
	EBEMemory       int64
	ProtobufMemory  int64
	MemoryDifference int64
	MemoryRatio     float64
}

// ComparisonFramework manages the benchmarking and comparison process
type ComparisonFramework struct {
	iterations int
	warmupRuns int
	results    []BenchmarkResult
}

// NewComparisonFramework creates a new benchmark framework
func NewComparisonFramework(iterations, warmupRuns int) *ComparisonFramework {
	return &ComparisonFramework{
		iterations: iterations,
		warmupRuns: warmupRuns,
		results:    make([]BenchmarkResult, 0),
	}
}

// SerializationTest represents a single test case
type SerializationTest struct {
	Name    string
	EBEData interface{}
	PBData  proto.Message
}

// RunComparison executes a complete comparison between EBE and Protocol Buffers
func (cf *ComparisonFramework) RunComparison(test SerializationTest) BenchmarkResult {
	// Warmup runs
	for i := 0; i < cf.warmupRuns; i++ {
		cf.benchmarkEBE(test.EBEData)
		cf.benchmarkProtobuf(test.PBData)
		runtime.GC() // Force garbage collection between warmup runs
	}

	// Actual benchmark runs
	ebeResults := make([]BenchmarkMeasurement, cf.iterations)
	pbResults := make([]BenchmarkMeasurement, cf.iterations)

	for i := 0; i < cf.iterations; i++ {
		runtime.GC() // Ensure clean state for each measurement
		
		ebeResults[i] = cf.benchmarkEBE(test.EBEData)
		runtime.GC()
		
		pbResults[i] = cf.benchmarkProtobuf(test.PBData)
		runtime.GC()
	}

	// Calculate statistics
	ebeStats := calculateStats(ebeResults)
	pbStats := calculateStats(pbResults)

	result := BenchmarkResult{
		TestName:         test.Name,
		EBESize:          ebeStats.Size,
		ProtobufSize:     pbStats.Size,
		SizeDifference:   ebeStats.Size - pbStats.Size,
		SizeRatio:        float64(ebeStats.Size) / float64(pbStats.Size),
		EBETime:          ebeStats.Time,
		ProtobufTime:     pbStats.Time,
		TimeDifference:   ebeStats.Time - pbStats.Time,
		TimeRatio:        float64(ebeStats.Time.Nanoseconds()) / float64(pbStats.Time.Nanoseconds()),
		EBEMemory:        ebeStats.Memory,
		ProtobufMemory:   pbStats.Memory,
		MemoryDifference: ebeStats.Memory - pbStats.Memory,
		MemoryRatio:      float64(ebeStats.Memory) / float64(pbStats.Memory),
	}

	cf.results = append(cf.results, result)
	return result
}

// BenchmarkMeasurement represents a single benchmark measurement
type BenchmarkMeasurement struct {
	Size   int
	Time   time.Duration
	Memory int64
}

// Stats represents aggregated statistics from multiple measurements
type Stats struct {
	Size   int           // Size should be consistent, so we take the first one
	Time   time.Duration // Average time
	Memory int64         // Average memory
}

// benchmarkEBE measures EBE serialization performance
func (cf *ComparisonFramework) benchmarkEBE(data interface{}) BenchmarkMeasurement {
	var memStats1, memStats2 runtime.MemStats
	runtime.ReadMemStats(&memStats1)

	var buf bytes.Buffer
	start := time.Now()
	
	err := serialize.Serialize(data, &buf)
	if err != nil {
		panic(fmt.Sprintf("EBE serialization failed: %v", err))
	}
	
	elapsed := time.Since(start)
	runtime.ReadMemStats(&memStats2)

	return BenchmarkMeasurement{
		Size:   buf.Len(),
		Time:   elapsed,
		Memory: int64(memStats2.TotalAlloc - memStats1.TotalAlloc),
	}
}

// benchmarkProtobuf measures Protocol Buffers serialization performance
func (cf *ComparisonFramework) benchmarkProtobuf(data proto.Message) BenchmarkMeasurement {
	var memStats1, memStats2 runtime.MemStats
	runtime.ReadMemStats(&memStats1)

	start := time.Now()
	
	serialized, err := proto.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Protobuf serialization failed: %v", err))
	}
	
	elapsed := time.Since(start)
	runtime.ReadMemStats(&memStats2)

	return BenchmarkMeasurement{
		Size:   len(serialized),
		Time:   elapsed,
		Memory: int64(memStats2.TotalAlloc - memStats1.TotalAlloc),
	}
}

// calculateStats computes average statistics from multiple measurements
func calculateStats(measurements []BenchmarkMeasurement) Stats {
	if len(measurements) == 0 {
		return Stats{}
	}

	var totalTime time.Duration
	var totalMemory int64
	size := measurements[0].Size // Size should be consistent

	for _, m := range measurements {
		totalTime += m.Time
		totalMemory += m.Memory
	}

	return Stats{
		Size:   size,
		Time:   totalTime / time.Duration(len(measurements)),
		Memory: totalMemory / int64(len(measurements)),
	}
}

// GetResults returns all benchmark results
func (cf *ComparisonFramework) GetResults() []BenchmarkResult {
	return cf.results
}

// PrintResult prints a single benchmark result
func (cf *ComparisonFramework) PrintResult(result BenchmarkResult) {
	fmt.Printf("=== %s ===\n", result.TestName)
	fmt.Printf("Size:   EBE: %6d bytes | Protobuf: %6d bytes | Diff: %+6d bytes (%.2fx)\n",
		result.EBESize, result.ProtobufSize, result.SizeDifference, result.SizeRatio)
	fmt.Printf("Time:   EBE: %8s | Protobuf: %8s | Diff: %+8s (%.2fx)\n",
		result.EBETime, result.ProtobufTime, result.TimeDifference, result.TimeRatio)
	fmt.Printf("Memory: EBE: %6d bytes | Protobuf: %6d bytes | Diff: %+6d bytes (%.2fx)\n",
		result.EBEMemory, result.ProtobufMemory, result.MemoryDifference, result.MemoryRatio)
	
	// Indicate which is better
	if result.SizeDifference < 0 {
		fmt.Printf("Size winner: EBE (%.1f%% smaller)\n", 
			(1.0-result.SizeRatio)*100)
	} else if result.SizeDifference > 0 {
		fmt.Printf("Size winner: Protobuf (%.1f%% smaller)\n", 
			(1.0-1.0/result.SizeRatio)*100)
	} else {
		fmt.Printf("Size winner: Tie\n")
	}
	
	fmt.Println()
}

// PrintSummary prints a summary of all results
func (cf *ComparisonFramework) PrintSummary() {
	if len(cf.results) == 0 {
		fmt.Println("No results to summarize")
		return
	}

	fmt.Printf("\n=== SUMMARY (%d tests) ===\n", len(cf.results))
	
	var totalEBESize, totalPBSize int
	var ebeWins, pbWins, ties int
	
	for _, result := range cf.results {
		totalEBESize += result.EBESize
		totalPBSize += result.ProtobufSize
		
		if result.SizeDifference < 0 {
			ebeWins++
		} else if result.SizeDifference > 0 {
			pbWins++
		} else {
			ties++
		}
	}
	
	fmt.Printf("Total serialized size:\n")
	fmt.Printf("  EBE:      %10d bytes\n", totalEBESize)
	fmt.Printf("  Protobuf: %10d bytes\n", totalPBSize)
	fmt.Printf("  Ratio:    %.3fx (EBE/Protobuf)\n", float64(totalEBESize)/float64(totalPBSize))
	
	fmt.Printf("\nSize comparison wins:\n")
	fmt.Printf("  EBE wins:      %d\n", ebeWins)
	fmt.Printf("  Protobuf wins: %d\n", pbWins)
	fmt.Printf("  Ties:          %d\n", ties)
}