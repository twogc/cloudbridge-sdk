# CloudBridge SDK API Reference

Version: 0.1.0
Last Updated: November 2025

## Table of Contents

- [Client](#client)
- [Connection](#connection)
- [Tunnel](#tunnel)
- [Mesh](#mesh)
- [Configuration](#configuration)
- [Errors](#errors)
- [Types](#types)

## Client

### NewClient

Creates a new CloudBridge client with the specified options.

```go
func NewClient(opts ...Option) (*Client, error)
```

**Parameters:**
- `opts` - Variadic configuration options

**Returns:**
- `*Client` - Configured client instance
- `error` - Configuration or initialization error

**Example:**
```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("token"),
    cloudbridge.WithRegion("eu-central"),
)
```

### Client.Connect

Establishes a P2P connection to the specified peer.

```go
func (c *Client) Connect(ctx context.Context, peerID string) (Connection, error)
```

**Parameters:**
- `ctx` - Context for cancellation and timeout
- `peerID` - Target peer identifier

**Returns:**
- `Connection` - Active connection to peer
- `error` - Connection error

**Errors:**
- `ErrAuth` - Authentication failure
- `ErrPeerNotFound` - Peer not available
- `ErrTimeout` - Connection timeout
- `ErrNetwork` - Network error

**Example:**
```go
conn, err := client.Connect(ctx, "peer-123")
if err != nil {
    return err
}
defer conn.Close()
```

### Client.CreateTunnel

Creates a secure tunnel with the specified configuration.

```go
func (c *Client) CreateTunnel(ctx context.Context, config TunnelConfig) (Tunnel, error)
```

**Parameters:**
- `ctx` - Context for cancellation
- `config` - Tunnel configuration

**Returns:**
- `Tunnel` - Active tunnel instance
- `error` - Tunnel creation error

**Example:**
```go
tunnel, err := client.CreateTunnel(ctx, cloudbridge.TunnelConfig{
    LocalPort:  8080,
    RemotePeer: "peer-123",
    RemotePort: 3000,
})
```

### Client.JoinMesh

Joins a mesh network with the specified name.

```go
func (c *Client) JoinMesh(ctx context.Context, networkName string) (Mesh, error)
```

**Parameters:**
- `ctx` - Context for cancellation
- `networkName` - Name of mesh network

**Returns:**
- `Mesh` - Mesh network instance
- `error` - Join error

**Example:**
```go
mesh, err := client.JoinMesh(ctx, "my-network")
if err != nil {
    return err
}
defer mesh.Leave()
```

### Client.RegisterService

Registers a service for discovery.

```go
func (c *Client) RegisterService(ctx context.Context, config ServiceConfig) error
```

**Parameters:**
- `ctx` - Context for cancellation
- `config` - Service configuration

**Returns:**
- `error` - Registration error

**Example:**
```go
err := client.RegisterService(ctx, cloudbridge.ServiceConfig{
    Name: "my-api",
    Port: 8080,
    Tags: []string{"http", "api"},
})
```

### Client.DiscoverServices

Discovers services by name.

```go
func (c *Client) DiscoverServices(ctx context.Context, serviceName string) ([]Service, error)
```

**Parameters:**
- `ctx` - Context for cancellation
- `serviceName` - Service name to discover

**Returns:**
- `[]Service` - List of discovered services
- `error` - Discovery error

**Example:**
```go
services, err := client.DiscoverServices(ctx, "my-api")
```

### Client.Health

Checks the health of the client connection.

```go
func (c *Client) Health(ctx context.Context) (*Health, error)
```

**Parameters:**
- `ctx` - Context for cancellation

**Returns:**
- `*Health` - Health status
- `error` - Health check error

**Example:**
```go
health, err := client.Health(ctx)
if err != nil {
    return err
}
log.Printf("Status: %s, Latency: %v", health.Status, health.Latency)
```

### Client.OnConnect

Registers a callback for connection events.

```go
func (c *Client) OnConnect(callback func(peer string))
```

**Parameters:**
- `callback` - Function called when connected

**Example:**
```go
client.OnConnect(func(peer string) {
    log.Printf("Connected to: %s", peer)
})
```

### Client.OnDisconnect

Registers a callback for disconnection events.

```go
func (c *Client) OnDisconnect(callback func(peer string, err error))
```

**Parameters:**
- `callback` - Function called when disconnected

**Example:**
```go
client.OnDisconnect(func(peer string, err error) {
    log.Printf("Disconnected from %s: %v", peer, err)
})
```

### Client.OnReconnect

Registers a callback for reconnection events.

```go
func (c *Client) OnReconnect(callback func(peer string))
```

**Parameters:**
- `callback` - Function called when reconnected

**Example:**
```go
client.OnReconnect(func(peer string) {
    log.Printf("Reconnected to: %s", peer)
})
```

### Client.Close

Closes the client and releases all resources.

```go
func (c *Client) Close() error
```

**Returns:**
- `error` - Close error

**Example:**
```go
defer client.Close()
```

## Connection

### Connection.Read

Reads data from the connection.

```go
func (c *Connection) Read(b []byte) (int, error)
```

**Parameters:**
- `b` - Buffer to read into

**Returns:**
- `int` - Number of bytes read
- `error` - Read error

### Connection.Write

Writes data to the connection.

```go
func (c *Connection) Write(b []byte) (int, error)
```

**Parameters:**
- `b` - Data to write

**Returns:**
- `int` - Number of bytes written
- `error` - Write error

### Connection.Close

Closes the connection.

```go
func (c *Connection) Close() error
```

**Returns:**
- `error` - Close error

### Connection.PeerID

Returns the peer identifier.

```go
func (c *Connection) PeerID() string
```

**Returns:**
- `string` - Peer identifier

### Connection.Metrics

Returns connection metrics.

```go
func (c *Connection) Metrics() (*ConnectionMetrics, error)
```

**Returns:**
- `*ConnectionMetrics` - Connection metrics
- `error` - Metrics error

**Example:**
```go
metrics, err := conn.Metrics()
if err != nil {
    return err
}
log.Printf("RTT: %v, Bytes sent: %d", metrics.RTT, metrics.BytesSent)
```

### Connection.SetDeadline

Sets the read and write deadlines.

```go
func (c *Connection) SetDeadline(t time.Time) error
```

**Parameters:**
- `t` - Deadline time

**Returns:**
- `error` - Set deadline error

### Connection.SetReadDeadline

Sets the read deadline.

```go
func (c *Connection) SetReadDeadline(t time.Time) error
```

### Connection.SetWriteDeadline

Sets the write deadline.

```go
func (c *Connection) SetWriteDeadline(t time.Time) error
```

## Tunnel

### Tunnel.RemotePeer

Returns the remote peer ID.

```go
func (t *Tunnel) RemotePeer() string
```

**Returns:**
- `string` - Remote peer identifier

### Tunnel.LocalPort

Returns the local port.

```go
func (t *Tunnel) LocalPort() int
```

**Returns:**
- `int` - Local port number

### Tunnel.RemotePort

Returns the remote port.

```go
func (t *Tunnel) RemotePort() int
```

**Returns:**
- `int` - Remote port number

### Tunnel.Close

Closes the tunnel.

```go
func (t *Tunnel) Close() error
```

**Returns:**
- `error` - Close error

## Mesh

### Mesh.NetworkName

Returns the network name.

```go
func (m *Mesh) NetworkName() string
```

**Returns:**
- `string` - Network name

### Mesh.Broadcast

Broadcasts a message to all peers.

```go
func (m *Mesh) Broadcast(ctx context.Context, data []byte) error
```

**Parameters:**
- `ctx` - Context for cancellation
- `data` - Message data

**Returns:**
- `error` - Broadcast error

**Example:**
```go
err := mesh.Broadcast(ctx, []byte("Hello mesh!"))
```

### Mesh.Send

Sends a message to a specific peer.

```go
func (m *Mesh) Send(ctx context.Context, peerID string, data []byte) error
```

**Parameters:**
- `ctx` - Context for cancellation
- `peerID` - Target peer identifier
- `data` - Message data

**Returns:**
- `error` - Send error

**Example:**
```go
err := mesh.Send(ctx, "peer-123", []byte("Direct message"))
```

### Mesh.Messages

Returns a channel for receiving messages.

```go
func (m *Mesh) Messages() <-chan Message
```

**Returns:**
- `<-chan Message` - Message receive channel

**Example:**
```go
for msg := range mesh.Messages() {
    log.Printf("From %s: %s", msg.From, string(msg.Data))
}
```

### Mesh.Peers

Returns a list of connected peers.

```go
func (m *Mesh) Peers() []string
```

**Returns:**
- `[]string` - List of peer identifiers

**Example:**
```go
peers := mesh.Peers()
log.Printf("Connected to %d peers", len(peers))
```

### Mesh.Leave

Leaves the mesh network.

```go
func (m *Mesh) Leave() error
```

**Returns:**
- `error` - Leave error

## Configuration

### WithToken

Sets the authentication token.

```go
func WithToken(token string) Option
```

**Parameters:**
- `token` - API authentication token

### WithRegion

Sets the preferred region.

```go
func WithRegion(region string) Option
```

**Parameters:**
- `region` - Region identifier (e.g., "eu-central")

### WithTimeout

Sets the operation timeout.

```go
func WithTimeout(timeout time.Duration) Option
```

**Parameters:**
- `timeout` - Timeout duration

### WithLogLevel

Sets the log level.

```go
func WithLogLevel(level string) Option
```

**Parameters:**
- `level` - Log level ("debug", "info", "warn", "error")

### WithRetryPolicy

Sets the retry policy.

```go
func WithRetryPolicy(policy RetryPolicy) Option
```

**Parameters:**
- `policy` - Retry policy configuration

### WithProtocols

Sets the protocol preferences.

```go
func WithProtocols(protocols ...Protocol) Option
```

**Parameters:**
- `protocols` - Ordered list of protocols

### WithInsecureSkipVerify

Disables TLS certificate verification.

```go
func WithInsecureSkipVerify(skip bool) Option
```

**Parameters:**
- `skip` - Whether to skip verification

**Warning:** Not recommended for production use.

## Errors

### IsAuthError

Checks if an error is an authentication error.

```go
func IsAuthError(err error) bool
```

### IsNetworkError

Checks if an error is a network error.

```go
func IsNetworkError(err error) bool
```

### IsPeerNotFoundError

Checks if an error is a peer not found error.

```go
func IsPeerNotFoundError(err error) bool
```

### IsTimeoutError

Checks if an error is a timeout error.

```go
func IsTimeoutError(err error) bool
```

## Types

### Health

```go
type Health struct {
    Status         string
    Latency        time.Duration
    ConnectedPeers int
}
```

### ConnectionMetrics

```go
type ConnectionMetrics struct {
    BytesSent     uint64
    BytesReceived uint64
    RTT           time.Duration
    Connected     bool
    ConnectedAt   time.Time
}
```

### TunnelConfig

```go
type TunnelConfig struct {
    LocalPort  int
    RemotePeer string
    RemotePort int
    Protocol   Protocol
}
```

### ServiceConfig

```go
type ServiceConfig struct {
    Name string
    Port int
    Tags []string
}
```

### Service

```go
type Service struct {
    Name    string
    Address string
    Port    int
    Tags    []string
    Healthy bool
}
```

### Message

```go
type Message struct {
    From string
    Data []byte
}
```

### Protocol

```go
type Protocol string

const (
    ProtocolQUIC      Protocol = "quic"
    ProtocolGRPC      Protocol = "grpc"
    ProtocolWebSocket Protocol = "websocket"
    ProtocolTCP       Protocol = "tcp"
)
```

### RetryPolicy

```go
type RetryPolicy struct {
    MaxRetries   int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}
```

## Environment Variables

- `CLOUDBRIDGE_TOKEN` - Authentication token
- `CLOUDBRIDGE_REGION` - Preferred region
- `CLOUDBRIDGE_LOG_LEVEL` - Log level
- `CLOUDBRIDGE_TIMEOUT` - Operation timeout

## Notes

- All context-aware methods respect cancellation
- Timeouts are configurable per operation
- Automatic retry on transient failures
- Connection pooling for performance
- Metrics collection for observability
