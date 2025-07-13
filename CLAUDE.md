# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is EBE (Efficient Binary Encoding), a Go library for compact binary serialization with support for multiple data types including integers, floats, strings, buffers, arrays, maps, and structs. The library uses a type-aware binary format with optimized encoding for small values and includes performance-optimized fast paths for common data types.

## Project Structure

- **Main EBE Library**: `/Users/craigsymonds/dev/ebe/library/` - Core serialization library (no external dependencies)
- **Comparison Tool**: `/Users/craigsymonds/dev/ebe/` - Standalone application comparing EBE vs Protocol Buffers performance

## Development Commands

### Testing
```bash
go test ./...                    # Run all tests
go test ./test                   # Run specific test package
go test -v ./test               # Run tests with verbose output
go test -run TestUNibbles ./test # Run specific test function
```

### Building
```bash
go build ./...                  # Build all packages
go mod tidy                     # Clean up dependencies
```

### Module Management
```bash
go mod init ebe-library         # Initialize module (already done)
go mod download                 # Download dependencies
```

## Architecture

### Core Components

**serialize/** - Main serialization/deserialization logic
- `Serialize.go` - Main entry point for serialization, handles type detection and routing
- `Deserialize.go` - Main entry point for deserialization with type-safe reflection
- `Map.go` - Map serialization/deserialization with fast paths for common types (map[string]int, map[string]string, etc.)
- Type-specific files: `Uint.go`, `Sint.go`, `Float.go`, `String.go`, `Buffer.go`, `Array.go`, `Struct.go`, `Boolean.go`, `Json.go`

**types/** - Type system and constants
- `types.go` - Defines the Types enum (UNibble, SNibble, SInt, UInt, Float, Boolean, String, Buffer, Array, Map, Json)
- `header.go` - Header byte manipulation utilities

**utils/** - Shared utilities
- `utils.go` - Byte I/O helpers, value comparison, type conversion, debug printing

**test/** - Comprehensive test suite
- `value_test.go` - Core serialization/deserialization tests for all types
- `map_test.go`, `map_integration_test.go`, `map_benchmark_test.go`, `map_performance_test.go` - Map serialization tests and benchmarks
- `array_test.go`, `struct_test.go`, `json_test.go`, `multi_value_test.go` - Specialized tests

### Binary Format Design

The library uses a compact binary format where the first 4 bits of each value's header indicate the type, and the remaining 4 bits store either the value itself (for small types like UNibble/SNibble) or metadata like length.

**Type Optimization:**
- UNibble (0-15): Value stored directly in header
- SNibble (-7 to 0): Signed values with sign bit and magnitude
- UInt/SInt: Variable-length encoding based on value size
- String/Buffer: Length-prefixed with overflow handling for large sizes

### Key Patterns

1. **Type Detection**: `Serialize()` uses Go's type system and reflection to route values to appropriate serializers
2. **Reflection-based Deserialization**: `Deserialize()` takes a pointer to the target type and uses reflection for type-safe assignment
3. **Value Conversion**: `utils.SetValueWithConversion()` handles type conversions between compatible types
4. **Streaming I/O**: All operations work with `io.Reader`/`io.Writer` interfaces for memory efficiency

### Testing Strategy

Tests are organized by value type and include:
- Boundary value testing (min/max for each type)
- Type conversion validation 
- Round-trip serialization/deserialization verification
- Special value handling (NaN, infinity for floats)
- Unicode and binary data support for strings/buffers