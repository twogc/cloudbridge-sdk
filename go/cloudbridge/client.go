package cloudbridge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Client represents a CloudBridge SDK client
type Client struct {
	config    *Config
	transport *transport
	conn      *connection
	mu        sync.RWMutex
	closed    bool
	services  map[string]Service

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
		config:       config,
		services:     make(map[string]Service),
		onConnect:    config.OnConnect,
		onDisconnect: config.OnDisconnect,
		onReconnect:  config.OnReconnect,
	}

	// Initialize transport
	tr, err := newTransport(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}
	client.transport = tr

	// Initialize transport context
	ctx := context.Background()
	if err := tr.initialize(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize transport: %w", err)
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

	// Use transport to connect
	conn, err := c.transport.connectToPeer(ctx, peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to peer %s: %w", peerID, err)
	}
	conn.client = c
	c.conn = conn

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

	// Store service locally
	serviceID := fmt.Sprintf("%s-%s-%d", config.Name, c.config.Region, config.Port)
	service := Service{
		ID:       serviceID,
		Name:     config.Name,
		Port:     config.Port,
		Tags:     config.Tags,
		Healthy:  true,
		PeerID:   c.transport.bridge.GetPeerID(),
		Metadata: map[string]string{"region": c.config.Region},
	}

	c.services[serviceID] = service
	
	// TODO: Broadcast service registration to mesh
	// c.transport.broadcast(...)

	return nil
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

	// Return local services matching the name
	var services []Service
	for _, s := range c.services {
		if s.Name == serviceName {
			services = append(services, s)
		}
	}

	// TODO: Query remote peers for services

	return services, nil
}

// DeregisterService deregisters a service
func (c *Client) DeregisterService(ctx context.Context, serviceID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("client is closed")
	}

	if _, exists := c.services[serviceID]; !exists {
		return fmt.Errorf("service not found: %s", serviceID)
	}

	delete(c.services, serviceID)
	return nil
}

// Health checks the health of the client connection
func (c *Client) Health(ctx context.Context) (*Health, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, errors.New("client is closed")
	}
	c.mu.RUnlock()

	peers := 0
	if c.conn != nil {
		peers = 1
	}

	return &Health{
		Status:         "healthy",
		Latency:        15 * time.Millisecond, // TODO: Get actual latency from transport
		ConnectedPeers: peers,
	}, nil
}

// Serve starts accepting incoming connections and handling tunnels
func (c *Client) Serve(ctx context.Context) error {
	// TODO: Implement real listener from transport
	// For now, we simulate a listener or block until context is done
	// In a real implementation, c.transport should expose an Accept() method
	
	<-ctx.Done()
	return nil
}

// HandleIncomingConnection handles an incoming P2P connection
// This should be called by the transport when a new stream is accepted
func (c *Client) HandleIncomingConnection(conn net.Conn) {
	go func() {
		defer conn.Close()

		// Read handshake
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("failed to read handshake: %v\n", err)
			return
		}

		var handshake struct {
			Type string `json:"type"`
			Port int    `json:"port"`
		}

		if err := json.Unmarshal(buf[:n], &handshake); err != nil {
			// Not a tunnel handshake, treat as generic app connection
			// TODO: Handle generic app connection
			return
		}

		if handshake.Type == "tunnel" {
			// Connect to local service
			localConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", handshake.Port))
			if err != nil {
				fmt.Printf("failed to connect to local service on port %d: %v\n", handshake.Port, err)
				return
			}
			defer localConn.Close()

			// Bidirectional copy
			go io.Copy(localConn, conn)
			io.Copy(conn, localConn)
		}
	}()
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
			// Log error but continue closing transport
			fmt.Printf("failed to close connection: %v\n", err)
		}
	}

	if c.transport != nil {
		if err := c.transport.close(); err != nil {
			return fmt.Errorf("failed to close transport: %w", err)
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
