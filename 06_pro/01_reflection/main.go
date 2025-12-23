package main

import (
	"fmt"
	"reflect"
)

// DEEP DIVE: Reflection
// Ability to inspect types and values at runtime.
// Use sparingly! It's slow and bypasses compile-time safety.
// Used by: JSON serialization, ORMs, fmt.Printf.

type Config struct {
	IP   string `env:"IP_ADDR"`
	Port int    `env:"PORT_NUM"`
}

func main() {
	// 1. TypeOf and ValueOf
	var x float64 = 3.14
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	fmt.Println("Type:", t)
	fmt.Println("Value:", v)
	fmt.Println("Kind:", t.Kind()) // float64

	// 2. Modifying Values (Setters)
	// You can only modify "Settable" values.
	// A value is settable if it's an addressable pointer.
	
	// v.SetFloat(7.1) // PANIC: reflect: reflect.Value.SetFloat using unaddressable value

	p := reflect.ValueOf(&x) // Pointer to x
	elem := p.Elem()         // Dereference
	if elem.CanSet() {
		elem.SetFloat(7.1)
		fmt.Println("New Value (via reflect):", x)
	}

	// 3. Struct Introspection
	cfg := Config{IP: "127.0.0.1", Port: 8080}
	inspectStruct(cfg)
}

func inspectStruct(s interface{}) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct")
		return
	}

	fmt.Println("\n--- Struct Inspection ---")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)      // Field metadata
		value := v.Field(i)      // Field value
		tag := field.Tag.Get("env")

		fmt.Printf("Field: %s, Type: %s, Value: %v, Tag: %s\n",
			field.Name, field.Type, value, tag)
	}
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Panic on Unexported Fields:
   Error: "reflect: reflect.Value.Interface: cannot return value obtained from unexported field"
   Why: You tried to read keys of a private field (lowercase) via reflection.
   Fix: Reflection ignores unexported fields unless you use unsafe (dangerous). Use exported fields.

2. Panic on Wrong Type:
   Error: "panic: reflect: Call using mismatching types"
   Why: Calling SetInt on a String field, etc.
   Fix: Always check Kind() before modifying.

3. Performance Cost:
   Fact: Reflection is 10-100x slower than direct access.
   Rule: Use only for specific needs (serialization, generic tools), never main business logic loops.
*/
