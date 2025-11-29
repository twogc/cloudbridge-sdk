package cloudbridge

import (
	"testing"
)

func TestServiceConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServiceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ServiceConfig{
				Name: "test-service",
				Port: 8080,
				Tags: []string{"api", "v1"},
			},
			wantErr: false,
		},
		{
			name: "valid config without tags",
			config: ServiceConfig{
				Name: "test-service",
				Port: 8080,
				Tags: nil,
			},
			wantErr: false,
		},
		{
			name: "empty service name",
			config: ServiceConfig{
				Name: "",
				Port: 8080,
			},
			wantErr: true,
		},
		{
			name: "invalid port - zero",
			config: ServiceConfig{
				Name: "test-service",
				Port: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid port - negative",
			config: ServiceConfig{
				Name: "test-service",
				Port: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid port - too high",
			config: ServiceConfig{
				Name: "test-service",
				Port: 70000,
			},
			wantErr: true,
		},
		{
			name: "valid port - 1",
			config: ServiceConfig{
				Name: "test-service",
				Port: 1,
			},
			wantErr: false,
		},
		{
			name: "valid port - 65535",
			config: ServiceConfig{
				Name: "test-service",
				Port: 65535,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceStruct(t *testing.T) {
	service := Service{
		ID:       "service-123",
		Name:     "test-service",
		Address:  "192.168.1.1",
		Port:     8080,
		Tags:     []string{"api", "v1"},
		Healthy:  true,
		PeerID:   "peer-123",
		Metadata: map[string]string{"region": "eu-central"},
	}

	if service.ID != "service-123" {
		t.Errorf("Service.ID = %v, want %v", service.ID, "service-123")
	}

	if service.Name != "test-service" {
		t.Errorf("Service.Name = %v, want %v", service.Name, "test-service")
	}

	if service.Port != 8080 {
		t.Errorf("Service.Port = %v, want %v", service.Port, 8080)
	}

	if !service.Healthy {
		t.Error("Service.Healthy = false, want true")
	}

	if len(service.Tags) != 2 {
		t.Errorf("Service.Tags length = %v, want %v", len(service.Tags), 2)
	}

	if service.Metadata["region"] != "eu-central" {
		t.Errorf("Service.Metadata[region] = %v, want %v", service.Metadata["region"], "eu-central")
	}
}

