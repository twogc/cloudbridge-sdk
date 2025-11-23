package cloudbridge

import "errors"

// Service represents a discovered service
type Service struct {
	ID       string
	Name     string
	Address  string
	Port     int
	Tags     []string
	Healthy  bool
	PeerID   string
	Metadata map[string]string
}

// ServiceConfig holds configuration for service registration
type ServiceConfig struct {
	Name string
	Port int
	Tags []string
}

// validate checks if the service configuration is valid
func (sc *ServiceConfig) validate() error {
	if sc.Name == "" {
		return errors.New("service name cannot be empty")
	}

	if sc.Port <= 0 || sc.Port > 65535 {
		return errors.New("invalid port")
	}

	return nil
}
