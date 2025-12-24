package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// =========================================================================
// DEEP DIVE: Mutex Internals
// =========================================================================
// A sync.Mutex is just 2 fields!
// type Mutex struct {
//    state int32
//    sema  uint32
// }

// We mirror it here to inspect the state bits
type MutexState struct {
	State int32
	Sema  uint32
}

const (
	mutexLocked = 1 << iota // 1
	mutexWoken              // 2
	mutexStarving           // 4
)

func main() {
	var mu sync.Mutex
	
	// Direct access to internal layout
	// Note: We use unsafe.Pointer(&mu) to get *Mutex, then cast to *MutexState
	internal := (*MutexState)(unsafe.Pointer(&mu))

	fmt.Printf("1. Initial State: %d (Locked=%t, Woken=%t)\n", 
		internal.State, 
		internal.State&mutexLocked != 0,
		internal.State&mutexWoken != 0)

	// Step 1: Lock
	mu.Lock()
	fmt.Printf("2. Locked:        %d (Locked=%t)\n", internal.State, internal.State&mutexLocked != 0)

	// Step 2: Contention Loop
	// We want to see the 'Woken' bit or 'Waiter' count.
	// Waiters are stored in higher bits: state >> 3.
	
	go func() {
		// Try to lock (will block because main holds it)
		// This should increment the waiter count in 'state'
		fmt.Println("   [Goroutine] Trying to lock...")
		mu.Lock()
		fmt.Println("   [Goroutine] Acquired lock!")
		mu.Unlock()
	}()

	time.Sleep(100 * time.Millisecond) // Give goroutine time to park

	state := internal.State
	waiters := state >> 3
	fmt.Printf("3. Contention:    %d (Locked=%t, Woken=%t, Waiters=%d)\n", 
		state, 
		state&mutexLocked != 0,
		state&mutexWoken != 0,
		waiters)

	// Step 3: Unlock
	// When we unlock, if there are waiters, the runtime might set Woken bit on the next waiter wakeup
	fmt.Println("   [Main] Unlocking...")
	mu.Unlock()
	
	time.Sleep(100 * time.Millisecond) // Wait for Goroutine to finish
	
	fmt.Printf("4. Final State:   %d (Locked=%t)\n", internal.State, internal.State&mutexLocked != 0)
}
