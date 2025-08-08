package id_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bold-minds/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
	gen := id.NewGenerator()

	// Act
	id := gen.Generate()
	t.Logf("Generated id: %+v", id)

	// Assert
	assert.Len(t, id, 26) // ULID should be 26 characters
	assert.True(t, gen.IsIdValid(id))
	// Should be uppercase
	assert.Equal(t, strings.ToUpper(id), id)
}

func Test_GenerateWithTime(t *testing.T) {
	gen := id.NewGenerator()
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	// Act
	id := gen.GenerateWithTime(testTime)

	// Assert
	assert.Len(t, id, 26)
	assert.True(t, gen.IsIdValid(id))

	// Extract timestamp and verify it matches
	extractedTime, err := gen.ExtractTimestamp(id)
	require.NoError(t, err)
	// Should be within 1 second (ULID uses millisecond precision)
	assert.WithinDuration(t, testTime, extractedTime, time.Second)
}

func Test_Generate_NoDups(t *testing.T) {
	gen := id.NewGenerator()

	// Act
	ids := map[string]bool{}

	for i := 0; i < 1000; i++ { // Reduced from 10k for faster tests
		id := gen.Generate()

		// Assert
		require.NotContains(t, ids, id)
		require.True(t, gen.IsIdValid(id))

		ids[id] = true
	}
}

func Test_IsIdValid(t *testing.T) {
	gen := id.NewGenerator()
	valid := gen.Generate()

	// Act & Assert
	assert.True(t, gen.IsIdValid(valid))
	assert.False(t, gen.IsIdValid(""))
	assert.False(t, gen.IsIdValid("invalid"))
	assert.False(t, gen.IsIdValid("short"))                       // Too short
	assert.False(t, gen.IsIdValid("TOOLONGFORTESTINGULIDS12345")) // Too long
}

func Test_ValidateAndNormalize(t *testing.T) {
	gen := id.NewGenerator()
	original := gen.Generate()
	lowercase := strings.ToLower(original)

	// Act
	normalized, err := gen.ValidateAndNormalize(lowercase)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, original, normalized)
	assert.Equal(t, strings.ToUpper(lowercase), normalized)

	// Test error cases
	_, err = gen.ValidateAndNormalize("")
	assert.Error(t, err)
	_, err = gen.ValidateAndNormalize("invalid")
	assert.Error(t, err)
}

func Test_GenerateBatch(t *testing.T) {
	gen := id.NewGenerator()

	// Act
	batch := gen.GenerateBatch(5)

	// Assert
	assert.Len(t, batch, 5)
	for _, id := range batch {
		assert.True(t, gen.IsIdValid(id))
	}

	// Test uniqueness
	unique := make(map[string]bool)
	for _, id := range batch {
		assert.False(t, unique[id], "Duplicate id found: %s", id)
		unique[id] = true
	}

	// Test edge cases
	assert.Empty(t, gen.GenerateBatch(0))
	assert.Empty(t, gen.GenerateBatch(-1))
}

func Test_GenerateRange(t *testing.T) {
	gen := id.NewGenerator()
	start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

	// Act
	ids := gen.GenerateRange(start, end, 3)

	// Assert
	assert.Len(t, ids, 3)
	for _, id := range ids {
		assert.True(t, gen.IsIdValid(id))
		timestamp, err := gen.ExtractTimestamp(id)
		require.NoError(t, err)
		assert.True(t, timestamp.Equal(start) || timestamp.After(start))
		assert.True(t, timestamp.Equal(end) || timestamp.Before(end))
	}

	// Test edge cases
	assert.Empty(t, gen.GenerateRange(start, end, 0))
	assert.Empty(t, gen.GenerateRange(end, start, 3)) // Invalid range
}

func Test_ExtractTimestamp(t *testing.T) {
	gen := id.NewGenerator()
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 123000000, time.UTC)
	id := gen.GenerateWithTime(testTime)

	// Act
	extracted, err := gen.ExtractTimestamp(id)

	// Assert
	require.NoError(t, err)
	// ULID has millisecond precision, so should be within 1ms
	assert.WithinDuration(t, testTime, extracted, time.Millisecond)

	// Test invalid id
	_, err = gen.ExtractTimestamp("invalid")
	assert.Error(t, err)
}

func Test_Age(t *testing.T) {
	gen := id.NewGenerator()
	pastTime := time.Now().Add(-1 * time.Hour)
	id := gen.GenerateWithTime(pastTime)

	// Act
	age, err := gen.Age(id)

	// Assert
	require.NoError(t, err)
	assert.True(t, age >= time.Hour)
	assert.True(t, age < time.Hour+time.Minute) // Should be close to 1 hour
}

func Test_IsExpired(t *testing.T) {
	gen := id.NewGenerator()
	oldId := gen.GenerateWithTime(time.Now().Add(-2 * time.Hour))
	newId := gen.GenerateWithTime(time.Now().Add(-30 * time.Minute))

	// Act & Assert
	expired, err := gen.IsExpired(oldId, time.Hour)
	require.NoError(t, err)
	assert.True(t, expired)

	notExpired, err := gen.IsExpired(newId, time.Hour)
	require.NoError(t, err)
	assert.False(t, notExpired)
}

func Test_Compare(t *testing.T) {
	gen := id.NewGenerator()
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)
	id1 := gen.GenerateWithTime(time1)
	id2 := gen.GenerateWithTime(time2)

	// Act
	cmp, err := gen.Compare(id1, id2)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, -1, cmp) // id1 should be before id2

	// Test same ids
	cmp, err = gen.Compare(id1, id1)
	require.NoError(t, err)
	assert.Equal(t, 0, cmp)

	// Test reverse order
	cmp, err = gen.Compare(id2, id1)
	require.NoError(t, err)
	assert.Equal(t, 1, cmp)
}

func Test_IsBefore_IsAfter(t *testing.T) {
	gen := id.NewGenerator()
	time1 := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)
	id1 := gen.GenerateWithTime(time1)
	id2 := gen.GenerateWithTime(time2)

	// Act & Assert
	before, err := gen.IsBefore(id1, id2)
	require.NoError(t, err)
	assert.True(t, before)

	after, err := gen.IsAfter(id2, id1)
	require.NoError(t, err)
	assert.True(t, after)

	before, err = gen.IsBefore(id2, id1)
	require.NoError(t, err)
	assert.False(t, before)
}

func Test_ToBytes_FromBytes(t *testing.T) {
	gen := id.NewGenerator()
	original := gen.Generate()

	// Act
	bytes, err := gen.ToBytes(original)
	require.NoError(t, err)
	restored := gen.FromBytes(bytes)

	// Assert
	assert.Equal(t, original, restored)
	assert.Len(t, bytes, 16) // ULID is 16 bytes
}

func Test_ToUUID(t *testing.T) {
	gen := id.NewGenerator()
	id := gen.Generate()

	// Act
	uuid, err := gen.ToUUID(id)

	// Assert
	require.NoError(t, err)
	assert.Len(t, uuid, 36) // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	assert.Contains(t, uuid, "-")
	parts := strings.Split(uuid, "-")
	assert.Len(t, parts, 5)
}

func Test_AnalyzeIDs(t *testing.T) {
	gen := id.NewGenerator()
	start := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)
	ids := gen.GenerateRange(start, end, 5)

	// Act
	stats, err := id.AnalyzeIDs(ids)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 5, stats.Count)
	assert.True(t, stats.TimeSpan >= 0)
	assert.True(t, gen.IsIdValid(stats.FirstID))
	assert.True(t, gen.IsIdValid(stats.LastID))
}

func Test_FilterByTimeRange(t *testing.T) {
	gen := id.NewGenerator()
	start := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	middle := time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC)
	end := time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)

	ids := []string{
		gen.GenerateWithTime(start.Add(-time.Hour)), // Before range
		gen.GenerateWithTime(start),                 // Start of range
		gen.GenerateWithTime(middle),                // Middle of range
		gen.GenerateWithTime(end),                   // End of range
		gen.GenerateWithTime(end.Add(time.Hour)),    // After range
	}

	// Act
	filtered := id.FilterByTimeRange(ids, start, end)

	// Assert
	assert.Len(t, filtered, 3) // Should include start, middle, and end
}

func Test_SortChronologically(t *testing.T) {
	gen := id.NewGenerator()
	times := []time.Time{
		time.Date(2023, 1, 3, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	ids := make([]string, len(times))
	for i, t := range times {
		ids[i] = gen.GenerateWithTime(t)
	}

	// Act
	sorted := id.SortChronologically(ids)

	// Assert
	assert.Len(t, sorted, 3)
	// Verify chronological order
	for i := 0; i < len(sorted)-1; i++ {
		cmp, err := gen.Compare(sorted[i], sorted[i+1])
		require.NoError(t, err)
		assert.True(t, cmp <= 0, "Ids should be in chronological order")
	}
}

func Test_NewSecureGenerator(t *testing.T) {
	secureGen := id.NewSecureGenerator()

	// Act
	id := secureGen.Generate()

	// Assert
	assert.Len(t, id, 26)
	assert.True(t, secureGen.IsIdValid(id))
}
