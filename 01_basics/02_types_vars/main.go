// ============================================================================
// GO FUNDAMENTALS: TYPES & VARIABLES
// ============================================================================
// This file provides a comprehensive guide to Go's type system, variable
// declarations, constants, and the critical concept of zero values.
// ============================================================================

package main

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// ============================================================================
// TOPIC: Package-Level Variables
// ============================================================================
// Variables declared outside functions have package scope.
// They are visible to all files in the same package.

var (
	// Package-level variables can use var, NOT :=
	PackageVar = "I am visible throughout the package"

	// Multiple related variables can be grouped
	maxRetries = 3
	timeout    = 30
)

// ============================================================================
// TOPIC: Constants
// ============================================================================
// Constants are evaluated at COMPILE TIME.
// They can only be character, string, boolean, or numeric values.

const (
	// Untyped constants have arbitrary precision until assigned
	Pi       = 3.14159265358979323846264338327950288419716939937510582097494459
	MaxInt64 = 1<<63 - 1 // Bit shifting in constants

	// Typed constants
	TypedPi float64 = 3.14

	// iota - the constant generator
	// iota starts at 0 and increments by 1 for each constant in the block
	Sunday    = iota // 0
	Monday           // 1
	Tuesday          // 2
	Wednesday        // 3
	Thursday         // 4
	Friday           // 5
	Saturday         // 6
)

// Advanced iota patterns
const (
	_  = iota             // Skip 0
	KB = 1 << (10 * iota) // 1 << 10 = 1024
	MB                    // 1 << 20 = 1048576
	GB                    // 1 << 30
	TB                    // 1 << 40
)

// Bit flags with iota
const (
	FlagRead  = 1 << iota // 1 (001)
	FlagWrite             // 2 (010)
	FlagExec              // 4 (100)
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              GO TYPES & VARIABLES                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Variable Declaration Styles
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Variable Declaration Styles")
	fmt.Println("─────────────────────────────────────────")

	// DECLARATION STYLES COMPARISON:
	// ┌──────────────────────────────────────────────────────────────────────┐
	// │ Style                    │ Example           │ Where Allowed        │
	// ├──────────────────────────┼───────────────────┼──────────────────────┤
	// │ var with type            │ var x int = 10    │ Anywhere             │
	// │ var type inference       │ var x = 10        │ Anywhere             │
	// │ var zero value           │ var x int         │ Anywhere             │
	// │ short declaration        │ x := 10           │ Inside functions     │
	// │ multiple vars            │ var a, b int      │ Anywhere             │
	// │ var block                │ var ( ... )       │ Package level        │
	// └──────────────────────────────────────────────────────────────────────┘

	// Style 1: Full declaration with explicit type
	var explicitInt int = 42
	fmt.Printf("Explicit: %d (type: %T)\n", explicitInt, explicitInt)

	// Style 2: Type inference (compiler determines type)
	var inferredInt = 42     // Inferred as int
	var inferredFloat = 3.14 // Inferred as float64
	fmt.Printf("Inferred: %d (type: %T), %f (type: %T)\n",
		inferredInt, inferredInt, inferredFloat, inferredFloat)

	// Style 3: Short declaration (most common inside functions)
	shortInt := 42
	shortString := "hello"
	fmt.Printf("Short: %d, %s\n", shortInt, shortString)

	// Style 4: Multiple variables
	var a, b, c int = 1, 2, 3
	x, y, z := 4.0, "five", true
	fmt.Printf("Multiple: %d, %d, %d | %v, %v, %v\n", a, b, c, x, y, z)

	// ========================================================================
	// SECTION 2: Zero Values
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Zero Values (Critical Concept)")
	fmt.Println("─────────────────────────────────────────")

	// ZERO VALUES TABLE:
	// ┌─────────────────────┬─────────────────────────────────────────────────┐
	// │ Type                │ Zero Value                                      │
	// ├─────────────────────┼─────────────────────────────────────────────────┤
	// │ bool                │ false                                           │
	// │ int, float, complex │ 0 (or 0.0, 0+0i)                                │
	// │ string              │ "" (empty string)                               │
	// │ pointer             │ nil                                             │
	// │ slice               │ nil (but len=0, cap=0, usable with append)      │
	// │ map                 │ nil (CANNOT write to nil map!)                  │
	// │ channel             │ nil (blocks forever on send/receive)            │
	// │ interface           │ nil                                             │
	// │ function            │ nil                                             │
	// │ struct              │ All fields set to their zero values             │
	// │ array               │ All elements set to their zero values           │
	// └─────────────────────┴─────────────────────────────────────────────────┘

	var (
		zeroInt       int
		zeroFloat     float64
		zeroBool      bool
		zeroString    string
		zeroPointer   *int
		zeroSlice     []int
		zeroMap       map[string]int
		zeroChannel   chan int
		zeroInterface interface{}
		zeroFunc      func()
	)

	fmt.Printf("int:       %v (zero: %t)\n", zeroInt, zeroInt == 0)
	fmt.Printf("float64:   %v (zero: %t)\n", zeroFloat, zeroFloat == 0)
	fmt.Printf("bool:      %v (zero: %t)\n", zeroBool, zeroBool == false)
	fmt.Printf("string:    %q (zero: %t)\n", zeroString, zeroString == "")
	fmt.Printf("*int:      %v (nil: %t)\n", zeroPointer, zeroPointer == nil)
	fmt.Printf("[]int:     %v (nil: %t, len: %d)\n", zeroSlice, zeroSlice == nil, len(zeroSlice))
	fmt.Printf("map:       %v (nil: %t)\n", zeroMap, zeroMap == nil)
	fmt.Printf("chan:      %v (nil: %t)\n", zeroChannel, zeroChannel == nil)
	fmt.Printf("interface: %v (nil: %t)\n", zeroInterface, zeroInterface == nil)
	fmt.Printf("func:      %v (nil: %t)\n", zeroFunc, zeroFunc == nil)

	// ⚠️ CRITICAL: nil map vs empty map
	// Writing to a nil map causes a PANIC!
	// zeroMap["key"] = 1  // PANIC: assignment to entry in nil map

	// Safe: Create map first
	safeMap := make(map[string]int)
	safeMap["key"] = 1
	fmt.Printf("Safe map: %v\n", safeMap)

	// ========================================================================
	// SECTION 3: Numeric Types In-Depth
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Numeric Types In-Depth")
	fmt.Println("─────────────────────────────────────────")

	// INTEGER TYPES:
	// ┌─────────────┬─────────┬────────────────────────────────────────────────┐
	// │ Type        │ Size    │ Range                                          │
	// ├─────────────┼─────────┼────────────────────────────────────────────────┤
	// │ int8        │ 8 bits  │ -128 to 127                                    │
	// │ int16       │ 16 bits │ -32,768 to 32,767                              │
	// │ int32       │ 32 bits │ -2.1B to 2.1B                                  │
	// │ int64       │ 64 bits │ -9.2 quintillion to 9.2 quintillion            │
	// │ int         │ varies  │ 32 bits on 32-bit, 64 bits on 64-bit systems   │
	// ├─────────────┼─────────┼────────────────────────────────────────────────┤
	// │ uint8/byte  │ 8 bits  │ 0 to 255                                       │
	// │ uint16      │ 16 bits │ 0 to 65,535                                    │
	// │ uint32      │ 32 bits │ 0 to 4.2B                                      │
	// │ uint64      │ 64 bits │ 0 to 18.4 quintillion                          │
	// │ uint        │ varies  │ Platform dependent                             │
	// │ uintptr     │ varies  │ Large enough to hold pointer bits              │
	// └─────────────┴─────────┴────────────────────────────────────────────────┘

	var (
		i8  int8   = 127
		i16 int16  = 32767
		i32 int32  = 2147483647
		i64 int64  = 9223372036854775807
		u8  uint8  = 255
		u64 uint64 = 18446744073709551615
	)

	fmt.Printf("int8:   %d (size: %d bytes)\n", i8, unsafe.Sizeof(i8))
	fmt.Printf("int16:  %d (size: %d bytes)\n", i16, unsafe.Sizeof(i16))
	fmt.Printf("int32:  %d (size: %d bytes)\n", i32, unsafe.Sizeof(i32))
	fmt.Printf("int64:  %d (size: %d bytes)\n", i64, unsafe.Sizeof(i64))
	fmt.Printf("uint8:  %d (size: %d bytes)\n", u8, unsafe.Sizeof(u8))
	fmt.Printf("uint64: %d (size: %d bytes)\n", u64, unsafe.Sizeof(u64))

	// int vs int64: CRITICAL DISTINCTION
	// They are DIFFERENT TYPES even if same size!
	var platformInt int = 100
	var explicitInt64 int64 = 100
	// platformInt = explicitInt64  // COMPILE ERROR!
	platformInt = int(explicitInt64) // Must cast explicitly
	fmt.Printf("\nint size on this platform: %d bytes\n", unsafe.Sizeof(platformInt))

	// FLOATING POINT TYPES:
	// ┌─────────────┬─────────┬──────────────────────────────────────────────────┐
	// │ Type        │ Size    │ Precision                                        │
	// ├─────────────┼─────────┼──────────────────────────────────────────────────┤
	// │ float32     │ 32 bits │ ~6-7 decimal digits                              │
	// │ float64     │ 64 bits │ ~15-16 decimal digits (DEFAULT for literals)     │
	// └─────────────┴─────────┴──────────────────────────────────────────────────┘

	var f32 float32 = 3.14159265358979323846
	var f64 float64 = 3.14159265358979323846
	fmt.Printf("\nfloat32: %.20f (precision loss!)\n", f32)
	fmt.Printf("float64: %.20f\n", f64)

	// COMPLEX TYPES:
	var c64 complex64 = complex(1, 2)
	var c128 complex128 = complex(1, 2)
	fmt.Printf("\ncomplex64:  %v\n", c64)
	fmt.Printf("complex128: %v\n", c128)
	fmt.Printf("Real: %f, Imag: %f\n", real(c128), imag(c128))

	// ========================================================================
	// SECTION 4: Type Conversions (NO Implicit Casting!)
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Type Conversions")
	fmt.Println("─────────────────────────────────────────")

	// Go requires EXPLICIT type conversion. No implicit casting.
	// This prevents subtle bugs but requires more typing.

	// Integer to Float
	intVal := 42
	floatVal := float64(intVal)
	fmt.Printf("int to float64: %d -> %f\n", intVal, floatVal)

	// Float to Integer (TRUNCATES, doesn't round!)
	pi := 3.99999
	truncated := int(pi)
	fmt.Printf("float to int (TRUNCATES): %f -> %d\n", pi, truncated)

	// Proper rounding
	rounded := int(math.Round(pi))
	fmt.Printf("Proper rounding: %f -> %d\n", pi, rounded)

	// String to bytes and vice versa
	str := "Hello"
	bytes := []byte(str)
	runes := []rune(str)
	fmt.Printf("string -> []byte: %v\n", bytes)
	fmt.Printf("string -> []rune: %v\n", runes)
	fmt.Printf("[]byte -> string: %s\n", string(bytes))

	// ========================================================================
	// SECTION 5: Variable Shadowing (Common Bug Source!)
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Variable Shadowing")
	fmt.Println("─────────────────────────────────────────")

	// Shadowing occurs when an inner scope declares a variable with
	// the same name as an outer scope. The inner variable "shadows" the outer.

	n := 10
	fmt.Printf("Outer n: %d (address: %p)\n", n, &n)

	if n > 5 {
		// SHADOWING: This creates a NEW variable 'n'
		n := 100 // := creates new variable, doesn't assign to outer n
		fmt.Printf("Inner n (shadowed): %d (address: %p)\n", n, &n)
		n++ // Only affects inner n
	}

	fmt.Printf("Outer n after if: %d (unchanged!)\n", n)

	// MEMORY VISUALIZATION:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ OUTER SCOPE                                                            │
	// │ ┌─────────────────────┐                                                │
	// │ │ n: 10 (0xABC)       │ ◄─── This is the original 'n'                  │
	// │ └─────────────────────┘                                                │
	// │     │                                                                  │
	// │     │ if n > 5 {                                                       │
	// │     ▼                                                                  │
	// │     INNER SCOPE                                                        │
	// │     ┌─────────────────────┐                                            │
	// │     │ n: 100 (0xDEF)     │ ◄─── NEW variable, shadows outer 'n'        │
	// │     └─────────────────────┘                                            │
	// │     }                                                                  │
	// │     │                                                                  │
	// │     ▼                                                                  │
	// │ Back to outer scope - original 'n' is still 10!                        │
	// └─────────────────────────────────────────────────────────────────────────┘

	// THE FIX: Use = instead of := when you want to assign to existing variable
	m := 10
	if m > 5 {
		m = 100 // = assigns to existing variable
	}
	fmt.Printf("Correct: m after if: %d\n", m)

	// COMMON SHADOWING BUG WITH ERROR HANDLING:
	// var err error
	// if condition {
	//     result, err := someFunc()  // BUG! 'err' is shadowed
	//     // outer 'err' is still nil
	// }
	// if err != nil { ... }  // This checks outer 'err', always nil!

	// ========================================================================
	// SECTION 6: Type Definitions vs Type Aliases
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Type Definitions vs Type Aliases")
	fmt.Println("─────────────────────────────────────────")

	// TYPE DEFINITION (type X T):
	// Creates a NEW type with same underlying type
	// Used for type safety, adding methods, documentation
	type UserID int
	type Temperature float64

	// TYPE ALIAS (type X = T):
	// Creates another name for the SAME type
	// Used for gradual refactoring, compatibility
	type LegacyID = int

	var uid UserID = 100
	var lid LegacyID = 100
	var plainInt int = 100

	// TYPE COMPARISON:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │                          Type System                                    │
	// │                                                                          │
	// │   ┌─────────────┐         ┌─────────────┐         ┌─────────────┐       │
	// │   │   UserID    │         │   int       │         │  LegacyID   │       │
	// │   │  (new type) │         │ (base type) │         │  (alias)    │       │
	// │   └──────┬──────┘         └──────┬──────┘         └──────┬──────┘       │
	// │          │                       │                       │               │
	// │          │    DIFFERENT TYPES    │      SAME TYPE        │               │
	// │          │ ◄───────────────────► │ ◄───────────────────► │               │
	// │          │    (need conversion)  │    (interchangeable)  │               │
	// └──────────────────────────────────────────────────────────────────────────┘

	// plainInt = uid        // COMPILE ERROR! Different types
	plainInt = int(uid) // OK: explicit conversion
	plainInt = lid      // OK: aliases are the same type

	fmt.Printf("UserID type:   %T (value: %d)\n", uid, uid)
	fmt.Printf("LegacyID type: %T (value: %d)\n", lid, lid)
	fmt.Printf("int type:      %T (value: %d)\n", plainInt, plainInt)

	// Why use type definitions?
	// 1. Type safety: Can't mix up UserID with ProductID
	// 2. Can add methods to your custom type
	// 3. Self-documenting code

	// ========================================================================
	// SECTION 7: Constants Deep Dive
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Constants Deep Dive")
	fmt.Println("─────────────────────────────────────────")

	// Constants have ARBITRARY PRECISION until assigned to a variable
	const hugeNumber = 1e1000 // This is valid! (untyped)
	// var x float64 = hugeNumber  // ERROR: constant overflow

	// You can do constant arithmetic with arbitrary precision
	const (
		big   = 1 << 100
		small = big >> 99
	)
	fmt.Printf("1<<100 >> 99 = %d (fits in int now)\n", small)

	// Untyped vs Typed Constants
	const untypedPi = 3.14159       // Untyped - can be used as float32 or float64
	const typedPi float64 = 3.14159 // Typed - must be float64

	var f32Result float32 = untypedPi // OK: untyped adapts
	// var f32Error float32 = typedPi  // ERROR: cannot use float64 as float32
	fmt.Printf("float32 from untyped: %f\n", f32Result)

	// iota Examples
	fmt.Printf("\nDays: Sun=%d, Mon=%d, Fri=%d\n", Sunday, Monday, Friday)
	fmt.Printf("Sizes: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
	fmt.Printf("Flags: Read=%d, Write=%d, Exec=%d\n", FlagRead, FlagWrite, FlagExec)

	// Combining flags with bitwise OR
	readWrite := FlagRead | FlagWrite
	fmt.Printf("Read+Write: %d (binary: %b)\n", readWrite, readWrite)

	// ========================================================================
	// SECTION 8: Type Information with Reflection
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Type Information")
	fmt.Println("─────────────────────────────────────────")

	values := []interface{}{
		42,
		3.14,
		"hello",
		true,
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		struct{ Name string }{"Go"},
	}

	for _, v := range values {
		t := reflect.TypeOf(v)
		fmt.Printf("Value: %-20v Type: %-20s Kind: %s\n", v, t, t.Kind())
	}

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Types & Variables Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. SHORT DECLARATION OUTSIDE FUNCTION
   ─────────────────────────────────────────────────────────────────────────
   Error: "syntax error: non-declaration statement outside function body"
   Why: := is only valid inside functions
   Fix: Use 'var' at package level

   WRONG:
   x := 10  // at package level

   CORRECT:
   var x = 10

2. UNUSED VARIABLES
   ─────────────────────────────────────────────────────────────────────────
   Error: "x declared but not used"
   Why: Go requires all declared variables to be used
   Fix: Use the variable or assign to blank identifier _

3. MIXED TYPES IN OPERATIONS
   ─────────────────────────────────────────────────────────────────────────
   Error: "invalid operation: x + y (mismatched types int and float64)"
   Why: Go doesn't implicitly convert types
   Fix: Explicitly convert: float64(x) + y

4. WRITING TO NIL MAP
   ─────────────────────────────────────────────────────────────────────────
   Error: "panic: assignment to entry in nil map"
   Why: A nil map has no storage to write to
   Fix: Initialize with make(): m := make(map[string]int)

5. OVERFLOW
   ─────────────────────────────────────────────────────────────────────────
   var x int8 = 127
   x++ // x becomes -128 (overflow wraps around)

   Constants catch overflow at compile time:
   const x int8 = 128  // ERROR: constant 128 overflows int8

6. SHADOWING WITH :=
   ─────────────────────────────────────────────────────────────────────────
   var err error
   if condition {
       result, err := someFunc()  // Creates new 'err', shadows outer
   }
   // outer 'err' is still nil!

   Fix: Declare result separately, use = for err
   var result Type
   result, err = someFunc()
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. PREFER := INSIDE FUNCTIONS
   Use short declaration for local variables; it's more concise.

2. USE EXPLICIT TYPES FOR CLARITY
   When the type isn't obvious from the value:
   var timeout time.Duration = 30 * time.Second

3. GROUP RELATED DECLARATIONS
   var (
       name    string
       age     int
       balance float64
   )

4. USE TYPE DEFINITIONS FOR DOMAIN CONCEPTS
   type UserID int
   type Money int64  // Cents to avoid float issues
   type Temperature float64

5. PREFER int FOR GENERAL INTEGERS
   Use sized types (int32, int64) only when necessary:
   - Interoperating with specific APIs
   - Memory-constrained environments
   - Bit manipulation

6. NEVER COMPARE FLOATS FOR EQUALITY
   WRONG:  if f1 == f2 { ... }
   RIGHT:  if math.Abs(f1-f2) < epsilon { ... }

7. DOCUMENT UNITS IN NAMES OR TYPES
   var timeoutMS int      // milliseconds in name
   var timeout time.Duration  // or use proper type

8. CHECK FOR OVERFLOW IN USER INPUT
   Validate that converted values fit in target type
*/
