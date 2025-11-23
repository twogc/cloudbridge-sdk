# Echo Server Example

This example demonstrates how to build a simple echo server using CloudBridge SDK that listens for peer connections and echoes back received data.

## What This Example Shows

1. **Server Initialization** - How to set up a server with connection callbacks
2. **Connection Callbacks** - How to handle `onConnect`, `onDisconnect`, and `onReconnect` events
3. **Echo Logic** - How to handle incoming data and echo it back
4. **Graceful Shutdown** - How to properly close connections and cleanup resources
5. **Concurrent Connections** - How to handle multiple peers simultaneously

## Running the Example

```bash
cd examples/echo_server
go run main.go
```

## Expected Output

```
CloudBridge SDK - Echo Server Example
=====================================

✓ Echo server initialized
  Waiting for peer connections...
  Press Ctrl+C to stop

=== Echo Logic Demo ===

Handling connection from: example-peer-123
  [1] Received 14 bytes: Hello, server!
  [1] Echoing back 14 bytes
  [2] Received 16 bytes: Test message 2
  [2] Echoing back 16 bytes
  [3] Received 8 bytes: Goodbye!
  [3] Echoing back 8 bytes
Connection from example-peer-123 handled successfully

To connect to this echo server from another peer:
  1. Get this peer's ID from discovery
  2. Use client.Connect(ctx, peerID)
  3. Write data with conn.Write()
  4. Read echoed data with conn.Read()
```

## Code Structure

### Echo Handler Function

The core echo logic is in the `echoHandler()` function:

```go
func echoHandler(conn io.ReadWriteCloser, peerID string) {
    defer conn.Close()

    buffer := make([]byte, 4096)

    for {
        // Read data
        n, err := conn.Read(buffer)
        if err != nil {
            break
        }

        // Echo back
        conn.Write(buffer[:n])
    }
}
```

### Connection Callbacks

Set up handlers for connection events:

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
    cloudbridge.WithOnConnect(func(peerID string) {
        fmt.Printf("Peer connected: %s\n", peerID)
        // Start echo handler for this peer
    }),
    cloudbridge.WithOnDisconnect(func(peerID string, err error) {
        fmt.Printf("Peer disconnected: %s\n", peerID)
        // Cleanup resources
    }),
    cloudbridge.WithOnReconnect(func(peerID string) {
        fmt.Printf("Peer reconnected: %s\n", peerID)
    }),
)
```

## Production Implementation

For a production echo server, you would:

### 1. Accept Incoming Connections

```go
// Listen for connections (pseudo-code)
listener, err := client.Listen()
if err != nil {
    log.Fatalf("Failed to listen: %v", err)
}

for {
    conn, err := listener.Accept()
    if err != nil {
        log.Printf("Accept error: %v", err)
        continue
    }

    // Handle connection in goroutine
    go echoHandler(conn, conn.PeerID())
}
```

### 2. Handle Multiple Concurrent Connections

```go
type EchoServer struct {
    client      *cloudbridge.Client
    connections map[string]io.ReadWriteCloser
    mu          sync.RWMutex
}

func (s *EchoServer) handleConnection(conn io.ReadWriteCloser) {
    peerID := conn.PeerID()

    // Track connection
    s.mu.Lock()
    s.connections[peerID] = conn
    s.mu.Unlock()

    defer func() {
        s.mu.Lock()
        delete(s.connections, peerID)
        s.mu.Unlock()
        conn.Close()
    }()

    s.echoHandler(conn)
}
```

### 3. Add Graceful Shutdown

```go
func (s *EchoServer) Shutdown(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Close all connections
    for peerID, conn := range s.connections {
        fmt.Printf("Closing connection with %s\n", peerID)
        conn.Close()
    }

    s.connections = make(map[string]io.ReadWriteCloser)
    return s.client.Close()
}
```

### 4. Add Monitoring and Metrics

```go
type EchoStats struct {
    TotalConnections    int64
    ActiveConnections   int64
    BytesReceived       int64
    BytesSent           int64
}

func (s *EchoServer) Stats() EchoStats {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return EchoStats{
        ActiveConnections: int64(len(s.connections)),
    }
}
```

## Client-Side Usage

To connect to this echo server from another peer:

```go
client, _ := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
)

conn, _ := client.Connect(ctx, serverPeerID)
defer conn.Close()

// Send data
message := []byte("Hello, echo server!")
conn.Write(message)

// Receive echo
response := make([]byte, len(message))
n, _ := conn.Read(response)
fmt.Printf("Echo response: %s\n", response[:n])
```

## Performance Considerations

1. **Buffer Size** - Use appropriate buffer size (4KB is typical)
   ```go
   buffer := make([]byte, 4096)
   ```

2. **Concurrent Goroutines** - Limit concurrent connections
   ```go
   semaphore := make(chan struct{}, maxConnections)
   ```

3. **Timeouts** - Set read/write deadlines
   ```go
   conn.SetReadDeadline(time.Now().Add(30 * time.Second))
   ```

4. **Memory Management** - Reuse buffers with sync.Pool
   ```go
   bufferPool := sync.Pool{
       New: func() interface{} {
           return make([]byte, 4096)
       },
   }
   ```

## Testing

Test the echo server:

```bash
# Terminal 1: Start echo server
go run main.go

# Terminal 2: Connect and test
go run ../simple_connection/main.go
```

## Troubleshooting

### Server doesn't accept connections
- Ensure relay server is running
- Check that firewall allows P2P connections
- Verify token is valid

### Data not being echoed
- Check that Read() is not blocking indefinitely
- Verify buffer size is sufficient
- Add logging to see what data is received

### Memory leak
- Ensure all connections are properly closed
- Use `defer conn.Close()` in handlers
- Remove entries from connection map on disconnect

## Next Steps

- See [Tunnel Example](../tunnel/) for port forwarding
- See [Mesh Chat Example](../mesh_chat/) for multi-peer messaging
- Read the [API Reference](../../docs/API_REFERENCE.md)
- Explore [Performance Tuning](../../docs/ARCHITECTURE.md#performance)
