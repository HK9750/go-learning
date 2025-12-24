package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// =========================================================================
// DEEP DIVE: Atomics vs Data Race
// =========================================================================
// Incrementing an integer (i++) is NOT atomic. It is 3 CPU instructions:
// 1. LOAD  reg, [addr]
// 2. ADD   reg, 1
// 3. STORE [addr], reg
//
// If two goroutines interleave these steps, updates are lost.
// Atomics use hardware instructions (LOCK prefix on x86) to make this indivisible.

func main() {
	var wg sync.WaitGroup
	var unsafeCounter int32
	var atomicCounter int32

	const goroutines = 50
	const increments = 1000

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				// 1. Unsafe Increment (Race Condition)
				// Race: Read(5) -> ContextSwitch -> OtherWrite(6) -> Write(6) (Overwrites 6!)
				unsafeCounter++ 

				// 2. Atomic Increment (Safe)
				atomic.AddInt32(&atomicCounter, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Println("=== Atomics vs Data Race Visualization ===")
	fmt.Printf("Expected Count : %d\n", goroutines*increments)
	fmt.Printf("Atomic Count   : %d (Values match!)\n", atomicCounter)
	fmt.Printf("Unsafe Count   : %d (Values lost due to races)\n", unsafeCounter)

	/*
	VISUALIZATION OF A RACE:
	
	Gorountine A    |   Goroutine B
	----------------|----------------
	LOAD X (Val: 5) |
	                |   LOAD X (Val: 5)
	ADD 1  (Reg: 6) |
	                |   ADD 1  (Reg: 6)
	STORE X (6)     |
	                |   STORE X (6)
	                
	Result: X is 6, but should be 7! One update is lost.
	*/
}
