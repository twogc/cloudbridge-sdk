package jwt

import (
	"testing"
)

func TestParseToken(t *testing.T) {
	// Valid JWT token (example, not a real token)
	// Header: {"alg":"RS256","typ":"JWT"}
	// Payload: {"sub":"user-123","tenant_id":"tenant-456","exp":1735689600,"iat":1735603200}
	validToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid format - 2 parts",
			token:   "part1.part2",
			wantErr: true,
		},
		{
			name:    "invalid format - 4 parts",
			token:   "part1.part2.part3.part4",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			token:   "invalid.invalid.invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims == nil {
				t.Error("ParseToken() returned nil claims for valid token")
			}
		})
	}
}

func TestExtractTenantID(t *testing.T) {
	// Valid token with tenant_id
	validToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid token with tenant_id",
			token:   validToken,
			want:    "tenant-456",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token format",
			token:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractTenantID(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractTenantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractTenantID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSubject(t *testing.T) {
	validToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyIsInRlbmFudF9pZCI6InRlbmFudC00NTYiLCJleHAiOjE3MzU2ODk2MDAsImlhdCI6MTczNTYwMzIwMH0.c2lnbmF0dXJl"

	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid token with subject",
			token:   validToken,
			want:    "user-123",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractSubject(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractSubject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExtractSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateBasicFormat(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid format",
			token:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEyMyJ9.c2lnbmF0dXJl",
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid format - 2 parts",
			token:   "part1.part2",
			wantErr: true,
		},
		{
			name:    "invalid format - 4 parts",
			token:   "part1.part2.part3.part4",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			token:   "!!!.!!!.!!!",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBasicFormat(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBasicFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
