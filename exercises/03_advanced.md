# Part 3: Advanced Patterns & Resilience (Exercises 21-30)

## Exercise 21: High-Performance LRU Cache
**Topics**: `container/list`, Generics, Sharding, Cache Eviction

### ❓ Why is this necessary?
Caches are essential for reducing latency and database load. An LRU (Least Recently Used) policy ensures the most relevant data stays in memory. In specific high-load scenarios (like a CDN or DB buffer pool), a single mutex becomes a bottleneck, requiring "sharding" to allow concurrent reads/writes.

### 🌍 Real-World Application
*   **Redis/Memcached**: In-memory key-value stores.
*   **Database Buffer Pools**: Keeping frequently accessed disk pages in RAM.

### 📦 Production Requirements
1.  Create `exercises/21_lru_cache/main.go`.
2.  **Generics**: Implement `LRUCache[K comparable, V any]`.
3.  **Sharding**: Instead of one big lock, create `type ShardedCache struct { shards []*LRUCache }`.
    - Hash the key to select a shard.
    - Each shard has its own `sync.RWMutex`.
4.  **TTL (Time To Live)**: Add an optional expiration time for items.
5.  **Eviction Callback**: Allow registering a function `OnEvicted(key K, value V)` (useful for metrics or cleanup).
6.  **Benchmark**: Compare `SingleLock` vs `ShardedLock` with heavy concurrent writes.

---

## Exercise 22: Asynchronous Event Bus
**Topics**: Observer Pattern, Dynamic Channels, Non-blocking Publish

### ❓ Why is this necessary?
Decoupling components allows them to evolve independently. An Event Bus allows a "User Service" to say "UserCreated" without knowing that the "Email Service", "Audit Service", and "Analytics Service" are all listening. It prevents spaghetti code dependencies.

### 🌍 Real-World Application
*   **Kubernetes Controllers**: Watching for resource changes.
*   **GUI Applications**: Button clicks processing.

### 📦 Production Requirements
1.  Create `exercises/22_event_bus/main.go`.
2.  **Dynamic Subscription**: `Subscribe(topic string, handler func(Event))` should return a subscription ID/Handle for easy unsubscription.
3.  **Non-Blocking Publish**: Use buffered channels or a dispatch goroutine so that a slow subscriber doesn't block the publisher.
4.  **Wildcards**: Support subscribing to `order.*` to receive `order.created` and `order.shipped`.
5.  **Graceful Shutdown**: `Close()` should wait for all pending events to be processed before returning.

---

## Exercise 23: Resilience Patterns (Retry, Backoff, Breaker)
**Topics**: `github.com/cenkalti/backoff`, Circuit Breaker, Jitter

### ❓ Why is this necessary?
Transient failures (network blips) are guaranteed in distributed systems. A simple retry loop can DDOS a struggling service. You need **Exponential Backoff** (wait 1s, 2s, 4s...) and **Jitter** (randomness) to spread load. **Circuit Breakers** stop requests entirely when a service is down to allow it to recover.

### 🌍 Real-World Application
*   **Microservices Communication**: Resilience against partial outages.
*   **AWS SDKs**: Built-in retries for throttling errors.

### 📦 Production Requirements
1.  Create `exercises/23_retry/main.go`.
2.  **Library Approach**: Implement a `ResilianceClient` that wraps an `http.Client`.
3.  **Retry Policy**: Use a customizable policy (MaxRetries, InitialInterval). Add Jitter.
4.  **Circuit Breaker State Machine**:
    - **Closed** (Normal): Requests go through.
    - **Open** (Broken): Fail immediately.
    - **Half-Open**: Allow 1 probe request. If success -> Closed. If fail -> Open.
5.  **Integration**: Simulate a flaky server (fails 80% of time) and prove that your client eventually succeeds without overwhelming it.

---

## Exercise 24: Production gRPC Service
**Topics**: Protocol Buffers, Interceptors, Status Codes, Metadata

### ❓ Why is this necessary?
gRPC is the standard for internal microservice communication due to its performance (Protobuf binary) and strict typing. Features like streaming and interceptors (middleware) make it powerful. You must know how to handle errors using standard gRPC codes, not HTTP codes.

### 🌍 Real-World Application
*   **Internal Microservices**: High-throughput service-to-service communication.
*   **Mobile Apps**: Efficient data transfer over shaky networks.

### 📦 Production Requirements
1.  Create `exercises/24_grpc/`.
2.  **Proto Definition**: Define `TodoService` with proper Protobuf style guide (PascalCase messages, snake_case fields).
3.  **Interceptors**: Implement:
    - **Logging**: Log RPC method, time, and error.
    - **Auth**: Extract JWT from `metadata` context.
4.  **Error Handling**: Return `codes.NotFound` for missing IDs, `codes.InvalidArgument` for bad input.
5.  **Bidirectional Streaming**: Implement a `Chat` RPC where client and server send messages freely to each other.

---

## Exercise 25: Modular Monolith Framework
**Topics**: Dependency Injection, Router Design, Middleware Stacks

### ❓ Why is this necessary?
While you shouldn't build a massive framework like Django, understanding how to structure a large application is vital. Use standard libraries where possible, but organize code into logical modules (Auth, Users, Billing) that share a common core (Router, DB, Logger).

### 🌍 Real-World Application
*   **Enterprise Go Apps**: Structuring large codebases (100k+ lines) to remain maintainable.
*   **Platform Engineering**: Building a "Service Template" for other teams.

### 📦 Production Requirements
1.  Create `exercises/25_mini_framework/`.
2.  **App Struct**: `type App struct { Router, DB, Config, Logger }`.
3.  **Module Interface**: `type Module interface { RegisterRoutes(r *Router) }`.
4.  **Trie Router**: Implement a prefix-tree (Trie) based router for efficient `GET /users/:id` matching (Don't use a regex loop).
5.  **Group Middleware**: `v1 := app.Group("/v1"); v1.Use(AuthMiddleware)`.
6.  **Context Helper**: Custom `Context` struct that handles JSON binding and validation automatically.

---

## Exercise 26: Functional Options Pattern
**Topics**: API Usability, Variadic Functions, Config Objects

### ❓ Why is this necessary?
Constructors with many parameters (`NewServer("localhost", 8080, 30, true, nil)`) are brittle and unreadable. The Functional Options pattern creates APIs that are easy to use (defaults work) and easy to extend (just add a new option) without breaking backward compatibility.

### 🌍 Real-World Application
*   **gRPC Client**: `grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())`.
*   **CLI Libraries**: Configuring command behavior.

### 📦 Production Requirements
1.  Create `exercises/26_functional_options/main.go`.
2.  **Scenario**: A `DatabaseConnection` pool.
3.  **Options**:
    - `WithMaxOpenConns(int)`
    - `WithMaxIdleConns(int)`
    - `WithConnMaxLifetime(duration)`
    - `WithLogger(Logger)`
4.  **Safety**: Ensure options can't corrupt the internal state (return error if invalid).
5.  **Immutability**: The config should be read-only after initialization.

---

## Exercise 27: Lock-Free Atomic Operations
**Topics**: `sync/atomic`, Compare-And-Swap (CAS), Memory Barriers

### ❓ Why is this necessary?
Mutexes involve OS context switches, which can be slow in ultra-high-performance hot paths. Atomics (CPU instructions) are orders of magnitude faster. They are the building blocks of synchronization primitives like Mutexes and Channels.

### 🌍 Real-World Application
*   **Metrics Counters**: Incrementing request counts in middleware (Mutex is overkill).
*   **Ring Buffers**: High-performance queues (LMAX Disruptor pattern).

### 📦 Production Requirements
1.  Create `exercises/27_atomic/main.go`.
2.  **Spinlock**: Implement a `SpinLock` using `atomic.CompareAndSwapInt32`.
    - Loop until you successfully swap 0 -> 1.
    - Unlock swaps 1 -> 0.
3.  **Atomic Config**: Use `atomic.Pointer[Config]` (Go 1.19+) to implementing "Hot Reloading" of global configuration without a Mutex lock on reads.
4.  **Benchmark**: spinlock vs `sync.Mutex` under high contention.

---

## Exercise 28: Memory Escape Analysis
**Topics**: Compiler Optimizations, Heap vs Stack, GC Pressure

### ❓ Why is this necessary?
In Go, variables on the **stack** are "free" (reclaimed when function returns). Variables on the **heap** must be tracked by the Garbage Collector (expensive). Understanding "escape analysis" helps you write zero-allocation code, critical for high-throughput systems.

### 🌍 Real-World Application
*   **High-Frequency Trading**: latency sensitivity requires zero GC pauses.
*   **Networking Libraries**: Reusing byte buffers to avoid allocations (`sync.Pool`).

### 📦 Production Requirements
1.  Create `exercises/28_escape_analysis/main.go`.
2.  **Analysis**: Write code that *accidentally* causes escapes:
    - Returning a pointer to a loop variable.
    - Passing a value to `fmt.Println` (interfaces escape).
    - Closing over a variable in a goroutine.
3.  **Fix**: Refactor the code to eliminate heap allocations where possible.
4.  **Proof**: Run `go build -gcflags="-m -m"` and interpret the output in comments. Benchmark the "before" and "after".

---

## Exercise 29: Semantic Difference & Deep Equality
**Topics**: `reflect`, Diff Algorithms, Test Helpers

### ❓ Why is this necessary?
Testing often requires comparing complex structs. `reflect.DeepEqual` is the standard but gives no info on *what* is different. Writing a customized "Diff" tool teaches you about the type system and recursion.

### 🌍 Real-World Application
*   **Terraform/Kubernetes**: Calculating the "diff" between desired state and actual state.
*   **Test Assertion Libraries**: `stretchr/testify` uses this to print pretty diffs.

### 📦 Production Requirements
1.  Create `exercises/29_deep_equal/main.go`.
2.  **Diff Function**: `Diff(a, b any) []string`. Return a list of differences (e.g., `".Author.Name: expected 'Bob', got 'Alice'"`).
3.  **Recursion**: Handle Slices, Maps, and Structs recursively.
4.  **Pointer Cycle Detection**: (Advanced) Handle cases where struct A points to B, and B points back to A (prevent infinite recursion).
5.  **Epsilon**: For floats, allow a small margin of error.

---

## Exercise 30: WebAssembly Integration
**Topics**: `syscall/js`, DOM Manipulation, Go-to-JS Bridge

### ❓ Why is this necessary?
WebAssembly (Wasm) allows you to run Go code in the browser. This enables sharing logic (validation, cryptography, business rules) between your Go backend and your Frontend, ensuring 100% consistency.

### 🌍 Real-World Application
*   **Video Editors**: Running heavy processing in browser.
*   **Cryptography**: End-to-end encryption libraries shared between server/client.

### 📦 Production Requirements
1.  Create `exercises/30_wasm/`.
2.  **DOM Interaction**: Create a button in HTML. Stick a Go Event Listener on it. When clicked, modify the DOM to show "Clicked!".
3.  **Async/Promises**: Expose a Go function that returns a Javascript `Promise` (using `js.FuncOf` and callbacks).
4.  **Build Optimization**: Use **TinyGo** to compile. Compare the binary size vs standard Go compiler (Production Wasm must be small).
