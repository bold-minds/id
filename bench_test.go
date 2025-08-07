package id_test

import (
	"testing"
	"time"

	"github.com/bold-minds/id"
)

func BenchmarkGenerate(b *testing.B) {
	gen := id.NewGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}
}

func BenchmarkGenerateSecure(b *testing.B) {
	gen := id.NewSecureGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.Generate()
	}
}

func BenchmarkGenerateBatch(b *testing.B) {
	gen := id.NewGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.GenerateBatch(100)
	}
}

func BenchmarkIsKeyValid(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.IsKeyValid(ulid)
	}
}

func BenchmarkExtractTimestamp(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gen.ExtractTimestamp(ulid)
	}
}

func BenchmarkAge(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gen.Age(ulid)
	}
}

func BenchmarkCompare(b *testing.B) {
	gen := id.NewGenerator()
	ulid1 := gen.Generate()
	ulid2 := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gen.Compare(ulid1, ulid2)
	}
}

func BenchmarkToBytes(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gen.ToBytes(ulid)
	}
}

func BenchmarkFromBytes(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	bytes, _ := gen.ToBytes(ulid)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = gen.FromBytes(bytes)
	}
}

func BenchmarkToUUID(b *testing.B) {
	gen := id.NewGenerator()
	ulid := gen.Generate()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = gen.ToUUID(ulid)
	}
}

func BenchmarkSortChronologically(b *testing.B) {
	gen := id.NewGenerator()
	ulids := make([]string, 1000)
	for i := range ulids {
		ulids[i] = gen.Generate()
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.SortChronologically(ulids)
	}
}

func BenchmarkAnalyzeIDs(b *testing.B) {
	gen := id.NewGenerator()
	ulids := make([]string, 100)
	for i := range ulids {
		ulids[i] = gen.Generate()
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = id.AnalyzeIDs(ulids)
	}
}

func BenchmarkFilterByTimeRange(b *testing.B) {
	gen := id.NewGenerator()
	start := time.Now().Add(-time.Hour)

	ulids := make([]string, 1000)
	for i := range ulids {
		ulids[i] = gen.GenerateWithTime(start.Add(time.Duration(i) * time.Minute))
	}

	filterStart := start.Add(15 * time.Minute)
	filterEnd := start.Add(45 * time.Minute)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.FilterByTimeRange(ulids, filterStart, filterEnd)
	}
}
