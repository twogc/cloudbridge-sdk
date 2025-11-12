package cloudbridge

import (
	"errors"
	"os"
	"time"
)

// Config holds the configuration for the CloudBridge client
type Config struct {
	// Authentication
	Token string

	// Region specifies the preferred CloudBridge region
	Region string

	// Timeout for operations
	Timeout time.Duration

	// Log level (debug, info, warn, error)
	LogLevel string

	// Retry policy
	RetryPolicy RetryPolicy

	// Protocol preferences
	Protocols []Protocol

	// TLS configuration
	InsecureSkipVerify bool
}

// RetryPolicy defines retry behavior for failed operations
type RetryPolicy struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// Protocol represents a connection protocol
type Protocol string

const (
	ProtocolQUIC      Protocol = "quic"
	ProtocolGRPC      Protocol = "grpc"
	ProtocolWebSocket Protocol = "websocket"
	ProtocolTCP       Protocol = "tcp"
)

// Option is a functional option for configuring the client
type Option func(*Config)

// WithToken sets the authentication token
func WithToken(token string) Option {
	return func(c *Config) {
		c.Token = token
	}
}

// WithRegion sets the preferred region
func WithRegion(region string) Option {
	return func(c *Config) {
		c.Region = region
	}
}

// WithTimeout sets the operation timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithLogLevel sets the log level
func WithLogLevel(level string) Option {
	return func(c *Config) {
		c.LogLevel = level
	}
}

// WithRetryPolicy sets the retry policy
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(c *Config) {
		c.RetryPolicy = policy
	}
}

// WithProtocols sets the protocol preferences
func WithProtocols(protocols ...Protocol) Option {
	return func(c *Config) {
		c.Protocols = protocols
	}
}

// WithInsecureSkipVerify disables TLS certificate verification (not recommended for production)
func WithInsecureSkipVerify(skip bool) Option {
	return func(c *Config) {
		c.InsecureSkipVerify = skip
	}
}

// defaultConfig returns a configuration with default values
func defaultConfig() *Config {
	return &Config{
		Token:    os.Getenv("CLOUDBRIDGE_TOKEN"),
		Region:   getEnvOrDefault("CLOUDBRIDGE_REGION", "eu-central"),
		Timeout:  30 * time.Second,
		LogLevel: getEnvOrDefault("CLOUDBRIDGE_LOG_LEVEL", "info"),
		RetryPolicy: RetryPolicy{
			MaxRetries:   3,
			InitialDelay: time.Second,
			MaxDelay:     time.Minute,
			Multiplier:   2.0,
		},
		Protocols: []Protocol{
			ProtocolQUIC,
			ProtocolGRPC,
			ProtocolWebSocket,
		},
		InsecureSkipVerify: false,
	}
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.Token == "" {
		return errors.New("token is required")
	}

	if c.Region == "" {
		return errors.New("region is required")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}

	if c.RetryPolicy.MaxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}

	if c.RetryPolicy.InitialDelay <= 0 {
		return errors.New("initial delay must be positive")
	}

	if c.RetryPolicy.MaxDelay <= 0 {
		return errors.New("max delay must be positive")
	}

	if c.RetryPolicy.Multiplier <= 0 {
		return errors.New("multiplier must be positive")
	}

	if len(c.Protocols) == 0 {
		return errors.New("at least one protocol must be specified")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.LogLevel] {
		return errors.New("invalid log level")
	}

	return nil
}

// getEnvOrDefault returns the value of the environment variable or the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
