package cloudbridge

import (
	"testing"
	"time"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: &Config{
				Token:    "test-token",
				Region:   "eu-central",
				Timeout:  30 * time.Second,
				LogLevel: "info",
				RetryPolicy: RetryPolicy{
					MaxRetries:   3,
					InitialDelay: time.Second,
					MaxDelay:     time.Minute,
					Multiplier:   2.0,
				},
				Protocols: []Protocol{ProtocolQUIC},
			},
			wantErr: false,
		},
		{
			name: "missing token",
			config: &Config{
				Region:  "eu-central",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "missing region",
			config: &Config{
				Token:   "test-token",
				Timeout: 30 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				Token:   "test-token",
				Region:  "eu-central",
				Timeout: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			config: &Config{
				Token:    "test-token",
				Region:   "eu-central",
				Timeout:  30 * time.Second,
				LogLevel: "invalid",
			},
			wantErr: true,
		},
		{
			name: "empty protocols",
			config: &Config{
				Token:     "test-token",
				Region:    "eu-central",
				Timeout:   30 * time.Second,
				LogLevel:  "info",
				Protocols: []Protocol{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.RetryPolicy == (RetryPolicy{}) {
				tt.config.RetryPolicy = RetryPolicy{
					MaxRetries:   3,
					InitialDelay: time.Second,
					MaxDelay:     time.Minute,
					Multiplier:   2.0,
				}
			}
			if tt.config.LogLevel == "" && !tt.wantErr {
				tt.config.LogLevel = "info"
			}
			if len(tt.config.Protocols) == 0 && !tt.wantErr {
				tt.config.Protocols = []Protocol{ProtocolQUIC}
			}

			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithOptions(t *testing.T) {
	config := defaultConfig()

	// Test WithToken
	WithToken("new-token")(config)
	if config.Token != "new-token" {
		t.Errorf("WithToken() did not set token correctly")
	}

	// Test WithRegion
	WithRegion("us-west")(config)
	if config.Region != "us-west" {
		t.Errorf("WithRegion() did not set region correctly")
	}

	// Test WithTimeout
	WithTimeout(60 * time.Second)(config)
	if config.Timeout != 60*time.Second {
		t.Errorf("WithTimeout() did not set timeout correctly")
	}

	// Test WithLogLevel
	WithLogLevel("debug")(config)
	if config.LogLevel != "debug" {
		t.Errorf("WithLogLevel() did not set log level correctly")
	}

	// Test WithProtocols
	WithProtocols(ProtocolGRPC, ProtocolWebSocket)(config)
	if len(config.Protocols) != 2 {
		t.Errorf("WithProtocols() did not set protocols correctly")
	}
	if config.Protocols[0] != ProtocolGRPC || config.Protocols[1] != ProtocolWebSocket {
		t.Errorf("WithProtocols() protocols in wrong order")
	}
}

func TestRetryPolicyValidation(t *testing.T) {
	tests := []struct {
		name    string
		policy  RetryPolicy
		wantErr bool
	}{
		{
			name: "valid policy",
			policy: RetryPolicy{
				MaxRetries:   3,
				InitialDelay: time.Second,
				MaxDelay:     time.Minute,
				Multiplier:   2.0,
			},
			wantErr: false,
		},
		{
			name: "negative max retries",
			policy: RetryPolicy{
				MaxRetries:   -1,
				InitialDelay: time.Second,
				MaxDelay:     time.Minute,
				Multiplier:   2.0,
			},
			wantErr: true,
		},
		{
			name: "negative initial delay",
			policy: RetryPolicy{
				MaxRetries:   3,
				InitialDelay: -1 * time.Second,
				MaxDelay:     time.Minute,
				Multiplier:   2.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				Token:       "test-token",
				Region:      "eu-central",
				Timeout:     30 * time.Second,
				LogLevel:    "info",
				RetryPolicy: tt.policy,
				Protocols:   []Protocol{ProtocolQUIC},
			}

			err := config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
