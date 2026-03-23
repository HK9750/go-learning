// ============================================================================
// GO FUNDAMENTALS: POINTERS
// ============================================================================
// This file provides a comprehensive guide to pointers in Go, including
// memory addresses, dereferencing, pointer receivers, and common patterns.
// ============================================================================

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                     GO POINTERS                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Pointer Basics
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Pointer Basics")
	fmt.Println("─────────────────────────────────────────")

	// WHAT IS A POINTER?
	// A pointer holds the MEMORY ADDRESS of a value.
	//
	// VISUALIZATION:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ MEMORY                                                                 │
	// │                                                                        │
	// │ Address      Value       Variable                                      │
	// │ ┌──────────┬───────────┬────────────┐                                 │
	// │ │ 0x1000   │    42     │    x       │  ◄── int variable               │
	// │ ├──────────┼───────────┼────────────┤                                 │
	// │ │ 0x1008   │  0x1000   │    p       │  ◄── pointer to x (*int)        │
	// │ └──────────┴───────────┴────────────┘                                 │
	// │                                                                        │
	// │ p holds the ADDRESS 0x1000, which is where x's VALUE (42) is stored   │
	// └────────────────────────────────────────────────────────────────────────┘

	// OPERATORS:
	// & (address-of): Gets the memory address of a variable
	// * (dereference): Gets the value at a memory address

	x := 42
	p := &x // p is a pointer to x (type: *int)

	fmt.Printf("x value:      %d\n", x)
	fmt.Printf("x address:    %p\n", &x)
	fmt.Printf("p value:      %p (same as x's address)\n", p)
	fmt.Printf("*p (deref):   %d (same as x's value)\n", *p)
	fmt.Printf("p type:       %T\n", p)

	// Modifying through pointer
	*p = 100 // Changes x through the pointer
	fmt.Printf("\nAfter *p = 100:\n")
	fmt.Printf("x value:      %d\n", x)
	fmt.Printf("*p value:     %d\n", *p)

	// ========================================================================
	// SECTION 2: Pointer Zero Value (nil)
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Pointer Zero Value (nil)")
	fmt.Println("─────────────────────────────────────────")

	// The zero value of a pointer is nil
	var nilPtr *int
	fmt.Printf("Nil pointer:  %v\n", nilPtr)
	fmt.Printf("Is nil:       %t\n", nilPtr == nil)

	// CRITICAL: Dereferencing nil pointer causes PANIC!
	// *nilPtr = 10  // PANIC: runtime error: invalid memory address

	// Always check for nil before dereferencing
	if nilPtr != nil {
		fmt.Println("Value:", *nilPtr)
	} else {
		fmt.Println("Pointer is nil, cannot dereference")
	}

	// Safe pattern with default
	value := safeDeref(nilPtr, 0)
	fmt.Printf("Safe deref of nil: %d\n", value)

	// ========================================================================
	// SECTION 3: Creating Pointers
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Creating Pointers")
	fmt.Println("─────────────────────────────────────────")

	// Method 1: Address of existing variable
	a := 10
	p1 := &a
	fmt.Printf("&a:     %p -> %d\n", p1, *p1)

	// Method 2: new() function
	// new(T) allocates zeroed storage for T and returns *T
	p2 := new(int) // Allocates int, returns *int, value is 0
	fmt.Printf("new(int): %p -> %d\n", p2, *p2)
	*p2 = 20
	fmt.Printf("After assignment: %p -> %d\n", p2, *p2)

	// Method 3: Pointer to literal (composite literals only)
	p3 := &Person{Name: "Alice", Age: 30}
	fmt.Printf("&literal: %p -> %+v\n", p3, *p3)

	// NOTE: Cannot take address of non-composite literals
	// p4 := &42  // COMPILE ERROR!
	// p5 := &"hello"  // COMPILE ERROR!

	// Method 4: Make pointer type explicit
	var p4 *int = new(int)
	*p4 = 42
	fmt.Printf("Explicit type: %p -> %d\n", p4, *p4)

	// ========================================================================
	// SECTION 4: Pointer Arithmetic (NOT Allowed!)
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Pointer Arithmetic (Restricted)")
	fmt.Println("─────────────────────────────────────────")

	// Go does NOT allow pointer arithmetic (unlike C/C++)
	// p++ is not allowed
	// p + 1 is not allowed
	// p1 - p2 is not allowed

	arr := [3]int{10, 20, 30}
	arrPtr := &arr[0]
	fmt.Printf("arr[0] address: %p\n", arrPtr)

	// WRONG (doesn't compile):
	// nextPtr := arrPtr + 1  // ERROR!

	// To access array elements, use indexing:
	fmt.Printf("arr[1]: %d\n", arr[1])

	// The unsafe package allows pointer arithmetic, but is... unsafe
	fmt.Println("\nUsing unsafe (not recommended):")
	unsafePtr := unsafe.Pointer(&arr[0])
	nextPtr := unsafe.Pointer(uintptr(unsafePtr) + unsafe.Sizeof(arr[0]))
	nextVal := *(*int)(nextPtr)
	fmt.Printf("arr[1] via unsafe: %d\n", nextVal)

	// ========================================================================
	// SECTION 5: Pointers and Functions
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Pointers and Functions")
	fmt.Println("─────────────────────────────────────────")

	// PASS BY VALUE VS PASS BY POINTER:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Pass by Value                 │ Pass by Pointer                       │
	// │ func foo(x int)               │ func bar(x *int)                      │
	// │ - Copy of value passed        │ - Address passed                      │
	// │ - Original unchanged          │ - Can modify original                 │
	// │ - Safe, no side effects       │ - Explicit mutation                   │
	// │ - May copy large structs      │ - Only copies 8 bytes (address)       │
	// └────────────────────────────────────────────────────────────────────────┘

	num := 10
	fmt.Printf("Before functions: num = %d\n", num)

	// Pass by value - cannot modify original
	incrementByValue(num)
	fmt.Printf("After incrementByValue: num = %d (unchanged)\n", num)

	// Pass by pointer - can modify original
	incrementByPointer(&num)
	fmt.Printf("After incrementByPointer: num = %d (changed!)\n", num)

	// When to use pointers for function parameters:
	// 1. You need to modify the argument
	// 2. The argument is large (avoid copying)
	// 3. The argument may be nil (optional parameter)
	// 4. Consistency with method receivers

	// Example: Swap function requires pointers
	a1, b1 := 5, 10
	fmt.Printf("\nBefore swap: a=%d, b=%d\n", a1, b1)
	swap(&a1, &b1)
	fmt.Printf("After swap: a=%d, b=%d\n", a1, b1)

	// ========================================================================
	// SECTION 6: Pointers to Structs
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Pointers to Structs")
	fmt.Println("─────────────────────────────────────────")

	// Go provides syntactic sugar for struct pointers

	person := Person{Name: "Bob", Age: 25}
	personPtr := &person

	// Accessing fields through pointer
	// Explicit dereference:
	fmt.Printf("(*personPtr).Name = %s\n", (*personPtr).Name)

	// Implicit dereference (syntactic sugar - preferred):
	fmt.Printf("personPtr.Name = %s\n", personPtr.Name)

	// Both are equivalent! Go automatically dereferences.

	// Modifying through pointer
	personPtr.Age = 26
	fmt.Printf("After personPtr.Age = 26: person = %+v\n", person)

	// Common pattern: Constructor returning pointer
	p5 := NewPerson("Charlie", 35)
	fmt.Printf("NewPerson result: %+v\n", *p5)

	// ========================================================================
	// SECTION 7: Method Receivers (Value vs Pointer)
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Method Receivers")
	fmt.Println("─────────────────────────────────────────")

	// VALUE RECEIVER: Method gets COPY of the struct
	// POINTER RECEIVER: Method gets POINTER to the struct

	// CHOOSING RECEIVER TYPE:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Use POINTER receiver when:                                              │
	// │ 1. Method needs to modify the receiver                                  │
	// │ 2. Receiver is large (avoid copying)                                    │
	// │ 3. Consistency: if any method needs pointer, use pointer for all        │
	// │                                                                          │
	// │ Use VALUE receiver when:                                                │
	// │ 1. Receiver is small (int, small struct)                                │
	// │ 2. Method doesn't modify receiver                                       │
	// │ 3. Receiver is a map, func, or chan (already reference types)           │
	// │ 4. Receiver is a slice that doesn't need to be resliced                 │
	// └──────────────────────────────────────────────────────────────────────────┘

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Rectangle: %+v\n", rect)
	fmt.Printf("Area (value receiver): %.2f\n", rect.Area())

	rect.Scale(2) // Pointer receiver - modifies rect
	fmt.Printf("After Scale(2): %+v\n", rect)

	// Go automatically takes address/dereferences as needed
	// Both work:
	fmt.Printf("rect.Area(): %.2f\n", rect.Area())
	fmt.Printf("(&rect).Area(): %.2f\n", (&rect).Area())

	rectPtr := &rect
	fmt.Printf("rectPtr.Area(): %.2f\n", rectPtr.Area())
	fmt.Printf("(*rectPtr).Area(): %.2f\n", (*rectPtr).Area())

	// ========================================================================
	// SECTION 8: Pointers to Slices and Maps
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Pointers to Slices and Maps")
	fmt.Println("─────────────────────────────────────────")

	// Slices and maps are ALREADY reference types!
	// They contain internal pointers to underlying data.
	// Usually, you don't need pointers to slices/maps.

	// SLICE HEADER:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Slice                                                                  │
	// │ ┌──────────────────┐     ┌─────────────────────────────────┐          │
	// │ │ Pointer ─────────┼────►│ underlying array data           │          │
	// │ │ Length           │     └─────────────────────────────────┘          │
	// │ │ Capacity         │                                                   │
	// │ └──────────────────┘                                                   │
	// └────────────────────────────────────────────────────────────────────────┘

	slice := []int{1, 2, 3}
	modifySliceElements(slice) // Elements modified
	fmt.Printf("After modifySliceElements: %v\n", slice)

	// But to append or reslice, you need pointer or return new slice
	appendToSlice(&slice, 4, 5)
	fmt.Printf("After appendToSlice: %v\n", slice)

	// Maps are similar - already reference type
	m := map[string]int{"a": 1}
	modifyMap(m)
	fmt.Printf("After modifyMap: %v\n", m)

	// ========================================================================
	// SECTION 9: Pointer Comparison
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Pointer Comparison")
	fmt.Println("─────────────────────────────────────────")

	// Pointers can be compared with == and !=
	// They're equal if they point to same memory location

	x1, x2 := 10, 10
	p6, p7 := &x1, &x1
	p8 := &x2

	fmt.Printf("p6 == p7 (same var): %t\n", p6 == p7)
	fmt.Printf("p6 == p8 (diff var): %t\n", p6 == p8)
	fmt.Printf("*p6 == *p8 (same val): %t\n", *p6 == *p8)

	// nil comparison
	var nilP *int
	fmt.Printf("nilP == nil: %t\n", nilP == nil)

	// ========================================================================
	// SECTION 10: Double Pointers
	// ========================================================================
	fmt.Println("\n▶ SECTION 10: Double Pointers")
	fmt.Println("─────────────────────────────────────────")

	// Pointer to pointer (rarely needed in Go)
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │      x           p           pp                                       │
	// │ ┌────────┐  ┌────────┐  ┌────────┐                                    │
	// │ │   42   │◄─┤ &x     │◄─┤ &p     │                                    │
	// │ └────────┘  └────────┘  └────────┘                                    │
	// │    int        *int        **int                                       │
	// └────────────────────────────────────────────────────────────────────────┘

	xVal := 42
	ptrX := &xVal
	ptrPtr := &ptrX // **int

	fmt.Printf("xVal:      %d\n", xVal)
	fmt.Printf("*ptrX:     %d\n", *ptrX)
	fmt.Printf("**ptrPtr:  %d\n", **ptrPtr)

	// Modifying through double pointer
	**ptrPtr = 100
	fmt.Printf("After **ptrPtr = 100, xVal = %d\n", xVal)

	// Use case: Changing which object a pointer points to
	var node *Node
	createNode(&node, 42)
	fmt.Printf("After createNode: node = %+v\n", *node)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Pointers Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// Supporting Types and Functions
// ============================================================================

type Person struct {
	Name string
	Age  int
}

// Constructor pattern - returns pointer
func NewPerson(name string, age int) *Person {
	return &Person{Name: name, Age: age}
}

type Rectangle struct {
	Width, Height float64
}

// Value receiver - doesn't modify
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Pointer receiver - can modify
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

type Node struct {
	Value int
	Next  *Node
}

// Pass by value - cannot modify original
func incrementByValue(x int) {
	x++
	fmt.Printf("  Inside incrementByValue: x = %d\n", x)
}

// Pass by pointer - can modify original
func incrementByPointer(x *int) {
	*x++
	fmt.Printf("  Inside incrementByPointer: *x = %d\n", *x)
}

// Swap requires pointers
func swap(a, b *int) {
	*a, *b = *b, *a
}

// Safe dereference with default
func safeDeref(p *int, defaultVal int) int {
	if p == nil {
		return defaultVal
	}
	return *p
}

// Modifying slice elements (works without pointer)
func modifySliceElements(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

// Appending to slice requires pointer or return
func appendToSlice(s *[]int, vals ...int) {
	*s = append(*s, vals...)
}

// Modifying map (works without pointer)
func modifyMap(m map[string]int) {
	m["b"] = 2
}

// Creating node through double pointer
func createNode(node **Node, value int) {
	*node = &Node{Value: value}
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. NIL POINTER DEREFERENCE
   ─────────────────────────────────────────────────────────────────────────
   Error: "panic: runtime error: invalid memory address or nil pointer dereference"

   var p *int
   *p = 10  // PANIC!

   Fix: Initialize pointer or check for nil before dereferencing.

2. RETURNING POINTER TO LOCAL VARIABLE (Actually SAFE in Go!)
   ─────────────────────────────────────────────────────────────────────────
   In C/C++, this is a bug. In Go, it's SAFE:

   func getPointer() *int {
       x := 42
       return &x  // Go's escape analysis moves x to heap
   }

   Go's compiler detects this and allocates on heap instead of stack.

3. POINTER TO LOOP VARIABLE (Pre-Go 1.22)
   ─────────────────────────────────────────────────────────────────────────
   BAD (older Go):
   var ptrs []*int
   for i := 0; i < 3; i++ {
       ptrs = append(ptrs, &i)  // All point to same variable!
   }

   Fix: Create local copy:
   for i := 0; i < 3; i++ {
       i := i  // Shadow
       ptrs = append(ptrs, &i)
   }

4. FORGETTING TO DEREFERENCE
   ─────────────────────────────────────────────────────────────────────────
   p := &someInt
   p = 42  // ERROR: cannot use 42 as *int

   Fix: *p = 42

5. COMPARING POINTERS WHEN YOU MEAN VALUES
   ─────────────────────────────────────────────────────────────────────────
   a, b := 10, 10
   if &a == &b {  // Always false! Different addresses
       // ...
   }

   Fix: Compare values: if a == b { ... }

6. UNNECESSARY POINTER TO POINTER
   ─────────────────────────────────────────────────────────────────────────
   Go code rarely needs **T. Consider:
   - Returning new pointer instead
   - Using a struct wrapper
   - Returning error along with result

7. POINTER TO INTERFACE (ALMOST NEVER NEEDED)
   ─────────────────────────────────────────────────────────────────────────
   var w *io.Writer  // Almost always wrong!

   Interfaces already hold a pointer internally.
   Use: var w io.Writer
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. CHECK NIL BEFORE DEREFERENCING
   if p != nil {
       value := *p
   }

2. USE POINTER RECEIVERS FOR MUTATION
   If a method modifies the receiver, use pointer receiver.

3. BE CONSISTENT WITH RECEIVER TYPES
   If one method uses pointer receiver, all should.

4. RETURN *T FOR CONSTRUCTORS
   func NewThing() *Thing { return &Thing{} }

5. PREFER VALUE FOR SMALL IMMUTABLE TYPES
   Small structs (< 64 bytes) are often faster to copy.

6. USE POINTERS FOR OPTIONAL PARAMETERS
   func Process(config *Config)  // nil means use defaults

7. AVOID POINTER TO INTERFACE
   Interfaces already contain a pointer internally.
   *io.Reader is almost never what you want.

8. USE new() FOR ZERO-VALUE ALLOCATION
   p := new(int)  // Clearer than: var x int; p := &x

9. DOCUMENT NIL BEHAVIOR
   // Process handles the data. If data is nil, it uses defaults.
   func Process(data *Data) { ... }

10. AVOID STORING POINTERS TO LOOP VARIABLES
    Always shadow or use Go 1.22+ where this is fixed.
*/
