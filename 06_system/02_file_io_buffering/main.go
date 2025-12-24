package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// =========================================================================
// SYSTEM: Buffered I/O (bufio)
// =========================================================================
// Loading large files (e.g., 1GB configs, logs) into memory (os.ReadFile) crashes apps.
// Solution: Stream it using bufio.Scanner or bufio.Reader.

const filename = "large_test.log"
const linesToWrite = 100000 // 100k lines

func main() {
	defer os.Remove(filename) // Cleanup

	// 1. Generate a large file EFFICIENTLY (Buffered Write)
	// If we used fmt.Fprintf in a loop directly to file, we'd do 100k syscalls.
	// bufio.Writer aggregates writes into 4KB chunks.
	fmt.Println("Generating test file...")
	start := time.Now()
	
	f, _ := os.Create(filename)
	bw := bufio.NewWriter(f) // 4KB buffer default
	
	for i := 0; i < linesToWrite; i++ {
		// Just writing random data
		bw.WriteString("INFO: This is a log line number ")
		bw.WriteString(fmt.Sprint(i))
		bw.WriteString(" with some extra payload data to make it longer.\n")
	}
	bw.Flush() // CRITICAL: Don't forget to flush the last chunk!
	f.Close()
	
	fmt.Printf("File generated in %v. Size: %.2f MB\n", 
		time.Since(start), 
		getFileSize(filename))

	// 2. Bad Way: os.ReadFile
	// fmt.Println("Reading all at once (Memory Hog)...")
	// data, _ := os.ReadFile(filename) // Allocates MBs immediately
	// _ = data

	// 3. Good Way: bufio.Scanner (Streaming)
	fmt.Println("\nStreaming file line-by-line...")
	start = time.Now()
	
	f2, _ := os.Open(filename)
	defer f2.Close()
	
	scanner := bufio.NewScanner(f2)
	// scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // Increase buffer if lines > 64KB
	
	lineCount := 0
	byteCount := 0
	
	for scanner.Scan() {
		line := scanner.Text() // Allocates string. Use scanner.Bytes() to avoid alloc if parsing.
		lineCount++
		byteCount += len(line)
		
		if lineCount % 20000 == 0 {
			fmt.Printf("   Processed %d lines...\n", lineCount)
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Processed %d lines (%d strings) in %v.\n", lineCount, byteCount, time.Since(start))
	fmt.Println("Memory usage stays constant regardless of file size!")
}

func getFileSize(name string) float64 {
	fi, _ := os.Stat(name)
	return float64(fi.Size()) / 1024 / 1024
}
