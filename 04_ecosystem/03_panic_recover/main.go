package main

import "fmt"

// DEEP DIVE: Panic & Recover
// Panic is like throwing a RuntimeException. It stops normal flow.
// Recover is like catch, but ONLY works inside a DEFERRED function.
//
// VISUALIZATION (Stack Unwinding):
// FAIL: panic() called!
//   |
//   v
// [ Stack Frame N   ] -> Deferred functions run? -> NO (Destroyed)
//   |
//   v
// [ ...             ]
//   |
//   v
// [ Stack Frame 1   ] -> defer func() { recover() } -> YES! (Stop unwinding, resume here)
//   |
// [ main()          ]

func riskyOperation(i int) {
	defer func() {
		// RECOVER must be called inside defer
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic in riskyOperation(%d): %v\n", i, r)
			// You can choose to re-panic here if the error is fatal
			// panic(r)
		}
	}()

	if i > 5 {
		panic("Something went wrong! value too high")
	}
	fmt.Println("Operation success:", i)
}

func main() {
	// 1. Safe boundary
	fmt.Println("Start processing...")
	
	// This loop will NOT crash the program, because riskyOperation recovers itself.
	for i := 0; i <= 7; i++ {
		riskyOperation(i)
	}

	fmt.Println("End of program. We survived.")

	// 2. Real world usage
	// HTTP Servers: Each request runs in a goroutine.
	// You always wrap the handler in a recover() middleware so one bad request
	// doesn't crash the entire server.
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Cross-Goroutine Panics:
   Critical: Recover ONLY works in the SAME goroutine that panicked.
   Scenario: Main calls 'go func() { panic() }'. Main's defer/recover CANNOT catch this.
   Result: Program crash.
   Fix: Every goroutine needs its own recover if you want to be safe.

2. Recover outside Defer:
   Code: 'if r := recover(); ...' inside normal control flow.
   Result: Returns nil, does nothing. Must be in 'defer'.
*/
