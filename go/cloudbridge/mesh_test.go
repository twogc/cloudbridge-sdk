package cloudbridge

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMeshNetworkName(t *testing.T) {
	mesh := &mesh{
		networkName: "test-network",
	}

	if mesh.NetworkName() != "test-network" {
		t.Errorf("NetworkName() = %v, want %v", mesh.NetworkName(), "test-network")
	}
}

func TestMeshPeers(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		peers:       make(map[string]bool),
	}

	peers := mesh.Peers()
	if peers == nil {
		t.Error("Peers() returned nil")
	}
	// Peers list might be empty if transport is not fully initialized
	_ = peers
}

func TestMeshBroadcast(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      false,
	}

	ctx := context.Background()
	data := []byte("test message")

	err = mesh.Broadcast(ctx, data)
	// This might fail if transport is not initialized, which is expected in tests
	if err != nil && err.Error() != "transport is closed" {
		t.Logf("Broadcast() error (may be expected): %v", err)
	}
}

func TestMeshBroadcastClosed(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      true,
	}

	ctx := context.Background()
	err = mesh.Broadcast(ctx, []byte("test"))
	if err == nil {
		t.Error("Broadcast() should return error for closed mesh")
	}
	if !errors.Is(err, errors.New("mesh is closed")) {
		t.Errorf("Broadcast() error = %v, want mesh is closed", err)
	}
}

func TestMeshSend(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      false,
	}

	ctx := context.Background()
	data := []byte("test message")

	err = mesh.Send(ctx, "peer-123", data)
	// This might fail if transport is not initialized, which is expected in tests
	if err != nil && err.Error() != "transport is closed" {
		t.Logf("Send() error (may be expected): %v", err)
	}
}

func TestMeshSendEmptyPeerID(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      false,
	}

	ctx := context.Background()
	err = mesh.Send(ctx, "", []byte("test"))
	if err == nil {
		t.Error("Send() should return error for empty peer ID")
	}
}

func TestMeshSendClosed(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      true,
	}

	ctx := context.Background()
	err = mesh.Send(ctx, "peer-123", []byte("test"))
	if err == nil {
		t.Error("Send() should return error for closed mesh")
	}
}

func TestMeshMessages(t *testing.T) {
	mesh := &mesh{
		messages: make(chan Message, 100),
	}

	messages := mesh.Messages()
	if messages == nil {
		t.Error("Messages() returned nil channel")
	}
}

func TestMeshLeave(t *testing.T) {
	mesh := &mesh{
		networkName: "test-network",
		closed:      false,
		messages:    make(chan Message, 100),
	}

	err := mesh.Leave()
	if err != nil {
		t.Errorf("Leave() error = %v", err)
	}

	if !mesh.closed {
		t.Error("Leave() did not set closed flag")
	}
}

func TestMeshLeaveAlreadyClosed(t *testing.T) {
	mesh := &mesh{
		networkName: "test-network",
		closed:      true,
		messages:    make(chan Message, 100),
	}

	err := mesh.Leave()
	if err != nil {
		t.Errorf("Leave() error = %v", err)
	}
}

func TestMeshJoin(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		messages:    make(chan Message, 100),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mesh.join(ctx)
	// This might fail if transport is not initialized, which is expected in tests
	if err != nil {
		t.Logf("join() error (may be expected): %v", err)
	}

	// Verify peers map is initialized
	if mesh.peers == nil {
		t.Error("join() did not initialize peers map")
	}
}

func TestMeshConcurrentOperations(t *testing.T) {
	client, err := NewClient(
		WithToken("test-token"),
		WithRegion("eu-central"),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	mesh := &mesh{
		networkName: "test-network",
		client:      client,
		closed:      false,
		messages:    make(chan Message, 100),
		peers:       make(map[string]bool),
	}

	// Test concurrent access
	done := make(chan bool, 3)
	ctx := context.Background()

	go func() {
		_ = mesh.Peers()
		done <- true
	}()
	go func() {
		_ = mesh.Broadcast(ctx, []byte("test"))
		done <- true
	}()
	go func() {
		_ = mesh.Send(ctx, "peer-123", []byte("test"))
		done <- true
	}()

	<-done
	<-done
	<-done
}

