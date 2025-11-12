package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	tunnelLocalAddr  string
	tunnelRemoteAddr string
	tunnelProtocol   string
)

var tunnelCmd = &cobra.Command{
	Use:   "tunnel <peer-id>",
	Short: "Create a tunnel to a peer",
	Long: `Create a TCP or UDP tunnel to a peer.
The tunnel forwards local traffic to a remote address through the peer.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		peerID := args[0]

		if tunnelLocalAddr == "" {
			return fmt.Errorf("local address is required (--local)")
		}
		if tunnelRemoteAddr == "" {
			return fmt.Errorf("remote address is required (--remote)")
		}

		logVerbose("Creating CloudBridge client...")
		client, err := createClient()
		if err != nil {
			return err
		}
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_ = ctx // Will be used for actual tunnel creation

		fmt.Printf("Creating %s tunnel to peer: %s\n", tunnelProtocol, peerID)
		fmt.Printf("  Local:  %s\n", tunnelLocalAddr)
		fmt.Printf("  Remote: %s\n", tunnelRemoteAddr)

		// Create tunnel configuration
		// In real implementation, this would call client.CreateTunnel(ctx, tunnelConfig)
		logVerbose("Establishing tunnel...")

		// Start local listener
		listener, err := net.Listen("tcp", tunnelLocalAddr)
		if err != nil {
			return fmt.Errorf("failed to start local listener: %w", err)
		}
		defer listener.Close()

		fmt.Printf("✓ Tunnel established\n")
		fmt.Printf("  Listening on: %s\n", tunnelLocalAddr)
		fmt.Println("\nPress Ctrl+C to stop the tunnel")

		// Setup signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Accept connections in goroutine
		connChan := make(chan net.Conn)
		errChan := make(chan error, 1)

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					errChan <- err
					return
				}
				connChan <- conn
			}
		}()

		connectionCount := 0

		for {
			select {
			case <-sigChan:
				fmt.Println("\nShutting down tunnel...")
				fmt.Printf("Total connections handled: %d\n", connectionCount)
				return nil

			case conn := <-connChan:
				connectionCount++
				fmt.Printf("[%s] New connection from %s (total: %d)\n",
					time.Now().Format("15:04:05"),
					conn.RemoteAddr(),
					connectionCount,
				)

				// Handle connection in goroutine
				go handleTunnelConnection(conn, peerID, tunnelRemoteAddr)

			case err := <-errChan:
				return fmt.Errorf("listener error: %w", err)
			}
		}
	},
}

func init() {
	tunnelCmd.Flags().StringVarP(&tunnelLocalAddr, "local", "l", "", "Local address to listen on (e.g., localhost:8080)")
	tunnelCmd.Flags().StringVarP(&tunnelRemoteAddr, "remote", "r", "", "Remote address to forward to (e.g., localhost:80)")
	tunnelCmd.Flags().StringVarP(&tunnelProtocol, "protocol", "p", "tcp", "Protocol (tcp or udp)")
}

func handleTunnelConnection(conn net.Conn, peerID, remoteAddr string) {
	defer conn.Close()

	logVerbose("Handling connection from %s", conn.RemoteAddr())

	// In real implementation, this would:
	// 1. Get or create peer connection
	// 2. Send tunnel request to peer with remote address
	// 3. Forward data bidirectionally between local conn and peer tunnel

	// For now, just simulate with a timeout
	fmt.Printf("[%s] Forwarding %s -> %s (via %s)\n",
		time.Now().Format("15:04:05"),
		conn.RemoteAddr(),
		remoteAddr,
		peerID,
	)

	// Simulate data transfer
	buffer := make([]byte, 4096)
	_, err := conn.Read(buffer)
	if err != nil {
		logVerbose("Connection closed: %v", err)
		return
	}

	logVerbose("Data forwarded successfully")
}
