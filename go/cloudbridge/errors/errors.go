package errors

import (
	"errors"
	"fmt"
)

// Error types
var (
	ErrAuth         = errors.New("authentication error")
	ErrNetwork      = errors.New("network error")
	ErrPeerNotFound = errors.New("peer not found")
	ErrTimeout      = errors.New("operation timeout")
	ErrClosed       = errors.New("connection closed")
	ErrInvalidInput = errors.New("invalid input")
)

// AuthError represents an authentication error
type AuthError struct {
	Message string
	Err     error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("authentication error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NetworkError represents a network error
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("network error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// PeerNotFoundError represents a peer not found error
type PeerNotFoundError struct {
	PeerID string
}

func (e *PeerNotFoundError) Error() string {
	return fmt.Sprintf("peer not found: %s", e.PeerID)
}

// TimeoutError represents a timeout error
type TimeoutError struct {
	Operation string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("operation timeout: %s", e.Operation)
}

// IsAuthError checks if an error is an authentication error
func IsAuthError(err error) bool {
	var authErr *AuthError
	return errors.As(err, &authErr) || errors.Is(err, ErrAuth)
}

// IsNetworkError checks if an error is a network error
func IsNetworkError(err error) bool {
	var netErr *NetworkError
	return errors.As(err, &netErr) || errors.Is(err, ErrNetwork)
}

// IsPeerNotFoundError checks if an error is a peer not found error
func IsPeerNotFoundError(err error) bool {
	var peerErr *PeerNotFoundError
	return errors.As(err, &peerErr) || errors.Is(err, ErrPeerNotFound)
}

// IsTimeoutError checks if an error is a timeout error
func IsTimeoutError(err error) bool {
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr) || errors.Is(err, ErrTimeout)
}

// NewAuthError creates a new authentication error
func NewAuthError(message string, err error) error {
	return &AuthError{
		Message: message,
		Err:     err,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(message string, err error) error {
	return &NetworkError{
		Message: message,
		Err:     err,
	}
}

// NewPeerNotFoundError creates a new peer not found error
func NewPeerNotFoundError(peerID string) error {
	return &PeerNotFoundError{
		PeerID: peerID,
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(operation string) error {
	return &TimeoutError{
		Operation: operation,
	}
}
