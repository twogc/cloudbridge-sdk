# CloudBridge SDK - Quick Start Guide

Get started with CloudBridge SDK in 5 minutes!

## Installation

### Option 1: Go Module

```bash
go get github.com/twogc/cloudbridge-sdk/go
```

Then in your code:

```go
import "github.com/twogc/cloudbridge-sdk/go/cloudbridge"
```

### Option 2: Clone Repository

```bash
git clone https://github.com/twogc/cloudbridge-sdk.git
cd cloudbridge-sdk/go
go mod download
```

## Your First Program (2 minutes)

Create a file `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
    // 1. Create client
    client, err := cloudbridge.NewClient(
        cloudbridge.WithToken(os.Getenv("CLOUDBRIDGE_TOKEN")),
        cloudbridge.WithRegion("eu-central"),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    fmt.Println("[OK] Client created successfully!")

    // 2. Connect to a peer
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    conn, err := client.Connect(ctx, "your-peer-id")
    if err != nil {
        log.Printf("Connection: %v", err)
        return
    }
    defer conn.Close()

    fmt.Println("[OK] Connected to peer!")

    // 3. Send data
    message := []byte("Hello, peer!")
    n, err := conn.Write(message)
    fmt.Printf("[OK] Sent %d bytes\n", n)

    // 4. Receive data
    buffer := make([]byte, 1024)
    n, err = conn.Read(buffer)
    fmt.Printf("[OK] Received: %s\n", buffer[:n])
}
```

Run it:

```bash
export CLOUDBRIDGE_TOKEN="your-jwt-token"
go run main.go
```

## Common Tasks

### Task 1: Create a Client

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
    cloudbridge.WithRegion("eu-central"),      // Optional
    cloudbridge.WithTimeout(30*time.Second),   // Optional
)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
defer client.Close()
```

### Task 2: Connect to a Peer

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

conn, err := client.Connect(ctx, peerID)
if err != nil {
    log.Printf("Connection failed: %v", err)
    return
}
defer conn.Close()

fmt.Printf("Connected to: %s\n", conn.PeerID())
```

### Task 3: Send Data

```go
message := []byte("Hello!")
n, err := conn.Write(message)
if err != nil {
    log.Printf("Write failed: %v", err)
    return
}
fmt.Printf("Sent %d bytes\n", n)
```

### Task 4: Receive Data

```go
buffer := make([]byte, 4096)
n, err := conn.Read(buffer)
if err != nil {
    log.Printf("Read failed: %v", err)
    return
}
fmt.Printf("Received: %s\n", buffer[:n])
```

### Task 5: Discover Peers

```go
peers, err := client.DiscoverPeers(ctx)
if err != nil {
    log.Printf("Discovery failed: %v", err)
    return
}

for _, peer := range peers {
    fmt.Printf("Peer: %s (Region: %s)\n", peer.ID, peer.Region)
}
```

### Task 6: Join a Mesh Network

```go
mesh, err := client.JoinMesh(ctx, "my-network")
if err != nil {
    log.Printf("Join mesh failed: %v", err)
    return
}
defer mesh.Leave()

peers, err := mesh.Peers()
fmt.Printf("Peers in mesh: %v\n", peers)
```

### Task 7: Check Health

```go
health, err := client.Health(ctx)
if err != nil {
    log.Printf("Health check failed: %v", err)
    return
}
fmt.Printf("Client status: %+v\n", health)
```

### Task 8: Create a Tunnel

```go
tunnelConfig := cloudbridge.TunnelConfig{
    PeerID:      "peer-id",
    LocalAddr:   "localhost:8080",
    RemoteAddr:  "localhost:80",
}

tunnel, err := client.CreateTunnel(ctx, tunnelConfig)
if err != nil {
    log.Printf("Tunnel creation failed: %v", err)
    return
}
defer tunnel.Close()

fmt.Printf("Tunnel listening on: %s\n", tunnel.LocalAddr())
```

## Configuration

### Environment Variables

```bash
# Required
export CLOUDBRIDGE_TOKEN="your-jwt-token"

# Optional
export CLOUDBRIDGE_REGION="eu-central"           # Default: eu-central
export CLOUDBRIDGE_TIMEOUT="30s"                 # Default: 30s
export CLOUDBRIDGE_LOG_LEVEL="info"              # Default: info
```

### Programmatic Configuration

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("token"),
    cloudbridge.WithRegion("eu-central"),
    cloudbridge.WithTimeout(30*time.Second),
    cloudbridge.WithLogLevel("debug"),
    cloudbridge.WithInsecureSkipVerify(false),
)
```

## Examples

Ready to dive deeper? Check out our examples:

- **[Simple Connection](go/examples/simple_connection/)** - Basic peer-to-peer
- **[Echo Server](go/examples/echo_server/)** - Server implementation
- **[Mesh Chat](go/examples/mesh_chat/)** - Group communication

Run them:

```bash
cd go/examples/simple_connection
go run main.go
```

## Learning Path

1. **Start here** - Read this Quick Start
2. **Try examples** - Run the examples
3. **Read API Reference** - [docs/API_REFERENCE.md](docs/API_REFERENCE.md)
4. **Understand architecture** - [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
5. **Explore authentication** - [docs/AUTHENTICATION.md](docs/AUTHENTICATION.md)
6. **Build your app** - Use what you learned

## FAQ

### Q: How do I get a JWT token?

**A:** You need an auth provider (Zitadel, Keycloak, etc.). Generate a JWT with these claims:

```json
{
  "sub": "your-user-id",
  "tenant_id": "your-tenant-id",
  "iat": 1735603200,
  "exp": 1735689600
}
```

### Q: Do I need to run a relay server?

**A:** For testing and learning: No, SDK has mock implementations.
For production: Yes, you need a CloudBridge relay server.

### Q: Can I use this with my existing Go app?

**A:** Yes! The SDK is just a library. Import it and use it in your code:

```go
import "github.com/twogc/cloudbridge-sdk/go/cloudbridge"
```

### Q: How do I handle errors?

**A:** Check for errors on all operations:

```go
conn, err := client.Connect(ctx, peerID)
if err != nil {
    // Handle error
    log.Printf("Error: %v", err)
}
```

### Q: How do I close connections properly?

**A:** Always use `defer`:

```go
client, _ := cloudbridge.NewClient(...)
defer client.Close()

conn, _ := client.Connect(ctx, peerID)
defer conn.Close()
```

### Q: Can I use this in production?

**A:** The SDK is currently in Alpha (v0.1.0). Some features are not yet implemented.
Check [SDK_STATUS.md](SDK_STATUS.md) for current state.

### Q: What's the difference between connection, tunnel, and mesh?

**A:**
- **Connection** - Direct point-to-point communication with one peer
- **Tunnel** - Forward traffic from local port to remote peer's port
- **Mesh** - Join a network where you can communicate with multiple peers

### Q: How do I contribute?

**A:** See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Troubleshooting

### Issue: "token is required"

```bash
export CLOUDBRIDGE_TOKEN="your-token"
go run main.go
```

### Issue: "connection failed: not implemented"

This is expected without a relay server. In production, ensure:
- CloudBridge relay server is running
- Token is valid
- Peer ID exists in your network

### Issue: "context deadline exceeded"

Increase timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
defer cancel()
```

### Issue: "memory leak" or "too many goroutines"

Ensure you close resources:

```go
defer client.Close()
defer conn.Close()
defer mesh.Leave()
```

## Documentation

- **[API Reference](docs/API_REFERENCE.md)** - Complete API documentation
- **[Architecture](docs/ARCHITECTURE.md)** - System design and internals
- **[Authentication](docs/AUTHENTICATION.md)** - Auth and security details
- **[CLI Tool](docs/CLI.md)** - Using the CLI tool
- **[SDK Status](SDK_STATUS.md)** - Current implementation status

## Resources

- **GitHub:** https://github.com/twogc/cloudbridge-sdk
- **Issues:** https://github.com/twogc/cloudbridge-sdk/issues
- **Discussions:** https://github.com/twogc/cloudbridge-sdk/discussions

## Next Steps

1. **Clone examples:** `git clone https://github.com/twogc/cloudbridge-sdk.git`
2. **Run a sample:** `cd go/examples/simple_connection && go run main.go`
3. **Read docs:** Check [docs/](docs/) directory
4. **Build something:** Use the examples as a template
5. **Share feedback:** We'd love to hear from you!

---

Happy coding!
