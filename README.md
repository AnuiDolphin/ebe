# EBE (Efficient Binary Encoding)

A comprehensive project for efficient binary serialization in Go, including both the core library and performance comparison tools.

## Project Structure

```
ebe/
├── library/           # Core EBE serialization library
│   ├── serialize/     # Main serialization/deserialization logic
│   ├── types/         # Type system and constants
│   ├── utils/         # Shared utilities
│   ├── test/          # Comprehensive test suite
│   └── CLAUDE.md      # Development documentation
└── compare/           # EBE vs Protocol Buffers comparison tool
    ├── proto/         # Protocol Buffer schemas
    ├── benchmark.go   # Benchmarking framework
    ├── testdata.go    # Test data generation
    └── README.md      # Comparison tool documentation
```

## Quick Start

### Using the EBE Library

```bash
cd library
go test ./test/        # Run all tests
```

### Running Performance Comparisons

```bash
go build              # Build comparison tool
./ebe-compare          # Run all comparisons
./ebe-compare -group primitives -output csv > results.csv
```

## Core Features

- **Compact Binary Format**: Efficient encoding with type-aware headers
- **Zero External Dependencies**: Core library has no dependencies
- **Fast Path Optimizations**: 49% faster serialization and 73% faster deserialization for common map types
- **Comprehensive Type Support**: Integers, floats, strings, buffers, arrays, maps, structs
- **Streaming I/O**: Works with `io.Reader`/`io.Writer` interfaces

## Performance Highlights

Recent benchmarks show EBE achieves:
- **14.2% smaller** serialized size for user profiles vs Protocol Buffers
- **4.8% smaller** serialized size for configurations vs Protocol Buffers
- **Fast path optimizations** for common data patterns reduce reflection overhead

## Documentation

- **Library Documentation**: See `library/CLAUDE.md` for development guidance
- **Comparison Tool**: See comparison tool README for usage instructions
- **Architecture**: Core serialization library with modular design
- **Testing**: Comprehensive test suite with benchmark comparisons