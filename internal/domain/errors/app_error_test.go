package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		reason   string
		cause    error
		expected string
	}{
		{
			name:     "without cause",
			reason:   "test error",
			cause:    nil,
			expected: "test error",
		},
		{
			name:     "with cause",
			reason:   "test error",
			cause:    errors.New("root cause"),
			expected: "test error (cause: root cause)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAppError(tt.reason, tt.cause)
			actual := err.Error()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestNewAppError(t *testing.T) {
	reason := "test reason"
	cause := errors.New("test cause")

	err := NewAppError(reason, cause)

	if err.reason != reason {
		t.Errorf("expected reason %q, got %q", reason, err.reason)
	}

	if !errors.Is(cause, err.cause) {
		t.Errorf("expected cause %v, got %v", cause, err.cause)
	}
}
