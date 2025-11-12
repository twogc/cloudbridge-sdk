# CloudBridge SDK

Official Software Development Kit for CloudBridge Global Network - Enterprise-grade P2P mesh networking platform.

## Overview

CloudBridge SDK provides developers with easy-to-use libraries for integrating CloudBridge's P2P mesh networking capabilities into their applications. The SDK supports multiple programming languages and abstracts the complexity of QUIC, MASQUE, and authentication protocols.

## Supported Languages

- **Go** - Primary SDK for backend services
- **Python** - For AI/ML and data science applications
- **JavaScript/TypeScript** - For web and Node.js applications

## Features

- P2P mesh networking with automatic peer discovery
- Secure tunneling (TCP port forwarding)
- L3-overlay networks with WireGuard integration
- Multi-protocol support (QUIC, gRPC, WebSocket)
- OIDC/JWT authentication with Zitadel integration
- Multi-tenant isolation
- Automatic failover and reconnection
- Prometheus metrics and health checks
- NAT traversal (ICE/STUN/TURN)

## Quick Start

### Go

```go
package main

import "github.com/twogc/cloudbridge-sdk/go/cloudbridge"

func main() {
    client, err := cloudbridge.NewClient(
        cloudbridge.WithToken("your-token"),
        cloudbridge.WithRegion("eu-central"),
    )
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Create P2P connection
    conn, err := client.Connect("peer-id-123")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```

### Python

```python
from cloudbridge import Client

client = Client(token="your-token", region="eu-central")

# Create P2P connection
conn = client.connect("peer-id-123")
conn.close()
```

### JavaScript

```javascript
import { CloudBridge } from '@twogc/cloudbridge-sdk';

const client = new CloudBridge({
  token: 'your-token',
  region: 'eu-central'
});

// Create P2P connection
const conn = await client.connect('peer-id-123');
await conn.close();
```

## Installation

### Go

```bash
go get github.com/twogc/cloudbridge-sdk/go/cloudbridge
```

### Python

```bash
pip install cloudbridge-sdk
```

### JavaScript

```bash
npm install @twogc/cloudbridge-sdk
```

## Documentation

- [Go SDK Documentation](./go/README.md)
- [Python SDK Documentation](./python/README.md)
- [JavaScript SDK Documentation](./javascript/README.md)
- [API Reference](./docs/API_REFERENCE.md)
- [Architecture Overview](./docs/ARCHITECTURE.md)
- [Authentication Guide](./docs/AUTHENTICATION.md)

## Architecture

CloudBridge SDK integrates with the CloudBridge Global Network platform:

```
Client Application (Your Code)
        |
        v
CloudBridge SDK (This Library)
        |
        v
CloudBridge Relay Client
        |
        v
CloudBridge Global Network
    - DNS Network (GeoDNS + Anycast)
    - Control Plane (Zitadel OIDC)
    - DDoS Protection (ML-based)
    - Scalable Relay (QUIC/MASQUE/WireGuard)
    - Monitoring (Prometheus + Grafana)
    - AI Service (Route optimization)
```

## Use Cases

### Secure Service-to-Service Communication

```go
// Connect microservices across regions
client.ConnectServices("eu-backend", "us-frontend")
```

### Remote Access and Tunneling

```python
# Create secure tunnel to remote service
tunnel = client.create_tunnel(
    local_port=8080,
    remote_peer="server-123",
    remote_port=3000
)
```

### P2P Mesh Networks

```javascript
// Join mesh network
const mesh = await client.joinMesh('my-network');
await mesh.broadcast({ type: 'hello', data: 'world' });
```

### AI/ML Model Communication

```python
# Secure data sharing between ML models
client.share_dataset(
    source="model-1",
    destination="model-2",
    encrypted=True
)
```

## Requirements

- Go 1.25+ (for Go SDK)
- Python 3.9+ (for Python SDK)
- Node.js 18+ (for JavaScript SDK)
- CloudBridge account and API token

## Getting API Token

1. Visit [CloudBridge Dashboard](https://dashboard.cloudbridge.global)
2. Navigate to Settings > API Tokens
3. Click "Generate New Token"
4. Copy and use in your SDK configuration

## Support

- Documentation: [https://docs.cloudbridge.global](https://docs.cloudbridge.global)
- GitHub Issues: [https://github.com/twogc/cloudbridge-sdk/issues](https://github.com/twogc/cloudbridge-sdk/issues)
- Community: [https://community.cloudbridge.global](https://community.cloudbridge.global)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](./LICENSE) for details.

## Security

For security issues, please email security@2gc.ru instead of using the issue tracker.

## Project Status

Current version: 0.1.0 (Alpha)

- Go SDK: In Development
- Python SDK: Planned
- JavaScript SDK: Planned

## Roadmap

### Phase 1 (Q1 2026)
- Go SDK Core functionality
- Authentication (JWT + OIDC)
- P2P connection management
- Basic tunneling

### Phase 2 (Q2 2026)
- Python SDK release
- JavaScript SDK release
- Advanced features (mesh networking)
- WebSocket fallback for browsers

### Phase 3 (Q3 2026)
- Mobile SDKs (iOS, Android)
- Advanced routing options
- Performance optimizations
- Enterprise features

## Related Projects

- [cloudbridge-docs](https://github.com/twogc/cloudbridge-docs) - Architecture documentation

## Copyright

Copyright (c) 2025 2GC CloudBridge Global Network
