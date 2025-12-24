package main

import "fmt"

// DEEP DIVE: Methods
// A method is just a function with a special "receiver" argument.
// Receiver can be (T) or (*T).
//
// VISUALIZATION (Method Sets):
// +----------------+      +-------------------------+
// |  Variable u    |      |  Method Set             |
// +----------------+      +-------------------------+
// |  u of type T   | ---> |  Receivers (t T)        |
// +----------------+      +-------------------------+
// |  u of type *T  | ---> |  Receivers (t T)        |
// |                |      |  Receivers (t *T)       |
// +----------------+      +-------------------------+

type User struct {
	Name string
	Logins int
}

// Value Receiver: 'u' is a COPY of the caller.
// Mutations here DO NOT affect the original.
func (u User) Display() {
	fmt.Printf("User: %s, Logins: %d\n", u.Name, u.Logins)
}

// Pointer Receiver: 'u' is a POINTER to the caller.
// Mutations here AFFECT the original.
// Prefer this for large structs (avoids copy) or when mutation is needed.
func (u *User) IncrementLogin() {
	if u == nil {
		fmt.Println("Warning: Nil receiver called!")
		return
	}
	u.Logins++
}

func main() {
	u := User{Name: "Gopher", Logins: 0}

	// 1. Value Receiver Call
	u.Display()

	// 2. Pointer Receiver Call (Go automatically takes address: (&u).IncrementLogin())
	u.IncrementLogin()
	u.Display() // Logins should be 1

	// 3. Method Expressions (Functional Style)
	// You can assign methods to variables.
	
	// Value receiver method expression
	// Signature: func(User)
	displayFunc := User.Display 
	displayFunc(u) // Pass 'u' explicitly

	// Pointer receiver method expression
	// Signature: func(*User)
	incFunc := (*User).IncrementLogin
	incFunc(&u) // Pass '&u' explicitly
	u.Display() // Logins: 2

	// 4. Nil Receivers
	// Methods CAN be called on nil pointers!
	var nilUser *User
	nilUser.IncrementLogin() // Safe, because we handled 'if u == nil' inside.
	// This is unlike method calls in Java/C++ which would NPE/Segfault immediately.
	// This is unlike method calls in Java/C++ which would NPE/Segfault immediately.
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Modifying Value Receiver:
   Ineffective: 'func (u User) update() { u.Logins++ }'
   Result: Updates the COPY of 'u'. Original object unchanged.
   Fix: Use pointer receiver 'func (u *User) update()'.

2. Nil Pointer Panic:
   Error: "panic: runtime error: invalid memory address or nil pointer dereference"
   Code: If receiver function body accesses field 'u.Logins' without checking 'if u == nil'.
   Fix: Guard clause 'if u == nil { return }'.

3. Mixing Receiver Types:
   Best Practice: If ONE method needs a pointer, ALL methods should use pointers for consistency (and to satisfy interfaces correctly).
*/
