# Part 2: Intermediate Data Structures & Algorithms (Exercises 11-20)

## Exercise 11: Production Database Patterns (SQLite)
**Topics**: `database/sql`, Repository Pattern, Transactions, `context`

### ❓ Why is this necessary?
Directly writing SQL in handlers leads to unmaintainable code. The **Repository Pattern** abstracts data access, allowing you to swap databases or mock them for testing. You must also handle **SQL Injection** (via prepared statements) and **timeouts** (via Context).

### 🌍 Real-World Application
*   **User Microservice**: Separating SQL logic from HTTP logic so tests can run without a real DB.
*   **Financial Ledger**: Using atomic transactions to ensure money isn't lost if a server crashes mid-transfer.

### 📦 Production Requirements
1.  Create `exercises/11_sqlite_crud/main.go`.
2.  **Schema Migration**: Create a function `InitDB(db *sql.DB)` that creates the `users` table if it doesn't exist.
3.  **Repository Interface**:
    ```go
    type UserRepository interface {
        Create(ctx context.Context, u User) error
        GetByID(ctx context.Context, id int64) (*User, error)
        Update(ctx context.Context, u User) error
        Delete(ctx context.Context, id int64) error
    }
    ```
4.  **Context Usage**: **ALL** SQL methods must take `context.Context` and use `db.ExecContext`, `db.QueryRowContext`, etc.
5.  **Transactions**: Implement a method `helper.RunInTx(db, func(tx) error)` to safely handle commit/rollback.
6.  **Validation**: Ensure email is unique (handle `sqlite3` constraint violation error explicitly).

---

## Exercise 12: Professional CLI Tooling
**Topics**: `github.com/spf13/cobra`, `spf13/viper`, Configuration, graceful exit

### ❓ Why is this necessary?
Ops tools, deploy scripts, and dev utilities are all CLIs. A production CLI needs structure (subcommands), configuration management (env vars/files), and standardized help output. Cobra is the industry standard (used by Kubernetes `kubectl`, Docker, Hugo).

### 🌍 Real-World Application
*   **`kubectl`**: Managing K8s clusters.
*   **`docker`**: Managing containers.
*   **Database Migrators**: CLI tools to apply schema changes.

### 📦 Production Requirements
1.  Create `exercises/12_cobra_cli/`.
2.  **Structure**: Use `cmd/root.go`, `cmd/add.go`, `cmd/list.go` pattern.
3.  **Viper Integration**: Allow configuration of the data file path via:
    - Flag: `--config ~/.my-tasks.json`
    - Environment Variable: `TASK_CONFIG=...`
    - Default: `$HOME/.tasks.json`
4.  **Output Formatting**: Use `text/tabwriter` to print aligned tables for the `list` command.
5.  **Signal Handling**: If the user presses Ctrl+C during a long operation, print "Cancelling..." and exit cleanly.

---

## Exercise 13: Robust WebSocket Chat
**Topics**: `github.com/gorilla/websocket`, Ping/Pong, Concurrency Safety

### ❓ Why is this necessary?
WebSockets are persistent connections. Unlike HTTP, you must actively manage the connection state: open, close, broken pipe, and "zombie" connections. You also need to safeguard shared memory (the list of connected clients) using mutexes or channels.

### 🌍 Real-World Application
*   **Slack/Discord**: Real-time messaging implementation.
*   **Live Sports Scores**: Pushing updates to thousands of connected clients.

### 📦 Production Requirements
1.  Create `exercises/13_websocket_chat/main.go`.
2.  **Hub Pattern**: Create a `Hub` struct that manages `register`, `unregister`, and `broadcast` channels to avoid mutex locking in the main loop.
3.  **Ping/Pong**: Implement a heartbeat.
    - Server sends "Ping" every 30s.
    - If Client doesn't reply "Pong", disconnect them.
4.  **Clean Exit**: On server shutdown, send a "Going Away" close frame to all clients before stopping.
5.  **Message Protocol**: JSON payloads: `{"type": "msg", "data": "..."}` or `{"type": "typing"}`.

---

## Exercise 14: Fan-Out/Fan-In Pipeline
**Topics**: Channels, Buffering, Error Propagation, `sync.ErrGroup`

### ❓ Why is this necessary?
Data processing pipelines (ETL) often need to run in parallel stages. Go's channels make this natural. Production pipelines must handle **backpressure** (buffered channels), **error propagation** (if one stage fails, stop everything), and **cancellation**.

### 🌍 Real-World Application
*   **Image Processing**: Resize -> Compress -> Upload to S3 (in parallel).
*   **Log Ingestion**: Read file -> Parse JSON -> Batch insert to Elasticsearch.

### 📦 Production Requirements
1.  Create `exercises/14_pipeline/main.go`.
2.  **Stages**:
    - `Generator(ctx)`: Stream path names to a channel.
    - `Processor(ctx, paths)`: Read file, calculate SHA256 hash (slow task). Start **5 concurrent workers**.
    - `Collector(ctx, results)`: Aggregate hashes into a map.
3.  **Cancellation**: Use `context.WithCancel`. If any worker encounters a "bad file" error, cancel the context and stop all other workers immediately.
4.  **Graceful Error Handling**: Return the *first* error encountered to `main`.

---

## Exercise 15: Context Timeouts & Deadlines
**Topics**: `context`, Network Resilience, Cascading Failure Prevention

### ❓ Why is this necessary?
A slow dependency (database, external API) should not hang your entire system. Contexts allow you to set a hard deadline. If an operation takes too long, you cut it off to free up resources (`408 Request Timeout`). This is critical for preventing "pile-ups".

### 🌍 Real-World Application
*   **Microservices**: Service A calls Service B. Service A sets a 100ms timeout so the user isn't left waiting forever.
*   **Database Queries**: Killing long-running "runaway" queries that are slowing down the DB.

### 📦 Production Requirements
1.  Create `exercises/15_context_timeout/main.go`.
2.  **Mock Server**: A handler that sleeps for `X` milliseconds based on a query param.
3.  **Client with Timeout**:
    - Make a request with `context.WithTimeout(ctx, 100*time.Millisecond)`.
    - Handle `url.Error` properly. Check `ctx.Err() == context.DeadlineExceeded`.
4.  **Cascading Context**: creating a child context from a parent. Cancel the parent and verify the child also cancels immediately.
5.  **Resource Cleanup**: Ensure `defer cancel()` is called to avoid context leaks.

---

## Exercise 16: Type-Safe Generics
**Topics**: Type Parameters (`[T any]`), Constraints, Slice Utilities

### ❓ Why is this necessary?
Before Go 1.18, we used `interface{}`/reflection for generic algorithms, losing type safety. Generics allow us to write reusable logic (Map, Filter, Sets) that is checked at compile time. This reduces runtime panics and code duplication.

### 🌍 Real-World Application
*   **ORM Libraries**: Returning `[]User` instead of `[]interface{}`.
*   **Data Structure Pools**: A `Pool[T]` that manages specific types.

### 📦 Production Requirements
1.  Create `exercises/16_generics/main.go`.
2.  **Result Type**: Implement a generic `Result[T any]` struct that contains `Value T` and `Error error`.
3.  **Concurrent Map**: Implement `func MapAsync[T, U any](ctx context.Context, data []T, fn func(T) U) ([]U, error)`.
    - It should verify safe concurrency (mutex or channels).
    - It must respect the context.
4.  **Constraints**: Write a function `Max[T constraints.Ordered](a, b T) T`.
5.  **Use Case**: Fetch 5 URLs concurrently using `MapAsync` and return the body lengths.

---

## Exercise 17: Reflection & Struct Validation
**Topics**: `reflect`, Tag Parsing, Recursive Validation

### ❓ Why is this necessary?
While reflection is slow, it is indispensable for libraries like serialization (JSON), configuration loading, and validation. Understanding it demystifies how `json.Unmarshal` works. In production, you might build a custom validator for config files.

### 🌍 Real-World Application
*   **`go-playground/validator`**: The standard validation library for structs.
*   **Config Loaders**: Mapping environment variables `MY_APP_PORT` to struct fields `Port`.

### 📦 Production Requirements
1.  Create `exercises/17_reflection/main.go`.
2.  **Library**: Write a function `Validate(v any) error` that supports:
    - `val:"required"` (strings not empty, numbers not zero).
    - `val:"min=10"` (numbers >= 10, string length >= 10).
3.  **Recursion**: If a struct contains a nested struct, validate that too.
4.  **Performance**: Add a benchmark. Compare your reflection validator vs a hardcoded `if` check validator (Observe the performance penalty!).

---

## Exercise 18: Rigorous Testing
**Topics**: `testing`, Table-Driven, Subtests, Mocks, Race Detector

### ❓ Why is this necessary?
Code without tests is technical debt. In Go, table-driven tests are the idiom. You must also know how to mock dependencies (like an external API) so your unit tests are fast and deterministic. The race detector is mandatory for concurrent code.

### 🌍 Real-World Application
*   **CI/CD**: Tests run on every Pull Request.
*   **Refactoring Safety**: Confidence to change code without breaking features.

### 📦 Production Requirements
1.  Create `exercises/18_testing/`
2.  **Interface Mocking**: Define `WeatherService` interface. Generate a mock manually or use `mockery`/`gomock` (bonus).
3.  **Table-Driven**: Test a `GetTemperature(city string)` function that uses the service.
    - Test cases: "Success", "API Error", "Invalid City".
4.  **Subtests**: Use `t.Run("case name", func(t *testing.T) {...})` for nice output.
5.  **Race Validations**: Write a test that sparks race conditions (shared map write), run with `go test -race`, fail, then fix it with a Mutex.

---

## Exercise 19: Composition over Inheritance
**Topics**: Embedding, Interface Segregation, Flexible Design

### ❓ Why is this necessary?
Go has no inheritance. You build complex objects by composing smaller ones. This promotes "Has-A" relationships over "Is-A". Understanding interface satisfaction via embedding is key to writing modular, testable systems (e.g. standard library `io.ReadWriter`).

### 🌍 Real-World Application
*   **Middleware**: Wrapping an `http.ResponseWriter` to capture status codes.
*   **Cross-Platform Systems**: A `File` struct that embeds platform-specific logic but exposes a common interface.

### 📦 Production Requirements
1.  Create `exercises/19_composition/main.go`.
2.  **Core Domain**: A `BaseEntity` struct with `ID`, `CreatedAt`, `UpdatedAt`.
3.  **Embedding**: `User` embeds `BaseEntity`. `Product` embeds `BaseEntity`.
4.  **Interface Segregation**: Define `Auditable` interface with `GetCreatedAt()`. Write a generic function that accepts any `Auditable`.
5.  **Decorator**: Create a `LoggingWriter` struct that embeds `io.Writer`. Override `Write()` to count bytes written, then call the inner writer. Test this with `os.Stdout`.

---

## Exercise 20: Robust Worker Pool
**Topics**: Worker Pattern, Job Queues, Panic Recovery, Monitoring

### ❓ Why is this necessary?
Launching a goroutine for every HTTP request can exhaust memory (OOM). A Worker Pool limits concurrency to a fixed number (e.g., 50 workers). This safeguards your system against spikes in load. You must also handle panics inside workers so one crash doesn't kill the process.

### 🌍 Real-World Application
*   **Video Transcoding**: Limiting expensive CPU tasks.
*   **Email Sending**: Throttling requests to an SMTP server.

### 📦 Production Requirements
1.  Create `exercises/20_worker_pool/main.go`.
2.  **Dispatcher**: Start `N` workers.
3.  **Panic Recovery**: Wrap the worker logic in a `defer func() { if r := recover(); r != nil { ... } }`. Log the stack trace but keep the worker (or restart it) alive.
4.  **Metrics**: Track `JobsProcessed`, `JobsFailed`, `ActiveWorkers` using `sync/atomic`.
5.  **Drain**: Implement a `Stop()` method that closes the job channel and waits for all workers to finish (`sync.WaitGroup`) before returning.
