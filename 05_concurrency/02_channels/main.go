package main

import (
	"fmt"
	"time"
)

// DEEP DIVE: Channels
// Channels are typed conduits for synchronization and messaging.
// "Don't communicate by sharing memory; share memory by communicating."

func main() {
	// 1. Unbuffered Channel (Synchronous)
	// Sends BLOCK until a receiver is ready.
	// Receives BLOCK until a sender is ready.
	ch := make(chan string)

	go func() {
		fmt.Println("Sender: Sending...")
		ch <- "Ping" // BLOCKS here until main receives
		fmt.Println("Sender: Sent!")
	}()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Receiver: Waiting...")
	msg := <-ch // BLOCKS here until sender sends
	fmt.Println("Receiver: Got", msg)

	// 2. Buffered Channel (Asynchronous)
	// Doesn't block send unless buffer is FULL.
	// Doesn't block receive unless buffer is EMPTY.
	bufCh := make(chan int, 2) // Capacity 2
	bufCh <- 1 // Won't block
	bufCh <- 2 // Won't block
	// bufCh <- 3 // WOULD BLOCK (Buffer full)
	
	fmt.Println("Buffered read:", <-bufCh)
	fmt.Println("Buffered read:", <-bufCh)

	// 3. Closing & Range
	queue := make(chan int, 3)
	queue <- 10
	queue <- 20
	close(queue) // Signals "no more data".
	
	// Panic: send on closed channel
	// queue <- 30 // PANIC!

	// Reading closed channel works: it returns zero value and false.
	val, ok := <-queue
	fmt.Printf("Read: %d, Open: %v\n", val, ok) // 10, true
	val, ok = <-queue
	fmt.Printf("Read: %d, Open: %v\n", val, ok) // 20, true
	val, ok = <-queue
	fmt.Printf("Read: %d, Open: %v\n", val, ok) // 0, FALSE (Closed & Empty)

	// Iterating
	queue2 := make(chan int, 3)
	queue2 <- 1
	queue2 <- 2
	close(queue2)
	
	fmt.Println("--- Range ---")
	for v := range queue2 { // Automatically stops when closed
		fmt.Println(v)
	}

	// 4. Nil Selection Blocking
	// A nil channel BLOCKS FOREVER on both send and receive.
	// Useful in 'select' statements to disable cases.
	var nilCh chan int // nil
	_ = nilCh
	// <-nilCh // Deadlock!
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Send on Closed Channel:
   Error: "panic: send on closed channel"
   Fix: Only the SENDER should close the channel. Never close from the receiver side.
   Pattern: Use a 'done' channel to signal workers to stop, but let the main coordinator close data channels.

2. Nil Channel Deadlock:
   Error: Blocks forever.
   Pitfall: 'var ch chan int' (nil). Sending/Receiving hangs.
   Fix: Always generic with 'make(chan int)'.

3. Unbuffered Deadlock:
   Error: 'ch := make(chan int); ch <- 1' in main thread blocks forever because no ONE is reading yet.
   Fix: Use a buffer OR start a goroutine to read BEFORE sending.
*/
