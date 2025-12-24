package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// DEEP DIVE: Arrays vs Slices
// Array: Fixed size, value type (passing it copies the whole array!).
// Slice: Dynamic window onto an underlying array. Reference-like light struct.

// The Slice Header (Internally) looks like this:
// type SliceHeader struct {
//     Data uintptr // Pointer to the underlying array
//     Len  int     // Number of elements in the slice
//     Cap  int     // Capacity (elements allocated in underlying array starting from Data)
// }
//
// VISUALIZATION:
//
//    Variable 'slice' (Stack)          Underlying Array (Heap/Stack)
//   +-----------------+               +----+----+----+----+----+
//   | Ptr: 0xAddr     | ------------> | 10 | 20 | 30 | .. | .. |
//   | Len: 3          |               +----+----+----+----+----+
//   | Cap: 5          |               ^              ^
//   +-----------------+               |              |
//                                    Head           Cap Limit

func main() {
	// 1. Array (Value Type)
	var arr [3]int = [3]int{10, 20, 30}
	arrCopy := arr // COPIES the *entire* array.
	arrCopy[0] = 99
	fmt.Printf("Array Original: %v, Copy: %v (Different!)\n", arr, arrCopy)

	// 2. Slice (Reference-like)
	slice := []int{1, 2, 3}
	sliceRef := slice // Copies only the Header (ptr, len, cap)
	sliceRef[0] = 999
	fmt.Printf("Slice Original: %v, Ref: %v (Same backing array!)\n", slice, sliceRef)

	// 3. Length vs Capacity & Append
	// When you append to a slice, if len < cap, it just increases len.
	// If len == cap, Go allocates a NEW BIGGER array, copies data, and updates the pointer.
	// Growth rule: Generally doubles for small slices, grows by ~1.25x for large ones.
	
	fmt.Println("\n--- Append Mechanic ---")
	var s []int
	fmt.Printf("Init: Len:%d Cap:%d Ptr:%p\n", len(s), cap(s), s)
	
	for i := 0; i < 5; i++ {
		s = append(s, i) // Might reallocate underlying array
		fmt.Printf("Append %d: Len:%d Cap:%d Ptr:%p\n", i, len(s), cap(s), s)
	}
	// Notice 'Ptr' changes when Capacity jumps (e.g. 1 -> 2 -> 4 -> 8).

	// 4. Slicing Slices (The "Window")
	// Creating a sub-slice shares the same underlying array!
	// Be careful with memory leaks here.
	original := make([]int, 1000000) // 8MB array
	smallWindow := original[:2]      // Keeps the WHOLE 8MB array alive in memory!

	// Fix for memory leak: Copy what you need
	realSmall := make([]int, 2)
	copy(realSmall, smallWindow)
	// 'original' can now be GC'd if no other refs exist.

	fmt.Println("\n--- Slice Header Inspection ---")
	inspectSlice(smallWindow)

	// 5. Preallocation (Performance)
	// Always preallocate if you know the size to avoid reallocations.
	prealloc := make([]int, 0, 100) // len=0, cap=100
	fmt.Printf("Preallocated: Len:%d Cap:%d\n", len(prealloc), cap(prealloc))
}

func inspectSlice(s []int) {
	// Using unsafe to convert slice to its Header representation
	header := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	fmt.Printf("Slice Header -> Data: 0x%x, Len: %d, Cap: %d\n", header.Data, header.Len, header.Cap)
}
