package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// DEEP DIVE: Structs & Memory Layout

// 1. Memory Alignment & Padding
// Fields are aligned in memory based on their size.
// The CPU accesses memory in words (e.g., 8 bytes on 64-bit machine).
// Poor ordering of fields creates "Padding" (wasted space).

type BadStruct struct {
	Flag    bool    // 1 byte
	// Padding: 7 bytes wasted here to align Int64!
	Value   int64   // 8 bytes
	Small   int8    // 1 byte
	// Padding: 7 bytes wasted here at the end (struct size must be multiple of largest alignment)
} // Total: 24 bytes

type GoodStruct struct {
	Value   int64   // 8 bytes
	Flag    bool    // 1 byte
	Small   int8    // 1 byte
	// Padding: 6 bytes at end to round up to 16? usually alignment is word size.
	// Actually struct alignment depends on max field alignment (8).
	// So 8 + 1 + 1 = 10 -> padded to 16 (multiple of 8).
} // Total: 16 bytes

// 2. Embedding vs Inheritance (Composition)
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person // Embedded Field (Anonymous)
	ID     int
	Name   string // Shadowing the embedded Name!
}

// 3. Struct Tags (Metadata)
type User struct {
	Username string `json:"user_name" xml:"u_name"`
	Password string `json:"-"` // "-" means ignore in JSON
}

func main() {
	// A. Memory Optimizations checking
	b := BadStruct{}
	g := GoodStruct{}
	fmt.Printf("Size of BadStruct: %d bytes\n", unsafe.Sizeof(b))
	fmt.Printf("Size of GoodStruct: %d bytes (Saved 33%%!)\n", unsafe.Sizeof(g))

	// B. Embedding Promoted Fields
	emp := Employee{
		Person: Person{Name: "Bob", Age: 30},
		ID:     1234,
		Name:   "Alice", // Sets the OUTER Name
	}
	
	fmt.Println("\n--- Embedding ---")
	fmt.Println("Employee Name (Outer):", emp.Name)
	fmt.Println("Person Name (Inner):", emp.Person.Name)
	fmt.Println("Age (Promoted):", emp.Age) // 'Age' promoted from Person

	// C. Empty Struct
	// struct{} takes 0 bytes of storage.
	// Useful for signal channels or hash sets (map[string]struct{}).
	var signal struct{}
	fmt.Printf("\nSize of empty struct: %d bytes\n", unsafe.Sizeof(signal))

	// D. JSON Tags usage
	user := User{Username: "admin", Password: "secret123"}
	js, _ := json.Marshal(user)
	fmt.Println("\n--- JSON Tags ---")
	fmt.Println("JSON:", string(js)) // Check output!
}
