// ============================================================================
// GO DATA STRUCTURES: ARRAYS & SLICES
// ============================================================================
// This file provides a comprehensive guide to Go's arrays and slices,
// including internal structures, memory layouts, performance considerations,
// and common patterns.
// ============================================================================

package main

import (
	"fmt"
	"reflect"
	"sort"
	"unsafe"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              GO ARRAYS & SLICES                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Arrays
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Arrays")
	fmt.Println("─────────────────────────────────────────")

	// ARRAY CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Fixed size (size is part of the type)                               │
	// │ 2. Value type (copying creates a complete copy)                        │
	// │ 3. Stored contiguously in memory                                       │
	// │ 4. Size must be constant expression (known at compile time)            │
	// │ 5. Zero value is array of zero values                                  │
	// │ 6. Comparable (can use == if element type is comparable)               │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Array declaration and initialization
	var arr1 [5]int                 // Zero-initialized: [0 0 0 0 0]
	arr2 := [5]int{1, 2, 3, 4, 5}   // Literal initialization
	arr3 := [5]int{1, 2, 3}         // Partial: [1 2 3 0 0]
	arr4 := [...]int{1, 2, 3, 4, 5} // Size inferred: [5]int
	arr5 := [5]int{0: 10, 4: 50}    // Indexed: [10 0 0 0 50]

	fmt.Printf("arr1 (zero):     %v\n", arr1)
	fmt.Printf("arr2 (literal):  %v\n", arr2)
	fmt.Printf("arr3 (partial):  %v\n", arr3)
	fmt.Printf("arr4 (inferred): %v (type: %T)\n", arr4, arr4)
	fmt.Printf("arr5 (indexed):  %v\n", arr5)

	// CRITICAL: Size is part of the type!
	// [3]int and [5]int are DIFFERENT types!
	// var x [3]int = arr2  // COMPILE ERROR: cannot use [5]int as [3]int

	// Array comparison
	a1 := [3]int{1, 2, 3}
	a2 := [3]int{1, 2, 3}
	a3 := [3]int{1, 2, 4}
	fmt.Printf("\n[1,2,3] == [1,2,3]: %t\n", a1 == a2)
	fmt.Printf("[1,2,3] == [1,2,4]: %t\n", a1 == a3)

	// CRITICAL: Arrays are VALUE types - copying duplicates all elements!
	fmt.Println("\nArray value semantics:")
	original := [3]int{1, 2, 3}
	copied := original
	copied[0] = 999
	fmt.Printf("Original: %v (unchanged)\n", original)
	fmt.Printf("Copied:   %v (modified)\n", copied)

	// MEMORY LAYOUT:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Array [5]int in memory:                                                │
	// │                                                                         │
	// │ ┌───────┬───────┬───────┬───────┬───────┐                              │
	// │ │ [0]   │ [1]   │ [2]   │ [3]   │ [4]   │                              │
	// │ │ 8byte │ 8byte │ 8byte │ 8byte │ 8byte │                              │
	// │ └───────┴───────┴───────┴───────┴───────┘                              │
	// │ Total: 40 bytes (5 × 8 bytes) - contiguous in memory                   │
	// │ Excellent cache locality for sequential access                         │
	// └─────────────────────────────────────────────────────────────────────────┘

	fmt.Printf("\nSize of [5]int: %d bytes\n", unsafe.Sizeof(arr2))

	// ========================================================================
	// SECTION 2: Slices - The Dynamic View
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Slices - The Dynamic View")
	fmt.Println("─────────────────────────────────────────")

	// SLICE CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Dynamically-sized view into an array                                │
	// │ 2. Reference type (header is value, but references shared array)       │
	// │ 3. Has length (len) and capacity (cap)                                 │
	// │ 4. Zero value is nil (but len=0, cap=0, usable with append)           │
	// │ 5. NOT comparable (can only compare to nil)                            │
	// │ 6. Most common collection type in Go                                   │
	// └──────────────────────────────────────────────────────────────────────────┘

	// SLICE INTERNAL STRUCTURE:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Slice Header (24 bytes on 64-bit):                                     │
	// │ ┌─────────────────────┐                                                │
	// │ │ Data *T (8 bytes)   │──────► Points to element [0] of backing array │
	// │ │ Len int (8 bytes)   │        Number of elements in slice            │
	// │ │ Cap int (8 bytes)   │        Max elements before reallocation       │
	// │ └─────────────────────┘                                                │
	// │                                                                         │
	// │ Example: s := []int{10, 20, 30} with cap 5                             │
	// │                                                                         │
	// │   Slice Header           Backing Array                                 │
	// │ ┌─────────────┐        ┌────┬────┬────┬────┬────┐                      │
	// │ │ Data: 0xABC │───────►│ 10 │ 20 │ 30 │ -- │ -- │                      │
	// │ │ Len:  3     │        └────┴────┴────┴────┴────┘                      │
	// │ │ Cap:  5     │        │← len=3 →│←remaining→│                         │
	// │ └─────────────┘        │←────── cap=5 ──────→│                         │
	// └─────────────────────────────────────────────────────────────────────────┘

	// Slice creation methods
	var slice1 []int             // nil slice
	slice2 := []int{}            // empty slice (not nil!)
	slice3 := []int{1, 2, 3}     // literal
	slice4 := make([]int, 5)     // make with length (zeroed)
	slice5 := make([]int, 3, 10) // make with length and capacity
	slice6 := arr2[1:4]          // slice from array

	fmt.Printf("nil slice:   %v (nil: %t, len: %d, cap: %d)\n",
		slice1, slice1 == nil, len(slice1), cap(slice1))
	fmt.Printf("empty slice: %v (nil: %t, len: %d, cap: %d)\n",
		slice2, slice2 == nil, len(slice2), cap(slice2))
	fmt.Printf("literal:     %v (len: %d, cap: %d)\n",
		slice3, len(slice3), cap(slice3))
	fmt.Printf("make(5):     %v (len: %d, cap: %d)\n",
		slice4, len(slice4), cap(slice4))
	fmt.Printf("make(3,10):  %v (len: %d, cap: %d)\n",
		slice5, len(slice5), cap(slice5))
	fmt.Printf("arr[1:4]:    %v (len: %d, cap: %d)\n",
		slice6, len(slice6), cap(slice6))

	// ========================================================================
	// SECTION 3: Slice Operations
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Slice Operations")
	fmt.Println("─────────────────────────────────────────")

	// SLICING SYNTAX: slice[low:high:max]
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ low  : Starting index (inclusive), default 0                           │
	// │ high : Ending index (exclusive), default len                           │
	// │ max  : Cap limit (exclusive), default cap - optional                   │
	// │                                                                          │
	// │ Result: len = high - low                                                │
	// │         cap = max - low (or original_cap - low if max not specified)   │
	// └──────────────────────────────────────────────────────────────────────────┘

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("Original: %v (len=%d, cap=%d)\n", s, len(s), cap(s))

	// Basic slicing
	fmt.Printf("s[2:5]:   %v (len=%d, cap=%d)\n", s[2:5], len(s[2:5]), cap(s[2:5]))
	fmt.Printf("s[:5]:    %v (len=%d, cap=%d)\n", s[:5], len(s[:5]), cap(s[:5]))
	fmt.Printf("s[5:]:    %v (len=%d, cap=%d)\n", s[5:], len(s[5:]), cap(s[5:]))
	fmt.Printf("s[:]:     %v (len=%d, cap=%d)\n", s[:], len(s[:]), cap(s[:]))

	// Three-index slicing (limiting capacity)
	limited := s[2:5:7] // len=3, cap=5
	fmt.Printf("s[2:5:7]: %v (len=%d, cap=%d) - capacity limited!\n",
		limited, len(limited), cap(limited))

	// VISUALIZATION:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Original s: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]                             │
	// │             ^        ^        ^           ^                            │
	// │             0        2        5           10                           │
	// │                                                                         │
	// │ s[2:5]:     [      2, 3, 4               ]                             │
	// │                    ^        ^                                          │
	// │                    0        3   len=3, cap=8 (can see rest!)           │
	// │                                                                         │
	// │ s[2:5:7]:   [      2, 3, 4      ]                                      │
	// │                    ^        ^                                          │
	// │                    0        3   len=3, cap=5 (capacity LIMITED)        │
	// └─────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 4: Append Deep Dive
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Append Deep Dive")
	fmt.Println("─────────────────────────────────────────")

	// APPEND BEHAVIOR:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ If len < cap:                                                           │
	// │   - Adds element, increments len                                       │
	// │   - Returns slice with same backing array                              │
	// │                                                                          │
	// │ If len == cap:                                                          │
	// │   - Allocates NEW larger array                                         │
	// │   - Copies existing elements                                           │
	// │   - Adds new element                                                   │
	// │   - Returns slice with NEW backing array                               │
	// │                                                                          │
	// │ GROWTH STRATEGY (approximate):                                          │
	// │   - Small slices (<1024): double capacity                              │
	// │   - Large slices: grow by ~25%                                         │
	// │   - Always ensure enough for new elements                              │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Observing append behavior
	fmt.Println("Watching append growth:")
	var growing []int
	prevCap := 0
	for i := 0; i < 20; i++ {
		growing = append(growing, i)
		if cap(growing) != prevCap {
			fmt.Printf("  len=%2d, cap=%2d (grew!)\n", len(growing), cap(growing))
			prevCap = cap(growing)
		}
	}

	// CRITICAL: Append may change the backing array!
	fmt.Println("\nAppend may disconnect slices:")
	base := make([]int, 3, 5)
	base[0], base[1], base[2] = 1, 2, 3

	derived := base[:]

	// This fits in capacity - same backing array
	base = append(base, 4)
	fmt.Printf("After append(base, 4):\n")
	fmt.Printf("  base:    %v (cap=%d)\n", base, cap(base))
	fmt.Printf("  derived: %v (still connected)\n", derived)

	// This exceeds capacity - NEW backing array
	base = append(base, 5, 6, 7, 8, 9)
	fmt.Printf("After append(base, 5,6,7,8,9):\n")
	fmt.Printf("  base:    %v (cap=%d)\n", base, cap(base))

	base[0] = 999
	fmt.Printf("  base[0]=999, derived: %v (disconnected!)\n", derived)

	// Appending slices
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	combined := append(s1, s2...) // Use ... to unpack slice
	fmt.Printf("\nCombining slices: %v + %v = %v\n", s1, s2, combined)

	// ========================================================================
	// SECTION 5: Copy Function
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Copy Function")
	fmt.Println("─────────────────────────────────────────")

	// copy(dst, src) copies elements from src to dst
	// Returns number of elements copied (min of len(dst), len(src))

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)

	n := copy(dst, src)
	fmt.Printf("copy(dst[3], src[5]): copied %d elements\n", n)
	fmt.Printf("  src: %v\n", src)
	fmt.Printf("  dst: %v\n", dst)

	// Copy into larger destination
	dst2 := make([]int, 10)
	copy(dst2, src)
	fmt.Printf("copy(dst[10], src[5]): %v\n", dst2)

	// Copy to create independent slice
	original2 := []int{1, 2, 3}
	independent := make([]int, len(original2))
	copy(independent, original2)
	independent[0] = 999
	fmt.Printf("\nIndependent copy: original=%v, copy=%v\n", original2, independent)

	// Self-copy (overlapping) - handled correctly
	overlap := []int{1, 2, 3, 4, 5}
	copy(overlap[2:], overlap[:3])
	fmt.Printf("Self-copy overlap: %v\n", overlap)

	// ========================================================================
	// SECTION 6: Slice Gotchas & Memory Leaks
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Slice Gotchas & Memory Leaks")
	fmt.Println("─────────────────────────────────────────")

	// GOTCHA 1: Slicing doesn't copy - shares backing array
	fmt.Println("Gotcha 1: Shared backing array")
	arr := []int{1, 2, 3, 4, 5}
	sub := arr[1:3]
	sub[0] = 999
	fmt.Printf("  Modifying sub affects original: %v\n", arr)

	// GOTCHA 2: Large slice keeping memory alive
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ MEMORY LEAK SCENARIO:                                                  │
	// │                                                                         │
	// │ large := make([]byte, 1000000)  // 1MB allocated                       │
	// │ small := large[:10]              // small keeps ALL of large alive!    │
	// │ large = nil                      // Doesn't help - small still refs it │
	// │                                                                         │
	// │ FIX: Copy what you need                                                │
	// │ small := make([]byte, 10)                                              │
	// │ copy(small, large[:10])          // Now large can be GC'd              │
	// └─────────────────────────────────────────────────────────────────────────┘

	fmt.Println("\nGotcha 2: Memory leak from slicing")
	demonstrateMemoryLeak()

	// GOTCHA 3: Append can affect other slices
	fmt.Println("\nGotcha 3: Append side effects")
	original3 := make([]int, 3, 6)
	original3[0], original3[1], original3[2] = 1, 2, 3

	view1 := original3[:]
	view2 := append(original3, 4) // Writes to original's backing array!

	fmt.Printf("  original after append: %v\n", original3)
	fmt.Printf("  view1: %v (might see changes!)\n", view1)
	fmt.Printf("  view2: %v\n", view2)

	// Fix: Use three-index slice to limit capacity
	safe := make([]int, 3, 6)
	safe[0], safe[1], safe[2] = 1, 2, 3
	limitedView := safe[0:3:3] // cap=3, so append will reallocate
	extended := append(limitedView, 4, 5)
	fmt.Printf("\n  Safe: original unchanged: %v\n", safe)
	fmt.Printf("  extended: %v\n", extended)

	// ========================================================================
	// SECTION 7: Slice Patterns
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Common Slice Patterns")
	fmt.Println("─────────────────────────────────────────")

	// Pattern 1: Filter in place
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	evens := filterInPlace(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter evens: %v\n", evens)

	// Pattern 2: Remove element at index
	items := []string{"a", "b", "c", "d", "e"}
	items = removeAt(items, 2) // Remove "c"
	fmt.Printf("Remove at 2: %v\n", items)

	// Pattern 3: Insert at index
	items = insertAt(items, 2, "X")
	fmt.Printf("Insert X at 2: %v\n", items)

	// Pattern 4: Reverse
	toReverse := []int{1, 2, 3, 4, 5}
	reverse(toReverse)
	fmt.Printf("Reversed: %v\n", toReverse)

	// Pattern 5: Deduplicate sorted slice
	withDups := []int{1, 1, 2, 2, 2, 3, 4, 4, 5}
	unique := dedup(withDups)
	fmt.Printf("Dedup: %v\n", unique)

	// Pattern 6: Stack using slice
	fmt.Println("\nStack operations:")
	var stack []int
	stack = push(stack, 1)
	stack = push(stack, 2)
	stack = push(stack, 3)
	fmt.Printf("  After pushes: %v\n", stack)
	stack, val := pop(stack)
	fmt.Printf("  Pop: %d, stack: %v\n", val, stack)

	// ========================================================================
	// SECTION 8: Slice Internals
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Slice Internals")
	fmt.Println("─────────────────────────────────────────")

	// Inspecting slice header
	inspectSlice := []int{10, 20, 30, 40, 50}
	header := (*reflect.SliceHeader)(unsafe.Pointer(&inspectSlice))
	fmt.Printf("Slice: %v\n", inspectSlice)
	fmt.Printf("Header: Data=0x%x, Len=%d, Cap=%d\n",
		header.Data, header.Len, header.Cap)
	fmt.Printf("Size of slice header: %d bytes\n", unsafe.Sizeof(inspectSlice))

	// Nil vs Empty slice
	fmt.Println("\nNil vs Empty slice:")
	var nilSlice []int
	emptySlice := []int{}
	makeSlice := make([]int, 0)

	fmt.Printf("nil slice:   %v, nil=%t, len=%d, cap=%d\n",
		nilSlice, nilSlice == nil, len(nilSlice), cap(nilSlice))
	fmt.Printf("empty slice: %v, nil=%t, len=%d, cap=%d\n",
		emptySlice, emptySlice == nil, len(emptySlice), cap(emptySlice))
	fmt.Printf("make slice:  %v, nil=%t, len=%d, cap=%d\n",
		makeSlice, makeSlice == nil, len(makeSlice), cap(makeSlice))

	// All work with append, len, range - no special handling needed
	nilSlice = append(nilSlice, 1, 2, 3)
	fmt.Printf("After append to nil: %v\n", nilSlice)

	// ========================================================================
	// SECTION 9: Performance Tips
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Performance Tips")
	fmt.Println("─────────────────────────────────────────")

	// TIP 1: Pre-allocate when size is known
	fmt.Println("Tip 1: Pre-allocate to avoid reallocations")
	// Bad: grows multiple times
	// Good: single allocation
	efficient := make([]int, 0, 1000)
	fmt.Printf("  Pre-allocated: len=%d, cap=%d\n", len(efficient), cap(efficient))

	// TIP 2: Use copy for large slices
	fmt.Println("\nTip 2: Use copy() for performance-critical code")

	// TIP 3: Consider array for fixed-size collections
	fmt.Println("\nTip 3: Use [N]T for fixed size (stays on stack)")

	// TIP 4: Clear slice by setting length to zero
	fmt.Println("\nTip 4: Clear slice efficiently")
	data := []int{1, 2, 3, 4, 5}
	data = data[:0] // Keeps capacity, sets length to 0
	fmt.Printf("  Cleared: %v (cap=%d)\n", data, cap(data))

	// TIP 5: Use sort.Slice for custom sorting
	fmt.Println("\nTip 5: sort.Slice for custom sorting")
	people := []struct {
		Name string
		Age  int
	}{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Printf("  Sorted by age: %v\n", people)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Arrays & Slices Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// Helper Functions
// ============================================================================

func demonstrateMemoryLeak() {
	// Simulated - in real code this would be problematic
	large := make([]byte, 1000000) // 1MB
	_ = large[0]                   // Use it

	// BAD: small keeps large alive
	// small := large[:10]

	// GOOD: Copy what you need
	small := make([]byte, 10)
	copy(small, large[:10])

	fmt.Printf("  Created independent 10-byte slice from 1MB slice\n")
	_ = small
}

func filterInPlace(s []int, keep func(int) bool) []int {
	n := 0
	for _, v := range s {
		if keep(v) {
			s[n] = v
			n++
		}
	}
	return s[:n]
}

func removeAt[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

func insertAt[T any](s []T, i int, v T) []T {
	s = append(s, *new(T)) // Grow by one
	copy(s[i+1:], s[i:])   // Shift right
	s[i] = v               // Insert
	return s
}

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func dedup[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}
	n := 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			s[n] = s[i]
			n++
		}
	}
	return s[:n]
}

func push[T any](s []T, v T) []T {
	return append(s, v)
}

func pop[T any](s []T) ([]T, T) {
	return s[:len(s)-1], s[len(s)-1]
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. SLICE BOUNDS OUT OF RANGE
   ─────────────────────────────────────────────────────────────────────────
   Error: "panic: runtime error: index out of range"

   s := []int{1, 2, 3}
   _ = s[5]  // PANIC!

   Fix: Check bounds before accessing.

2. MODIFYING SLICE DURING ITERATION
   ─────────────────────────────────────────────────────────────────────────
   for i, v := range s {
       if shouldRemove(v) {
           s = append(s[:i], s[i+1:]...)  // BUG: indices shift!
       }
   }

   Fix: Iterate backwards or collect indices first.

3. APPEND NOT REASSIGNED
   ─────────────────────────────────────────────────────────────────────────
   s := []int{1, 2, 3}
   append(s, 4)  // WARNING: result ignored!

   Fix: s = append(s, 4)

4. COMPARING SLICES WITH ==
   ─────────────────────────────────────────────────────────────────────────
   Error: "invalid operation: s1 == s2 (slice can only be compared to nil)"

   Fix: Use slices.Equal() or manual comparison.

5. FORGETTING SLICE IS REFERENCE-LIKE
   ─────────────────────────────────────────────────────────────────────────
   func modify(s []int) {
       s[0] = 999  // Modifies original!
   }

   If unintended: pass copy.

6. NIL SLICE PANIC ASSUMPTION
   ─────────────────────────────────────────────────────────────────────────
   var s []int  // nil
   len(s)       // OK, returns 0
   cap(s)       // OK, returns 0
   s = append(s, 1)  // OK, creates new backing array
   s[0]         // OK after append

   nil slices are safe for read operations!
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. PRE-ALLOCATE WHEN SIZE IS KNOWN
   s := make([]T, 0, expectedSize)

2. USE copy() FOR INDEPENDENT SLICES
   Prevents unintended sharing of backing array.

3. USE THREE-INDEX SLICING FOR SAFETY
   s[low:high:max] limits capacity.

4. PREFER nil OVER EMPTY FOR ZERO-VALUE
   var s []int  // nil is idiomatic for "no slice"
   s := []int{} // only when empty slice specifically needed (JSON)

5. PASS LARGE SLICES AS IS (NOT BY POINTER)
   Slice header is only 24 bytes.

6. CLEAR SLICE BY RESLICING
   s = s[:0]  // Keeps capacity

7. BE CAREFUL WITH APPEND IN SHARED SLICES
   Document when slices share backing array.

8. COPY SMALL PORTIONS OF LARGE SLICES
   Prevents memory leaks from keeping large backing array alive.

9. USE sort.Slice FOR CUSTOM SORTING
   More flexible than implementing sort.Interface.

10. CONSIDER GENERICS FOR REUSABLE SLICE OPERATIONS
    Go 1.18+ allows generic slice utilities.
*/
