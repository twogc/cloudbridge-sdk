package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	fmt.Println("CloudBridge SDK - Simple Connection Example")
	fmt.Println("============================================")

	// Create a new CloudBridge client
	// In production, use a real JWT token from your auth provider
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"),
		cloudbridge.WithRegion("eu-central"),
		cloudbridge.WithTimeout(30*time.Second),
		cloudbridge.WithLogLevel("info"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ Client created successfully")
	fmt.Println()

	// Example 1: Connect to a peer
	fmt.Println("Example 1: Connecting to a peer...")
	ctx := context.Background()

	// The peer ID would be obtained from discovery or provided by user
	peerID := "peer-example-123"

	conn, err := client.Connect(ctx, peerID)
	if err != nil {
		// In this example, connection will fail because we don't have a real relay
		// This demonstrates the API usage
		fmt.Printf("⚠ Connection failed (expected without relay): %v\n", err)
		fmt.Println("   In production, this would establish a real P2P connection")
	} else {
		defer conn.Close()

		fmt.Printf("✓ Connected to peer: %s\n", peerID)

		// Get connection metrics
		metrics, err := conn.Metrics()
		if err == nil {
			fmt.Printf("  Peer ID: %s\n", conn.PeerID())
			fmt.Printf("  Connected: %v\n", metrics.Connected)
			fmt.Printf("  Connected At: %s\n", metrics.ConnectedAt.Format(time.RFC3339))
			fmt.Printf("  RTT: %s\n", metrics.RTT)
		}

		// Send data to peer
		message := []byte("Hello from CloudBridge SDK!")
		n, err := conn.Write(message)
		if err == nil {
			fmt.Printf("✓ Sent %d bytes to peer\n", n)
		}

		// Read response from peer
		buffer := make([]byte, 1024)
		n, err = conn.Read(buffer)
		if err == nil {
			fmt.Printf("✓ Received %d bytes from peer: %s\n", n, string(buffer[:n]))
		}
	}

	fmt.Println()

	// Example 2: Health check
	fmt.Println("Example 2: Checking client health...")
	health, err := client.Health(ctx)
	if err != nil {
		fmt.Printf("⚠ Health check failed: %v\n", err)
	} else {
		fmt.Printf("✓ Client health: %+v\n", health)
	}

	fmt.Println()
	fmt.Println("Example completed!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Set up a CloudBridge relay server")
	fmt.Println("2. Obtain a valid JWT token from your auth provider")
	fmt.Println("3. Replace the example token with your real token")
	fmt.Println("4. Update the peer ID with a real peer in your network")
	fmt.Println()
	fmt.Println("For more examples, see: https://github.com/twogc/cloudbridge-sdk/tree/main/go/examples")
}
