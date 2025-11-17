# System Patterns

> **Derived from**: projectbrief.md
> **Purpose**: Documents how the system is architected and key technical patterns

## Architecture Overview
### High-Level Structure
```
┌─────────────┐
│   Browser   │ (or any HTTP client)
│  Terminal   │
│    Slack    │
└──────┬──────┘
       │ HTTP Request: g/zap/issues/42
       ↓
┌──────────────────────────────────────┐
│    Zap HTTP Service (localhost:80)   │
│  ┌────────────────────────────────┐  │
│  │  HTTP Router (httprouter)      │  │
│  │  - / (redirect handler)        │  │
│  │  - /healthz (health check)     │  │
│  │  - /varz (config dump)         │  │
│  └──────────┬─────────────────────┘  │
│             ↓                         │
│  ┌────────────────────────────────┐  │
│  │  URL Expansion Engine          │  │
│  │  - Parse path tokens           │  │
│  │  - Recursive tree traversal    │  │
│  │  - Apply expansion rules       │  │
│  └──────────┬─────────────────────┘  │
│             ↓                         │
│  ┌────────────────────────────────┐  │
│  │  Configuration Layer           │  │
│  │  - YAML → JSON parser          │  │
│  │  - File watcher (fsnotify)     │  │
│  │  - Thread-safe updates         │  │
│  └──────────┬─────────────────────┘  │
└─────────────┼────────────────────────┘
              ↓
       ┌──────────────┐
       │   c.yml      │ (config file)
       └──────────────┘

Response: HTTP 302 → https://github.com/issmirnov/zap/issues/42
```

### Component Breakdown
- **HTTP Router** (cmd/zap/web.go): Lightweight routing using julienschmidt/httprouter, handles incoming requests
- **URL Expansion Engine** (cmd/zap/text.go): Core business logic - tokenizes paths and recursively expands shortcuts
- **Configuration Layer** (cmd/zap/config.go): Parses YAML, validates, watches for changes, hot reloads
- **Context Manager** (cmd/zap/structs.go): Thread-safe config access with sync.RWMutex
- **Main** (cmd/main.go): Entry point, initializes service, starts HTTP server

### Data Flow
```
HTTP Request → Router → Path Tokenization → Tree Traversal → URL Construction → HTTP 302 Redirect

Example: GET /g/zap/issues/42
  ↓ tokenize: ["g", "zap", "issues", "42"]
  ↓ lookup "g" → "github.com"
  ↓ lookup "zap" → "issmirnov/zap"
  ↓ lookup "issues" → "issues"
  ↓ append remaining: "42"
  ↓ construct: "https://github.com/issmirnov/zap/issues/42"
  ↓ HTTP 302 redirect
```

## Design Patterns

### Pattern: Configuration as Single Source of Truth
**Where Used**: cmd/zap/config.go, entire system
**Why**: Eliminates need for database, ensures config file is authoritative
**Example**:
```go
// Config is loaded once, then hot-reloaded on file changes
func loadConfig(fs afero.Fs, configPath string) (*gabs.Container, error) {
    yamlData, err := afero.ReadFile(fs, configPath)
    jsonData, _ := yaml.YAMLToJSON(yamlData)
    return gabs.ParseJSON(jsonData)
}
```

### Pattern: Hot Reload with File Watching
**Where Used**: cmd/zap/config.go:watchConfig()
**Why**: Zero-downtime config updates for rapid iteration
**Example**:
```go
watcher, _ := fsnotify.NewWatcher()
go func() {
    for event := range watcher.Events {
        if event.Op&fsnotify.Write == fsnotify.Write {
            loadConfig() // Reload on file write
        }
    }
}()
```

### Pattern: Recursive Tree Traversal
**Where Used**: cmd/zap/text.go:getURL()
**Why**: Supports unlimited nesting depth elegantly
**Example**:
```go
// Recursively traverse config tree by splitting path on "/"
func getURL(c *gabs.Container, path string) string {
    tokens := strings.Split(path, "/")
    return recurse(c, tokens, 0, "")
}
```

### Pattern: Thread-Safe Config Access
**Where Used**: cmd/zap/structs.go:Context, all handlers
**Why**: Hot reload requires mutex protection for concurrent reads/writes
**Example**:
```go
type Context struct {
    Mux    *sync.RWMutex
    Config *gabs.Container
}
// Handlers acquire read lock before accessing config
ctx.Mux.RLock()
defer ctx.Mux.RUnlock()
```

### Pattern: Filesystem Abstraction
**Where Used**: afero.Fs throughout codebase
**Why**: Enables testing without real filesystem I/O
**Example**:
```go
// Production uses real filesystem, tests use in-memory
fs := afero.NewOsFs()           // Production
fs := afero.NewMemMapFs()       // Tests
```

## Key Technical Decisions

### Decision: YAML → JSON → gabs.Container
- **Context**: Need to parse hierarchical config format
- **Options Considered**: Native YAML parsing, custom DSL, JSON only
- **Decision**: Convert YAML to JSON, parse with gabs
- **Rationale**: YAML is user-friendly, JSON has better Go tooling, gabs provides path traversal
- **Consequences**: Two-step parsing adds complexity but improves UX and developer experience

### Decision: julienschmidt/httprouter over net/http
- **Context**: Need high-performance routing for 150k+ QPS goal
- **Options Considered**: stdlib mux, gorilla/mux, httprouter
- **Decision**: httprouter for performance and simplicity
- **Rationale**: Zero allocations, fastest router available, minimal API
- **Consequences**: Slight learning curve but massive performance gain

### Decision: HTTP 302 Redirects (not 301)
- **Context**: How to redirect users to expanded URLs
- **Options Considered**: 301 permanent, 302 temporary, 307/308
- **Decision**: 302 temporary redirects
- **Rationale**: Shortcuts may change over time, don't want browser caching
- **Consequences**: Slightly more requests but ensures latest config is used

### Decision: Stateless Service (No Database)
- **Context**: Where to store shortcut definitions
- **Options Considered**: SQLite, Redis, PostgreSQL, flat file
- **Decision**: YAML config file only, no database
- **Rationale**: Simplicity, version control, no deployment complexity
- **Consequences**: Limited to single-machine state but aligns with local-first design

### Decision: localhost Binding by Default
- **Context**: Security concerns for redirect service
- **Options Considered**: Bind to 0.0.0.0, require auth, localhost only
- **Decision**: Default to 127.0.0.1, allow override
- **Rationale**: Redirect service on public internet is security risk
- **Consequences**: Network-wide access requires explicit configuration

## Component Relationships

### HTTP Router ↔ URL Expansion Engine
- **Interaction**: Router extracts path from request, passes to expansion engine
- **Dependencies**: Router needs expansion result to construct redirect
- **Interface**: `getURL(config, path) → (url string, error)`

### Configuration Layer ↔ URL Expansion Engine
- **Interaction**: Expansion engine reads from config container
- **Dependencies**: Engine needs current config, config layer provides thread-safe access
- **Interface**: Shared `*gabs.Container` protected by RWMutex

### File Watcher ↔ Configuration Layer
- **Interaction**: Watcher triggers config reload on file changes
- **Dependencies**: Watcher needs file path, config layer handles reload
- **Interface**: fsnotify events → `loadConfig()` call

## Code Organization

### Directory Structure
```
zap2/
├── cmd/
│   ├── main.go              # Entry point, HTTP server setup
│   └── zap/
│       ├── config.go        # YAML parsing, validation, file watching
│       ├── config_test.go   # Config tests with in-memory filesystem
│       ├── structs.go       # Core data structures (Context, etc.)
│       ├── text.go          # URL expansion logic
│       ├── text_test.go     # Expansion algorithm tests
│       ├── web.go           # HTTP handlers
│       └── web_test.go      # HTTP integration tests
├── c.yml                    # Default config with examples
├── e2e.sh                   # End-to-end bash tests
├── go.mod / go.sum          # Go module definition
├── goreleaser.yml           # Multi-platform release config
├── .github/workflows/       # CI/CD automation
├── README.md                # User documentation
└── CONTRIBUTING.md          # Developer documentation
```

### Module Responsibilities
- **cmd/main.go**: Bootstrap only - parse flags, load config, start HTTP server
- **cmd/zap/config.go**: All configuration concerns - parsing, validation, file watching, /etc/hosts updates
- **cmd/zap/text.go**: Pure URL expansion logic - no I/O, no HTTP, just transformation
- **cmd/zap/web.go**: HTTP layer only - routing, request handling, response construction
- **cmd/zap/structs.go**: Shared data structures - Context, CtxWrapper for middleware

## Conventions

### Naming
- **Files**: Lowercase, descriptive (config.go, text.go, web.go)
- **Functions**: camelCase, verb-noun pattern (getURL, loadConfig, watchConfig)
- **Types**: PascalCase (Context, CtxWrapper)
- **Config keywords**: lowercase (expand, query, port, ssl_off, schema)

### Error Handling
- Return errors up the stack, don't panic
- Use hashicorp/go-multierror for aggregating multiple validation errors
- Log errors but continue with old config if reload fails (graceful degradation)
- HTTP handlers return appropriate status codes (404 for missing shortcuts)

### State Management
- All mutable state in `Context` struct
- `sync.RWMutex` protects config during hot reload
- Handlers acquire `RLock()` for read access, config reload acquires `Lock()` for write
- No global variables except Context instance

### Testing
- **BDD style**: smartystreets/goconvey for readable test descriptions
- **Filesystem abstraction**: afero.MemMapFs for isolated config tests
- **Table-driven tests**: text_test.go uses test tables for URL expansion cases
- **Integration tests**: web_test.go uses httptest for full request/response cycles
- **E2E tests**: e2e.sh for bash-level verification of running service
- **Race detector**: All tests run with `-race` flag in CI

## Critical Paths

### URL Expansion (Core Feature)
**File**: cmd/zap/text.go:getURL
**Flow**:
1. Split incoming path by "/" into tokens
2. Start with root config container
3. For each token:
   - Check for special keywords (expand, query, port, ssl_off, schema, *)
   - Navigate to child container if exists
   - Accumulate URL parts
4. Apply schema (https:// default, http:// if ssl_off, custom if schema)
5. Construct final URL with all accumulated parts
6. Return expanded URL

**Gotchas**:
- YAML reserved words (on, off, yes, no) must be quoted in config
- Wildcard (*) preserves remaining path segments without config lookup
- Query expansion removes trailing "/" for search compatibility

### Hot Reload (Reliability Feature)
**File**: cmd/zap/config.go:watchConfig
**Flow**:
1. fsnotify.Watcher monitors config file
2. On Write event, trigger reload
3. Parse new YAML → JSON → gabs.Container
4. Validate new config (check for errors)
5. If valid: acquire write lock, swap config, release lock
6. If invalid: log error, keep old config (fail-safe)
7. Update /etc/hosts if permissions allow

**Gotchas**:
- Editors may trigger multiple write events (debouncing may be needed)
- Parse errors don't crash service, old config remains active
- /etc/hosts updates require root, silently skip if permission denied

### HTTP Request Handling
**File**: cmd/zap/web.go + cmd/main.go
**Flow**:
1. Client sends request: GET /g/zap/issues/42
2. httprouter matches catch-all route
3. Handler acquires read lock on config
4. Extract path from request (handle X-Forwarded-Host header)
5. Call getURL(config, path)
6. Release read lock
7. If expansion succeeds: HTTP 302 redirect to expanded URL
8. If fails: HTTP 404 with error message

**Gotchas**:
- X-Forwarded-Host header allows reverse proxy setups
- Must release lock before sending response (avoid deadlock)
- 302 (not 301) to avoid browser caching

## Anti-Patterns

### Don't: Use Global Mutable State
Instead: Pass Context explicitly to functions that need it

### Don't: Panic on Config Errors
Instead: Return errors, log them, keep old config working

### Don't: Block on File I/O in HTTP Handlers
Instead: Hot reload runs in separate goroutine, handlers read cached config

### Don't: Expose Service to Public Internet Without Auth
Instead: Bind to localhost or use reverse proxy with auth

### Don't: Use 301 Permanent Redirects
Instead: Use 302 temporary redirects since shortcuts change

### Don't: Parse Config on Every Request
Instead: Parse once, cache in memory, reload only on file changes
