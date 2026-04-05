// Package examples_test provides runnable examples for the bold-minds/id
// library. Every example here is compiled and executed by `go test` and
// rendered on pkg.go.dev, so they double as living documentation.
//
// ULID generation is non-deterministic by design (random entropy), so the
// // Output: assertions here only cover deterministic aspects: lengths,
// round-trips, comparisons, counts, and extracted timestamps from fixed
// input times. The actual ULID string values are never asserted.
//
// Example functions are named after the package-level constructor
// (NewGenerator / NewSecureGenerator) with a lowercase `_suffix` tag that
// picks out the behavior being demonstrated, because the concrete
// generator type is unexported and cannot be named directly in example
// function identifiers.
package examples_test

import (
	"fmt"
	"time"

	"github.com/bold-minds/id"
)

// Fixed timestamps used across examples so that comparison, sort, filter,
// and extraction results are deterministic.
var (
	t1 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
)

// ExampleNewGenerator shows the fast default generator. The output is a
// 26-character Crockford Base32 ULID.
func ExampleNewGenerator() {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	fmt.Println(len(ulid))
	// Output: 26
}

// ExampleNewSecureGenerator shows the cryptographically-secure generator,
// appropriate for user-visible identifiers.
func ExampleNewSecureGenerator() {
	gen := id.NewSecureGenerator()
	ulid := gen.Generate()
	fmt.Println(len(ulid))
	// Output: 26
}

// ExampleNewGenerator_generateWithTime shows deterministic timestamp
// extraction — the timestamp embedded in the ULID round-trips losslessly
// at millisecond granularity.
func ExampleNewGenerator_generateWithTime() {
	gen := id.NewGenerator()
	ulid := gen.GenerateWithTime(t1)
	extracted, _ := gen.ExtractTimestamp(ulid)
	fmt.Println(extracted.UTC().Format(time.RFC3339))
	// Output: 2024-01-01T12:00:00Z
}

// ExampleNewGenerator_generateBatch shows high-throughput batch generation.
// Every ID in the batch is valid and unique.
func ExampleNewGenerator_generateBatch() {
	gen := id.NewGenerator()
	batch := gen.GenerateBatch(5)
	fmt.Println(len(batch))
	// Output: 5
}

// ExampleNewGenerator_generateRange produces a specified number of ULIDs
// spaced across a time window. Useful for backfilling historical data or
// generating test fixtures spanning a period.
func ExampleNewGenerator_generateRange() {
	gen := id.NewGenerator()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC)
	ulids := gen.GenerateRange(start, end, 3)
	fmt.Println(len(ulids))
	// Output: 3
}

// ExampleNewGenerator_isIdValid shows validation — accepts only well-formed
// Crockford Base32 ULIDs of the exact expected length.
func ExampleNewGenerator_isIdValid() {
	gen := id.NewGenerator()
	valid := gen.Generate()
	fmt.Println(gen.IsIdValid(valid))
	fmt.Println(gen.IsIdValid("not a ulid"))
	fmt.Println(gen.IsIdValid(""))
	// Output:
	// true
	// false
	// false
}

// ExampleNewGenerator_validateAndNormalize accepts mixed-case input and
// returns the canonical uppercase Crockford Base32 form.
func ExampleNewGenerator_validateAndNormalize() {
	gen := id.NewGenerator()
	lower := "01arz3ndektsv4rrffq69g5fav"
	normalized, err := gen.ValidateAndNormalize(lower)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(normalized)
	// Output: 01ARZ3NDEKTSV4RRFFQ69G5FAV
}

// ExampleNewGenerator_compare returns -1, 0, or 1 based on chronological
// order.
func ExampleNewGenerator_compare() {
	gen := id.NewGenerator()
	earlier := gen.GenerateWithTime(t1)
	later := gen.GenerateWithTime(t2)
	cmp, _ := gen.Compare(earlier, later)
	fmt.Println(cmp)
	// Output: -1
}

// ExampleNewGenerator_isBefore is the outcome-named equivalent of
// Compare < 0.
func ExampleNewGenerator_isBefore() {
	gen := id.NewGenerator()
	earlier := gen.GenerateWithTime(t1)
	later := gen.GenerateWithTime(t2)
	before, _ := gen.IsBefore(earlier, later)
	fmt.Println(before)
	// Output: true
}

// ExampleNewGenerator_isAfter is the symmetric counterpart to IsBefore.
func ExampleNewGenerator_isAfter() {
	gen := id.NewGenerator()
	earlier := gen.GenerateWithTime(t1)
	later := gen.GenerateWithTime(t2)
	after, _ := gen.IsAfter(later, earlier)
	fmt.Println(after)
	// Output: true
}

// ExampleNewGenerator_toBytes shows lossless conversion between the
// 26-character string form and the compact 16-byte representation.
func ExampleNewGenerator_toBytes() {
	gen := id.NewGenerator()
	original := gen.Generate()
	bytes, _ := gen.ToBytes(original)
	restored := gen.FromBytes(bytes)
	fmt.Println(original == restored)
	// Output: true
}

// ExampleNewGenerator_toUUID renders a ULID as a canonical 36-character
// UUID-shaped string (hyphens included). Useful for emitting IDs into
// systems that expect UUID format.
func ExampleNewGenerator_toUUID() {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	uuid, _ := gen.ToUUID(ulid)
	fmt.Println(len(uuid))
	// Output: 36
}

// ExampleSortChronologically re-orders a slice by embedded timestamp,
// ascending.
func ExampleSortChronologically() {
	gen := id.NewGenerator()
	a := gen.GenerateWithTime(t1)
	b := gen.GenerateWithTime(t2)
	c := gen.GenerateWithTime(t3)

	sorted := id.SortChronologically([]string{c, a, b})
	for _, ulid := range sorted {
		ts, _ := gen.ExtractTimestamp(ulid)
		fmt.Println(ts.UTC().Format("2006-01-02"))
	}
	// Output:
	// 2024-01-01
	// 2024-02-01
	// 2024-03-01
}

// ExampleSortChronologicallyReverse is the descending-order variant.
func ExampleSortChronologicallyReverse() {
	gen := id.NewGenerator()
	a := gen.GenerateWithTime(t1)
	b := gen.GenerateWithTime(t2)
	c := gen.GenerateWithTime(t3)

	sorted := id.SortChronologicallyReverse([]string{a, c, b})
	for _, ulid := range sorted {
		ts, _ := gen.ExtractTimestamp(ulid)
		fmt.Println(ts.UTC().Format("2006-01-02"))
	}
	// Output:
	// 2024-03-01
	// 2024-02-01
	// 2024-01-01
}

// ExampleAnalyzeIDs reports count and time span across a collection.
func ExampleAnalyzeIDs() {
	gen := id.NewGenerator()
	ids := []string{
		gen.GenerateWithTime(t1),
		gen.GenerateWithTime(t2),
		gen.GenerateWithTime(t3),
	}
	stats, _ := id.AnalyzeIDs(ids)
	fmt.Println(stats.Count)
	fmt.Println(stats.TimeSpan)
	// Output:
	// 3
	// 1440h0m0s
}

// ExampleFilterByTimeRange keeps only the ids whose embedded timestamp
// falls within [start, end].
func ExampleFilterByTimeRange() {
	gen := id.NewGenerator()
	ids := []string{
		gen.GenerateWithTime(t1),
		gen.GenerateWithTime(t2),
		gen.GenerateWithTime(t3),
	}
	start := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 28, 23, 59, 59, 0, time.UTC)
	filtered := id.FilterByTimeRange(ids, start, end)
	fmt.Println(len(filtered))
	// Output: 1
}
