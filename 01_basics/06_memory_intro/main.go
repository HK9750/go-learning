package main

import "fmt"

// DEEP DIVE: Memory Management (Stack vs Heap)
// Go has a Garbage Collector (GC), but understanding allocation is pro-level.

// Stack: Fast allocation/deallocation (moving a pointer). Local variables usually live here.
// Heap: Slower. Shared variables, dynamic size, or those that "escape" live here. GC scans this.

// 1. Escape Analysis
// The compiler determines where to store a variable.
// If a variable's reference is returned from a function, it "escapes" to the heap.

func leakingPointer() *int {
	x := 100 // 'x' would normally be on stack...
	return &x // But we return a pointer to it.
	// In C/C++, this is a dangling pointer bug (stack frame destroyed).
	// In Go, the compiler sees this and allocates 'x' on the HEAP. Safe.
}

func stayOnStack() int {
	x := 200
	return x // Value copied. 'x' stays on stack.
}

func main() {
	// Pointers refresher
	a := 42
	p := &a         // '&' generates a pointer (address)
	fmt.Println(*p) // '*' dereferences the pointer (reads value)
	*p = 21         // Change value through pointer
	fmt.Println(a)  // a is now 21

	// Heap vs Stack
	ptr := leakingPointer()
	fmt.Printf("Value from heap: %d\n", *ptr)
	
	val := stayOnStack()
	fmt.Printf("Value from stack: %d\n", val)

	// How to see this?
	// Run: go build -gcflags="-m" main.go
	// You will see output like:
	// "./main.go:15:2: moved to heap: x"
}

/*
PRO TIP:
Don't fear the Heap, but respect it. Use pointers (*) to share data or modify state,
not just to avoid copying struct.
Copying small structs on stack is often FASTER than passing pointers (which cause heap allocations + GC pressure).
*/

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Nil Pointer Dereference:
   Error: "panic: runtime error: invalid memory address or nil pointer dereference"
   Code: var p *int; *p = 10
   Why: 'p' is nil. You cannot write to address 0.
   Fix: Initialize pointer (p = new(int) or p = &val) or check for nil.

2. Premature Optimization:
   Anti-Pattern: Passing everything by pointer to "save memory".
   Reality: Passing small structs by value is often faster (stack copy) than pointer (heap alloc + GC).
   Rule: Use pointers for Sharing (mutation) or Large Structs (>64 bytes).

3. Stack Overflow:
   Error: "runtime: goroutine stack exceeds 1000000000-byte limit"
   Why: Infinite recursion or massive stack allocation.
*/
