package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ==========================================================================
// SYSTEM: Graceful Shutdown
// =========================================================================
// Production apps must handle SIGINT (Ctrl+C) and SIGTERM (Kubernetes stops).
// We use a buffered channel to catch these signals.
//
// VISUALIZATION (Graceful Flow):
// [ OS Signal (SIGINT) ] -> [ Go Runtime ] -> [ sigChan ] (Buffered)
//                                                  |
//                                                  v
// [ Main Goroutine ] Unblocks! --------------------+
//       |
//       +--> Call cancel() -> [ Context Canceled ]
//                                     |
//       +--> [ Workers ] DETECT <-Done() -- stop processing
//       |
//       +--> cleanup() -> [ Exit 0 ]

func main() {
	// 1. Create a channel to listen for signals
	// We use a buffered channel of size 1 to prevent blocking the signal sender
	sigChan := make(chan os.Signal, 1)

	// 2. Register for SIGINT and SIGTERM
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("App Started. Process ID:", os.Getpid())
	fmt.Println("Waiting for SIGINT (Ctrl+C)...")

	// 3. Simulate work (Goroutine)
	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cleanup

	go worker(ctx)

	// 4. AUTO-TEST: Send Signal to Self after 2 seconds
	// In a real app, YOU would press Ctrl+C. Here we simulate it.
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n[Auto-Test] Sending SIGINT to self...")
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	// 5. Block until signal received
	receivedSig := <-sigChan
	fmt.Printf("\nSignal Received: %s\n", receivedSig)
	fmt.Println("Starting Graceful Shutdown...")

	// 6. Signal cleanup
	cancel() // Tell workers to stop

	// Simulate cleanup time (e.g., closing DB connections)
	time.Sleep(1 * time.Second)
	fmt.Println("Cleanup Complete. Exiting.")
}

func worker(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("   [Worker] Context cancelled, stopping work.")
			return
		case t := <-ticker.C:
			fmt.Println("   [Worker] Working...", t.Format("15:04:05"))
		}
	}
}
