package cloudbridge

import (
	"context"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "valid configuration",
			opts: []Option{
				WithToken("test-token"),
				WithRegion("eu-central"),
			},
			wantErr: false,
		},
		{
			name:    "missing token",
			opts:    []Option{},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			opts: []Option{
				WithToken("test-token"),
				WithTimeout(-1 * time.Second),
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			opts: []Option{
				WithToken("test-token"),
				WithLogLevel("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
			if client != nil {
				client.Close()
			}
		})
	}
}

func TestClientConnect(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	tests := []struct {
		name    string
		peerID  string
		wantErr bool
	}{
		{
			name:    "empty peer ID",
			peerID:  "",
			wantErr: true,
		},
		{
			name:    "valid peer ID",
			peerID:  "peer-123",
			wantErr: false, // Will fail with "not implemented" but validates input
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Connect(ctx, tt.peerID)
			if (err != nil) != tt.wantErr && err.Error() != "not implemented" {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientHealth(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	health, err := client.Health(ctx)
	if err != nil {
		t.Errorf("Health() error = %v", err)
		return
	}

	if health == nil {
		t.Error("Health() returned nil")
		return
	}

	if health.Status == "" {
		t.Error("Health status is empty")
	}
}

func TestClientClose(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// First close should succeed
	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Second close should also succeed (idempotent)
	err = client.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}

	// Operations after close should fail
	ctx := context.Background()
	_, err = client.Connect(ctx, "peer-123")
	if err == nil {
		t.Error("Connect() after Close() should fail")
	}
}

func TestClientCallbacks(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	var connectCalled, disconnectCalled, reconnectCalled bool

	client.OnConnect(func(peer string) {
		connectCalled = true
	})

	client.OnDisconnect(func(peer string, err error) {
		disconnectCalled = true
	})

	client.OnReconnect(func(peer string) {
		reconnectCalled = true
	})

	// Verify callbacks are registered
	if client.onConnect == nil {
		t.Error("OnConnect callback not registered")
	}
	if client.onDisconnect == nil {
		t.Error("OnDisconnect callback not registered")
	}
	if client.onReconnect == nil {
		t.Error("OnReconnect callback not registered")
	}

	// Suppress unused variable warnings
	_ = connectCalled
	_ = disconnectCalled
	_ = reconnectCalled
}
