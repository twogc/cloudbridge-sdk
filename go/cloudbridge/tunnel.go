package cloudbridge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", t.config.LocalPort))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	go func() {
		defer listener.Close()
		for {
			t.mu.RLock()
			if t.closed {
				t.mu.RUnlock()
				return
			}
			t.mu.RUnlock()

			conn, err := listener.Accept()
			if err != nil {
				if !t.closed {
					fmt.Printf("failed to accept connection: %v\n", err)
				}
				return
			}

			go t.handleConnection(ctx, conn)
		}
	}()

	return nil
}

func (t *tunnel) handleConnection(ctx context.Context, localConn net.Conn) {
	defer localConn.Close()

	remoteConn, err := t.client.Connect(ctx, t.config.RemotePeer)
	if err != nil {
		fmt.Printf("failed to connect to remote peer: %v\n", err)
		return
	}
	defer remoteConn.Close()

	// Send handshake
	handshake := fmt.Sprintf(`{"type":"tunnel","port":%d}`, t.config.RemotePort)
	if _, err := remoteConn.Write([]byte(handshake)); err != nil {
		fmt.Printf("failed to send handshake: %v\n", err)
		return
	}

	// Bidirectional copy
	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(localConn, remoteConn)
		errChan <- err
	}()

	<-errChan
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
