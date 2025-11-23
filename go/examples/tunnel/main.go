package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

func main() {
	mode := flag.String("mode", "server", "Mode: server (service provider) or client (tunnel initiator)")
	token := flag.String("token", "", "CloudBridge API Token")
	peerID := flag.String("peer", "", "Peer ID to connect to (client mode)")
	localPort := flag.Int("local-port", 8080, "Local port to listen on")
	remotePort := flag.Int("remote-port", 8080, "Remote port to forward to (client mode)")
	region := flag.String("region", "eu-central", "CloudBridge Region")
	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required")
	}

	// Initialize client
	client, err := cloudbridge.NewClient(
		cloudbridge.WithToken(*token),
		cloudbridge.WithRegion(*region),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	if *mode == "server" {
		runServer(ctx, client, *localPort)
	} else {
		if *peerID == "" {
			log.Fatal("Peer ID is required for client mode")
		}
		runClient(ctx, client, *peerID, *localPort, *remotePort)
	}
}

func runServer(ctx context.Context, client *cloudbridge.Client, port int) {
	// Start a simple HTTP server to act as the service
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from CloudBridge Tunnel Service running on port %d!", port)
		})
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}
		log.Printf("Starting local HTTP service on port %d...", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	log.Println("Waiting for incoming tunnel connections...")
	// In a real implementation, client.Serve(ctx) would block and handle incoming connections
	// For this example, we simulate the blocking behavior
	if err := client.Serve(ctx); err != nil {
		log.Printf("Serve error: %v", err)
	}
}

func runClient(ctx context.Context, client *cloudbridge.Client, peerID string, localPort, remotePort int) {
	log.Printf("Creating tunnel: localhost:%d -> %s:%d", localPort, peerID, remotePort)

	tunnel, err := client.CreateTunnel(ctx, cloudbridge.TunnelConfig{
		LocalPort:  localPort,
		RemotePeer: peerID,
		RemotePort: remotePort,
		Protocol:   cloudbridge.ProtocolTCP,
	})
	if err != nil {
		log.Fatalf("Failed to create tunnel: %v", err)
	}
	defer tunnel.Close()

	log.Printf("Tunnel established! Access the service at http://localhost:%d", localPort)
	
	<-ctx.Done()
}
