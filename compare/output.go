package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// outputCSV writes benchmark results to a CSV file
func outputCSV(results []BenchmarkResult) error {
	file, err := os.Create("benchmark_results.csv")
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"TestName",
		"EBESize",
		"ProtobufSize", 
		"SizeDifference",
		"SizeRatio",
		"EBETime_ns",
		"ProtobufTime_ns",
		"TimeDifference_ns",
		"TimeRatio",
		"EBEMemory",
		"ProtobufMemory",
		"MemoryDifference",
		"MemoryRatio",
	}
	
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, result := range results {
		row := []string{
			result.TestName,
			strconv.Itoa(result.EBESize),
			strconv.Itoa(result.ProtobufSize),
			strconv.Itoa(result.SizeDifference),
			fmt.Sprintf("%.6f", result.SizeRatio),
			strconv.FormatInt(result.EBETime.Nanoseconds(), 10),
			strconv.FormatInt(result.ProtobufTime.Nanoseconds(), 10),
			strconv.FormatInt(result.TimeDifference.Nanoseconds(), 10),
			fmt.Sprintf("%.6f", result.TimeRatio),
			strconv.FormatInt(result.EBEMemory, 10),
			strconv.FormatInt(result.ProtobufMemory, 10),
			strconv.FormatInt(result.MemoryDifference, 10),
			fmt.Sprintf("%.6f", result.MemoryRatio),
		}
		
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	fmt.Printf("Results written to benchmark_results.csv\n")
	return nil
}

// outputJSON writes benchmark results to a JSON file
func outputJSON(results []BenchmarkResult) error {
	file, err := os.Create("benchmark_results.json")
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	// Create structured output
	output := struct {
		Summary struct {
			TotalTests     int     `json:"total_tests"`
			TotalEBESize   int     `json:"total_ebe_size"`
			TotalPBSize    int     `json:"total_pb_size"`
			OverallRatio   float64 `json:"overall_ratio"`
			EBEWins        int     `json:"ebe_wins"`
			ProtobufWins   int     `json:"protobuf_wins"`
			Ties           int     `json:"ties"`
		} `json:"summary"`
		Results []BenchmarkResult `json:"results"`
	}{
		Results: results,
	}

	// Calculate summary statistics
	var totalEBESize, totalPBSize int
	var ebeWins, pbWins, ties int
	
	for _, result := range results {
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
	
	output.Summary.TotalTests = len(results)
	output.Summary.TotalEBESize = totalEBESize
	output.Summary.TotalPBSize = totalPBSize
	output.Summary.OverallRatio = float64(totalEBESize) / float64(totalPBSize)
	output.Summary.EBEWins = ebeWins
	output.Summary.ProtobufWins = pbWins
	output.Summary.Ties = ties

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Printf("Results written to benchmark_results.json\n")
	return nil
}