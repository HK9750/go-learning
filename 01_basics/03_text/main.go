package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// DEEP DIVE: Strings in Go
// A string is a read-only slice of bytes.
// It is NOT necessarily a sequence of characters.
// This is crucial: strings are immutable. You cannot change s[i].

func main() {
	// 1. Strings vs Bytes vs Runes
	const s = "Hello, 世界" // Contains 2 ASCII chars and 2 Chinese chars (multi-byte)

	fmt.Printf("String: %s\n", s)
	fmt.Printf("Length (len): %d bytes\n", len(s)) 
	// Output: 13 bytes.
	// Breakdown: 
	// "Hello, " (7 chars * 1 byte) = 7
	// "世界"    (2 chars * 3 bytes) = 6
	// Total = 13 bytes.
	// Note: len(s) returns number of BYTES, not CHARACTERS.
	// 世界 -> UTF-8 usually 3 bytes per char. 
	// Let's verify.

	// 2. Iterating strings
	// THE WRONG WAY (if you expect chars):
	fmt.Println("\n--- Byte Iteration (Wrong for unicode) ---")
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x ", s[i]) // Prints hex of each BYTE
	}
	// You will see multi-byte sequences split up.

	// THE RIGHT WAY (Runes):
	fmt.Println("\n\n--- Rune Iteration (Range loop) ---")
	// 'range' on a string decodes UTF-8 automatically!
	for idx, runeValue := range s {
		fmt.Printf("%d: %c (Bytes: %d)\n", idx, runeValue, utf8.RuneLen(runeValue))
	}
	// Notice 'idx' jumps by more than 1 for the Chinese characters.

	// 3. What is a 'rune'?
	// type rune = int32
	// It represents a Unicode Code Point.
	var myRune rune = 'A' // Single quotes for runes/chars
	fmt.Printf("\nRune A: %T, Value: %d\n", myRune, myRune)

	// 4. String Immutability & Modification
	str := "cat"
	// str[0] = 'b' // COMPILER ERROR: cannot assign to str[0]

	// To modify, convert to []byte or []rune
	byteSlice := []byte(str) // Copies memory!
	byteSlice[0] = 'b'
	str2 := string(byteSlice) // Copies memory again!
	fmt.Println("New String:", str2)

	// 5. Efficient String Building
	// Since strings are immutable, "a" + "b" creates a NEW string.
	// Doing this in a loop is O(N^2) memory wise.
	// Use strings.Builder for efficiency.
	fmt.Println("\n--- Efficient Builder ---")
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		sb.WriteString("Go ")
	}
	result := sb.String()
	fmt.Println("Built string:", result)
	// strings.Builder minimizes memory copying.

	// 6. Multiline Strings (Backticks)
	raw := `
This is a "raw" string.
Escapes like \n don't work here.
Great for HTML, JSON, or Regex.
`
	fmt.Println("Raw string:", raw)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Indexing Strings:
   Risk: s[i] accesses the BYTE at index i, not the character.
   Why: In UTF-8, characters (runes) can be 1-4 bytes long.
   Consequence: You might slice a character in half, resulting in valid but garbage bytes aka "REPLACEMENT CHARACTER" .

2. String Immutability:
   Error: "cannot assign to s[0]"
   Why: Strings are read-only.
   Fix: Convert to []byte, modify, convert back to string.

3. Range Loop vs Len Loop:
   Best Practice: Always use 'for i, r := range str' to iterate over characters safely.
   Avoid: 'for i := 0; i < len(str); i++' unless you specifically want to process bytes.
*/
