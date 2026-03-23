// ============================================================================
// GO FUNDAMENTALS: CONTROL FLOW
// ============================================================================
// This file provides a comprehensive guide to Go's control structures,
// including conditionals, loops, switches, and jump statements.
// Go has fewer control structures than most languages but they're more flexible.
// ============================================================================

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  GO CONTROL FLOW                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: If Statements
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: If Statements")
	fmt.Println("─────────────────────────────────────────")

	// GO IF STATEMENT RULES:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. No parentheses around condition (they're optional but discouraged)    │
	// │ 2. Braces {} are REQUIRED (even for single statements)                   │
	// │ 3. Opening brace MUST be on same line (Go's automatic semicolon rule)    │
	// │ 4. Can have initialization statement before condition                    │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Basic if
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	}

	// if-else
	if x > 15 {
		fmt.Println("x is greater than 15")
	} else {
		fmt.Println("x is not greater than 15")
	}

	// if-else if-else chain
	score := 85
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: F")
	}

	// CRITICAL: If with initialization statement
	// Variables declared here are ONLY visible within the if-else block
	if n := computeValue(); n > 10 {
		fmt.Printf("Computed value %d is greater than 10\n", n)
	} else {
		fmt.Printf("Computed value %d is not greater than 10\n", n)
	}
	// n is NOT accessible here!

	// Common pattern: Error handling
	if result, err := mightFail(false); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	// VISUALIZATION: Scope of if-init variables
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ if val := compute(); val > 0 {                                        │
	// │     ┌─────────────────────────────────────────────────────────────┐   │
	// │     │ 'val' is visible here                                       │   │
	// │     └─────────────────────────────────────────────────────────────┘   │
	// │ } else {                                                              │
	// │     ┌─────────────────────────────────────────────────────────────┐   │
	// │     │ 'val' is also visible here                                  │   │
	// │     └─────────────────────────────────────────────────────────────┘   │
	// │ }                                                                     │
	// │ // 'val' is NOT visible here!                                         │
	// └────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 2: The For Loop (Go's Only Loop!)
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: The For Loop")
	fmt.Println("─────────────────────────────────────────")

	// Go has ONLY 'for' - no 'while' or 'do-while'!
	// But 'for' can express all loop types.

	// FORM 1: Classic C-style for loop
	fmt.Print("C-style: ")
	for i := 0; i < 5; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// FORM 2: While-style (condition only)
	fmt.Print("While-style: ")
	count := 0
	for count < 5 {
		fmt.Print(count, " ")
		count++
	}
	fmt.Println()

	// FORM 3: Infinite loop
	// for { ... } is equivalent to while(true) or for(;;)
	fmt.Print("Infinite (with break): ")
	i := 0
	for {
		if i >= 5 {
			break
		}
		fmt.Print(i, " ")
		i++
	}
	fmt.Println()

	// FORM 4: Range loop (for collections)
	fmt.Print("Range over slice: ")
	numbers := []int{10, 20, 30, 40, 50}
	for index, value := range numbers {
		fmt.Printf("[%d]=%d ", index, value)
	}
	fmt.Println()

	// Range variations
	fmt.Println("\nRange variations:")

	// Index only (ignore value)
	fmt.Print("  Index only: ")
	for i := range numbers {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Value only (ignore index)
	fmt.Print("  Value only: ")
	for _, v := range numbers {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// Range over string (iterates runes, not bytes!)
	fmt.Print("  Range over string: ")
	for i, r := range "Go语言" {
		fmt.Printf("[%d]='%c' ", i, r)
	}
	fmt.Println()

	// Range over map
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	fmt.Print("  Range over map: ")
	for key, value := range m {
		fmt.Printf("%s=%d ", key, value)
	}
	fmt.Println()
	// NOTE: Map iteration order is RANDOMIZED by design!

	// Range over channel
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch) // Must close to exit range loop

	fmt.Print("  Range over channel: ")
	for v := range ch {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// ========================================================================
	// SECTION 3: Loop Control (break, continue)
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Loop Control")
	fmt.Println("─────────────────────────────────────────")

	// break - exits the innermost loop
	fmt.Print("Break at 3: ")
	for i := 0; i < 10; i++ {
		if i == 3 {
			break
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// continue - skips to next iteration
	fmt.Print("Skip evens: ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// ========================================================================
	// SECTION 4: Labels (break/continue with labels)
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Labels")
	fmt.Println("─────────────────────────────────────────")

	// Labels allow breaking out of nested loops
	// SYNTAX: LabelName: for ... { }

	fmt.Println("Breaking outer loop at i=1, j=1:")
OuterLoop:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("  Breaking OuterLoop!")
				break OuterLoop // Breaks the OUTER loop
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}
	fmt.Println("After outer loop")

	// Continue with label
	fmt.Println("\nContinue outer loop at i=1, j=1:")
ContinueOuter:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("  Continuing to next i!")
				continue ContinueOuter // Continues the OUTER loop
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}

	// VISUALIZATION: Label scope
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Outer:                                                                 │
	// │ for i := 0; i < 3; i++ {      ◄── Label attached to this loop         │
	// │     Inner:                                                             │
	// │     for j := 0; j < 3; j++ {  ◄── Label attached to this loop         │
	// │         if condition {                                                 │
	// │             break Outer  ───────────────► Exits outer loop            │
	// │             break Inner  ───────────────► Exits inner loop (default)  │
	// │             break        ───────────────► Same as break Inner         │
	// │         }                                                              │
	// │     }                                                                  │
	// │ }                                                                      │
	// └────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 5: Switch Statement
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Switch Statement")
	fmt.Println("─────────────────────────────────────────")

	// GO SWITCH DIFFERENCES FROM C/JAVA:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. NO automatic fallthrough (implicit break after each case)             │
	// │ 2. Cases don't need to be constants                                      │
	// │ 3. Can switch on any comparable type (not just int)                      │
	// │ 4. Can have multiple values in one case                                  │
	// │ 5. Can have no expression (boolean switch)                               │
	// │ 6. Can have initialization statement                                     │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Basic switch
	day := "Tuesday"
	switch day {
	case "Monday":
		fmt.Println("Start of work week")
	case "Tuesday", "Wednesday", "Thursday": // Multiple values in one case
		fmt.Println("Midweek")
	case "Friday":
		fmt.Println("Almost weekend!")
	case "Saturday", "Sunday":
		fmt.Println("Weekend!")
	default:
		fmt.Println("Invalid day")
	}

	// Switch with initialization
	switch n := rand.Intn(10); {
	case n < 3:
		fmt.Printf("Random %d is low\n", n)
	case n < 7:
		fmt.Printf("Random %d is medium\n", n)
	default:
		fmt.Printf("Random %d is high\n", n)
	}

	// Tagless switch (switch true) - cleaner than if-else chain
	hour := time.Now().Hour()
	switch {
	case hour < 6:
		fmt.Println("Night time")
	case hour < 12:
		fmt.Println("Morning")
	case hour < 17:
		fmt.Println("Afternoon")
	case hour < 21:
		fmt.Println("Evening")
	default:
		fmt.Println("Night time")
	}

	// Fallthrough keyword (explicit fallthrough)
	num := 1
	fmt.Print("Fallthrough example: ")
	switch num {
	case 1:
		fmt.Print("one ")
		fallthrough // Explicitly fall through to next case
	case 2:
		fmt.Print("two ")
		fallthrough
	case 3:
		fmt.Print("three")
	}
	fmt.Println()
	// Note: fallthrough always executes next case, doesn't check condition!

	// Type switch (covered more in interfaces)
	printType(42)
	printType("hello")
	printType(3.14)
	printType([]int{1, 2, 3})

	// ========================================================================
	// SECTION 6: Defer Statement
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Defer Statement")
	fmt.Println("─────────────────────────────────────────")

	// defer schedules a function call to run when the function returns
	// Deferred calls are executed in LIFO (Last In, First Out) order

	fmt.Println("Defer execution order:")
	deferDemo()

	// CRITICAL: Defer arguments are evaluated IMMEDIATELY
	fmt.Println("\nDefer argument evaluation:")
	argValue := 10
	defer fmt.Printf("  Deferred print (value was %d when scheduled)\n", argValue)
	argValue = 20
	fmt.Printf("  Current value: %d\n", argValue)
	// The deferred call will print 10, not 20!

	// Common use: Resource cleanup
	// file, err := os.Open("file.txt")
	// if err != nil { return err }
	// defer file.Close()  // Guaranteed to run when function returns

	// DEFER STACK VISUALIZATION:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Function execution:                                                     │
	// │                                                                          │
	// │ defer A()  ───┐                                                         │
	// │ defer B()  ───┼── Added to defer stack                                  │
	// │ defer C()  ───┘                                                         │
	// │ ...                                                                      │
	// │ return      ◄── Function about to return                                │
	// │                                                                          │
	// │ Defer Stack:           Execution Order:                                 │
	// │ ┌─────────┐            1. C() runs first                                │
	// │ │ C()     │ ◄── Top    2. B() runs second                               │
	// │ │ B()     │            3. A() runs last                                 │
	// │ │ A()     │ ◄── Bottom                                                  │
	// │ └─────────┘                                                             │
	// └──────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 7: Goto Statement
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Goto Statement")
	fmt.Println("─────────────────────────────────────────")

	// Go has goto, but with strict rules to prevent spaghetti code:
	// 1. Cannot jump into a block
	// 2. Cannot jump over variable declarations
	// 3. Target label must be in same function

	// Legitimate use case: Error cleanup in generated code
	fmt.Println("Goto example:")
	gotoDemo()

	// In practice, goto is rarely needed. Use loops, break, continue instead.

	// ========================================================================
	// SECTION 8: Loop Variable Capture (Go 1.22+ Change!)
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Loop Variable Capture")
	fmt.Println("─────────────────────────────────────────")

	// HISTORICAL BUG (Pre-Go 1.22):
	// Loop variables were shared across iterations in closures.
	//
	// OLD BEHAVIOR:
	// for i := 0; i < 3; i++ {
	//     go func() { fmt.Println(i) }()  // Would print "3 3 3"
	// }
	//
	// NEW BEHAVIOR (Go 1.22+):
	// Each iteration has its own copy of the loop variable.

	fmt.Println("Go 1.22+ behavior (each iteration has own variable):")
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			fmt.Printf("  Goroutine sees i=%d\n", i)
			done <- true
		}()
	}
	// Wait for all goroutines
	<-done
	<-done
	<-done

	// For compatibility with older Go or explicit behavior:
	fmt.Println("\nExplicit copy (works in all Go versions):")
	for i := 0; i < 3; i++ {
		i := i // Shadow loop variable with local copy
		go func() {
			fmt.Printf("  Goroutine sees i=%d\n", i)
		}()
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Control Flow Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// Helper functions for demonstrations

func computeValue() int {
	return 15
}

func mightFail(fail bool) (string, error) {
	if fail {
		return "", fmt.Errorf("operation failed")
	}
	return "success", nil
}

func deferDemo() {
	defer fmt.Println("  Third (last defer, runs first)")
	defer fmt.Println("  Second")
	defer fmt.Println("  First (first defer, runs last)")
	fmt.Println("  Function body executing...")
}

func printType(v interface{}) {
	switch t := v.(type) {
	case int:
		fmt.Printf("Type switch: %v is int\n", t)
	case string:
		fmt.Printf("Type switch: %v is string\n", t)
	case float64:
		fmt.Printf("Type switch: %v is float64\n", t)
	default:
		fmt.Printf("Type switch: %v is unknown type %T\n", v, v)
	}
}

func gotoDemo() {
	i := 0
Loop:
	if i < 3 {
		fmt.Printf("  i=%d\n", i)
		i++
		goto Loop
	}
	fmt.Println("  Loop complete")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. MISSING BRACES
   ─────────────────────────────────────────────────────────────────────────
   Error: "syntax error: unexpected newline, expecting { after if clause"

   WRONG:
   if x > 5
       fmt.Println("big")  // No braces!

   RIGHT:
   if x > 5 {
       fmt.Println("big")
   }

   Go ALWAYS requires braces, even for single statements.

2. USING 'while' OR 'do-while'
   ─────────────────────────────────────────────────────────────────────────
   Error: "syntax error: unexpected while"

   Go has NO 'while' keyword. Use 'for' instead:

   WRONG: while x < 10 { ... }
   RIGHT: for x < 10 { ... }

3. EXPECTING SWITCH FALLTHROUGH
   ─────────────────────────────────────────────────────────────────────────
   In Go, each case automatically breaks.

   UNEXPECTED:
   switch n {
   case 1:
       fmt.Println("one")
   case 2:
       fmt.Println("two")  // Never executes when n=1
   }

   Use 'fallthrough' keyword if you want C-style behavior.

4. FALLTHROUGH DOESN'T CHECK CONDITION
   ─────────────────────────────────────────────────────────────────────────
   switch n {
   case 1:
       fallthrough
   case 2:  // Executes even if n != 2 (after fallthrough from case 1)
       fmt.Println("reached")
   }

   'fallthrough' unconditionally executes next case body.

5. SHADOWING IN IF INITIALIZATION
   ─────────────────────────────────────────────────────────────────────────
   var err error
   if err := doSomething(); err != nil {
       // This 'err' is a NEW variable, shadows outer err
   }
   // Outer 'err' is still nil!

6. INFINITE LOOP WITHOUT EXIT
   ─────────────────────────────────────────────────────────────────────────
   for {
       // No break, return, or os.Exit
       // Program hangs forever!
   }

   Always ensure infinite loops have an exit condition.

7. RANGE OVER NIL COLLECTION
   ─────────────────────────────────────────────────────────────────────────
   This is actually SAFE - it just iterates zero times:
   var slice []int  // nil slice
   for _, v := range slice {  // No panic, just skips
       fmt.Println(v)
   }

8. MODIFYING COLLECTION DURING RANGE
   ─────────────────────────────────────────────────────────────────────────
   For slices: Appending during range may cause unexpected behavior.
   For maps: Adding/deleting during range is actually defined behavior,
             but items may or may not appear in the iteration.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. PREFER switch OVER LONG if-else CHAINS
   Cleaner, more readable, and sometimes faster.

2. USE TAGLESS switch FOR COMPLEX CONDITIONS
   switch {
   case a && b:
   case c || d:
   }
   More readable than nested if-else.

3. USE break WITH LABELS FOR NESTED LOOPS
   Instead of flags or goto, use labeled break.

4. ALWAYS USE defer FOR CLEANUP
   Guarantees cleanup even if function panics.
   Place defer right after resource acquisition.

5. AVOID goto IN NEW CODE
   Use loops and breaks instead.
   goto may be acceptable in generated code or specific patterns.

6. USE for range FOR COLLECTIONS
   More readable and less error-prone than index-based loops.

7. UNDERSTAND LOOP VARIABLE CAPTURE
   In older Go or when explicit about behavior, shadow loop variables
   before capturing in closures: i := i

8. KEEP SWITCH CASES SIMPLE
   If a case has complex logic, extract to a function.

9. USE DEFAULT CASE IN SWITCHES
   Handle unexpected values explicitly.

10. AVOID DEEP NESTING
    If you have deeply nested if/for statements, refactor.
    Extract inner logic to functions or use early returns.
*/
