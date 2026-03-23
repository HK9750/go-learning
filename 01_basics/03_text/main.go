// ============================================================================
// GO FUNDAMENTALS: STRINGS, RUNES, AND UNICODE
// ============================================================================
// This file provides a comprehensive guide to Go's string handling, including
// the critical distinction between bytes and runes, Unicode support, and
// efficient string manipulation techniques.
// ============================================================================

package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║          GO STRINGS, RUNES & UNICODE                     ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: String Fundamentals
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: String Fundamentals")
	fmt.Println("─────────────────────────────────────────")

	// CRITICAL CONCEPT: A Go string is a READ-ONLY slice of BYTES.
	// It is NOT a slice of characters!
	//
	// MEMORY LAYOUT:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ String Header (16 bytes on 64-bit)                                     │
	// │ ┌──────────────────┬──────────────────┐                                │
	// │ │ Pointer (8 bytes)│ Length (8 bytes) │                                │
	// │ └────────┬─────────┴──────────────────┘                                │
	// │          │                                                              │
	// │          ▼                                                              │
	// │ ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐                 │
	// │ │ H │ e │ l │ l │ o │ , │   │ 世│ 世│ 世│ 界│ 界│ 界│   (bytes)       │
	// │ │[0]│[1]│[2]│[3]│[4]│[5]│[6]│[7]│[8]│[9]│[10][11][12]│                 │
	// │ └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘                 │
	// │ │←─────── ASCII (1 byte each) ──────→│←─ UTF-8 (3 bytes each) ─→│     │
	// └─────────────────────────────────────────────────────────────────────────┘

	str := "Hello, 世界"

	fmt.Printf("String: %s\n", str)
	fmt.Printf("len(str) = %d (bytes, NOT characters!)\n", len(str))
	fmt.Printf("utf8.RuneCountInString(str) = %d (actual characters)\n",
		utf8.RuneCountInString(str))

	// Breakdown:
	// "Hello, " = 7 bytes (7 ASCII characters × 1 byte)
	// "世界"    = 6 bytes (2 Chinese characters × 3 bytes in UTF-8)
	// Total    = 13 bytes

	// ========================================================================
	// SECTION 2: Bytes vs Runes vs Characters
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Bytes vs Runes vs Characters")
	fmt.Println("─────────────────────────────────────────")

	// TERMINOLOGY:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ Term      │ Go Type   │ Description                                   │
	// ├───────────┼───────────┼───────────────────────────────────────────────┤
	// │ Byte      │ byte/uint8│ 8-bit value (0-255)                           │
	// │ Rune      │ rune/int32│ Unicode code point (U+0000 to U+10FFFF)       │
	// │ Character │ (concept) │ What humans perceive as a character           │
	// │ Grapheme  │ (concept) │ Visual unit (may be multiple runes!)          │
	// └────────────────────────────────────────────────────────────────────────┘

	// A rune is a Unicode code point
	var r rune = '世' // Single quotes for runes
	fmt.Printf("\nRune '世':\n")
	fmt.Printf("  Value: %d\n", r)
	fmt.Printf("  Unicode: U+%04X\n", r)
	fmt.Printf("  UTF-8 bytes: % X\n", []byte(string(r)))
	fmt.Printf("  UTF-8 length: %d bytes\n", utf8.RuneLen(r))

	// ASCII characters are 1 byte = 1 rune
	var asciiRune rune = 'A'
	fmt.Printf("\nRune 'A':\n")
	fmt.Printf("  Value: %d\n", asciiRune)
	fmt.Printf("  Unicode: U+%04X\n", asciiRune)
	fmt.Printf("  UTF-8 bytes: % X\n", []byte(string(asciiRune)))

	// IMPORTANT: Grapheme clusters can span multiple runes!
	// Example: "é" can be represented two ways:
	// 1. Single rune: U+00E9 (Latin Small Letter E with Acute)
	// 2. Two runes: U+0065 (e) + U+0301 (Combining Acute Accent)
	emoji := "👨‍👩‍👧‍👦" // Family emoji - multiple code points!
	fmt.Printf("\nFamily emoji: %s\n", emoji)
	fmt.Printf("  Bytes: %d\n", len(emoji))
	fmt.Printf("  Runes: %d\n", utf8.RuneCountInString(emoji))
	fmt.Printf("  Graphemes: 1 (what you see)\n")

	// ========================================================================
	// SECTION 3: Iterating Over Strings
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Iterating Over Strings")
	fmt.Println("─────────────────────────────────────────")

	testStr := "Go语言"

	// METHOD 1: Byte iteration (WRONG for Unicode!)
	fmt.Println("\n❌ BYTE ITERATION (usually wrong):")
	fmt.Println("   Iterates over raw bytes, splits multi-byte characters")
	for i := 0; i < len(testStr); i++ {
		fmt.Printf("   [%d] byte: 0x%02X\n", i, testStr[i])
	}

	// METHOD 2: Range loop (CORRECT for characters)
	fmt.Println("\n✓ RANGE LOOP (correct for Unicode):")
	fmt.Println("   Automatically decodes UTF-8 into runes")
	for i, r := range testStr {
		fmt.Printf("   [%d] rune: '%c' (U+%04X, %d bytes)\n",
			i, r, r, utf8.RuneLen(r))
	}
	// Notice: index jumps by more than 1 for multi-byte characters!

	// METHOD 3: Convert to []rune (for random access)
	fmt.Println("\n✓ []RUNE CONVERSION (for random access):")
	runes := []rune(testStr)
	fmt.Printf("   Rune slice: %v\n", runes)
	fmt.Printf("   Third character: '%c'\n", runes[2])

	// MEMORY IMPLICATIONS:
	// []byte - same memory as string (copy)
	// []rune - 4 bytes per character (may be larger!)

	// ========================================================================
	// SECTION 4: String Immutability
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: String Immutability")
	fmt.Println("─────────────────────────────────────────")

	// Strings are IMMUTABLE. You cannot change individual bytes.
	original := "hello"
	// original[0] = 'H'  // COMPILE ERROR!

	// To modify, convert to []byte or []rune, modify, convert back
	// This COPIES the data!

	// Modify as bytes (safe only for ASCII)
	bytes := []byte(original)
	bytes[0] = 'H'
	modified := string(bytes)
	fmt.Printf("Original: %s (unchanged)\n", original)
	fmt.Printf("Modified: %s\n", modified)

	// Modify as runes (safe for any Unicode)
	runeSlice := []rune("hello")
	runeSlice[0] = 'H'
	modifiedRune := string(runeSlice)
	fmt.Printf("Modified (runes): %s\n", modifiedRune)

	// MEMORY COST:
	// Converting string ↔ []byte/[]rune always allocates new memory
	// string(bytes) creates new string header + copies data
	// []byte(str) creates new byte slice + copies data

	// ========================================================================
	// SECTION 5: String Literals
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: String Literals")
	fmt.Println("─────────────────────────────────────────")

	// INTERPRETED STRINGS (double quotes)
	// Support escape sequences
	interpreted := "Hello\tWorld\nLine 2"
	fmt.Printf("Interpreted: %q\n", interpreted)
	fmt.Printf("Displayed:\n%s\n", interpreted)

	// ESCAPE SEQUENCES:
	// ┌─────────┬───────────────────────────────────────────────────────────────┐
	// │ Escape  │ Description                                                   │
	// ├─────────┼───────────────────────────────────────────────────────────────┤
	// │ \n      │ Newline                                                       │
	// │ \r      │ Carriage return                                               │
	// │ \t      │ Tab                                                           │
	// │ \\      │ Backslash                                                     │
	// │ \"      │ Double quote                                                  │
	// │ \xhh    │ Hex byte (e.g., \x41 = 'A')                                   │
	// │ \uhhhh  │ Unicode code point (16-bit)                                   │
	// │ \Uhhhhhhhh│ Unicode code point (32-bit)                                 │
	// └─────────┴───────────────────────────────────────────────────────────────┘

	fmt.Printf("Unicode escapes: %s\n", "\u4e16\u754c") // 世界

	// RAW STRINGS (backticks)
	// No escape processing, can span multiple lines
	raw := `This is a raw string.
Newlines are literal.
Escapes like \n don't work.
Great for:
- Regular expressions: ^[a-z]+$
- File paths: C:\Users\Name
- SQL queries
- JSON templates`
	fmt.Printf("\nRaw string:\n%s\n", raw)

	// Note: You cannot include a backtick in a raw string!
	// Workaround: concatenation
	withBacktick := `Hello ` + "`" + `World` + "`"
	fmt.Printf("With backtick: %s\n", withBacktick)

	// ========================================================================
	// SECTION 6: String Operations
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: String Operations")
	fmt.Println("─────────────────────────────────────────")

	// CONCATENATION
	s1 := "Hello"
	s2 := "World"

	// + operator (creates new string)
	concat := s1 + ", " + s2
	fmt.Printf("Concatenation: %s\n", concat)

	// COMPARISON (lexicographic, byte-by-byte)
	fmt.Printf("\"apple\" < \"banana\": %t\n", "apple" < "banana")
	fmt.Printf("\"Go\" == \"Go\": %t\n", "Go" == "Go")

	// INDEXING (returns byte, not rune!)
	fmt.Printf("s1[0] = %c (byte at index 0)\n", s1[0])

	// SLICING (operates on bytes!)
	fmt.Printf("s1[0:3] = %s\n", s1[0:3])

	// ⚠️ WARNING: Slicing can break multi-byte characters!
	chinese := "中文"
	// fmt.Println(chinese[0:2])  // Produces invalid UTF-8!

	// Safe substring with runes
	safeSubstring := string([]rune(chinese)[0:1])
	fmt.Printf("Safe substring: %s\n", safeSubstring)

	// ========================================================================
	// SECTION 7: strings Package (Essential Functions)
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: strings Package")
	fmt.Println("─────────────────────────────────────────")

	sample := "  Hello, World! Hello, Go!  "

	// SEARCHING
	fmt.Printf("Contains 'World': %t\n", strings.Contains(sample, "World"))
	fmt.Printf("HasPrefix 'Hello': %t\n", strings.HasPrefix(strings.TrimSpace(sample), "Hello"))
	fmt.Printf("HasSuffix 'Go!': %t\n", strings.HasSuffix(strings.TrimSpace(sample), "Go!"))
	fmt.Printf("Index of 'World': %d\n", strings.Index(sample, "World"))
	fmt.Printf("Count of 'Hello': %d\n", strings.Count(sample, "Hello"))

	// TRANSFORMATION
	fmt.Printf("ToUpper: %s\n", strings.ToUpper(sample))
	fmt.Printf("ToLower: %s\n", strings.ToLower(sample))
	fmt.Printf("Title: %s\n", strings.Title(sample)) // Deprecated, use cases.Title
	fmt.Printf("TrimSpace: '%s'\n", strings.TrimSpace(sample))
	fmt.Printf("Trim '!': %s\n", strings.Trim(sample, " !"))

	// REPLACEMENT
	fmt.Printf("Replace first: %s\n", strings.Replace(sample, "Hello", "Hi", 1))
	fmt.Printf("Replace all: %s\n", strings.ReplaceAll(sample, "Hello", "Hi"))

	// SPLITTING AND JOINING
	csv := "apple,banana,cherry"
	parts := strings.Split(csv, ",")
	fmt.Printf("Split: %v\n", parts)
	joined := strings.Join(parts, " | ")
	fmt.Printf("Join: %s\n", joined)

	// SPLITTING WITH LIMITS
	limited := strings.SplitN("a:b:c:d", ":", 2)
	fmt.Printf("SplitN(2): %v\n", limited) // ["a", "b:c:d"]

	// FIELDS (split by whitespace)
	sentence := "  one   two   three  "
	words := strings.Fields(sentence)
	fmt.Printf("Fields: %v\n", words) // ["one", "two", "three"]

	// ========================================================================
	// SECTION 8: strings.Builder (Efficient String Building)
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: strings.Builder")
	fmt.Println("─────────────────────────────────────────")

	// PROBLEM: String concatenation in a loop is O(n²) in memory
	// Each concatenation creates a new string!
	//
	// BAD:
	// result := ""
	// for i := 0; i < 1000; i++ {
	//     result += strconv.Itoa(i)  // Creates new string each time!
	// }

	// GOOD: Use strings.Builder
	var builder strings.Builder

	// Pre-allocate if you know approximate size
	builder.Grow(100) // Allocate 100 bytes upfront

	builder.WriteString("Hello")
	builder.WriteString(", ")
	builder.WriteString("World")
	builder.WriteByte('!')
	builder.WriteRune('🎉')

	result := builder.String()
	fmt.Printf("Builder result: %s\n", result)
	fmt.Printf("Builder length: %d\n", builder.Len())
	fmt.Printf("Builder capacity: %d\n", builder.Cap())

	// Reset for reuse
	builder.Reset()
	builder.WriteString("Reused!")
	fmt.Printf("After reset: %s\n", builder.String())

	// MEMORY LAYOUT:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ strings.Builder                                                        │
	// │ ┌─────────────────────────────────────────────────────────────┐        │
	// │ │ Internal byte slice (grows as needed)                       │        │
	// │ │ ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───────────────┐  │        │
	// │ │ │ H │ e │ l │ l │ o │ , │   │ W │...│ ! │   (capacity)   │  │        │
	// │ │ └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───────────────┘  │        │
	// │ │ |←──────── len ──────────→|←─── available capacity ────→|  │        │
	// │ └─────────────────────────────────────────────────────────────┘        │
	// │                                                                        │
	// │ .String() returns string pointing to same underlying array (no copy!)  │
	// │ ⚠️ Do NOT use Builder after calling String() (undefined behavior)      │
	// └─────────────────────────────────────────────────────────────────────────┘

	// ========================================================================
	// SECTION 9: Unicode Package
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Unicode Package")
	fmt.Println("─────────────────────────────────────────")

	testRunes := []rune{'A', 'a', '1', '中', ' ', '\n', '!', 'α'}

	fmt.Println("Character classification:")
	for _, r := range testRunes {
		fmt.Printf("  '%c': ", r)
		props := []string{}
		if unicode.IsLetter(r) {
			props = append(props, "Letter")
		}
		if unicode.IsDigit(r) {
			props = append(props, "Digit")
		}
		if unicode.IsUpper(r) {
			props = append(props, "Upper")
		}
		if unicode.IsLower(r) {
			props = append(props, "Lower")
		}
		if unicode.IsSpace(r) {
			props = append(props, "Space")
		}
		if unicode.IsPunct(r) {
			props = append(props, "Punct")
		}
		if unicode.IsControl(r) {
			props = append(props, "Control")
		}
		if len(props) == 0 {
			props = append(props, "Other")
		}
		fmt.Println(strings.Join(props, ", "))
	}

	// Case conversion
	fmt.Printf("\nCase conversion:\n")
	fmt.Printf("  ToUpper('a') = '%c'\n", unicode.ToUpper('a'))
	fmt.Printf("  ToLower('A') = '%c'\n", unicode.ToLower('A'))
	fmt.Printf("  ToTitle('a') = '%c'\n", unicode.ToTitle('a'))

	// ========================================================================
	// SECTION 10: utf8 Package
	// ========================================================================
	fmt.Println("\n▶ SECTION 10: utf8 Package")
	fmt.Println("─────────────────────────────────────────")

	testString := "Hello, 世界!"

	// Validation
	fmt.Printf("Valid UTF-8: %t\n", utf8.ValidString(testString))

	// Counting
	fmt.Printf("Byte count: %d\n", len(testString))
	fmt.Printf("Rune count: %d\n", utf8.RuneCountInString(testString))

	// Decoding runes manually
	fmt.Println("\nManual UTF-8 decoding:")
	remaining := testString
	for len(remaining) > 0 {
		r, size := utf8.DecodeRuneInString(remaining)
		fmt.Printf("  '%c' (U+%04X) takes %d byte(s)\n", r, r, size)
		remaining = remaining[size:]
	}

	// UTF-8 encoding sizes
	fmt.Println("\nUTF-8 encoding sizes:")
	fmt.Printf("  ASCII (U+0000-U+007F): 1 byte\n")
	fmt.Printf("  Latin/Greek/etc (U+0080-U+07FF): 2 bytes\n")
	fmt.Printf("  CJK/Common (U+0800-U+FFFF): 3 bytes\n")
	fmt.Printf("  Emoji/Rare (U+10000-U+10FFFF): 4 bytes\n")

	// Example sizes
	examples := []rune{'A', 'é', '中', '🚀'}
	for _, r := range examples {
		fmt.Printf("  '%c' (U+%04X): %d byte(s)\n", r, r, utf8.RuneLen(r))
	}

	// ========================================================================
	// SECTION 11: Common String Tasks
	// ========================================================================
	fmt.Println("\n▶ SECTION 11: Common String Tasks")
	fmt.Println("─────────────────────────────────────────")

	// Reverse a string (rune-aware)
	original2 := "Hello, 世界!"
	reversed := reverseString(original2)
	fmt.Printf("Reverse: '%s' -> '%s'\n", original2, reversed)

	// Check palindrome (rune-aware)
	palindrome := "上海自来水来自海上"
	fmt.Printf("Palindrome '%s': %t\n", palindrome, isPalindrome(palindrome))

	// Truncate with ellipsis (rune-aware)
	longText := "This is a very long text that needs truncation"
	truncated := truncateString(longText, 20)
	fmt.Printf("Truncate: '%s'\n", truncated)

	// Count words
	text := "  Hello   World  "
	wordCount := len(strings.Fields(text))
	fmt.Printf("Word count in '%s': %d\n", text, wordCount)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Strings, Runes & Unicode Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// reverseString reverses a string in a Unicode-safe manner
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// isPalindrome checks if a string is a palindrome (Unicode-aware)
func isPalindrome(s string) bool {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}

// truncateString truncates a string to maxRunes characters with ellipsis
func truncateString(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	if maxRunes <= 3 {
		return "..."
	}
	return string(runes[:maxRunes-3]) + "..."
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. INDEXING STRINGS BY BYTE INSTEAD OF RUNE
   ─────────────────────────────────────────────────────────────────────────
   WRONG:
   s := "世界"
   firstChar := s[0]  // This is a BYTE, not '世'!

   RIGHT:
   firstChar := []rune(s)[0]  // Now it's '世'

   Or use range loop for iteration.

2. SLICING MULTI-BYTE CHARACTERS
   ─────────────────────────────────────────────────────────────────────────
   WRONG:
   s := "Hello, 世界"
   sub := s[0:8]  // May slice '世' in half, producing invalid UTF-8

   RIGHT:
   runes := []rune(s)
   sub := string(runes[0:8])  // Safe substring

3. len() RETURNS BYTES, NOT CHARACTERS
   ─────────────────────────────────────────────────────────────────────────
   s := "世界"
   len(s)                     // Returns 6 (bytes)
   utf8.RuneCountInString(s)  // Returns 2 (characters)

4. COMPARING STRINGS WITH DIFFERENT NORMALIZATIONS
   ─────────────────────────────────────────────────────────────────────────
   s1 := "é"        // Single rune: U+00E9
   s2 := "é"        // Two runes: U+0065 + U+0301
   s1 == s2         // May be FALSE even though they look the same!

   Solution: Use golang.org/x/text/unicode/norm for normalization.

5. MODIFYING STRING IN PLACE
   ─────────────────────────────────────────────────────────────────────────
   Error: "cannot assign to s[0]"
   Strings are immutable. Convert to []byte or []rune first.

6. INEFFICIENT STRING CONCATENATION
   ─────────────────────────────────────────────────────────────────────────
   WRONG (O(n²) memory):
   result := ""
   for _, s := range strings { result += s }

   RIGHT:
   var builder strings.Builder
   for _, s := range strings { builder.WriteString(s) }
   result := builder.String()
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. USE range LOOP FOR STRING ITERATION
   Automatically decodes UTF-8 and handles multi-byte characters.

2. USE strings.Builder FOR CONCATENATION
   Especially in loops. Pre-allocate with Grow() if size is known.

3. USE utf8.RuneCountInString FOR LENGTH
   When you need character count, not byte count.

4. USE []rune FOR RANDOM ACCESS
   When you need to access characters by index or modify string.

5. PREFER RAW STRINGS FOR SPECIAL CHARACTERS
   Regular expressions, file paths, SQL queries.

6. VALIDATE USER INPUT
   Use utf8.ValidString() for user-provided strings.

7. BE AWARE OF GRAPHEME CLUSTERS
   Some visual "characters" are multiple runes (emojis, combining marks).
   For proper text segmentation, use golang.org/x/text/unicode/segmentation.

8. NORMALIZE UNICODE FOR COMPARISON
   If comparing user input, normalize first using golang.org/x/text/unicode/norm.

9. USE strings PACKAGE FUNCTIONS
   They're optimized and handle edge cases properly.

10. DOCUMENT ENCODING EXPECTATIONS
    Be explicit about whether strings are UTF-8, ASCII, or other encoding.
*/
