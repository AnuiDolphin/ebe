package main

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

// TestData defines common test structures for comparing EBE vs Protocol Buffers
type TestData struct {
	// Primitive types
	Primitives PrimitiveTypes
	
	// Collections
	Collections CollectionTypes
	
	// Complex structures
	Complex ComplexTypes
	
	// Edge cases
	EdgeCases EdgeCaseTypes
	
	// Real-world scenarios
	RealWorld RealWorldTypes
}

// PrimitiveTypes covers basic data types
type PrimitiveTypes struct {
	// Integers - various sizes and signs
	UInt8   uint8
	UInt16  uint16
	UInt32  uint32
	UInt64  uint64
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	
	// Floating point
	Float32 float32
	Float64 float64
	
	// Boolean
	Bool    bool
	
	// String
	String  string
	
	// Binary data
	Bytes   []byte
}

// CollectionTypes covers arrays and slices
type CollectionTypes struct {
	// Small arrays
	SmallIntArray    []int32
	SmallStringArray []string
	
	// Medium arrays (100 elements)
	MediumIntArray   []int32
	MediumFloatArray []float64
	
	// Large arrays (10000 elements)
	LargeIntArray    []int64
	LargeStringArray []string
	
	// Mixed type arrays (using interface{} for EBE, oneof for protobuf)
	MixedArray       []interface{}
	
	// Nested arrays
	NestedIntArray   [][]int32
}

// ComplexTypes covers nested structures
type ComplexTypes struct {
	// Nested structs
	Person      Person
	Company     Company
	
	// Struct arrays
	People      []Person
	Departments []Department
	
	// Maps (will use repeated key-value pairs in protobuf)
	StringMap   map[string]string
	IntMap      map[string]int32
}

// EdgeCaseTypes covers boundary conditions
type EdgeCaseTypes struct {
	// Empty values
	EmptyString     string
	EmptyBytes      []byte
	EmptyArray      []int32
	
	// Nil values
	NilBytes        []byte
	
	// Large values
	LargeString     string   // 1MB string
	LargeBytes      []byte   // 1MB binary data
	
	// Unicode and special characters
	UnicodeString   string
	SpecialChars    string
	
	// Extreme numbers
	MaxInt64        int64
	MinInt64        int64
	MaxUint64       uint64
	SmallFloat      float64
	LargeFloat      float64
	NaNFloat        float64
	InfFloat        float64
}

// RealWorldTypes represents realistic data structures
type RealWorldTypes struct {
	// User profile
	UserProfile     UserProfile
	
	// Configuration data
	Config          Configuration
	
	// Time series data
	TimeSeries      []TimeSeriesPoint
	
	// Log entries
	LogEntries      []LogEntry
	
	// API responses
	APIResponse     APIResponse
}

// Supporting structures for complex types

type Person struct {
	ID       uint64
	Name     string
	Email    string
	Age      int32
	Active   bool
	Tags     []string
	Metadata map[string]string
}

type Company struct {
	ID          uint64
	Name        string
	Founded     int32
	Employees   []Person
	Departments []Department
	Revenue     float64
}

type Department struct {
	ID       uint32
	Name     string
	Manager  Person
	Budget   float64
	Projects []string
}

type UserProfile struct {
	UserID      uint64
	Username    string
	FullName    string
	Email       string
	PhoneNumber string
	BirthDate   time.Time
	Preferences map[string]interface{}
	Friends     []uint64
	Groups      []string
	LastLogin   time.Time
	IsActive    bool
	ProfilePic  []byte
}

type Configuration struct {
	AppName     string
	Version     string
	Environment string
	Features    map[string]bool
	Limits      map[string]int32
	Endpoints   []string
	Timeouts    map[string]int32
	Debug       bool
}

type TimeSeriesPoint struct {
	Timestamp time.Time
	Value     float64
	Tags      map[string]string
	Source    string
}

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Source    string
	ThreadID  uint32
	Data      map[string]interface{}
}

type APIResponse struct {
	Status    int32
	Message   string
	Data      interface{}
	Timestamp time.Time
	RequestID string
	Errors    []string
	Meta      map[string]interface{}
}

// GenerateTestData creates comprehensive test data for comparison
func GenerateTestData() *TestData {
	return &TestData{
		Primitives: generatePrimitives(),
		Collections: generateCollections(),
		Complex: generateComplex(),
		EdgeCases: generateEdgeCases(),
		RealWorld: generateRealWorld(),
	}
}

func generatePrimitives() PrimitiveTypes {
	return PrimitiveTypes{
		UInt8:   255,
		UInt16:  65535,
		UInt32:  4294967295,
		UInt64:  18446744073709551615,
		Int8:    -128,
		Int16:   -32768,
		Int32:   -2147483648,
		Int64:   -9223372036854775808,
		Float32: 3.14159,
		Float64: 2.718281828459045,
		Bool:    true,
		String:  "Hello, World! 🌍",
		Bytes:   []byte{0x00, 0x01, 0x02, 0x03, 0xFF},
	}
}

func generateCollections() CollectionTypes {
	// Small arrays
	smallInts := make([]int32, 10)
	for i := range smallInts {
		smallInts[i] = int32(i * i)
	}
	
	smallStrings := []string{"apple", "banana", "cherry", "date", "elderberry"}
	
	// Medium arrays
	mediumInts := make([]int32, 100)
	for i := range mediumInts {
		mediumInts[i] = int32(i * 2)
	}
	
	mediumFloats := make([]float64, 100)
	for i := range mediumFloats {
		mediumFloats[i] = float64(i) * 0.1
	}
	
	// Large arrays
	largeInts := make([]int64, 10000)
	for i := range largeInts {
		largeInts[i] = int64(i * i * i)
	}
	
	largeStrings := make([]string, 1000)
	for i := range largeStrings {
		largeStrings[i] = fmt.Sprintf("string_%d", i)
	}
	
	// Mixed array
	mixedArray := []interface{}{
		int32(42),
		"mixed string",
		3.14159,
		true,
		[]byte{0xDE, 0xAD, 0xBE, 0xEF},
	}
	
	// Nested arrays
	nestedInts := [][]int32{
		{1, 2, 3},
		{4, 5, 6, 7},
		{8, 9},
		{10, 11, 12, 13, 14},
	}
	
	return CollectionTypes{
		SmallIntArray:    smallInts,
		SmallStringArray: smallStrings,
		MediumIntArray:   mediumInts,
		MediumFloatArray: mediumFloats,
		LargeIntArray:    largeInts,
		LargeStringArray: largeStrings,
		MixedArray:       mixedArray,
		NestedIntArray:   nestedInts,
	}
}

func generateComplex() ComplexTypes {
	person1 := Person{
		ID:       1,
		Name:     "Alice Johnson",
		Email:    "alice@example.com",
		Age:      30,
		Active:   true,
		Tags:     []string{"developer", "golang", "senior"},
		Metadata: map[string]string{"team": "backend", "location": "remote"},
	}
	
	person2 := Person{
		ID:       2,
		Name:     "Bob Smith",
		Email:    "bob@example.com",
		Age:      25,
		Active:   true,
		Tags:     []string{"designer", "ui", "junior"},
		Metadata: map[string]string{"team": "frontend", "location": "office"},
	}
	
	dept := Department{
		ID:       100,
		Name:     "Engineering",
		Manager:  person1,
		Budget:   500000.0,
		Projects: []string{"project-a", "project-b", "project-c"},
	}
	
	company := Company{
		ID:          1001,
		Name:        "Tech Corp",
		Founded:     2020,
		Employees:   []Person{person1, person2},
		Departments: []Department{dept},
		Revenue:     1000000.0,
	}
	
	return ComplexTypes{
		Person:      person1,
		Company:     company,
		People:      []Person{person1, person2},
		Departments: []Department{dept},
		StringMap:   map[string]string{"key1": "value1", "key2": "value2"},
		IntMap:      map[string]int32{"count": 42, "limit": 100},
	}
}

func generateEdgeCases() EdgeCaseTypes {
	// Generate large string (1MB)
	var largeString bytes.Buffer
	for i := 0; i < 1024*1024; i++ {
		largeString.WriteByte(byte('A' + (i % 26)))
	}
	
	// Generate large bytes (1MB)
	largeBytes := make([]byte, 1024*1024)
	for i := range largeBytes {
		largeBytes[i] = byte(i % 256)
	}
	
	return EdgeCaseTypes{
		EmptyString:     "",
		EmptyBytes:      []byte{},
		EmptyArray:      []int32{},
		NilBytes:        nil,
		LargeString:     largeString.String(),
		LargeBytes:      largeBytes,
		UnicodeString:   "Hello 世界 🌍 Γειά σου κόσμε ¡Hola mundo!",
		SpecialChars:    "!@#$%^&*()_+-=[]{}|;':\",./<>?`~",
		MaxInt64:        9223372036854775807,
		MinInt64:        -9223372036854775808,
		MaxUint64:       18446744073709551615,
		SmallFloat:      1e-100,
		LargeFloat:      1e100,
		NaNFloat:        math.NaN(),
		InfFloat:        math.Inf(1),
	}
}

func generateRealWorld() RealWorldTypes {
	// Generate realistic user profile
	userProfile := UserProfile{
		UserID:      12345,
		Username:    "alice_johnson",
		FullName:    "Alice Johnson",
		Email:       "alice.johnson@example.com",
		PhoneNumber: "+1-555-0123",
		BirthDate:   time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
		Preferences: map[string]interface{}{
			"theme":         "dark",
			"notifications": true,
			"language":      "en",
			"timezone":      "UTC-5",
			"font_size":     14,
		},
		Friends:     []uint64{11111, 22222, 33333, 44444},
		Groups:      []string{"developers", "golang-users", "open-source"},
		LastLogin:   time.Now().Add(-2 * time.Hour),
		IsActive:    true,
		ProfilePic:  []byte{0x89, 0x50, 0x4E, 0x47}, // PNG header bytes
	}

	// Generate configuration
	config := Configuration{
		AppName:     "EBE Benchmark",
		Version:     "1.0.0",
		Environment: "production",
		Features: map[string]bool{
			"metrics_enabled":   true,
			"debug_mode":        false,
			"cache_enabled":     true,
			"compression":       true,
			"rate_limiting":     true,
		},
		Limits: map[string]int32{
			"max_connections":   1000,
			"max_requests_per_second": 100,
			"max_file_size":     10485760, // 10MB
			"timeout_seconds":   30,
		},
		Endpoints: []string{
			"https://api.example.com/v1",
			"https://backup-api.example.com/v1",
			"https://analytics.example.com/v1",
		},
		Timeouts: map[string]int32{
			"read_timeout":  5000,
			"write_timeout": 5000,
			"idle_timeout":  60000,
		},
		Debug: false,
	}

	// Generate time series data
	timeSeries := make([]TimeSeriesPoint, 100)
	baseTime := time.Now().Add(-24 * time.Hour)
	for i := range timeSeries {
		timeSeries[i] = TimeSeriesPoint{
			Timestamp: baseTime.Add(time.Duration(i) * 15 * time.Minute),
			Value:     50.0 + 10.0*math.Sin(float64(i)*0.1) + 2.0*float64(i%10),
			Tags: map[string]string{
				"metric":     "cpu_usage",
				"host":       fmt.Sprintf("server-%02d", i%5),
				"datacenter": "us-west-1",
				"env":        "production",
			},
			Source: "monitoring-system",
		}
	}

	// Generate log entries
	logEntries := make([]LogEntry, 50)
	levels := []string{"INFO", "DEBUG", "WARN", "ERROR"}
	sources := []string{"api-server", "database", "cache", "worker"}
	for i := range logEntries {
		logEntries[i] = LogEntry{
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
			Level:     levels[i%len(levels)],
			Message:   fmt.Sprintf("Log message %d with some details", i),
			Source:    sources[i%len(sources)],
			ThreadID:  uint32(1000 + i%10),
			Data: map[string]interface{}{
				"request_id": fmt.Sprintf("req-%06d", i),
				"user_id":    uint64(10000 + i%100),
				"duration":   float64(i%1000) + 0.5,
				"success":    i%10 != 0, // 90% success rate
			},
		}
	}

	// Generate API response
	apiResponse := APIResponse{
		Status:    200,
		Message:   "Request processed successfully",
		Data:      map[string]interface{}{
			"total_count": 1234,
			"results":     []string{"item1", "item2", "item3"},
			"has_more":    true,
		},
		Timestamp: time.Now(),
		RequestID: "req-abc123def456",
		Errors:    []string{}, // No errors for successful response
		Meta: map[string]interface{}{
			"api_version": "1.2.3",
			"rate_limit_remaining": 95,
			"processing_time_ms": 45.7,
		},
	}

	return RealWorldTypes{
		UserProfile: userProfile,
		Config:      config,
		TimeSeries:  timeSeries,
		LogEntries:  logEntries,
		APIResponse: apiResponse,
	}
}