package main

import (
	"fmt"
)

// TOPIC: VARIABLE SCOPES
var packageVar = "I am available throughout the package"

// DEEP DIVE: Shadowing
// Shadowing occurs when a variable declared in an inner scope has the same name as a variable in an outer scope.

func main() {
	// 1. Declaration styles
	var explicit int = 10         // Explicit typing
	inferred := 20                // Short declaration (type inferred as int). ONLY inside functions.
	var zeroValueString string    // Declared but not initialized.
	var zeroValueInt int          // Declared but not initialized.
	var zeroValueBool bool        // Declared but not initialized.
	var zeroValuePtr *int         // Declared but not initialized.

	fmt.Println(inferred)
	fmt.Printf("Zero Values -> String: '%q', Int: %d, Bool: %v, Ptr: %v\n", 
		zeroValueString, zeroValueInt, zeroValueBool, zeroValuePtr)
	// Output: String: '', Int: 0, Bool: false, Ptr: <nil>

	// 2. Type System Strictness
	// Go is statically typed. You cannot assign a float to an int without explicit conversion.
	var myFloat float64 = 3.14 
	// explicit = myFloat // COMPILER ERROR: cannot use myFloat (type float64) as type int in assignment
	explicit = int(myFloat) // Correct: Explicit Type Conversion (truncates, doesn't round)
	fmt.Println("Truncated float:", explicit)

	// 3. Numeric Types Deep Dive
	// int vs int64: 'int' is platform dependent (32-bit on 32-bit systems, 64-bit on 64-bit systems).
	// 'int64' is always 64-bit. They are DIFFERENT types.
	var x int = 10
	var y int64 = 10
	// if x == y { } // COMPILER ERROR: mismatched types int and int64

	fmt.Printf("Type of x: %T, Type of y: %T\n", x, y) // %T prints the type information

	// CRITICAL: Type Conversion
	// Go does not support implicit type casting. You must cast explicitly.
	// Bad: var f float64 = 10; var i int = f; (Error)
	// Good: var i int = int(f)

	// 4. Constants
	// Constants are evaluated at COMPILE TIME.
	// They have arbitrary precision until assigned to a variable.
	const Pi = 3.14159265358979323846264338327950288419716939937510582097494459
	// Implicitly, 'Pi' is an 'untyped float' with high precision.
	fmt.Println("Pi:", Pi)
	// const invalid = math.Sin(1.57) // ERROR: math.Sin is a runtime function call, constants must be compile-time.

	// 5. Variable Shadowing Pitfall
	n := 10
	if n > 5 {
		fmt.Println("\n--- Shadowing Zone ---")
		n := 5 // New variable 'n' shadows the outer 'n'
		
		// MEMORY VISUALIZATION:
		// Stack Frame:
		// | ...            |
		// | 0xAddr1: 10    | <- Outer 'n' (still exists, but hidden)
		// | 0xAddr2: 5     | <- Inner 'n' (active in this scope)
		// | ...            |
		
		fmt.Printf("Inner n: %d (Address: %p)\n", n, &n)
		n++ // Affects inner 'n' only
	}
	fmt.Printf("Outer n: %d (Address: %p)\n", n, &n) // Outer 'n' remains 10
	// Notice the memory addresses are different!

	// 6. Type Aliasing vs Type Definition
	//
	// VISUALIZATION:
	// type UserID int      ->  [ UserID Type ] --(distinct from)--> [ int Type ]
	// type LegacyID = int  ->  [ LegacyID Tag ] ------------------> [ int Type ] (Same thing)
	
	type UserID int         // NEW TYPE: 'UserID' is distinct from 'int'.
	type LegacyID = int     // ALIAS: 'LegacyID' is just another name for 'int'.

	var myID UserID = 100
	var oldID LegacyID = 100
	var plainInt int = 100

	// plainInt = myID // ERROR: different types
	plainInt = int(myID) // OK: conversion required
	plainInt = oldID     // OK: aliases are identical types

	fmt.Println(plainInt)
	fmt.Println("\n--- Type Identity ---")
	fmt.Printf("UserID type: %T\n", myID)
	fmt.Printf("LegacyID type: %T\n", oldID)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Short Declaration (:=) outside functions:
   Error: "syntax error: non-declaration statement outside function body"
   Why: ':=' is only valid inside functions. At package level, use 'var'.

2. Accidental Shadowing:
   Risk: Logic bugs where you think you're updating a variable but are updating a new local one.
   Pattern to Avoid:
     var err error
     if condition {
         err := doSomething() // Shadows the outer 'err'!
         // Outer 'err' remains nil.
     }
   Fix: Use '=' instead of ':=' if you intend to reuse the variable.
     err = doSomething()

3. Unused Variables:
   Error: "x declared but not used"
   Why: Go compilers are strict.
   Fix: Use the variable or assign to blank identifier '_'.
     _ = x
*/
