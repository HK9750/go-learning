package main

import (
	"testing"
)

// DEEP DIVE: Testing & Benchmarks

// The function we want to test
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 1. Unit Test
// Run: go test .
func TestSum(t *testing.T) {
	got := Sum(1, 2, 3)
	want := 6
	if got != want {
		t.Errorf("Sum(1,2,3) = %d; want %d", got, want)
	}
}

// 2. Table-Driven Tests (Idiomatic Go)
func TestSumTable(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		want    int
	}{
		{"Basic", []int{1, 2}, 3},
		{"Empty", []int{}, 0},
		{"Negative", []int{-1, -1}, -2},
		{"Mixed", []int{1, -1}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sum(tt.input...)
			if got != tt.want {
				t.Errorf("%s: got %d, want %d", tt.name, got, tt.want)
			}
		})
	}
}

// 3. Benchmarks
// Run: go test -bench=. -benchmem
// Output shows: iterations, time/op, B/op (bytes allocated per op), allocs/op
func BenchmarkSum(b *testing.B) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	// Reset timer to ignore setup costs
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		Sum(nums...)
	}
}

// 4. Fuzzing (Go 1.18+)
// Generates random inputs to find crashes.
// Run: go test -fuzz=FuzzSum
func FuzzSum(f *testing.F) {
	f.Add(1, 2) // Seed corpus

	f.Fuzz(func(t *testing.T, a int, b int) {
		// Property based testing principles:
		// Sum(a, b) should equal Sum(b, a) (Commutative)
		if Sum(a, b) != Sum(b, a) {
			t.Errorf("Commutativity failed for %d, %d", a, b)
		}
	})
}

/*
BEST PRACTICES & COMMON ERRORS:

1. Race Detector:
   Command: 'go test -race'
   Why: Detects concurrent map writes and data races during tests. Essential for concurrency code.

2. t.Parallel() Pitfall:
   Error: Loop variable capture in subtests.
   Code: for _, tc := range tests { t.Run(tc.name, func(t *testing.T) { t.Parallel(); ... }) }
   Fix: tc := tc (re-declare variable inside loop) before t.Run or inside t.Run (Go < 1.22).

3. Benchmark Accuracy:
   Warning: Compiler optimizations might effectively delete your benchmark code if result isn't used.
   Fix: Assign result to a global variable or exported field to prevent dead code elimination.
*/
