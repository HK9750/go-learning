package main

import "fmt"

// =============================================================================
// WHAT YOU DID: Interface-Based Polymorphism Pattern
// =============================================================================
//
// You implemented the "Strategy Pattern" using Go interfaces. This allows
// different notification strategies (Email, SMS) to be used interchangeably.
//
// =============================================================================
// WHY YOU DID IT:
// =============================================================================
// 1. LOOSE COUPLING: The Broadcast function doesn't care HOW notifications
//    are sent. It only knows that something can "Notify". This means you can
//    add PushNotifier, SlackNotifier, etc. without changing Broadcast.
//
// 2. TESTABILITY: You can create a MockNotifier for testing without sending
//    real emails/SMS.
//
// 3. OPEN/CLOSED PRINCIPLE: Your code is "open for extension" (add new
//    notifiers) but "closed for modification" (don't change existing code).
//
// =============================================================================
// WHAT HAPPENS AT SYSTEM LEVEL:
// =============================================================================
//
// 1. INTERFACE TYPE (Notifier):
//    - At runtime, Go stores interfaces as a pair: (type, value)
//    - The "type" part contains a pointer to type metadata (called itab)
//    - The "value" part contains the actual data (pointer to EmailNotifier/SmsNotifier)
//
// 2. METHOD DISPATCH (Dynamic Dispatch):
//    - When you call n.Notify(), Go looks up the method in the interface's
//      method table (similar to a vtable in C++)
//    - This is slightly slower than direct function calls (1-2 nanoseconds)
//    - The CPU does: dereference interface → find method pointer → call method
//
// 3. MEMORY LAYOUT:
//    - Notifier interface = 16 bytes (2 pointers on 64-bit: itab + data)
//    - []Notifier slice = 24 bytes (pointer + length + capacity)
//    - Each struct is allocated on heap because we take its address (&)
//
// 4. IMPLICIT INTERFACE SATISFACTION:
//    - Go uses "structural typing" - no "implements" keyword needed
//    - At compile time, Go verifies EmailNotifier has Notify() method
//    - This is checked during compilation, not runtime
//
// =============================================================================

// Notifier defines a contract: any type with a Notify method is a Notifier.
// This is Go's way of achieving polymorphism without inheritance.
// System: This creates an interface type descriptor stored in the binary.
type Notifier interface {
	Notify(message string) error
}

// EmailNotifier is a concrete implementation of Notifier.
// Empty struct = 0 bytes (Go optimizes this).
type EmailNotifier struct{}

// This method has a POINTER RECEIVER (*EmailNotifier).
// System: The method is registered in the type's method set.
// Using pointer receiver because:
// - Convention for consistency
// - Would be required if struct had fields to modify
func (e *EmailNotifier) Notify(message string) error {
	fmt.Println("Email Notifier message: ", message)
	return nil
}

type SmsNotifier struct{}

func (s *SmsNotifier) Notify(message string) error {
	fmt.Println("Sms Notifier message", message)
	return nil
}

// Broadcast accepts ANY type that satisfies Notifier interface.
// This is polymorphism - one function, multiple behaviors.
// System: Each element in 'notifiers' is an interface value (itab + data pair).
// The loop performs dynamic dispatch for each call.
func Broadcast(notifiers []Notifier, message string) {
	for _, n := range notifiers {
		n.Notify(message) // Dynamic dispatch happens here
	}
}

func main() {
	// & creates a pointer to the struct (allocates on heap due to escape analysis)
	// System: Go's escape analysis determines these escape to heap because
	// they're stored in an interface (which could outlive this scope)
	email := &EmailNotifier{}
	sms := &SmsNotifier{}

	// Type conversion: *EmailNotifier and *SmsNotifier → Notifier interface
	// System: Go creates interface values wrapping each pointer
	Broadcast([]Notifier{email, sms}, "Hello its me")
}

// =============================================================================
// JAVASCRIPT EQUIVALENT:
// =============================================================================
//
// JavaScript doesn't have interfaces, but we can achieve similar behavior
// using "duck typing" (if it walks like a duck and quacks like a duck...).
//
// ```javascript
// // In JS, we rely on convention rather than compile-time checks
//
// class EmailNotifier {
//     notify(message) {
//         console.log("Email Notifier message:", message);
//         return null; // equivalent to nil error
//     }
// }
//
// class SmsNotifier {
//     notify(message) {
//         console.log("SMS Notifier message:", message);
//         return null;
//     }
// }
//
// // Broadcast works with ANY object that has a notify() method
// // This is duck typing - no interface declaration needed
// function broadcast(notifiers, message) {
//     for (const notifier of notifiers) {
//         notifier.notify(message); // Runtime error if notify doesn't exist
//     }
// }
//
// // Usage
// const email = new EmailNotifier();
// const sms = new SmsNotifier();
//
// broadcast([email, sms], "Hello its me");
// ```
//
// KEY DIFFERENCES:
// -----------------------------------------------------------------------------
// | Aspect              | Go                      | JavaScript              |
// |---------------------|-------------------------|-------------------------|
// | Type checking       | Compile-time            | Runtime                 |
// | Interface keyword   | Explicit                | None (duck typing)      |
// | Error if wrong type | Compile error           | Runtime TypeError       |
// | Performance         | Optimized vtable lookup | Prototype chain lookup  |
// | Memory overhead     | 16 bytes per interface  | Object overhead varies  |
// -----------------------------------------------------------------------------
//
// =============================================================================
