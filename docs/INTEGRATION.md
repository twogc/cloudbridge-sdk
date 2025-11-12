# CloudBridge SDK Integration Guide

Version: 0.1.0
Last Updated: November 2025

## Overview

CloudBridge SDK integrates with the CloudBridge Relay Client to provide simplified P2P networking capabilities. This document describes the integration architecture and how the SDK leverages the existing relay client infrastructure.

## Integration Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│                  (Your Application Code)                     │
└─────────────────────────────────────────────────────────────┘
                              |
                              v
┌─────────────────────────────────────────────────────────────┐
│                    CloudBridge SDK                           │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  High-Level API                                        │  │
│  │  - Client                                              │  │
│  │  - Connection                                          │  │
│  │  - Tunnel                                              │  │
│  │  - Mesh                                                │  │
│  └───────────────────────────────────────────────────────┘  │
│                              |                               │
│                              v                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Transport Layer (transport.go)                        │  │
│  │  - Protocol selection                                  │  │
│  │  - Connection management                               │  │
│  │  - Error handling                                      │  │
│  └───────────────────────────────────────────────────────┘  │
│                              |                               │
│                              v                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Bridge Layer (internal/bridge)                        │  │
│  │  - ClientBridge                                        │  │
│  │  - Protocol adaptation                                 │  │
│  │  - State management                                    │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              |
                              v
┌─────────────────────────────────────────────────────────────┐
│            CloudBridge Relay Client (pkg/)                   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  P2P Manager (pkg/p2p)                                 │  │
│  │  - Connection establishment                            │  │
│  │  - Mesh networking                                     │  │
│  │  - Peer discovery                                      │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  QUIC Connection (pkg/quic)                            │  │
│  │  - Stream management                                   │  │
│  │  - TLS configuration                                   │  │
│  │  - Connection monitoring                               │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  API Manager (pkg/api)                                 │  │
│  │  - REST API client                                     │  │
│  │  - Peer registration                                   │  │
│  │  - Status updates                                      │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Auth Manager (pkg/auth)                               │  │
│  │  - JWT token management                                │  │
│  │  - OIDC integration                                    │  │
│  │  - Token refresh                                       │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              |
                              v
┌─────────────────────────────────────────────────────────────┐
│             CloudBridge Global Network                       │
│  - Relay Servers                                             │
│  - Control Plane                                             │
│  - DDoS Protection                                           │
│  - Monitoring                                                │
└─────────────────────────────────────────────────────────────┘
```

## Component Interactions

### 1. SDK Client Initialization

When an SDK client is created:

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("token"),
    cloudbridge.WithRegion("eu-central"),
)
```

**Internal flow:**
1. SDK validates configuration
2. Creates transport layer
3. Transport creates ClientBridge
4. Bridge initializes:
   - Auth Manager (JWT token handling)
   - API Manager (REST API client)
   - P2P Manager (peer networking)
   - QUIC Connection (transport protocol)

### 2. Peer Connection

When connecting to a peer:

```go
conn, err := client.Connect(ctx, "peer-id")
```

**Internal flow:**
1. SDK validates peer ID
2. Transport requests connection from bridge
3. Bridge uses P2P Manager to:
   - Discover peer information
   - Negotiate connection (ICE/STUN/TURN)
   - Establish QUIC connection
   - Open bidirectional stream
4. Bridge wraps connection as PeerConnection
5. SDK wraps as Connection interface
6. Returns to application

### 3. Data Transmission

When reading/writing data:

```go
n, err := conn.Write(data)
n, err := conn.Read(buffer)
```

**Internal flow:**
1. SDK Connection validates state
2. Calls PeerConnection Read/Write
3. PeerConnection uses QUIC stream
4. Data flows through:
   - QUIC encryption layer
   - TLS 1.3 encryption
   - Network transport
5. CloudBridge relay (if needed)
6. Peer's QUIC connection
7. Peer's application

## Bridge Architecture

### ClientBridge

The `ClientBridge` is the integration layer between SDK and relay client.

**Responsibilities:**
- Initialize relay client components
- Manage component lifecycle
- Translate SDK calls to relay client operations
- Handle error translation
- Manage state synchronization

**Key Methods:**

```go
type ClientBridge struct {
    config      *BridgeConfig
    p2pManager  *p2p.Manager
    quicConn    *quic.QUICConnection
    apiManager  *api.Manager
    authManager *auth.AuthManager
}

func (b *ClientBridge) Initialize(ctx context.Context) error
func (b *ClientBridge) ConnectToPeer(ctx context.Context, peerID string) (*PeerConnection, error)
func (b *ClientBridge) DiscoverPeers(ctx context.Context) ([]*api.Peer, error)
func (b *ClientBridge) Close() error
```

### Transport Layer

The `transport` layer manages protocol selection and connection establishment.

**Responsibilities:**
- Create and manage bridge instance
- Protocol fallback logic (QUIC → gRPC → WebSocket)
- Connection pooling
- Retry logic
- Health monitoring

**Key Methods:**

```go
type transport struct {
    config *Config
    bridge *bridge.ClientBridge
}

func newTransport(config *Config) (*transport, error)
func (t *transport) initialize(ctx context.Context) error
func (t *transport) connectToPeer(ctx context.Context, peerID string) (*connection, error)
func (t *transport) close() error
```

## Configuration Mapping

SDK configuration maps to relay client configuration:

| SDK Config | Bridge Config | Relay Client Component |
|-----------|---------------|----------------------|
| Token | Token | auth.Config.Token |
| Region | RelayServerURL | api.ManagerConfig.BaseURL |
| Timeout | Timeout | api.ManagerConfig.Timeout |
| InsecureSkipVerify | InsecureSkipVerify | quic.TLS.InsecureSkipVerify |
| Protocols | EnableP2P, EnableMesh | p2p.P2PConfig |
| RetryPolicy | MaxRetries, BackoffMultiplier | api.ManagerConfig |

## Authentication Flow

```
SDK Token → Bridge → Auth Manager → Relay Client → CloudBridge Network

1. SDK receives token from application
2. Bridge extracts tenant_id from JWT
3. Auth Manager validates token format
4. API Manager uses token for REST calls
5. P2P Manager uses token for peer auth
6. CloudBridge Network validates via Zitadel
```

## Error Handling

SDK errors are translated from relay client errors:

| Relay Client Error | SDK Error |
|-------------------|-----------|
| auth.ErrInvalidToken | errors.ErrAuth |
| api.ErrNetworkTimeout | errors.ErrTimeout |
| p2p.ErrPeerNotFound | errors.ErrPeerNotFound |
| quic.ErrConnectionFailed | errors.ErrNetwork |

## State Management

### Connection State

```
SDK Connection State:
- connected: bool
- connectedAt: time.Time
- bytesSent: uint64
- bytesReceived: uint64

Bridge State:
- PeerConnection
- QUIC Stream
- Connection metrics

Relay Client State:
- P2P peer connections
- QUIC connection
- ICE agent state
```

### Lifecycle Management

```
SDK Client.Close()
    → Transport.close()
        → Bridge.Close()
            → P2P Manager.Stop()
            → API Manager.Stop()
            → QUIC Connection.Close()
```

## Dependencies

### Go Module Dependencies

SDK depends on relay client as a module:

```go
// go.mod
require (
    github.com/2gc-dev/cloudbridge-client v1.4.20
    github.com/quic-go/quic-go v0.40.1
    google.golang.org/grpc v1.60.1
)
```

### Package Dependencies

```
cloudbridge (SDK)
    ├── internal/bridge
    │   └── imports: github.com/2gc-dev/cloudbridge-client/pkg/*
    ├── transport.go
    │   └── imports: internal/bridge
    ├── client.go
    │   └── imports: transport
    └── connection.go
        └── imports: transport, internal/bridge
```

## Testing Strategy

### Unit Tests

SDK unit tests mock the bridge:

```go
type mockBridge struct {
    connectFunc func(ctx context.Context, peerID string) (*PeerConnection, error)
}

func TestClientConnect(t *testing.T) {
    bridge := &mockBridge{
        connectFunc: func(ctx, peerID) (*PeerConnection, error) {
            return &PeerConnection{PeerID: peerID}, nil
        },
    }
    // Test SDK client with mock bridge
}
```

### Integration Tests

Integration tests use real relay client components:

```go
func TestIntegrationConnect(t *testing.T) {
    // Create SDK client
    client := cloudbridge.NewClient(
        cloudbridge.WithToken(testToken),
    )

    // Connect to test peer
    conn, err := client.Connect(ctx, testPeerID)

    // Verify real connection established
    metrics, _ := conn.Metrics()
    assert.True(t, metrics.Connected)
}
```

### End-to-End Tests

E2E tests connect to real CloudBridge network:

```go
func TestE2EConnect(t *testing.T) {
    // Requires: Real token, real relay server
    client := cloudbridge.NewClient(
        cloudbridge.WithToken(os.Getenv("CLOUDBRIDGE_TOKEN")),
        cloudbridge.WithRegion("eu-central"),
    )

    // Real peer discovery
    peers, err := client.DiscoverPeers(ctx)

    // Real peer connection
    conn, err := client.Connect(ctx, peers[0].PeerID)

    // Real data transmission
    n, err := conn.Write([]byte("test"))
}
```

## Performance Considerations

### Connection Pooling

Bridge maintains connection pool to reduce overhead:
- Reuse existing P2P manager instances
- Cache peer connections
- Lazy initialization of QUIC connections

### Protocol Selection

Transport layer selects optimal protocol:
1. Try QUIC (fastest, lowest latency)
2. Fallback to gRPC (if QUIC blocked)
3. Fallback to WebSocket (if both blocked)

### Memory Management

- Connection cleanup on close
- Stream recycling in QUIC
- Periodic garbage collection
- Resource limits enforcement

## Troubleshooting

### Common Issues

**Bridge initialization fails:**
- Check token validity
- Verify relay server URL
- Check network connectivity
- Review logs from auth manager

**Connection timeout:**
- Increase SDK timeout config
- Check peer availability
- Verify firewall rules
- Review ICE/STUN/TURN settings

**Authentication errors:**
- Verify JWT token format
- Check token expiration
- Validate tenant_id claim
- Review Zitadel configuration

### Debug Mode

Enable debug logging:

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
    cloudbridge.WithLogLevel("debug"),
)
```

This enables verbose logging in:
- SDK client
- Transport layer
- Bridge
- Relay client components

## Future Enhancements

### Planned Improvements

1. Connection multiplexing
2. Advanced retry strategies
3. Protocol optimization hints
4. Enhanced metrics collection
5. Distributed tracing integration
6. Custom transport plugins

### API Stability

Current integration is alpha (0.1.0):
- Bridge API may change
- Transport interface may evolve
- Configuration options may expand
- Semantic versioning after 1.0.0

## Related Documentation

- [SDK Architecture](./ARCHITECTURE.md)
- [API Reference](./API_REFERENCE.md)
- [CloudBridge Client Architecture](https://github.com/2gc-dev/cloudbridge-client/blob/main/Architecture.md)
- [QUIC Protocol Documentation](https://github.com/twogc/cloudbridge-docs/blob/main/LAB/QUIC_Laboratory_Research_Report.md)
