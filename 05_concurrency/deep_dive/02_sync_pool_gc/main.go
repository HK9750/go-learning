package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// =========================================================================
// DEEP DIVE: Sync.Pool & GC
// =========================================================================
// sync.Pool caches items efficiently.
// KEY BEHAVIOR: Items in the pool are CLEARED during Garbage Collection.
// This prevents the Pool from becoming a memory leak.

func main() {
	fmt.Println("=== Sync.Pool GC Interaction ===")

	var counter int
	pool := sync.Pool{
		New: func() interface{} {
			counter++
			fmt.Printf("   [Pool] Allocating new object (Generated ID: %d)\n", counter)
			return counter
		},
	}

	// 1. Put an item
	pool.Put(100)
	fmt.Println("1. Put 100 into the pool.")

	// 2. Get it back (Should be 100, no allocation)
	val := pool.Get().(int)
	fmt.Printf("2. Got from pool: %d (Expected 100)\n", val)
	
	// 3. Put it back
	pool.Put(val)
	fmt.Println("3. Returned 100 to pool.")

	// 4. Force GC
	// The sync.Pool implementation registers a cleanup function with the runtime.
	// When GC runs, pools are often drained.
	fmt.Println("4. *** Triggering Garbage Collection ***")
	runtime.GC()
	// Debug note: Sometimes one GC isn't enough depending on runtime version/phases,
	// but generally for sync.Pool it connects to runtime.
	// Give it a moment.
	time.Sleep(10 * time.Millisecond) 

	// 5. Get again
	// If GC cleared the pool, New() will be called.
	val2 := pool.Get().(int)
	if val2 == 100 {
		fmt.Printf("5. Got from pool: %d (Survived GC - unexpected but possible)\n", val2)
	} else {
		// Expect 'counter' to increment -> 1 (since the first 100 was manual)
		// Wait, New() uses 'counter'. First New() call!
		// Logic: 
		// - We did manual Put(100).
		// - We did Get() -> 100.
		// - We did Put(100).
		// - GC happens. Pool clears 100.
		// - We did Get(). Pool empty. New() called.
		// - New() increments counter (0 -> 1) and returns 1.
		fmt.Printf("5. Got from pool: %d (New Allocation! Pool was drained)\n", val2)
	}
}
