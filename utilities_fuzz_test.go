package apirouter

import (
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/mrz1836/go-parameters"
	"github.com/stretchr/testify/require"
)

// FuzzSnakeCase tests the SnakeCase function with various string inputs
// to ensure it handles Unicode, special characters, and edge cases properly
func FuzzSnakeCase(f *testing.F) {
	// Seed corpus with representative test cases
	testCases := []string{
		"testCamelCase",
		"TestCamelCase",
		"TEstCamelCase",
		"testCamelCASE",
		"testCamelAPI",
		"testCamelIP",
		"testCamelURL",
		"testCamelJSON",
		"APIKey",
		"JSONData",
		"URLPath",
		"IPAddress",
		"",
		"a",
		"A",
		"123",
		"test_snake_case",
		"ALLCAPS",
		"lowercase",
		"MixedCASE123",
		"test with spaces",
		"test-with-dashes",
		"test.with.dots",
		"test@with@symbols",
		"тест", // Cyrillic
		"テスト",  // Japanese
		"🔥test🔥",
		"test\n\t\r",
		"test\x00null",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// Ensure the function doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SnakeCase panicked with input %q: %v", input, r)
			}
		}()

		result := SnakeCase(input)

		// Basic validation: result should be valid UTF-8
		if !utf8.ValidString(result) {
			t.Errorf("SnakeCase(%q) produced invalid UTF-8: %q", input, result)
		}

		// Result should not be longer than reasonable bounds
		// (accounting for potential underscore insertions)
		maxExpectedLength := len(input) * 2
		if len(result) > maxExpectedLength {
			t.Errorf("SnakeCase(%q) produced unexpectedly long result: %q (len=%d)", input, result, len(result))
		}

		// If input is empty, result should be empty
		if input == "" && result != "" {
			t.Errorf("SnakeCase(%q) should return empty string, got %q", input, result)
		}

		// Result should be lowercase (excluding underscores)
		for _, r := range result {
			if r != '_' && r >= 'A' && r <= 'Z' {
				t.Errorf("SnakeCase(%q) should return lowercase result, got %q", input, result)
				break
			}
		}
	})
}

// FuzzGetClientIPAddress tests IP address extraction from HTTP requests
// with various header combinations and malformed inputs
func FuzzGetClientIPAddress(f *testing.F) {
	// Seed corpus with representative IP addresses and headers
	testCases := []struct {
		xForwardedFor string
		remoteAddr    string
	}{
		{"192.168.1.1", "10.0.0.1:8080"},
		{"203.0.113.1,198.51.100.1", "192.168.1.1:8080"},
		{"2001:db8::1", "[2001:db8::1]:8080"},
		{"", "192.168.1.1:8080"},
		{"", "[::1]:8080"},
		{"invalid-ip", "192.168.1.1:8080"},
		{"", "invalid:addr"},
		{"", ":8080"},
		{"", ""},
		{"127.0.0.1", ""},
		{"10.0.0.1,", "192.168.1.1:8080"},
		{",10.0.0.1", "192.168.1.1:8080"},
		{"10.0.0.1,,192.168.1.1", "203.0.113.1:8080"},
		{"999.999.999.999", "192.168.1.1:8080"},
		{"192.168.1", "10.0.0.1:8080"},
		{"192.168.1.1.1", "10.0.0.1:8080"},
		{"192.168.1.1:8080", "10.0.0.1:8080"},
		{"[192.168.1.1]", "10.0.0.1:8080"},
		{"::ffff:192.168.1.1", "[::1]:8080"},
	}

	for _, tc := range testCases {
		f.Add(tc.xForwardedFor, tc.remoteAddr)
	}

	f.Fuzz(func(t *testing.T, xForwardedFor, remoteAddr string) {
		// Create a test request
		req := httptest.NewRequest("GET", "/test", nil)

		// Set headers based on fuzz inputs
		if xForwardedFor != "" {
			req.Header.Set("X-Forwarded-For", xForwardedFor)
		}
		req.RemoteAddr = remoteAddr

		// Ensure the function doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GetClientIPAddress panicked with X-Forwarded-For=%q, RemoteAddr=%q: %v", xForwardedFor, remoteAddr, r)
			}
		}()

		result := GetClientIPAddress(req)

		// Result should be valid UTF-8
		if !utf8.ValidString(result) {
			t.Errorf("GetClientIPAddress produced invalid UTF-8: %q", result)
		}

		// Result should either be empty or a reasonable IP address format
		if result != "" {
			// Basic sanity check: should not contain obviously invalid characters for IP
			if strings.ContainsAny(result, "\n\r\t\x00") {
				t.Errorf("GetClientIPAddress produced result with control characters: %q", result)
			}

			// Should not be excessively long
			if len(result) > 45 { // Max IPv6 length is about 45 chars
				t.Errorf("GetClientIPAddress produced unexpectedly long result: %q (len=%d)", result, len(result))
			}
		}
	})
}

// FuzzFilterMap tests parameter filtering with various input combinations
// to ensure sensitive data is properly masked
func FuzzFilterMap(f *testing.F) {
	// Seed corpus with representative parameter names and values
	testCases := []struct {
		paramKey   string
		paramValue string
		filterKey  string
	}{
		{"password", "secret123", "password"},
		{"apiKey", "key123", "apiKey"},
		{"token", "jwt123", "token"},
		{"email", "test@example.com", "password"},
		{"", "value", ""},
		{"key", "", "key"},
		{"normalParam", "normalValue", "sensitiveParam"},
		{"user_id", "12345", "password"},
		{"api-key", "secret", "api-key"},
		{"API_KEY", "secret", "API_KEY"},
		{"test key with spaces", "value", "test key with spaces"},
		{"тест", "значение", "тест"},
		{"🔑key", "🔒secret", "🔑key"},
		{"key\nwith\nnewlines", "value\nwith\nnewlines", "key\nwith\nnewlines"},
		{"key\x00null", "value\x00null", "key\x00null"},
	}

	for _, tc := range testCases {
		f.Add(tc.paramKey, tc.paramValue, tc.filterKey)
	}

	f.Fuzz(func(t *testing.T, paramKey, paramValue, filterKey string) {
		// Create parameters
		params := &parameters.Params{
			Values: make(map[string]interface{}),
		}
		params.Values[paramKey] = paramValue

		filterOutFields := []string{filterKey}

		// Ensure the function doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("FilterMap panicked with paramKey=%q, paramValue=%q, filterKey=%q: %v", paramKey, paramValue, filterKey, r)
			}
		}()

		result := FilterMap(params, filterOutFields)

		// Result should not be nil
		require.NotNil(t, result)
		require.NotNil(t, result.Values)

		// Result should have the same number of parameters
		if len(result.Values) != len(params.Values) {
			t.Errorf("FilterMap changed parameter count: expected %d, got %d", len(params.Values), len(result.Values))
		}

		// Check if filtering worked correctly
		if val, exists := result.Values[paramKey]; exists {
			if paramKey == filterKey {
				// Should be filtered (replaced with PROTECTED)
				if filtered, ok := val.([]string); ok {
					if len(filtered) != 1 || filtered[0] != "PROTECTED" {
						t.Errorf("Expected filtered value to be [PROTECTED], got %v", val)
					}
				} else {
					t.Errorf("Expected filtered value to be []string, got %T", val)
				}
			} else {
				// Should not be filtered (should be original value)
				if val != paramValue {
					t.Errorf("Expected unfiltered value %v, got %v", paramValue, val)
				}
			}
		}
	})
}
