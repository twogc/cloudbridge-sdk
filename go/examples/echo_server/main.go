package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	fmt.Println("CloudBridge SDK - Echo Server Example")
	fmt.Println("=====================================")
	fmt.Println()

	// Get token from environment or use example token
	token := os.Getenv("CLOUDBRIDGE_TOKEN")
	if token == "" {
		token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"
		fmt.Println("⚠ Using example token. Set CLOUDBRIDGE_TOKEN for production.")
	}

	// Create client with callbacks
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken(token),
		cloudbridge.WithRegion("eu-central"),
		cloudbridge.WithTimeout(30*time.Second),
		cloudbridge.WithLogLevel("info"),
		cloudbridge.WithOnConnect(func(peerID string) {
			fmt.Printf("🔗 Peer connected: %s\n", peerID)
		}),
		cloudbridge.WithOnDisconnect(func(peerID string, err error) {
			if err != nil {
				fmt.Printf("❌ Peer disconnected: %s (error: %v)\n", peerID, err)
			} else {
				fmt.Printf("👋 Peer disconnected: %s\n", peerID)
			}
		}),
		cloudbridge.WithOnReconnect(func(peerID string) {
			fmt.Printf("🔄 Peer reconnected: %s\n", peerID)
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("✓ Echo server initialized")
	fmt.Println("  Waiting for peer connections...")
	fmt.Println("  Press Ctrl+C to stop")
	fmt.Println()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// In a real implementation, we would:
	// 1. Listen for incoming connections
	// 2. Handle multiple peers concurrently
	// 3. Echo back any received data

	// Simulate waiting for connections
	ctx := context.Background()

	// Example: Accept connections (pseudo-code, actual implementation would be different)
	go func() {
		// This would be implemented when the SDK supports server mode
		fmt.Println("Server mode not yet implemented in SDK")
		fmt.Println("Once implemented, the server will:")
		fmt.Println("  1. Listen for incoming peer connections")
		fmt.Println("  2. Accept connections automatically")
		fmt.Println("  3. Echo back any received data")
		fmt.Println("  4. Handle multiple peers concurrently")
	}()

	// Example of what client-side connection would look like
	fmt.Println("\nTo connect to this echo server from another peer:")
	fmt.Println("  1. Get this peer's ID from discovery")
	fmt.Println("  2. Use client.Connect(ctx, peerID)")
	fmt.Println("  3. Write data with conn.Write()")
	fmt.Println("  4. Read echoed data with conn.Read()")
	fmt.Println()

	// Demonstrate handling a mock connection
	demonstrateEchoLogic()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\n\nShutting down echo server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Close all connections
	// In real implementation:
	// - Close all active peer connections
	// - Wait for pending operations
	// - Cleanup resources

	select {
	case <-shutdownCtx.Done():
		fmt.Println("⚠ Shutdown timeout, forcing exit")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("✓ Server stopped gracefully")
	}
}

// demonstrateEchoLogic shows how echo functionality would work
func demonstrateEchoLogic() {
	fmt.Println("=== Echo Logic Demo ===")
	fmt.Println()

	// Simulate handling a connection
	handleConnection := func(peerID string) {
		fmt.Printf("Handling connection from: %s\n", peerID)

		// In real implementation, this would be the actual connection
		// For demo, we'll use mock data
		mockData := [][]byte{
			[]byte("Hello, server!"),
			[]byte("Test message 2"),
			[]byte("Goodbye!"),
		}

		for i, data := range mockData {
			// Simulate receiving data
			fmt.Printf("  [%d] Received %d bytes: %s\n", i+1, len(data), string(data))

			// Echo back
			fmt.Printf("  [%d] Echoing back %d bytes\n", i+1, len(data))

			// In real implementation:
			// n, err := conn.Write(data)
			// if err != nil {
			//     log.Printf("Failed to echo: %v", err)
			//     return
			// }
		}

		fmt.Printf("Connection from %s handled successfully\n", peerID)
	}

	// Demonstrate with mock peer
	handleConnection("example-peer-123")
	fmt.Println()
}

// echoHandler handles echoing data for a single connection
func echoHandler(conn io.ReadWriteCloser, peerID string) {
	defer conn.Close()

	fmt.Printf("Echo handler started for peer: %s\n", peerID)

	buffer := make([]byte, 4096)

	for {
		// Read data from peer
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error from %s: %v", peerID, err)
			}
			break
		}

		if n == 0 {
			continue
		}

		data := buffer[:n]
		fmt.Printf("📨 Received %d bytes from %s: %s\n", n, peerID, string(data))

		// Echo back
		written, err := conn.Write(data)
		if err != nil {
			log.Printf("Write error to %s: %v", peerID, err)
			break
		}

		fmt.Printf("📤 Echoed %d bytes to %s\n", written, peerID)
	}

	fmt.Printf("Echo handler stopped for peer: %s\n", peerID)
}
