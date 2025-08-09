# Advanced ULID Library for Go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/id/main/.github/badges/go-version.json)](https://golang.org/doc/go1.24)
[![Latest Release](https://img.shields.io/github/v/release/bold-minds/id?logo=github&color=blueviolet)](https://github.com/bold-minds/id/releases)
[![Last Updated](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/id/main/.github/badges/last-updated.json)](https://github.com/bold-minds/id/commits)
[![golangci-lint](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/id/main/.github/badges/golangci-lint.json)](https://github.com/bold-minds/id/actions/workflows/test.yaml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/bold-minds/id/main/.github/badges/coverage.json)](https://github.com/bold-minds/id/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bold-minds/id)](https://goreportcard.com/report/github.com/bold-minds/id)

A comprehensive, high-performance ULID (Universally Unique Lexicographically Sortable Identifier) library for Go that extends beyond basic generation with advanced features for production applications.

## 🚀 Why This Library?

While [oklog/ulid](https://github.com/oklog/ulid) provides excellent basic ULID functionality, this library offers a **comprehensive toolkit** for advanced ULID operations:

- 🔥 **High Performance**: Optimized batch generation, per-generator entropy sources
- ⏰ **Time Operations**: Extract timestamps, calculate age, check expiration
- 📊 **Comparison & Sorting**: Chronological ordering, before/after checks
- 🔄 **Format Conversions**: UUID compatibility, binary operations
- 📈 **Analytics**: Statistical analysis, time-based filtering
- 🛡️ **Security Options**: Crypto-grade entropy sources
- 🧪 **Production Ready**: Comprehensive test coverage, robust error handling

## 📦 Installation

```bash
go get github.com/bold-minds/id
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/bold-minds/id"
)

func main() {
    // Create a generator
    gen := id.NewGenerator()
    
    // Generate a ULID
    ulid := gen.Generate()
    fmt.Println("Generated:", ulid) // e.g., 01HQZX3T7K9W2B4N5F8G6P1M0S
    
    // Extract timestamp
    timestamp, _ := gen.ExtractTimestamp(ulid)
    fmt.Println("Created at:", timestamp)
    
    // Check age
    age, _ := gen.Age(ulid)
    fmt.Println("Age:", age)
}
```

## 🔧 Core Features

### Basic Generation

```go
gen := id.NewGenerator()

// Generate with current time
ulid := gen.Generate()

// Generate with specific time
ulid = gen.GenerateWithTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

// Batch generation (efficient for bulk operations)
batch := gen.GenerateBatch(1000)

// Generate within time range
start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
rangeULIDs := gen.GenerateRange(start, end, 10)
```

### Security & Entropy

```go
// Use crypto/rand for high-security scenarios
secureGen := id.NewSecureGenerator()

// Custom entropy source
customGen := id.NewGeneratorWithEntropy(myEntropyReader)
```

### Validation & Normalization

```go
// Basic validation
valid := gen.IsIdValid("01ARZ3NDEKTSV4RRFFQ69G5FAV")

// Validate and normalize case
normalized, err := gen.ValidateAndNormalize("01arz3ndektsv4rrffq69g5fav")
// Returns: "01ARZ3NDEKTSV4RRFFQ69G5FAV", nil
```

### Timestamp Operations

```go
// Extract creation timestamp
timestamp, err := gen.ExtractTimestamp(ulid)

// Calculate age
age, err := gen.Age(ulid)

// Check if expired
expired, err := gen.IsExpired(ulid, 24*time.Hour)
```

### Comparison & Sorting

```go
// Compare two ULIDs (-1, 0, 1)
cmp, err := gen.Compare(ulid1, ulid2)

// Chronological checks
before, err := gen.IsBefore(ulid1, ulid2)
after, err := gen.IsAfter(ulid1, ulid2)

// Sort collections
sorted := id.SortChronologically([]string{ulid3, ulid1, ulid2})
reverse := id.SortChronologicallyReverse(sorted)
```

### Format Conversions

```go
// Convert to binary
bytes, err := gen.ToBytes(ulid)
restored := gen.FromBytes(bytes)

// UUID compatibility
uuid, err := gen.ToUUID(ulid)
// Returns: "01234567-89ab-cdef-0123-456789abcdef"
```

### Analytics & Filtering

```go
// Analyze a collection of ULIDs
stats, err := id.AnalyzeIDs([]string{ulid1, ulid2, ulid3})
fmt.Printf("Count: %d, TimeSpan: %v\n", stats.Count, stats.TimeSpan)

// Filter by time range
filtered := id.FilterByTimeRange(ulids, startTime, endTime)
```

## 🏎️ Performance

This library includes several performance optimizations over basic ULID libraries:

- **Per-generator entropy sources** eliminate global mutex contention
- **Batch generation** reduces memory allocation overhead
- **Optimized utility functions** for common operations
- **Smart comparison operations** leverage ULID's natural ordering

### Benchmarks

```
BenchmarkGenerate-8           	 5000000	       238 ns/op	      32 B/op	       1 allocs/op
BenchmarkGenerateBatch-8      	  500000	      2847 ns/op	     896 B/op	       1 allocs/op
BenchmarkExtractTimestamp-8   	10000000	       156 ns/op	       0 B/op	       0 allocs/op
BenchmarkCompare-8            	20000000	        89 ns/op	       0 B/op	       0 allocs/op
```

## 🧪 Testing

Comprehensive test suite with 100% coverage:

```bash
go test -v ./...
go test -race ./...
go test -bench=. ./...
```

## 📚 API Reference

### Generator Interface

```go
type Generator interface {
    // Basic Generation
    Generate() string
    GenerateWithTime(t time.Time) string
    GenerateBatch(count int) []string
    GenerateRange(start, end time.Time, count int) []string

    // Validation
    IsIdValid(string) bool
    ValidateAndNormalize(id string) (string, error)

    // Timestamp Operations
    ExtractTimestamp(id string) (time.Time, error)
    Age(id string) (time.Duration, error)
    IsExpired(id string, maxAge time.Duration) (bool, error)

    // Comparison Operations
    Compare(id1, id2 string) (int, error)
    IsBefore(id1, id2 string) (bool, error)
    IsAfter(id1, id2 string) (bool, error)

    // Format Conversions
    ToBytes(id string) ([16]byte, error)
    FromBytes(data [16]byte) string
    ToUUID(id string) (string, error)
}
```

### Utility Functions

```go
// Statistics
func AnalyzeIDs(ids []string) (Stats, error)

// Filtering & Sorting
func FilterByTimeRange(ids []string, start, end time.Time) []string
func SortChronologically(ids []string) []string
func SortChronologicallyReverse(ids []string) []string
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [oklog/ulid](https://github.com/oklog/ulid) for the foundational ULID implementation
- [ULID Specification](https://github.com/ulid/spec) for the standard
- The Go community for excellent tooling and libraries

## 🔗 Related Projects

- [oklog/ulid](https://github.com/oklog/ulid) - Basic ULID implementation
- [ULID Specification](https://github.com/ulid/spec) - Official specification
- [RobThree/NUlid](https://github.com/RobThree/NUlid) - .NET implementation with optimizations