package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Security utility functions to prevent injection attacks

// SanitizeLogString removes or replaces characters that could be used for log injection attacks
// This includes newlines, carriage returns, and other control characters
func SanitizeLogString(input string) string {
	if !utf8.ValidString(input) {
		// Replace invalid UTF-8 with replacement character
		input = strings.ToValidUTF8(input, "�")
	}

	// Remove null bytes, newlines, carriage returns, tabs, and other control characters
	input = strings.ReplaceAll(input, "\x00", " ")
	input = strings.ReplaceAll(input, "\n", " ")
	input = strings.ReplaceAll(input, "\r", " ")
	input = strings.ReplaceAll(input, "\t", " ")
	
	// Remove other dangerous control characters (0x01-0x1F except space)
	var result strings.Builder
	for _, r := range input {
		if r >= 0x20 || r == ' ' {
			result.WriteRune(r)
		} else {
			result.WriteString(" ")
		}
	}

	// Truncate if too long (prevent log flooding)
	output := result.String()
	const maxLogLength = 1000
	if len(output) > maxLogLength {
		output = output[:maxLogLength] + "...[truncated]"
	}

	return output
}

// SanitizeUserInput provides general user input sanitization
func SanitizeUserInput(input string) string {
	if !utf8.ValidString(input) {
		input = strings.ToValidUTF8(input, "�")
	}

	// Remove null bytes and replace with spaces
	input = strings.ReplaceAll(input, "\x00", " ")
	
	// Trim whitespace
	return strings.TrimSpace(input)
}

// ValidateSlug ensures a slug contains only safe characters for MongoDB queries
// Prevents NoSQL injection by allowing only alphanumeric, dash, and underscore
func ValidateSlug(slug string) bool {
	if slug == "" || len(slug) > 100 {
		return false
	}
	
	// Only allow alphanumeric characters, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, slug)
	return matched
}

// ValidateServiceName validates service names to prevent injection attacks
func ValidateServiceName(name string) bool {
	if name == "" || len(name) > 255 {
		return false
	}
	
	// Allow alphanumeric, spaces, hyphens, underscores, and dots
	// But prevent patterns that could be problematic in logs or queries
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\s._-]+$`, name)
	if !matched {
		return false
	}
	
	// Prevent names that start or end with dots, spaces, or special chars
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 || strings.HasPrefix(trimmed, ".") || strings.HasSuffix(trimmed, ".") {
		return false
	}
	
	return true
}

// ValidateURL performs basic URL validation without being overly restrictive
func ValidateURL(url string) bool {
	if url == "" || len(url) > 2048 {
		return false
	}
	
	// Basic URL pattern - allow http/https URLs
	urlPattern := regexp.MustCompile(`^https?://[^\s<>"{}|\\^` + "`" + `]+$`)
	return urlPattern.MatchString(url)
}

// ValidateStatusValue ensures status values are from allowed set
func ValidateStatusValue(status string) bool {
	allowedStatuses := map[string]bool{
		"operational": true,
		"degraded":    true,
		"down":        true,
		"maintenance": true,
	}
	return allowedStatuses[status]
}

// SanitizeMap sanitizes all string values in a map
func SanitizeMap(input map[string]interface{}) map[string]interface{} {
	if input == nil {
		return nil
	}
	
	result := make(map[string]interface{})
	for k, v := range input {
		sanitizedKey := SanitizeLogString(k)
		if len(sanitizedKey) > 100 {
			continue // Skip keys that are too long
		}
		
		switch val := v.(type) {
		case string:
			result[sanitizedKey] = SanitizeLogString(val)
		case map[string]interface{}:
			result[sanitizedKey] = SanitizeMap(val)
		default:
			result[sanitizedKey] = val
		}
	}
	
	return result
}

// ValidateObjectID checks if a string could be a valid MongoDB ObjectID
func ValidateObjectID(id string) bool {
	if len(id) != 24 {
		return false
	}
	
	// Check if it contains only hexadecimal characters
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]{24}$`, id)
	return matched
}

// TruncateString safely truncates a string to a maximum length
func TruncateString(input string, maxLen int) string {
	if len(input) <= maxLen {
		return input
	}
	
	return input[:maxLen] + "..."
}