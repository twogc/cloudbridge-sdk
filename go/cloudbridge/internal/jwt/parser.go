package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Claims represents JWT claims
type Claims struct {
	Sub        string `json:"sub"`
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name,omitempty"`
	OrgID      string `json:"organization_id,omitempty"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
	Iss        string `json:"iss"`
	Aud        string `json:"aud"`
}

// ParseToken parses JWT token without verification
// Note: This is for extracting claims only, not for security validation
func ParseToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format: expected 3 parts, got %d", len(parts))
	}

	// Decode payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	return &claims, nil
}

// ExtractTenantID extracts tenant ID from JWT token
func ExtractTenantID(token string) (string, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}

	if claims.TenantID == "" {
		return "", fmt.Errorf("token does not contain tenant_id claim")
	}

	return claims.TenantID, nil
}

// ExtractSubject extracts subject (user ID) from JWT token
func ExtractSubject(token string) (string, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}

	if claims.Sub == "" {
		return "", fmt.Errorf("token does not contain sub claim")
	}

	return claims.Sub, nil
}

// ValidateBasicFormat validates basic JWT format without signature verification
func ValidateBasicFormat(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid token format: expected 3 parts separated by dots")
	}

	// Validate each part is base64 encoded
	for i, part := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			return fmt.Errorf("invalid base64 encoding in part %d: %w", i+1, err)
		}
	}

	return nil
}
