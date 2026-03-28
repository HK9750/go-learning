// ============================================================================
// GO DATA STRUCTURES: STRUCTS
// ============================================================================
// This file provides a comprehensive guide to Go's struct type,
// including memory layout, embedding, tags, initialization patterns,
// and best practices.
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

// ============================================================================
// STRUCT DEFINITIONS
// ============================================================================

// Basic struct definition
type Person struct {
	Name    string
	Age     int
	Email   string
	private string // unexported field (lowercase)
}

// Struct with various field types
type ComplexStruct struct {
	ID        int
	Name      string
	Tags      []string          // slice field
	Metadata  map[string]string // map field
	CreatedAt time.Time         // struct field
	UpdatedAt *time.Time        // pointer to struct
	Callback  func()            // function field
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    GO STRUCTS                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Struct Basics
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Struct Basics")
	fmt.Println("─────────────────────────────────────────")

	// STRUCT CHARACTERISTICS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Collection of named fields with types                               │
	// │ 2. Value type (copying creates a complete copy)                        │
	// │ 3. Comparable if all fields are comparable                             │
	// │ 4. Zero value has all fields set to their zero values                  │
	// │ 5. Fields can be any type including other structs                      │
	// │ 6. Export controlled by field name capitalization                      │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Initialization methods
	var p1 Person                               // Zero value
	p2 := Person{}                              // Explicit zero value
	p3 := Person{Name: "Alice", Age: 30}        // Named fields (partial)
	p4 := Person{"Bob", 25, "bob@mail.com", ""} // Positional (all fields required)
	p5 := new(Person)                           // Pointer to zero-value struct

	fmt.Printf("p1 (zero):     %+v\n", p1)
	fmt.Printf("p2 (explicit): %+v\n", p2)
	fmt.Printf("p3 (named):    %+v\n", p3)
	fmt.Printf("p4 (positional): %+v\n", p4)
	fmt.Printf("p5 (new):      %+v (pointer: %p)\n", *p5, p5)

	// Field access
	p3.Email = "alice@example.com"
	fmt.Printf("\nAfter setting email: %+v\n", p3)

	// Pointer to struct - automatic dereference
	ptr := &p3
	ptr.Age = 31 // Same as (*ptr).Age = 31
	fmt.Printf("Modified via pointer: %+v\n", p3)

	// STRUCT COMPARISON (if all fields comparable)
	a := Person{Name: "Alice", Age: 30}
	b := Person{Name: "Alice", Age: 30}
	c := Person{Name: "Bob", Age: 30}
	fmt.Printf("\na == b: %t\n", a == b)
	fmt.Printf("a == c: %t\n", a == c)

	// ========================================================================
	// SECTION 2: Memory Layout & Alignment
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Memory Layout & Alignment")
	fmt.Println("─────────────────────────────────────────")

	// MEMORY ALIGNMENT RULES:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Each field aligned to its natural alignment                         │
	// │    - int64/float64/pointers: 8-byte aligned                           │
	// │    - int32/float32: 4-byte aligned                                    │
	// │    - int16: 2-byte aligned                                            │
	// │    - int8/bool/byte: 1-byte aligned                                   │
	// │ 2. Struct size rounded up to alignment of largest field               │
	// │ 3. Padding inserted to maintain alignment                              │
	// └──────────────────────────────────────────────────────────────────────────┘

	// BAD LAYOUT (lots of padding):
	type BadLayout struct {
		a bool  // 1 byte + 7 padding (next field needs 8-byte alignment)
		b int64 // 8 bytes
		c bool  // 1 byte + 7 padding
		d int64 // 8 bytes
	}

	// GOOD LAYOUT (minimal padding):
	type GoodLayout struct {
		b int64 // 8 bytes
		d int64 // 8 bytes
		a bool  // 1 byte
		c bool  // 1 byte + 6 padding (round up to 8)
	}

	// VISUALIZATION:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ BadLayout (32 bytes):                                                  │
	// │ ┌───┬─────────────────┬────────────────┬───┬─────────────────┬────────┐│
	// │ │ a │    PADDING      │       b        │ c │    PADDING      │   d    ││
	// │ │1B │      7B         │      8B        │1B │      7B         │   8B   ││
	// │ └───┴─────────────────┴────────────────┴───┴─────────────────┴────────┘│
	// │                                                                         │
	// │ GoodLayout (24 bytes):                                                 │
	// │ ┌────────────────┬────────────────┬───┬───┬──────────────────┐         │
	// │ │       b        │       d        │ a │ c │     PADDING      │         │
	// │ │      8B        │      8B        │1B │1B │       6B         │         │
	// │ └────────────────┴────────────────┴───┴───┴──────────────────┘         │
	// └─────────────────────────────────────────────────────────────────────────┘

	fmt.Printf("BadLayout size:  %d bytes\n", unsafe.Sizeof(BadLayout{}))
	fmt.Printf("GoodLayout size: %d bytes\n", unsafe.Sizeof(GoodLayout{}))
	fmt.Println("Rule: Order fields from largest to smallest!")

	// Field offsets
	fmt.Println("\nField offsets in GoodLayout:")
	fmt.Printf("  b offset: %d\n", unsafe.Offsetof(GoodLayout{}.b))
	fmt.Printf("  d offset: %d\n", unsafe.Offsetof(GoodLayout{}.d))
	fmt.Printf("  a offset: %d\n", unsafe.Offsetof(GoodLayout{}.a))
	fmt.Printf("  c offset: %d\n", unsafe.Offsetof(GoodLayout{}.c))

	// ========================================================================
	// SECTION 3: Anonymous Structs
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Anonymous Structs")
	fmt.Println("─────────────────────────────────────────")

	// Anonymous structs - defined inline, no type name
	// Useful for one-off structures

	// Inline anonymous struct
	point := struct {
		X, Y int
	}{10, 20}
	fmt.Printf("Anonymous struct: %+v\n", point)

	// Common use: Test cases
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero", 0, 0},
		{"positive", 5, 25},
		{"negative", -3, 9},
	}
	fmt.Printf("Test cases: %+v\n", tests)

	// Common use: JSON parsing of unknown structure
	jsonStr := `{"name": "Alice", "score": 95}`
	var result struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
	json.Unmarshal([]byte(jsonStr), &result)
	fmt.Printf("Parsed JSON: %+v\n", result)

	// ========================================================================
	// SECTION 4: Struct Embedding (Composition)
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Struct Embedding")
	fmt.Println("─────────────────────────────────────────")

	// Go uses COMPOSITION instead of inheritance
	// Embedding provides field and method promotion

	type Address struct {
		Street string
		City   string
		Zip    string
	}

	type Contact struct {
		Phone string
		Email string
	}

	type Customer struct {
		Name    string
		Address // Embedded (anonymous) - fields promoted
		Contact // Embedded (anonymous) - fields promoted
	}

	// EMBEDDING VISUALIZATION:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Customer struct:                                                       │
	// │                                                                         │
	// │ ┌─────────────────────────────────────────────────────────────────┐    │
	// │ │ Name: string                                                     │    │
	// │ ├─────────────────────────────────────────────────────────────────┤    │
	// │ │ Address (embedded):                                              │    │
	// │ │   ├── Street: string  ◄── Promoted to Customer.Street           │    │
	// │ │   ├── City: string    ◄── Promoted to Customer.City             │    │
	// │ │   └── Zip: string     ◄── Promoted to Customer.Zip              │    │
	// │ ├─────────────────────────────────────────────────────────────────┤    │
	// │ │ Contact (embedded):                                              │    │
	// │ │   ├── Phone: string   ◄── Promoted to Customer.Phone            │    │
	// │ │   └── Email: string   ◄── Promoted to Customer.Email            │    │
	// │ └─────────────────────────────────────────────────────────────────┘    │
	// └─────────────────────────────────────────────────────────────────────────┘

	cust := Customer{
		Name: "Alice",
		Address: Address{
			Street: "123 Main St",
			City:   "NYC",
			Zip:    "10001",
		},
		Contact: Contact{
			Phone: "555-1234",
			Email: "alice@example.com",
		},
	}

	// Accessing promoted fields
	fmt.Printf("Name: %s\n", cust.Name)
	fmt.Printf("City (promoted): %s\n", cust.City)   // Same as cust.Address.City
	fmt.Printf("Email (promoted): %s\n", cust.Email) // Same as cust.Contact.Email

	// Can still access via embedded type
	fmt.Printf("Address.City: %s\n", cust.Address.City)

	// FIELD SHADOWING
	type Employee struct {
		Name   string // Shadows embedded Name
		Person        // Embedded
	}

	emp := Employee{
		Name:   "Manager Name", // Employee.Name
		Person: Person{Name: "Person Name", Age: 30},
	}

	fmt.Printf("\nField shadowing:\n")
	fmt.Printf("  emp.Name: %s (Employee's field)\n", emp.Name)
	fmt.Printf("  emp.Person.Name: %s (Person's field)\n", emp.Person.Name)

	// ========================================================================
	// SECTION 5: Struct Tags
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Struct Tags")
	fmt.Println("─────────────────────────────────────────")

	// Tags are metadata attached to struct fields
	// Format: `key1:"value1" key2:"value2"`

	type User struct {
		ID        int       `json:"id" db:"user_id"`
		Username  string    `json:"username" validate:"required,min=3"`
		Email     string    `json:"email" validate:"required,email"`
		Password  string    `json:"-"`                    // "-" means skip in JSON
		CreatedAt time.Time `json:"created_at,omitempty"` // omitempty skips zero values
		Age       int       `json:",string"`              // serialize as string
	}

	// TAG CONVENTIONS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ json       │ JSON encoding/decoding (encoding/json)                     │
	// │ xml        │ XML encoding/decoding (encoding/xml)                       │
	// │ db         │ Database column mapping (various ORM packages)             │
	// │ yaml       │ YAML encoding/decoding                                     │
	// │ validate   │ Validation rules (validator package)                       │
	// │ form       │ HTML form binding                                          │
	// │ mapstructure │ Map to struct conversion                                 │
	// └──────────────────────────────────────────────────────────────────────────┘

	user := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret123",
		Age:      30,
	}

	jsonBytes, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("JSON output:")
	fmt.Println(string(jsonBytes))
	// Note: Password is excluded, Age is a string

	// Reading tags with reflection
	fmt.Println("\nReading tags via reflection:")
	t := reflect.TypeOf(user)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		fmt.Printf("  %s: json=%q\n", field.Name, jsonTag)
	}

	// ========================================================================
	// SECTION 6: Constructor Patterns
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Constructor Patterns")
	fmt.Println("─────────────────────────────────────────")

	// Go doesn't have constructors, but we use factory functions

	// Pattern 1: Simple constructor returning value
	type Point struct {
		X, Y float64
	}

	NewPoint := func(x, y float64) Point {
		return Point{X: x, Y: y}
	}
	p := NewPoint(3.0, 4.0)
	fmt.Printf("NewPoint: %+v\n", p)

	// Pattern 2: Constructor returning pointer
	type Config struct {
		Host    string
		Port    int
		Timeout time.Duration
	}

	NewConfig := func(host string, port int) *Config {
		return &Config{
			Host:    host,
			Port:    port,
			Timeout: 30 * time.Second, // Default value
		}
	}
	cfg := NewConfig("localhost", 8080)
	fmt.Printf("NewConfig: %+v\n", cfg)

	// Pattern 3: Functional options (flexible constructor)
	type Server struct {
		host     string
		port     int
		timeout  time.Duration
		maxConns int
	}

	type ServerOption func(*Server)

	WithPort := func(port int) ServerOption {
		return func(s *Server) { s.port = port }
	}
	WithTimeout := func(d time.Duration) ServerOption {
		return func(s *Server) { s.timeout = d }
	}
	WithMaxConns := func(n int) ServerOption {
		return func(s *Server) { s.maxConns = n }
	}
	WithMaxConns(5)
	NewServer := func(host string, opts ...ServerOption) *Server {
		s := &Server{
			host:     host,
			port:     80,
			timeout:  30 * time.Second,
			maxConns: 100,
		}
		for _, opt := range opts {
			opt(s)
		}
		return s
	}

	srv := NewServer("example.com",
		WithPort(8080),
		WithTimeout(60*time.Second),
	)
	fmt.Printf("Functional options: %+v\n", srv)

	// ========================================================================
	// SECTION 7: Empty Struct
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Empty Struct")
	fmt.Println("─────────────────────────────────────────")

	// struct{} takes ZERO bytes of memory
	// Useful for:
	// 1. Set implementation: map[T]struct{}
	// 2. Signal channels: chan struct{}
	// 3. Method-only types

	fmt.Printf("Size of struct{}: %d bytes\n", unsafe.Sizeof(struct{}{}))

	// Set implementation
	set := make(map[string]struct{})
	set["apple"] = struct{}{}
	set["banana"] = struct{}{}
	fmt.Printf("Set: %v\n", set)
	if _, exists := set["apple"]; exists {
		fmt.Println("'apple' is in set")
	}

	// Signal channel
	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		close(done) // Signal completion
	}()
	<-done // Wait for signal
	fmt.Println("Signal received via empty struct channel")

	// ========================================================================
	// SECTION 8: Struct Copying
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Struct Copying")
	fmt.Println("─────────────────────────────────────────")

	// Structs are VALUE types - copying creates independent copy
	// But beware of fields that are reference types!

	type Data struct {
		ID    int
		Items []string          // Slice - reference type!
		Meta  map[string]string // Map - reference type!
	}

	original := Data{
		ID:    1,
		Items: []string{"a", "b", "c"},
		Meta:  map[string]string{"key": "value"},
	}

	// Shallow copy - reference fields share underlying data
	shallow := original
	shallow.ID = 2                   // Independent
	shallow.Items[0] = "X"           // Affects original!
	shallow.Meta["key"] = "modified" // Affects original!

	fmt.Printf("Original after shallow copy modification:\n")
	fmt.Printf("  ID: %d (unchanged)\n", original.ID)
	fmt.Printf("  Items: %v (CHANGED!)\n", original.Items)
	fmt.Printf("  Meta: %v (CHANGED!)\n", original.Meta)

	// Deep copy - need to manually copy reference types
	deepCopy := func(d Data) Data {
		newItems := make([]string, len(d.Items))
		copy(newItems, d.Items)

		newMeta := make(map[string]string)
		for k, v := range d.Meta {
			newMeta[k] = v
		}

		return Data{
			ID:    d.ID,
			Items: newItems,
			Meta:  newMeta,
		}
	}

	original2 := Data{
		ID:    1,
		Items: []string{"a", "b", "c"},
		Meta:  map[string]string{"key": "value"},
	}
	deep := deepCopy(original2)
	deep.Items[0] = "X"
	deep.Meta["key"] = "modified"

	fmt.Printf("\nOriginal after deep copy modification:\n")
	fmt.Printf("  Items: %v (unchanged!)\n", original2.Items)
	fmt.Printf("  Meta: %v (unchanged!)\n", original2.Meta)

	// ========================================================================
	// SECTION 9: Unexported Fields & Packages
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Unexported Fields")
	fmt.Println("─────────────────────────────────────────")

	// VISIBILITY RULES:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Capitalized names (ID, Name, DoSomething):                              │
	// │   - EXPORTED: accessible from other packages                           │
	// │                                                                          │
	// │ lowercase names (id, name, doSomething):                                │
	// │   - unexported: only accessible within same package                    │
	// └──────────────────────────────────────────────────────────────────────────┘

	type account struct { // unexported type
		balance int // unexported field
	}

	type Account struct { // Exported type
		ID      int // Exported field
		balance int // unexported field - accessible in same package
	}

	// From another package, only Account and ID would be accessible
	// balance would not be accessible

	acc := Account{ID: 1, balance: 1000}
	fmt.Printf("Account: ID=%d, balance=%d (same package access)\n", acc.ID, acc.balance)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Structs Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. COMPARING STRUCTS WITH UNCOMPARABLE FIELDS
   ─────────────────────────────────────────────────────────────────────────
   Error: "invalid operation: a == b (struct containing []string cannot be compared)"

   type Data struct {
       Items []string
   }
   a, b := Data{}, Data{}
   a == b  // COMPILE ERROR!

   Fix: Compare field by field or use reflect.DeepEqual.

2. SHALLOW COPY OF REFERENCE FIELDS
   ─────────────────────────────────────────────────────────────────────────
   copy := original  // Slices and maps share underlying data!

   Fix: Manually deep copy reference type fields.

3. TAKING ADDRESS OF LITERAL FIELD
   ─────────────────────────────────────────────────────────────────────────
   Error: "cannot take the address of Point{X: 1, Y: 2}.X"

   p := &Point{X: 1, Y: 2}.X  // ERROR

   Fix: Store in variable first.
   pt := Point{X: 1, Y: 2}
   p := &pt.X  // OK

4. POSITIONAL INITIALIZATION FRAGILITY
   ─────────────────────────────────────────────────────────────────────────
   p := Person{"Alice", 30, "alice@mail.com", ""}

   If struct changes, this breaks silently.
   Fix: Use named fields: Person{Name: "Alice", Age: 30}

5. FORGETTING TO INITIALIZE EMBEDDED FIELDS
   ─────────────────────────────────────────────────────────────────────────
   emp := Employee{Name: "Bob"}  // Person fields are zero values!

   Fix: Initialize all embedded structs explicitly.

6. JSON UNEXPORTED FIELDS
   ─────────────────────────────────────────────────────────────────────────
   type Data struct {
       name string `json:"name"`  // lowercase - NOT exported, NOT marshaled!
   }

   Fix: Capitalize field names for JSON serialization.

7. MUTEX IN STRUCT BY VALUE
   ─────────────────────────────────────────────────────────────────────────
   type SafeCounter struct {
       mu sync.Mutex
       n  int
   }
   c2 := c1  // COPIES the mutex - deadlock risk!

   Fix: Use pointer receiver or embed *sync.Mutex.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. USE NAMED FIELDS IN LITERALS
   p := Person{Name: "Alice", Age: 30}
   Not: p := Person{"Alice", 30, "", ""}

2. ORDER FIELDS BY SIZE (DESCENDING)
   Minimizes padding, reduces memory usage.

3. USE FUNCTIONAL OPTIONS FOR COMPLEX CONSTRUCTORS
   Flexible, backward-compatible, self-documenting.

4. EMBED FOR "IS-A" RELATIONSHIPS CAREFULLY
   Prefer explicit fields for "HAS-A" relationships.

5. RETURN *T FROM CONSTRUCTORS
   func New() *Config { return &Config{...} }

6. USE EMPTY STRUCT FOR SIGNALS/SETS
   chan struct{}, map[T]struct{}

7. DOCUMENT UNEXPORTED FIELDS
   They're internal but still deserve comments.

8. CONSIDER JSON TAG "-" FOR SENSITIVE DATA
   `json:"-"` excludes from serialization.

9. USE omitempty FOR OPTIONAL FIELDS
   `json:"field,omitempty"`

10. DEEP COPY WHEN SHARING STRUCTS WITH REFERENCE FIELDS
    Prevents unintended data sharing.
*/
