package cloudbridge

import (
	"errors"
	"io"
	"testing"
	"time"
)

// mockPeerConnection is a mock implementation of bridge.PeerConnection
type mockPeerConnection struct {
	readData  []byte
	readErr   error
	writeErr  error
	closeErr  error
	closed    bool
	readIndex int
}

func (m *mockPeerConnection) Read(b []byte) (int, error) {
	if m.closed {
		return 0, errors.New("connection closed")
	}
	if m.readErr != nil {
		return 0, m.readErr
	}
	if m.readIndex >= len(m.readData) {
		return 0, io.EOF
	}
	n := copy(b, m.readData[m.readIndex:])
	m.readIndex += n
	return n, nil
}

func (m *mockPeerConnection) Write(b []byte) (int, error) {
	if m.closed {
		return 0, errors.New("connection closed")
	}
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(b), nil
}

func (m *mockPeerConnection) Close() error {
	if m.closed {
		return nil
	}
	m.closed = true
	return m.closeErr
}

func TestConnectionPeerID(t *testing.T) {
	conn := &connection{
		peerID: "test-peer-123",
	}

	if conn.PeerID() != "test-peer-123" {
		t.Errorf("PeerID() = %v, want %v", conn.PeerID(), "test-peer-123")
	}
}

func TestConnectionRead(t *testing.T) {
	tests := []struct {
		name      string
		conn      *connection
		data      []byte
		wantErr   bool
		wantBytes int
	}{
		{
			name: "no bridge connection",
			conn: &connection{
				peerID:     "peer-123",
				bridgeConn: nil,
			},
			data:    []byte("test data"),
			wantErr: true,
		},
		{
			name: "closed connection",
			conn: &connection{
				peerID: "peer-123",
				closed: true,
			},
			data:    []byte("test"),
			wantErr: true,
		},
		{
			name: "no bridge connection",
			conn: &connection{
				peerID:     "peer-123",
				bridgeConn: nil,
			},
			data:    []byte("test"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 1024)
			n, err := tt.conn.Read(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && n == 0 {
				t.Error("Read() returned 0 bytes without error")
			}
		})
	}
}

func TestConnectionWrite(t *testing.T) {
	tests := []struct {
		name    string
		conn    *connection
		data    []byte
		wantErr bool
	}{
		{
			name: "closed connection",
			conn: &connection{
				peerID: "peer-123",
				closed: true,
			},
			data:    []byte("test"),
			wantErr: true,
		},
		{
			name: "no bridge connection",
			conn: &connection{
				peerID:     "peer-123",
				bridgeConn: nil,
			},
			data:    []byte("test"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := tt.conn.Write(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && n == 0 {
				t.Error("Write() returned 0 bytes without error")
			}
		})
	}
}

func TestConnectionClose(t *testing.T) {
	tests := []struct {
		name    string
		conn    *connection
		wantErr bool
	}{
		{
			name: "close open connection",
			conn: &connection{
				peerID:  "peer-123",
				closed:  false,
				connected: true,
			},
			wantErr: false,
		},
		{
			name: "close already closed connection",
			conn: &connection{
				peerID:  "peer-123",
				closed:  true,
				connected: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.conn.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.conn.closed {
				t.Error("Close() did not set closed flag")
			}
		})
	}
}

func TestConnectionMetrics(t *testing.T) {
	now := time.Now()
	conn := &connection{
		peerID:       "peer-123",
		connected:    true,
		connectedAt:  now,
		bytesSent:    1000,
		bytesReceived: 2000,
		closed:       false,
	}

	metrics, err := conn.Metrics()
	if err != nil {
		t.Errorf("Metrics() error = %v", err)
		return
	}

	if metrics == nil {
		t.Error("Metrics() returned nil")
		return
	}

	if metrics.BytesSent != 1000 {
		t.Errorf("Metrics().BytesSent = %v, want %v", metrics.BytesSent, 1000)
	}

	if metrics.BytesReceived != 2000 {
		t.Errorf("Metrics().BytesReceived = %v, want %v", metrics.BytesReceived, 2000)
	}

	if !metrics.Connected {
		t.Error("Metrics().Connected = false, want true")
	}

	if metrics.ConnectedAt != now {
		t.Errorf("Metrics().ConnectedAt = %v, want %v", metrics.ConnectedAt, now)
	}

	// Test closed connection
	conn.closed = true
	_, err = conn.Metrics()
	if err == nil {
		t.Error("Metrics() should return error for closed connection")
	}
}

func TestConnectionSetDeadline(t *testing.T) {
	conn := &connection{
		peerID: "peer-123",
	}

	err := conn.SetDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Error("SetDeadline() should return error (not implemented)")
	}
}

func TestConnectionSetReadDeadline(t *testing.T) {
	conn := &connection{
		peerID: "peer-123",
	}

	err := conn.SetReadDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Error("SetReadDeadline() should return error (not implemented)")
	}
}

func TestConnectionSetWriteDeadline(t *testing.T) {
	conn := &connection{
		peerID: "peer-123",
	}

	err := conn.SetWriteDeadline(time.Now().Add(time.Second))
	if err == nil {
		t.Error("SetWriteDeadline() should return error (not implemented)")
	}
}

func TestConnectionBytesTracking(t *testing.T) {
	conn := &connection{
		peerID:     "peer-123",
		connected: true,
	}

	// Simulate bytes sent and received
	conn.mu.Lock()
	conn.bytesSent = 100
	conn.bytesReceived = 200
	conn.mu.Unlock()

	metrics, err := conn.Metrics()
	if err != nil {
		t.Fatalf("Metrics() error = %v", err)
	}

	if metrics.BytesSent != 100 {
		t.Errorf("BytesSent = %v, want %v", metrics.BytesSent, 100)
	}

	if metrics.BytesReceived != 200 {
		t.Errorf("BytesReceived = %v, want %v", metrics.BytesReceived, 200)
	}
}

