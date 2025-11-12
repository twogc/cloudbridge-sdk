# CloudBridge CLI

Command-line tool for testing and interacting with the CloudBridge SDK.

## Quick Start

### Build

```bash
# From the go/ directory
make cli

# Or manually
cd cmd/cloudbridge
go build -o cloudbridge .
```

### Run

```bash
# Set your token
export CLOUDBRIDGE_TOKEN="your-token-here"

# Check health
./cloudbridge health

# Discover peers
./cloudbridge discover

# Connect to a peer
./cloudbridge connect peer-id

# Get help
./cloudbridge --help
```

## Installation

### From Source

```bash
cd go/
make install
```

This will build and install the CLI to `/usr/local/bin/cloudbridge`.

### Using Go

```bash
go install github.com/twogc/cloudbridge-sdk/go/cmd/cloudbridge@latest
```

## Commands

- `connect <peer-id>` - Connect to a peer
- `discover` - Discover available peers
- `tunnel <peer-id>` - Create a tunnel to a peer
- `health` - Check system health
- `version` - Print version information

For detailed documentation, see [CLI.md](../../../docs/CLI.md).

## Examples

### Basic Health Check

```bash
cloudbridge health
```

### Interactive Peer Connection

```bash
cloudbridge connect --interactive peer-abc123
```

### Create TCP Tunnel

```bash
cloudbridge tunnel --local localhost:8080 --remote localhost:80 peer-abc123
```

### Watch Peer Discovery

```bash
cloudbridge discover --watch
```

## Configuration

### Environment Variables

- `CLOUDBRIDGE_TOKEN` - Authentication token (required)
- `CLOUDBRIDGE_REGION` - CloudBridge region (default: eu-central)
- `CLOUDBRIDGE_TIMEOUT` - Operation timeout (default: 30s)
- `CLOUDBRIDGE_LOG_LEVEL` - Log level (debug, info, warn, error)

### Flags

Global flags available for all commands:

- `--token`, `-t` - Authentication token
- `--region`, `-r` - CloudBridge region
- `--timeout` - Operation timeout
- `--verbose`, `-v` - Verbose output
- `--insecure-skip-verify` - Skip TLS verification (testing only)

## Development

### Running Tests

```bash
cd ../../
make test
```

### Building for Release

```bash
cd ../../
make release
```

This creates binaries for multiple platforms in `dist/`.

## Documentation

- [Full CLI Documentation](../../../docs/CLI.md)
- [SDK API Reference](../../../docs/API_REFERENCE.md)
- [Integration Guide](../../../docs/INTEGRATION.md)

## Support

- Report issues: https://github.com/twogc/cloudbridge-sdk/issues
- SDK Documentation: https://github.com/twogc/cloudbridge-sdk
