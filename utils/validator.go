package utils

import (
	"fmt"
	"net"
	"net/mail"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateIP validates an IP address
func ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("IP address is required")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address format")
	}

	return nil
}

// ValidateRequired validates that a string is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateMinLength validates minimum string length
func ValidateMinLength(value, fieldName string, minLength int) error {
	if len(value) < minLength {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLength)
	}
	return nil
}

// ValidateMaxLength validates maximum string length
func ValidateMaxLength(value, fieldName string, maxLength int) error {
	if len(value) > maxLength {
		return fmt.Errorf("%s must not exceed %d characters", fieldName, maxLength)
	}
	return nil
}

// ValidateAPIKey validates an API key format
func ValidateAPIKey(apiKey string) error {
	if err := ValidateRequired(apiKey, "API key"); err != nil {
		return err
	}

	if len(apiKey) < 16 {
		return fmt.Errorf("API key must be at least 16 characters")
	}

	return nil
}

// ValidateEndpoint validates an endpoint format
func ValidateEndpoint(endpoint string) error {
	if err := ValidateRequired(endpoint, "endpoint"); err != nil {
		return err
	}

	if !strings.HasPrefix(endpoint, "/") {
		return fmt.Errorf("endpoint must start with '/'")
	}

	return nil
}

// SanitizeString removes dangerous characters from a string
func SanitizeString(input string) string {
	// Remove null bytes and trim spaces
	sanitized := strings.ReplaceAll(input, "\x00", "")
	return strings.TrimSpace(sanitized)
}
