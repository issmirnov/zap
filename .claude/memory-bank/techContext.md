# Technical Context

> **Derived from**: projectbrief.md
> **Purpose**: Documents technologies, tools, and technical setup

## Technology Stack

### Core Technologies
- **Language**: Go 1.24+
- **Runtime**: Native compiled binary (no runtime dependencies)
- **HTTP Router**: julienschmidt/httprouter (zero-allocation router)
- **Database**: None (config file is state)

### Major Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/julienschmidt/httprouter | v1.3.0+ | High-performance HTTP routing |
| github.com/Jeffail/gabs/v2 | v2.7.0+ | JSON parsing and path traversal |
| github.com/ghodss/yaml | v1.0.0+ | YAML to JSON conversion |
| github.com/fsnotify/fsnotify | v1.8.0+ | File system event watching |
| github.com/spf13/afero | v1.15.0+ | Filesystem abstraction for testing |
| github.com/hashicorp/go-multierror | v1.1.1+ | Error aggregation for validation |
| github.com/smartystreets/goconvey | Latest | BDD-style testing framework |

### Development Tools
- **Build System**: `go build` (standard Go toolchain)
- **Package Manager**: Go modules (go.mod)
- **Linter**: golangci-lint (comprehensive Go linter)
- **Testing**: `go test` with smartystreets/goconvey
- **Release**: GoReleaser (cross-platform binary builds)
- **CI/CD**: GitHub Actions

## Development Setup

### Prerequisites
```bash
# Required installations
go 1.24+          # Go toolchain
make (optional)   # For build automation
```

### Installation
```bash
# Steps to set up development environment
git clone https://github.com/issmirnov/zap
cd zap2

# Download dependencies
go mod download

# Verify setup
go test ./...
```

### Configuration
- **Environment Variables**:
  - No environment variables required
  - All configuration is in `c.yml` file

- **Config Files**:
  - `c.yml`: Main configuration file with shortcut definitions
  - Default locations:
    - macOS: `/usr/local/etc/zap/c.yml`
    - Linux: `/etc/zap/c.yml`
    - Custom: Use `-config` flag

### Running Locally
```bash
# Build binary
go build -o zap cmd/main.go

# Run with default config
sudo ./zap  # Requires port 80 access

# Run on non-privileged port
./zap -host 127.0.0.1:8927

# Run with custom config
./zap -config /path/to/c.yml

# Validate config without starting server
./zap -validate

# Run tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run E2E tests (requires running server)
./e2e.sh

# Run linter
golangci-lint run
```

## Technical Constraints

### Performance Requirements
- **Target QPS**: 150,000+ requests per second
- **Latency**: <15ms per request under load
- **Memory**: Minimal footprint (stateless service, config cached in memory)
- **Startup Time**: Near-instant (<100ms)

### Platform Support
- **Operating Systems**:
  - Linux (amd64, 386, arm64, arm)
  - macOS (Darwin) (amd64, arm64)
  - Potential Windows support (not officially tested)
- **Deployment Modes**:
  - Standalone (requires port 80 access)
  - Behind reverse proxy (nginx, Caddy)
  - systemd service (Linux)
  - launchd service (macOS)
  - Network-wide with dnsmasq

### Security Requirements
- **Default**: localhost-only binding (127.0.0.1)
- **No Authentication**: Service trusts all requests (local-first design)
- **Config File**: Must be readable by zap process
- **/etc/hosts**: Requires root for modification (optional feature)
- **Public Deployment**: NOT recommended without reverse proxy + auth

## Infrastructure

### Deployment

- **Installation Methods**:
  1. **Homebrew** (macOS):
     ```bash
     brew install issmirnov/apps/zap
     brew services start zap
     ```
  2. **Ansible** (cross-platform):
     ```bash
     ansible-galaxy install issmirnov.zap
     ```
  3. **Manual** (Linux with systemd):
     ```bash
     # Download binary from GitHub releases
     # Copy to /usr/local/bin/zap
     # Install systemd unit file
     systemctl enable --now zap
     ```
  4. **From Source**:
     ```bash
     go install github.com/issmirnov/zap/cmd@latest
     ```

- **CI/CD**: GitHub Actions
  - **Workflow**: `.github/workflows/ci.yml`
    - Format check
    - Module download
    - go vet
    - Unit tests (with race detector, coverage)
    - golangci-lint
    - E2E tests
    - GoReleaser dry-run
  - **Release Workflow**: `.github/workflows/release.yml`
    - Triggered on version tags (`v*.*.*`)
    - Cross-platform builds (GoReleaser)
    - GitHub releases with checksums
    - Homebrew formula update (if configured)

- **Environments**: Single production deployment (self-hosted)

### Monitoring
- **Logging**: Stdout/stderr (captured by systemd/launchd)
- **Health Check**: `GET /healthz` endpoint (returns 200 OK)
- **Config Inspection**: `GET /varz` endpoint (returns JSON config)
- **Error Tracking**: No external service, logs to stdout
- **Metrics**: Basic (future enhancement: /metrics endpoint)

## External Integrations

### None at Runtime
- Zap has zero external dependencies at runtime
- All functionality is self-contained in single binary
- Config file is the only external resource

### Development/Build Integrations

#### GitHub
- **Purpose**: Source control, issue tracking, CI/CD
- **Authentication**: GITHUB_TOKEN for Actions
- **Workflows**: ci.yml, release.yml

#### GoReleaser
- **Purpose**: Cross-platform binary builds and releases
- **Configuration**: goreleaser.yml
- **Platforms**: Linux (multiple arch), macOS (Intel + ARM)

#### Homebrew
- **Purpose**: Package distribution for macOS
- **Tap**: issmirnov/apps
- **Formula**: Auto-updated on release

#### Ansible Galaxy
- **Purpose**: Automated deployment
- **Role**: issmirnov.zap
- **Installation**: `ansible-galaxy install issmirnov.zap`

## Technical Debt

### Minor Items
- **Config Validation**: Could be more comprehensive (e.g., URL format validation)
- **Metrics**: No built-in metrics endpoint (Prometheus-style)
- **Windows Support**: Not officially tested or documented
- **Config Reload Debouncing**: Multiple rapid writes may trigger multiple reloads
- **Structured Logging**: Uses plain log.Printf, could use structured logger

### Non-Issues (Intentional Choices)
- **No Database**: Config file is sufficient for use case
- **No Authentication**: Local-first design assumes trusted environment
- **No GUI**: YAML file is the interface by design

## Version History

### Current Version
- **Go Version**: 1.24 (recent upgrade from 1.23)
- **Major Dependencies**: Stable versions, regularly updated via Dependabot

### Build Configuration
- **CGO**: Disabled (`CGO_ENABLED=0`) for static binaries
- **Compilation Flags**: Standard optimizations
- **Module Mode**: Enabled (go.mod present)

### Recent Maintenance
- Regular dependency updates (Dependabot PRs)
- Go 1.24 compatibility maintained
- CI/CD pipeline stable and comprehensive
