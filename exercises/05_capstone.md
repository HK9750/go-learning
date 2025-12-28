# Part 5: Distributed Systems Capstone (Exercises 41-50)

## Exercise 41: Event Sourcing Engine
**Topics**: Append-Only Logs, State Replay, Snapshots

### ❓ Why is this necessary?
In systems like Banking or Inventory, "Current Balance" is just a projection of history. Event Eourcing stores *what happened* (MoneyDeposited, MoneyWithdrawn), ensuring perfect auditability.

### 🌍 Real-World Application
*   **Bank Ledgers**: Creating a statement history.
*   **Git**: Reconstructing the file state from a series of commits.

### 📦 Production Requirements
1.  Create `exercises/41_event_sourcing/main.go`.
2.  **Store**: An append-only file (or DB table) `events.log`.
3.  **Aggregator**: `func Rehydrate(id string) (*Account, error)` that reads all events for an ID and applies them.
4.  **Snapshotting**: Every 100 events, save the current `Account` state to disk. Optimization: `Rehydrate` should load the snapshot and only play events since then.
5.  **Concurrency**: Ensure two users can't append conflicting events at the same time (Optimistic Concurrency Control with version numbers).

---

## Exercise 42: CQRS Implementation
**Topics**: Command Query Responsibility Segregation, Read Models

### ❓ Why is this necessary?
Optimizing for Writes (Normalization) and Reads (Denormalization/JOINs) often conflicts. CQRS splits them. You write to a "Command" model (Event Store) and it asynchronously updates a "Query" model (Elasticsearch/Redis).

### 🌍 Real-World Application
*   **High-Traffic Search**: Writes go to Postgres. A worker syncs changes to Elasticsearch for fast searching.

### 📦 Production Requirements
1.  Create `exercises/42_cqrs/`.
2.  **Command Side**: `POST /orders` -> Validates and writes `OrderCreated` event.
3.  **Query Side**: `GET /orders` -> Reads from a simple in-memory map.
4.  **Projector**: A background worker that listens for `OrderCreated` events and updates the Read Map.
5.  **Eventual Consistency**: Simulate a delay in the projector. Show that `POST` returns 200 OK immediately, but `GET` might be stale for 500ms.

---

## Exercise 43: Consistent Hashing Ring
**Topics**: Distributed Key-Value Store, Partitioning, Virtual Nodes

### ❓ Why is this necessary?
When you add/remove a cache server, you don't want to invalidate *all* keys (modulo hashing). Consistent hashing only moves `K/N` keys. This is fundamental to scalable databases (Cassandra, DynamoDB).

### 🌍 Real-World Application
*   **Load Balancers**: Routing sticky sessions to backend servers.
*   **Sharded Databases**: Determining which node holds User ID 12345.

### 📦 Production Requirements
1.  Create `exercises/43_consistent_hash/main.go`.
2.  **Ring**: Implement a standard Consistent Hash ring (Sorted Map of Hashes).
3.  **Virtual Nodes**: Each physical node (e.g., "Server A") should exist at 20 different points on the ring to prevent "hot spots".
4.  **Replication**: `GetNodes(key, n=3)` should return the primary owner and 2 replicas (the next 2 distinct nodes on the ring).
5.  **Visualization**: Print the distribution of 10,000 keys across 5 nodes. It should be roughly even (20% each).

---

## Exercise 44: Leader Election Simulation
**Topics**: Consensus, `sync/atomic`, Heartbeats, Split Brain

### ❓ Why is this necessary?
In a cluster, you usually want only *one* node doing the scheduling/writing to avoid conflicts. Choosing that node (Leader Election) requires robustness.

### 🌍 Real-World Application
*   **Kubernetes Scheduler**: Only one active scheduler.
*   **Etcd/Zookeeper**: Managing cluster state.

### 📦 Production Requirements
1.  Create `exercises/44_leader_election/main.go`.
2.  **Nodes**: Simulate 5 nodes (goroutines).
3.  **Bully Algorithm** (or simplfiied Raft):
    - Nodes have IDs. Highest ID wins.
    - If Leader doesn't send heartbeat for 500ms, start election.
4.  **Split Brain**: Simulate a network partition where Node 5 is cut off. It should declare itself leader (if naive), or the cluster should prefer the majority partition.
5.  **Visual Log**: "Node 1: I accept Node 5 as leader".

---

## Exercise 45: Blockchain Prototype
**Topics**: Cryptography (SHA256), Merkle Trees, Proof-of-Work

### ❓ Why is this necessary?
Understanding Blockchains teaches you about immutable ledgers, hashing, and distributed consensus (Nakamoto Consensus) without the hype.

### 🌍 Real-World Application
*   **Git**: Uses SHA1 hashes to chain commits.
*   **Audit Trails**: Verifiable logs that no admin can tamper with.

### 📦 Production Requirements
1.  Create `exercises/45_blockchain/main.go`.
2.  **Block Struct**: `Index`, `Timestamp`, `Data`, `PrevHash`, `Hash`, `Nonce`.
3.  **PoW**: Implement `Mine()` function that finds a hash starting with "0000" (Difficulty).
4.  **Verification**: Function to validate the entire chain (checking previous hashes match).
5.  **Tampering**: Manually modify data in Block 2 and prove that Block 3's hash is now invalid.

---

## Exercise 46: Interpreter - The Lexer
**Topics**: Compilers, State Machines, Tokenization

### ❓ Why is this necessary?
Parsing text (Config files, Domain Specific Languages) is a common advanced task. A Lexer turns "source code" into "tokens".

### 🌍 Real-World Application
*   **SQL Parsers**: Turning `SELECT * FROM` into tokens.
*   **JSON Parsers**: Reading `{ "key": "value" }`.

### 📦 Production Requirements
1.  Create `exercises/46_interpreter/lexer/`.
2.  **Language**: "Monkey" (from "Writing an Interpreter in Go") or simple Calculator expressions (`let x = 5 + 10;`).
3.  **Tokens**: Define `ILLEGAL`, `EOF`, `IDENT`, `INT`, `ASSIGN`, `PLUS`, etc.
4.  **State Machine**: `NextToken()` method that reads chars and advances position.
5.  **Test**: Comprehensive unit tests covering all edge cases.

---

## Exercise 47: Interpreter - The Parser
**Topics**: Recursive Descent, Pratt Parsing, AST

### ❓ Why is this necessary?
A Parser turns Tokens into a Tree Structure (AST) that the computer can understand. Pratt Parsers are elegant and handle operator precedence (`*` binds tighter than `+`) easily.

### 🌍 Real-World Application
*   **template/html**: Parsing HTML templates.
*   **Expression Evaluation**: Rules engines (e.g., "If Age > 18 AND Score > 50").

### 📦 Production Requirements
1.  Create `exercises/46_interpreter/parser/`.
2.  **ASTNodes**: `LetStatement`, `ReturnStatement`, `ExpressionStatement`, `IntegerLiteral`.
3.  **Pratt Parsing**: Register prefix and infix parse functions.
4.  **Precedence**: Ensure `5 * 5 + 10` is parsed as `(5 * 5) + 10`.

---

## Exercise 48: Interpreter - The Evaluator
**Topics**: Tree Walking, Environments, Recursion

### ❓ Why is this necessary?
The Evaluator actually runs the code. You'll learn about "Environments" (scopes) - how `x` in a function is different from global `x`.

### 🌍 Real-World Application
*   **Scripting Languages**: Lua embedded in games.
*   **Policy Engines**: Open Policy Agent (Rego).

### 📦 Production Requirements
1.  Create `exercises/46_interpreter/evaluator/`.
2.  **Object System**: `Integer`, `Boolean`, `Null`, `ReturnValue`, `Error`.
3.  **Eval**: Recursive function `Eval(node ASTNode, env *Environment) Object`.
4.  **Environment**: Enable nested scopes (function calls extending the outer scope).
5.  **REPL**: A `main.go` loop that accepts input, lexes, parses, evaluates, and prints result.

---

## Exercise 49: P2P File Sharing System
**Topics**: BitTorrent Protocol, Choking, Piece Selection

### ❓ Why is this necessary?
P2P systems are the ultimate decentralization. You deal with peers that lie, disconnect, or are slow. You need robust error handling.

### 🌍 Real-World Application
*   **Container Distribution**: Dragonfly (P2P Docker registry).
*   **Software Updates**: Windows Update Delivery Optimization.

### 📦 Production Requirements
1.  Create `exercises/49_p2p/`.
2.  **Handshake**: Implement a TCP handshake exchanging "Peer ID".
3.  **Bitfield**: Map of which "pieces" of the file a peer has.
4.  **Downloading**: Request pieces from peers that have them.
5.  **Integrity**: Verify pieces against a SHA1 hash.
6.  **Concurrency**: Download from 4 peers simultaneously.

---

## Exercise 50: Microservices Capstone
**Topics**: Distributed Tracing, API Gateway, Service Discovery, Logging

### ❓ Why is this necessary?
This puts everything together. You are building a "Platform".

### 🌍 Real-World Application
*   **Uber/Netflix backend**: Thousands of services talking to each other.

### 📦 Production Requirements
1.  Create `exercises/50_microservices/`.
2.  **Architecture**:
    - **Gateway**: HTTP entry point, forwards requests.
    - **Auth Service**: Issues JWTs.
    - **Order Service**: Creates orders (SQL).
    - **Payment Service**: Processes payments (Mock).
3.  **Infrastructure**:
    - **Tracing**: Pass a `X-Trace-ID` header through **all** services. Log it in every message.
    - **Service Discovery**: Services shouldn't know localhost URLs. They should ask a "Registry" (simulate with a map).
4.  **Docker Compose**: Write a `docker-compose.yml` to spin up all 4 services + Postgres + Redis.

---

## Bonus Challenges
*   **Kubernetes Operator**: Write a custom controller.
*   **eBPF**: Write a verified BPF program in Go.
*   **Compiler**: Compile your Interpreter language to x86 assembly.
