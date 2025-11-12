package cloudbridge

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Client represents a CloudBridge SDK client
type Client struct {
	config *Config
	conn   *connection
	mu     sync.RWMutex
	closed bool

	// Callbacks
	onConnect    func(peer string)
	onDisconnect func(peer string, err error)
	onReconnect  func(peer string)
}

// NewClient creates a new CloudBridge client with the given options
func NewClient(opts ...Option) (*Client, error) {
	config := defaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	client := &Client{
		config: config,
	}

	return client, nil
}

// Connect establishes a P2P connection to the specified peer
func (c *Client) Connect(ctx context.Context, peerID string) (Connection, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	if peerID == "" {
		return nil, errors.New("peer ID cannot be empty")
	}

	conn := &connection{
		peerID: peerID,
		client: c,
	}

	if err := conn.dial(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
	}

	if c.onConnect != nil {
		c.onConnect(peerID)
	}

	return conn, nil
}

// CreateTunnel creates a secure tunnel with the specified configuration
func (c *Client) CreateTunnel(ctx context.Context, config TunnelConfig) (Tunnel, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid tunnel configuration: %w", err)
	}

	tunnel := &tunnel{
		config: config,
		client: c,
	}

	if err := tunnel.start(ctx); err != nil {
		return nil, fmt.Errorf("failed to create tunnel: %w", err)
	}

	return tunnel, nil
}

// JoinMesh joins a mesh network with the specified name
func (c *Client) JoinMesh(ctx context.Context, networkName string) (Mesh, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	if networkName == "" {
		return nil, errors.New("network name cannot be empty")
	}

	mesh := &mesh{
		networkName: networkName,
		client:      c,
		messages:    make(chan Message, 100),
	}

	if err := mesh.join(ctx); err != nil {
		return nil, fmt.Errorf("failed to join mesh network %s: %w", networkName, err)
	}

	return mesh, nil
}

// RegisterService registers a service for discovery
func (c *Client) RegisterService(ctx context.Context, config ServiceConfig) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return errors.New("client is closed")
	}
	c.mu.RUnlock()

	if err := config.validate(); err != nil {
		return fmt.Errorf("invalid service configuration: %w", err)
	}

	// TODO: Implement service registration
	return errors.New("not implemented")
}

// DiscoverServices discovers services by name
func (c *Client) DiscoverServices(ctx context.Context, serviceName string) ([]Service, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	if serviceName == "" {
		return nil, errors.New("service name cannot be empty")
	}

	// TODO: Implement service discovery
	return nil, errors.New("not implemented")
}

// Health checks the health of the client connection
func (c *Client) Health(ctx context.Context) (*Health, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	// TODO: Implement health check
	return &Health{
		Status:         "healthy",
		Latency:        10 * time.Millisecond,
		ConnectedPeers: 0,
	}, nil
}

// OnConnect registers a callback for connection events
func (c *Client) OnConnect(callback func(peer string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onConnect = callback
}

// OnDisconnect registers a callback for disconnection events
func (c *Client) OnDisconnect(callback func(peer string, err error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onDisconnect = callback
}

// OnReconnect registers a callback for reconnection events
func (c *Client) OnReconnect(callback func(peer string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReconnect = callback
}

// Close closes the client and releases all resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	if c.conn != nil {
		if err := c.conn.close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

// Health represents the health status of the client
type Health struct {
	Status         string
	Latency        time.Duration
	ConnectedPeers int
}
