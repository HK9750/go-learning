package main

import (
	"fmt"
	"runtime"
	"sync"
)

// =========================================================================
// DEEP DIVE: GMP Scheduler & Context Switching
// =========================================================================
// G: Goroutine (Stack, PC)
// M: Machine (OS Thread)
// P: Processor (Local Run Queue of Gs)
//
// We will force concurrency on a SINGLE Processor (P=1) to see the scheduler work.

func main() {
	// 1. Force Single Processor
	// This ensures only ONE goroutine runs at a time physically.
	// Parallelism = 0, Concurrency = 1.
	prev := runtime.GOMAXPROCS(1)
	fmt.Printf("Previous MaxProcs: %d, Now: 1\n", prev)
	fmt.Println("=== Scheduler Trace (Single P) ===")

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine 1: "A"
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Printf("A%d ", i)
			// Yield the Processor!
			// This tells the scheduler: "Put me in the back of the generic queue, run someone else."
			if i == 2 {
				fmt.Print("[A YIELD] ")
				runtime.Gosched() 
			}
		}
	}()

	// Goroutine 2: "B"
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Printf("B%d ", i)
			// Yield the Processor
			if i == 2 {
				fmt.Print("[B YIELD] ")
				runtime.Gosched()
			}
		}
	}()

	wg.Wait()
	fmt.Println("\n=== Done ===")
	
	/*
	EXPECTED BEHAVIOR (GMP):
	Because we have P=1, they cannot run in parallel.
	
	VISUALIZATION (Time Slicing):
	Time --->
	[ P1 ] : [ A ] [ A ] [ A ] (Yield) -> [ B ] [ B ] [ B ] (Yield) -> [ A ] ...
	         |
	         Global/Local Queue logic manages this "Context Switch".
	*/
}
