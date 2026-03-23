// ============================================================================
// GO FUNDAMENTALS: PRINTING & FORMATTING
// ============================================================================
// This file provides a comprehensive guide to Go's fmt package, covering all
// formatting verbs, width/precision specifiers, and I/O functions.
// ============================================================================

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Custom type for demonstrating Stringer interface
type Person struct {
	Name string
	Age  int
	City string
}

// Implementing fmt.Stringer interface
// This controls how the type appears with %v and %s
func (p Person) String() string {
	return fmt.Sprintf("%s (%d) from %s", p.Name, p.Age, p.City)
}

// Custom type implementing fmt.GoStringer for %#v
type Color struct {
	R, G, B uint8
}

func (c Color) GoString() string {
	return fmt.Sprintf("Color{R: 0x%02X, G: 0x%02X, B: 0x%02X}", c.R, c.G, c.B)
}

// Custom type implementing fmt.Formatter for complete control
type HexInt int

func (h HexInt) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v', 'd':
		fmt.Fprintf(f, "%d", int(h))
	case 'x':
		fmt.Fprintf(f, "0x%x", int(h))
	case 'X':
		fmt.Fprintf(f, "0x%X", int(h))
	default:
		fmt.Fprintf(f, "%%!%c(HexInt=%d)", verb, int(h))
	}
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              GO PRINTING & FORMATTING                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Print Functions Overview
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Print Functions Overview")
	fmt.Println("─────────────────────────────────────────")

	// FUNCTION FAMILIES:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Function      │ Output To    │ Newline │ Format │ Returns              │
	// ├───────────────┼──────────────┼─────────┼────────┼──────────────────────┤
	// │ Print         │ stdout       │ No      │ No     │ bytes written, error │
	// │ Println       │ stdout       │ Yes     │ No     │ bytes written, error │
	// │ Printf        │ stdout       │ No      │ Yes    │ bytes written, error │
	// ├───────────────┼──────────────┼─────────┼────────┼──────────────────────┤
	// │ Fprint        │ io.Writer    │ No      │ No     │ bytes written, error │
	// │ Fprintln      │ io.Writer    │ Yes     │ No     │ bytes written, error │
	// │ Fprintf       │ io.Writer    │ No      │ Yes    │ bytes written, error │
	// ├───────────────┼──────────────┼─────────┼────────┼──────────────────────┤
	// │ Sprint        │ string       │ No      │ No     │ string               │
	// │ Sprintln      │ string       │ Yes     │ No     │ string               │
	// │ Sprintf       │ string       │ No      │ Yes    │ string               │
	// ├───────────────┼──────────────┼─────────┼────────┼──────────────────────┤
	// │ Errorf        │ error        │ No      │ Yes    │ error                │
	// └─────────────────────────────────────────────────────────────────────────┘

	// Print - no newline, adds spaces between args only if neither is a string
	fmt.Print("Print: ")
	fmt.Print("Hello", " ", "World")
	fmt.Print("\n")

	// Println - adds newline, always adds spaces between args
	fmt.Println("Println:", "Hello", "World")

	// Printf - formatted output, no automatic newline
	name := "Gopher"
	age := 10
	fmt.Printf("Printf: %s is %d years old\n", name, age)

	// Sprint family - returns string instead of printing
	greeting := fmt.Sprintf("Hello, %s!", name)
	fmt.Println("Sprintf result:", greeting)

	// Fprint family - writes to any io.Writer
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Written to buffer: %d\n", 42)
	fmt.Print("Fprint to buffer: ", buf.String())

	// Writing to stderr
	fmt.Fprintln(os.Stderr, "This goes to stderr (error output)")

	// Errorf - creates formatted error
	err := fmt.Errorf("user %s not found (id: %d)", "Alice", 123)
	fmt.Println("Errorf result:", err)

	// ========================================================================
	// SECTION 2: General Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: General Verbs")
	fmt.Println("─────────────────────────────────────────")

	// GENERAL VERBS:
	// ┌────────┬─────────────────────────────────────────────────────────────┐
	// │ Verb   │ Description                                                 │
	// ├────────┼─────────────────────────────────────────────────────────────┤
	// │ %v     │ Default format (uses type's natural representation)        │
	// │ %+v    │ For structs: includes field names                          │
	// │ %#v    │ Go-syntax representation (useful for debugging)            │
	// │ %T     │ Type of the value                                          │
	// │ %%     │ Literal percent sign                                       │
	// └────────┴─────────────────────────────────────────────────────────────┘

	type Point struct {
		X, Y int
		Name string
	}
	p := Point{10, 20, "origin"}

	fmt.Printf("%%v  : %v\n", p)
	fmt.Printf("%%+v : %+v\n", p)
	fmt.Printf("%%#v : %#v\n", p)
	fmt.Printf("%%T  : %T\n", p)
	fmt.Printf("%%%%  : %%\n")

	// Custom Stringer implementation
	person := Person{"Alice", 30, "NYC"}
	fmt.Printf("Person with Stringer: %v\n", person)

	// Custom GoStringer implementation
	color := Color{255, 128, 0}
	fmt.Printf("Color with GoStringer: %#v\n", color)

	// ========================================================================
	// SECTION 3: Boolean Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Boolean Verbs")
	fmt.Println("─────────────────────────────────────────")

	// %t - the word "true" or "false"
	fmt.Printf("%%t: true=%t, false=%t\n", true, false)

	// ========================================================================
	// SECTION 4: Integer Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Integer Verbs")
	fmt.Println("─────────────────────────────────────────")

	// INTEGER VERBS:
	// ┌────────┬─────────────────────────────────────────────────────────────┐
	// │ Verb   │ Description                                                 │
	// ├────────┼─────────────────────────────────────────────────────────────┤
	// │ %b     │ Binary (base 2)                                             │
	// │ %c     │ Character (Unicode code point)                              │
	// │ %d     │ Decimal (base 10)                                           │
	// │ %o     │ Octal (base 8)                                              │
	// │ %O     │ Octal with 0o prefix                                        │
	// │ %q     │ Quoted character literal (with escapes)                     │
	// │ %x     │ Hexadecimal (base 16, lowercase a-f)                        │
	// │ %X     │ Hexadecimal (base 16, uppercase A-F)                        │
	// │ %U     │ Unicode format: U+0041                                      │
	// └────────┴─────────────────────────────────────────────────────────────┘

	num := 65 // ASCII for 'A'
	negNum := -42

	fmt.Printf("Value: %d\n", num)
	fmt.Printf("%%b (binary):      %b\n", num)
	fmt.Printf("%%c (character):   %c\n", num)
	fmt.Printf("%%d (decimal):     %d\n", num)
	fmt.Printf("%%o (octal):       %o\n", num)
	fmt.Printf("%%O (octal 0o):    %O\n", num)
	fmt.Printf("%%q (quoted char): %q\n", num)
	fmt.Printf("%%x (hex lower):   %x\n", num)
	fmt.Printf("%%X (hex upper):   %X\n", num)
	fmt.Printf("%%U (unicode):     %U\n", num)
	fmt.Printf("Negative %%d:      %d\n", negNum)

	// Using # flag for alternative format
	fmt.Printf("\nWith # flag (alternate format):\n")
	fmt.Printf("%%#b: %#b\n", num)
	fmt.Printf("%%#o: %#o\n", num)
	fmt.Printf("%%#x: %#x\n", num)
	fmt.Printf("%%#X: %#X\n", num)

	// ========================================================================
	// SECTION 5: Floating-Point Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Floating-Point Verbs")
	fmt.Println("─────────────────────────────────────────")

	// FLOAT VERBS:
	// ┌────────┬─────────────────────────────────────────────────────────────┐
	// │ Verb   │ Description                                                 │
	// ├────────┼─────────────────────────────────────────────────────────────┤
	// │ %b     │ Scientific with power of 2 exponent                         │
	// │ %e     │ Scientific notation (lowercase e)                           │
	// │ %E     │ Scientific notation (uppercase E)                           │
	// │ %f     │ Decimal point, no exponent                                  │
	// │ %F     │ Same as %f                                                  │
	// │ %g     │ %e for large exponents, %f otherwise                        │
	// │ %G     │ %E for large exponents, %F otherwise                        │
	// │ %x     │ Hexadecimal notation (lowercase)                            │
	// │ %X     │ Hexadecimal notation (uppercase)                            │
	// └────────┴─────────────────────────────────────────────────────────────┘

	pi := 3.14159265358979323846
	large := 123456789.0
	small := 0.000000123

	fmt.Printf("Pi:\n")
	fmt.Printf("%%f:  %f\n", pi)
	fmt.Printf("%%e:  %e\n", pi)
	fmt.Printf("%%E:  %E\n", pi)
	fmt.Printf("%%g:  %g\n", pi)
	fmt.Printf("%%x:  %x\n", pi)

	fmt.Printf("\nLarge number:\n")
	fmt.Printf("%%f:  %f\n", large)
	fmt.Printf("%%e:  %e\n", large)
	fmt.Printf("%%g:  %g (chooses compact form)\n", large)

	fmt.Printf("\nSmall number:\n")
	fmt.Printf("%%f:  %f\n", small)
	fmt.Printf("%%e:  %e\n", small)
	fmt.Printf("%%g:  %g (chooses compact form)\n", small)

	// ========================================================================
	// SECTION 6: String and Byte Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: String and Byte Verbs")
	fmt.Println("─────────────────────────────────────────")

	// STRING/BYTE VERBS:
	// ┌────────┬─────────────────────────────────────────────────────────────┐
	// │ Verb   │ Description                                                 │
	// ├────────┼─────────────────────────────────────────────────────────────┤
	// │ %s     │ Plain string                                                │
	// │ %q     │ Double-quoted string with escapes                           │
	// │ %x     │ Hex dump (lowercase)                                        │
	// │ %X     │ Hex dump (uppercase)                                        │
	// └────────┴─────────────────────────────────────────────────────────────┘

	str := "Hello\tWorld\n"
	byteSlice := []byte{72, 101, 108, 108, 111} // "Hello"

	fmt.Printf("String:\n")
	fmt.Printf("%%s: %s", str)
	fmt.Printf("%%q: %q\n", str) // Shows escape sequences
	fmt.Printf("%%x: %x\n", str)
	fmt.Printf("%%X: %X\n", str)

	fmt.Printf("\nByte slice:\n")
	fmt.Printf("%%s: %s\n", byteSlice)
	fmt.Printf("%%q: %q\n", byteSlice)
	fmt.Printf("%%x: %x\n", byteSlice)
	fmt.Printf("%% x (with space): % x\n", byteSlice) // Space between bytes

	// ========================================================================
	// SECTION 7: Pointer Verbs
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Pointer Verbs")
	fmt.Println("─────────────────────────────────────────")

	x := 42
	ptr := &x

	fmt.Printf("%%p (pointer): %p\n", ptr)
	fmt.Printf("%%#p (no 0x): %#p\n", ptr) // Without 0x prefix

	// Slices and maps also show address with %p
	slice := []int{1, 2, 3}
	fmt.Printf("Slice %%p: %p\n", slice)

	// ========================================================================
	// SECTION 8: Width and Precision
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Width and Precision")
	fmt.Println("─────────────────────────────────────────")

	// WIDTH AND PRECISION SYNTAX: %[flags][width][.precision]verb
	//
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Syntax        │ Meaning                                               │
	// ├───────────────┼───────────────────────────────────────────────────────┤
	// │ %5d           │ Minimum width 5, right-aligned, space-padded          │
	// │ %-5d          │ Minimum width 5, left-aligned                         │
	// │ %05d          │ Minimum width 5, zero-padded                          │
	// │ %.2f          │ Precision 2 decimal places for float                  │
	// │ %8.2f         │ Width 8, precision 2                                  │
	// │ %*d           │ Width from argument                                   │
	// │ %.*f          │ Precision from argument                               │
	// │ %*.*f         │ Width and precision from arguments                    │
	// └────────────────────────────────────────────────────────────────────────┘

	// Width examples
	fmt.Println("\nWidth (integers):")
	fmt.Printf("|%5d|\n", 42)  // Right-aligned
	fmt.Printf("|%-5d|\n", 42) // Left-aligned
	fmt.Printf("|%05d|\n", 42) // Zero-padded
	fmt.Printf("|%5d|\n", -42) // Negative
	fmt.Printf("|%+5d|\n", 42) // Always show sign
	fmt.Printf("|% 5d|\n", 42) // Space for positive

	// Precision examples
	fmt.Println("\nPrecision (floats):")
	fmt.Printf("|%f|\n", pi)     // Default precision
	fmt.Printf("|%.2f|\n", pi)   // 2 decimal places
	fmt.Printf("|%.0f|\n", pi)   // No decimal places
	fmt.Printf("|%8.2f|\n", pi)  // Width 8, 2 decimals
	fmt.Printf("|%-8.2f|\n", pi) // Left-aligned
	fmt.Printf("|%08.2f|\n", pi) // Zero-padded

	// Precision for strings (truncation)
	fmt.Println("\nPrecision (strings):")
	longStr := "Hello, World!"
	fmt.Printf("|%s|\n", longStr)      // Full string
	fmt.Printf("|%.5s|\n", longStr)    // First 5 chars
	fmt.Printf("|%10.5s|\n", longStr)  // Width 10, max 5 chars
	fmt.Printf("|%-10.5s|\n", longStr) // Left-aligned

	// Dynamic width and precision with *
	fmt.Println("\nDynamic width/precision:")
	width := 10
	precision := 3
	fmt.Printf("|%*d|\n", width, 42)              // Width from variable
	fmt.Printf("|%.*f|\n", precision, pi)         // Precision from variable
	fmt.Printf("|%*.*f|\n", width, precision, pi) // Both from variables

	// ========================================================================
	// SECTION 9: Flags
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Flags")
	fmt.Println("─────────────────────────────────────────")

	// FLAGS:
	// ┌────────┬────────────────────────────────────────────────────────────────┐
	// │ Flag   │ Description                                                    │
	// ├────────┼────────────────────────────────────────────────────────────────┤
	// │ -      │ Left-justify (default is right-justify)                        │
	// │ +      │ Always show sign for numeric values                            │
	// │ (space)│ Leave space for positive sign                                  │
	// │ #      │ Alternate format (0x for hex, etc.)                            │
	// │ 0      │ Pad with zeros instead of spaces                               │
	// └────────┴────────────────────────────────────────────────────────────────┘

	val := 42

	fmt.Printf("No flags:     |%8d|\n", val)
	fmt.Printf("-  (left):    |%-8d|\n", val)
	fmt.Printf("+  (sign):    |%+8d|\n", val)
	fmt.Printf("   (space):   |% 8d|\n", val)
	fmt.Printf("0  (zero):    |%08d|\n", val)
	fmt.Printf("#  (alt hex): |%#8x|\n", val)

	// Combining flags
	fmt.Printf("Combined:     |%+08d|\n", val) // Sign + zero-pad
	fmt.Printf("Combined:     |%-+8d|\n", val) // Left + sign

	// ========================================================================
	// SECTION 10: Argument Indexing
	// ========================================================================
	fmt.Println("\n▶ SECTION 10: Argument Indexing")
	fmt.Println("─────────────────────────────────────────")

	// EXPLICIT ARGUMENT INDICES: %[index]verb
	// Indices are 1-based!

	fmt.Printf("Reordering: %[2]s %[1]s\n", "World", "Hello")
	fmt.Printf("Repeat: %[1]d + %[1]d = %d\n", 5, 10)

	// Useful for internationalization
	fmt.Printf("User %[1]s has %[2]d points. Congrats %[1]s!\n", "Alice", 100)

	// ========================================================================
	// SECTION 11: Error Formatting
	// ========================================================================
	fmt.Println("\n▶ SECTION 11: Error Formatting")
	fmt.Println("─────────────────────────────────────────")

	// fmt.Errorf creates formatted errors
	baseErr := errors.New("connection refused")

	// Simple error
	err1 := fmt.Errorf("failed to connect to server: %v", baseErr)
	fmt.Printf("Simple: %v\n", err1)

	// Error wrapping with %w (Go 1.13+)
	err2 := fmt.Errorf("database operation failed: %w", baseErr)
	fmt.Printf("Wrapped: %v\n", err2)
	fmt.Printf("Unwrapped: %v\n", errors.Unwrap(err2))
	fmt.Printf("Is connection refused? %t\n", errors.Is(err2, baseErr))

	// ========================================================================
	// SECTION 12: Scanning (Input)
	// ========================================================================
	fmt.Println("\n▶ SECTION 12: Scanning (Input)")
	fmt.Println("─────────────────────────────────────────")

	// SCAN FUNCTIONS:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ Function      │ Input From   │ Format │ Returns                        │
	// ├───────────────┼──────────────┼────────┼────────────────────────────────┤
	// │ Scan          │ stdin        │ No     │ items scanned, error           │
	// │ Scanf         │ stdin        │ Yes    │ items scanned, error           │
	// │ Scanln        │ stdin        │ No     │ items scanned, error (newline) │
	// ├───────────────┼──────────────┼────────┼────────────────────────────────┤
	// │ Fscan         │ io.Reader    │ No     │ items scanned, error           │
	// │ Fscanf        │ io.Reader    │ Yes    │ items scanned, error           │
	// │ Fscanln       │ io.Reader    │ No     │ items scanned, error           │
	// ├───────────────┼──────────────┼────────┼────────────────────────────────┤
	// │ Sscan         │ string       │ No     │ items scanned, error           │
	// │ Sscanf        │ string       │ Yes    │ items scanned, error           │
	// │ Sscanln       │ string       │ No     │ items scanned, error           │
	// └─────────────────────────────────────────────────────────────────────────┘

	// Scanning from strings
	input := "Alice 30 175.5"
	var scanName string
	var scanAge int
	var scanHeight float64

	n, scanErr := fmt.Sscan(input, &scanName, &scanAge, &scanHeight)
	if scanErr != nil {
		fmt.Printf("Scan error: %v\n", scanErr)
	} else {
		fmt.Printf("Scanned %d items: name=%s, age=%d, height=%.1f\n",
			n, scanName, scanAge, scanHeight)
	}

	// Scanf with format
	dateInput := "2024-01-15"
	var year, month, day int
	fmt.Sscanf(dateInput, "%d-%d-%d", &year, &month, &day)
	fmt.Printf("Parsed date: year=%d, month=%d, day=%d\n", year, month, day)

	// ========================================================================
	// SECTION 13: Custom Formatters
	// ========================================================================
	fmt.Println("\n▶ SECTION 13: Custom Formatters")
	fmt.Println("─────────────────────────────────────────")

	// Custom Formatter implementation (see HexInt type above)
	h := HexInt(255)
	fmt.Printf("HexInt %%v: %v\n", h)
	fmt.Printf("HexInt %%x: %x\n", h)
	fmt.Printf("HexInt %%X: %X\n", h)
	fmt.Printf("HexInt %%s: %s\n", h)

	// ========================================================================
	// SECTION 14: Performance Considerations
	// ========================================================================
	fmt.Println("\n▶ SECTION 14: Performance Considerations")
	fmt.Println("─────────────────────────────────────────")

	// PERFORMANCE TIPS:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ Situation                        │ Recommendation                        │
	// ├──────────────────────────────────┼───────────────────────────────────────┤
	// │ Building strings in loop         │ Use strings.Builder                   │
	// │ Logging hot paths                │ Use structured logging (slog)         │
	// │ Known format at compile time     │ Consider direct concatenation         │
	// │ Avoiding allocations             │ Use Fprintf to reusable buffer        │
	// │ High-throughput logging          │ Use zerolog/zap instead of fmt        │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Bad: String concatenation in loop (creates many allocations)
	// result := ""
	// for i := 0; i < 1000; i++ {
	//     result += fmt.Sprintf("%d ", i)
	// }

	// Good: Use strings.Builder
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		fmt.Fprintf(&sb, "%d ", i)
	}
	fmt.Println("Builder result:", sb.String())

	// Pre-allocate when size is known
	sb2 := strings.Builder{}
	sb2.Grow(100) // Pre-allocate 100 bytes
	sb2.WriteString("Pre-allocated buffer")
	fmt.Println("Pre-allocated:", sb2.String())

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Printing & Formatting Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. WRONG NUMBER OF ARGUMENTS
   ─────────────────────────────────────────────────────────────────────────
   Warning: "Printf format %s reads arg #2, but call has 1 arg"
   fmt.Printf("Name: %s, Age: %d", name)  // Missing age argument

   Go's vet tool catches these at compile time. Use 'go vet'.

2. WRONG TYPE FOR VERB
   ─────────────────────────────────────────────────────────────────────────
   fmt.Printf("%d", "hello")  // Prints: %!d(string=hello)

   The %!d(...) format indicates type mismatch. Use correct verb.

3. FORGETTING NEWLINE
   ─────────────────────────────────────────────────────────────────────────
   fmt.Printf("Hello")
   fmt.Printf("World")  // Outputs: HelloWorld (no newline between)

   Printf doesn't add newline. Use \n or Println.

4. INFINITE RECURSION IN STRINGER
   ─────────────────────────────────────────────────────────────────────────
   func (p Person) String() string {
       return fmt.Sprintf("Person: %v", p)  // INFINITE RECURSION!
   }

   %v calls String(), which calls %v, which calls String()...
   Fix: Use explicit field access: fmt.Sprintf("Person: %s", p.Name)

5. POINTER VS VALUE IN FORMAT
   ─────────────────────────────────────────────────────────────────────────
   ptr := &value
   fmt.Printf("%v", ptr)   // Prints pointer address
   fmt.Printf("%v", *ptr)  // Prints dereferenced value

6. PRECISION LOSS WITH %g
   ─────────────────────────────────────────────────────────────────────────
   %g removes trailing zeros, which may cause confusion:
   fmt.Printf("%g", 1.500)  // Prints: 1.5 (not 1.500)

   Use %f with explicit precision for consistent output.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. USE go vet
   Always run 'go vet' to catch format string errors at compile time.

2. PREFER Printf FOR FORMATTED OUTPUT
   More readable than string concatenation:
   AVOID: fmt.Println("User " + name + " has " + strconv.Itoa(points) + " points")
   USE:   fmt.Printf("User %s has %d points\n", name, points)

3. IMPLEMENT Stringer FOR CUSTOM TYPES
   Makes debugging easier and fmt output cleaner.

4. USE %q FOR DEBUG OUTPUT
   %q shows string boundaries and escapes:
   fmt.Printf("Input: %q\n", userInput)  // Shows hidden whitespace

5. USE %#v FOR DEEP DEBUGGING
   Shows complete Go syntax representation:
   fmt.Printf("Config: %#v\n", config)

6. USE strings.Builder FOR STRING BUILDING
   Much more efficient than repeated concatenation.

7. LOG STRUCTURED DATA
   For production, prefer structured logging (log/slog in Go 1.21+):
   slog.Info("user logged in", "user", name, "ip", ipAddr)

8. DOCUMENT FORMAT STRINGS
   If format strings are complex, add comments explaining them.
*/
