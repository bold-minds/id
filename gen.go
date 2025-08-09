package id

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	mathrand "math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

var (
	entropyMu sync.Mutex
	// Default entropy uses math/rand for performance. Use NewSecureGenerator() for crypto-secure randomness.
	entropy = ulid.Monotonic(mathrand.New(mathrand.NewSource(time.Now().UnixNano())), 0) //nolint:gosec // G404: Intentional use of math/rand for performance; crypto/rand available via NewSecureGenerator()
)

// Core id generation
type Generator interface {
	Generate() string
	IsIdValid(string) bool
}

// Id generation for batch operations
type Batcher interface {
	Generator
	GenerateWithTime(t time.Time) string
	GenerateBatch(count int) []string
	GenerateRange(start, end time.Time, count int) []string
}

// Validation and normalization
type Validator interface {
	IsIdValid(string) bool
	ValidateAndNormalize(id string) (string, error)
}

// Time-based operations
type Timestamper interface {
	ExtractTimestamp(id string) (time.Time, error)
	Age(id string) (time.Duration, error)
	IsExpired(id string, maxAge time.Duration) (bool, error)
}

// Comparison operations
type Comparator interface {
	Compare(id1, id2 string) (int, error)
	IsBefore(id1, id2 string) (bool, error)
	IsAfter(id1, id2 string) (bool, error)
}

// Format conversions
type Converter interface {
	ToBytes(id string) ([16]byte, error)
	FromBytes(data [16]byte) string
	ToUUID(id string) (string, error)
}

// Composite interface with everything
type Provider interface {
	Generator
	Batcher
	Validator
	Timestamper
	Comparator
	Converter
}

// generator ensures valid ids for records
type generator struct {
	entropySource io.Reader
}

// NewGenerator creates a new generator with default entropy
func NewGenerator() *generator {
	return &generator{
		entropySource: entropy,
	}
}

// NewGeneratorWithEntropy creates a generator with custom entropy source
func NewGeneratorWithEntropy(entropySource io.Reader) *generator {
	return &generator{
		entropySource: entropySource,
	}
}

// NewSecureGenerator creates a generator using crypto/rand for high-security scenarios
func NewSecureGenerator() *generator {
	return &generator{
		entropySource: rand.Reader,
	}
}

// Basic Generation Methods

// Generate provides a new globally unique URL safe id for a record
func (g *generator) Generate() string {
	return g.GenerateWithTime(time.Now())
}

// GenerateWithTime generates a ULID with a specific timestamp
func (g *generator) GenerateWithTime(t time.Time) string {
	entropyMu.Lock()
	defer entropyMu.Unlock()
	id := ulid.MustNew(ulid.Timestamp(t), g.entropySource)
	return id.String()
}

// GenerateBatch creates multiple ULIDs efficiently
func (g *generator) GenerateBatch(count int) []string {
	if count <= 0 {
		return []string{}
	}

	result := make([]string, count)
	entropyMu.Lock()
	defer entropyMu.Unlock()

	for i := 0; i < count; i++ {
		id := ulid.MustNew(ulid.Timestamp(time.Now()), g.entropySource)
		result[i] = id.String()
	}
	return result
}

// GenerateRange creates ULIDs within a time range
func (g *generator) GenerateRange(start, end time.Time, count int) []string {
	if count <= 0 || end.Before(start) {
		return []string{}
	}

	result := make([]string, count)
	duration := end.Sub(start)
	entropyMu.Lock()
	defer entropyMu.Unlock()

	for i := 0; i < count; i++ {
		// Distribute timestamps evenly across the range
		offset := time.Duration(int64(duration) * int64(i) / int64(count))
		timestamp := start.Add(offset)
		id := ulid.MustNew(ulid.Timestamp(timestamp), g.entropySource)
		result[i] = id.String()
	}
	return result
}

// Validation Methods

// IsIdValid validates that the provided id is a valid ULID
func (g *generator) IsIdValid(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// ValidateAndNormalize checks and normalizes a ULID string
func (g *generator) ValidateAndNormalize(id string) (string, error) {
	if id == "" {
		return "", errors.New("empty ULID string")
	}

	// Normalize case (ULIDs should be uppercase)
	normalized := strings.ToUpper(id)

	// Validate the normalized ULID
	parsed, err := ulid.Parse(normalized)
	if err != nil {
		return "", fmt.Errorf("invalid ULID: %w", err)
	}

	return parsed.String(), nil
}

// Timestamp Operations

// ExtractTimestamp returns the timestamp component of a ULID
func (g *generator) ExtractTimestamp(id string) (time.Time, error) {
	parsed, err := ulid.Parse(id)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid ULID: %w", err)
	}

	timestamp := parsed.Time()
	// ULID timestamp is milliseconds since Unix epoch
	// Safe conversion to avoid integer overflow (gosec G115)
	timestampMs := timestamp
	const maxInt64 = uint64(9223372036854775807) // math.MaxInt64 as uint64
	if timestampMs > maxInt64 {
		return time.Time{}, fmt.Errorf("timestamp too large: %d", timestampMs)
	}
	seconds := int64(timestampMs) / 1000
	nanoseconds := (int64(timestampMs) % 1000) * int64(time.Millisecond)
	return time.Unix(seconds, nanoseconds), nil
}

// Age returns how old a ULID is
func (g *generator) Age(id string) (time.Duration, error) {
	timestamp, err := g.ExtractTimestamp(id)
	if err != nil {
		return 0, err
	}

	return time.Since(timestamp), nil
}

// IsExpired checks if ULID is older than maxAge
func (g *generator) IsExpired(id string, maxAge time.Duration) (bool, error) {
	age, err := g.Age(id)
	if err != nil {
		return false, err
	}

	return age > maxAge, nil
}

// Comparison Operations

// Compare returns -1, 0, or 1 for chronological ordering
func (g *generator) Compare(id1, id2 string) (int, error) {
	ulid1, err := ulid.Parse(id1)
	if err != nil {
		return 0, fmt.Errorf("invalid first ULID: %w", err)
	}

	ulid2, err := ulid.Parse(id2)
	if err != nil {
		return 0, fmt.Errorf("invalid second ULID: %w", err)
	}

	return ulid1.Compare(ulid2), nil
}

// IsBefore checks if id1 was generated before id2
func (g *generator) IsBefore(id1, id2 string) (bool, error) {
	cmp, err := g.Compare(id1, id2)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

// IsAfter checks if id1 was generated after id2
func (g *generator) IsAfter(id1, id2 string) (bool, error) {
	cmp, err := g.Compare(id1, id2)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

// Format Conversions

// ToBytes returns the binary representation of a ULID
func (g *generator) ToBytes(id string) ([16]byte, error) {
	parsed, err := ulid.Parse(id)
	if err != nil {
		return [16]byte{}, fmt.Errorf("invalid ULID: %w", err)
	}

	// Convert ULID to byte array
	var result [16]byte
	copy(result[:], parsed[:])
	return result, nil
}

// FromBytes creates ULID string from binary representation
func (g *generator) FromBytes(data [16]byte) string {
	var u ulid.ULID
	copy(u[:], data[:])
	return u.String()
}

// ToUUID converts ULID to UUID format (for compatibility)
func (g *generator) ToUUID(id string) (string, error) {
	bytes, err := g.ToBytes(id)
	if err != nil {
		return "", err
	}

	// Format as UUID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16]), nil
}

// Utility Functions

// Stats provides statistics about a collection of ULIDs
type Stats struct {
	Count     int
	TimeSpan  time.Duration
	FirstID   string
	LastID    string
	FirstTime time.Time
	LastTime  time.Time
}

// AnalyzeIDs provides generation statistics for a slice of ULIDs
func AnalyzeIDs(ids []string) (Stats, error) {
	if len(ids) == 0 {
		return Stats{}, nil
	}

	// Create a temporary generator for parsing
	g := NewGenerator()

	// Parse all timestamps
	timestamps := make([]time.Time, 0, len(ids))
	validIDs := make([]string, 0, len(ids))

	for _, id := range ids {
		if timestamp, err := g.ExtractTimestamp(id); err == nil {
			timestamps = append(timestamps, timestamp)
			validIDs = append(validIDs, id)
		}
	}

	if len(validIDs) == 0 {
		return Stats{}, errors.New("no valid ULIDs found")
	}

	// Sort by timestamp to find first and last
	sort.Slice(validIDs, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	firstTime := timestamps[0]
	lastTime := timestamps[len(timestamps)-1]

	return Stats{
		Count:     len(validIDs),
		TimeSpan:  lastTime.Sub(firstTime),
		FirstID:   validIDs[0],
		LastID:    validIDs[len(validIDs)-1],
		FirstTime: firstTime,
		LastTime:  lastTime,
	}, nil
}

// FilterByTimeRange filters ULIDs within time bounds
func FilterByTimeRange(ids []string, start, end time.Time) []string {
	g := NewGenerator()
	result := make([]string, 0, len(ids))

	for _, id := range ids {
		if timestamp, err := g.ExtractTimestamp(id); err == nil {
			if (timestamp.Equal(start) || timestamp.After(start)) &&
				(timestamp.Equal(end) || timestamp.Before(end)) {
				result = append(result, id)
			}
		}
	}

	return result
}

// SortChronologically sorts ULIDs by their timestamp component
func SortChronologically(ids []string) []string {
	if len(ids) <= 1 {
		return ids
	}

	g := NewGenerator()
	result := make([]string, len(ids))
	copy(result, ids)

	sort.Slice(result, func(i, j int) bool {
		cmp, err := g.Compare(result[i], result[j])
		if err != nil {
			return false // Keep original order if comparison fails
		}
		return cmp < 0
	})

	return result
}

// SortChronologicallyReverse sorts ULIDs by timestamp in reverse order (newest first)
func SortChronologicallyReverse(ids []string) []string {
	sorted := SortChronologically(ids)

	// Reverse the slice
	for i := 0; i < len(sorted)/2; i++ {
		j := len(sorted) - 1 - i
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}

	return sorted
}
