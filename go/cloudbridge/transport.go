package cloudbridge

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge/internal/bridge"
	"github.com/twogc/cloudbridge-sdk/go/cloudbridge/internal/jwt"
)

// transport manages the underlying transport layer using bridge
type transport struct {
	config *Config
	bridge *bridge.ClientBridge
	logger *defaultLogger
	mu     sync.RWMutex
	closed bool
}

// newTransport creates a new transport layer
func newTransport(config *Config) (*transport, error) {
	logger := &defaultLogger{}

	bridgeConfig := &bridge.BridgeConfig{
		Token:              config.Token,
		RelayServerURL:     fmt.Sprintf("https://relay.%s.2gc.ru", config.Region),
		TenantID:           extractTenantID(config.Token),
		InsecureSkipVerify: config.InsecureSkipVerify,
		Timeout:            config.Timeout,
		EnableP2P:          true,
		EnableMesh:         true,
	}

	clientBridge, err := bridge.NewClientBridge(bridgeConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create bridge: %w", err)
	}

	return &transport{
		config: config,
		bridge: clientBridge,
		logger: logger,
	}, nil
}

// initialize initializes the transport
func (t *transport) initialize(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	if err := t.bridge.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize bridge: %w", err)
	}

	return nil
}

// connectToPeer connects to a peer
func (t *transport) connectToPeer(ctx context.Context, peerID string) (*connection, error) {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return nil, fmt.Errorf("transport is closed")
	}
	t.mu.RUnlock()

	peerConn, err := t.bridge.ConnectToPeer(ctx, peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer: %w", err)
	}

	conn := &connection{
		peerID:      peerID,
		connected:   true,
		connectedAt: peerConn.ConnectedAt,
		bridgeConn:  peerConn,
	}

	return conn, nil
}

// close closes the transport
func (t *transport) close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true

	if t.bridge != nil {
		return t.bridge.Close()
	}

	return nil
}

// broadcast sends data to all connected peers
func (t *transport) broadcast(ctx context.Context, data []byte) error {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return fmt.Errorf("transport is closed")
	}
	t.mu.RUnlock()

	return t.bridge.Broadcast(ctx, data)
}

// send sends data to a specific peer
func (t *transport) send(ctx context.Context, peerID string, data []byte) error {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return fmt.Errorf("transport is closed")
	}
	t.mu.RUnlock()

	return t.bridge.Send(ctx, peerID, data)
}

// getMeshPeers returns a list of connected peers in the mesh
func (t *transport) getMeshPeers() []string {
	t.mu.RLock()
	if t.closed {
		t.mu.RUnlock()
		return []string{}
	}
	t.mu.RUnlock()

	return t.bridge.GetMeshPeers()
}

// extractTenantID extracts tenant ID from JWT token
func extractTenantID(token string) string {
	tenantID, err := jwt.ExtractTenantID(token)
	if err != nil {
		// Log error but return default to allow initialization
		log.Printf("Failed to extract tenant ID from token: %v, using default", err)
		return "default-tenant"
	}
	return tenantID
}

// defaultLogger implements a simple logger
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Printf("[INFO] %s %v", msg, fields)
	} else {
		log.Printf("[INFO] %s", msg)
	}
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Printf("[ERROR] %s %v", msg, fields)
	} else {
		log.Printf("[ERROR] %s", msg)
	}
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Printf("[DEBUG] %s %v", msg, fields)
	} else {
		log.Printf("[DEBUG] %s", msg)
	}
}

func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		log.Printf("[WARN] %s %v", msg, fields)
	} else {
		log.Printf("[WARN] %s", msg)
	}
}
