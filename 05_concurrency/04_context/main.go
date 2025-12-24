package main

import (
	"context"
	"fmt"
	"time"
)

// DEEP DIVE: Context
// Context carries deadlines, cancellation signals, and request-scoped values.
// Rule 1: Always pass context as the first argument to functions (ctx context.Context).
// Rule 1: Always pass context as the first argument to functions (ctx context.Context).
// Rule 2: Never store context in a struct; pass it explicitly.
//
// VISUALIZATION (Cancellation Propagation):
// [ Background ] (Root)
//       |
//       +-> [ WithCancel (ctx1) ] --+-- [ Worker A ]
//       |      ^ Cancel triggers    |
//       |                           +-- [ Worker B ]
//       |
//       +-> [ WithTimeout (ctx2) ] ---> [ DB Query ] (Cancels on timeout)

func main() {
	// 1. WithCancel (Manually trigger cancellation)
	// 'Background' is the root of the context tree.
	ctx, cancel := context.WithCancel(context.Background())
	
	go worker(ctx, "Worker 1")

	time.Sleep(1 * time.Second)
	fmt.Println("Main: Cancelling Worker 1...")
	cancel() // Sends signal to Done() channel
	time.Sleep(500 * time.Millisecond)

	// 2. WithTimeout (Deadlines)
	// Cancels automatically after duration.
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelTimeout() // IMPORTANT: Always defer cancel to release resources!

	slowOperation(ctxTimeout)

	// 3. WithValue (Request Scoped Data)
	// Use loosely typed keys. Prone to type errors, use sparingly (e.g., TraceIDs, Users).
	key := "userID"
	ctxVal := context.WithValue(context.Background(), key, 42)
	processRequest(ctxVal)
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done(): // Calls when cancel() is called OR parent cancels.
			fmt.Printf("%s stopping: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s working...\n", name)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func slowOperation(ctx context.Context) {
	select {
	case <-time.After(3 * time.Second): // Takes 3s
		fmt.Println("Operation finished (success)")
	case <-ctx.Done(): // Timeout triggers after 2s
		fmt.Println("Operation aborted:", ctx.Err()) // context deadline exceeded
	}
}

func processRequest(ctx context.Context) {
	val := ctx.Value("userID")
	if id, ok := val.(int); ok {
		fmt.Println("Processing for UserID:", id)
	} else {
		fmt.Println("No user found")
	}
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Context Leak:
   Warning: Forgetting 'defer cancel()'.
   Why: The parent context retains a reference to the child until cancelled/timed out. Leaks memory.
   Fix: Always call cancel() (usually via defer) for WithCancel/Timeout/Deadline.

2. Passing Nil Context:
   Error: "panic: runtime error: invalid memory address..." inside libraries.
   Fix: Never pass 'nil'. Use 'context.TODO()' if you are unsure, or 'context.Background()' at root.

3. Context Values Abuse:
   Anti-Pattern: Passing optional function parameters via Context.
   Why: Loss of type safety. Hidden dependencies.
   Rule: Only use for request-scoped data (TraceID, AuthToken) that traverses APIs.
*/
