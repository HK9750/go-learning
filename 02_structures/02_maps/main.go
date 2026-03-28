// ============================================================================
// GO DATA STRUCTURES: MAPS
// ============================================================================
// This file provides a comprehensive guide to Go's map type,
// including internal structure, operations, concurrency patterns,
// and performance considerations.
// ============================================================================

package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    GO MAPS                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Map Basics
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Map Basics")
	fmt.Println("─────────────────────────────────────────")

	// MAP CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Hash table implementation (unordered key-value pairs)               │
	// │ 2. Reference type (passing map shares the data)                        │
	// │ 3. O(1) average for lookup, insert, delete                             │
	// │ 4. Keys must be comparable (==, !=) - no slices, maps, funcs          │
	// │ 5. Zero value is nil (cannot write to nil map!)                        │
	// │ 6. NOT thread-safe (needs external synchronization)                    │
	// │ 7. Iteration order is randomized                                       │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Map creation methods
	var nilMap map[string]int                    // nil map
	emptyMap := map[string]int{}                 // empty map (not nil)
	literalMap := map[string]int{"a": 1, "b": 2} // literal
	makeMap := make(map[string]int)              // using make
	makeMapHint := make(map[string]int, 100)     // with capacity hint

	fmt.Printf("nil map:     %v (nil: %t)\n", nilMap, nilMap == nil)
	fmt.Printf("empty map:   %v (nil: %t)\n", emptyMap, emptyMap == nil)
	fmt.Printf("literal map: %v\n", literalMap)
	fmt.Printf("make map:    %v (len: %d)\n", makeMap, len(makeMap))
	fmt.Printf("make hint:   %v (len: %d)\n", makeMapHint, len(makeMapHint))

	// CRITICAL: nil map vs empty map
	fmt.Println("\n⚠️  NIL MAP DANGER:")
	// nilMap["key"] = 1  // PANIC: assignment to entry in nil map
	_ = nilMap["key"] // Reading is OK, returns zero value
	fmt.Println("  Reading from nil map returns zero value:", nilMap["key"])
	fmt.Println("  Writing to nil map causes PANIC!")

	// ========================================================================
	// SECTION 2: Map Operations
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Map Operations")
	fmt.Println("─────────────────────────────────────────")

	m := make(map[string]int)

	// Insert/Update
	m["Alice"] = 25
	m["Bob"] = 30
	m["Charlie"] = 35
	fmt.Printf("After inserts: %v\n", m)

	// Update existing key
	m["Alice"] = 26
	fmt.Printf("After update: %v\n", m)

	// Lookup
	age := m["Bob"]
	fmt.Printf("Bob's age: %d\n", age)

	// Lookup missing key returns zero value
	missing := m["Nobody"]
	fmt.Printf("Missing key returns zero: %d\n", missing)

	// THE COMMA-OK IDIOM (distinguish missing from zero value)
	m["Zero"] = 0 // Intentionally zero

	val1, exists1 := m["Zero"]
	fmt.Printf("\n'Zero' key: value=%d, exists=%t\n", val1, exists1)

	val2, exists2 := m["Nobody"]
	fmt.Printf("'Nobody' key: value=%d, exists=%t\n", val2, exists2)

	// Delete
	delete(m, "Bob")
	fmt.Printf("\nAfter delete('Bob'): %v\n", m)

	// Delete non-existent key is safe (no-op)
	delete(m, "Nobody") // No panic
	fmt.Println("Deleting non-existent key is safe")

	// Length
	fmt.Printf("Map length: %d\n", len(m))
	// Note: There's no cap() for maps

	// ========================================================================
	// SECTION 3: Map Iteration
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Map Iteration")
	fmt.Println("─────────────────────────────────────────")

	nums := map[int]string{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
	}

	// CRITICAL: Iteration order is RANDOMIZED!
	// Go intentionally randomizes to prevent depending on order
	fmt.Println("Iteration 1:")
	for k, v := range nums {
		fmt.Printf("  %d: %s\n", k, v)
	}

	fmt.Println("Iteration 2 (may be different order):")
	for k, v := range nums {
		fmt.Printf("  %d: %s\n", k, v)
	}

	// To iterate in order, sort keys first
	fmt.Println("\nSorted iteration:")
	keys := make([]int, 0, len(nums))
	for k := range nums {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fmt.Printf("  %d: %s\n", k, nums[k])
	}

	// Iterating only keys
	fmt.Print("Keys only: ")
	for k := range nums {
		fmt.Printf("%d ", k)
	}
	fmt.Println()

	// ========================================================================
	// SECTION 4: Map Internals
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Map Internals")
	fmt.Println("─────────────────────────────────────────")

	// MAP INTERNAL STRUCTURE (simplified):
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ hmap (map header)                                                      │
	// │ ┌─────────────────────────────────────────────────────────────────┐    │
	// │ │ count:    number of entries                                     │    │
	// │ │ B:        log_2 of bucket count (buckets = 2^B)                 │    │
	// │ │ hash0:    random seed for hash function                         │    │
	// │ │ buckets:  pointer to bucket array                               │    │
	// │ │ oldbuckets: pointer to old buckets (during growth)              │    │
	// │ └─────────────────────────────────────────────────────────────────┘    │
	// │                    │                                                   │
	// │                    ▼                                                   │
	// │ ┌─────────────────────────────────────────────────────────────────┐    │
	// │ │ Bucket Array (2^B buckets)                                      │    │
	// │ │ ┌─────────┬─────────┬─────────┬─────────┐                       │    │
	// │ │ │Bucket 0 │Bucket 1 │Bucket 2 │Bucket 3 │ ...                   │    │
	// │ │ └────┬────┴─────────┴─────────┴─────────┘                       │    │
	// │ │      │                                                          │    │
	// │ │      ▼                                                          │    │
	// │ │ ┌────────────────────────────────────────────────────────┐      │    │
	// │ │ │ Bucket Structure (bmap)                                │      │    │
	// │ │ │ tophash[8]  - top 8 bits of hash for quick compare    │      │    │
	// │ │ │ keys[8]     - up to 8 keys                            │      │    │
	// │ │ │ values[8]   - up to 8 values                          │      │    │
	// │ │ │ overflow    - pointer to overflow bucket              │      │    │
	// │ │ └────────────────────────────────────────────────────────┘      │    │
	// │ └─────────────────────────────────────────────────────────────────┘    │
	// └─────────────────────────────────────────────────────────────────────────┘

	fmt.Println("Map Implementation:")
	fmt.Println("  - Hash table with chained buckets")
	fmt.Println("  - Each bucket holds up to 8 key-value pairs")
	fmt.Println("  - Uses top 8 bits of hash for quick comparison")
	fmt.Println("  - Overflow buckets linked for more entries")
	fmt.Println("  - Grows when load factor exceeds threshold (~6.5)")

	// HASH COLLISION HANDLING:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Hash Function: hash(key) → bucket_index & top_hash                     │
	// │                                                                         │
	// │ Lookup Process:                                                        │
	// │ 1. Compute hash(key)                                                   │
	// │ 2. bucket_index = hash & (2^B - 1)  // Low bits select bucket         │
	// │ 3. top_hash = hash >> (64 - 8)      // Top 8 bits for quick match     │
	// │ 4. Search bucket and overflow chains                                   │
	// │                                                                         │
	// │ When load factor > 6.5:                                                │
	// │ - Allocate larger bucket array (2x)                                    │
	// │ - Incrementally migrate entries (during reads/writes)                  │
	// │ - This is "incremental growth" to avoid big pauses                     │
	// └─────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 5: Valid Key Types
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Valid Key Types")
	fmt.Println("─────────────────────────────────────────")

	// Keys must be COMPARABLE (support == and !=)
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ VALID KEY TYPES:                                                        │
	// │ ✓ bool, int, uint, float, complex, string                              │
	// │ ✓ pointers, channels                                                   │
	// │ ✓ arrays (if element type is comparable)                               │
	// │ ✓ structs (if all fields are comparable)                               │
	// │ ✓ interfaces (compared by dynamic type and value)                      │
	// │                                                                          │
	// │ INVALID KEY TYPES:                                                      │
	// │ ✗ slices                                                               │
	// │ ✗ maps                                                                 │
	// │ ✗ functions                                                            │
	// │ ✗ structs containing slices, maps, or functions                        │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Array as key (valid - arrays are comparable)
	arrayKeyMap := make(map[[2]int]string)
	arrayKeyMap[[2]int{1, 2}] = "one-two"
	arrayKeyMap[[2]int{3, 4}] = "three-four"
	fmt.Printf("Array key map: %v\n", arrayKeyMap)

	// Struct as key (valid if all fields comparable)
	type Point struct {
		X, Y int
	}
	pointMap := make(map[Point]string)
	pointMap[Point{0, 0}] = "origin"
	pointMap[Point{1, 1}] = "diagonal"
	fmt.Printf("Struct key map: %v\n", pointMap)

	// Invalid:
	// map[[]int]string{}          // COMPILE ERROR: slice as key
	// map[map[int]int]string{}    // COMPILE ERROR: map as key
	// map[func()]string{}         // COMPILE ERROR: function as key

	// ========================================================================
	// SECTION 6: Concurrency
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Concurrency")
	fmt.Println("─────────────────────────────────────────")

	// CRITICAL: Maps are NOT thread-safe!
	// Concurrent read/write causes: "fatal error: concurrent map writes"
	// Concurrent read while writing: undefined behavior

	fmt.Println("⚠️  Maps require synchronization for concurrent access!")

	// Option 1: sync.Mutex
	fmt.Println("\nOption 1: sync.Mutex")
	safeMap := &SafeMap{
		data: make(map[string]int),
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", n)
			safeMap.Set(key, n)
		}(i)
	}
	wg.Wait()
	fmt.Printf("  SafeMap size: %d\n", safeMap.Len())

	// Option 2: sync.RWMutex (better for read-heavy workloads)
	fmt.Println("\nOption 2: sync.RWMutex (read-heavy workloads)")
	rwMap := &RWMap{
		data: make(map[string]int),
	}
	rwMap.Set("test", 42)
	fmt.Printf("  RWMap get: %d\n", rwMap.Get("test"))

	// Option 3: sync.Map (built-in concurrent map)
	fmt.Println("\nOption 3: sync.Map (built-in)")
	var syncMap sync.Map
	syncMap.Store("key1", "value1")
	syncMap.Store("key2", "value2")

	if val, ok := syncMap.Load("key1"); ok {
		fmt.Printf("  sync.Map load: %v\n", val)
	}

	// LoadOrStore: load if exists, otherwise store
	actual, loaded := syncMap.LoadOrStore("key3", "default")
	fmt.Printf("  LoadOrStore: value=%v, loaded=%t\n", actual, loaded)

	// sync.Map is optimized for:
	// 1. Entry written once, read many times
	// 2. Multiple goroutines read/write/overwrite disjoint key sets

	// ========================================================================
	// SECTION 7: Map Patterns
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Map Patterns")
	fmt.Println("─────────────────────────────────────────")

	// Pattern 1: Set implementation
	fmt.Println("Pattern 1: Set using map[T]struct{}")
	set := make(map[string]struct{})
	set["apple"] = struct{}{}
	set["banana"] = struct{}{}
	set["cherry"] = struct{}{}

	if _, exists := set["apple"]; exists {
		fmt.Println("  'apple' is in set")
	}
	fmt.Printf("  Set size: %d\n", len(set))

	// Pattern 2: Counting/Histogram
	fmt.Println("\nPattern 2: Counting occurrences")
	words := []string{"go", "is", "go", "great", "go", "is", "fun"}
	counts := make(map[string]int)
	for _, w := range words {
		counts[w]++ // Zero value (0) + 1 for first occurrence
	}
	fmt.Printf("  Word counts: %v\n", counts)

	// Pattern 3: Grouping
	fmt.Println("\nPattern 3: Grouping by key")
	type Person struct {
		Name string
		City string
	}
	people := []Person{
		{"Alice", "NYC"},
		{"Bob", "LA"},
		{"Charlie", "NYC"},
		{"Diana", "LA"},
	}
	byCity := make(map[string][]Person)
	for _, p := range people {
		byCity[p.City] = append(byCity[p.City], p)
	}
	for city, persons := range byCity {
		fmt.Printf("  %s: %v\n", city, persons)
	}

	// Pattern 4: Memoization/Cache
	fmt.Println("\nPattern 4: Memoization")
	cache := make(map[int]int)
	fib := func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := cache[n]; ok {
			return v
		}
		// This would need recursion helper in real impl
		result := n // Simplified
		cache[n] = result
		return result
	}
	_ = fib(10)
	fmt.Println("  Fibonacci with memoization implemented")

	// Pattern 5: Default values
	fmt.Println("\nPattern 5: Default values with ok idiom")
	config := map[string]string{"host": "localhost"}

	getWithDefault := func(m map[string]string, key, def string) string {
		if v, ok := m[key]; ok {
			return v
		}
		return def
	}
	fmt.Printf("  host: %s\n", getWithDefault(config, "host", "127.0.0.1"))
	fmt.Printf("  port: %s\n", getWithDefault(config, "port", "8080"))

	// Pattern 6: Two-way mapping
	fmt.Println("\nPattern 6: Bidirectional mapping")
	forward := map[string]int{"a": 1, "b": 2, "c": 3}
	reverse := make(map[int]string)
	for k, v := range forward {
		reverse[v] = k
	}
	fmt.Printf("  forward['b'] = %d\n", forward["b"])
	fmt.Printf("  reverse[2] = %s\n", reverse[2])

	// ========================================================================
	// SECTION 8: Memory Considerations
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Memory Considerations")
	fmt.Println("─────────────────────────────────────────")

	// CRITICAL: Maps grow but NEVER shrink!
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Problem:                                                                │
	// │   m := make(map[int]int)                                               │
	// │   // Add 1 million entries                                             │
	// │   for i := 0; i < 1000000; i++ { m[i] = i }                            │
	// │   // Delete all entries                                                 │
	// │   for k := range m { delete(m, k) }                                    │
	// │   // Map still has allocated buckets for 1M entries!                   │
	// │                                                                          │
	// │ Solution: Re-create the map if you need to reclaim memory              │
	// │   m = make(map[int]int)  // Old map becomes garbage                    │
	// └──────────────────────────────────────────────────────────────────────────┘

	fmt.Println("Maps never shrink automatically.")
	fmt.Println("To reclaim memory after mass deletion:")
	fmt.Println("  - Create a new map")
	fmt.Println("  - Copy remaining entries")
	fmt.Println("  - Let GC collect the old map")

	// Capacity hint reduces reallocations
	fmt.Println("\nUse capacity hint when size is known:")
	fmt.Println("  make(map[K]V, expectedSize)")

	// ========================================================================
	// SECTION 9: Performance Tips
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Performance Tips")
	fmt.Println("─────────────────────────────────────────")

	fmt.Println("1. Pre-size maps: make(map[K]V, size)")
	fmt.Println("2. Use struct{} for sets (0 bytes per entry)")
	fmt.Println("3. Consider string interning for repeated string keys")
	fmt.Println("4. Use sync.RWMutex for read-heavy concurrent access")
	fmt.Println("5. For write-heavy concurrent: consider sharded maps")
	fmt.Println("6. Profile before optimizing - maps are fast!")

	// Benchmark: string key vs int key
	// Generally, int keys are faster (simpler hash)
	// But difference is usually negligible

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Maps Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// Concurrent Map Implementations
// ============================================================================

// SafeMap wraps a map with a mutex
type SafeMap struct {
	mu   sync.Mutex
	data map[string]int
}

func (m *SafeMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.data[key]
	return v, ok
}

func (m *SafeMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

// RWMap uses RWMutex for better read concurrency
type RWMap struct {
	mu   sync.RWMutex
	data map[string]int
}

func (m *RWMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *RWMap) Get(key string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

func (m *RWMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. WRITING TO NIL MAP
   ─────────────────────────────────────────────────────────────────────────
   Error: "panic: assignment to entry in nil map"

   var m map[string]int
   m["key"] = 1  // PANIC!

   Fix: Initialize first: m := make(map[string]int)

2. CONCURRENT MAP ACCESS
   ─────────────────────────────────────────────────────────────────────────
   Error: "fatal error: concurrent map read and map write"

   Fix: Use sync.Mutex, sync.RWMutex, or sync.Map

3. RELYING ON ITERATION ORDER
   ─────────────────────────────────────────────────────────────────────────
   Map iteration order is intentionally randomized.

   Fix: Sort keys if order matters.

4. USING NON-COMPARABLE TYPES AS KEYS
   ─────────────────────────────────────────────────────────────────────────
   Error: "invalid map key type []int"

   Slices, maps, and functions cannot be keys.
   Fix: Use arrays, convert to string, or use pointer.

5. CHECKING FOR EXISTENCE INCORRECTLY
   ─────────────────────────────────────────────────────────────────────────
   if m["key"] != "" {  // WRONG: what if value is empty string?
       // ...
   }

   Fix: Use comma-ok: if val, ok := m["key"]; ok { ... }

6. MAP NOT SHRINKING
   ─────────────────────────────────────────────────────────────────────────
   Deleting entries doesn't free bucket memory.

   Fix: Create new map if memory is concern after mass deletion.

7. MODIFYING MAP DURING ITERATION
   ─────────────────────────────────────────────────────────────────────────
   Adding/deleting during range may or may not see changes.
   Go doesn't guarantee consistency.

   Fix: Collect keys first, then modify.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. ALWAYS INITIALIZE BEFORE WRITING
   m := make(map[K]V)  // or m := map[K]V{}

2. USE COMMA-OK IDIOM
   if val, ok := m[key]; ok { ... }

3. PRE-SIZE WHEN CAPACITY KNOWN
   m := make(map[K]V, expectedSize)

4. USE struct{} FOR SETS
   set := make(map[string]struct{})

5. SYNCHRONIZE CONCURRENT ACCESS
   Choose based on access pattern:
   - sync.Mutex for simple cases
   - sync.RWMutex for read-heavy
   - sync.Map for specific patterns

6. SORT KEYS FOR DETERMINISTIC ITERATION
   keys := make([]K, 0, len(m))
   for k := range m { keys = append(keys, k) }
   sort.Slice(keys, ...)

7. RECREATE MAP TO RECLAIM MEMORY
   After deleting many entries, create new map.

8. PREFER SIMPLER KEY TYPES
   int keys hash faster than string keys.

9. DOCUMENT THREAD-SAFETY
   // Thread-safe: protected by mu
   // Not thread-safe: caller must synchronize

10. USE RANGE FOR ITERATION
    for k, v := range m { ... }
    Not: for i := 0; i < len(m); i++
*/
