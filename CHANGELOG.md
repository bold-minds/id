# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-08-07

### Added
- Comprehensive ULID generation and manipulation library
- Basic generation with `Generate()` and `GenerateWithTime()`
- Batch operations with `GenerateBatch()` and `GenerateRange()`
- Advanced validation with `IsKeyValid()` and `ValidateAndNormalize()`
- Timestamp operations: `ExtractTimestamp()`, `Age()`, `IsExpired()`
- Comparison operations: `Compare()`, `IsBefore()`, `IsAfter()`
- Format conversions: `ToBytes()`, `FromBytes()`, `ToUUID()`
- Utility functions: `AnalyzeIDs()`, `FilterByTimeRange()`, `SortChronologically()`
- Security options: `NewSecureGenerator()`, `NewGeneratorWithEntropy()`
- Comprehensive test suite with 15+ test cases
- Performance benchmarks for all major operations
- Complete documentation and examples

### Performance
- Per-generator entropy sources eliminate global mutex contention
- Optimized batch generation reduces allocation overhead
- Smart comparison operations leverage ULID's natural ordering
- Efficient utility functions for common operations

### Documentation
- Comprehensive README with usage examples
- Contributing guidelines and code of conduct
- MIT license
- API reference documentation
- Basic and advanced usage examples

