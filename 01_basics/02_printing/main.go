package main

import (
	"fmt"
)

// TOPIC: fmt Formatting Verbs
// Reference for all commonly used % verbs in Go.

func main() {
	// 1. General Verbs
	// %v   : The value in a default format
	// %+v  : Adds field names (for structs)
	// %#v  : Go-syntax representation (useful for debugging)
	// %T   : Prints the type of the value
	// %%   : Prints a literal percent sign

	fmt.Println("--- General ---")
	str := "Hello"
	num := 42
	fmt.Printf("%%v (value): %v, %v\n", str, num)
	fmt.Printf("%%T (type): %T, %T\n", str, num)

	type User struct {
		Name string
		Age  int
	}
	u := User{"Alice", 30}
	fmt.Printf("Struct default (%%v): %v\n", u)
	fmt.Printf("Struct with fields (%%+v): %+v\n", u)
	fmt.Printf("Go syntax (%%#v): %#v\n", u)

	// 2. Boolean
	// %t   : true or false
	fmt.Println("\n--- Boolean ---")
	yes := true
	fmt.Printf("%%t: %t\n", yes)

	// 3. Integer
	// %d   : Base 10 (decimal)
	// %b   : Base 2 (binary)
	// %o   : Base 8 (octal)
	// %x   : Base 16 (hex, lowercase)
	// %X   : Base 16 (hex, uppercase)
	// %c   : Character represented by the code point
	// %q   : Quoted character literal
	fmt.Println("\n--- Integer ---")
	x := 65 // ASCII for 'A'
	fmt.Printf("Decimal (%%d): %d\n", x)
	fmt.Printf("Binary (%%b): %b\n", x)
	fmt.Printf("Octal (%%o): %o\n", x)
	fmt.Printf("Hex (%%x): %x\n", x)
	fmt.Printf("Character (%%c): %c\n", x)
	fmt.Printf("Quoted char (%%q): %q\n", x)

	// 4. Floating Point
	// %f   : Decimal point, no exponent
	// %.2f : Decimal with precision (e.g. 2 decimal places)
	// %e   : Scientific notation
	// %g   : Compact (uses %f or %e depending on value)
	fmt.Println("\n--- Float ---")
	pi := 3.14159
	fmt.Printf("Default (%%f): %f\n", pi)
	fmt.Printf("Precision 2 (%%.2f): %.2f\n", pi)
	fmt.Printf("Scientific (%%e): %e\n", pi)
	fmt.Printf("Compact (%%g): %g\n", pi)

	// 5. String & Bytes
	// %s   : Basic string
	// %q   : Double-quoted string (safely escaped)
	// %x   : Hex dump of bytes
	fmt.Println("\n--- String ---")
	s := "Go\tLang"
	fmt.Printf("String (%%s): %s\n", s)
	fmt.Printf("Quoted (%%q): %q\n", s) // Shows escapes like \t
	fmt.Printf("Hex (%%x): %x\n", s)
	// 6. Pointer
	// %p   : Base 16 notation with leading 0x
	fmt.Println("\n--- Pointer ---")
	ptr := &x
	fmt.Printf("Pointer (%%p): %p\n", ptr)

	// 7. Width & Padding
	// %5d  : Min width 5, right aligned
	// %-5d : Min width 5, left aligned
	// %05d : Min width 5, zero padded
	fmt.Println("\n--- Width & Padding ---")
	n := 12
	fmt.Printf("Right aligned (%%5d): |%5d|\n", n)
	fmt.Printf("Left aligned (%%-5d): |%-5d|\n", n)
	fmt.Printf("Zero padded (%%05d): |%05d|\n", n)
}
