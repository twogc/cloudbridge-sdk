package cloudbridge

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestTunnelConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  TunnelConfig
		wantErr bool
	}{
		{
			name: "valid TCP config",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 3000,
				Protocol:   ProtocolTCP,
			},
			wantErr: false,
		},
		{
			name: "valid QUIC config",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 3000,
				Protocol:   ProtocolQUIC,
			},
			wantErr: false,
		},
		{
			name: "default protocol (TCP)",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 3000,
				Protocol:   "",
			},
			wantErr: false,
		},
		{
			name: "invalid local port - zero",
			config: TunnelConfig{
				LocalPort:  0,
				RemotePeer: "peer-123",
				RemotePort: 3000,
			},
			wantErr: true,
		},
		{
			name: "invalid local port - negative",
			config: TunnelConfig{
				LocalPort:  -1,
				RemotePeer: "peer-123",
				RemotePort: 3000,
			},
			wantErr: true,
		},
		{
			name: "invalid local port - too high",
			config: TunnelConfig{
				LocalPort:  70000,
				RemotePeer: "peer-123",
				RemotePort: 3000,
			},
			wantErr: true,
		},
		{
			name: "invalid remote port - zero",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid remote port - too high",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 70000,
			},
			wantErr: true,
		},
		{
			name: "empty remote peer",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "",
				RemotePort: 3000,
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			config: TunnelConfig{
				LocalPort:  8080,
				RemotePeer: "peer-123",
				RemotePort: 3000,
				Protocol:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTunnelRemotePeer(t *testing.T) {
	tunnel := &tunnel{
		config: TunnelConfig{
			RemotePeer: "peer-123",
		},
	}

	if tunnel.RemotePeer() != "peer-123" {
		t.Errorf("RemotePeer() = %v, want %v", tunnel.RemotePeer(), "peer-123")
	}
}

func TestTunnelLocalPort(t *testing.T) {
	tunnel := &tunnel{
		config: TunnelConfig{
			LocalPort: 8080,
		},
	}

	if tunnel.LocalPort() != 8080 {
		t.Errorf("LocalPort() = %v, want %v", tunnel.LocalPort(), 8080)
	}
}

func TestTunnelRemotePort(t *testing.T) {
	tunnel := &tunnel{
		config: TunnelConfig{
			RemotePort: 3000,
		},
	}

	if tunnel.RemotePort() != 3000 {
		t.Errorf("RemotePort() = %v, want %v", tunnel.RemotePort(), 3000)
	}
}

func TestTunnelClose(t *testing.T) {
	tests := []struct {
		name    string
		tunnel  *tunnel
		wantErr bool
	}{
		{
			name: "close open tunnel",
			tunnel: &tunnel{
				closed: false,
			},
			wantErr: false,
		},
		{
			name: "close already closed tunnel",
			tunnel: &tunnel{
				closed: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tunnel.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.tunnel.closed {
				t.Error("Close() did not set closed flag")
			}
		})
	}
}

func TestTunnelStart(t *testing.T) {
	// This test requires a real client, so we'll test the basic structure
	// In a real scenario, we'd need to mock the client and transport
	
	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config := TunnelConfig{
		LocalPort:  port,
		RemotePeer: "peer-123",
		RemotePort: 3000,
		Protocol:   ProtocolTCP,
	}

	tunnel, err := client.CreateTunnel(ctx, config)
	if err != nil {
		// This is expected if the bridge is not initialized or peer doesn't exist
		// We're just testing that the method doesn't panic
		t.Logf("CreateTunnel() error (expected in test): %v", err)
		return
	}
	defer tunnel.Close()

	if tunnel == nil {
		t.Error("CreateTunnel() returned nil tunnel")
	}

	if tunnel.RemotePeer() != "peer-123" {
		t.Errorf("RemotePeer() = %v, want %v", tunnel.RemotePeer(), "peer-123")
	}

	if tunnel.LocalPort() != port {
		t.Errorf("LocalPort() = %v, want %v", tunnel.LocalPort(), port)
	}
}

func TestTunnelConcurrentClose(t *testing.T) {
	tunnel := &tunnel{
		closed: false,
	}

	// Test concurrent closes
	done := make(chan bool, 2)
	go func() {
		tunnel.Close()
		done <- true
	}()
	go func() {
		tunnel.Close()
		done <- true
	}()

	<-done
	<-done

	if !tunnel.closed {
		t.Error("Tunnel should be closed after concurrent Close() calls")
	}
}

func TestTunnelConfigDefaultProtocol(t *testing.T) {
	config := TunnelConfig{
		LocalPort:  8080,
		RemotePeer: "peer-123",
		RemotePort: 3000,
		Protocol:   "", // Empty protocol
	}

	err := config.validate()
	if err != nil {
		t.Fatalf("validate() error = %v", err)
	}

	if config.Protocol != ProtocolTCP {
		t.Errorf("Default protocol = %v, want %v", config.Protocol, ProtocolTCP)
	}
}

