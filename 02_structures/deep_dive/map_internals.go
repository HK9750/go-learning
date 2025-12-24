package main

import (
	"fmt"
	"unsafe"
)

// =========================================================================
// DEEP DIVE: Map Internals (The hmap struct)
// =========================================================================
// A Go map is a pointer to a 'hmap' struct in the runtime.
// When you pass a map to a function, you are passing this pointer (8 bytes).
// That is why maps effectively act like references.

// Mirroring the runtime 'hmap' struct (simplified for 64-bit architectures)
// Source: src/runtime/map.go
type hmap struct {
	count     int            // Live cells == len(map)
	flags     uint8          // Iterator flags, writing flags
	B         uint8          // log_2 of # of buckets (can hold up to loadFactor * 2^B items)
	noverflow uint16         // Approximate number of overflow buckets
	hash0     uint32         // Hash seed (randomized per map)
	buckets   unsafe.Pointer // Array of 2^B buckets. may be nil if count==0
	oldbuckets unsafe.Pointer // Previous bucket array of half the size, non-nil only when growing
	nevacuate  uintptr       // Progress counter for evacuation (low-order bits used)
	extra      unsafe.Pointer // Optional fields
}

// bmap (Bucket Map) is not exposed directly, but it looks like this conceptually:
// type bmap struct {
//    tophash [8]uint8 // Top 8 bits of hash for each key
//    keys    [8]Key
//    values  [8]Value
//    overflow *bmap   // Overflow bucket
// }

func main() {
	fmt.Println("=== Map Internals Visualization ===")

	// 1. Create a map and populate it
	// Force B to grow: Default load factor is 6.5.
	// We need > 8 items to likely trigger > 1 bucket if B starts small.
	m := make(map[string]int, 10) 
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("k%d", i)
		m[key] = i * 100
	}

	// 2. Inspect the hmap
	// A map variable 'm' is effectively *hmap.
	// We cast the value of 'm' (which is the address of hmap) to *hmap.
	// Note: We use unsafe.Pointer(&m) to get pointer to variable m, then dereference
	//       to get the value stored in m (which is the pointer to hmap).
	//       Wait, no. 'm' IS the pointer. 
	//       But we can't cast 'map[string]int' to '*hmap' directly.
	//       We must cast via unsafe.Pointer.
	
	// Create a generic pointer to the underlying hmap
	pointToHmap := *(*unsafe.Pointer)(unsafe.Pointer(&m))
	hm := (*hmap)(pointToHmap)

	fmt.Printf("\nMap Address (hmap ptr): %p\n", pointToHmap)
	fmt.Printf("Count (len)           : %d\n", hm.count)
	fmt.Printf("B (log2 buckets)      : %d (Values: 2^%d = %d buckets)\n", hm.B, hm.B, 1<<hm.B)
	fmt.Printf("Flags                 : %08b\n", hm.flags) 
	fmt.Printf("Hash0 (Seed)          : %d\n", hm.hash0)
	fmt.Printf("Buckets Ptr           : %p\n", hm.buckets)

	// 3. Visualization
	fmt.Println("\nVisual Representation:")
	fmt.Println("m (var) -> [ hmap * ]")
	fmt.Println("             |")
	fmt.Println("             v")
	fmt.Println("          +----------+")
	fmt.Printf("          | count: %d |\n", hm.count)
	fmt.Printf("          | B    : %d |\n", hm.B)
	fmt.Println("          | buckets  | --+--> [ Bucket 0 ]")
	fmt.Println("          +----------+   |    [ Bucket 1 ]")
	fmt.Println("                         |    [ ...      ]")
	fmt.Println("                         +--> [ Bucket N ]")

	// 4. Growth Trigger
	// Modify map to force growth
	fmt.Println("\n--- Adding more items to trigger potential growth ---")
	for i := 5; i < 20; i++ {
		m[fmt.Sprintf("k%d", i)] = i
	}
	
	fmt.Printf("New Count : %d\n", hm.count)
	fmt.Printf("New B     : %d (Buckets: %d)\n", hm.B, 1<<hm.B)
	
	if hm.oldbuckets != nil {
		fmt.Printf("GROWTH IN PROGRESS! oldbuckets ptr: %p\n", hm.oldbuckets)
		fmt.Println("(The map is gradually moving keys from oldbuckets to buckets)")
	} else {
		fmt.Println("Growth stable (no oldbuckets).")
	}
}
