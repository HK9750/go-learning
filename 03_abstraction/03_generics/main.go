package main

import "fmt"

// DEEP DIVE: Generics (Go 1.18+)
// Parametric polymorphism.

// 1. Generic Function
// [T any] is the Type Parameter List.
// 'comparable' is a built-in constraint (supports == and !=).
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// 2. Custom Constraints (Type Sets)
type Number interface {
	int64 | float64 // Union of types
}

func Sum[V Number](m []V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// 3. Generic Types
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	n := len(s.items)
	item := s.items[n-1]
	s.items = s.items[:n-1]
	return item
}

func main() {
	// A. Generic Functions
	si := []int{10, 20, 15, -10}
	fmt.Println("Index of 15:", Index(si, 15))

	ss := []string{"foo", "bar", "baz"}
	fmt.Println("Index of bar:", Index(ss, "bar"))

	// B. Generic Constraints
	floats := []float64{1.1, 2.2, 3.3}
	fmt.Println("Sum Floats:", Sum(floats))
	
	ints := []int64{100, 200, 300}
	fmt.Println("Sum Ints:", Sum(ints))

	// C. Generic Struct
	strStack := Stack[string]{}
	strStack.Push("Hello")
	strStack.Push("Generics")
	fmt.Println("Popped:", strStack.Pop())
	fmt.Println("Popped:", strStack.Pop())
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Generic Methods Limitation:
   Error: "method must have no type parameters"
   Why: Go does not support generic methods on non-generic types (yet).
   Fix: Use a generic function instead or make the struct generic.

2. Constraints Satisfaction:
   Error: "T does not satisfy comparable"
   Why: Tries to use '==' on a type that doesn't support it (e.g., slice, map, func).
   Fix: Use strict constraints like 'comparable' or custom interfaces.
*/
