# Part 1: Basics (Exercises 1-10)

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

## Exercise 6: JSON Config Loader
**Topics**: `encoding/json`, `os`, File I/O, Error Wrapping

### Requirements
1.  Create `exercises/06_json_config/main.go` and a `config.json` file.
2.  Define a struct:
    ```go
    type Config struct {
        AppName     string   `json:"app_name"`
        Port        int      `json:"port"`
        Debug       bool     `json:"debug"`
        AllowedIPs  []string `json:"allowed_ips"`
    }
    ```
3.  Create `config.json` with sample data.
4.  Write a function `LoadConfig(path string) (*Config, error)` that:
    - Opens and reads the file.
    - Unmarshals JSON into the struct.
    - Uses `fmt.Errorf` with `%w` for error wrapping.
5.  Validate: Port must be between 1-65535, AppName cannot be empty.
6.  Print the loaded config in a formatted way.

---

## Exercise 7: Custom Error Types
**Topics**: `error` interface, Custom Errors, `errors.Is`, `errors.As`, `fmt.Errorf`

### ❓ Why is this necessary?
In production systems, simply returning string errors (`errors.New("failed")`) is insufficient. You need to distinguish between **client errors** (400 Bad Request), **server errors** (500 Internal Server Error), and **not found errors** (404). Custom error types allow layers of your application to inspect *why* something failed without parsing error strings, leading to robust robust retry logic and meaningful API responses.

### 🌍 Real-World Application
*   **Kubernetes API**: Uses custom error types to signaling if a resource exists, is conflict, or if the user is unauthorized.
*   **Database Drivers**: Return specific error types for connection timeouts vs constraint violations, allowing the application to decide whether to retry.

### 📦 Production Requirements
1.  Create `exercises/07_custom_errors/main.go`.
2.  Define custom error structs:
    ```go
    type ValidationError struct {
        Field   string
        Message string
    }
    
    type NotFoundError struct {
        Resource string
        ID       int
    }
    ```
3.  Implement the `Error() string` method for both.
4.  Create a function `GetUser(id int) (string, error)` that:
    - Returns `NotFoundError` wrapped with context if `id < 0`.
    - Returns `ValidationError` wrapped with context if `id == 0`.
    - Returns a fake username otherwise.
5.  **Use `errors.As`** in `main` to detect the error type.
6.  **Structured Logging**: Print different exit codes or logs depending on the error type (e.g., "User error" vs "System failure").

---

## Exercise 8: Production-Ready HTTP REST API
**Topics**: `net/http`, structured json, context, middleware

### ❓ Why is this necessary?
Writing a robust API involves more than just request matching. You must handle JSON parsing safely, protect shared state, manage request deadlines, and ensure your API behaves predictably under load. This exercise simulates the core of any microservice backend.

### 🌍 Real-World Application
*   **User Management Service**: Handling CRUD operations for user profiles.
*   **Inventory System**: Tracking stock levels with high concurrent reads/writes.

### 📦 Production Requirements
1.  Create `exercises/08_production_api/main.go`.
2.  **State Management**: Use `sync.RWMutex` to protect a `map[int]Task`.
3.  **Handlers**:
    - `GET /tasks`: Return JSON list.
    - `POST /tasks`: Validate input (Title cannot be empty). Return `400` if invalid.
    - `DELETE /tasks/{id}`: handling non-existent IDs with `404`.
4.  **JSON Best Practices**: Use `json.NewEncoder(w).Encode(v)` and check errors.
5.  **Graceful Error Handling**: Don't just `http.Error` strings; returns a JSON error object: `{"error": "message"}`.
6.  **Concurrency Safety**: Launch a goroutine that writes to the map every 100ms while running `k6` or `hey` to stress test reads.

---

## Exercise 9: Middleware Chain & Interceptors
**Topics**: `http.Handler`, Middleware Pattern, Decorators, Context Value

### ❓ Why is this necessary?
Middleware allows you to decouple cross-cutting concerns (logging, auth, metrics) from business logic. In production, you never write these inside every handler. You write them once and chain them. This pattern is the backbone of frameworks like Gin, Chi, and Echo.

### 🌍 Real-World Application
*   **Auth Gateways**: Verifying JWT tokens before passing requests to services.
*   **Observability agents**: Automatically timing every request and sending stats to Prometheus/Datadog.

### 📦 Production Requirements
1.  Create `exercises/09_middleware/main.go`.
2.  Implement standardized middleware:
    - **RequestID**: Generates a UUID and adds it to `request.Context()`.
    - **Logger**: Logs method, path, duration, and *RequestID*.
    - **Recovery**: Catches `panic()` to prevent server crash, logs the stack trace, and returns 500.
3.  **Chain Helper**: Write a function `Chain(h http.Handler, m ...Middleware) http.Handler` to compose them cleanly.
4.  **Demonstration**: Create a handler that simulates a panic, and verify the server stays alive and logs the partial stack.

---

## Exercise 10: Distributed-Ready Rate Limiter
**Topics**: `time.Ticker`, `sync.Mutex`, Token Bucket, Refill Strategies

### ❓ Why is this necessary?
Without rate limiting, a single bad actor or bug in a client script can bring down your entire platform (DDoS). Rate limiters protect your resources and ensure fair usage. Implementing one from scratch teaches you about time, locking, and atomic state management.

### 🌍 Real-World Application
*   **Public APIs**: GitHub/Twitter API limits (e.g., 5000 requests/hour).
*   **Login Endpoints**: Preventing brute-force password attacks.

### 📦 Production Requirements
1.  Create `exercises/10_rate_limiter/main.go`.
2.  Implement a **Token Bucket** struct.
3.  **Configuration**: Make `MaxTokens` and `RefillRate` configurable.
4.  **Per-IP Limiting**: Maintain a map of `map[string]*RateLimiter`.
5.  **Cleanup**: Run a background goroutine to clean up limiters for old IPs (Memory leak prevention).
6.  **HTTP Integration**: Middleware that checks IP, calls `Allow()`, and sets `X-RateLimit-Remaining` headers.
7.  Return `429 Too Many Requests` if limit exceeded.
