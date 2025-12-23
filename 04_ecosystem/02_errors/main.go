package main

import (
	"errors"
	"fmt"
	"os"
)

// DEEP DIVE: Error Handling
// Go treats errors as values, not exceptions.

// 1. Sentinel Errors (Pre-defined errors)
// Convention: var ErrSomeName = ...
var ErrNotFound = errors.New("item not found")
var ErrInvalid = errors.New("invalid input")

// 2. Custom Error Types (Structure)
type FileError struct {
	Path string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("file error at %s: %v", e.Path, e.Err)
}

// Ensure it verifies the 'Unwrap' interface
func (e *FileError) Unwrap() error {
	return e.Err
}

func openConfig() error {
	// 3. Wrapping Errors (Go 1.13+)
	// Use %w verb to wrap an error inside another.
	_, err := os.Open("missing_config.json")
	if err != nil {
		// We return a Custom Error wrapping the original os.PathError
		return &FileError{Path: "missing_config.json", Err: err}
	}
	return nil
}

func main() {
	err := openConfig()

	if err != nil {
		fmt.Println("Error occurred:", err)

		// 4. Checking Sentinel Errors (errors.Is)
		// Checks if the chain of errors contains a specific error.
		// Even if wrapped!
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(">> Detected wrapped 'not exist' error!")
		}

		// 5. Asserting Types (errors.As)
		// Checks if the error chain contains a specific TYPE.
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			fmt.Println(">> This is an os.PathError!")
			fmt.Println(">> Op:", pathErr.Op)
			fmt.Println(">> Path:", pathErr.Path)
		}
	}
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Comparing Wrapped Errors:
   Error: 'if err == os.ErrNotExist' fails if error is wrapped.
   Why: Wrapping adds a layer. Equality checks the outer layer.
   Fix: Always use 'errors.Is(err, target)' for sentinels.

2. Custom Error Pointer Equality:
   Pitfall: A custom error pointer 'var e *MyErr = nil' passed as 'error' interface is NOT nil (Typed Nil).
   Fix: Return 'error' interface explicitly as nil from function.
*/
