package main

import "fmt"

// DEEP DIVE: Functions & Closures

// 1. Naked Returns
// Named return parameters are initialized to their zero values.
// A "naked" return (return without arguments) returns those values.
// Use sparingly for long functions; it hurts readability.
func split(sum int) (x, y int) { // x, y initialized to 0
	x = sum * 4 / 9
	y = sum - x
	y = sum - x
	return // Returns x, y. 
	// BEST PRACTICE: Avoid naked returns in large functions. 
	// They obscure what is actually being returned.
}

// 2. Variadic Functions
// 'nums' is a []int inside the function.
func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func main() {
	// A. Calling basic functions
	x, y := split(17)
	fmt.Printf("Split 17: %d, %d\n", x, y)

	// B. Variadic
	fmt.Println("Sum:", sum(1, 2, 3))
	// Unpacking a slice
	s := []int{10, 20, 30}
	fmt.Println("Sum slice:", sum(s...))

	// C. Anonymous Functions & Closures
	// A closure is a function value that references variables from outside its body.
	fmt.Println("\n--- Closures ---")
	
	multiplier := func(factor int) func(int) int {
		return func(val int) int {
			return val * factor // 'factor' is captured here
		}
	}

	timesTwo := multiplier(2)
	timesTen := multiplier(10)

	fmt.Println("5 * 2 =", timesTwo(5))
	fmt.Println("5 * 10 =", timesTen(5))

	// DEEP DIVE: Defer Arguments Evaluation
	// Arguments to defer are evaluated IMMEDIATELY, not when the function runs.
	fmt.Println("\n--- Defer Magic ---")
	a := 10
	defer fmt.Println("Deferred Print (Value of a):", a) // Captures a=10 NOW.
	a = 20
	fmt.Println("Current Value of a:", a)
	// 1. Current Value: 20
	// 2. Deferred Print: 10 (Because it was captured at scheduling time)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Nil Function Call:
   Error: "panic: runtime error: invalid memory address or nil pointer dereference"
   Why: Calling a function variable that is nil.
   Fix: Check if funcVar != nil before calling.

2. Variable shadowing in Named Returns:
   Risk: Declaring a new variable with same name as return parameter inside function blocks.
   func example() (count int) {
       count := 10 // NEW variable, not the return one!
       return // Returns 0 (the outer 'count'), not 10.
   }

3. Defer Loop Pitfall:
   Warning: Defer in a tight loop can cause stack overflow or memory buildup (defers run at *function* exit, not loop exit).
   Fix: Wrap loop body in an anonymous function if you need immediate cleanup.
*/
