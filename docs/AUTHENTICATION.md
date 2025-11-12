# CloudBridge SDK Authentication Guide

Version: 0.1.0
Last Updated: November 2025

## Overview

CloudBridge SDK uses JWT-based authentication with OIDC (OpenID Connect) integration through Zitadel. This guide covers authentication setup, token management, and security best practices.

## Authentication Flow

```
Application
    |
    | 1. Initialize SDK with token
    v
SDK Client
    |
    | 2. Validate token format
    v
Token Validation
    |
    | 3. Connect to CloudBridge
    v
CloudBridge Network
    |
    | 4. Verify JWT with Zitadel
    v
Zitadel OIDC
    |
    | 5. Validate signature & claims
    v
Authorization Decision
    |
    | 6. Grant/Deny access
    v
Connection Established
```

## Getting an API Token

### Via Dashboard

1. Log in to CloudBridge Dashboard
2. Navigate to Settings > API Tokens
3. Click "Generate New Token"
4. Set token name and permissions
5. Copy token (shown only once)
6. Store securely

### Via CLI

```bash
cloudbridge auth login
cloudbridge auth token create --name "my-app" --scope "connect,tunnel"
```

### Token Format

CloudBridge tokens are JWT tokens with the following structure:

```
Header:
{
  "alg": "RS256",
  "typ": "JWT",
  "kid": "key-id"
}

Payload:
{
  "iss": "https://auth.cloudbridge.global",
  "sub": "user-id",
  "aud": "cloudbridge-api",
  "exp": 1735689600,
  "iat": 1735603200,
  "tenant_id": "tenant-123",
  "scopes": ["connect", "tunnel", "mesh"]
}

Signature:
RS256(base64(header) + "." + base64(payload), private_key)
```

## Using Tokens in SDK

### Method 1: Direct Token

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."),
)
```

### Method 2: Environment Variable

```bash
export CLOUDBRIDGE_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

```go
// Token automatically loaded from environment
client, err := cloudbridge.NewClient()
```

### Method 3: Configuration File

```yaml
# config.yaml
cloudbridge:
  token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
  region: eu-central
```

```go
// Load from file
config := loadConfig("config.yaml")
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(config.Token),
)
```

### Method 4: External Secret Manager

```go
// AWS Secrets Manager
token, err := getSecretValue("cloudbridge/api-token")
if err != nil {
    return err
}

client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
)
```

## Token Scopes

CloudBridge tokens support granular permissions:

### Available Scopes

- `connect` - Establish P2P connections
- `tunnel` - Create secure tunnels
- `mesh` - Join mesh networks
- `discover` - Service discovery
- `admin` - Administrative operations

### Scope Examples

**Read-only connection:**
```
Scopes: ["connect"]
```

**Full access:**
```
Scopes: ["connect", "tunnel", "mesh", "discover"]
```

**Service account:**
```
Scopes: ["connect", "discover"]
```

## Token Lifecycle

### Token Expiration

Tokens have configurable expiration:

- Short-lived: 1 hour (interactive use)
- Medium-lived: 24 hours (applications)
- Long-lived: 90 days (service accounts)

### Checking Token Expiration

```go
import "github.com/golang-jwt/jwt/v5"

func isTokenExpired(tokenString string) (bool, error) {
    token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
    if err != nil {
        return false, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return false, errors.New("invalid claims")
    }

    exp, ok := claims["exp"].(float64)
    if !ok {
        return false, errors.New("no expiration")
    }

    return time.Now().Unix() > int64(exp), nil
}
```

### Token Refresh

```go
// Automatic refresh before expiration
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(initialToken),
    cloudbridge.WithTokenRefresh(func() (string, error) {
        // Fetch new token from your auth system
        return fetchNewToken()
    }),
)
```

## OIDC Integration

### Using OIDC Flow

For interactive applications, use OIDC flow:

```go
import "github.com/zitadel/oidc/v3/pkg/client/rp"

// Configure OIDC client
config := &oauth2.Config{
    ClientID:     "cloudbridge-app",
    ClientSecret: "secret",
    Endpoint: oauth2.Endpoint{
        AuthURL:  "https://auth.cloudbridge.global/oauth/authorize",
        TokenURL: "https://auth.cloudbridge.global/oauth/token",
    },
    RedirectURL: "http://localhost:8080/callback",
    Scopes:      []string{"openid", "profile", "cloudbridge"},
}

// Get authorization URL
authURL := config.AuthCodeURL("state")

// Exchange code for token
token, err := config.Exchange(ctx, code)

// Use token with SDK
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token.AccessToken),
)
```

### Device Code Flow

For CLI applications:

```go
// Request device code
resp, err := http.PostForm(
    "https://auth.cloudbridge.global/oauth/device",
    url.Values{"client_id": {"cloudbridge-cli"}},
)

// Show user code
fmt.Printf("Visit: %s\n", deviceCode.VerificationURI)
fmt.Printf("Enter code: %s\n", deviceCode.UserCode)

// Poll for token
token, err := pollForToken(deviceCode.DeviceCode)

// Use token
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
)
```

## Multi-Tenancy

CloudBridge supports multi-tenant isolation via JWT claims.

### Tenant Identification

Token includes tenant information:

```json
{
  "tenant_id": "tenant-123",
  "tenant_name": "Acme Corp",
  "organization_id": "org-456"
}
```

### Tenant Isolation

SDK automatically enforces tenant isolation:

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token), // Contains tenant_id
)

// All operations scoped to tenant
conn, err := client.Connect(ctx, "peer-in-same-tenant")
```

### Cross-Tenant Communication

Requires explicit permission:

```json
{
  "tenant_id": "tenant-123",
  "allowed_tenants": ["tenant-456", "tenant-789"]
}
```

## Security Best Practices

### Do

1. Store tokens in secure locations
   - OS keyring (macOS Keychain, Windows Credential Manager)
   - Environment variables (ephemeral)
   - Secret management systems (Vault, AWS Secrets Manager)

2. Use short-lived tokens for interactive sessions

3. Rotate tokens regularly

4. Monitor token usage and revoke suspicious tokens

5. Use minimal required scopes

6. Enable TLS certificate verification

### Do Not

1. Hardcode tokens in source code
2. Commit tokens to version control
3. Share tokens between applications
4. Use long-lived tokens unnecessarily
5. Disable TLS verification in production
6. Log tokens in application logs

### Token Storage Examples

**macOS Keychain:**
```go
import "github.com/zalando/go-keyring"

// Store token
err := keyring.Set("cloudbridge", "api-token", token)

// Retrieve token
token, err := keyring.Get("cloudbridge", "api-token")
```

**Environment Variable (Docker):**
```dockerfile
# Use secrets, not ENV in production
RUN --mount=type=secret,id=cloudbridge_token \
    CLOUDBRIDGE_TOKEN=$(cat /run/secrets/cloudbridge_token) \
    ./app
```

**Kubernetes Secret:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloudbridge-token
type: Opaque
data:
  token: ZXlKaGJHY2lPaUpTVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5...
```

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    env:
    - name: CLOUDBRIDGE_TOKEN
      valueFrom:
        secretKeyRef:
          name: cloudbridge-token
          key: token
```

## Error Handling

### Authentication Errors

```go
import "github.com/twogc/cloudbridge-sdk/go/cloudbridge/errors"

client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
)

conn, err := client.Connect(ctx, peerID)
if err != nil {
    if errors.IsAuthError(err) {
        // Token expired or invalid
        newToken, err := refreshToken()
        if err != nil {
            return err
        }

        // Retry with new token
        client, err = cloudbridge.NewClient(
            cloudbridge.WithToken(newToken),
        )
    }
}
```

### Common Error Scenarios

**Token Expired:**
```
Error: authentication error: token expired
Solution: Refresh token and retry
```

**Invalid Token Format:**
```
Error: invalid token format
Solution: Check token is complete and properly formatted
```

**Insufficient Permissions:**
```
Error: insufficient permissions: missing scope 'tunnel'
Solution: Request token with required scopes
```

**Invalid Signature:**
```
Error: token signature verification failed
Solution: Token may be tampered, request new token
```

## Compliance

### GDPR Compliance

- Tokens contain minimal personal information
- Token revocation available
- Audit logs for token usage
- Data deletion on request

### SOC 2 Compliance

- Token encryption in transit (TLS 1.3)
- Token expiration enforcement
- Access control via scopes
- Audit trail of authentication events

## Testing

### Mock Authentication

For testing without real tokens:

```go
// Test mode (not for production)
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken("test-token"),
    cloudbridge.WithInsecureSkipVerify(true),
)
```

### Token Generation for Tests

```go
func generateTestToken() string {
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
        "sub": "test-user",
        "tenant_id": "test-tenant",
        "scopes": []string{"connect", "tunnel"},
        "exp": time.Now().Add(time.Hour).Unix(),
    })

    tokenString, _ := token.SignedString(testPrivateKey)
    return tokenString
}
```

## Troubleshooting

### Debug Authentication

Enable debug logging:

```go
client, err := cloudbridge.NewClient(
    cloudbridge.WithToken(token),
    cloudbridge.WithLogLevel("debug"),
)
```

### Verify Token

```bash
# Decode JWT (without verifying signature)
echo "eyJhbGc..." | base64 -d | jq .
```

### Test Token Validity

```go
health, err := client.Health(ctx)
if err != nil {
    if errors.IsAuthError(err) {
        log.Println("Token is invalid or expired")
    }
}
```

## Additional Resources

- [Zitadel Documentation](https://zitadel.com/docs)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OAuth 2.0 Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)
