package main

import (
	"fmt"
	"unsafe"
)

// =============================================================================
// METHODS IN GO - COMPREHENSIVE DEEP DIVE
// =============================================================================
//
// TOPIC: Methods - Functions with Receivers
//
// Go does not have classes. Instead, you can define methods on types.
// A method is a function with a special "receiver" argument that appears
// between the 'func' keyword and the method name.
//
// Why methods instead of classes?
// - Simplicity: No inheritance hierarchies to manage
// - Composition over inheritance: Embed types to reuse functionality
// - Clear ownership: Methods are defined near their types
// - Interface satisfaction: Methods enable polymorphism
//
// =============================================================================
// VISUALIZATION: Method vs Function
// =============================================================================
//
//  REGULAR FUNCTION:
//  +------------------------------------------------------+
//  |  func PrintUser(u User) { ... }                      |
//  |       └── User passed as regular parameter           |
//  +------------------------------------------------------+
//
//  METHOD:
//  +------------------------------------------------------+
//  |  func (u User) Print() { ... }                       |
//  |       └── User is the "receiver" (special position)  |
//  +------------------------------------------------------+
//
//  Under the hood, they compile to the same thing!
//  The receiver is just syntactic sugar for the first parameter.
//
// =============================================================================
// VISUALIZATION: Method Sets
// =============================================================================
//
// A "method set" determines which methods can be called on a value.
// This is CRITICAL for interface satisfaction!
//
//  +-----------------------+--------------------------------+
//  |  Receiver Type        |  Method Set Contains           |
//  +-----------------------+--------------------------------+
//  |  Value (T)            |  Methods with (t T) receivers  |
//  +-----------------------+--------------------------------+
//  |  Pointer (*T)         |  Methods with (t T) receivers  |
//  |                       |  Methods with (t *T) receivers |
//  +-----------------------+--------------------------------+
//
//  WHY? A pointer can always dereference to get a value,
//  but a value cannot always get its address (e.g., map values).
//
// =============================================================================
// VISUALIZATION: Value vs Pointer Receiver (Memory View)
// =============================================================================
//
//  VALUE RECEIVER - Copy is made:
//
//  Original:                   Copy (inside method):
//  +-------------------+       +-------------------+
//  |  u User           |       |  u User (COPY)    |
//  |  Name: "Alice"    |  -->  |  Name: "Alice"    |
//  |  Age:  25         |       |  Age:  25         |
//  +-------------------+       +-------------------+
//       [0x1000]                    [0x2000]
//                                     ↓
//                              Modifications here
//                              DO NOT affect original!
//
//  POINTER RECEIVER - No copy, direct access:
//
//  Original:                   Inside method:
//  +-------------------+       +-------------------+
//  |  u User           |  <--  |  u *User          |
//  |  Name: "Alice"    |       |  (points to       |
//  |  Age:  25         |       |   original)       |
//  +-------------------+       +-------------------+
//       [0x1000]                 [value: 0x1000]
//                                     ↓
//                              Modifications here
//                              AFFECT original!
//
// =============================================================================

// -----------------------------------------------------------------------------
// SECTION 1: Basic Types for Methods
// -----------------------------------------------------------------------------

// User demonstrates methods on a struct type
type User struct {
	Name   string
	Email  string
	Age    int
	Logins int
}

// Counter demonstrates methods on a simple type alias
type Counter int

// MyFloat demonstrates methods on any named type
type MyFloat float64

// Point demonstrates geometric methods
type Point struct {
	X, Y float64
}

// -----------------------------------------------------------------------------
// SECTION 2: Value Receivers
// -----------------------------------------------------------------------------
//
// Value receivers receive a COPY of the value.
//
// When to use value receivers:
// 1. When the method does NOT need to modify the receiver
// 2. When the receiver is small (primitives, small structs)
// 3. When you want to work with a snapshot of the data
// 4. When the type is inherently a value (time.Time, small structs)
//

// String returns a string representation of User (value receiver)
// The receiver 'u' is a COPY of the original User
func (u User) String() string {
	return fmt.Sprintf("User{Name: %q, Email: %q, Age: %d, Logins: %d}",
		u.Name, u.Email, u.Age, u.Logins)
}

// FullInfo demonstrates that value receiver gets a copy
func (u User) FullInfo() string {
	// This modification only affects the local copy!
	u.Name = "MODIFIED (but only locally)"
	return fmt.Sprintf("%s <%s> - Age: %d", u.Name, u.Email, u.Age)
}

// Value method on Counter - works with the counter value
func (c Counter) Value() int {
	return int(c)
}

// Double returns a new Counter with doubled value (value receiver)
// Notice: returns new value, doesn't modify receiver
func (c Counter) Double() Counter {
	return c * 2
}

// Abs returns the absolute value of MyFloat
func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

// Distance calculates distance from origin
func (p Point) Distance() float64 {
	return (p.X*p.X + p.Y*p.Y) // simplified, would use math.Sqrt in real code
}

// -----------------------------------------------------------------------------
// SECTION 3: Pointer Receivers
// -----------------------------------------------------------------------------
//
// Pointer receivers receive a POINTER to the original value.
//
// When to use pointer receivers:
// 1. When the method needs to MODIFY the receiver
// 2. When the receiver is large (avoids copying)
// 3. When consistency is needed (if one method needs pointer, use for all)
// 4. When the type contains fields that should not be copied (sync.Mutex)
//

// IncrementLogin modifies the original User (pointer receiver)
func (u *User) IncrementLogin() {
	// IMPORTANT: Always check for nil receivers!
	if u == nil {
		fmt.Println("Warning: IncrementLogin called on nil User")
		return
	}
	u.Logins++
}

// SetEmail modifies the email field
func (u *User) SetEmail(email string) {
	if u == nil {
		return
	}
	u.Email = email
}

// Birthday increases age by 1
func (u *User) Birthday() {
	if u == nil {
		return
	}
	u.Age++
}

// Reset clears login count
func (u *User) Reset() {
	if u == nil {
		return
	}
	u.Logins = 0
}

// Increment adds 1 to the Counter (pointer receiver)
func (c *Counter) Increment() {
	if c == nil {
		return
	}
	*c++
}

// Add adds n to the Counter
func (c *Counter) Add(n int) {
	if c == nil {
		return
	}
	*c += Counter(n)
}

// Scale multiplies the point coordinates
func (p *Point) Scale(factor float64) {
	if p == nil {
		return
	}
	p.X *= factor
	p.Y *= factor
}

// Translate moves the point
func (p *Point) Translate(dx, dy float64) {
	if p == nil {
		return
	}
	p.X += dx
	p.Y += dy
}

// -----------------------------------------------------------------------------
// SECTION 4: Nil Receivers - A Unique Go Feature
// -----------------------------------------------------------------------------
//
// Unlike most languages, Go allows calling methods on nil pointers!
// This is safe IF the method handles the nil case.
//
// VISUALIZATION: Nil Receiver Call
//
//  var u *User = nil
//  u.IncrementLogin()  // This WORKS! (if method handles nil)
//
//  +-------------------+
//  |  u *User = nil    |      Method receives:
//  |  [value: 0x0]     | -->  u *User = nil (0x0)
//  +-------------------+
//                             if u == nil { return }  // SAFE!
//                             u.Logins++              // Would panic!
//

// NilSafeString demonstrates nil-safe method
func (u *User) NilSafeString() string {
	if u == nil {
		return "User(nil)"
	}
	return u.String()
}

// -----------------------------------------------------------------------------
// SECTION 5: Method Expressions and Method Values
// -----------------------------------------------------------------------------
//
// Methods can be extracted and used as regular functions.
//
// VISUALIZATION: Method Expression vs Method Value
//
//  METHOD EXPRESSION (from type):
//  +------------------------------------------+
//  |  f := User.String                        |
//  |  // f has signature: func(User) string   |
//  |  f(someUser)  // must pass receiver      |
//  +------------------------------------------+
//
//  METHOD VALUE (from instance):
//  +------------------------------------------+
//  |  f := someUser.String                    |
//  |  // f has signature: func() string       |
//  |  f()  // receiver already bound          |
//  +------------------------------------------+
//

// -----------------------------------------------------------------------------
// SECTION 6: Methods on Non-Struct Types
// -----------------------------------------------------------------------------

// ByteSlice demonstrates methods on a slice type
type ByteSlice []byte

// Append adds bytes to the slice (pointer receiver because it may reallocate)
func (b *ByteSlice) Append(data ...byte) {
	*b = append(*b, data...)
}

// String implements Stringer interface
func (b ByteSlice) String() string {
	return fmt.Sprintf("ByteSlice(%d bytes): %v", len(b), []byte(b))
}

// StringList demonstrates methods on a slice of strings
type StringList []string

// Filter returns a new list with only matching elements
func (s StringList) Filter(predicate func(string) bool) StringList {
	result := make(StringList, 0)
	for _, item := range s {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map transforms each element
func (s StringList) Map(transform func(string) string) StringList {
	result := make(StringList, len(s))
	for i, item := range s {
		result[i] = transform(item)
	}
	return result
}

// IntMap demonstrates methods on a map type
type IntMap map[string]int

// GetOrDefault returns value or default if key not found
func (m IntMap) GetOrDefault(key string, defaultVal int) int {
	if val, ok := m[key]; ok {
		return val
	}
	return defaultVal
}

// Keys returns all keys
func (m IntMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// -----------------------------------------------------------------------------
// SECTION 7: Embedding and Method Promotion
// -----------------------------------------------------------------------------
//
// When you embed a type, its methods are "promoted" to the outer type.
//
// VISUALIZATION: Method Promotion
//
//  type Inner struct { ... }
//  func (i Inner) Hello() { ... }
//
//  type Outer struct {
//      Inner  // embedded (anonymous field)
//  }
//
//  o := Outer{}
//  o.Hello()  // Works! Promoted from Inner
//
//  Compiler transforms:  o.Hello()  -->  o.Inner.Hello()
//

// Base provides common functionality
type Base struct {
	ID   int
	Name string
}

func (b Base) Identify() string {
	return fmt.Sprintf("ID=%d, Name=%s", b.ID, b.Name)
}

func (b *Base) SetName(name string) {
	b.Name = name
}

// Admin embeds Base, inheriting its methods
type Admin struct {
	Base        // Embedded - methods are promoted
	Permissions []string
}

// Override Identify to include admin-specific info
func (a Admin) Identify() string {
	return fmt.Sprintf("ADMIN: %s, Permissions: %v", a.Base.Identify(), a.Permissions)
}

// Employee also embeds Base
type Employee struct {
	Base
	Department string
	Salary     float64
}

// Additional method specific to Employee
func (e Employee) PayInfo() string {
	return fmt.Sprintf("%s works in %s, earns $%.2f", e.Name, e.Department, e.Salary)
}

// -----------------------------------------------------------------------------
// SECTION 8: Chaining Methods (Fluent Interface)
// -----------------------------------------------------------------------------

// Builder demonstrates method chaining pattern
type Builder struct {
	data []byte
}

func NewBuilder() *Builder {
	return &Builder{data: make([]byte, 0)}
}

// Write adds bytes and returns self for chaining
func (b *Builder) Write(data string) *Builder {
	b.data = append(b.data, []byte(data)...)
	return b
}

// WriteByte adds a single byte
func (b *Builder) WriteByte(c byte) *Builder {
	b.data = append(b.data, c)
	return b
}

// WriteInt adds an integer as string
func (b *Builder) WriteInt(n int) *Builder {
	b.data = append(b.data, []byte(fmt.Sprintf("%d", n))...)
	return b
}

// Build returns the final string
func (b *Builder) Build() string {
	return string(b.data)
}

// -----------------------------------------------------------------------------
// SECTION 9: Understanding Method Dispatch - Under the Hood
// -----------------------------------------------------------------------------
//
// DEEP DIVE: How Go Dispatches Methods
//
// When you call u.Method(), Go looks up the method in the type's method table.
// For concrete types, this is a direct call (static dispatch).
// For interfaces, this goes through an indirect call (dynamic dispatch).
//
// VISUALIZATION: Method Table
//
//  type User struct { ... }
//  func (u User) String() string { ... }
//  func (u *User) SetEmail(s string) { ... }
//
//  COMPILE TIME: Go builds method tables
//
//  User Method Table (*typeinfo):     *User Method Table:
//  +-------------------------+        +-------------------------+
//  | String: func_addr_1     |        | String: func_addr_1     |
//  +-------------------------+        | SetEmail: func_addr_2   |
//                                     +-------------------------+
//
//  RUNTIME: Direct function call
//  u.String() --> call func_addr_1(u)
//

// -----------------------------------------------------------------------------
// MAIN FUNCTION - Demonstrations
// -----------------------------------------------------------------------------

func main() {
	fmt.Println("=== GO METHODS - COMPREHENSIVE GUIDE ===\n")

	// =========================================================================
	// DEMO 1: Value vs Pointer Receivers
	// =========================================================================
	fmt.Println("--- DEMO 1: Value vs Pointer Receivers ---")

	user := User{
		Name:   "Alice",
		Email:  "alice@example.com",
		Age:    28,
		Logins: 0,
	}

	fmt.Println("Original:", user.String())

	// Call value receiver method - gets a copy
	info := user.FullInfo()
	fmt.Println("FullInfo (modified copy):", info)
	fmt.Println("After FullInfo (original unchanged):", user.String())

	// Call pointer receiver method - modifies original
	user.IncrementLogin()
	user.IncrementLogin()
	fmt.Println("After 2 IncrementLogin calls:", user.String())

	user.Birthday()
	fmt.Println("After Birthday:", user.String())

	// =========================================================================
	// DEMO 2: Automatic Pointer/Value Conversion
	// =========================================================================
	fmt.Println("\n--- DEMO 2: Automatic Conversion ---")

	// Go automatically takes address when needed
	v := User{Name: "Bob", Age: 30}
	v.IncrementLogin() // Go does: (&v).IncrementLogin()
	fmt.Println("Value calling pointer method:", v.Logins)

	// Go automatically dereferences when needed
	p := &User{Name: "Carol", Age: 25}
	fmt.Println("Pointer calling value method:", p.String()) // Go does: (*p).String()

	// BUT: This only works for addressable values!
	// User{Name: "X"}.IncrementLogin() // ERROR: cannot take address

	// =========================================================================
	// DEMO 3: Nil Receivers
	// =========================================================================
	fmt.Println("\n--- DEMO 3: Nil Receivers ---")

	var nilUser *User = nil

	// This works because IncrementLogin handles nil!
	nilUser.IncrementLogin() // Output: Warning message

	// This also works
	fmt.Println("Nil user string:", nilUser.NilSafeString())

	// DANGER: If method doesn't check nil:
	// nilUser.UnsafeMethod() // Would panic!

	// =========================================================================
	// DEMO 4: Method Expressions
	// =========================================================================
	fmt.Println("\n--- DEMO 4: Method Expressions ---")

	// Method Expression from VALUE type
	// Signature: func(User) string
	stringFunc := User.String
	fmt.Printf("Method Expression type: %T\n", stringFunc)
	fmt.Println("Called:", stringFunc(user))

	// Method Expression from POINTER type
	// Signature: func(*User)
	incFunc := (*User).IncrementLogin
	fmt.Printf("Pointer Method Expression type: %T\n", incFunc)
	incFunc(&user)
	fmt.Println("After incFunc:", user.Logins)

	// =========================================================================
	// DEMO 5: Method Values (Bound Methods)
	// =========================================================================
	fmt.Println("\n--- DEMO 5: Method Values ---")

	// Method Value - receiver is already bound
	boundString := user.String // Captures 'user'
	fmt.Printf("Method Value type: %T\n", boundString)
	fmt.Println("Called:", boundString()) // No argument needed!

	// Useful for callbacks
	printFuncs := []func() string{
		User{Name: "X"}.String,
		User{Name: "Y"}.String,
		User{Name: "Z"}.String,
	}
	for i, f := range printFuncs {
		fmt.Printf("  [%d] %s\n", i, f())
	}

	// =========================================================================
	// DEMO 6: Methods on Non-Struct Types
	// =========================================================================
	fmt.Println("\n--- DEMO 6: Methods on Non-Struct Types ---")

	// Counter type
	var count Counter = 5
	fmt.Println("Counter value:", count.Value())
	fmt.Println("Counter doubled:", count.Double())
	count.Increment()
	count.Add(10)
	fmt.Println("After increment and add(10):", count)

	// MyFloat type
	mf := MyFloat(-3.14)
	fmt.Println("MyFloat abs:", mf.Abs())

	// ByteSlice type
	var bs ByteSlice
	bs.Append(65, 66, 67) // A, B, C
	fmt.Println(bs.String())

	// StringList with functional methods
	names := StringList{"Alice", "Bob", "Charlie", "Diana"}
	longNames := names.Filter(func(s string) bool {
		return len(s) > 4
	})
	fmt.Println("Long names:", longNames)

	upperNames := names.Map(func(s string) string {
		return s + "!"
	})
	fmt.Println("With exclamation:", upperNames)

	// IntMap type
	scores := IntMap{"Alice": 95, "Bob": 87}
	fmt.Println("Carol's score:", scores.GetOrDefault("Carol", 0))
	fmt.Println("Keys:", scores.Keys())

	// =========================================================================
	// DEMO 7: Embedding and Method Promotion
	// =========================================================================
	fmt.Println("\n--- DEMO 7: Embedding and Method Promotion ---")

	admin := Admin{
		Base:        Base{ID: 1, Name: "SuperAdmin"},
		Permissions: []string{"read", "write", "delete"},
	}

	// Promoted method from Base (overridden)
	fmt.Println("Admin.Identify():", admin.Identify())

	// Can still access Base's method directly
	fmt.Println("Admin.Base.Identify():", admin.Base.Identify())

	// Promoted method (not overridden)
	admin.SetName("MegaAdmin") // Calls Base.SetName
	fmt.Println("After SetName:", admin.Identify())

	employee := Employee{
		Base:       Base{ID: 42, Name: "John"},
		Department: "Engineering",
		Salary:     75000,
	}
	fmt.Println("Employee.Identify():", employee.Identify()) // From Base
	fmt.Println("Employee.PayInfo():", employee.PayInfo())

	// =========================================================================
	// DEMO 8: Method Chaining (Fluent Interface)
	// =========================================================================
	fmt.Println("\n--- DEMO 8: Method Chaining ---")

	result := NewBuilder().
		Write("Hello, ").
		Write("World").
		WriteByte('!').
		Write(" Count: ").
		WriteInt(42).
		Build()

	fmt.Println("Built string:", result)

	// =========================================================================
	// DEMO 9: Method Internals (Advanced)
	// =========================================================================
	fmt.Println("\n--- DEMO 9: Method Internals (Advanced) ---")

	// Show that method calls are just function calls
	u := User{Name: "Test", Email: "test@test.com", Age: 30, Logins: 5}

	// These are equivalent:
	fmt.Println("Method call:      ", u.String())
	fmt.Println("Function call:    ", User.String(u))

	// Memory layout demonstration
	fmt.Printf("\nUser struct size: %d bytes\n", unsafe.Sizeof(u))
	fmt.Printf("User struct alignment: %d\n", unsafe.Alignof(u))

	// Field offsets
	fmt.Printf("Offset of Name:   %d\n", unsafe.Offsetof(u.Name))
	fmt.Printf("Offset of Email:  %d\n", unsafe.Offsetof(u.Email))
	fmt.Printf("Offset of Age:    %d\n", unsafe.Offsetof(u.Age))
	fmt.Printf("Offset of Logins: %d\n", unsafe.Offsetof(u.Logins))

	// =========================================================================
	// DEMO 10: Common Patterns
	// =========================================================================
	fmt.Println("\n--- DEMO 10: Common Patterns ---")

	// Constructor pattern
	user2 := NewUser("Frank", "frank@email.com", 35)
	fmt.Println("New user:", user2.String())

	// Validate pattern
	invalid := &User{Name: "", Age: -5}
	if err := invalid.Validate(); err != nil {
		fmt.Println("Validation error:", err)
	}

	// Clone pattern
	original := &User{Name: "Original", Email: "orig@test.com", Age: 25, Logins: 10}
	clone := original.Clone()
	clone.Name = "Clone"
	fmt.Println("Original:", original.Name)
	fmt.Println("Clone:", clone.Name)
}

// -----------------------------------------------------------------------------
// SECTION 10: Common Patterns
// -----------------------------------------------------------------------------

// NewUser is a constructor function (not a method, but commonly used pattern)
func NewUser(name, email string, age int) *User {
	return &User{
		Name:   name,
		Email:  email,
		Age:    age,
		Logins: 0,
	}
}

// Validate checks if User has valid data
func (u *User) Validate() error {
	if u == nil {
		return fmt.Errorf("user is nil")
	}
	if u.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if u.Age < 0 {
		return fmt.Errorf("age cannot be negative: %d", u.Age)
	}
	return nil
}

// Clone creates a deep copy of User
func (u *User) Clone() *User {
	if u == nil {
		return nil
	}
	return &User{
		Name:   u.Name,
		Email:  u.Email,
		Age:    u.Age,
		Logins: u.Logins,
	}
}

// Equal checks if two users are equal
func (u *User) Equal(other *User) bool {
	if u == nil || other == nil {
		return u == other
	}
	return u.Name == other.Name &&
		u.Email == other.Email &&
		u.Age == other.Age &&
		u.Logins == other.Logins
}

/*
===============================================================================
ERROR ANALYSIS & COMMON MISTAKES
===============================================================================

1. MODIFYING VALUE RECEIVER (Most Common Bug!)
   -------------------------
   WRONG:
   func (u User) SetName(name string) {
       u.Name = name  // Modifies COPY, original unchanged!
   }

   CORRECT:
   func (u *User) SetName(name string) {
       u.Name = name  // Modifies original
   }

2. NIL POINTER DEREFERENCE
   ------------------------
   WRONG:
   func (u *User) GetName() string {
       return u.Name  // Panics if u is nil!
   }

   CORRECT:
   func (u *User) GetName() string {
       if u == nil {
           return ""
       }
       return u.Name
   }

3. MIXING RECEIVER TYPES INCONSISTENTLY
   ------------------------------------
   WRONG:
   func (u User) Method1() {}    // value receiver
   func (u *User) Method2() {}   // pointer receiver
   // Inconsistent! Makes interface satisfaction confusing

   CORRECT: Pick one style and stick with it
   func (u *User) Method1() {}   // all pointer
   func (u *User) Method2() {}   // all pointer

   RULE: If ANY method needs pointer receiver, use pointer for ALL methods.

4. CALLING METHODS ON NON-ADDRESSABLE VALUES
   -----------------------------------------
   ERROR: "cannot take address of User{...}"

   WRONG:
   User{Name: "X"}.SetName("Y")  // Cannot take address!

   CORRECT:
   u := User{Name: "X"}
   u.SetName("Y")
   // OR
   (&User{Name: "X"}).SetName("Y")

5. FORGETTING THAT METHOD EXPRESSIONS NEED RECEIVER
   ------------------------------------------------
   WRONG:
   f := User.String
   f()  // Missing receiver!

   CORRECT:
   f := User.String
   f(someUser)  // Pass receiver as first argument

6. SHADOWING EMBEDDED METHOD ACCIDENTALLY
   --------------------------------------
   type Inner struct{}
   func (Inner) Process() { fmt.Println("Inner") }

   type Outer struct {
       Inner
   }
   func (Outer) Process() { fmt.Println("Outer") }

   o := Outer{}
   o.Process()       // Prints "Outer" (shadowed)
   o.Inner.Process() // Prints "Inner" (explicit)

===============================================================================
BEST PRACTICES
===============================================================================

1. RECEIVER NAMING
   - Use short, consistent names (1-2 letters)
   - Use the first letter(s) of the type name
   - WRONG: func (this *User), func (self *User), func (user *User)
   - CORRECT: func (u *User)

2. RECEIVER TYPE CHOICE
   Use POINTER receiver when:
   - Method modifies the receiver
   - Struct is large (>40 bytes as rule of thumb)
   - Type contains sync.Mutex or similar (must not be copied)
   - Consistency (if any method needs pointer, use for all)

   Use VALUE receiver when:
   - Type is small immutable type (like time.Time)
   - Type is a map, func, or chan (already reference types)
   - Method doesn't modify receiver

3. NIL SAFETY
   - Always handle nil receivers in pointer methods
   - Document whether nil is valid
   - Consider returning zero values for nil receivers

4. CONSTRUCTOR FUNCTIONS
   - Use NewTypeName() for complex initialization
   - Return pointer for mutable types
   - Validate inputs in constructor

5. METHOD GROUPING
   - Group methods by receiver type
   - Place constructors before methods
   - Order: New/Make, Getters, Setters, Actions, Helpers

6. DOCUMENTATION
   - Document the receiver in method comments
   - Note if method is nil-safe
   - Document if method has side effects

===============================================================================
PERFORMANCE NOTES
===============================================================================

1. Value receiver copies the entire struct on each call
   - For large structs (>64 bytes), use pointer receivers
   - Copying small structs is often faster than pointer indirection

2. Pointer receivers involve one level of indirection
   - May cause cache misses for very hot code paths
   - But allows mutation without copies

3. Method calls have zero overhead compared to functions
   - Go compiles methods to regular function calls
   - The receiver becomes the first parameter

4. Interface method calls have overhead
   - Involves dynamic dispatch (vtable lookup)
   - Discussed in interfaces section

===============================================================================
*/
