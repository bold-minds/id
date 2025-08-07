package main

import (
	"fmt"
	"time"

	"github.com/bold-minds/id"
)

func main() {
	// Create a generator
	gen := id.NewGenerator()

	fmt.Println("=== Basic ULID Generation ===")

	// Generate a ULID
	ulid := gen.Generate()
	fmt.Printf("Generated ULID: %s\n", ulid)

	// Extract timestamp
	timestamp, err := gen.ExtractTimestamp(ulid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created at: %s\n", timestamp.Format(time.RFC3339))

	// Check age
	age, err := gen.Age(ulid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Age: %v\n", age)

	fmt.Println("\n=== Batch Generation ===")

	// Generate multiple ULIDs efficiently
	batch := gen.GenerateBatch(5)
	for i, id := range batch {
		fmt.Printf("Batch[%d]: %s\n", i, id)
	}

	fmt.Println("\n=== Time Range Generation ===")

	// Generate ULIDs within a specific time range
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC)
	rangeULIDs := gen.GenerateRange(start, end, 3)

	for i, id := range rangeULIDs {
		ts, _ := gen.ExtractTimestamp(id)
		fmt.Printf("Range[%d]: %s (created: %s)\n", i, id, ts.Format("15:04:05"))
	}

	fmt.Println("\n=== Validation & Normalization ===")

	// Validate ULIDs
	fmt.Printf("Valid ULID: %t\n", gen.IsKeyValid(ulid))
	fmt.Printf("Invalid ULID: %t\n", gen.IsKeyValid("invalid"))

	// Normalize case
	lowercase := "01k23cg6gn6xgjwz1bd7wh3zg5"
	normalized, err := gen.ValidateAndNormalize(lowercase)
	if err != nil {
		fmt.Printf("Normalization failed: %v\n", err)
	} else {
		fmt.Printf("Normalized: %s -> %s\n", lowercase, normalized)
	}
}
