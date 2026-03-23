// ============================================================================
// GO FUNDAMENTALS: MEMORY MODEL - STACK VS HEAP
// ============================================================================
// This file provides a comprehensive guide to Go's memory management,
// including stack vs heap allocation, escape analysis, garbage collection,
// and memory optimization techniques.
// ============================================================================

package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

// ============================================================================
// GLOBAL VARIABLES (Stored in data segment, not stack or heap)
// ============================================================================
var globalVar = "I live in the data segment"

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║           GO MEMORY MODEL: STACK VS HEAP                 ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Memory Regions Overview
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Memory Regions Overview")
	fmt.Println("─────────────────────────────────────────")

	// GO PROGRAM MEMORY LAYOUT:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │                        MEMORY LAYOUT                                   │
	// │                                                                        │
	// │  ┌─────────────────────────────────────────────────────────────────┐  │
	// │  │                        TEXT (Code)                              │  │
	// │  │              Machine instructions (read-only)                   │  │
	// │  ├─────────────────────────────────────────────────────────────────┤  │
	// │  │                     DATA/BSS Segment                            │  │
	// │  │         Global variables, constants, string literals            │  │
	// │  ├─────────────────────────────────────────────────────────────────┤  │
	// │  │                          HEAP                                   │  │
	// │  │        ↓ Grows downward (toward higher addresses)              │  │
	// │  │  - Dynamic allocation (new, make, &T{})                        │  │
	// │  │  - Managed by Garbage Collector                                │  │
	// │  │  - Slower allocation, slower access                            │  │
	// │  │  - Long-lived objects                                           │  │
	// │  │                          ...                                    │  │
	// │  │                          ...                                    │  │
	// │  │                          ...                                    │  │
	// │  │                         STACK                                   │  │
	// │  │        ↑ Grows upward (toward lower addresses)                 │  │
	// │  │  - Function call frames                                        │  │
	// │  │  - Local variables                                              │  │
	// │  │  - Fast allocation (just move stack pointer)                   │  │
	// │  │  - Automatic cleanup on function return                        │  │
	// │  │  - Goroutine-specific (each has own stack)                     │  │
	// │  └─────────────────────────────────────────────────────────────────┘  │
	// │                                                                        │
	// │  Note: In Go, the actual layout may differ; this is conceptual.       │
	// └────────────────────────────────────────────────────────────────────────┘

	fmt.Println("Memory regions:")
	fmt.Println("  1. TEXT:  Program code (instructions)")
	fmt.Println("  2. DATA:  Global variables, string literals")
	fmt.Println("  3. HEAP:  Dynamic allocations (GC managed)")
	fmt.Println("  4. STACK: Function call frames, local vars")

	// ========================================================================
	// SECTION 2: Stack Allocation
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Stack Allocation")
	fmt.Println("─────────────────────────────────────────")

	// STACK CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Pros:                                                                   │
	// │ ✓ Extremely fast allocation (just adjust stack pointer)                │
	// │ ✓ Automatic deallocation when function returns                         │
	// │ ✓ Great cache locality                                                 │
	// │ ✓ No GC overhead                                                       │
	// │                                                                          │
	// │ Cons:                                                                   │
	// │ ✗ Limited size (starts at 2KB per goroutine, can grow)                 │
	// │ ✗ Cannot outlive function scope (unless escaped to heap)               │
	// │ ✗ Not suitable for large allocations                                   │
	// └──────────────────────────────────────────────────────────────────────────┘

	// This stays on stack - returned by value
	val := staysOnStack()
	fmt.Printf("Value from stack: %d\n", val)

	// STACK FRAME VISUALIZATION:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Stack during function calls:                                          │
	// │                                                                        │
	// │ main() calls funcA() calls funcB():                                   │
	// │                                                                        │
	// │ LOW ADDRESS                                                            │
	// │ ┌──────────────────────────┐                                          │
	// │ │ funcB's stack frame      │ ◄── Stack pointer (SP)                   │
	// │ │ - local variables        │                                          │
	// │ │ - return address         │                                          │
	// │ ├──────────────────────────┤                                          │
	// │ │ funcA's stack frame      │                                          │
	// │ │ - local variables        │                                          │
	// │ │ - return address         │                                          │
	// │ ├──────────────────────────┤                                          │
	// │ │ main's stack frame       │                                          │
	// │ │ - local variables        │                                          │
	// │ └──────────────────────────┘                                          │
	// │ HIGH ADDRESS                                                           │
	// │                                                                        │
	// │ When funcB returns, SP moves down. funcB's frame is "freed" instantly.│
	// └────────────────────────────────────────────────────────────────────────┘

	// Demonstrate stack frame behavior
	demonstrateStackFrames()

	// ========================================================================
	// SECTION 3: Heap Allocation
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Heap Allocation")
	fmt.Println("─────────────────────────────────────────")

	// HEAP CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Pros:                                                                   │
	// │ ✓ Large capacity                                                       │
	// │ ✓ Objects can outlive function scope                                   │
	// │ ✓ Suitable for shared/long-lived data                                  │
	// │                                                                          │
	// │ Cons:                                                                   │
	// │ ✗ Slower allocation (need to find free space)                          │
	// │ ✗ Garbage Collection overhead                                          │
	// │ ✗ Poorer cache locality                                                │
	// │ ✗ Potential fragmentation                                              │
	// └──────────────────────────────────────────────────────────────────────────┘

	// This escapes to heap - returned by pointer
	ptr := escapesToHeap()
	fmt.Printf("Value from heap: %d (at address %p)\n", *ptr, ptr)

	// Allocations that escape to heap:
	// 1. Returning pointer to local variable
	// 2. Storing pointer in global variable
	// 3. Sending pointer to channel
	// 4. Closure capturing pointer
	// 5. Interface conversion (sometimes)
	// 6. Slice/map with unpredictable size

	// ========================================================================
	// SECTION 4: Escape Analysis
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Escape Analysis")
	fmt.Println("─────────────────────────────────────────")

	// Escape analysis is performed at COMPILE TIME.
	// The compiler decides whether variables go on stack or heap.
	//
	// Run: go build -gcflags="-m" main.go
	// This shows escape analysis decisions.
	//
	// Output examples:
	// ./main.go:XX: moved to heap: x
	// ./main.go:XX: x does not escape

	fmt.Println("To see escape analysis:")
	fmt.Println("  go build -gcflags=\"-m\" main.go")
	fmt.Println("  go build -gcflags=\"-m -m\" main.go  (more verbose)")
	fmt.Println()

	// Examples of escape scenarios
	fmt.Println("ESCAPE SCENARIOS:")
	fmt.Println("─────────────────────────────────────────")

	// Case 1: Value returned - STAYS ON STACK
	_ = noEscape()
	fmt.Println("1. Return by value: STACK")

	// Case 2: Pointer returned - ESCAPES TO HEAP
	_ = pointerEscapes()
	fmt.Println("2. Return pointer: HEAP (escapes)")

	// Case 3: Pointer passed down - may stay on stack
	y := 42
	usePointer(&y)
	fmt.Println("3. Pointer passed to callee: depends on callee")

	// Case 4: Interface conversion
	_ = toInterface(42)
	fmt.Println("4. Interface conversion: often HEAP")

	// Case 5: Closure capturing variable
	_ = closureCapture()
	fmt.Println("5. Closure capture: often HEAP")

	// Case 6: Large allocation
	_ = make([]byte, 1024*1024) // 1MB
	fmt.Println("6. Large slice: HEAP")

	// ========================================================================
	// SECTION 5: Goroutine Stacks
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Goroutine Stacks")
	fmt.Println("─────────────────────────────────────────")

	// GOROUTINE STACK FACTS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Each goroutine has its own stack                                     │
	// │ 2. Initial size: 2KB (much smaller than OS thread stack ~1-8MB)        │
	// │ 3. Stacks GROW dynamically as needed                                   │
	// │ 4. Maximum size: 1GB on 64-bit, 250MB on 32-bit                        │
	// │ 5. Growth is done by COPYING to a larger stack (contiguous)            │
	// │ 6. Shrinking happens during GC                                         │
	// └──────────────────────────────────────────────────────────────────────────┘

	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

	// Creating many goroutines - each has tiny 2KB stack
	done := make(chan bool, 1000)
	for i := 0; i < 1000; i++ {
		go func() {
			time.Sleep(10 * time.Millisecond)
			done <- true
		}()
	}

	// Wait for some to start
	time.Sleep(5 * time.Millisecond)
	fmt.Printf("After spawning 1000 goroutines: %d active\n", runtime.NumGoroutine())

	// Wait for completion
	for i := 0; i < 1000; i++ {
		<-done
	}

	// Stack growth demonstration (recursive function)
	fmt.Println("\nStack growth with deep recursion:")
	depth := deepRecursion(0, 10000)
	fmt.Printf("Reached depth: %d (stack grew as needed)\n", depth)

	// ========================================================================
	// SECTION 6: Garbage Collection
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Garbage Collection")
	fmt.Println("─────────────────────────────────────────")

	// GO'S GC CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Type: Concurrent, Tri-color Mark and Sweep                              │
	// │                                                                          │
	// │ Phases:                                                                 │
	// │ 1. Mark Setup (STW): Very short stop-the-world to enable write barrier │
	// │ 2. Marking (Concurrent): Traverse heap, mark reachable objects         │
	// │ 3. Mark Termination (STW): Complete marking, disable write barrier     │
	// │ 4. Sweeping (Concurrent): Reclaim unmarked objects                     │
	// │                                                                          │
	// │ STW = Stop The World (all goroutines paused)                           │
	// │ Go 1.8+: STW pauses typically < 1ms                                    │
	// └──────────────────────────────────────────────────────────────────────────┘

	// GC can be triggered by:
	// 1. Heap size reaches 2x previous post-GC size (GOGC=100 default)
	// 2. Manual trigger with runtime.GC()
	// 3. Every 2 minutes (if no other trigger)

	fmt.Println("GC Statistics:")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("  Heap Alloc:    %d KB\n", m.HeapAlloc/1024)
	fmt.Printf("  Heap Sys:      %d KB\n", m.HeapSys/1024)
	fmt.Printf("  Heap Objects:  %d\n", m.HeapObjects)
	fmt.Printf("  GC Cycles:     %d\n", m.NumGC)
	fmt.Printf("  Total Alloc:   %d KB\n", m.TotalAlloc/1024)

	// Allocate some memory
	allocateMemory()

	runtime.ReadMemStats(&m)
	fmt.Println("\nAfter allocations:")
	fmt.Printf("  Heap Alloc:    %d KB\n", m.HeapAlloc/1024)
	fmt.Printf("  Heap Objects:  %d\n", m.HeapObjects)

	// Force GC
	runtime.GC()

	runtime.ReadMemStats(&m)
	fmt.Println("\nAfter GC:")
	fmt.Printf("  Heap Alloc:    %d KB\n", m.HeapAlloc/1024)
	fmt.Printf("  GC Cycles:     %d\n", m.NumGC)

	// ========================================================================
	// SECTION 7: Memory Optimization Tips
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Memory Optimization Tips")
	fmt.Println("─────────────────────────────────────────")

	// OPTIMIZATION STRATEGIES:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Reduce allocations                                                   │
	// │    - Reuse buffers with sync.Pool                                      │
	// │    - Pre-allocate slices: make([]T, 0, expectedCap)                    │
	// │    - Use value types where appropriate                                 │
	// │                                                                          │
	// │ 2. Keep objects on stack                                               │
	// │    - Return values, not pointers (for small types)                     │
	// │    - Pass large structs by pointer (avoid copy AND allocation)         │
	// │    - Use arrays instead of slices for fixed sizes                      │
	// │                                                                          │
	// │ 3. Reduce GC pressure                                                  │
	// │    - Pool frequently allocated objects                                 │
	// │    - Avoid creating many short-lived objects in hot loops              │
	// │    - Use GOGC to tune GC frequency                                     │
	// │                                                                          │
	// │ 4. Struct field ordering                                               │
	// │    - Order fields from largest to smallest to minimize padding         │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Struct padding demonstration
	fmt.Println("\nStruct padding example:")
	demonstratePadding()

	// ========================================================================
	// SECTION 8: Profiling Memory
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Profiling Memory")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println("Memory profiling commands:")
	fmt.Println("  go test -memprofile=mem.prof -bench=.")
	fmt.Println("  go tool pprof mem.prof")
	fmt.Println()
	fmt.Println("Runtime debugging:")
	fmt.Println("  GODEBUG=gctrace=1 ./program")
	fmt.Println("  GODEBUG=allocfreetrace=1 ./program")
	fmt.Println()
	fmt.Println("Escape analysis:")
	fmt.Println("  go build -gcflags=\"-m -m\" ./...")

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Memory Model Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// Supporting Functions
// ============================================================================

// This stays on stack - returned by value
func staysOnStack() int {
	x := 42
	return x // Value is copied to caller's stack frame
}

// This escapes to heap - returned by pointer
func escapesToHeap() *int {
	x := 42
	return &x // x must be moved to heap so it survives function return
}

func demonstrateStackFrames() {
	fmt.Println("Calling nested functions:")
	level1()
}

func level1() {
	x := 1
	fmt.Printf("  level1: x = %d (on level1's stack frame)\n", x)
	level2()
}

func level2() {
	y := 2
	fmt.Printf("  level2: y = %d (on level2's stack frame)\n", y)
	level3()
}

func level3() {
	z := 3
	fmt.Printf("  level3: z = %d (on level3's stack frame)\n", z)
	// When level3 returns, its stack frame is instantly reclaimed
}

// Escape analysis examples
func noEscape() int {
	x := 42
	return x // x stays on stack
}

func pointerEscapes() *int {
	x := 42
	return &x // x escapes to heap
}

func usePointer(p *int) {
	*p++ // Modifies caller's variable
}

func toInterface(x int) interface{} {
	return x // May escape depending on usage
}

func closureCapture() func() int {
	x := 42
	return func() int {
		return x // x is captured, may escape
	}
}

// Deep recursion to demonstrate stack growth
func deepRecursion(current, max int) int {
	if current >= max {
		return current
	}
	return deepRecursion(current+1, max)
}

// Allocate memory for GC demonstration
func allocateMemory() {
	// Create many allocations
	allocations := make([]*[1024]byte, 1000)
	for i := range allocations {
		allocations[i] = new([1024]byte)
	}
	// allocations goes out of scope, eligible for GC
}

// Struct padding demonstration
type BadLayout struct {
	a bool  // 1 byte + 7 padding
	b int64 // 8 bytes
	c bool  // 1 byte + 7 padding
	d int64 // 8 bytes
	// Total: 32 bytes
}

type GoodLayout struct {
	b int64 // 8 bytes
	d int64 // 8 bytes
	a bool  // 1 byte
	c bool  // 1 byte + 6 padding
	// Total: 24 bytes
}

func demonstratePadding() {
	var bad BadLayout
	var good GoodLayout

	fmt.Printf("  BadLayout size:  %d bytes\n", unsafe.Sizeof(bad))
	fmt.Printf("  GoodLayout size: %d bytes\n", unsafe.Sizeof(good))
	fmt.Println("  (Order fields from largest to smallest)")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. PREMATURE OPTIMIZATION
   ─────────────────────────────────────────────────────────────────────────
   Don't obsess over stack vs heap without profiling first.
   The Go compiler is smart about escape analysis.
   Profile, then optimize.

2. ASSUMING SMALL = STACK
   ─────────────────────────────────────────────────────────────────────────
   Even small variables can escape to heap:
   - If pointer is returned
   - If stored in interface
   - If captured by closure

3. LARGE STACK ALLOCATIONS
   ─────────────────────────────────────────────────────────────────────────
   var arr [10000000]int  // May cause stack overflow

   Use slice for large arrays:
   arr := make([]int, 10000000)  // Heap allocated

4. MEMORY LEAKS FROM SLICES
   ─────────────────────────────────────────────────────────────────────────
   original := make([]int, 1000000)
   subset := original[:10]  // Keeps whole backing array alive!

   Fix: Copy to new slice if keeping small part of large slice.
   subset := make([]int, 10)
   copy(subset, original[:10])

5. GOROUTINE LEAKS
   ─────────────────────────────────────────────────────────────────────────
   go func() {
       <-ch  // Blocks forever if ch is never written/closed
   }()

   This goroutine's stack and memory are leaked.

6. UNBOUNDED ALLOCATIONS
   ─────────────────────────────────────────────────────────────────────────
   for {
       data := make([]byte, userSize)  // User controls size!
   }

   Validate sizes, set limits.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. PROFILE BEFORE OPTIMIZING
   Use pprof to identify actual bottlenecks:
   go test -bench=. -memprofile=mem.out

2. PRE-ALLOCATE SLICES WHEN SIZE IS KNOWN
   slice := make([]T, 0, knownCapacity)

3. USE sync.Pool FOR FREQUENTLY ALLOCATED OBJECTS
   var bufPool = sync.Pool{
       New: func() interface{} { return new(bytes.Buffer) },
   }

4. RETURN VALUES, NOT POINTERS FOR SMALL TYPES
   Keeps allocations on stack.

5. ORDER STRUCT FIELDS BY SIZE (DESCENDING)
   Minimizes padding, reduces memory footprint.

6. USE ARRAYS FOR FIXED-SIZE COLLECTIONS
   [N]T stays on stack; []T header on stack, backing array may escape.

7. AVOID ALLOCATIONS IN HOT LOOPS
   Move allocations outside the loop, reuse buffers.

8. SET REASONABLE GOGC
   GOGC=100 (default): GC when heap doubles
   GOGC=200: Less frequent GC, more memory
   GOGC=50: More frequent GC, less memory

9. CLEAR REFERENCES TO ENABLE GC
   Set pointers to nil when objects are no longer needed.

10. COPY SMALL SUBSETS OF LARGE SLICES
    Prevents keeping the entire backing array alive.
*/
