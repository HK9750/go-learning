package main

// ^ 'package main' tells the Go compiler that this file should compile as an executable program
// rather than a shared library.

import (
	"fmt" // 'fmt' stands for format. It implements formatted I/O with functions analogous to C's printf and scanf.
	"os"
)

// DEEP DIVE: The Initialization Order
// 1. Importing packages (recursive).
// 2. Package-level variables are initialized.
// 3. 'init()' functions are run (we'll see these later).
// 4. 'main()' function is executed.

// main is the entry point of the executable.
// It takes no arguments and returns no values.
// To read args, use 'os.Args'. To return status codes, use 'os.Exit()'.
func main() {
	// Standard Output
	fmt.Println("Hello, Deep Dive World!")

	// DEEP DIVE: println vs fmt.Println
	// 'println' is a built-in function that writes to reference standard error (stderr).
	// It is intended for bootstrapping and debugging. usage in production is discouraged.
	// 'fmt.Println' is the standard way to write to stdout.
	println("This is a built-in println (debug only, goes to stderr mostly)")

	// Parsing Arguments
	// os.Args is a slice of strings.
	// os.Args[0] is strictly the name of the program itself.
	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		fmt.Printf("Arguments provided: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided.")
	}

	// Exit Codes
	// Go programs exit with code 0 by default.
	// Uncommenting the line below would force an exit with code 1 (error).
	os.Exit(1)
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Unused Imports:
   Error: "imported and not used: 'fmt'"
   Why: Go prohibits unused dependencies to keep build times fast and code clean.
   Fix: Remove the import or use it (e.g., _ = fmt.Println).

2. The Brace Style:
   Error: "syntax error: unexpected newline, expecting {"
   Why: Go enforces 'K&R style'. The opening brace '{' MUST be on the same line as 'func' or 'if'.
   Example (Wrong):
     func main()
     { ... }

3. Semicolons:
   Error: Rare, but logic errors can occur.
   Why: Go inserts semicolons automatically at end of lines.
   Best Practice: Don't add semicolons manually unless putting multiple statements on one line.
*/

/*
   HOW TO RUN:
   1. go run main.go
      - Compiles and runs in a temporary directory. Fast for dev.

   2. go build
      - Compiles and creates a binary executable in the current directory.
      - ./01_syntax (on Linux/Mac) or 01_syntax.exe (on Windows)
*/
