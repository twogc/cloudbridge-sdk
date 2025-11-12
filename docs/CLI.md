# CloudBridge CLI Tool

Version: 0.1.0
Last Updated: November 2025

## Overview

The CloudBridge CLI is a command-line tool for testing and interacting with the CloudBridge SDK. It provides an easy way to connect to peers, discover available peers, create tunnels, and monitor system health.

## Installation

### From Source

```bash
cd cloudbridge-sdk/go
go build -o cloudbridge ./cmd/cloudbridge
sudo mv cloudbridge /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/twogc/cloudbridge-sdk/go/cmd/cloudbridge@latest
```

## Quick Start

### Authentication

The CLI requires a CloudBridge authentication token. You can provide it in two ways:

1. **Using flag:**
   ```bash
   cloudbridge --token "your-token-here" health
   ```

2. **Using environment variable:**
   ```bash
   export CLOUDBRIDGE_TOKEN="your-token-here"
   cloudbridge health
   ```

### Basic Health Check

```bash
cloudbridge health
```

Output:
```
✓ Overall Status: healthy
  Checked at: 2025-11-12T10:30:00Z

Components:
  NAME    STATUS      LATENCY  MESSAGE
  client  ✓ healthy   2ms      CloudBridge client initialized
  auth    ✓ healthy   -        Authentication token valid
  p2p     ✓ healthy   -        P2P manager ready
  quic    ✓ healthy   1ms      QUIC transport available

Connectivity:
  Relay Server: ✓ OK
  Internet:     ✓ OK
  DNS:          ✓ OK
  Latency:      15ms
```

## Commands

### Global Flags

All commands support the following global flags:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--token` | `-t` | - | Authentication token (or set `CLOUDBRIDGE_TOKEN`) |
| `--region` | `-r` | `eu-central` | CloudBridge region |
| `--timeout` | - | `30s` | Operation timeout |
| `--insecure-skip-verify` | - | `false` | Skip TLS certificate verification |
| `--verbose` | `-v` | `false` | Verbose output |

### connect

Connect to a peer and optionally send/receive messages.

**Usage:**
```bash
cloudbridge connect [flags] <peer-id>
```

**Flags:**
- `--interactive`, `-i`: Interactive mode for sending messages
- `--message`, `-m`: Single message to send

**Examples:**

1. **Simple connection:**
   ```bash
   cloudbridge connect peer-123
   ```

2. **Send a single message:**
   ```bash
   cloudbridge connect --message "Hello, peer!" peer-123
   ```

3. **Interactive mode:**
   ```bash
   cloudbridge connect --interactive peer-123
   ```

   In interactive mode, you can:
   - Type messages to send
   - Use commands:
     - `/quit` or `/exit` - Exit interactive mode
     - `/metrics` - Display connection metrics

**Output:**
```
Connecting to peer: peer-123
✓ Connected to peer peer-123
  Connected at: 2025-11-12T10:30:00Z
  Latency: 12ms

Interactive mode - type messages to send (Ctrl+C to exit)
Commands: /quit, /exit, /metrics

> Hello!
[peer-123] Hi there!
> /metrics
Metrics: {Connected:true ConnectedAt:2025-11-12 10:30:00 BytesSent:42 BytesReceived:58}
> /quit
Exiting...
```

### discover

Discover and list available peers in the CloudBridge network.

**Usage:**
```bash
cloudbridge discover [flags]
```

**Flags:**
- `--json`: Output as JSON
- `--watch`, `-w`: Watch for peer changes
- `--filter`, `-f`: Filter peers by ID pattern

**Examples:**

1. **List all peers:**
   ```bash
   cloudbridge discover
   ```

2. **JSON output:**
   ```bash
   cloudbridge discover --json
   ```

3. **Watch for changes:**
   ```bash
   cloudbridge discover --watch
   ```

**Output (Table):**
```
Discovering peers in CloudBridge network...

PEER ID           STATUS   REGION       LATENCY   PROTOCOL   LAST SEEN
example-peer-1    online   eu-central   12ms      QUIC       5m0s ago
example-peer-2    online   us-east      45ms      QUIC       2m0s ago

Total peers: 2
```

**Output (JSON):**
```json
[
  {
    "peer_id": "example-peer-1",
    "status": "online",
    "region": "eu-central",
    "last_seen": "2025-11-12T10:25:00Z",
    "latency": "12ms",
    "protocol": "QUIC",
    "public_addr": "203.0.113.1:4433",
    "capabilities": ["tunnel", "mesh"]
  },
  {
    "peer_id": "example-peer-2",
    "status": "online",
    "region": "us-east",
    "last_seen": "2025-11-12T10:28:00Z",
    "latency": "45ms",
    "protocol": "QUIC",
    "public_addr": "203.0.113.2:4433",
    "capabilities": ["tunnel"]
  }
]
```

### tunnel

Create a TCP or UDP tunnel to a peer.

**Usage:**
```bash
cloudbridge tunnel [flags] <peer-id>
```

**Flags:**
- `--local`, `-l`: Local address to listen on (required)
- `--remote`, `-r`: Remote address to forward to (required)
- `--protocol`, `-p`: Protocol - `tcp` or `udp` (default: `tcp`)

**Examples:**

1. **TCP tunnel:**
   ```bash
   cloudbridge tunnel --local localhost:8080 --remote localhost:80 peer-123
   ```

2. **UDP tunnel:**
   ```bash
   cloudbridge tunnel --local localhost:5353 --remote localhost:53 --protocol udp peer-123
   ```

**Output:**
```
Creating tcp tunnel to peer: peer-123
  Local:  localhost:8080
  Remote: localhost:80
✓ Tunnel established
  Listening on: localhost:8080

Press Ctrl+C to stop the tunnel

[10:30:00] New connection from 127.0.0.1:54321 (total: 1)
[10:30:00] Forwarding 127.0.0.1:54321 -> localhost:80 (via peer-123)
[10:30:05] New connection from 127.0.0.1:54322 (total: 2)
[10:30:05] Forwarding 127.0.0.1:54322 -> localhost:80 (via peer-123)

^C
Shutting down tunnel...
Total connections handled: 2
```

### health

Check the health of CloudBridge client and connectivity.

**Usage:**
```bash
cloudbridge health [flags]
```

**Flags:**
- `--json`: Output as JSON
- `--watch`, `-w`: Watch health status continuously
- `--verbose`: Show detailed health information

**Examples:**

1. **Basic health check:**
   ```bash
   cloudbridge health
   ```

2. **Verbose health check:**
   ```bash
   cloudbridge health --verbose
   ```

3. **Watch mode:**
   ```bash
   cloudbridge health --watch
   ```

4. **JSON output:**
   ```bash
   cloudbridge health --json
   ```

**Output (Table):**
```
✓ Overall Status: healthy
  Checked at: 2025-11-12T10:30:00Z

Components:
  NAME    STATUS      LATENCY  MESSAGE
  client  ✓ healthy   2ms      CloudBridge client initialized
  auth    ✓ healthy   -        Authentication token valid
  p2p     ✓ healthy   -        P2P manager ready
  quic    ✓ healthy   1ms      QUIC transport available

Connectivity:
  Relay Server: ✓ OK
  Internet:     ✓ OK
  DNS:          ✓ OK
  Latency:      15ms

Performance:
  Memory Usage: 12.5 MB
  Goroutines:   8
  Connections:  0
```

**Output (JSON):**
```json
{
  "timestamp": "2025-11-12T10:30:00Z",
  "overall": "healthy",
  "components": {
    "client": {
      "status": "healthy",
      "message": "CloudBridge client initialized",
      "latency": "2ms"
    },
    "auth": {
      "status": "healthy",
      "message": "Authentication token valid"
    },
    "p2p": {
      "status": "healthy",
      "message": "P2P manager ready"
    },
    "quic": {
      "status": "healthy",
      "message": "QUIC transport available",
      "latency": "1ms"
    }
  },
  "connectivity": {
    "relay_server": true,
    "internet": true,
    "dns": true,
    "latency": "15ms"
  },
  "performance": {
    "memory_usage": "12.5 MB",
    "goroutines": 8,
    "connections": 0
  }
}
```

### version

Print version information.

**Usage:**
```bash
cloudbridge version
```

**Output:**
```
CloudBridge CLI version 0.1.0
SDK version: 0.1.0
Build date: November 2025
```

## Use Cases

### 1. Testing Peer Connection

```bash
# Check health
cloudbridge health

# Discover available peers
cloudbridge discover

# Connect to a specific peer
cloudbridge connect --interactive peer-abc123
```

### 2. Setting Up a Tunnel

```bash
# Create a tunnel to forward local port 8080 to remote port 80
cloudbridge tunnel --local localhost:8080 --remote localhost:80 peer-abc123

# In another terminal, test the tunnel
curl http://localhost:8080
```

### 3. Monitoring System Health

```bash
# Continuous health monitoring
cloudbridge health --watch --verbose
```

### 4. Automated Testing

```bash
#!/bin/bash

# Check if CloudBridge is healthy
if cloudbridge health --json | jq -e '.overall == "healthy"'; then
    echo "System is healthy"

    # Connect to peer and send test message
    cloudbridge connect --message "test-message" peer-abc123
else
    echo "System is not healthy"
    exit 1
fi
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `CLOUDBRIDGE_TOKEN` | Authentication token | - |
| `CLOUDBRIDGE_REGION` | CloudBridge region | `eu-central` |
| `CLOUDBRIDGE_TIMEOUT` | Operation timeout | `30s` |
| `CLOUDBRIDGE_LOG_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |

## Configuration File

You can create a configuration file at `~/.cloudbridge/config.yaml`:

```yaml
token: your-token-here
region: eu-central
timeout: 30s
log_level: info
insecure_skip_verify: false
```

## Troubleshooting

### Connection Timeout

**Issue:** Connection to peer times out

**Solutions:**
- Increase timeout: `--timeout 60s`
- Check peer is online: `cloudbridge discover`
- Verify network connectivity: `cloudbridge health`
- Check firewall rules

### Authentication Error

**Issue:** "token is required" or authentication fails

**Solutions:**
- Verify token is set: `echo $CLOUDBRIDGE_TOKEN`
- Check token format (should be valid JWT)
- Verify token hasn't expired
- Use `--verbose` for detailed error messages

### TLS Certificate Error

**Issue:** TLS certificate verification fails

**Solutions:**
- Use `--insecure-skip-verify` for testing (not recommended for production)
- Update system CA certificates
- Contact CloudBridge support for valid certificates

### Peer Not Found

**Issue:** Cannot connect to peer

**Solutions:**
- List available peers: `cloudbridge discover`
- Verify peer ID is correct
- Check peer is in the same region
- Verify peer is online

## Advanced Usage

### Using with Docker

```dockerfile
FROM golang:1.25.3-alpine
WORKDIR /app
COPY . .
RUN go build -o cloudbridge ./cmd/cloudbridge
ENTRYPOINT ["/app/cloudbridge"]
```

Build and run:
```bash
docker build -t cloudbridge-cli .
docker run --rm -e CLOUDBRIDGE_TOKEN="your-token" cloudbridge-cli health
```

### CI/CD Integration

**GitHub Actions:**
```yaml
- name: Test CloudBridge Connection
  run: |
    cloudbridge health
    cloudbridge discover
  env:
    CLOUDBRIDGE_TOKEN: ${{ secrets.CLOUDBRIDGE_TOKEN }}
```

**GitLab CI:**
```yaml
test_cloudbridge:
  script:
    - cloudbridge health --json
    - cloudbridge discover --json
  variables:
    CLOUDBRIDGE_TOKEN: $CLOUDBRIDGE_TOKEN
```

### Scripting Examples

**Bash script to monitor peers:**
```bash
#!/bin/bash

while true; do
    PEER_COUNT=$(cloudbridge discover --json | jq 'length')
    echo "$(date): $PEER_COUNT peers online"

    if [ "$PEER_COUNT" -lt 2 ]; then
        echo "Warning: Low peer count!"
    fi

    sleep 60
done
```

**Python script using CLI:**
```python
import subprocess
import json

def get_peers():
    result = subprocess.run(
        ['cloudbridge', 'discover', '--json'],
        capture_output=True,
        text=True
    )
    return json.loads(result.stdout)

peers = get_peers()
for peer in peers:
    print(f"Peer: {peer['peer_id']} - {peer['status']}")
```

## Performance Tips

1. **Reduce latency**: Use `--region` closest to your location
2. **Batch operations**: Use scripts to automate multiple commands
3. **Enable verbose logging**: Only when debugging (`--verbose`)
4. **Reuse connections**: Keep interactive sessions open for multiple messages
5. **Monitor health**: Use `cloudbridge health --watch` for continuous monitoring

## Security Best Practices

1. **Never commit tokens**: Use environment variables or secure vaults
2. **Rotate tokens regularly**: Update tokens periodically
3. **Use TLS verification**: Avoid `--insecure-skip-verify` in production
4. **Limit tunnel exposure**: Only expose necessary ports
5. **Monitor access**: Use `cloudbridge health` to detect unauthorized access

## Getting Help

- **Built-in help**: `cloudbridge --help`
- **Command help**: `cloudbridge connect --help`
- **Report issues**: https://github.com/twogc/cloudbridge-sdk/issues
- **Documentation**: https://github.com/twogc/cloudbridge-sdk/tree/main/docs

## Related Documentation

- [SDK API Reference](./API_REFERENCE.md)
- [Integration Guide](./INTEGRATION.md)
- [Architecture Overview](./ARCHITECTURE.md)
- [Authentication Guide](./AUTHENTICATION.md)
