package main

import (
	"fmt"
	"unsafe"
)

// DEEP DIVE: Unsafe Pointer Arithmetic
// Bypass Go's type safety and memory safety.
// Use cases: Interacting with C code (Cgo), extreme optimizations, serialization internals.
// Risks: Segfaults, GC corruption, non-portable code.

type Data struct {
	Count int64 // 8 bytes
	Flag  bool  // 1 byte
	// Padding: 7 bytes
	Value float64 // 8 bytes
}

func main() {
	d := Data{Count: 100, Flag: true, Value: 99.9}

	fmt.Printf("Original: %+v\n", d)

	// 1. Getting size and alignment
	fmt.Println("Size:", unsafe.Sizeof(d)) // 24 bytes
	fmt.Println("Align:", unsafe.Alignof(d))

	// 2. Pointer Arithmetic
	// Access 'Value' field manually using offsets.

	// Start address of struct
	startPtr := unsafe.Pointer(&d)

	// Offset of 'Value' field
	offsetValue := unsafe.Offsetof(d.Value) // Should be 16
	fmt.Println("Offsetof Value:", offsetValue)

	// Calculate address of 'Value'
	// uintptr is an integer representation of a pointer (allows math).
	// unsafe.Pointer is a void* (allows casting).

	valPtr := (*float64)(unsafe.Pointer(uintptr(startPtr) + offsetValue))

	// Read
	fmt.Println("Read Value via unsafe:", *valPtr)

	// Write
	*valPtr = 123.456
	fmt.Printf("Modified via unsafe: %+v\n", d)

	// 3. String to ByteSlice (Zero Allocation Cast)
	// Typical conversion []byte(str) allocates memory.
	// This unsafe cast reuses the string's backing array correctly.
	// CAUTION: Resulting slice MUST NOT be modified if string is read-only memory.
	str := "Top Secret"

	// StringHeader (Internal structure)
	// struct { Data uintptr; Len int }

	// SliceHeader (Internal structure)
	// struct { Data uintptr; Len int; Cap int }

	nb := stringToBytes(str)
	fmt.Printf("Zero-copy bytes: %v\n", nb)
}

// Unsafe conversion from string to []byte without allocation.
func stringToBytes(s string) []byte {
	// This is the "old way" (pre-Go 1.20).
	// Go 1.20+ provides unsafe.String and unsafe.Slice for safer usage.
	return unsafe.Slice(unsafe.StringData(s), len(s))
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

/*
WARNINGS & SAFETY RULES:

1. The uintptr Trap:
   Error: Storing a pointer as uintptr, then GC runs, then converting back.
   Why: GC treats uintptr as just an integer. It doesn't track it as a reference. If the object moves or is collected, the uintptr becomes invalid/dangling.
   Rule: Convert unsafe.Pointer to uintptr ONLY for immediate arithmetic and back to unsafe.Pointer in the ONE expression.

2. Struct Layout Changes:
   Risk: Assuming fixed offsets.
   Why: Struct padding/layout can change between Go versions or architectures (32 vs 64 bit).
   Fix: Always use unsafe.Offsetof, never hardcoded numbers (e.g., 16).
*/
