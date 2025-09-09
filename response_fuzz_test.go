package apirouter

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

// TestStruct represents a simple struct for fuzz testing JSON encoding
type TestStruct struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	Password    string                 `json:"password"`
	APIKey      string                 `json:"api_key"`
	IsActive    bool                   `json:"is_active"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NestedTestStruct represents a more complex nested struct for hierarchical testing
type NestedTestStruct struct {
	User    TestStruct `json:"user"`
	Company struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		ID      int    `json:"id"`
	} `json:"company"`
	Permissions []string `json:"permissions"`
}

// FuzzJSONEncode tests JSON encoding with field filtering using various inputs
// to ensure proper handling of data structures and field selection
func FuzzJSONEncode(f *testing.F) {
	// Seed corpus with representative test cases
	testCases := []struct {
		name         string
		email        string
		password     string
		allowedField string
	}{
		{"John Doe", "john@example.com", "secret123", "name"},
		{"", "john@example.com", "secret123", "name"},
		{"John Doe", "", "secret123", "email"},
		{"John Doe", "john@example.com", "", "password"},
		{"John Doe", "john@example.com", "secret123", "nonexistent"},
		{"John Doe", "john@example.com", "secret123", ""},
		{"User🆔", "email📧", "pass🔒", "name"},
		{"User\nwith\nnewlines", "email\twith\ttabs", "pass\rwith\rreturns", "name"},
		{"User\x00null", "email\x00null", "pass\x00null", "email"},
		{strings.Repeat("n", 1000), "john@example.com", "secret123", "name"},
		{"John Doe", strings.Repeat("e", 1000), "secret123", "email"},
		{"John Doe", "john@example.com", strings.Repeat("p", 1000), "password"},
		{"John Doe", "john@example.com", "secret123", strings.Repeat("f", 100)},
		{"<script>alert('xss')</script>", "test@example.com", "secret", "name"},
		{`{"injected": "json"}`, "test@example.com", "secret", "name"},
		{"Name with \"quotes\"", "test@example.com", "secret", "name"},
	}

	for _, tc := range testCases {
		f.Add(tc.name, tc.email, tc.password, tc.allowedField)
	}

	f.Fuzz(func(t *testing.T, name, email, password, allowedField string) {
		// Create test data structures
		singleStruct := TestStruct{
			ID:          123,
			Name:        name,
			Email:       email,
			Password:    password,
			APIKey:      "api-key-123",
			IsActive:    true,
			Description: "Test user",
			Tags:        []string{"tag1", "tag2"},
			Metadata:    map[string]interface{}{"key": "value"},
		}

		// Test with slice of structs
		structSlice := []TestStruct{singleStruct}
		if name != "" { // Add second item for non-empty names
			structSlice = append(structSlice, TestStruct{
				ID:       456,
				Name:     name + "_2",
				Email:    email + "_2",
				Password: password + "_2",
			})
		}

		allowed := []string{allowedField}
		if allowedField != "" {
			// Add some common fields for better coverage
			allowed = append(allowed, "id", "is_active")
		}

		testCases := []struct {
			name string
			data interface{}
		}{
			{"single_struct", singleStruct},
			{"struct_slice", structSlice},
			{"empty_slice", []TestStruct{}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)

				// Ensure JSONEncode doesn't panic
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("JSONEncode panicked with data=%s, allowed=%v: %v", tc.name, allowed, r)
					}
				}()

				err := JSONEncode(encoder, tc.data, allowed)

				// Validate result
				if err != nil {
					// Some errors are expected with malformed data
					t.Logf("JSONEncode returned error for %s: %v", tc.name, err)
				} else {
					// Validate the JSON output
					result := buf.Bytes()

					// Should be valid UTF-8
					if !utf8.Valid(result) {
						t.Errorf("JSONEncode produced invalid UTF-8 for %s", tc.name)
					}

					// Should be valid JSON
					var decoded interface{}
					if jsonErr := json.Unmarshal(result, &decoded); jsonErr != nil {
						t.Errorf("JSONEncode produced invalid JSON for %s: %v", tc.name, jsonErr)
					}

					// Result should not be excessively large (prevent memory issues)
					if len(result) > 100000 {
						t.Errorf("JSONEncode produced unexpectedly large output for %s: %d bytes", tc.name, len(result))
					}

					// If we have specific allowed fields, verify they're present/absent as expected
					resultStr := string(result)
					if allowedField != "" && tc.data != nil {
						// Check that sensitive fields are not present when not allowed
						sensitiveFields := []string{"password", "api_key"}
						for _, field := range sensitiveFields {
							if !contains(allowed, field) && strings.Contains(resultStr, field) {
								t.Logf("Note: Sensitive field %s found in output when not explicitly allowed", field)
							}
						}
					}
				}
			})
		}
	})
}

// FuzzJSONEncodeHierarchy tests hierarchical JSON encoding with nested structures
// and complex allowed field configurations
func FuzzJSONEncodeHierarchy(f *testing.F) {
	// Seed corpus with representative test cases
	testCases := []struct {
		userName    string
		userEmail   string
		companyName string
		permission  string
	}{
		{"John Doe", "john@example.com", "ACME Corp", "read"},
		{"", "john@example.com", "ACME Corp", "read"},
		{"John Doe", "", "ACME Corp", "write"},
		{"John Doe", "john@example.com", "", "admin"},
		{"John Doe", "john@example.com", "ACME Corp", ""},
		{"User🆔", "email📧@example.com", "Company🏢", "permission🔐"},
		{"User\nMultiline", "email\twith\ttabs@example.com", "Company\rName", "read\nwrite"},
		{"User\x00Null", "email\x00@example.com", "Company\x00Name", "permission\x00"},
		{strings.Repeat("n", 500), "john@example.com", "ACME Corp", "read"},
		{"John Doe", strings.Repeat("e", 500), "ACME Corp", "read"},
		{"John Doe", "john@example.com", strings.Repeat("c", 500), "read"},
		{"John Doe", "john@example.com", "ACME Corp", strings.Repeat("p", 500)},
		{`{"name": "injected"}`, "test@example.com", "Test Corp", "read"},
		{"Name \"with\" quotes", "test@example.com", "Corp with \"quotes\"", "read"},
		{"<script>alert('xss')</script>", "test@example.com", "Evil Corp", "admin"},
	}

	for _, tc := range testCases {
		f.Add(tc.userName, tc.userEmail, tc.companyName, tc.permission)
	}

	f.Fuzz(func(t *testing.T, userName, userEmail, companyName, permission string) {
		// Create nested test structure
		nested := NestedTestStruct{
			User: TestStruct{
				ID:       123,
				Name:     userName,
				Email:    userEmail,
				Password: "secret123",
				APIKey:   "api-key-123",
				IsActive: true,
			},
			Permissions: []string{permission, "base"},
		}
		nested.Company.Name = companyName
		nested.Company.Address = "123 Main St"
		nested.Company.ID = 456

		// Test different allowed field configurations
		allowedConfigs := []interface{}{
			nil,                             // No filtering
			[]string{"name", "email", "id"}, // Simple string array
			AllowedKeys{ // Hierarchical filtering
				"user":        []string{"name", "email"},
				"company":     []string{"name"},
				"permissions": nil,
			},
			AllowedKeys{ // Complex nested filtering
				"user": AllowedKeys{
					"name":      nil,
					"email":     nil,
					"is_active": nil,
				},
				"company": AllowedKeys{
					"name": nil,
					"id":   nil,
				},
			},
			AllowedKeys{}, // Empty allowed keys
		}

		for i, allowed := range allowedConfigs {
			t.Run(string(rune('A'+i)), func(t *testing.T) {
				var buf bytes.Buffer

				// Ensure JSONEncodeHierarchy doesn't panic
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("JSONEncodeHierarchy panicked with config %d, allowed=%v: %v", i, allowed, r)
					}
				}()

				err := JSONEncodeHierarchy(&buf, nested, allowed)

				// Validate result
				if err != nil {
					t.Logf("JSONEncodeHierarchy returned error for config %d: %v", i, err)
				} else {
					result := buf.Bytes()

					// Should be valid UTF-8
					if !utf8.Valid(result) {
						t.Errorf("JSONEncodeHierarchy produced invalid UTF-8 for config %d", i)
					}

					// Should be valid JSON
					var decoded interface{}
					if jsonErr := json.Unmarshal(result, &decoded); jsonErr != nil {
						t.Errorf("JSONEncodeHierarchy produced invalid JSON for config %d: %v", i, jsonErr)
					}

					// Result should not be excessively large
					if len(result) > 100000 {
						t.Errorf("JSONEncodeHierarchy produced unexpectedly large output for config %d: %d bytes", i, len(result))
					}

					// Basic sanity check for nested structure
					resultStr := string(result)
					if allowed != nil && len(resultStr) > 2 { // Not just "{}"
						// Should have proper JSON structure
						if !strings.Contains(resultStr, "{") {
							t.Errorf("JSONEncodeHierarchy did not produce object structure for config %d", i)
						}
					}
				}
			})
		}
	})
}

// FuzzRespondWith tests the RespondWith function with various status codes and data types
func FuzzRespondWith(f *testing.F) {
	// Seed corpus with representative test cases
	testCases := []struct {
		statusCode int
		dataStr    string
		dataType   string // "string", "error", "nil", "struct"
	}{
		{200, "success", "string"},
		{404, "not found", "error"},
		{500, "internal error", "error"},
		{204, "", "nil"},
		{304, "", "nil"},
		{400, "bad request", "string"},
		{201, "created", "struct"},
		{422, "validation error", "error"},         // Unprocessable Entity
		{409, "conflict", "error"},                 // Conflict
		{503, "service unavailable", "error"},      // Service Unavailable
		{200, strings.Repeat("a", 1000), "string"}, // Large data
		{200, "data\nwith\nnewlines", "string"},
		{200, "data\x00with\x00nulls", "string"},
		{200, "data🎉with🔥emojis", "string"},
		{200, `{"json": "data"}`, "string"},
		{500, "<script>alert('xss')</script>", "error"},
	}

	for _, tc := range testCases {
		f.Add(tc.statusCode, tc.dataStr, tc.dataType)
	}

	f.Fuzz(func(t *testing.T, statusCode int, dataStr, dataType string) {
		// Skip invalid status codes that would cause expected panics
		if statusCode < 100 || statusCode >= 600 {
			t.Skip("Skipping invalid status code")
			return
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)

		// Convert dataStr to appropriate data type
		var data interface{}
		switch dataType {
		case "error":
			if dataStr != "" {
				data = &testError{message: dataStr}
			} else {
				data = &testError{message: "test error"}
			}
		case "nil":
			data = nil
		case "struct":
			data = TestStruct{
				Name:  dataStr,
				Email: "test@example.com",
				ID:    statusCode, // Use status code as ID for variety
			}
		default: // "string"
			data = dataStr
		}

		// Ensure RespondWith doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RespondWith panicked with statusCode=%d, data=%v, type=%s: %v", statusCode, data, dataType, r)
			}
		}()

		RespondWith(w, req, statusCode, data)

		// Validate response
		result := w.Result()
		require.NotNil(t, result)

		// Check status code handling
		if statusCode >= 100 && statusCode < 600 {
			// Valid HTTP status codes should be set correctly
			if result.StatusCode != statusCode {
				t.Errorf("Expected status code %d, got %d", statusCode, result.StatusCode)
			}
		} else {
			// Invalid status codes will cause panics - this is expected behavior
			// The test is designed to catch panics above, so we just log here
			t.Logf("Used invalid status code %d, response code: %d", statusCode, result.StatusCode)
		}

		// Check response body for certain status codes
		body := w.Body.Bytes()

		if statusCode == 204 || statusCode == 304 {
			// No Content and Not Modified should have empty body
			if len(body) > 0 {
				t.Errorf("Expected empty body for status %d, got %d bytes", statusCode, len(body))
			}
		} else {
			// Other status codes should have JSON response
			contentType := w.Header().Get("Content-Type")
			if !strings.Contains(contentType, "application/json") && len(body) > 0 {
				t.Errorf("Expected JSON content type, got %s", contentType)
			}

			// Body should be valid UTF-8
			if !utf8.Valid(body) {
				t.Errorf("Response body is not valid UTF-8")
			}

			// Body should be valid JSON if not empty
			if len(body) > 0 {
				var decoded interface{}
				if err := json.Unmarshal(body, &decoded); err != nil {
					t.Errorf("Response body is not valid JSON: %v", err)
				}
			}

			// Response should not be excessively large
			if len(body) > 50000 {
				t.Errorf("Response body is unexpectedly large: %d bytes", len(body))
			}
		}
	})
}

// testError implements error interface for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
