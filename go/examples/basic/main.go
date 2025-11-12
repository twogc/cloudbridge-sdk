package main

import (
	"context"
	"log"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	// Create client with authentication token
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken("your-api-token"),
		cloudbridge.WithRegion("eu-central"),
		cloudbridge.WithTimeout(30*time.Second),
		cloudbridge.WithLogLevel("info"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	log.Println("CloudBridge client created successfully")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: Connect to peer
	log.Println("Connecting to peer...")
	conn, err := client.Connect(ctx, "peer-id-123")
	if err != nil {
		log.Fatalf("Failed to connect to peer: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to peer: %s", conn.PeerID())

	// Get connection metrics
	metrics, err := conn.Metrics()
	if err != nil {
		log.Printf("Failed to get metrics: %v", err)
	} else {
		log.Printf("Connection metrics: Bytes sent: %d, Bytes received: %d, RTT: %v",
			metrics.BytesSent, metrics.BytesReceived, metrics.RTT)
	}

	// Example 2: Create tunnel
	log.Println("Creating tunnel...")
	tunnel, err := client.CreateTunnel(ctx, cloudbridge.TunnelConfig{
		LocalPort:  8080,
		RemotePeer: "peer-id-456",
		RemotePort: 3000,
		Protocol:   cloudbridge.ProtocolTCP,
	})
	if err != nil {
		log.Fatalf("Failed to create tunnel: %v", err)
	}
	defer tunnel.Close()

	log.Printf("Tunnel created: localhost:%d -> %s:%d",
		tunnel.LocalPort(), tunnel.RemotePeer(), tunnel.RemotePort())

	// Example 3: Health check
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("Failed to get health: %v", err)
	}

	log.Printf("Health: Status=%s, Latency=%v, Connected Peers=%d",
		health.Status, health.Latency, health.ConnectedPeers)

	log.Println("Example completed successfully")
}
