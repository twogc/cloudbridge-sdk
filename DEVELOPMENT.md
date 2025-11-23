# CloudBridge SDK - Development Guide

This guide is for developers working on the CloudBridge SDK itself.

## Project Structure

```
cloudbridge-sdk/
├── go/                           # Go SDK implementation
│   ├── cloudbridge/              # Core SDK package
│   │   ├── client.go             # Main Client interface
│   │   ├── config.go             # Configuration
│   │   ├── connection.go         # Connection interface
│   │   ├── tunnel.go             # Tunnel interface
│   │   ├── mesh.go               # Mesh networking
│   │   ├── service.go            # Service discovery
│   │   ├── errors/               # Error types
│   │   ├── internal/
│   │   │   ├── bridge/           # Relay client integration
│   │   │   └── jwt/              # JWT token handling
│   │   └── *_test.go             # Unit tests
│   ├── cmd/cloudbridge/          # CLI application
│   ├── examples/                 # Usage examples
│   └── Makefile                  # Build automation
│
├── python/                       # Python SDK (coming soon)
├── javascript/                   # JavaScript SDK (coming soon)
│
├── docs/                         # Documentation
│   ├── API_REFERENCE.md          # API documentation
│   ├── ARCHITECTURE.md           # System design
│   ├── AUTHENTICATION.md         # Auth details
│   ├── INTEGRATION.md            # Relay integration
│   └── CLI.md                    # CLI documentation
│
├── QUICKSTART.md                 # Quick start guide
├── SDK_STATUS.md                 # Implementation status
├── CONTRIBUTING.md               # Contribution guidelines
└── README.md                     # Project overview
```

## Development Setup

### Prerequisites

- **Go 1.25.3+**
- **Git**
- **Make** (optional, for Makefile)

### Clone Repository

```bash
git clone https://github.com/twogc/cloudbridge-sdk.git
cd cloudbridge-sdk/go
```

### Install Dependencies

```bash
go mod download
go mod verify
```

### Verify Setup

```bash
go build ./...
go test ./...
```

## Architecture Overview

### Core Components

1. **Client** - Main entry point for SDK users
   - Manages connections
   - Handles callbacks
   - Coordinates with transport layer

2. **Connection** - P2P connection to a peer
   - Implements `io.ReadWriteCloser`
   - Tracks metrics
   - Handles data transfer

3. **Tunnel** - TCP/UDP port forwarding
   - Local listener
   - Remote forwarding
   - Bidirectional proxying

4. **Mesh** - Multi-peer networking
   - Peer discovery
   - Broadcast messaging
   - Network topology

5. **Bridge** - Integration with CloudBridge Relay
   - Wraps relay client
   - Provides P2P connectivity
   - Manages streams

## Code Style Guide

### Naming Conventions

```go
// Exported (public)
type Client struct { ... }
func (c *Client) Connect(ctx context.Context, peerID string) (Connection, error) { ... }

// Unexported (internal)
type connection struct { ... }
func (c *connection) dial(ctx context.Context) error { ... }
```

### Error Handling

```go
// Always return error as last value
func (c *Client) Connect(ctx context.Context, peerID string) (Connection, error)

// Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
}

// Use custom error types for important errors
if errors.Is(err, ErrPeerNotFound) {
    // Handle specific case
}
```

### Documentation

```go
// Package comment for each file
// connection.go implements the Connection interface for CloudBridge SDK

// Type comment
// Connection represents a P2P connection to a peer
type Connection interface {
    // Method comment
    // Read reads data from the peer connection
    Read(b []byte) (n int, err error)
}
```

### Testing

```go
// Test function naming
func TestClientConnect(t *testing.T) { ... }
func TestClientConnect_Success(t *testing.T) { ... }
func TestClientConnect_Timeout(t *testing.T) { ... }

// Table-driven tests preferred
func TestConnection(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    interface{}
        wantErr bool
    }{
        {
            name:    "success case",
            input:   "test",
            want:    "expected",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test code
        })
    }
}
```

## Testing

### Run All Tests

```bash
make test
```

### Run Specific Tests

```bash
go test -v ./cloudbridge -run TestClientConnect
go test -v ./cloudbridge/internal/jwt -run TestParseToken
```

### Test Coverage

```bash
make test-coverage
# Opens coverage.html in browser
```

### Add New Tests

1. Create `*_test.go` file in same package
2. Use table-driven test pattern
3. Test both success and error cases
4. Aim for >80% coverage

Example:

```go
func TestNewClient(t *testing.T) {
    tests := []struct {
        name    string
        opts    []Option
        wantErr bool
    }{
        {
            name:    "valid config",
            opts:    []Option{WithToken("token")},
            wantErr: false,
        },
        {
            name:    "missing token",
            opts:    []Option{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := NewClient(tt.opts...)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
            }
            if client != nil {
                client.Close()
            }
        })
    }
}
```

## Building

### Build SDK Library

```bash
make build
```

### Build CLI Tool

```bash
make cli
```

### Build Everything

```bash
make all
```

### Build for Specific Platform

```bash
GOOS=linux GOARCH=amd64 go build ./cmd/cloudbridge
GOOS=darwin GOARCH=arm64 go build ./cmd/cloudbridge
GOOS=windows GOARCH=amd64 go build ./cmd/cloudbridge
```

## Debugging

### Enable Debug Logging

```bash
export CLOUDBRIDGE_LOG_LEVEL=debug
go run main.go
```

### Using Delve Debugger

```bash
go install github.com/go-delve/delve/cmd/dlv@latest

dlv debug ./cmd/cloudbridge
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print variableName
```

### Print Debugging

```go
import "log"

log.Printf("Debug: %+v", variable)
```

## Documentation Generation

### Generate Go Documentation

```bash
godoc -h localhost:6060 &
# Visit http://localhost:6060/pkg/github.com/twogc/cloudbridge-sdk/go/cloudbridge/
```

### Update Documentation

1. Update comments in code
2. Update markdown files in `docs/`
3. Run tests to ensure examples work

## Publishing

### Version Bumping

Update version in:
- `cmd/cloudbridge/main.go` - `Version = "0.x.x"`
- `go.mod` - Go version if needed
- `CHANGELOG.md` - Add entry

### Creating Release

```bash
git tag -a v0.2.0 -m "Release version 0.2.0"
git push origin v0.2.0
```

### Publishing to pkg.go.dev

Automatic when tag is pushed to GitHub.

## Contributing

### Code Review Checklist

- [ ] Code follows style guide
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Commits are atomic
- [ ] Commit messages are clear

### Pull Request Process

1. Fork repository
2. Create feature branch: `git checkout -b feature/my-feature`
3. Make changes
4. Write tests
5. Update documentation
6. Commit: `git commit -am "feat: add my feature"`
7. Push: `git push origin feature/my-feature`
8. Create Pull Request

### Commit Message Format

```
<type>: <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`

Example:
```
feat: add mesh peer discovery

Implement automatic peer discovery in mesh networks
using the relay server's peer list API.

Fixes #123
```

## Security Considerations

### Token Handling

- Never log tokens
- Use environment variables
- Validate token claims
- Check token expiry

```go
// Good
token := os.Getenv("CLOUDBRIDGE_TOKEN")

// Bad
token := "hardcoded-token"
log.Printf("Using token: %s", token)
```

### Connection Security

- Always verify peer identity
- Use TLS when available
- Implement rate limiting
- Validate input data

### Dependencies

- Keep dependencies updated
- Check for vulnerabilities
- Review dependency changes

```bash
go get -u
go mod verify
```

## Metrics and Monitoring

### Instrumentation Points

- Connection creation/closing
- Data transfer (bytes/messages)
- Error rates
- Latency measurements

### Adding Metrics

```go
// Example metric
var (
    connectionsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "sdk_connections_total",
            Help: "Total number of connections",
        },
        []string{"status"},
    )
)

// Usage
connectionsTotal.WithLabelValues("success").Inc()
```

## Development Priorities

### Current Sprint

1. **Connection Implementation** - Real P2P support
2. **Bridge Integration** - Full relay client integration
3. **Tunnel Functionality** - Port forwarding
4. **Comprehensive Tests** - Unit and integration tests

### Next Sprint

1. **Mesh Networking** - Multi-peer support
2. **Service Discovery** - Service registry
3. **Python SDK** - Python implementation
4. **JavaScript SDK** - JavaScript implementation

See [SDK_STATUS.md](SDK_STATUS.md) for detailed roadmap.

## Issue Reporting

When reporting issues:

1. **Provide minimal reproduction** - Code that reproduces the issue
2. **Include environment** - Go version, OS, SDK version
3. **Describe expected behavior** - What should happen
4. **Describe actual behavior** - What actually happens
5. **Attach logs** - With debug logging enabled

Example:

```
## Description
Connection times out when connecting to peer

## Steps to Reproduce
1. Create client with token
2. Call client.Connect() with peer ID
3. Wait 30 seconds

## Expected Behavior
Connection should succeed within 10 seconds

## Actual Behavior
Connection times out after 30 seconds

## Environment
- Go version: 1.25.3
- SDK version: 0.1.0
- OS: macOS 12.x
- Network: Fiber, low latency

## Logs
[Paste debug logs here]
```

## Getting Help

- **Questions:** Use GitHub Discussions
- **Bugs:** Open GitHub Issues
- **Features:** Start a Discussion first
- **Security:** Email security@2gc.dev

## References

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Test Best Practices](https://golang.org/doc/effective_go#testing)

## Quick Commands

```bash
# Development workflow
make test                          # Run tests
make test-coverage                # Test with coverage
make lint                          # Run linter
make fmt                           # Format code
make vet                           # Run go vet
make cli                           # Build CLI
make run-cli ARGS="health"         # Run CLI with args

# Building
make build                         # Build SDK
make release                       # Build release binaries
make clean                         # Clean build artifacts

# Other
make help                          # Show all targets
make tidy                          # Tidy dependencies
make dev-setup                     # Setup dev environment
```

---

Happy coding!

For questions or issues, please open a GitHub issue or discussion.
