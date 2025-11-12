package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	connectInteractive bool
	connectMessage     string
)

var connectCmd = &cobra.Command{
	Use:   "connect <peer-id>",
	Short: "Connect to a peer",
	Long: `Connect to a peer by ID and establish a bidirectional connection.
You can send messages interactively or pass a single message to send.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		peerID := args[0]

		logVerbose("Creating CloudBridge client...")
		client, err := createClient()
		if err != nil {
			return err
		}
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		fmt.Printf("Connecting to peer: %s\n", peerID)
		logVerbose("Establishing connection...")

		conn, err := client.Connect(ctx, peerID)
		if err != nil {
			return fmt.Errorf("failed to connect to peer: %w", err)
		}
		defer conn.Close()

		fmt.Printf("✓ Connected to peer %s\n", peerID)

		// Display connection metrics
		metrics, err := conn.Metrics()
		if err == nil {
			fmt.Printf("  Connected at: %s\n", metrics.ConnectedAt.Format(time.RFC3339))
			if metrics.RTT > 0 {
				fmt.Printf("  RTT: %s\n", metrics.RTT)
			}
		}

		// Handle single message or interactive mode
		if connectMessage != "" {
			return sendSingleMessage(conn, connectMessage)
		}

		if connectInteractive {
			return runInteractiveMode(conn, peerID)
		}

		fmt.Println("Connection established. Use --interactive to send messages or --message to send a single message.")
		return nil
	},
}

func init() {
	connectCmd.Flags().BoolVarP(&connectInteractive, "interactive", "i", false, "Interactive mode for sending messages")
	connectCmd.Flags().StringVarP(&connectMessage, "message", "m", "", "Single message to send")
}

func sendSingleMessage(conn io.ReadWriter, message string) error {
	logVerbose("Sending message: %s", message)

	n, err := conn.Write([]byte(message + "\n"))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("✓ Sent %d bytes\n", n)

	// Wait for response
	fmt.Println("Waiting for response...")

	buffer := make([]byte, 4096)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	var readErr error

	go func() {
		n, readErr = conn.Read(buffer)
		close(done)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("No response received (timeout)")
	case <-done:
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("failed to read response: %w", readErr)
		}
		if n > 0 {
			fmt.Printf("✓ Received response (%d bytes):\n%s\n", n, string(buffer[:n]))
		}
	}

	return nil
}

func runInteractiveMode(conn io.ReadWriter, peerID string) error {
	fmt.Println("\nInteractive mode - type messages to send (Ctrl+C to exit)")
	fmt.Println("Commands: /quit, /exit, /metrics")
	fmt.Println()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Channel for incoming messages
	incoming := make(chan []byte, 10)
	errChan := make(chan error, 1)

	// Start reader goroutine
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				return
			}
			if n > 0 {
				msg := make([]byte, n)
				copy(msg, buffer[:n])
				incoming <- msg
			}
		}
	}()

	// Start input reader
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return nil

		case msg := <-incoming:
			fmt.Printf("\n[%s] %s\n", peerID, string(msg))
			fmt.Print("> ")

		case err := <-errChan:
			return fmt.Errorf("connection error: %w", err)

		default:
			fmt.Print("> ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("input error: %w", err)
				}
				return nil
			}

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			// Handle commands
			if strings.HasPrefix(input, "/") {
				if err := handleCommand(input, conn); err != nil {
					if err.Error() == "quit" {
						return nil
					}
					fmt.Printf("Error: %v\n", err)
				}
				continue
			}

			// Send message
			_, err := conn.Write([]byte(input + "\n"))
			if err != nil {
				return fmt.Errorf("failed to send message: %w", err)
			}

			logVerbose("Sent: %s", input)
		}
	}
}

func handleCommand(cmd string, conn io.ReadWriter) error {
	switch strings.ToLower(cmd) {
	case "/quit", "/exit":
		fmt.Println("Exiting...")
		return fmt.Errorf("quit")

	case "/metrics":
		if metricsConn, ok := conn.(interface{ Metrics() (interface{}, error) }); ok {
			metrics, err := metricsConn.Metrics()
			if err != nil {
				return err
			}
			fmt.Printf("Metrics: %+v\n", metrics)
		} else {
			fmt.Println("Metrics not available")
		}

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}

	return nil
}
