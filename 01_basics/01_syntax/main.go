// ============================================================================
// GO FUNDAMENTALS: SYNTAX & PROGRAM STRUCTURE
// ============================================================================
// This file covers the foundational syntax and structure of Go programs.
// Understanding these concepts is essential for every Go developer.
// ============================================================================

package main

// TOPIC: Package Declaration
// ============================================================================
// Every Go file MUST start with a package declaration.
//
// 'package main' tells the Go compiler this file should compile as an
// executable program rather than a shared library.
//
// IMPORTANT RULES:
// 1. The package name must match the directory name (except for 'main')
// 2. All files in the same directory must belong to the same package
// 3. Only 'package main' can have a main() function (entry point)
//
// PACKAGE TYPES:
// ┌─────────────────┬────────────────────────────────────────────────┐
// │ package main    │ Executable program - must have main() func    │
// │ package foo     │ Library package - can be imported by others   │
// │ package foo_test│ Test package - for black-box testing          │
// └─────────────────┴────────────────────────────────────────────────┘

import (
	"fmt"     // Format package: formatted I/O (Printf, Println, etc.)
	"os"      // Operating system functionality (args, env, files)
	"runtime" // Runtime information (Go version, architecture)
	"time"    // Time and duration utilities
)

// TOPIC: Import Declarations
// ============================================================================
// Imports bring external packages into your code.
//
// IMPORT STYLES:
// ┌────────────────────────────────────────────────────────────────────────┐
// │ import "fmt"                    // Single import                       │
// │ import (                        // Grouped imports (preferred)         │
// │     "fmt"                                                              │
// │     "os"                                                               │
// │ )                                                                      │
// │ import . "fmt"                  // Dot import (use Println directly)   │
// │ import _ "net/http/pprof"       // Blank import (side effects only)    │
// │ import myio "io"                // Aliased import (myio.Reader)        │
// └────────────────────────────────────────────────────────────────────────┘
//
// IMPORT PATHS:
// - Standard library: "fmt", "os", "net/http"
// - External: "github.com/user/repo/package"
// - Internal: "./internal/mypackage" or relative to module root
//
// CRITICAL: Go prohibits unused imports. Your code won't compile if you
// import something and don't use it. This keeps builds fast and code clean.

// ============================================================================
// DEEP DIVE: The Initialization Order
// ============================================================================
// Go programs have a strict, deterministic initialization order:
//
// INITIALIZATION SEQUENCE:
// ┌─────────────────────────────────────────────────────────────────────────┐
// │  1. IMPORTS          2. PACKAGE VARS      3. init()        4. main()   │
// │  ┌─────────┐         ┌─────────────┐      ┌─────────┐      ┌─────────┐ │
// │  │ Imports │ ──────► │ var x = ... │ ───► │ init()  │ ───► │ main()  │ │
// │  │ (recur- │         │ const y = ..│      │ init()  │      │ (entry) │ │
// │  │  sive)  │         │ (in order)  │      │ (multi) │      │         │ │
// │  └─────────┘         └─────────────┘      └─────────┘      └─────────┘ │
// └─────────────────────────────────────────────────────────────────────────┘
//
// KEY POINTS:
// - Imported packages initialize FIRST (recursively, depth-first)
// - Package-level variables initialize in declaration order
// - Multiple init() functions are allowed (per file AND across files)
// - init() functions run in the order they appear in source
// - main() is called LAST, only after ALL init() complete
//
// DEPENDENCY GRAPH EXAMPLE:
//
//     main imports A, B
//     A imports C
//     B imports C
//
//     Initialization order: C → A → B → main
//     (C only initialized once, even though both A and B import it)

// Package-level variables - initialized before any init() or main()
var (
	programStart = time.Now()        // Initialized at program start
	goVersion    = runtime.Version() // Go runtime version
	numCPU       = runtime.NumCPU()  // Available CPU cores
	goos         = runtime.GOOS      // Operating system
	goarch       = runtime.GOARCH    // Architecture
)

// TOPIC: init() Function
// ============================================================================
// init() is a special function that:
// - Takes no arguments and returns nothing
// - Cannot be called directly
// - Runs automatically before main()
// - Can appear multiple times (even in same file)
// - Runs in the order defined
//
// USE CASES:
// - Verify/repair program state
// - Register with other packages (e.g., database drivers)
// - Compute complex package-level variables
// - Initialize things that can't be done with simple var declarations
//
// WARNING: Avoid heavy work in init(). It makes programs slow to start
// and makes testing difficult (init always runs).

func init() {
	// This runs BEFORE main()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  INIT #1: Program initializing...")
	fmt.Printf("  Go Version: %s | OS: %s | Arch: %s | CPUs: %d\n",
		goVersion, goos, goarch, numCPU)
	fmt.Println("═══════════════════════════════════════════════════════════")
}

func init() {
	// Yes, you can have multiple init() functions!
	// They run in the order they appear in the file.
	fmt.Println("  INIT #2: Second init function running...")
}

// ============================================================================
// main() - The Entry Point
// ============================================================================
// main() is the entry point for executable programs.
//
// RULES:
// - Must be in package main
// - Takes no arguments (use os.Args for command line args)
// - Returns nothing (use os.Exit() for exit codes)
// - When main() returns, the program exits immediately
//   (all goroutines are killed without cleanup!)

func main() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║           GO SYNTAX & PROGRAM STRUCTURE                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Basic Output
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Basic Output Functions")
	fmt.Println("─────────────────────────────────────")

	// fmt.Println() - Print with newline
	fmt.Println("Hello, World!")

	// fmt.Printf() - Formatted print (no automatic newline)
	name := "Gopher"
	age := 10
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	// fmt.Print() - Print without newline
	fmt.Print("This ")
	fmt.Print("is ")
	fmt.Print("continuous.\n")

	// DEEP DIVE: println vs fmt.Println
	// ────────────────────────────────────────────────────────────────────────
	// println() is a built-in function that:
	// - Writes to STDERR (not stdout!)
	// - Is NOT guaranteed to remain in the language
	// - Has inconsistent spacing behavior
	// - Is intended ONLY for bootstrapping and debugging
	//
	// NEVER use println() in production code. Always use fmt.Println().
	// ────────────────────────────────────────────────────────────────────────
	println("[Built-in println - goes to stderr, debug only]")

	// ========================================================================
	// SECTION 2: Command Line Arguments
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Command Line Arguments")
	fmt.Println("─────────────────────────────────────")

	// os.Args is a slice of strings containing command line arguments
	//
	// MEMORY LAYOUT:
	// ┌──────────────────────────────────────────────────────────────┐
	// │ os.Args = []string                                          │
	// │ ┌─────────────┬────────────┬────────────┬────────────┐      │
	// │ │ [0] program │ [1] arg1   │ [2] arg2   │ [3] arg3   │ ...  │
	// │ │ name/path   │ (first arg)│ (second)   │ (third)    │      │
	// │ └─────────────┴────────────┴────────────┴────────────┘      │
	// └──────────────────────────────────────────────────────────────┘
	//
	// Run with: go run main.go arg1 arg2 arg3

	fmt.Printf("Program name: %s\n", os.Args[0])
	fmt.Printf("Number of arguments: %d\n", len(os.Args))

	if len(os.Args) > 1 {
		fmt.Printf("Arguments provided: %v\n", os.Args[1:])
	} else {
		fmt.Println("No arguments provided. Try: go run main.go hello world")
	}

	// For complex argument parsing, use the 'flag' package:
	// import "flag"
	// verbose := flag.Bool("v", false, "verbose output")
	// flag.Parse()

	// ========================================================================
	// SECTION 3: Exit Codes
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Exit Codes")
	fmt.Println("─────────────────────────────────────")

	// EXIT CODE CONVENTIONS:
	// ┌─────────┬────────────────────────────────────────────────────────┐
	// │ Code    │ Meaning                                                │
	// ├─────────┼────────────────────────────────────────────────────────┤
	// │ 0       │ Success (default when main() returns normally)        │
	// │ 1       │ General error                                         │
	// │ 2       │ Misuse of command (e.g., invalid arguments)           │
	// │ 126     │ Command not executable                                │
	// │ 127     │ Command not found                                     │
	// │ 128+N   │ Fatal signal N (e.g., 130 = SIGINT, 137 = SIGKILL)   │
	// └─────────┴────────────────────────────────────────────────────────┘
	//
	// IMPORTANT: os.Exit() is IMMEDIATE
	// - defer statements do NOT run
	// - Cleanup code does NOT execute
	// - Goroutines are killed without warning
	//
	// Best practice: Return errors up to main(), call os.Exit() only there.

	fmt.Println("Default exit code is 0 (success)")
	fmt.Println("Use os.Exit(1) for errors")
	fmt.Println("⚠️  os.Exit() skips all defers!")

	// Demonstrating that defer won't run with os.Exit()
	// (Commented out to let the rest of main() run)
	/*
		defer fmt.Println("This will NOT print if os.Exit() is called!")
		os.Exit(1) // Uncomment to test - defer above won't run
	*/

	// ========================================================================
	// SECTION 4: Code Organization Best Practices
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Code Organization")
	fmt.Println("─────────────────────────────────────")

	// GO FILE STRUCTURE (Recommended Order):
	// ┌────────────────────────────────────────────────────────────────────┐
	// │ 1. Package comment (for godoc)                                    │
	// │ 2. Package declaration                                            │
	// │ 3. Imports                                                        │
	// │ 4. Constants                                                      │
	// │ 5. Package-level variables                                        │
	// │ 6. Type definitions                                               │
	// │ 7. init() functions                                               │
	// │ 8. main() or exported functions                                   │
	// │ 9. Private functions                                              │
	// │ 10. Methods (grouped by type)                                     │
	// └────────────────────────────────────────────────────────────────────┘

	fmt.Println("✓ Package clause first")
	fmt.Println("✓ Imports grouped and organized")
	fmt.Println("✓ Constants before variables")
	fmt.Println("✓ Types before functions")
	fmt.Println("✓ Exported before unexported")

	// ========================================================================
	// SECTION 5: Semicolons & Brace Style
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Semicolons & Brace Style")
	fmt.Println("─────────────────────────────────────")

	// Go uses AUTOMATIC SEMICOLON INSERTION
	// The lexer automatically inserts semicolons at end of lines when:
	// - The line ends with an identifier, literal, or certain keywords
	// - The line ends with ), ], or }
	//
	// This is why the opening brace MUST be on the same line:
	//
	// CORRECT:
	// if true {        // semicolon NOT inserted after 'true'
	//     ...
	// }
	//
	// WRONG:
	// if true          // semicolon INSERTED here!
	// {                // This becomes a separate statement - syntax error
	//     ...
	// }

	fmt.Println("Go uses K&R brace style (opening brace on same line)")
	fmt.Println("This is ENFORCED by the language, not just convention")

	// Semicolons are legal but rarely used:
	x := 1
	y := 2
	fmt.Printf("Multiple statements: x=%d, y=%d\n", x, y)

	// ========================================================================
	// SECTION 6: Running Go Programs
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Running Go Programs")
	fmt.Println("─────────────────────────────────────")

	// EXECUTION METHODS:
	// ┌─────────────────────────────────────────────────────────────────────┐
	// │ Method              │ Command                │ Use Case            │
	// ├─────────────────────┼────────────────────────┼─────────────────────┤
	// │ Run directly        │ go run main.go         │ Development         │
	// │ Build then run      │ go build && ./program  │ Testing binary      │
	// │ Install globally    │ go install             │ CLI tools           │
	// │ Run tests           │ go test                │ Testing             │
	// │ Run with race check │ go run -race main.go   │ Concurrency debug   │
	// └─────────────────────────────────────────────────────────────────────┘

	fmt.Println("go run main.go      - Compile and run (dev)")
	fmt.Println("go build            - Create executable")
	fmt.Println("go install          - Build and install to GOPATH/bin")
	fmt.Println("go run -race main.go - Run with race detector")

	// ========================================================================
	// SECTION 7: Environment Information
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Environment Information")
	fmt.Println("─────────────────────────────────────")

	fmt.Printf("GOOS (target OS):     %s\n", runtime.GOOS)
	fmt.Printf("GOARCH (target arch): %s\n", runtime.GOARCH)
	fmt.Printf("Go version:           %s\n", runtime.Version())
	fmt.Printf("Num CPUs:             %d\n", runtime.NumCPU())
	fmt.Printf("GOMAXPROCS:           %d\n", runtime.GOMAXPROCS(0))

	// ========================================================================
	// SECTION 8: Program Timing
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Program Timing")
	fmt.Println("─────────────────────────────────────")

	elapsed := time.Since(programStart)
	fmt.Printf("Program started at: %s\n", programStart.Format(time.RFC3339))
	fmt.Printf("Time elapsed: %v\n", elapsed)

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Program completed successfully!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. UNUSED IMPORTS
   ────────────────────────────────────────────────────────────────────────────
   Error: "imported and not used: 'fmt'"
   Why: Go prohibits unused dependencies to keep build times fast.
   Fix: Remove the import, or use blank identifier: _ = fmt.Println

2. BRACE PLACEMENT
   ────────────────────────────────────────────────────────────────────────────
   Error: "syntax error: unexpected newline, expecting { after )"
   Why: Go's lexer inserts semicolons automatically.

   WRONG:                       CORRECT:
   func main()                  func main() {
   {                                ...
       ...                      }
   }

3. MISSING PACKAGE DECLARATION
   ────────────────────────────────────────────────────────────────────────────
   Error: "expected 'package', found 'import'"
   Why: Every Go file must start with a package declaration.
   Fix: Add 'package main' (or appropriate package name) as first line.

4. CALLING init() OR main()
   ────────────────────────────────────────────────────────────────────────────
   Error: "undefined: init" or "cannot call main"
   Why: init() and main() are special functions called by runtime.
   Fix: Don't call them directly; they're invoked automatically.

5. CIRCULAR IMPORTS
   ────────────────────────────────────────────────────────────────────────────
   Error: "import cycle not allowed"
   Why: Package A imports B, and B imports A (directly or indirectly).
   Fix: Refactor to break the cycle, often by creating a third package.

6. MULTIPLE main() FUNCTIONS
   ────────────────────────────────────────────────────────────────────────────
   Error: "main redeclared in this block"
   Why: Only one main() allowed per package main.
   Fix: Keep only one main(), move others to separate packages or remove.
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. FORMATTING
   - Always run 'go fmt' or 'gofmt' before committing
   - Use 'goimports' to automatically manage imports
   - Let your editor do this automatically on save

2. PACKAGE NAMING
   - Lowercase, single-word names
   - Avoid underscores and mixedCaps
   - Package name should match directory name
   - Good: http, json, bytes
   - Bad: httpUtil, Http_Utils

3. IMPORT ORGANIZATION
   - Group standard library imports separately from third-party
   - Use goimports to maintain consistent ordering

   Example:
   import (
       "fmt"           // Standard library
       "net/http"

       "github.com/pkg/errors"  // Third-party

       "mycompany/internal/config"  // Internal
   )

4. init() USAGE
   - Keep init() functions simple and fast
   - Avoid side effects that are hard to test
   - Consider explicit initialization functions instead
   - Never rely on init() order across packages (undefined behavior)

5. ERROR HANDLING
   - main() should handle all errors and call os.Exit() appropriately
   - Don't let panics escape main() in production code
   - Consider wrapping main logic in a run() function that returns error

   func main() {
       if err := run(); err != nil {
           fmt.Fprintf(os.Stderr, "error: %v\n", err)
           os.Exit(1)
       }
   }

   func run() error {
       // All your main logic here
       return nil
   }
*/
