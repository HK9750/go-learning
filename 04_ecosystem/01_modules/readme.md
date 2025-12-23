# Go Modules & Packages Deep Dive

## 1. The `go.mod` file
This file defines:
- The module path (unique namespace).
- The Go version.
- Dependencies (direct and indirect).

Example:
```go
module github.com/username/project

go 1.21

require (
    github.com/google/uuid v1.3.0
    github.com/gin-gonic/gin v1.9.0 // indirect
)
```

## 2. Package Visibility (Capitalization)
- **Capitalized** identifiers (Functions, Variables, Structs) are **EXPORTED** (Public).
- **Lowercased** identifiers are **UNEXPORTED** (Private, package-local).

## 3. The `internal` directory
- Any package named `internal/` or inside an `internal/` directory cannot be imported by code outside the *module base*.
- Enforces strict boundaries.
- `project/internal/foo` can ONLY be imported by `project/...`.
- `other-project` CANNOT import `project/internal/foo`.

## 4. Initialization Order
When your program starts:
1. All `import`ed packages are initialized.
2. Package-level `var` declarations are evaluated.
3. `init()` functions run.
   - You can have multiple `init()` functions per file.
   - They run in order of appearance.
4. Finally, `main()` runs.

## 5. Workspaces (`go.work`) (Go 1.18+)
Allows working on multiple modules simultaneously without pushing to remote.
File: `go.work`
```go
go 1.21

use (
    ./my-app
    ./my-lib
)
```
This overrides local `go.mod` requirements to point to local folders.

## 6. TROUBLESHOOTING & BEST PRACTICES

1. **"ambiguous import"**:
   - Happens when multiple modules provide the same package.
   - Fix: Use `go.mod` 'replace' directive to pin one.

2. **Indirect Dependencies**:
   - You'll see `// indirect` in `go.mod`.
   - These are deps of your deps. Do not edit manually.
   - Run `go mod tidy` to clean up unused ones.

3. **Caching issues**:
   - `go clean -modcache` can fix weird dependency states.

