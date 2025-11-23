# CloudBridge Go SDK

Official Go SDK for CloudBridge Global Network.

## Installation

```bash
go get github.com/twogc/cloudbridge-sdk/go/cloudbridge
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
    // Create client with authentication
    client, err := cloudbridge.NewClient(
        cloudbridge.WithToken("your-api-token"),
        cloudbridge.WithRegion("eu-central"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Connect to peer
    conn, err := client.Connect(ctx, "peer-id-123")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    log.Println("Connected successfully")
}
```

## Features

- P2P connection management
- Secure tunneling
- Mesh networking
- Automatic failover and reconnection
- JWT/OIDC authentication
- Multi-protocol support (QUIC, gRPC, WebSocket)
- Prometheus metrics
- Context-aware operations

## Configuration

### Environment Variables

```bash
CLOUDBRIDGE_TOKEN=your-api-token
CLOUDBRIDGE_REGION=eu-central
CLOUDBRIDGE_LOG_LEVEL=info
CLOUDBRIDGE_TIMEOUT=30s
```

### Programmatic Configuration

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("your-api-token"),
    cloudbridge.WithRegion("eu-central"),
    cloudbridge.WithTimeout(30 * time.Second),
    cloudbridge.WithLogLevel("debug"),
    cloudbridge.WithRetryPolicy(cloudbridge.RetryPolicy{
        MaxRetries: 3,
        InitialDelay: time.Second,
        MaxDelay: time.Minute,
    }),
)
```

## Usage Examples

### Basic P2P Connection

```go
// Connect to peer
conn, err := client.Connect(ctx, "peer-id")
if err != nil {
    return err
}
defer conn.Close()

// Send data
_, err = conn.Write([]byte("Hello, peer!"))
if err != nil {
    return err
}

// Read data
buf := make([]byte, 1024)
n, err := conn.Read(buf)
if err != nil {
    return err
}
log.Printf("Received: %s", buf[:n])
```

### TCP Tunnel

```go
// Create tunnel
tunnel, err := client.CreateTunnel(ctx, cloudbridge.TunnelConfig{
    LocalPort:  8080,
    RemotePeer: "peer-id",
    RemotePort: 3000,
    Protocol:   cloudbridge.ProtocolTCP,
})
if err != nil {
    return err
}
defer tunnel.Close()

log.Printf("Tunnel ready: localhost:8080 -> %s:3000", tunnel.RemotePeer())
```

### Mesh Network

```go
// Join mesh network
mesh, err := client.JoinMesh(ctx, "my-network")
if err != nil {
    return err
}
defer mesh.Leave()

// Broadcast message
err = mesh.Broadcast(ctx, []byte("Hello mesh!"))
if err != nil {
    return err
}

// Receive messages
for msg := range mesh.Messages() {
    log.Printf("Received from %s: %s", msg.From, msg.Data)
}
```

### Service Discovery

```go
// Register service
err := client.RegisterService(ctx, cloudbridge.ServiceConfig{
    Name: "my-api",
    Port: 8080,
    Tags: []string{"http", "api"},
})
if err != nil {
    return err
}

// Discover services
services, err := client.DiscoverServices(ctx, "my-api")
if err != nil {
    return err
}

for _, svc := range services {
    log.Printf("Found service: %s at %s:%d", svc.Name, svc.Address, svc.Port)
}
```

### Health Checks

```go
// Check connection health
health, err := client.Health(ctx)
if err != nil {
    return err
}

log.Printf("Status: %s, Latency: %v, Connected Peers: %d",
    health.Status, health.Latency, health.ConnectedPeers)
```

### Metrics

```go
// Get connection metrics
metrics, err := conn.Metrics()
if err != nil {
    return err
}

log.Printf("Bytes sent: %d, Bytes received: %d, RTT: %v",
    metrics.BytesSent, metrics.BytesReceived, metrics.RTT)
```

## Error Handling

```go
import "github.com/twogc/cloudbridge-sdk/go/cloudbridge/errors"

conn, err := client.Connect(ctx, "peer-id")
if err != nil {
    switch {
    case errors.IsAuthError(err):
        log.Println("Authentication failed")
    case errors.IsNetworkError(err):
        log.Println("Network error occurred")
    case errors.IsPeerNotFoundError(err):
        log.Println("Peer not found")
    case errors.IsTimeoutError(err):
        log.Println("Operation timed out")
    default:
        log.Printf("Unknown error: %v", err)
    }
    return err
}
```

## Advanced Features

### Custom Protocol Selection

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("token"),
    cloudbridge.WithProtocols(
        cloudbridge.ProtocolQUIC,
        cloudbridge.ProtocolGRPC,
        cloudbridge.ProtocolWebSocket,
    ),
)
```

### Connection Callbacks

```go
client.OnConnect(func(peer string) {
    log.Printf("Connected to peer: %s", peer)
})

client.OnDisconnect(func(peer string, err error) {
    log.Printf("Disconnected from peer: %s, error: %v", peer, err)
})

client.OnReconnect(func(peer string) {
    log.Printf("Reconnected to peer: %s", peer)
})
```

### Custom Middleware

```go
// Add custom middleware for logging
client.Use(func(next cloudbridge.Handler) cloudbridge.Handler {
    return func(ctx context.Context, req *cloudbridge.Request) (*cloudbridge.Response, error) {
        log.Printf("Request: %s", req.Method)
        resp, err := next(ctx, req)
        log.Printf("Response: %v, Error: %v", resp, err)
        return resp, err
    }
})
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

## Requirements

- Go 1.21 or higher
- CloudBridge API token

## Documentation

- [API Reference](https://pkg.go.dev/github.com/twogc/cloudbridge-sdk/go/cloudbridge)
- [Examples](./examples)
- [Architecture](../docs/ARCHITECTURE.md)

## Support

- GitHub Issues: https://github.com/twogc/cloudbridge-sdk/issues
- Documentation: https://docs.2gc.ru

## License

MIT License - see [LICENSE](../LICENSE) for details.
