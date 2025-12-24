package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// DEEP DIVE: The GMP Model (Goroutines, Machine Threads, Processors)
// G: Goroutine. Lightweight (~2KB stack).
// M: Machine (OS Thread). Expensive.
// P: Processor. Context for executing Go code. Default GOMAXPROCS = NumCPU.
//
// The Go Scheduler maps G -> M -> CPU via P.
// When a G blocks (e.g. syscall), the P detaches from M and moves to another M.
// This is "User Space Scheduling" (M:N scheduling).
//
// VISUALIZATION (GMP):
// [ Global Queue ]  <- Waiting Goroutines
//
//      P1 (Local Q)      P2 (Local Q)
//      |   [ G1 ]        |   [ G4 ]
//      v   [ G2 ]        v   [ G5 ]
//     [M1] (Thread)     [M2] (Thread)
//      |                 |
//     (CPU Core 1)      (CPU Core 2)

func main() {
	// 1. GOMAXPROCS
	// Controls how many OS threads can execute Go code simultaneously.
	fmt.Println("CPUs:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0)) // 0 means just query

	// 2. Goroutine Basics
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done() 
		fmt.Println("Goroutine 1 executing")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 2 executing")
	}()

	wg.Wait() // Wait for both to finish. Otherwise main exits and kills them!

	// 3. Concurrency vs Parallelism
	// Concurrency: Dealing with many things at once (Structure).
	// Parallelism: Doing many things at once (Execution).
	// One core = Concurrent but NOT Parallel.
	
	// 4. Gosched (Yielding)
	// Manually yielding core to other goroutines.
	// Only needed in tight CPU loops usually.
	time.Sleep(100 * time.Millisecond)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Working...")
			runtime.Gosched() // "I'm not done, but let others run"
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. WaitGroup Deadlock:
   Error: "fatal error: all goroutines are asleep - deadlock!"
   Why: If wg.Add(N) is > wg.Done() calls. The Wait() never unblocks.
   Fix: Ensure Done() is called in 'defer' to guarantee execution even during panic.

2. Goroutine Leak:
   Pitfall: Starting a goroutine that blocks forever (e.g., waiting on nil channel).
   Consequence: Memory leak. Goroutines (2KB) pile up.
   Fix: Ensure every goroutine has an exit condition (via Context or Channel close).
*/
