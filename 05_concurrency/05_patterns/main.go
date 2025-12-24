package main

import (
	"fmt"
	"time"
)

// DEEP DIVE: Concurrency Patterns
// 1. Worker Pools
// 2. Pipelines
// 3. Fan-in / Fan-out

func main() {
	// --- PATTERN 1: WORKER POOL ---
	// Distributes jobs N workers.
	//
	// VISUALIZATION:
	// [ Jobs Channel ] -> [ Worker 1 ] -> [ Results Channel ]
	//                  -> [ Worker 2 ] ->
	//                  -> [ Worker 3 ] ->
	
	const numJobs = 5
	const numWorkers = 3
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// A. Start Workers (Fan-Out)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// B. Send Jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Signal "no more jobs"

	// C. Collect Results (Fan-In/Wait)
	for a := 1; a <= numJobs; a++ {
		res := <-results
		fmt.Printf("Result: %d\n", res)
	}
	
	// --- PATTERN 2: GENERATOR / PIPELINE ---
	fmt.Println("\n--- Pipeline ---")
	// 1. Generator: Converts list of ints to a channel
	gen := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			for _, n := range nums {
				out <- n
			}
			close(out)
		}()
		return out
	}

	// 2. Stage 2: Squares the numbers
	sq := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			for n := range in {
				out <- n * n
			}
			close(out)
		}()
		return out
	}

	// Chaining stages
	for n := range sq(gen(2, 3, 4, 5)) {
		fmt.Printf("Pipeline output: %d\n", n)
	}

	// --- PATTERN 3: ERRGROUP (Conceptual) ---
	// "golang.org/x/sync/errgroup" is standard for managing groups of goroutines
	// where if ONE fails, ALL are cancelled (via Context).
	// We won't import external lib here but know it exists!
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(time.Millisecond * 200) // Simulate work
		results <- j * 2
	}
}

/*
ERROR ANALYSIS & BEST PRACTICES:

1. Pipeline Goroutine Leaks:
   Pitfall: If the consumer stops reading (e.g., error happens), the producer goroutines block forever on 'out <- n'.
   Fix: Pass a 'done' channel (or Context) to all generators so they can exit early.
   
   func gen(done <-chan struct{}, nums ...int) <-chan int {
       go func() {
           for _, n := range nums {
               select {
               case out <- n:
               case <-done: return // Exit!
               }
           }
           close(out)
       }()
       return out
   }
*/
