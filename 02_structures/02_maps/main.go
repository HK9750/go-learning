package main

import (
	"fmt"
	"sync"
)

// DEEP DIVE: Maps
// Maps are Hash Tables.
// O(1) average lookup, O(1) insert, O(1) delete.
// Internally implemented as "Buckets". Each bucket holds 8 key/value pairs.
// When buckets overflow, chaining or growing happens.
//
// VISUALIZATION (B=2 -> 4 Buckets):
// [ hmap ] -> [ buckets array ]
//             |
//             +-> [ Bucket 0 ] -> [ Tophash | Key1 | Key2.. | Val1 | Val2.. ] -> [ Overflow Bucket ]
//             +-> [ Bucket 1 ]
//             +-> [ Bucket 2 ]
//             +-> [ Bucket 3 ]

func main() {
	// 1. Initialization
	var nilMap map[string]int // nil. Reading OK (0), Writing PANIC!
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map
	// nilMap["key"] = 1 // PANIC: assignment to entry in nil map
	fmt.Println("Map is nil. Reading returns zero-value:", nilMap["foo"])

	m := make(map[string]int) // Initialized.
	m["Alice"] = 25
	m["Bob"] = 30

	// 2. The "Comma Ok" Idiom
	// Accessing a missing key returns the Zero Value (0 for int).
	// How to distinguish between "missing" and "value is actually 0"?
	m["ZeroVal"] = 0
	
	val, exists := m["NonExistent"]
	fmt.Printf("NonExistent: %d, Exists: %v\n", val, exists)

	val2, exists2 := m["ZeroVal"]
	fmt.Printf("ZeroVal: %d, Exists: %v\n", val2, exists2)

	// 3. Deletion is Safe
	delete(m, "Bob")
	delete(m, "NonExistent") // No-op, safe.

	// 4. Map iteration order is RANDOMIPED
	// To force randomness and prevent developers from relying on order,
	// Go randomizes the starting bucket during iteration.
	fmt.Println("\n--- Random Iteration Order ---")
	m2 := map[int]int{1:1, 2:2, 3:3, 4:4, 5:5}
	for k, v := range m2 {
		fmt.Printf("%d:%d ", k, v)
	}
	fmt.Println()

	// 5. Maps are NOT Thread-Safe
	// Concurrent Read/Write triggers fatal error: "concurrent map and/or map write"
	// Use sync.RWMutex or sync.Map.
	
	fmt.Println("\n--- Concurrent Map (Safe vs Unsafe) ---")
	// The below would crash if we ran it without mutex:
	var mu sync.Mutex
	safeMap := make(map[int]int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			mu.Lock()   // Lock before write
			safeMap[v] = v
			mu.Unlock() // Unlock after
		}(i)
	}
	wg.Wait()
	fmt.Println("SafeMap len:", len(safeMap))

	// 6. Memory Hint
	// Maps grow but never shrink (in terms of allocated buckets).
	// If you fill a map with 1M items and delete them all, it still consumes high RAM.
	// Solution: Re-make the map if you need to reclaim memory.
}
