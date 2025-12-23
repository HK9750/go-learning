package main

import (
	"fmt"
	"runtime"
)

// DEEP DIVE: Optimizations & Runtime Tuning

// 1. GOGC (Garbage Collector Pacing)
// Environment variable GOGC (default 100).
// Controls aggressive GC.
// GOGC=100 -> GC triggers when heap grows by 100% since last GC.
// GOGC=off -> Disable GC.
// Tuning: Raise GOGC (e.g., 200, 500) to trade RAM for CPU (fewer GC cycles).

// 2. Memory Ballast (Legacy trick)
// Allocating a huge byte array (e.g. 10GB) at startup to trick GC into running less often.
// Kept on heap but never accessed = OS doesn't use physical RAM (virtual only).
// Modern alternative: GOMEMLIMIT (Go 1.19+).

func main() {
	// 3. Stack vs Heap Escape Analysis check
	// Run: go build -gcflags="-m -l" main.go
	// Look for "escapes to heap".
	
	// 4. Trace Tool
	// Captures ms-level events (GC, Scheduler, Syscalls).
	// usage:
	// f, _ := os.Create("trace.out")
	// trace.Start(f)
	// defer trace.Stop()
	// ... workload ...
	// go tool trace trace.out

	fmt.Println("Simulating Allocations to trigger GC...")
	
	printMemStats()

	// Create garbage
	for i := 0; i < 10000; i++ {
		_ = make([]byte, 1024) // 1KB garbage
	}
	
	// Force GC
	runtime.GC()
	printMemStats()

	// 5. Zero-Allocation Tricks
	// - Use buffers (sync.Pool)
	// - Use scratch buffers for serialization
	// - Use strconv.AppendInt instead of fmt.Sprintf
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v KiB\tTotalAlloc = %v KiB\tSys = %v KiB\tNumGC = %v\n",
		m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
}

/*
PRO TIPS:
- Preallocate slices/maps (make(..., cap)).
- Avoid interfaces in hot loops (if possible, to avoid dynamic dispatch).
- Use Profile Guided Optimization (PGO) in Go 1.20+:
  1. Collect cpuprofile from prod.
  2. go build -pgo=default.pgo
  2. go build -pgo=default.pgo
*/

/*
OPTIMIZATION PITFALLS:

1. Guessing vs Measuring:
   Error: "I think this is slow so I'll parallelize it."
   Reality: Benchmarking often shows channel overhead/mutex contention makes it slower than single-threaded.
   Rule: Benchmark first. 'pprof' second. Optimize third.

2. Large Structs by Pointer:
   Myth: "Pointers are always faster."
   Reality: Pointers stress the GC. Passing values (up to even 100 bytes) stays on stack (zero cost).

3. Sync.Pool Abuse:
   Risk: Storing database connections or heavy resources in Sync.Pool.
   Why: Sync.Pool items are drained arbitrarily during GC. You lose control of lifecycle.
*/
