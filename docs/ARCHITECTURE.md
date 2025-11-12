# CloudBridge SDK Architecture

Version: 0.1.0
Last Updated: November 2025

## Overview

CloudBridge SDK provides a high-level interface for integrating CloudBridge Global Network capabilities into applications. The SDK abstracts the complexity of P2P networking, authentication, and protocol negotiation.

## Architecture Layers

```
Application Layer
    |
    v
CloudBridge SDK (This Library)
    |
    +-- Client Layer
    |   +-- Configuration Management
    |   +-- Connection Management
    |   +-- Lifecycle Management
    |
    +-- Protocol Layer
    |   +-- QUIC Transport
    |   +-- gRPC Transport
    |   +-- WebSocket Transport
    |   +-- Protocol Selection
    |
    +-- Security Layer
    |   +-- JWT Authentication
    |   +-- OIDC Integration
    |   +-- TLS 1.3 Encryption
    |   +-- Certificate Management
    |
    +-- Service Layer
    |   +-- Connection Service
    |   +-- Tunnel Service
    |   +-- Mesh Service
    |   +-- Discovery Service
    |
    +-- Utility Layer
        +-- Error Handling
        +-- Retry Logic
        +-- Metrics Collection
        +-- Logging
    |
    v
CloudBridge Relay Client (Native Implementation)
    |
    v
CloudBridge Global Network
```

## Core Components

### 1. Client

The main entry point for SDK operations.

**Responsibilities:**
- Configuration management
- Connection lifecycle
- Service coordination
- Event callbacks

**Key Methods:**
- `NewClient(options...)` - Creates configured client
- `Connect(ctx, peerID)` - Establishes P2P connection
- `CreateTunnel(ctx, config)` - Creates secure tunnel
- `JoinMesh(ctx, network)` - Joins mesh network
- `Health(ctx)` - Health check
- `Close()` - Cleanup resources

### 2. Connection

Represents an active P2P connection to a peer.

**Responsibilities:**
- Data transmission (Read/Write)
- Connection metrics
- Deadline management
- Error handling

**Interfaces:**
- `io.ReadWriteCloser` - Standard Go I/O
- `Connection` - Extended functionality

### 3. Tunnel

Provides TCP port forwarding capabilities.

**Responsibilities:**
- Local port listening
- Remote peer connection
- Bidirectional forwarding
- Connection multiplexing

**Configuration:**
- Local port
- Remote peer ID
- Remote port
- Protocol selection

### 4. Mesh

Enables mesh networking with multiple peers.

**Responsibilities:**
- Peer discovery
- Message broadcasting
- Direct messaging
- Peer management

**Features:**
- Automatic peer discovery
- Message routing
- Connection recovery
- Event notification

### 5. Error Handling

Structured error types for different failure scenarios.

**Error Categories:**
- Authentication errors
- Network errors
- Peer not found
- Timeout errors
- Invalid input

**Error Checking:**
```go
if errors.IsAuthError(err) {
    // Handle authentication failure
}
```

## Data Flow

### Connection Establishment

```
Application
    |
    | client.Connect(ctx, peerID)
    v
Client
    |
    | 1. Validate configuration
    | 2. Select protocol (QUIC/gRPC/WebSocket)
    v
Protocol Layer
    |
    | 3. Authenticate via JWT
    v
Security Layer
    |
    | 4. Establish encrypted connection
    v
CloudBridge Network
    |
    | 5. Peer lookup via DNS
    | 6. Route through optimal PoP
    v
Remote Peer
```

### Message Transmission

```
Application Data
    |
    | conn.Write(data)
    v
Connection
    |
    | 1. Buffer management
    | 2. Metrics update
    v
Protocol Layer
    |
    | 3. Framing
    | 4. Encryption
    v
Transport Layer
    |
    | 5. QUIC/gRPC/WebSocket
    v
Network
```

## Protocol Selection

The SDK automatically selects the best protocol based on:

1. Network conditions
2. Protocol availability
3. Configuration preferences
4. Fallback requirements

**Protocol Priority:**
1. QUIC (preferred) - Best performance, 0-RTT
2. gRPC - Good for service mesh
3. WebSocket - Browser compatibility

**Selection Algorithm:**
```
1. Try QUIC connection
   - If successful: use QUIC
   - If failed: try next protocol

2. Try gRPC connection
   - If successful: use gRPC
   - If failed: try next protocol

3. Try WebSocket connection
   - If successful: use WebSocket
   - If failed: return error

4. Connection established with selected protocol
```

## Authentication Flow

```
Client Creation
    |
    | WithToken("token")
    v
Token Validation
    |
    | 1. Check token format
    | 2. Verify expiration
    v
OIDC Integration
    |
    | 3. Validate with Zitadel
    | 4. Get JWT claims
    v
Connection Request
    |
    | 5. Attach JWT to request
    | 6. Send to CloudBridge Network
    v
Authorization
    |
    | 7. Verify JWT signature
    | 8. Check tenant permissions
    v
Connection Granted
```

## Configuration Management

### Configuration Sources (Priority Order)

1. **Explicit Options** (Highest Priority)
   ```go
   cloudbridge.WithToken("token")
   ```

2. **Environment Variables**
   ```
   CLOUDBRIDGE_TOKEN=token
   CLOUDBRIDGE_REGION=eu-central
   ```

3. **Default Values** (Lowest Priority)
   ```go
   defaultConfig()
   ```

### Configuration Validation

All configurations are validated on client creation:
- Required fields present
- Valid value ranges
- Protocol compatibility
- Credential format

## Error Recovery

### Retry Policy

Configurable retry behavior for transient failures:

```go
RetryPolicy{
    MaxRetries:   3,
    InitialDelay: 1s,
    MaxDelay:     1m,
    Multiplier:   2.0,
}
```

**Retry Algorithm:**
```
attempt = 0
delay = InitialDelay

while attempt < MaxRetries:
    try operation
    if success:
        return result

    wait(min(delay, MaxDelay))
    delay = delay * Multiplier
    attempt++

return error
```

### Automatic Reconnection

The SDK automatically reconnects on connection loss:

1. Detect connection failure
2. Trigger `OnDisconnect` callback
3. Wait for initial delay
4. Attempt reconnection
5. On success: trigger `OnReconnect`
6. On failure: retry with backoff

## Metrics and Observability

### Connection Metrics

Tracked per connection:
- Bytes sent
- Bytes received
- Round-trip time (RTT)
- Connection duration
- Error counts

### Client Metrics

Exposed via Prometheus:
- Active connections
- Connection attempts
- Failed connections
- Protocol distribution
- Average latency

### Logging

Structured logging with configurable levels:
- `debug` - Detailed traces
- `info` - Normal operations
- `warn` - Warnings
- `error` - Errors only

## Security Considerations

### Transport Security

- TLS 1.3 encryption for all connections
- Certificate validation (configurable)
- Perfect forward secrecy
- ALPN protocol negotiation

### Authentication Security

- JWT token with signature verification
- Token expiration enforcement
- Secure token storage recommendations
- OIDC integration with Zitadel

### Best Practices

1. Never hardcode tokens in source code
2. Use environment variables or secure vaults
3. Implement token rotation
4. Enable TLS verification in production
5. Monitor authentication failures
6. Use connection timeouts
7. Validate peer identities

## Performance Optimization

### Connection Pooling

Reuse connections when possible:
- Connection caching by peer ID
- Automatic cleanup of idle connections
- Configurable pool size

### Protocol Optimization

- QUIC 0-RTT for returning connections
- Connection migration on network change
- Multiplexing multiple streams
- Adaptive congestion control (BBRv3)

### Resource Management

- Bounded connection limits
- Memory pooling for buffers
- Graceful degradation on overload
- Background cleanup routines

## Testing Strategy

### Unit Tests

Test individual components:
- Configuration validation
- Error handling
- Protocol selection
- Retry logic

### Integration Tests

Test component interaction:
- End-to-end connections
- Protocol fallback
- Authentication flow
- Error recovery

### Performance Tests

Benchmark critical paths:
- Connection establishment
- Data throughput
- Latency measurements
- Resource usage

## Future Enhancements

### Planned Features

1. Connection multiplexing
2. Advanced routing options
3. Traffic shaping
4. Enhanced metrics
5. Distributed tracing
6. Mobile SDK variants

### API Stability

Current version: 0.1.0 (Alpha)
- API may change between releases
- Semantic versioning after 1.0.0
- Deprecation warnings before removal

## Related Documentation

- [API Reference](./API_REFERENCE.md)
- [Authentication Guide](./AUTHENTICATION.md)
- [CloudBridge Architecture](https://github.com/twogc/cloudbridge-docs)

## Support

For architecture questions:
- GitHub Discussions
- Technical documentation
- Architecture reviews
