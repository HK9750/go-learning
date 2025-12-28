# Part 4: Systems Engineering & Tooling (Exercises 31-40)

## Exercise 31: Fuzz Testing (Security)
**Topics**: `testing` (Go 1.18+), Random Input Generation, Crash Detection

### ❓ Why is this necessary?
Unit tests only check what you predict. Fuzzing throws random garbage at your code to find edge cases (off-by-one, UTF-8 decoding panics, memory exhaustion). It is mandatory for any code defining a protocol or parser in production.

### 🌍 Real-World Application
*   **Standard Library**: Go team fuzzes `encoding/json` and `net/http` to find security vulnerabilities.
*   **Parsers**: Finding exploits in image decoders or data deserializers.

### 📦 Production Requirements
1.  Create `exercises/31_fuzzing/`.
2.  **Target**: Write a complex function `ParseConfig(data []byte) (Config, error)` that expects a custom binary format (e.g., "[4 bytes length][content]").
3.  **Fuzz Target**: `FuzzParseConfig(f *testing.F)`. Seed it with valid samples.
4.  **Invariants**: Ensure `ParseConfig` *never* panics, even with malformed input.
5.  **Round-Trip**: Property test: `Encode(Decode(x)) == x`.

---

## Exercise 32: Profiling & Performance Engineering
**Topics**: `net/http/pprof`, Flame Graphs, Allocations per Op

### ❓ Why is this necessary?
"It feels slow" is not an engineering metric. You need to know exactly *where* CPU cycles are spent and *where* memory is allocated. Production services almost always expose pprof endpoints (protected by auth) to debug live latency spikes.

### 🌍 Real-World Application
*   **Latency Debugging**: Finding why the 99th percentile response time jumped from 20ms to 500ms.
*   **Memory Leaks**: Identifying a goroutine leak or a map that grows forever.

### 📦 Production Requirements
1.  Create `exercises/32_optimization/main.go`.
2.  **Bad Function**: Creating a "Log Processor" that concatenates strings using `+` in a tight loop and allocates massive buffers.
3.  **Pprof Integration**: Start the standard pprof server on `localhost:6060`.
4.  **Analysis**:
    - Run `go tool pprof -http=:8080 ...`
    - Capture a **Heap Profile** (allocations) and **CPU Profile**.
5.  **Optimization**: Rewrite using `strings.Builder` and `sync.Pool` to reuse buffers.
6.  **Benchmark**: Demonstrate a 10x improvement in `ns/op` and `B/op`.

---

## Exercise 33: Plugin Architecture (RPC-based)
**Topics**: `net/rpc`, HashiCorp `go-plugin`, Subprocesses

### ❓ Why is this necessary?
Go's native `plugin` package is limited (Linux only, same compiler version). The industry standard for extensible apps (like Terraform) is RPC-based plugins. The plugin runs as a separate process and communicates over stdout/TCP. This isolates crashes.

### 🌍 Real-World Application
*   **Terraform Providers**: Each provider (AWS, Google) is a separate binary.
*   **Linter Frameworks**: Static analysis tools running custom rule sets.

### 📦 Production Requirements
1.  Create `exercises/33_plugin_system/`.
2.  **Contract**: Define `type Greeter interface { Greet() string }` in a shared module.
3.  **Plugin Implementation**: A `main` package that listens on TCP (or uses `hashicorp/go-plugin` generic scaffolding).
4.  **Host Application**: Launches the plugin binary (`exec.Command`), parses the handshake, connects, and calls the RPC.
5.  **Resilience**: If the plugin crashes, the host should log an error (not panic) and potentially restart it.

---

## Exercise 34: Custom Static Analysis (Linter)
**Topics**: `golang.org/x/tools/go/analysis`, AST, Engineering Standards

### ❓ Why is this necessary?
Code reviews are for logic, not style. Linters enforce architectural rules automatically. You might want to ban `fmt.Println` in production code (enforce structured logging) or ensure every `struct` has a specific tag.

### 🌍 Real-World Application
*   **`golangci-lint`**: The de-facto linter runner.
*   **Uber/Google**: Internal linters to enforce "No external calls in transactions".

### 📦 Production Requirements
1.  Create `exercises/34_static_analysis/`.
2.  **Rules**: Write an `Analyzer` that reports an error if:
    - `context.Background()` is used inside a function (it should be passed in).
    - A JSON tag is missing on an API response struct.
3.  **Test**: Run it against a `testdata/` directory with "golden" files to verify it catches the bugs.

---

## Exercise 35: Code Instrumentation (AST Rewriting)
**Topics**: `go/ast`, `go/printer`, Automatic Tracing

### ❓ Why is this necessary?
Manually adding `TraceStart()` and `TraceEnd()` to 1000 functions is toil. AST manipulation allows you to refactor or instrument massive codebases instantly. This is how tools like OpenTelemetry auto-instrumentation work.

### 🌍 Real-World Application
*   **Observability**: Injecting tracing spans into every HTTP handler automatically.
*   **Refactoring**: mass-renaming a method across 500 files safely.

### 📦 Production Requirements
1.  Create `exercises/35_ast_rewrite/main.go`.
2.  **Goal**: Inject a `defer timer.Track(time.Now(), "funcName")` at the start of every function in a package.
3.  **Robustness**: Handle functions that already have comments or different formatting.
4.  **Imports**: Automatically add the `timer` package to `imports` if it’s missing.
5.  **Output**: Write the modified formatted code to `_instrumented.go`.

---

## Exercise 36: Code Generation (`go generate`)
**Topics**: `text/template`, Stringer, Embedding

### ❓ Why is this necessary?
Go prefers code generation over runtime magic (reflection). It is faster and type-safe. You use this to generate Enums, Mocks, or SQL boilerplate.

### 🌍 Real-World Application
*   **`mockgen`**: Generating test mocks.
*   **`sqlc`**: Generating type-safe Go SQL code from raw SQL queries.

### 📦 Production Requirements
1.  Create `exercises/36_code_gen/`.
2.  **Input**: A definition file (YAML or specifically tagged Go struct) defining "Features" (Feature Flags).
3.  **Generator**: Write a tool that reads the defs and outputs a `flags.go` file with methods:
    - `func (f Feature) IsEnabled() bool`
4.  **Integration**: Add `//go:generate go run tools/gen.go` at the top of main.
5.  **CI Check**: Ensure strictly that generated code is committed and up-to-date.

---

## Exercise 37: Low-Level TCP Protocol Design
**Topics**: `net`, Framing, Endianness, IO Buffers

### ❓ Why is this necessary?
Sometimes HTTP is too heavy. Understanding TCP (streams, not packets) is crucial for keeping persistent connections (gaming, trading). You must handle **Framing** (knowing where one message ends and the next begins).

### 🌍 Real-World Application
*   **Database Protocols**: MySQL/Postgres wire protocols.
*   **IoT Devices**: Sending sensor data over raw TCP to save bandwidth.

### 📦 Production Requirements
1.  Create `exercises/37_tcp_chat/main.go`.
2.  **Protocol**: Header `[4 bytes length]`, Body `[JSONPayload]`.
3.  **Splitter**: Implement a `bufio.SplitFunc` to handle the framing correctly.
4.  **Connection Manager**: Handle 10k concurrent connections (simulated).
5.  **Timeouts**: Set `SetDeadline` on connections to kill idle clients.

---

## Exercise 38: High-Performance UDP Ingestion
**Topics**: `net.ListenUDP`, Lossy Protocols, Ring Buffers

### ❓ Why is this necessary?
UDP is "fire and forget". It's used where speed > correctness (video, gaming, metrics). The OS buffer can overflow if you don't read fast enough.

### 🌍 Real-World Application
*   **StatsD/DogStatsD**: Metrics agents accept UDP packets to avoid slowing down the app.
*   **Multiplayer Games**: Player position updates (if you miss one, just use the next one).

### 📦 Production Requirements
1.  Create `exercises/38_udp_tracker/main.go`.
2.  **Worker Pool**: One reader goroutine simply copies packets from the socket to a buffered channel (Ring Buffer).
3.  **Processing**: Multiple workers parse the packet and update a map.
4.  **Metrics**: Count `PacketsDropped` (if channel is full).
5.  **Shedding**: Intentionally drop packets if the system is overloaded.

---

## Exercise 39: GraphQL Gateway
**Topics**: `gqlgen`, Schema Design, Dataloaders

### ❓ Why is this necessary?
REST over-fetches (gets too much data) or under-fetches (needs N+1 requests). GraphQL solves this but introduces the **N+1 DB Query problem**. Dataloaders are the mandatory fix in production.

### 🌍 Real-World Application
*   **GitHub API v4**: Exposes the entire graph of repos/issues via GraphQL.
*   **Mobile Apps**: Aggregating data from 5 microservices into 1 request.

### 📦 Production Requirements
1.  Create `exercises/39_graphql/`.
2.  **Schema**: `User -> Posts -> Comments`.
3.  **N+1 Problem**: Intentionally implement a naive resolver that makes a DB call per post. Observe the log flood.
4.  **Dataloader**: Implement a `UserLoader` that batches IDs and fetches 100 users in 1 query.
5.  **Complexity Limit**: Prevent malicious queries (`User { Friends { Friends ... } }`) by limiting depth.

---

## Exercise 40: Clean Architecture (Domain-Driven Design)
**Topics**: Dependency Injection, Layers, Hexagonal Arch

### ❓ Why is this necessary?
"Spaghetti code" happens when HTTP handlers talk directly to SQL. Clean Architecture separates **Domain** (Business Rules), **Application** (Use Cases), and **Infrastructure** (Db/Web). This makes the core logic independent of frameworks.

### 🌍 Real-World Application
*   **Banking Core**: The logic "Transfer Money" must be testable without starting a web server or database.

### 📦 Production Requirements
1.  Create `exercises/40_clean_arch/`.
2.  **Structure**:
    - `internal/domain`: Pure Go structs (`Order`), Errors. No external deps.
    - `internal/service`: Business logic (`PayOrder`). Accepts interfaces.
    - `internal/repository/postgres`: Implementation of `domain.Repository`.
    - `cmd/api`: Wiring everything together (DI).
3.  **Rule**: `domain` must NOT import `repository`.
4.  **Test**: unit test `service` using a Mock Repository.
