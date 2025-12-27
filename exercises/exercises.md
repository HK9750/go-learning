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
**Topics**: `error` interface, Custom Errors, `errors.Is`, `errors.As`

### Requirements
1.  Create `exercises/07_custom_errors/main.go`.
2.  Define custom error types:
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
    - Returns `NotFoundError` if `id < 0`.
    - Returns `ValidationError` if `id == 0`.
    - Returns a fake username otherwise.
5.  In `main`, use `errors.As` to check error types and handle each differently.

---

## Exercise 8: HTTP REST API
**Topics**: `net/http`, JSON responses, Routing, HTTP Methods

### Requirements
1.  Create `exercises/08_http_api/main.go`.
2.  Implement an in-memory "task" store with struct:
    ```go
    type Task struct {
        ID        int    `json:"id"`
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
    }
    ```
3.  Implement these endpoints:
    - `GET /tasks` - Return all tasks as JSON.
    - `GET /tasks/{id}` - Return a single task (parse ID from path).
    - `POST /tasks` - Create a new task from JSON body.
    - `DELETE /tasks/{id}` - Delete a task by ID.
4.  Return proper status codes: `200`, `201`, `404`, `400`.
5.  Use `sync.RWMutex` for thread-safe access to the task store.

---

## Exercise 9: Middleware Chain
**Topics**: `http.Handler`, Middleware Pattern, Function Closures

### Requirements
1.  Create `exercises/09_middleware/main.go`.
2.  Create an HTTP server with these middleware functions:
    - **Logger**: Logs request method, path, and duration.
    - **Recoverer**: Catches panics and returns 500 error.
    - **Auth**: Checks for `X-API-Key` header, returns 401 if missing.
3.  Each middleware should follow the pattern:
    ```go
    func Logger(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // ... before
            next.ServeHTTP(w, r)
            // ... after
        })
    }
    ```
4.  Chain them: `Logger(Recoverer(Auth(yourHandler)))`.
5.  Create a `/protected` endpoint that panics and verify recovery works.

---

## Exercise 10: Rate Limiter
**Topics**: `time.Ticker`, `sync.Mutex`, Token Bucket Algorithm

### Requirements
1.  Create `exercises/10_rate_limiter/main.go`.
2.  Implement a **Token Bucket** rate limiter:
    ```go
    type RateLimiter struct {
        tokens     int
        maxTokens  int
        refillRate time.Duration
        mu         sync.Mutex
    }
    ```
3.  Methods:
    - `NewRateLimiter(max int, refill time.Duration) *RateLimiter`
    - `Allow() bool` - Returns true if request allowed, consumes a token.
    - Start a background goroutine that refills tokens.
4.  Create an HTTP server with an endpoint `/api` that uses the rate limiter.
5.  If rate limited, return `429 Too Many Requests`.
6.  Test by sending rapid requests.

---

## Exercise 11: Database CRUD with SQLite
**Topics**: `database/sql`, SQLite, Prepared Statements, Transactions

### Requirements
1.  Create `exercises/11_sqlite_crud/main.go`.
2.  Install driver: `go get github.com/mattn/go-sqlite3`
3.  Create a `users` table: `id`, `name`, `email`, `created_at`.
4.  Implement a `UserRepository` with methods:
    - `Create(name, email string) (int64, error)`
    - `GetByID(id int64) (*User, error)`
    - `Update(id int64, name, email string) error`
    - `Delete(id int64) error`
    - `List() ([]User, error)`
5.  Use **prepared statements** for all queries.
6.  Implement a transaction that creates two users atomically (rollback if either fails).

---

## Exercise 12: CLI with Cobra
**Topics**: `github.com/spf13/cobra`, CLI Frameworks, Subcommands

### Requirements
1.  Create `exercises/12_cobra_cli/`.
2.  Install: `go get github.com/spf13/cobra`
3.  Build a `task` CLI tool with subcommands:
    - `task add "Buy groceries"` - Add a task.
    - `task list` - Show all tasks.
    - `task complete <id>` - Mark task as done.
    - `task delete <id>` - Remove a task.
4.  Store tasks in a local JSON file (`~/.tasks.json`).
5.  Add flags: `--all` for list (show completed), `--priority` for add.

---

## Exercise 13: WebSocket Chat Server
**Topics**: `gorilla/websocket`, Real-time Communication, Channels

### Requirements
1.  Create `exercises/13_websocket_chat/main.go`.
2.  Install: `go get github.com/gorilla/websocket`
3.  Implement a chat server:
    - Maintain a map of connected clients.
    - Broadcast messages to all connected clients.
    - Handle client connect/disconnect gracefully.
4.  Create a simple HTML page with JavaScript to connect and send messages.
5.  Each message should include: `username`, `message`, `timestamp`.

---

## Exercise 14: Pipeline Pattern
**Topics**: Fan-out/Fan-in, Pipeline Stages, Channel Orchestration

### Requirements
1.  Create `exercises/14_pipeline/main.go`.
2.  Build a data processing pipeline with stages:
    - **Generator**: Produces numbers 1-100.
    - **Filter**: Only passes even numbers.
    - **Square**: Squares each number.
    - **Aggregate**: Sums all results.
3.  Each stage should be a goroutine reading from input channel, writing to output.
4.  Implement **fan-out**: Use 3 workers for the Square stage.
5.  Implement **fan-in**: Merge results from workers into single channel.
6.  Print the final sum.

---

## Exercise 15: Context with Timeouts
**Topics**: `context.WithTimeout`, `context.WithDeadline`, Cancellation

### Requirements
1.  Create `exercises/15_context_timeout/main.go`.
2.  Simulate a slow external API call:
    ```go
    func slowAPI(ctx context.Context) (string, error) {
        select {
        case <-time.After(3 * time.Second):
            return "data from API", nil
        case <-ctx.Done():
            return "", ctx.Err()
        }
    }
    ```
3.  Create three scenarios:
    - Call with 5-second timeout (should succeed).
    - Call with 1-second timeout (should timeout).
    - Call with manual cancellation after 500ms.
4.  Properly handle `context.DeadlineExceeded` and `context.Canceled` errors.

---

## Exercise 16: Generics
**Topics**: Type Parameters, Constraints, Generic Data Structures

### Requirements
1.  Create `exercises/16_generics/main.go`.
2.  Implement generic functions:
    ```go
    func Map[T, U any](slice []T, fn func(T) U) []U
    func Filter[T any](slice []T, fn func(T) bool) []T
    func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U
    ```
3.  Implement a generic `Stack[T]` data structure with:
    - `Push(item T)`
    - `Pop() (T, bool)`
    - `Peek() (T, bool)`
    - `IsEmpty() bool`
4.  Test with different types: `int`, `string`, custom struct.

---

## Exercise 17: Reflection
**Topics**: `reflect` package, Runtime Type Inspection, Struct Tags

### Requirements
1.  Create `exercises/17_reflection/main.go`.
2.  Create a function `PrintStructInfo(v any)` that:
    - Prints struct name.
    - Lists all fields with: name, type, value, and tags.
3.  Create a `Validate` function using reflection:
    - Read `validate:"required"` tags.
    - Return error if required fields are zero-value.
4.  Create a `ToMap(v any) map[string]any` function that converts struct to map.
5.  Test with a complex nested struct.

---

## Exercise 18: Testing & Benchmarking
**Topics**: `testing`, Table-Driven Tests, Benchmarks, Mocks

### Requirements
1.  Create `exercises/18_testing/` with `calculator.go` and `calculator_test.go`.
2.  Implement a `Calculator` with `Add`, `Sub`, `Mul`, `Div` methods.
3.  Write **table-driven tests** for each operation:
    ```go
    func TestAdd(t *testing.T) {
        tests := []struct{
            name string
            a, b, want float64
        }{...}
    }
    ```
4.  Write **benchmarks** for each operation.
5.  Create an interface `DataFetcher` and mock it for testing.
6.  Achieve 100% code coverage.
7.  Run: `go test -v -cover -bench=.`

---

## Exercise 19: Embedded Structs & Composition
**Topics**: Embedding, Method Promotion, Interface Satisfaction

### Requirements
1.  Create `exercises/19_composition/main.go`.
2.  Build a permission system:
    ```go
    type Reader struct{}
    func (r Reader) Read() string { return "reading" }
    
    type Writer struct{}
    func (w Writer) Write(data string) string { return "writing: " + data }
    
    type Admin struct {
        Reader
        Writer
        Name string
    }
    ```
3.  Define interfaces: `Readable`, `Writable`, `ReadWriter`.
4.  Demonstrate that `Admin` satisfies `ReadWriter` via embedding.
5.  Create a function that accepts `ReadWriter` interface.
6.  Show method overriding: `Admin` can override `Read()`.

---

## Exercise 20: Worker Pool with Context
**Topics**: Advanced Concurrency, Context Propagation, Graceful Shutdown

### Requirements
1.  Create `exercises/20_worker_pool/main.go`.
2.  Build a robust worker pool:
    ```go
    type Pool struct {
        workers    int
        jobQueue   chan Job
        results    chan Result
        ctx        context.Context
        cancel     context.CancelFunc
        wg         sync.WaitGroup
    }
    ```
3.  Features:
    - Configurable number of workers.
    - Jobs can be cancelled via context.
    - Graceful shutdown: wait for in-flight jobs.
    - Error handling per job.
4.  Create a `Job` interface with `Execute(ctx context.Context) (any, error)`.
5.  Test with 100 jobs across 5 workers, cancel after 50 complete.

---

## Exercise 21: LRU Cache
**Topics**: `container/list`, Maps, Cache Eviction, Thread Safety

### Requirements
1.  Create `exercises/21_lru_cache/main.go`.
2.  Implement an LRU (Least Recently Used) cache:
    ```go
    type LRUCache struct {
        capacity int
        cache    map[string]*list.Element
        list     *list.List
        mu       sync.RWMutex
    }
    ```
3.  Methods:
    - `Get(key string) (any, bool)`
    - `Put(key string, value any)`
    - `Delete(key string)`
    - `Len() int`
4.  When capacity is exceeded, evict the least recently used item.
5.  Make it thread-safe with `sync.RWMutex`.
6.  Write tests to verify eviction order.

---

## Exercise 22: Event Bus / Pub-Sub
**Topics**: Observer Pattern, Channels, Goroutines, Thread Safety

### Requirements
1.  Create `exercises/22_event_bus/main.go`.
2.  Implement a publish-subscribe event system:
    ```go
    type EventBus struct {
        subscribers map[string][]chan Event
        mu          sync.RWMutex
    }
    
    type Event struct {
        Topic   string
        Payload any
    }
    ```
3.  Methods:
    - `Subscribe(topic string) <-chan Event`
    - `Unsubscribe(topic string, ch <-chan Event)`
    - `Publish(topic string, payload any)`
    - `Close()`
4.  Multiple subscribers per topic.
5.  Test: 3 subscribers on "orders", publish 5 events, verify all receive them.

---

## Exercise 23: Retry with Exponential Backoff
**Topics**: Error Handling, Time, Jitter, Circuit Breaker

### Requirements
1.  Create `exercises/23_retry/main.go`.
2.  Implement a retry mechanism:
    ```go
    type RetryConfig struct {
        MaxRetries  int
        BaseDelay   time.Duration
        MaxDelay    time.Duration
        Multiplier  float64
    }
    
    func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error
    ```
3.  Features:
    - Exponential backoff with jitter.
    - Respect context cancellation.
    - Different retry strategies for different error types.
4.  Add a simple **Circuit Breaker**:
    - Open after N consecutive failures.
    - Half-open after timeout, allow one request.
    - Close if request succeeds.
5.  Test with a flaky function that fails 70% of the time.

---

## Exercise 24: gRPC Service
**Topics**: Protocol Buffers, gRPC, Client-Server, Streaming

### Requirements
1.  Create `exercises/24_grpc/`.
2.  Install: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
3.  Define a `todo.proto`:
    ```protobuf
    service TodoService {
        rpc CreateTodo(CreateTodoRequest) returns (Todo);
        rpc ListTodos(Empty) returns (stream Todo);
        rpc WatchTodos(Empty) returns (stream Todo);
    }
    ```
4.  Generate Go code and implement:
    - Server with in-memory storage.
    - Client that creates and lists todos.
5.  Implement **server streaming** for `ListTodos`.
6.  Implement **real-time watching** for `WatchTodos`.

---

## Exercise 25: Build a Mini Framework
**Topics**: All Previous Topics Combined

### Requirements
1.  Create `exercises/25_mini_framework/`.
2.  Build a minimal web framework with:
    - **Router**: Path parameters, method-based routing.
    - **Middleware**: Logging, recovery, CORS.
    - **Context**: Request-scoped values, JSON helpers.
    - **Validation**: Struct tag-based validation.
3.  Example usage:
    ```go
    app := mini.New()
    app.Use(mini.Logger, mini.Recover)
    
    app.GET("/users/:id", func(c *mini.Context) {
        id := c.Param("id")
        c.JSON(200, map[string]string{"id": id})
    })
    
    app.Run(":8080")
    ```
4.  Support route groups and nested middleware.
5.  Write comprehensive tests.

---

## 🎯 Bonus Challenges

### Challenge A: Concurrent Map
Implement a thread-safe map without using `sync.Map`. Use sharding for better performance.

### Challenge B: Connection Pool
Build a database connection pool with: min/max connections, health checks, and timeout handling.

### Challenge C: Distributed Lock
Implement distributed locking using Redis or file-based locks.

### Challenge D: Log Aggregator
Build a service that collects logs from multiple sources, parses them, and provides search functionality.

---

## 📊 Progress Tracker

| # | Exercise | Status |
|---|----------|--------|
| 1 | CLI Calculator | ⬜ |
| 2 | Log Analyzer | ⬜ |
| 3 | Interface Notifier | ⬜ |
| 4 | URL Checker | ⬜ |
| 5 | Graceful Shutdown | ⬜ |
| 6 | JSON Config | ⬜ |
| 7 | Custom Errors | ⬜ |
| 8 | HTTP REST API | ⬜ |
| 9 | Middleware Chain | ⬜ |
| 10 | Rate Limiter | ⬜ |
| 11 | SQLite CRUD | ⬜ |
| 12 | Cobra CLI | ⬜ |
| 13 | WebSocket Chat | ⬜ |
| 14 | Pipeline Pattern | ⬜ |
| 15 | Context Timeout | ⬜ |
| 16 | Generics | ⬜ |
| 17 | Reflection | ⬜ |
| 18 | Testing | ⬜ |
| 19 | Composition | ⬜ |
| 20 | Worker Pool | ⬜ |
| 21 | LRU Cache | ⬜ |
| 22 | Event Bus | ⬜ |
| 23 | Retry/Backoff | ⬜ |
| 24 | gRPC Service | ⬜ |
| 25 | Mini Framework | ⬜ |

---

Good luck! 🚀 Each exercise builds on previous concepts. Take your time and understand each pattern deeply.
