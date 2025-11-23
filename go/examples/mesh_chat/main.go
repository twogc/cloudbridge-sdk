package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

// ChatMessage represents a message in the mesh network
type ChatMessage struct {
	From      string
	To        string
	Content   string
	Timestamp time.Time
}

func main() {
	fmt.Println("CloudBridge SDK - Mesh Chat Example")
	fmt.Println("====================================")
	fmt.Println()

	// Get token and user info
	token := os.Getenv("CLOUDBRIDGE_TOKEN")
	if token == "" {
		token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"
	}

	username := os.Getenv("CLOUDBRIDGE_USER")
	if username == "" {
		username = "anonymous"
	}

	// Create client
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken(token),
		cloudbridge.WithRegion("eu-central"),
		cloudbridge.WithLogLevel("info"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Printf("✓ Chat client created\n")
	fmt.Printf("  Username: %s\n", username)
	fmt.Println()

	// Join mesh network
	ctx := context.Background()
	meshName := "chat-network"

	fmt.Printf("Joining mesh network: %s\n", meshName)
	mesh, err := client.JoinMesh(ctx, meshName)
	if err != nil {
		fmt.Printf("⚠ Failed to join mesh (expected without relay): %v\n", err)
		fmt.Println("  In production, the client would join the mesh network")
	} else {
		defer mesh.Leave()

		fmt.Printf("✓ Joined mesh network\n")
		fmt.Println()

		// Show connected peers
		peers := mesh.Peers()
		if len(peers) > 0 {
			fmt.Printf("Connected peers (%d):\n", len(peers))
			for _, peer := range peers {
				fmt.Printf("  - %s\n", peer)
			}
		}
	}

	fmt.Println()
	fmt.Println("=== Chat Interface ===")
	fmt.Println("Commands:")
	fmt.Println("  /peers         - List connected peers")
	fmt.Println("  /send <msg>    - Broadcast message to all peers")
	fmt.Println("  /quit          - Exit chat")
	fmt.Println()
	fmt.Println("Type your message and press Enter to broadcast:")
	fmt.Println()

	// Start reading messages from stdin
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(input, "/") {
			handleCommand(input, mesh, username)
			continue
		}

		// Send message to mesh
		if mesh != nil {
			sendMessageToMesh(mesh, username, input)
		} else {
			fmt.Println("⚠ Not connected to mesh network")
		}
	}

	fmt.Println("\nGoodbye!")
}

func handleCommand(command string, mesh cloudbridge.Mesh, username string) {
	parts := strings.Fields(command)
	cmd := parts[0]

	switch cmd {
	case "/peers":
		if mesh == nil {
			fmt.Println("⚠ Not connected to mesh")
			return
		}

		peers := mesh.Peers()


		if len(peers) == 0 {
			fmt.Println("No connected peers")
		} else {
			fmt.Printf("Connected peers (%d):\n", len(peers))
			for _, peer := range peers {
				fmt.Printf("  - %s\n", peer)
			}
		}

	case "/send":
		if mesh == nil {
			fmt.Println("⚠ Not connected to mesh")
			return
		}

		if len(parts) < 2 {
			fmt.Println("Usage: /send <message>")
			return
		}

		message := strings.Join(parts[1:], " ")
		sendMessageToMesh(mesh, username, message)

	case "/quit", "/exit":
		fmt.Println("Exiting...")
		os.Exit(0)

	case "/help":
		fmt.Println("Available commands:")
		fmt.Println("  /peers         - List connected peers")
		fmt.Println("  /send <msg>    - Send message")
		fmt.Println("  /quit          - Exit chat")
		fmt.Println("  /help          - Show this help")

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}

func sendMessageToMesh(mesh cloudbridge.Mesh, from string, content string) {
	msg := ChatMessage{
		From:      from,
		Content:   content,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling message: %v\n", err)
		return
	}

	// Send to all peers in mesh
	err = mesh.Broadcast(context.Background(), data)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		return
	}

	fmt.Printf("✓ Message sent to mesh\n")
}

// demonstrateMeshChatLogic shows what would happen with real mesh
func demonstrateMeshChatLogic() {
	fmt.Println("=== Mesh Chat Demo ===")
	fmt.Println()

	// Simulate peers in mesh
	peers := []string{"peer-alice", "peer-bob", "peer-charlie"}

	fmt.Println("Simulated mesh network:")
	for _, peer := range peers {
		fmt.Printf("  - %s\n", peer)
	}
	fmt.Println()

	// Simulate sending a message
	message := ChatMessage{
		From:      "user-you",
		Content:   "Hello, mesh!",
		Timestamp: time.Now(),
	}

	fmt.Printf("Broadcasting message from %s:\n", message.From)
	fmt.Printf("  Content: %s\n", message.Content)
	fmt.Printf("  Timestamp: %s\n", message.Timestamp.Format(time.RFC3339))
	fmt.Println()

	// Simulate receiving responses
	fmt.Println("Simulated responses from peers:")
	responses := map[string]string{
		"peer-alice":   "Hi there!",
		"peer-bob":     "Hello!",
		"peer-charlie": "Nice to meet you!",
	}

	for peer, response := range responses {
		fmt.Printf("  %s: %s\n", peer, response)
	}
	fmt.Println()
}

// MessageBroadcaster handles sending messages to all peers
type MessageBroadcaster struct {
	mesh     cloudbridge.Mesh
	messages chan ChatMessage
	done     chan struct{}
}

// NewMessageBroadcaster creates a new broadcaster
func NewMessageBroadcaster(mesh cloudbridge.Mesh) *MessageBroadcaster {
	return &MessageBroadcaster{
		mesh:     mesh,
		messages: make(chan ChatMessage, 10),
		done:     make(chan struct{}),
	}
}

// Send queues a message for broadcasting
func (b *MessageBroadcaster) Send(msg ChatMessage) error {
	select {
	case b.messages <- msg:
		return nil
	case <-b.done:
		return fmt.Errorf("broadcaster is closed")
	default:
		return fmt.Errorf("message queue full")
	}
}

// Start begins the broadcasting loop
func (b *MessageBroadcaster) Start() {
	go func() {
		for {
			select {
			case msg := <-b.messages:
				// In real implementation, would send to all peers
				fmt.Printf("[%s] %s: %s\n",
					msg.Timestamp.Format("15:04:05"),
					msg.From,
					msg.Content,
				)

			case <-b.done:
				return
			}
		}
	}()
}

// Stop closes the broadcaster
func (b *MessageBroadcaster) Stop() {
	close(b.done)
}
