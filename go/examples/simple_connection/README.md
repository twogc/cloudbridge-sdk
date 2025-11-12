# Simple Connection Example

This example demonstrates how to create a basic P2P connection using the CloudBridge SDK.

## What This Example Shows

1. **Client Creation** - How to initialize a CloudBridge client with configuration
2. **Peer Connection** - How to connect to a peer using their peer ID
3. **Data Transfer** - How to send and receive data over the connection
4. **Connection Metrics** - How to retrieve connection statistics
5. **Health Checks** - How to monitor client health

## Running the Example

```bash
cd examples/simple_connection
go run main.go
```

## Expected Output

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

## Code Walkthrough

### 1. Create Client

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("your-jwt-token"),
    cloudbridge.WithRegion("eu-central"),
    cloudbridge.WithTimeout(30*time.Second),
)
```

**Configuration Options:**
- `WithToken()` - JWT authentication token from your auth provider
- `WithRegion()` - CloudBridge region (eu-central, us-east, etc.)
- `WithTimeout()` - Connection timeout duration
- `WithLogLevel()` - Logging level (debug, info, warn, error)

### 2. Connect to Peer

```go
conn, err := client.Connect(ctx, peerID)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()
```

**Peer ID** can be obtained from:
- Peer discovery (`client.DiscoverPeers()`)
- User input
- Configuration file
- Service registry

### 3. Send Data

```go
message := []byte("Hello, peer!")
n, err := conn.Write(message)
fmt.Printf("Sent %d bytes\n", n)
```

The `Connection` interface implements `io.ReadWriter`, so you can use standard Go I/O operations.

### 4. Receive Data

```go
buffer := make([]byte, 1024)
n, err := conn.Read(buffer)
fmt.Printf("Received: %s\n", buffer[:n])
```

### 5. Get Metrics

```go
metrics, err := conn.Metrics()
fmt.Printf("RTT: %s\n", metrics.RTT)
fmt.Printf("Bytes sent: %d\n", metrics.BytesSent)
fmt.Printf("Bytes received: %d\n", metrics.BytesReceived)
```

## Production Setup

To use this in production:

1. **Set up CloudBridge Relay Server**
   ```bash
   # Deploy relay server to your infrastructure
   docker run -p 5550:5550 cloudbridge/relay
   ```

2. **Configure Authentication**
   - Set up Zitadel or your OIDC provider
   - Generate JWT tokens for your users
   - Include `tenant_id` in token claims

3. **Update Configuration**
   ```go
   client, err := cloudbridge.NewClient(
       cloudbridge.WithToken(os.Getenv("CLOUDBRIDGE_TOKEN")),
       cloudbridge.WithRegion(os.Getenv("CLOUDBRIDGE_REGION")),
   )
   ```

4. **Discover Peers**
   ```go
   peers, err := client.DiscoverPeers(ctx)
   for _, peer := range peers {
       fmt.Printf("Found peer: %s\n", peer.ID)
   }
   ```

5. **Handle Errors**
   ```go
   conn, err := client.Connect(ctx, peerID)
   if err != nil {
       if errors.Is(err, cloudbridge.ErrPeerNotFound) {
           // Peer is offline or doesn't exist
       }
       if errors.Is(err, cloudbridge.ErrTimeout) {
           // Connection timed out
       }
       return err
   }
   ```

## Common Issues

### "token is required"
**Solution:** Provide a valid JWT token via `WithToken()` or `CLOUDBRIDGE_TOKEN` environment variable.

### "not implemented"
**Solution:** The SDK requires a running CloudBridge relay server. This example runs without one to demonstrate the API.

### "connection timeout"
**Solution:** Increase timeout with `WithTimeout()` or check network connectivity.

## Next Steps

- See [Echo Server Example](../echo_server/) for a complete server implementation
- See [Tunnel Example](../tunnel/) for port forwarding
- See [Mesh Chat Example](../mesh_chat/) for multi-peer communication
- Read the [API Reference](../../docs/API_REFERENCE.md)

## Learn More

- [CloudBridge Documentation](https://github.com/twogc/cloudbridge-sdk)
- [Authentication Guide](../../docs/AUTHENTICATION.md)
- [Architecture Overview](../../docs/ARCHITECTURE.md)
