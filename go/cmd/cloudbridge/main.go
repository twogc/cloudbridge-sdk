package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/twogc/cloudbridge-sdk/go/cloudbridge"
)

var (
	token              string
	region             string
	timeout            time.Duration
	insecureSkipVerify bool
	verbose            bool
)

var rootCmd = &cobra.Command{
	Use:   "cloudbridge",
	Short: "CloudBridge SDK CLI Tool",
	Long: `CloudBridge CLI is a command-line tool for testing and interacting with the CloudBridge SDK.
It provides commands to connect to peers, discover peers, create tunnels, and check health.`,
	Version: "0.1.0",
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "CloudBridge authentication token (or set CLOUDBRIDGE_TOKEN)")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "eu-central", "CloudBridge region")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Operation timeout")
	rootCmd.PersistentFlags().BoolVar(&insecureSkipVerify, "insecure-skip-verify", false, "Skip TLS certificate verification")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add commands
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(tunnelCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// getToken retrieves token from flag or environment
func getToken() (string, error) {
	if token != "" {
		return token, nil
	}

	token = os.Getenv("CLOUDBRIDGE_TOKEN")
	if token == "" {
		return "", fmt.Errorf("token is required (use --token or set CLOUDBRIDGE_TOKEN)")
	}

	return token, nil
}

// createClient creates a new CloudBridge client with configured options
func createClient() (*cloudbridge.Client, error) {
	tok, err := getToken()
	if err != nil {
		return nil, err
	}

	opts := []cloudbridge.Option{
		cloudbridge.WithToken(tok),
		cloudbridge.WithRegion(region),
		cloudbridge.WithTimeout(timeout),
		cloudbridge.WithInsecureSkipVerify(insecureSkipVerify),
	}

	if verbose {
		opts = append(opts, cloudbridge.WithLogLevel("debug"))
	}

	client, err := cloudbridge.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// logVerbose prints message only in verbose mode
func logVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("CloudBridge CLI version %s\n", rootCmd.Version)
		fmt.Println("SDK version: 0.1.0")
		fmt.Println("Build date: November 2025")
	},
}
