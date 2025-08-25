package utils

import (
	"strings"
	"testing"
)

func TestSanitizeLogString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean string",
			input:    "normal log message",
			expected: "normal log message",
		},
		{
			name:     "with newlines",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
		},
		{
			name:     "with tabs",
			input:    "col1\tcol2\tcol3",
			expected: "col1 col2 col3",
		},
		{
			name:     "with null bytes",
			input:    "data\x00more data",
			expected: "data more data",
		},
		{
			name:     "control characters",
			input:    "start\x01\x02\x1fend",
			expected: "start   end",
		},
		{
			name:     "long string truncation",
			input:    strings.Repeat("a", 1100),
			expected: strings.Repeat("a", 1000) + "...[truncated]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeLogString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeLogString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid slug",
			input: "my-service-1",
			valid: true,
		},
		{
			name:  "with underscores",
			input: "my_service_1",
			valid: true,
		},
		{
			name:  "alphanumeric only",
			input: "service123",
			valid: true,
		},
		{
			name:  "empty slug",
			input: "",
			valid: false,
		},
		{
			name:  "with spaces",
			input: "my service",
			valid: false,
		},
		{
			name:  "with special chars",
			input: "service@test",
			valid: false,
		},
		{
			name:  "with dots",
			input: "service.test",
			valid: false,
		},
		{
			name:  "too long",
			input: strings.Repeat("a", 101),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSlug(tt.input)
			if result != tt.valid {
				t.Errorf("ValidateSlug(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestValidateServiceName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid name",
			input: "My Service 1",
			valid: true,
		},
		{
			name:  "with hyphens and underscores",
			input: "My-Service_1",
			valid: true,
		},
		{
			name:  "with dots",
			input: "My.Service.API",
			valid: true,
		},
		{
			name:  "empty name",
			input: "",
			valid: false,
		},
		{
			name:  "only spaces",
			input: "   ",
			valid: false,
		},
		{
			name:  "starts with dot",
			input: ".hidden-service",
			valid: false,
		},
		{
			name:  "ends with dot",
			input: "service.",
			valid: false,
		},
		{
			name:  "with special chars",
			input: "service@test",
			valid: false,
		},
		{
			name:  "too long",
			input: strings.Repeat("a", 256),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateServiceName(tt.input)
			if result != tt.valid {
				t.Errorf("ValidateServiceName(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid http URL",
			input: "http://example.com",
			valid: true,
		},
		{
			name:  "valid https URL",
			input: "https://api.example.com/health",
			valid: true,
		},
		{
			name:  "with port",
			input: "https://api.example.com:8080/health",
			valid: true,
		},
		{
			name:  "with query params",
			input: "https://api.example.com/health?check=true",
			valid: true,
		},
		{
			name:  "empty URL",
			input: "",
			valid: false,
		},
		{
			name:  "no protocol",
			input: "example.com",
			valid: false,
		},
		{
			name:  "ftp protocol",
			input: "ftp://example.com",
			valid: false,
		},
		{
			name:  "with dangerous chars",
			input: "https://example.com/<script>",
			valid: false,
		},
		{
			name:  "too long",
			input: "https://" + strings.Repeat("a", 2050),
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateURL(tt.input)
			if result != tt.valid {
				t.Errorf("ValidateURL(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestValidateStatusValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "operational",
			input: "operational",
			valid: true,
		},
		{
			name:  "degraded",
			input: "degraded",
			valid: true,
		},
		{
			name:  "down",
			input: "down",
			valid: true,
		},
		{
			name:  "maintenance",
			input: "maintenance",
			valid: true,
		},
		{
			name:  "invalid status",
			input: "unknown",
			valid: false,
		},
		{
			name:  "empty status",
			input: "",
			valid: false,
		},
		{
			name:  "mixed case",
			input: "Operational",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateStatusValue(tt.input)
			if result != tt.valid {
				t.Errorf("ValidateStatusValue(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestValidateObjectID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid ObjectID",
			input: "507f1f77bcf86cd799439011",
			valid: true,
		},
		{
			name:  "valid ObjectID uppercase",
			input: "507F1F77BCF86CD799439011",
			valid: true,
		},
		{
			name:  "too short",
			input: "507f1f77bcf86cd79943901",
			valid: false,
		},
		{
			name:  "too long",
			input: "507f1f77bcf86cd7994390111",
			valid: false,
		},
		{
			name:  "invalid characters",
			input: "507f1f77bcf86cd79943901g",
			valid: false,
		},
		{
			name:  "empty string",
			input: "",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateObjectID(tt.input)
			if result != tt.valid {
				t.Errorf("ValidateObjectID(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func TestSanitizeMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "nil map",
			input:    nil,
			expected: nil,
		},
		{
			name: "clean map",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		},
		{
			name: "map with dirty strings",
			input: map[string]interface{}{
				"key\n1": "value\r\n1",
				"key2":   "value2",
			},
			expected: map[string]interface{}{
				"key 1": "value  1",
				"key2":  "value2",
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner\n": "value\r",
				},
			},
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner ": "value ",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeMap(tt.input)
			if !mapsEqual(result, tt.expected) {
				t.Errorf("SanitizeMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare maps for testing
func mapsEqual(a, b map[string]interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || !valuesEqual(v, bv) {
			return false
		}
	}
	return true
}

func valuesEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		if bv, ok := b.(string); ok {
			return av == bv
		}
	case int:
		if bv, ok := b.(int); ok {
			return av == bv
		}
	case map[string]interface{}:
		if bv, ok := b.(map[string]interface{}); ok {
			return mapsEqual(av, bv)
		}
	default:
		return a == b
	}
	return false
}