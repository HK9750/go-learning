// ============================================================================
// GO FUNDAMENTALS: FUNCTIONS
// ============================================================================
// This file provides a comprehensive guide to Go functions, including
// parameters, return values, variadic functions, closures, defer, and more.
// ============================================================================

package main

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// SECTION: Function Declarations
// ============================================================================

// Basic function with no parameters and no return
func sayHello() {
	fmt.Println("Hello!")
}

// Function with parameters
func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Function with multiple parameters of same type
func add(a, b int) int {
	return a + b
}

// Function with multiple parameters of different types
func describe(name string, age int, height float64) {
	fmt.Printf("%s is %d years old and %.1f cm tall\n", name, age, height)
}

// Function with multiple return values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Function with named return values
// The return values are initialized to their zero values
func splitSum(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return // "naked return" - returns x and y
}

// Function with named return values (explicit return)
// Better for longer functions - more readable
func rectangle(width, height float64) (area, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)
	return area, perimeter // Explicit return is clearer
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    GO FUNCTIONS                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// ========================================================================
	// SECTION 1: Basic Function Calls
	// ========================================================================
	fmt.Println("\n▶ SECTION 1: Basic Function Calls")
	fmt.Println("─────────────────────────────────────────")

	// FUNCTION SYNTAX:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ func functionName(param1 type1, param2 type2) (return1 type, ...) {   │
	// │     // function body                                                   │
	// │     return value1, value2, ...                                        │
	// │ }                                                                      │
	// └────────────────────────────────────────────────────────────────────────┘

	sayHello()
	greet("Gopher")

	result := add(5, 3)
	fmt.Printf("5 + 3 = %d\n", result)

	describe("Alice", 30, 165.5)

	// Multiple return values
	quotient, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", quotient)
	}

	// Handling error case
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Expected error:", err)
	}

	// Named returns
	x, y := splitSum(100)
	fmt.Printf("splitSum(100) = %d, %d\n", x, y)

	area, perimeter := rectangle(5, 3)
	fmt.Printf("Rectangle 5x3: area=%.1f, perimeter=%.1f\n", area, perimeter)

	// Ignoring return values with blank identifier
	quotient2, _ := divide(20, 4) // Ignore error
	fmt.Printf("20 / 4 = %.1f (error ignored)\n", quotient2)

	// ========================================================================
	// SECTION 2: Variadic Functions
	// ========================================================================
	fmt.Println("\n▶ SECTION 2: Variadic Functions")
	fmt.Println("─────────────────────────────────────────")

	// VARIADIC FUNCTION SYNTAX:
	// ┌────────────────────────────────────────────────────────────────────────┐
	// │ func name(prefix string, nums ...int)                                 │
	// │           └── regular params    └── variadic param (must be last)     │
	// │                                                                        │
	// │ Inside function: nums is []int                                        │
	// └────────────────────────────────────────────────────────────────────────┘

	// Calling with individual arguments
	fmt.Printf("sum(1, 2, 3) = %d\n", sum(1, 2, 3))
	fmt.Printf("sum(10, 20, 30, 40) = %d\n", sum(10, 20, 30, 40))
	fmt.Printf("sum() = %d\n", sum()) // No arguments is valid

	// Calling with a slice (unpack with ...)
	numbers := []int{5, 10, 15, 20}
	fmt.Printf("sum(numbers...) = %d\n", sum(numbers...))

	// Mixed regular and variadic parameters
	fmt.Println(concat("-", "a", "b", "c"))

	// ========================================================================
	// SECTION 3: Function Values (First-Class Functions)
	// ========================================================================
	fmt.Println("\n▶ SECTION 3: Function Values")
	fmt.Println("─────────────────────────────────────────")

	// Functions are first-class values in Go:
	// - Can be assigned to variables
	// - Can be passed as arguments
	// - Can be returned from functions
	// - Can be stored in data structures

	// Function as variable
	var operation func(int, int) int
	operation = add
	fmt.Printf("operation(5, 3) = %d\n", operation(5, 3))

	operation = multiply
	fmt.Printf("operation(5, 3) = %d\n", operation(5, 3))

	// Function as argument
	results := applyToEach([]int{1, 2, 3, 4}, square)
	fmt.Printf("Squares: %v\n", results)

	results = applyToEach([]int{1, 2, 3, 4}, double)
	fmt.Printf("Doubles: %v\n", results)

	// Function returning function
	addFive := makeAdder(5)
	addTen := makeAdder(10)
	fmt.Printf("addFive(3) = %d\n", addFive(3))
	fmt.Printf("addTen(3) = %d\n", addTen(3))

	// ========================================================================
	// SECTION 4: Anonymous Functions
	// ========================================================================
	fmt.Println("\n▶ SECTION 4: Anonymous Functions")
	fmt.Println("─────────────────────────────────────────")

	// Anonymous function assigned to variable
	greetAnon := func(name string) string {
		return "Hello, " + name + "!"
	}
	fmt.Println(greetAnon("World"))

	// Immediately invoked function expression (IIFE)
	result2 := func(a, b int) int {
		return a * b
	}(4, 5)
	fmt.Printf("IIFE result: %d\n", result2)

	// Anonymous function for one-time use
	numbers2 := []int{5, 2, 8, 1, 9}
	fmt.Printf("Max: %d\n", func(nums []int) int {
		max := nums[0]
		for _, n := range nums[1:] {
			if n > max {
				max = n
			}
		}
		return max
	}(numbers2))

	// ========================================================================
	// SECTION 5: Closures
	// ========================================================================
	fmt.Println("\n▶ SECTION 5: Closures")
	fmt.Println("─────────────────────────────────────────")

	// A closure is a function that references variables from outside its body.
	// The function "closes over" those variables.

	// CLOSURE VISUALIZATION:
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ func makeCounter() func() int {                                        │
	// │     count := 0  ◄─── This variable is "captured"                       │
	// │     return func() int {                                                │
	// │         count++  ◄─── Closure references 'count'                       │
	// │         return count                                                    │
	// │     }                                                                   │
	// │ }                                                                       │
	// │                                                                         │
	// │ MEMORY:                                                                 │
	// │ ┌───────────────┐     ┌───────────────────────┐                        │
	// │ │ counter1 func │────►│ Closure Environment   │                        │
	// │ └───────────────┘     │ ┌─────────────────┐   │                        │
	// │                       │ │ count: 0 → 1 → 2│   │                        │
	// │ ┌───────────────┐     │ └─────────────────┘   │                        │
	// │ │ counter2 func │────►│ Closure Environment   │ ◄── Separate state!   │
	// │ └───────────────┘     │ ┌─────────────────┐   │                        │
	// │                       │ │ count: 0 → 1    │   │                        │
	// │                       │ └─────────────────┘   │                        │
	// │                       └───────────────────────┘                        │
	// └─────────────────────────────────────────────────────────────────────────┘

	counter1 := makeCounter()
	counter2 := makeCounter()

	fmt.Printf("counter1: %d\n", counter1()) // 1
	fmt.Printf("counter1: %d\n", counter1()) // 2
	fmt.Printf("counter1: %d\n", counter1()) // 3
	fmt.Printf("counter2: %d\n", counter2()) // 1 (separate state!)
	fmt.Printf("counter1: %d\n", counter1()) // 4

	// Practical closure: Function factory
	double2 := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(5) = %d\n", double2(5))
	fmt.Printf("triple(5) = %d\n", triple(5))

	// Practical closure: Memoization
	fibMemo := memoize(slowFib)
	fmt.Printf("fib(10) = %d\n", fibMemo(10))
	fmt.Printf("fib(10) = %d (cached)\n", fibMemo(10))

	// ========================================================================
	// SECTION 6: Defer In-Depth
	// ========================================================================
	fmt.Println("\n▶ SECTION 6: Defer In-Depth")
	fmt.Println("─────────────────────────────────────────")

	// DEFER RULES:
	// ┌──────────────────────────────────────────────────────────────────────────┐
	// │ 1. Arguments evaluated when defer is scheduled, not when it runs        │
	// │ 2. Deferred functions run in LIFO (Last In, First Out) order           │
	// │ 3. Deferred functions can read/modify named return values               │
	// │ 4. Deferred functions run even if function panics                       │
	// └──────────────────────────────────────────────────────────────────────────┘

	// Rule 1: Arguments evaluated immediately
	fmt.Println("\nRule 1: Argument evaluation")
	deferArgDemo()

	// Rule 2: LIFO order
	fmt.Println("\nRule 2: LIFO execution")
	deferOrderDemo()

	// Rule 3: Modifying named return values
	fmt.Println("\nRule 3: Named return modification")
	result3 := deferModifyReturn()
	fmt.Printf("Result: %d\n", result3)

	// Practical use: Resource cleanup
	fmt.Println("\nPractical: Resource cleanup pattern")
	processFile() // Simulated file processing

	// Practical use: Timing
	fmt.Println("\nPractical: Function timing")
	timedOperation()

	// Practical use: Unlock
	fmt.Println("\nPractical: Mutex pattern")
	mutexDemo()

	// ========================================================================
	// SECTION 7: Recursion
	// ========================================================================
	fmt.Println("\n▶ SECTION 7: Recursion")
	fmt.Println("─────────────────────────────────────────")

	fmt.Printf("factorial(5) = %d\n", factorial(5))
	fmt.Printf("fibonacci(10) = %d\n", fibonacci(10))

	// Recursive data structure traversal
	tree := &TreeNode{
		Value: 1,
		Left: &TreeNode{
			Value: 2,
			Left:  &TreeNode{Value: 4},
			Right: &TreeNode{Value: 5},
		},
		Right: &TreeNode{
			Value: 3,
			Left:  &TreeNode{Value: 6},
			Right: &TreeNode{Value: 7},
		},
	}
	fmt.Print("Tree inorder: ")
	inorderTraversal(tree)
	fmt.Println()

	// Tail recursion (Go doesn't optimize, but pattern is useful)
	fmt.Printf("factorialTail(5) = %d\n", factorialTail(5, 1))

	// ========================================================================
	// SECTION 8: Method vs Function
	// ========================================================================
	fmt.Println("\n▶ SECTION 8: Methods Preview")
	fmt.Println("─────────────────────────────────────────")

	// Methods are functions with a receiver argument.
	// This is covered in detail in the Methods section.

	// METHOD SYNTAX:
	// func (r ReceiverType) MethodName(params) returnType { }

	p := Point{X: 3, Y: 4}
	fmt.Printf("Point: %v\n", p)
	fmt.Printf("Distance from origin: %.2f\n", p.Distance())
	p.Scale(2)
	fmt.Printf("After scaling: %v\n", p)

	// ========================================================================
	// SECTION 9: Function Type Aliases
	// ========================================================================
	fmt.Println("\n▶ SECTION 9: Function Types")
	fmt.Println("─────────────────────────────────────────")

	// You can create named function types
	// type HandlerFunc func(request string) string

	// This is useful for:
	// 1. Documentation
	// 2. Adding methods to function types
	// 3. Making signatures clearer

	var handler HandlerFunc = func(req string) string {
		return "Handled: " + req
	}

	fmt.Println(handler("test request"))
	fmt.Println(handler.Log("logged request"))

	// ========================================================================
	// SECTION 10: Init Functions
	// ========================================================================
	fmt.Println("\n▶ SECTION 10: Init Functions")
	fmt.Println("─────────────────────────────────────────")

	// init() functions are special:
	// - Called automatically before main()
	// - Can have multiple per file
	// - Cannot be called directly
	// - Run in order they appear in file
	// - Run after all imported packages init

	fmt.Println("(init() already ran before main - see top of output)")

	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("  Functions Complete!")
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// Supporting Functions and Types
// ============================================================================

// Variadic function
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Mixed regular and variadic parameters
func concat(separator string, items ...string) string {
	return strings.Join(items, separator)
}

// Function for use as value
func multiply(a, b int) int {
	return a * b
}

func square(x int) int {
	return x * x
}

func double(x int) int {
	return x * 2
}

// Higher-order function (takes function as argument)
func applyToEach(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

// Function returning function
func makeAdder(n int) func(int) int {
	return func(x int) int {
		return x + n
	}
}

// Closure: Counter factory
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// Closure: Multiplier factory
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// Memoization closure
func memoize(fn func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(n int) int {
		if v, ok := cache[n]; ok {
			fmt.Print("(cached) ")
			return v
		}
		result := fn(n)
		cache[n] = result
		return result
	}
}

func slowFib(n int) int {
	if n <= 1 {
		return n
	}
	return slowFib(n-1) + slowFib(n-2)
}

// Defer demonstrations
func deferArgDemo() {
	x := 10
	defer fmt.Printf("  Deferred x = %d (captured when defer was called)\n", x)
	x = 20
	fmt.Printf("  Current x = %d\n", x)
}

func deferOrderDemo() {
	defer fmt.Println("  Third (LIFO: runs first)")
	defer fmt.Println("  Second")
	defer fmt.Println("  First (LIFO: runs last)")
	fmt.Println("  Function body")
}

func deferModifyReturn() (result int) {
	defer func() {
		result += 10 // Modifies the named return value!
	}()
	return 5 // result = 5, then defer runs, result = 15
}

func processFile() {
	fmt.Println("  Opening file...")
	// file, _ := os.Open("file.txt")
	// defer file.Close()  // Guaranteed to run!
	defer fmt.Println("  File closed (cleanup)")

	fmt.Println("  Processing file...")
	// ... do work ...
}

func timedOperation() {
	start := time.Now()
	defer func() {
		fmt.Printf("  Operation took %v\n", time.Since(start))
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("  Expensive operation...")
}

var mu sync.Mutex

func mutexDemo() {
	mu.Lock()
	defer mu.Unlock() // Always unlock, even if panic

	fmt.Println("  Critical section protected by mutex")
}

// Recursion
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// Tail recursive (Go doesn't optimize, but pattern is still useful)
func factorialTail(n, acc int) int {
	if n <= 1 {
		return acc
	}
	return factorialTail(n-1, n*acc)
}

// Tree structure for recursion demo
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraversal(node *TreeNode) {
	if node == nil {
		return
	}
	inorderTraversal(node.Left)
	fmt.Printf("%d ", node.Value)
	inorderTraversal(node.Right)
}

// Method example (preview)
type Point struct {
	X, Y float64
}

func (p Point) Distance() float64 {
	return (p.X*p.X + p.Y*p.Y)
}

func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

// Function type alias with method
type HandlerFunc func(string) string

func (h HandlerFunc) Log(req string) string {
	fmt.Printf("  [LOG] Handling: %s\n", req)
	return h(req)
}

// Init function (runs before main)
func init() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  INIT: Functions module initialized (Go %s)\n", runtime.Version())
	fmt.Println("═══════════════════════════════════════════════════════════")
}

// ============================================================================
// ERROR ANALYSIS & COMMON MISTAKES
// ============================================================================
/*
1. NIL FUNCTION CALL
   ─────────────────────────────────────────────────────────────────────────
   Error: "panic: runtime error: invalid memory address or nil pointer dereference"

   var fn func()  // nil
   fn()           // PANIC!

   Fix: Check if fn != nil before calling.

2. NAKED RETURN CONFUSION
   ─────────────────────────────────────────────────────────────────────────
   func bad() (x int) {
       x := 10  // This creates NEW variable, shadows return value!
       return   // Returns 0, not 10!
   }

   Fix: Use = not := for named return values.

3. DEFER IN LOOP
   ─────────────────────────────────────────────────────────────────────────
   for _, file := range files {
       f, _ := os.Open(file)
       defer f.Close()  // BAD! All closes happen at function end, not loop end
   }

   Fix: Wrap in anonymous function or call Close() explicitly.

   for _, file := range files {
       func() {
           f, _ := os.Open(file)
           defer f.Close()
           // process file
       }()  // defer runs here
   }

4. CAPTURING LOOP VARIABLE (Pre-Go 1.22)
   ─────────────────────────────────────────────────────────────────────────
   for i := 0; i < 3; i++ {
       go func() {
           fmt.Println(i)  // May print "3 3 3" in older Go
       }()
   }

   Fix: Shadow the variable or pass as argument.

   for i := 0; i < 3; i++ {
       i := i  // Shadow
       go func() { fmt.Println(i) }()
   }

5. IGNORING ERRORS
   ─────────────────────────────────────────────────────────────────────────
   result, _ := mightFail()  // Ignoring error!

   Fix: Always handle or explicitly ignore with comment.

   result, err := mightFail()
   if err != nil {
       return err
   }

6. VARIADIC SLICE TYPE MISMATCH
   ─────────────────────────────────────────────────────────────────────────
   func printStrings(strs ...string) { }

   args := []interface{}{"a", "b"}
   printStrings(args...)  // ERROR: cannot use []interface{} as []string

   Fix: Create correct slice type or use interface{}.

7. DEFER ARGUMENT EVALUATION TIMING
   ─────────────────────────────────────────────────────────────────────────
   x := 10
   defer fmt.Println(x)  // Prints 10, not 20!
   x = 20

   If you want current value, use closure:
   defer func() { fmt.Println(x) }()
*/

// ============================================================================
// BEST PRACTICES
// ============================================================================
/*
1. PREFER SMALL, FOCUSED FUNCTIONS
   Each function should do one thing well.

2. RETURN EARLY FOR ERRORS
   if err != nil {
       return err
   }
   // happy path continues

3. USE NAMED RETURNS SPARINGLY
   Only when it aids documentation.
   Avoid naked returns in long functions.

4. DEFER FOR RESOURCE CLEANUP
   Place defer immediately after resource acquisition.

   file, err := os.Open(path)
   if err != nil { return err }
   defer file.Close()

5. PASS BY VALUE FOR SMALL TYPES
   Passing small structs by value is often faster than by pointer.
   Pointer only needed for: mutation, large structs, or interface requirements.

6. USE FUNCTION TYPES FOR CALLBACKS
   type Callback func(result string, err error)

   func fetchAsync(url string, cb Callback) { ... }

7. DOCUMENT FUNCTIONS WITH COMMENTS
   // ParseJSON parses the JSON data and returns the result.
   // It returns an error if the JSON is malformed.
   func ParseJSON(data []byte) (Result, error) { ... }

8. USE CLOSURES FOR STATE ENCAPSULATION
   Closures can encapsulate private state without needing a struct.

9. AVOID DEEP RECURSION
   Go doesn't have tail call optimization.
   Convert to iteration for deep recursion.

10. CHECK FUNCTION VALUES FOR NIL
    if handler != nil {
        handler()
    }
*/
