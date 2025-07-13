package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"compare/proto/testdata"
)

func main() {
	// Command-line flags
	var (
		iterations = flag.Int("iterations", 10, "Number of benchmark iterations per test")
		warmupRuns = flag.Int("warmup", 3, "Number of warmup runs before benchmarking")
		verbose    = flag.Bool("verbose", false, "Enable verbose output")
		testGroup  = flag.String("group", "all", "Test group to run: all, primitives, collections, complex, edge, realworld")
		output     = flag.String("output", "console", "Output format: console, csv, json")
	)
	flag.Parse()

	fmt.Printf("EBE vs Protocol Buffers Serialization Comparison\n")
	fmt.Printf("================================================\n")
	fmt.Printf("Iterations: %d, Warmup runs: %d\n", *iterations, *warmupRuns)
	fmt.Printf("Test group: %s, Output format: %s\n\n", *testGroup, *output)

	// Create comparison framework
	framework := NewComparisonFramework(*iterations, *warmupRuns)

	// Generate test data
	testData := GenerateTestData()

	// Run tests based on selected group
	switch *testGroup {
	case "all":
		runAllTests(framework, testData, *verbose)
	case "primitives":
		runPrimitiveTests(framework, testData, *verbose)
	case "collections":
		runCollectionTests(framework, testData, *verbose)
	case "complex":
		runComplexTests(framework, testData, *verbose)
	case "edge":
		runEdgeCaseTests(framework, testData, *verbose)
	case "realworld":
		runRealWorldTests(framework, testData, *verbose)
	default:
		fmt.Printf("Unknown test group: %s\n", *testGroup)
		os.Exit(1)
	}

	// Output results
	switch *output {
	case "console":
		framework.PrintSummary()
	case "csv":
		err := outputCSV(framework.GetResults())
		if err != nil {
			log.Fatalf("Failed to output CSV: %v", err)
		}
	case "json":
		err := outputJSON(framework.GetResults())
		if err != nil {
			log.Fatalf("Failed to output JSON: %v", err)
		}
	default:
		fmt.Printf("Unknown output format: %s\n", *output)
		os.Exit(1)
	}
}

// runAllTests executes all test categories
func runAllTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	runPrimitiveTests(framework, testData, verbose)
	runCollectionTests(framework, testData, verbose)
	runComplexTests(framework, testData, verbose)
	runEdgeCaseTests(framework, testData, verbose)
	runRealWorldTests(framework, testData, verbose)
}

// runPrimitiveTests runs tests for primitive data types
func runPrimitiveTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	fmt.Println("=== PRIMITIVE TYPES ===")

	test := SerializationTest{
		Name:    "Primitives",
		EBEData: testData.Primitives,
		PBData:  ConvertPrimitives(testData.Primitives),
	}

	result := framework.RunComparison(test)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}
}

// runCollectionTests runs tests for collection data types
func runCollectionTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	fmt.Println("=== COLLECTION TYPES ===")

	// Small arrays test
	smallTest := SerializationTest{
		Name: "Small Arrays",
		EBEData: struct {
			IntArray    []int32
			StringArray []string
		}{
			IntArray:    testData.Collections.SmallIntArray,
			StringArray: testData.Collections.SmallStringArray,
		},
		PBData: &testdata.CollectionTypes{
			SmallIntArray:    testData.Collections.SmallIntArray,
			SmallStringArray: testData.Collections.SmallStringArray,
		},
	}

	result := framework.RunComparison(smallTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Medium arrays test
	mediumTest := SerializationTest{
		Name: "Medium Arrays (100 elements)",
		EBEData: struct {
			IntArray   []int32
			FloatArray []float64
		}{
			IntArray:   testData.Collections.MediumIntArray,
			FloatArray: testData.Collections.MediumFloatArray,
		},
		PBData: &testdata.CollectionTypes{
			MediumIntArray:   testData.Collections.MediumIntArray,
			MediumFloatArray: testData.Collections.MediumFloatArray,
		},
	}

	result = framework.RunComparison(mediumTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Large arrays test
	largeTest := SerializationTest{
		Name: "Large Arrays (1000+ elements)",
		EBEData: struct {
			IntArray    []int64
			StringArray []string
		}{
			IntArray:    testData.Collections.LargeIntArray,
			StringArray: testData.Collections.LargeStringArray,
		},
		PBData: &testdata.CollectionTypes{
			LargeIntArray:    testData.Collections.LargeIntArray,
			LargeStringArray: testData.Collections.LargeStringArray,
		},
	}

	result = framework.RunComparison(largeTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}
}

// runComplexTests runs tests for complex nested structures
func runComplexTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	fmt.Println("=== COMPLEX STRUCTURES ===")

	// Person test
	personTest := SerializationTest{
		Name:    "Person Structure",
		EBEData: testData.Complex.Person,
		PBData:  ConvertPerson(testData.Complex.Person),
	}

	result := framework.RunComparison(personTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Company test
	companyTest := SerializationTest{
		Name:    "Company Structure",
		EBEData: testData.Complex.Company,
		PBData:  ConvertCompany(testData.Complex.Company),
	}

	result = framework.RunComparison(companyTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}
}

// runEdgeCaseTests runs tests for edge cases and boundary conditions
func runEdgeCaseTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	fmt.Println("=== EDGE CASES ===")

	// Empty values test
	emptyTest := SerializationTest{
		Name: "Empty Values",
		EBEData: struct {
			EmptyString string
			EmptyBytes  []byte
			EmptyArray  []int32
		}{
			EmptyString: testData.EdgeCases.EmptyString,
			EmptyBytes:  testData.EdgeCases.EmptyBytes,
			EmptyArray:  testData.EdgeCases.EmptyArray,
		},
		PBData: &testdata.EdgeCaseTypes{
			EmptyString: testData.EdgeCases.EmptyString,
			EmptyBytes:  testData.EdgeCases.EmptyBytes,
			EmptyArray:  testData.EdgeCases.EmptyArray,
		},
	}

	result := framework.RunComparison(emptyTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Large data test
	largeTest := SerializationTest{
		Name: "Large Data (1MB each)",
		EBEData: struct {
			LargeString string
			LargeBytes  []byte
		}{
			LargeString: testData.EdgeCases.LargeString,
			LargeBytes:  testData.EdgeCases.LargeBytes,
		},
		PBData: &testdata.EdgeCaseTypes{
			LargeString: testData.EdgeCases.LargeString,
			LargeBytes:  testData.EdgeCases.LargeBytes,
		},
	}

	result = framework.RunComparison(largeTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Unicode test
	unicodeTest := SerializationTest{
		Name: "Unicode Strings",
		EBEData: struct {
			UnicodeString string
			SpecialChars  string
		}{
			UnicodeString: testData.EdgeCases.UnicodeString,
			SpecialChars:  testData.EdgeCases.SpecialChars,
		},
		PBData: &testdata.EdgeCaseTypes{
			UnicodeString: testData.EdgeCases.UnicodeString,
			SpecialChars:  testData.EdgeCases.SpecialChars,
		},
	}

	result = framework.RunComparison(unicodeTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}
}

// runRealWorldTests runs tests for realistic data scenarios
func runRealWorldTests(framework *ComparisonFramework, testData *TestData, verbose bool) {
	fmt.Println("=== REAL WORLD SCENARIOS ===")

	// User profile test
	userTest := SerializationTest{
		Name:    "User Profile",
		EBEData: testData.RealWorld.UserProfile,
		PBData:  ConvertUserProfile(testData.RealWorld.UserProfile),
	}

	result := framework.RunComparison(userTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Configuration test
	configTest := SerializationTest{
		Name:    "Application Configuration",
		EBEData: testData.RealWorld.Config,
		PBData:  ConvertConfiguration(testData.RealWorld.Config),
	}

	result = framework.RunComparison(configTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Time series test
	timeSeriesTest := SerializationTest{
		Name:    "Time Series Data (100 points)",
		EBEData: testData.RealWorld.TimeSeries,
		PBData: func() *testdata.RealWorldTypes {
			pb := &testdata.RealWorldTypes{}
			pb.TimeSeries = make([]*testdata.TimeSeriesPoint, len(testData.RealWorld.TimeSeries))
			for i, ts := range testData.RealWorld.TimeSeries {
				pb.TimeSeries[i] = ConvertTimeSeriesPoint(ts)
			}
			return pb
		}(),
	}

	result = framework.RunComparison(timeSeriesTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}

	// Log entries test
	logTest := SerializationTest{
		Name:    "Log Entries (50 entries)",
		EBEData: testData.RealWorld.LogEntries,
		PBData: func() *testdata.RealWorldTypes {
			pb := &testdata.RealWorldTypes{}
			pb.LogEntries = make([]*testdata.LogEntry, len(testData.RealWorld.LogEntries))
			for i, log := range testData.RealWorld.LogEntries {
				pb.LogEntries[i] = ConvertLogEntry(log)
			}
			return pb
		}(),
	}

	result = framework.RunComparison(logTest)
	if verbose {
		framework.PrintResult(result)
	} else {
		printCompactResult(result)
	}
}

// printCompactResult prints a single line summary of a benchmark result
func printCompactResult(result BenchmarkResult) {
	winner := "="
	if result.SizeDifference < 0 {
		winner = "EBE"
	} else if result.SizeDifference > 0 {
		winner = "PB"
	}

	fmt.Printf("%-30s | EBE: %8d | PB: %8d | Ratio: %5.2fx | Winner: %s\n",
		result.TestName, result.EBESize, result.ProtobufSize, result.SizeRatio, winner)
}
