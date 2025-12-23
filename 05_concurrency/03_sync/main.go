package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// DEEP DIVE: Sync Package
// Synchronization primitives for memory access.

func main() {
	// 1. Mutex (Mutual Exclusion)
	// Use pointer to avoid copying lock value!
	var mu sync.Mutex
	counter := 0

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("Mutex Counter:", counter)

	// 2. RWMutex (Read/Write Mutex)
	// Many readers allowed, or ONE writer.
	// Prefer this over Mutex if reads >>> writes (e.g., config cache).
	var rw sync.RWMutex
	data := make(map[string]string)

	rw.Lock() // Write lock
	data["key"] = "value"
	rw.Unlock()

	rw.RLock() // Read lock (shared)
	fmt.Println("Read:", data["key"])
	rw.RUnlock()

	// 3. Sync.Once (Singleton Pattern)
	// Ensures a function runs EXACTLY once, even if called concurrently.
	var once sync.Once
	var config string
	
	initFn := func() {
		fmt.Println("Initializing config (Once)...")
		config = "LOADED"
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(initFn) // Safe concurrent initialization
		}()
	}
	wg.Wait()
	fmt.Println("Config:", config)

	// 4. Sync.Pool (GC Optimization)
	// A pool of temporary objects. Reduces allocation pressure on GC.
	// NOT a cache (items can be cleared by GC anytime).
	pool := sync.Pool{
		New: func() interface{} {
			fmt.Println("Allocating new buffer...")
			return new(strings.Builder)
		},
	}

	// Get from pool
	sb := pool.Get().(*strings.Builder)
	sb.WriteString("Hello Pool")
	fmt.Println("Used:", sb.String())
	sb.Reset() // Clean before put back
	pool.Put(sb) // Return to pool

	// Reuse
	sb2 := pool.Get().(*strings.Builder) // Likely reuses existing object
	fmt.Println("Got buffer from pool. Capacity:", sb2.Cap())

	// 5. Atomics (Low-level non-blocking)
	var ops int64
	atomic.AddInt64(&ops, 1) // Safe increment without Mutex
	fmt.Println("Atomic Ops:", atomic.LoadInt64(&ops))
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Copying Locks (CRITICAL):
   Error: Passing Mutex by value copies the internal state.
   Code: 'func foo(mu sync.Mutex)'
   Result: The function gets a NEW lock. Locking it does nothing to the original.
   Fix: ALWAYS pass pointers: 'func foo(mu *sync.Mutex)'.
   Tool: 'go vet' detects this ("call of foo copies lock value").

2. Unlock without Lock:
   Error: "fatal error: sync: unlock of unlocked mutex"
   Why: Double unlock or unlocking wrong path.
   Fix: defer mu.Unlock() immediately after Lock().
*/
