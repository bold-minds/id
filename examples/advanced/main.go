package main

import (
	"fmt"
	"time"

	"github.com/bold-minds/id"
)

func main() {
	gen := id.NewGenerator()

	fmt.Println("=== Comparison & Sorting ===")

	// Generate ULIDs at different times
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)
	time3 := time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC)

	ulid1 := gen.GenerateWithTime(time1)
	ulid2 := gen.GenerateWithTime(time2)
	ulid3 := gen.GenerateWithTime(time3)

	// Compare ULIDs
	cmp, _ := gen.Compare(ulid1, ulid2)
	fmt.Printf("Compare %s vs %s: %d\n", ulid1[:8]+"...", ulid2[:8]+"...", cmp)

	// Check chronological order
	before, _ := gen.IsBefore(ulid1, ulid2)
	fmt.Printf("Is %s before %s: %t\n", ulid1[:8]+"...", ulid2[:8]+"...", before)

	// Sort ULIDs chronologically
	unsorted := []string{ulid3, ulid1, ulid2}
	sorted := id.SortChronologically(unsorted)

	fmt.Println("Unsorted order:")
	for i, id := range unsorted {
		ts, _ := gen.ExtractTimestamp(id)
		fmt.Printf("  [%d]: %s... (%s)\n", i, id[:8], ts.Format("Jan 2"))
	}

	fmt.Println("Sorted order:")
	for i, id := range sorted {
		ts, _ := gen.ExtractTimestamp(id)
		fmt.Printf("  [%d]: %s... (%s)\n", i, id[:8], ts.Format("Jan 2"))
	}

	fmt.Println("\n=== Format Conversions ===")

	// Convert to binary and back
	bytes, _ := gen.ToBytes(ulid1)
	restored := gen.FromBytes(bytes)
	fmt.Printf("Original: %s\n", ulid1)
	fmt.Printf("Restored: %s\n", restored)
	fmt.Printf("Match: %t\n", ulid1 == restored)

	// Convert to UUID format
	uuid, _ := gen.ToUUID(ulid1)
	fmt.Printf("As UUID: %s\n", uuid)

	fmt.Println("\n=== Analytics & Statistics ===")

	// Generate a collection of ULIDs
	collection := []string{ulid1, ulid2, ulid3}
	stats, _ := id.AnalyzeIDs(collection)

	fmt.Printf("Collection stats:\n")
	fmt.Printf("  Count: %d\n", stats.Count)
	fmt.Printf("  Time span: %v\n", stats.TimeSpan)
	fmt.Printf("  First ID: %s... (%s)\n", stats.FirstID[:8], stats.FirstTime.Format("Jan 2"))
	fmt.Printf("  Last ID: %s... (%s)\n", stats.LastID[:8], stats.LastTime.Format("Jan 2"))

	fmt.Println("\n=== Time-based Filtering ===")

	// Filter by time range
	filterStart := time1.Add(12 * time.Hour)
	filterEnd := time3.Add(-12 * time.Hour)

	filtered := id.FilterByTimeRange(collection, filterStart, filterEnd)
	fmt.Printf("Filtered %d -> %d ULIDs within range\n", len(collection), len(filtered))

	for _, id := range filtered {
		ts, _ := gen.ExtractTimestamp(id)
		fmt.Printf("  %s... (%s)\n", id[:8], ts.Format("Jan 2 15:04"))
	}

	fmt.Println("\n=== Expiration Checking ===")

	// Check if ULIDs are expired
	oldULID := gen.GenerateWithTime(time.Now().Add(-2 * time.Hour))
	newULID := gen.GenerateWithTime(time.Now().Add(-30 * time.Minute))

	expired1, _ := gen.IsExpired(oldULID, time.Hour)
	expired2, _ := gen.IsExpired(newULID, time.Hour)

	fmt.Printf("Old ULID expired (1h limit): %t\n", expired1)
	fmt.Printf("New ULID expired (1h limit): %t\n", expired2)

	fmt.Println("\n=== Secure Generation ===")

	// Use cryptographically secure generator
	secureGen := id.NewSecureGenerator()
	secureULID := secureGen.Generate()
	fmt.Printf("Secure ULID: %s\n", secureULID)
}
