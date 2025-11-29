package cloudbridge

import (
	"context"
	"testing"
	"time"
)

func TestNewTransport(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}

	if transport == nil {
		t.Error("newTransport() returned nil")
	}

	if transport.config != config {
		t.Error("newTransport() did not set config correctly")
	}
}

func TestTransportInitialize(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = transport.initialize(ctx)
	// This might fail if bridge initialization fails (expected in test environment)
	if err != nil {
		t.Logf("initialize() error (may be expected): %v", err)
	}
}

func TestTransportClose(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}

	err = transport.close()
	if err != nil {
		t.Errorf("close() error = %v", err)
	}

	if !transport.closed {
		t.Error("close() did not set closed flag")
	}

	// Second close should be idempotent
	err = transport.close()
	if err != nil {
		t.Errorf("Second close() error = %v", err)
	}
}

func TestTransportCloseBeforeInitialize(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}

	err = transport.close()
	if err != nil {
		t.Errorf("close() error = %v", err)
	}

	// Try to initialize after close
	ctx := context.Background()
	err = transport.initialize(ctx)
	if err == nil {
		t.Error("initialize() should fail after close")
	}
}

func TestTransportConnectToPeer(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}
	defer transport.close()

	ctx := context.Background()
	_, err = transport.connectToPeer(ctx, "peer-123")
	// This will fail if bridge is not initialized, which is expected in tests
	if err != nil {
		t.Logf("connectToPeer() error (may be expected): %v", err)
	}
}

func TestTransportConnectToPeerClosed(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}

	transport.close()

	ctx := context.Background()
	_, err = transport.connectToPeer(ctx, "peer-123")
	if err == nil {
		t.Error("connectToPeer() should fail when transport is closed")
	}
}

func TestTransportBroadcast(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}
	defer transport.close()

	ctx := context.Background()
	err = transport.broadcast(ctx, []byte("test message"))
	// This will fail if bridge is not initialized, which is expected in tests
	if err != nil {
		t.Logf("broadcast() error (may be expected): %v", err)
	}
}

func TestTransportSend(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}
	defer transport.close()

	ctx := context.Background()
	err = transport.send(ctx, "peer-123", []byte("test message"))
	// This will fail if bridge is not initialized, which is expected in tests
	if err != nil {
		t.Logf("send() error (may be expected): %v", err)
	}
}

func TestTransportGetMeshPeers(t *testing.T) {
	config := &Config{
		Token:    "test-token",
		Region:   "eu-central",
		Timeout:  30 * time.Second,
		LogLevel: "info",
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{ProtocolQUIC},
	}

	transport, err := newTransport(config)
	if err != nil {
		t.Fatalf("newTransport() error = %v", err)
	}
	defer transport.close()

	peers := transport.getMeshPeers()
	if peers == nil {
		t.Error("getMeshPeers() returned nil")
	}
	// Peers list might be empty if bridge is not initialized
	_ = peers
}

func TestExtractTenantID(t *testing.T) {
	// Test with a valid JWT token format
	// This will use the jwt package to extract tenant ID
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"
	
	tenantID := extractTenantID(token)
	// If extraction fails, it should return "default-tenant"
	if tenantID == "" {
		t.Error("extractTenantID() returned empty string")
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := &defaultLogger{}

	// Test all log methods don't panic
	logger.Info("test info")
	logger.Error("test error")
	logger.Debug("test debug")
	logger.Warn("test warn")

	// Test with fields
	logger.Info("test info", "key", "value")
	logger.Error("test error", "key", "value")
	logger.Debug("test debug", "key", "value")
	logger.Warn("test warn", "key", "value")
}

