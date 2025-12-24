# Go Mastery Exercises

Complete these exercises in order. Create your own `main.go` files from scratch.

---

## Exercise 1: Robust CLI Calculator
**Topics**: `os.Args`, `strconv`, `switch`, Error Handling

### Requirements
1.  Create `exercises/01_calculator/main.go`.
2.  Run with: `go run main.go <operation> <val1> <val2>`
    - Example: `go run main.go add 10 5` → Output: `15`
3.  Supported operations: `add`, `sub`, `mul`, `div`.
4.  **Validation**:
    - If fewer than 3 arguments, print usage info and exit.
    - If `operation` is unknown, print an error.
    - If `val1` or `val2` are not valid numbers, print a specific error.
    - If dividing by zero, handle gracefully (don't panic).
5.  Print result as a float (e.g., `5 / 2 = 2.5`).

---

## Exercise 2: Log Analyzer
**Topics**: `structs`, `slices`, `maps`, `strings.Split`

### Requirements
1.  Create `exercises/02_log_analyzer/main.go`.
2.  Define a struct: `LogEntry { IP string, Path string, Code int }`.
3.  Hardcode this raw log data:
    ```
    "10.0.0.1,/home,200"
    "10.0.0.2,/about,404"
    "10.0.0.1,/contact,200"
    "10.0.0.3,/pricing,500"
    "10.0.0.2,/login,404"
    ```
4.  Parse each line into `LogEntry` structs.
5.  Use a `map[string]int` to count requests per IP.
6.  **Output**:
    - Total count of 404 errors.
    - The IP address with the most requests.

---

## Exercise 3: Interface-Based Notifier
**Topics**: `interfaces`, `polymorphism`, `methods`

### Requirements
1.  Create `exercises/03_notifier/main.go`.
2.  Define an interface: `Notifier` with method `Notify(message string) error`.
3.  Implement two structs that satisfy `Notifier`:
    - `EmailNotifier`: Prints `"[EMAIL] Sending: <message>"`
    - `SMSNotifier`: Prints `"[SMS] Sending: <message>"`
4.  Create a function `Broadcast(notifiers []Notifier, msg string)` that loops through the slice and calls `Notify` on each.
5.  In `main`, create a slice containing both notifier types and call `Broadcast` with a test message.

---

## Exercise 4: Concurrent URL Checker
**Topics**: `goroutines`, `channels`, `worker pools`, `sync.WaitGroup`

### Requirements
1.  Create `exercises/04_url_checker/main.go`.
2.  Define a struct: `Result { URL string, Status string }`.
3.  Hardcode a list of 10 fake URLs (e.g., `"http://site1.com"`, ...).
4.  Implement a **Worker Pool** pattern:
    - Start **3 worker goroutines**.
    - Create a `jobs` channel (for URLs) and a `results` channel (for Results).
5.  **Worker logic**:
    - Sleep for a random duration (100-400ms) to simulate network latency.
    - Randomly return `"200 OK"` or `"500 Error"` (50/50 chance).
6.  **Main logic**:
    - Send all URLs to the `jobs` channel, then close it.
    - Collect and print all results.

---

## Exercise 5: Graceful Shutdown Server
**Topics**: `context`, `os/signal`, `goroutines`, `time.Ticker`

### Requirements
1.  Create `exercises/05_graceful_shutdown/main.go`.
2.  Start a "fake server" goroutine that prints `"Server processing..."` every 500ms using a `time.Ticker`.
3.  Listen for `SIGINT` (Ctrl+C) using `signal.Notify`.
4.  When the signal is received:
    - Cancel the context.
    - The worker goroutine should detect `<-ctx.Done()` and print `"Server shutting down..."`.
    - Wait 1 second for "cleanup", then exit gracefully.

---

Good luck! 🚀
