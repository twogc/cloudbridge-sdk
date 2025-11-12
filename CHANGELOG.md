# Changelog

All notable changes to CloudBridge SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial SDK project structure
- Go SDK core implementation
- Client configuration and lifecycle management
- Connection management interface
- Tunnel interface for secure port forwarding
- Mesh networking interface
- Service discovery interface
- Comprehensive error handling
- Authentication support (JWT + OIDC)
- Multi-protocol support (QUIC, gRPC, WebSocket)
- Retry policy configuration
- Health check functionality
- Connection callbacks (OnConnect, OnDisconnect, OnReconnect)
- Environment variable configuration
- Comprehensive documentation:
  - README with quick start guide
  - API Reference documentation
  - Architecture documentation
  - Authentication guide
  - Contributing guidelines
- Example applications:
  - Basic P2P connection example
  - Mesh networking example
- Unit tests for core functionality
- MIT License

### Status
- Go SDK: Alpha (0.1.0)
- Python SDK: Planned
- JavaScript SDK: Planned

## [0.1.0] - 2025-11-12

### Added
- Initial project repository created
- Project structure and scaffolding
- Core Go SDK interfaces and types
- Basic documentation framework

[Unreleased]: https://github.com/twogc/cloudbridge-sdk/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/twogc/cloudbridge-sdk/releases/tag/v0.1.0
