package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func launchServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(500 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Server is running ninja")
		case <-ctx.Done():
			fmt.Println("Server Shutting down")
			time.Sleep(2 * time.Second)
			fmt.Println("Cleanup completed. Server exited gracefully.")
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go launchServer(ctx, &wg)

	fmt.Println("Do ctrl + C to exit")

	<-exit

	fmt.Println("Exit signal received shutting down the server")

	cancel()

	wg.Wait()

	fmt.Println("Graceful shut down of server completed")
}
