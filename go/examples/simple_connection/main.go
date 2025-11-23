package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	token := flag.String("token", "", "CloudBridge API Token")
	peerID := flag.String("peer", "", "Peer ID to connect to")
	region := flag.String("region", "eu-central", "CloudBridge Region")
	flag.Parse()

	if *token == "" || *peerID == "" {
		log.Fatal("Token and Peer ID are required")
	}

	// Initialize client
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken(*token),
		cloudbridge.WithRegion(*region),
		cloudbridge.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Connect to peer
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("Connecting to peer %s...\n", *peerID)
	conn, err := client.Connect(ctx, *peerID)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s!\n", conn.PeerID())

	// Send data
	msg := []byte("Hello from CloudBridge Go SDK!")
	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}
	fmt.Printf("Sent: %s\n", string(msg))

	// Read response
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
	} else {
		fmt.Printf("Received: %s\n", string(buf[:n]))
	}
}
