# CloudBridge SDK Examples

This directory contains practical examples demonstrating how to use the CloudBridge SDK for different use cases.

## Examples Overview

### 1. Simple Connection ([simple_connection/](simple_connection/))

**Demonstrates:** Basic peer-to-peer connection and data transfer

```bash
cd simple_connection && go run main.go
```

**What you'll learn:**
- How to create a CloudBridge client
- How to connect to a peer
- How to send and receive data
- How to retrieve connection metrics
- How to check client health

**Use case:** Simple point-to-point communication

### 2. Echo Server ([echo_server/](echo_server/))

**Demonstrates:** Server implementation that echoes received data

```bash
cd echo_server && go run main.go
```

**What you'll learn:**
- How to set up connection callbacks (onConnect, onDisconnect, onReconnect)
- How to handle multiple concurrent connections
- How to implement server logic
- How to manage connection lifecycle
- How to implement graceful shutdown

**Use case:** Simple server applications

### 3. Mesh Chat ([mesh_chat/](mesh_chat/))

**Demonstrates:** Multi-peer communication in a mesh network

```bash
cd mesh_chat && go run main.go
```

**What you'll learn:**
- How to join a mesh network
- How to discover peers in the network
- How to broadcast messages to all peers
- How to handle peer list
- How to implement chat-like applications

**Use case:** Group communication, broadcast messaging

## Running the Examples

### Prerequisites

1. **Go 1.25.3+** installed
2. **CLOUDBRIDGE_TOKEN** environment variable (optional, defaults to example token)
3. **CloudBridge Relay Server** (for real connections, not required for API demonstration)

### Example 1: Simple Connection

```bash
cd simple_connection
export CLOUDBRIDGE_TOKEN="your-jwt-token"
go run main.go
```

Expected output:
```
CloudBridge SDK - Simple Connection Example
============================================
✓ Client created successfully

Example 1: Connecting to a peer...
⚠ Connection failed (expected without relay): not implemented
   In production, this would establish a real P2P connection

Example 2: Checking client health...
✓ Client health: &{...}

Example completed!
```

### Example 2: Echo Server

```bash
cd echo_server
export CLOUDBRIDGE_TOKEN="your-jwt-token"
go run main.go
```

Expected output:
```
CloudBridge SDK - Echo Server Example
=====================================

✓ Echo server initialized
  Waiting for peer connections...
  Press Ctrl+C to stop
```

### Example 3: Mesh Chat

```bash
cd mesh_chat
export CLOUDBRIDGE_TOKEN="your-jwt-token"
export CLOUDBRIDGE_USER="alice"
go run main.go
```

Expected output:
```
CloudBridge SDK - Mesh Chat Example
====================================

✓ Chat client created
  Username: alice

Joining mesh network: chat-network
⚠ Failed to join mesh (expected without relay): not implemented
  In production, the client would join the mesh network

=== Chat Interface ===
Commands:
  /peers         - List connected peers
  /send <msg>    - Broadcast message to all peers
  /quit          - Exit chat

Type your message and press Enter to broadcast:

>
```

## Environment Variables

Configure examples with environment variables:

```bash
# Authentication
export CLOUDBRIDGE_TOKEN="your-jwt-token"

# Network
export CLOUDBRIDGE_REGION="eu-central"           # Default: eu-central
export CLOUDBRIDGE_TIMEOUT="30s"                 # Default: 30s

# User Information
export CLOUDBRIDGE_USER="alice"                  # For chat example

# Debugging
export CLOUDBRIDGE_LOG_LEVEL="debug"             # Default: info
```

## Configuration

### Token Format

CloudBridge uses JWT (JSON Web Token) for authentication. The token should include:

```json
{
  "sub": "user-id",
  "tenant_id": "tenant-id",
  "iat": 1735603200,
  "exp": 1735689600
}
```

Generate a valid token from your auth provider (Zitadel, Keycloak, etc.)

### Peer Discovery

To find peer IDs:

```go
peers, err := client.DiscoverPeers(ctx)
for _, peer := range peers {
    fmt.Printf("Peer: %s (Region: %s, Latency: %d ms)\n",
        peer.ID, peer.Region, peer.LatencyMs)
}
```

## Common Patterns

### Pattern 1: Client with Error Handling

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
    cloudbridge.WithRegion(region),
    cloudbridge.WithTimeout(timeout),
)
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
defer client.Close()
```

### Pattern 2: Connection with Cleanup

```go
conn, err := client.Connect(ctx, peerID)
if err != nil {
    return err
}
defer conn.Close()

// Use connection
data := make([]byte, 1024)
n, err := conn.Read(data)
```

### Pattern 3: Context with Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

conn, err := client.Connect(ctx, peerID)
if err != nil {
    return err
}
```

### Pattern 4: Graceful Shutdown

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

<-sigChan
fmt.Println("Shutting down...")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Close resources
client.Close()
```

## Troubleshooting

### Issue: "token is required"
**Solution:** Set `CLOUDBRIDGE_TOKEN` environment variable or pass token to `WithToken()`

### Issue: "not implemented"
**Solution:** This means relay server integration is not active. In production with relay running, this would work.

### Issue: "connection timeout"
**Solution:** Increase timeout with `WithTimeout()` or check network connectivity

### Issue: "peer not found"
**Solution:** Use `client.DiscoverPeers()` to find available peers first

## Best Practices

1. **Always use defer client.Close()** to cleanup resources
2. **Use context with timeout** for all operations
3. **Handle all errors** properly
4. **Use callbacks** for event-driven code
5. **Reuse connections** when possible
6. **Monitor metrics** for debugging
7. **Set appropriate log levels** for production
8. **Use environment variables** for configuration

## Performance Tips

1. **Connection pooling** - Reuse connections when possible
2. **Batch operations** - Send multiple messages together
3. **Monitor metrics** - Track latency and throughput
4. **Adjust buffer sizes** - Based on message sizes
5. **Use goroutines** - For concurrent operations
6. **Enable compression** - For large messages

## Next Steps

- Read the [API Reference](../docs/API_REFERENCE.md) for complete API documentation
- Check [Architecture](../docs/ARCHITECTURE.md) for system design
- See [Authentication](../docs/AUTHENTICATION.md) for security details
- Review [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines

## Advanced Examples

Looking for more advanced examples?

- **File Transfer** - See `file_transfer/` (coming soon)
- **Service Discovery** - See `service_discovery/` (coming soon)
- **Tunnel** - See `tunnel/` (coming soon)
- **Monitoring** - See `monitoring/` (coming soon)

## Support

- **Issues:** Report issues on [GitHub Issues](https://github.com/twogc/cloudbridge-sdk/issues)
- **Documentation:** Read the [full documentation](../docs/)
- **Community:** Join our community discussions

## License

All examples are licensed under MIT License. See [LICENSE](../LICENSE) for details.
