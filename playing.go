package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	// =========================================================================
	// 1. Arrays are Values (Copy Semantics)
	// =========================================================================
	// In Go, an array is a value. It represents the entire block of memory.
	//
	// VISUALIZATION:
	// 'a' variable in Memory:
	// +---+---+---+
	// | 1 | 2 | 3 |  (Address: 0x1000)
	// +---+---+---+
	a := [3]int{1, 2, 3}

	// When you assign an array to a new variable 'b', Go COPIES the entire memory.
	// 'b' is completely independent of 'a'.
	//
	// 'b' variable in Memory:
	// +---+---+---+
	// | 1 | 2 | 3 |  (Address: 0x2000)
	// +---+---+---+
	b := a

	// Modifying 'b' touches 0x2000, leaving 0x1000 ('a') untouched.
	b[0] = 99
	fmt.Println("array copy ->", a, b) // Output: [1 2 3] [99 2 3]

	// =========================================================================
	// 2. Slices are Descriptors (Reference Semantics)
	// =========================================================================
	// A slice is NOT the array itself. It is a tiny "header" structure describing a window.
	// Struct { Data: *ptr, Len: int, Cap: int }
	//
	// Underlying Array (Anonymous):
	// +---+---+---+
	// | 1 | 2 | 3 | (Address: 0x3000)
	// +---+---+---+
	//
	// 's' Header:
	// [ Ptr: 0x3000 | Len: 3 | Cap: 3 ]
	s := []int{1, 2, 3}

	// Assigning a slice copies only the HEADER (Ptr, Len, Cap).
	// 't' now points to the SAME 0x3000 address.
	//
	// 't' Header:
	// [ Ptr: 0x3000 | Len: 3 | Cap: 3 ]
	t := s

	// Modifying via 't' follows the pointer to 0x3000 and modifies the shared array.
	t[0] = 77
	fmt.Println("slice alias ->", s, t) // Output: [77 2 3] [77 2 3]

	// =========================================================================
	// 3. Append and Dynamic Reallocation
	// =========================================================================
	// Create a slice with Len=3, Cap=3.
	// Underlying: [ 10 | 20 | 30 ]
	s = make([]int, 3, 3)
	s[0], s[1], s[2] = 10, 20, 30

	fmt.Printf("before append ptr=%p len=%d cap=%d\n", unsafe.Pointer(&s[0]), len(s), cap(s))

	// VISUALIZATION BEFORE:
	// s -> [ PtrA | Len:3 | Cap:3 ] --> [ 10 | 20 | 30 ]  (Full!)

	// 's' is full. Appending requires more space.
	// Go performs a "Grow":
	// 1. Allocs NEW array (usually double size).
	// 2. Copies 10, 20, 30 to NEW array.
	// 3. Adds 40.
	// 4. Returns NEW slice header pointing to NEW array.
	s = append(s, 40)

	// VISUALIZATION AFTER:
	// Old Array: [ 10 | 20 | 30 ] (Garbage Collected eventually)
	// s -> [ PtrB | Len:4 | Cap:6 ] --> [ 10 | 20 | 30 | 40 | ? | ? ]
	fmt.Printf("after append ptr=%p len=%d cap=%d\n", unsafe.Pointer(&s[0]), len(s), cap(s))

	// =========================================================================
	// 4. Memory Leaks / Holding References
	// =========================================================================
	// Allocate a giant 1MB array.
	large := make([]byte, 1<<20)

	// 'sub' is a tiny slice (8 bytes), but it points into the GIANT array.
	sub := large[:8]
	fmt.Printf("large backing ptr=%#x sub ptr=%#x len=%d cap=%d\n", uintptr(unsafe.Pointer(&large[0])), uintptr(unsafe.Pointer(&sub[0])), len(sub), cap(sub))

	// Even if we nil 'large', the underlying 1MB array CANNOT be freed
	// because 'sub' is still looking at it.
	large = nil
	fmt.Printf("after large=nil sub still ptr=%#x len=%d cap=%d\n", uintptr(unsafe.Pointer(&sub[0])), len(sub), cap(sub))

	// SOLUTION: Copy what you need to a fresh, small array.
	// append([]byte(nil), ...) creates a new backing array just big enough.
	short := append([]byte(nil), sub...)
	fmt.Printf("copied short ptr=%#x len=%d cap=%d\n", uintptr(unsafe.Pointer(&short[0])), len(short), cap(short))
	// Now the 1MB array has zero references and will be GC'd.

	// =========================================================================
	// 5. Nil vs Empty Slices
	// =========================================================================
	// Uninitialized slice.
	// Header: [ Ptr: 0 | Len: 0 | Cap: 0 ]
	var nilSlice []int

	// Initialized literal with no elements.
	// Header: [ Ptr: 0xZerobase | Len: 0 | Cap: 0 ]
	// It points to a special sentinel address (non-nil), but holds no data.
	emptySlice := []int{}

	fmt.Printf("nilSlice==nil? %t len=%d cap=%d; emptySlice==nil? %t len=%d cap=%d\n", nilSlice == nil, len(nilSlice), cap(nilSlice), emptySlice == nil, len(emptySlice), cap(emptySlice))

	// =========================================================================
	// 6. Full Slice Expressions (Control Capacity)
	// =========================================================================
	u := []int{1, 2, 3, 4, 5}

	// Standard slice: v := u[1:3] would default Cap to (Cap(u) - 1) = 4.
	// v would be able to expand into u's data (3, 4) via append.
	//
	// Full Slice Expression: u[low : high : max]
	// Sets Cap = max - low.
	// Here: Cap = 3 - 1 = 2.
	v := u[1:3:3]

	// v -> [ Ptr: &u[1] | Len: 2 | Cap: 2 ]
	// 'v' is now "capped" strictly. Appending to 'v' requires reallocation,
	// protecting 'u' from accidental overwrites.
	fmt.Printf("u len=%d cap=%d v len=%d cap=%d\n", len(u), cap(u), len(v), cap(v))

	// =========================================================================
	// 7. Internals: Reflecting the Header
	// =========================================================================
	// reflect.SliceHeader mimics the actual C-struct layout of a slice in Go runtime.
	// type SliceHeader struct {
	//     Data uintptr
	//     Len  int
	//     Cap  int
	// }
	var sh reflect.SliceHeader
	_ = sh
	shdr := (*reflect.SliceHeader)(unsafe.Pointer(&u))
	fmt.Printf("slice header Data=%#x Len=%d Cap=%d\n", shdr.Data, shdr.Len, shdr.Cap)
}
