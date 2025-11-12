package main

import (
	"context"
	"log"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	// Create client
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken("your-api-token"),
		cloudbridge.WithRegion("eu-central"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Join mesh network
	log.Println("Joining mesh network...")
	mesh, err := client.JoinMesh(ctx, "my-network")
	if err != nil {
		log.Fatalf("Failed to join mesh: %v", err)
	}
	defer mesh.Leave()

	log.Printf("Joined mesh network: %s", mesh.NetworkName())

	// Start message receiver
	go func() {
		for msg := range mesh.Messages() {
			log.Printf("Received message from %s: %s", msg.From, string(msg.Data))
		}
	}()

	// Broadcast messages
	for i := 0; i < 5; i++ {
		message := []byte("Hello mesh!")
		if err := mesh.Broadcast(ctx, message); err != nil {
			log.Printf("Failed to broadcast: %v", err)
		} else {
			log.Printf("Broadcasted message: %s", message)
		}
		time.Sleep(2 * time.Second)
	}

	// List connected peers
	peers := mesh.Peers()
	log.Printf("Connected peers: %v", peers)

	// Send direct message to specific peer
	if len(peers) > 0 {
		peerID := peers[0]
		message := []byte("Direct message")
		if err := mesh.Send(ctx, peerID, message); err != nil {
			log.Printf("Failed to send to peer %s: %v", peerID, err)
		} else {
			log.Printf("Sent direct message to %s", peerID)
		}
	}

	log.Println("Mesh example completed")
}
