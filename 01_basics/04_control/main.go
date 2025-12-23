package main

import (
	"fmt"
	"time"
)

// DEEP DIVE: Control Structures
// Go has fewer control structures than languages like C++ or Java.
// No while, no do-while. Only 'for'.
// All loops are created using the 'for' keyword.

func main() {
	// 1. The Many Faces of 'For'
	// A. Standard C-style
	for i := 0; i < 3; i++ {
		fmt.Print(i, " ")
	}
	
	// B. While-style
	j := 0
	for j < 3 {
		j++ // Go doesn't have pre-increment (++j)
	}

	// C. Infinite loop
	// for { break }

	// 2. Switch Deep Dive
	// - Implicit break (no fallthrough by default)
	// - Cases don't need to be constants
	// - Can switch on types (in interfaces, later topic)
	fmt.Println("\n\n--- Switch ---")
	
	day := "Sat"
	switch day {
	case "Mon", "Tue", "Wed", "Thu", "Fri":
		fmt.Println("Workday")
	case "Sat", "Sun":
		fmt.Println("Weekend")
	default: // Optional
		fmt.Println("Invalid day")
	}

	// Fallthrough
	// Use 'fallthrough' keyword to force execution of next case
	num := 1
	switch num {
	case 1:
		fmt.Println("One")
		fallthrough
	case 2:
		fmt.Println("Two (Printed because of fallthrough)")
	}

	// Switch true (If-else replacement)
	// Clean way to write long if-else chains
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("Good morning")
	case t.Hour() < 17:
		fmt.Println("Good afternoon")
	default:
		fmt.Println("Good evening")
	}

	// 3. Labelled Break/Continue
	// Useful for breaking out of nested loops
	fmt.Println("\n--- Labelled Break ---")
OuterLoop:
	for i := 0; i < 800; i++ {
		for j := 0; j < 3800; j++ {
			if i == 700 && j == 700 {
				fmt.Println("Breaking OuterLoop at", i, j)
				break OuterLoop // Breaks the OUTER loop, not just inner
			}
			fmt.Printf("%d-%d ", i, j)
		}
		fmt.Println()
	}

	// 4. Defer (Deep Peek)
	// We'll cover this more in Functions, but know that 'defer' schedules a function call
	// to run immediately before the surrounding function returns.
	// It's LIFO (Last In First Out).
	defer fmt.Println("\nThis prints LAST (LIFO)")
	defer fmt.Println("This prints Second to Last")
	defer fmt.Println("This prints Second to Last")
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Loop Variable Capture (Pre-Go 1.22):
   Pitfall: Capturing loop variable address in a closure/goroutine.
   for i := 0; i < 3; i++ { go func() { fmt.Println(i) }() }
   Old Behavior: Printed "3 3 3" (shared variable).
   New Behavior (Go 1.22+): Prints "0 1 2" (loop var is per-iteration).
   
2. Switch Fallthrough:
   Gotcha: Unlike C/Java, Go breaks by default.
   Error: Expecting fallthrough without 'fallthrough' keyword.

3. Unreachable Code:
   Error: "unreachable code"
   Why: Placing code after a 'return' or 'break' that is always executed.
*/
