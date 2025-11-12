package cloudbridge

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Tunnel represents a secure tunnel
type Tunnel interface {
	// RemotePeer returns the remote peer ID
	RemotePeer() string

	// LocalPort returns the local port
	LocalPort() int

	// RemotePort returns the remote port
	RemotePort() int

	// Close closes the tunnel
	Close() error
}

// TunnelConfig holds configuration for a tunnel
type TunnelConfig struct {
	LocalPort  int
	RemotePeer string
	RemotePort int
	Protocol   Protocol
}

// validate checks if the tunnel configuration is valid
func (tc *TunnelConfig) validate() error {
	if tc.LocalPort <= 0 || tc.LocalPort > 65535 {
		return errors.New("invalid local port")
	}

	if tc.RemotePort <= 0 || tc.RemotePort > 65535 {
		return errors.New("invalid remote port")
	}

	if tc.RemotePeer == "" {
		return errors.New("remote peer cannot be empty")
	}

	if tc.Protocol == "" {
		tc.Protocol = ProtocolTCP
	}

	validProtocols := map[Protocol]bool{
		ProtocolTCP:  true,
		ProtocolQUIC: true,
	}

	if !validProtocols[tc.Protocol] {
		return fmt.Errorf("invalid protocol: %s", tc.Protocol)
	}

	return nil
}

// tunnel implements the Tunnel interface
type tunnel struct {
	config TunnelConfig
	client *Client
	mu     sync.RWMutex
	closed bool
}

// start starts the tunnel
func (t *tunnel) start(ctx context.Context) error {
	// TODO: Implement tunnel creation
	// This would involve:
	// 1. Establishing connection to remote peer
	// 2. Setting up local listener on LocalPort
	// 3. Forwarding traffic between local and remote
	// 4. Handling reconnection on failure

	return nil
}

// RemotePeer returns the remote peer ID
func (t *tunnel) RemotePeer() string {
	return t.config.RemotePeer
}

// LocalPort returns the local port
func (t *tunnel) LocalPort() int {
	return t.config.LocalPort
}

// RemotePort returns the remote port
func (t *tunnel) RemotePort() int {
	return t.config.RemotePort
}

// Close closes the tunnel
func (t *tunnel) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	// TODO: Close tunnel resources

	return nil
}
