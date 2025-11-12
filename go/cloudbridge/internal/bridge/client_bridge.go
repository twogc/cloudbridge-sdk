package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/2gc-dev/cloudbridge-client/pkg/api"
	"github.com/2gc-dev/cloudbridge-client/pkg/auth"
	"github.com/2gc-dev/cloudbridge-client/pkg/p2p"
	"github.com/2gc-dev/cloudbridge-client/pkg/quic"
)

// ClientBridge provides integration between SDK and CloudBridge Relay Client
type ClientBridge struct {
	config      *BridgeConfig
	p2pManager  *p2p.Manager
	quicConn    *quic.QUICConnection
	apiManager  *api.Manager
	authManager *auth.AuthManager
	logger      Logger
}

// BridgeConfig holds configuration for the bridge
type BridgeConfig struct {
	Token              string
	RelayServerURL     string
	TenantID           string
	InsecureSkipVerify bool
	Timeout            time.Duration
	EnableP2P          bool
	EnableMesh         bool
}

// Logger interface for bridge logging
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}

// NewClientBridge creates a new bridge to CloudBridge Relay Client
func NewClientBridge(config *BridgeConfig, logger Logger) (*ClientBridge, error) {
	if config == nil {
		return nil, fmt.Errorf("bridge configuration is required")
	}

	if config.Token == "" {
		return nil, fmt.Errorf("authentication token is required")
	}

	if config.TenantID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	return &ClientBridge{
		config: config,
		logger: logger,
	}, nil
}

// Initialize sets up the bridge components
func (b *ClientBridge) Initialize(ctx context.Context) error {
	b.logger.Info("Initializing CloudBridge client bridge")

	// Create auth manager
	authConfig := &auth.AuthConfig{
		Type:           "jwt",
		SkipValidation: b.config.InsecureSkipVerify,
	}

	authManager, err := auth.NewAuthManager(authConfig)
	if err != nil {
		return fmt.Errorf("failed to create auth manager: %w", err)
	}
	b.authManager = authManager

	// Create API manager
	apiConfig := &api.ManagerConfig{
		BaseURL:            b.config.RelayServerURL,
		InsecureSkipVerify: b.config.InsecureSkipVerify,
		Timeout:            b.config.Timeout,
		MaxRetries:         3,
		BackoffMultiplier:  2.0,
		MaxBackoff:         30 * time.Second,
		Token:              b.config.Token,
		TenantID:           b.config.TenantID,
		HeartbeatInterval:  30 * time.Second,
	}

	b.apiManager = api.NewManager(apiConfig, authManager, b.logger)

	// Start API manager
	if err := b.apiManager.Start(); err != nil {
		return fmt.Errorf("failed to start API manager: %w", err)
	}

	// Create QUIC connection if needed
	if b.config.EnableP2P {
		b.quicConn = quic.NewQUICConnection(b.logger)
		b.quicConn.SetInsecureSkipVerify(b.config.InsecureSkipVerify)
	}

	// Create P2P manager if enabled
	if b.config.EnableP2P || b.config.EnableMesh {
		p2pConfig := &p2p.P2PConfig{
			TenantID:          b.config.TenantID,
			ConnectionType:    "quic+ice",
			HeartbeatInterval: 30 * time.Second,
			MeshConfig: &p2p.MeshConfig{
				HeartbeatInterval: "30s",
			},
		}

		b.p2pManager = p2p.NewManagerWithAPI(p2pConfig, apiConfig, authManager, b.config.Token, b.logger)

		if err := b.p2pManager.Start(); err != nil {
			return fmt.Errorf("failed to start P2P manager: %w", err)
		}
	}

	b.logger.Info("CloudBridge client bridge initialized successfully")
	return nil
}

// ConnectToPeer establishes a connection to a peer
func (b *ClientBridge) ConnectToPeer(ctx context.Context, peerID string) (*PeerConnection, error) {
	if b.p2pManager == nil {
		return nil, fmt.Errorf("P2P manager not initialized")
	}

	b.logger.Info("Connecting to peer", "peer_id", peerID)

	// Use P2P manager to establish connection
	// This is a simplified implementation - actual connection logic will be more complex
	conn := &PeerConnection{
		PeerID:      peerID,
		ConnectedAt: time.Now(),
		bridge:      b,
	}

	return conn, nil
}

// DiscoverPeers discovers available peers in the network
func (b *ClientBridge) DiscoverPeers(ctx context.Context) ([]*api.Peer, error) {
	if b.apiManager == nil {
		return nil, fmt.Errorf("API manager not initialized")
	}

	resp, err := b.apiManager.GetClient().DiscoverPeers(ctx, b.config.TenantID, b.config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to discover peers: %w", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("peer discovery failed: %s", resp.Error)
	}

	return resp.Peers, nil
}

// GetPeerID returns the local peer ID
func (b *ClientBridge) GetPeerID() string {
	if b.apiManager != nil {
		return b.apiManager.GetPeerID()
	}
	return ""
}

// GetTenantID returns the tenant ID
func (b *ClientBridge) GetTenantID() string {
	return b.config.TenantID
}

// Close closes the bridge and releases resources
func (b *ClientBridge) Close() error {
	b.logger.Info("Closing CloudBridge client bridge")

	var errs []error

	if b.p2pManager != nil {
		if err := b.p2pManager.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("failed to stop P2P manager: %w", err))
		}
	}

	if b.apiManager != nil {
		b.apiManager.Stop()
	}

	if b.quicConn != nil {
		if err := b.quicConn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close QUIC connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing bridge: %v", errs)
	}

	return nil
}

// PeerConnection represents a connection to a peer through the bridge
type PeerConnection struct {
	PeerID      string
	ConnectedAt time.Time
	bridge      *ClientBridge
}

// Read reads data from the peer connection
func (pc *PeerConnection) Read(b []byte) (int, error) {
	// TODO: Implement actual read from QUIC stream
	return 0, fmt.Errorf("not implemented")
}

// Write writes data to the peer connection
func (pc *PeerConnection) Write(b []byte) (int, error) {
	// TODO: Implement actual write to QUIC stream
	return 0, fmt.Errorf("not implemented")
}

// Close closes the peer connection
func (pc *PeerConnection) Close() error {
	// TODO: Implement connection close
	return nil
}
