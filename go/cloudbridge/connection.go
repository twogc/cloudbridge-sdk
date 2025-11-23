package cloudbridge

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge/internal/bridge"
)

// Connection represents a P2P connection to a peer
type Connection interface {
	io.ReadWriteCloser

	// PeerID returns the peer identifier
	PeerID() string

	// Metrics returns connection metrics
	Metrics() (*ConnectionMetrics, error)

	// SetDeadline sets the read and write deadlines
	SetDeadline(t time.Time) error

	// SetReadDeadline sets the read deadline
	SetReadDeadline(t time.Time) error

	// SetWriteDeadline sets the write deadline
	SetWriteDeadline(t time.Time) error
}

// connection implements the Connection interface
type connection struct {
	peerID string
	client *Client
	mu     sync.RWMutex
	closed bool

	// Connection state
	connected   bool
	connectedAt time.Time

	// Metrics
	bytesSent     uint64
	bytesReceived uint64

	// Underlying bridge connection
	bridgeConn *bridge.PeerConnection
}

// dial establishes a connection to the peer
func (c *connection) dial(ctx context.Context) error {
	// This method is deprecated in favor of transport.connectToPeer
	return errors.New("use transport.connectToPeer instead")
}

// Read reads data from the connection
func (c *connection) Read(b []byte) (int, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return 0, errors.New("connection is closed")
	}
	if c.bridgeConn == nil {
		c.mu.RUnlock()
		return 0, errors.New("connection not established")
	}
	c.mu.RUnlock()

	n, err := c.bridgeConn.Read(b)

	c.mu.Lock()
	c.bytesReceived += uint64(n)
	c.mu.Unlock()

	return n, err
}

// Write writes data to the connection
func (c *connection) Write(b []byte) (int, error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return 0, errors.New("connection is closed")
	}
	if c.bridgeConn == nil {
		c.mu.RUnlock()
		return 0, errors.New("connection not established")
	}
	c.mu.RUnlock()

	n, err := c.bridgeConn.Write(b)

	c.mu.Lock()
	c.bytesSent += uint64(n)
	c.mu.Unlock()

	return n, err
}

// Close closes the connection
func (c *connection) Close() error {
	return c.close()
}

// close is the internal close method
func (c *connection) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	c.connected = false

	var err error
	if c.bridgeConn != nil {
		err = c.bridgeConn.Close()
	}

	if c.client != nil && c.client.onDisconnect != nil {
		c.client.onDisconnect(c.peerID, err)
	}

	return err
}

// PeerID returns the peer identifier
func (c *connection) PeerID() string {
	return c.peerID
}

// Metrics returns connection metrics
func (c *connection) Metrics() (*ConnectionMetrics, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, errors.New("connection is closed")
	}

	return &ConnectionMetrics{
		BytesSent:     c.bytesSent,
		BytesReceived: c.bytesReceived,
		RTT:           10 * time.Millisecond, // TODO: Get actual RTT
		Connected:     c.connected,
		ConnectedAt:   c.connectedAt,
	}, nil
}

// SetDeadline sets the read and write deadlines
func (c *connection) SetDeadline(t time.Time) error {
	// TODO: Implement deadline setting
	return errors.New("not implemented")
}

// SetReadDeadline sets the read deadline
func (c *connection) SetReadDeadline(t time.Time) error {
	// TODO: Implement read deadline setting
	return errors.New("not implemented")
}

// SetWriteDeadline sets the write deadline
func (c *connection) SetWriteDeadline(t time.Time) error {
	// TODO: Implement write deadline setting
	return errors.New("not implemented")
}

// ConnectionMetrics represents metrics for a connection
type ConnectionMetrics struct {
	BytesSent     uint64
	BytesReceived uint64
	RTT           time.Duration
	Connected     bool
	ConnectedAt   time.Time
}
