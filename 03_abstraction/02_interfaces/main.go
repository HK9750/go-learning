package main

import (
	"fmt"
	"unsafe"
)

// DEEP DIVE: Interfaces
// Interfaces are implicitly implemented. No "implements" keyword.
// Interface Value = (Type, Value) tuple.

type Speaker interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d *Dog) Speak() string {
	if d == nil {
		return "<nil dog>"
	}
	return "Woof!"
} 

type Cat struct {}
func (c Cat) Speak() string { return "Meow" }

// The "Typed Nil" Pitfall
// An interface holding a nil concrete pointer is NOT nil itself.
func getSpeaker(isDog bool) Speaker {
	if isDog {
		var d *Dog = nil // Pointer is nil
		return d         // returns interface{Type=*Dog, Value=nil}
	}
	return nil // returns interface{Type=nil, Value=nil}
}

func main() {
	// 1. Basic Polymorphism
	animals := []Speaker{&Dog{"Rex"}, Cat{}}
	for _, a := range animals {
		fmt.Println(a.Speak())
	}

	// 2. Empty Interface (Any)
	// interface{} or 'any' (Go 1.18+) can hold anything.
	var i interface{} = "hello"
	fmt.Printf("Dynamic Type: %T, Value: %v\n", i, i)

	// 3. Type Assertions
	// x.(T) asserts that x is not nil and the concrete value is of type T.
	s, ok := i.(string) // Safe assertion
	if ok {
		fmt.Println("It's a string:", s)
	}

	// 4. The Typed Nil Trap
	fmt.Println("\n--- Typed Nil Trap ---")
	s1 := getSpeaker(true)
	
	if s1 == nil {
		fmt.Println("s1 is nil")
	} else {
		fmt.Println("s1 is NOT nil (Wait, why?)")
		fmt.Printf("Deep inspection: Type=%T, Value=%v\n", s1, s1)
		// Explanation: s1 is (Type=*Dog, Value=nil).
		// An interface is nil ONLY if both Type and Value are nil.
	}
	
	// Safe to call method because *Dog implementation handles nil receiver!
	fmt.Println("Result:", s1.Speak()) 
	fmt.Println("Result:", s1.Speak()) 

	// =========================================================================
	// 5. DEEP DIVE: Interface Internals (Runtime View)
	// =========================================================================
	// Interfaces are just 2 words in memory: (TablePtr, DataPtr).
	// Let's look at them using unsafe!

	fmt.Println("\n--- Deep Dive: Interface Internals ---")
	
	// Create an interface
	var s2 Speaker = &Dog{"Buddy"}
	
	// Define a struct that matches the runtime layout of a non-empty interface (iface)
	// src/runtime/runtime2.go: type iface struct { tab *itab; data unsafe.Pointer }
	type ifaceHeader struct {
		Tab  uintptr
		Data uintptr
	}
	
	// Cast the interface variable to our header struct
	// We get the address of s2, cast to unsafe.Pointer, then to *ifaceHeader
	ptr := (*ifaceHeader)(unsafe.Pointer(&s2))
	
	fmt.Printf("Interface Variable 's2' (Address: %p):\n", &s2)
	fmt.Printf("  [Word 1] itab (Method Table) : %#x\n", ptr.Tab) // Points to function list + type info
	fmt.Printf("  [Word 2] data (Instance Ptr) : %#x\n", ptr.Data) // Points to the Dog struct
	
	// Verify Data ptr matches the actual Dog
	dogPtr := s2.(*Dog)
	fmt.Printf("  Actual Dog Pointer           : %p\n", dogPtr)
	
	// Visualizing the Dispatch
	// s2.Speak() -> essentially does: ptr.Tab.Fun[0](ptr.Data)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. The Typed Nil (Most Common Interface Bug):
   Error: Interface check 'if s1 != nil' is true, but concrete value is nil.
   Why: Interface is a tuple (Type, Value). (Dog*, nil) != (nil, nil).
   Fix: Always return 'nil' explicitly from functions returning interfaces, rather than a nil concrete pointer.

2. Interface Pollution:
   Anti-Pattern: Defining large interfaces (e.g., 20 methods).
   Best Practice: "The bigger the interface, the weaker the abstraction." - Rob Pike. Keep them small (1-3 methods).

3. Type Assertion Panic:
   Error: "panic: interface conversion: interface {} is int, not string"
   Code: s := i.(string) // without comma-ok
   Fix: Always use comma-ok idiom: 's, ok := i.(string); if !ok { ... }'
*/
