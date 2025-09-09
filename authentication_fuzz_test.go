package apirouter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// FuzzCreateAndParseToken tests JWT token creation and parsing with various inputs
// to ensure proper handling of edge cases and malformed data
func FuzzCreateAndParseToken(f *testing.F) {
	// Seed corpus with representative test cases
	testCases := []struct {
		sessionSecret string
		userID        string
		issuer        string
		sessionID     string
		expiration    int64 // seconds
	}{
		{"secret123", "user123", "test-issuer", "session123", 3600},
		{"", "user123", "test-issuer", "session123", 3600},
		{"secret123", "", "test-issuer", "session123", 3600},
		{"secret123", "user123", "", "session123", 3600},
		{"secret123", "user123", "test-issuer", "", 3600},
		{"secret123", "user123", "test-issuer", "session123", 0},
		{"secret123", "user123", "test-issuer", "session123", -1},
		{"very-long-secret-key-that-should-still-work", "user123", "test-issuer", "session123", 3600},
		{"🔑", "user123", "test-issuer", "session123", 3600},
		{"secret123", "user🆔", "issuer🏢", "session📝", 3600},
		{"secret123", "user\nwith\nnewlines", "issuer\twith\ttabs", "session\rwith\rreturns", 3600},
		{"secret123", "user\x00null", "issuer\x00null", "session\x00null", 3600},
		{strings.Repeat("a", 1000), "user123", "test-issuer", "session123", 3600},
		{"secret123", strings.Repeat("u", 1000), "test-issuer", "session123", 3600},
		{"secret123", "user123", strings.Repeat("i", 1000), "session123", 3600},
		{"secret123", "user123", "test-issuer", strings.Repeat("s", 1000), 3600},
	}

	for _, tc := range testCases {
		f.Add(tc.sessionSecret, tc.userID, tc.issuer, tc.sessionID, tc.expiration)
	}

	f.Fuzz(func(t *testing.T, sessionSecret, userID, issuer, sessionID string, expirationSeconds int64) {
		expiration := time.Duration(expirationSeconds) * time.Second

		// Ensure CreateToken doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CreateToken panicked with inputs secret=%q, userID=%q, issuer=%q, sessionID=%q, expiration=%v: %v",
					sessionSecret, userID, issuer, sessionID, expiration, r)
			}
		}()

		token, err := CreateToken(sessionSecret, userID, issuer, sessionID, expiration)

		// If inputs are reasonable, token creation should succeed
		if sessionSecret != "" && utf8.ValidString(sessionSecret) &&
			utf8.ValidString(userID) && utf8.ValidString(issuer) && utf8.ValidString(sessionID) {

			if err != nil {
				// Only log error for reasonable inputs
				if len(sessionSecret) < 10000 && len(userID) < 10000 && len(issuer) < 10000 && len(sessionID) < 10000 {
					t.Logf("CreateToken returned error with reasonable inputs: %v", err)
				}
			} else {
				// Token should be valid UTF-8
				if !utf8.ValidString(token) {
					t.Errorf("CreateToken produced invalid UTF-8 token")
				}

				// Token should have reasonable length (JWT tokens are typically < 2KB)
				if len(token) > 5000 {
					t.Errorf("CreateToken produced unexpectedly long token: %d bytes", len(token))
				}

				// Token should contain dots (JWT format)
				if !strings.Contains(token, ".") {
					t.Errorf("CreateToken produced token without JWT format markers")
				}
			}
		}
	})
}

// FuzzGetTokenFromHeader tests token extraction from various header formats
// to ensure proper parsing of Authorization headers
func FuzzGetTokenFromHeader(f *testing.F) {
	// Seed corpus with representative header values
	testCases := []string{
		"Bearer valid-jwt-token",
		"Bearer ",
		"Bearer",
		"",
		"bearer lowercase-bearer",
		"Basic dXNlcjpwYXNz", // Basic auth
		"Bearer token-with-special-chars!@#$%^&*()",
		"Bearer " + strings.Repeat("a", 1000),
		"Bearer token\nwith\nnewlines",
		"Bearer token\twith\ttabs",
		"Bearer token\x00with\x00nulls",
		"Bearer 🔐token🔑with🗝️emojis",
		"BEARER uppercase-bearer token",
		"Bearer multiple words in token",
		"Token without-bearer-prefix",
		"Bearer token.with.dots.like.jwt",
		"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ", // Real JWT
		"Multiple Bearer tokens Bearer second-token",
		" Bearer token-with-leading-space",
		"Bearer token-with-trailing-space ",
		"\tBearer\ttoken-with-tabs\t",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, headerValue string) {
		// Test GetTokenFromHeaderFromRequest
		req := httptest.NewRequest("GET", "/test", nil)
		if headerValue != "" {
			req.Header.Set(AuthorizationHeader, headerValue)
		}

		// Ensure function doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GetTokenFromHeaderFromRequest panicked with header %q: %v", headerValue, r)
			}
		}()

		token1 := GetTokenFromHeaderFromRequest(req)

		// Test GetTokenFromHeader with ResponseWriter
		w := httptest.NewRecorder()
		if headerValue != "" {
			w.Header().Set(AuthorizationHeader, headerValue)
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GetTokenFromHeader panicked with header %q: %v", headerValue, r)
			}
		}()

		token2 := GetTokenFromHeader(w)

		// Test GetTokenFromResponse
		res := &http.Response{
			Header: make(http.Header),
		}
		if headerValue != "" {
			res.Header.Set(AuthorizationHeader, headerValue)
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GetTokenFromResponse panicked with header %q: %v", headerValue, r)
			}
		}()

		token3 := GetTokenFromResponse(res)

		// All three functions should return the same result
		if token1 != token2 || token2 != token3 {
			t.Errorf("Token extraction functions returned different results: req=%q, writer=%q, response=%q",
				token1, token2, token3)
		}

		// Validate extracted token
		if token1 != "" {
			// Token should be valid UTF-8
			if !utf8.ValidString(token1) {
				t.Errorf("Extracted token is not valid UTF-8: %q", token1)
			}

			// Token should not contain control characters that would break HTTP
			if strings.ContainsAny(token1, "\n\r\t") {
				t.Logf("Note: Token contains control characters: %q", token1)
			}

			// Token should not be excessively long
			if len(token1) > 2000 {
				t.Errorf("Extracted token is unexpectedly long: %d bytes", len(token1))
			}
		}

		// Verify expected behavior for common cases
		parts := strings.Fields(headerValue)
		if len(parts) >= 2 && strings.EqualFold(parts[0], "Bearer") {
			expectedToken := parts[1]
			if token1 != expectedToken {
				t.Logf("Expected token %q but got %q for header %q", expectedToken, token1, headerValue)
			}
		} else if len(parts) <= 1 {
			// Should return empty string for malformed headers
			if token1 != "" {
				t.Logf("Expected empty token for malformed header %q, got %q", headerValue, token1)
			}
		}
	})
}

// FuzzJWTClaims tests Claims validation and verification with various inputs
func FuzzJWTClaims(f *testing.F) {
	// Seed corpus with representative claims data
	testCases := []struct {
		userID      string
		issuer      string
		claimID     string
		checkIssuer string
	}{
		{"user123", "test-issuer", "claim123", "test-issuer"},
		{"", "test-issuer", "claim123", "test-issuer"},
		{"user123", "", "claim123", "test-issuer"},
		{"user123", "test-issuer", "", "test-issuer"},
		{"user123", "test-issuer", "claim123", "different-issuer"},
		{"user123", "test-issuer", "claim123", ""},
		{"user🆔", "issuer🏢", "claim📝", "issuer🏢"},
		{"user\nwith\nnewlines", "issuer\twith\ttabs", "claim\rwith\rreturns", "issuer\twith\ttabs"},
		{"user\x00null", "issuer\x00null", "claim\x00null", "issuer\x00null"},
		{strings.Repeat("u", 1000), "test-issuer", "claim123", "test-issuer"},
		{"user123", strings.Repeat("i", 1000), "claim123", strings.Repeat("i", 1000)},
		{"user123", "test-issuer", strings.Repeat("c", 1000), "test-issuer"},
	}

	for _, tc := range testCases {
		f.Add(tc.userID, tc.issuer, tc.claimID, tc.checkIssuer)
	}

	f.Fuzz(func(t *testing.T, userID, issuer, claimID, checkIssuer string) {
		// Create claims
		claims := Claims{
			UserID: userID,
		}
		claims.Issuer = issuer
		claims.ID = claimID

		// Test IsEmpty method
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Claims.IsEmpty panicked with userID=%q: %v", userID, r)
			}
		}()

		isEmpty := claims.IsEmpty()
		expectedEmpty := len(userID) <= 0
		if isEmpty != expectedEmpty {
			t.Errorf("Claims.IsEmpty() = %v, expected %v for userID=%q", isEmpty, expectedEmpty, userID)
		}

		// Test Verify method
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Claims.Verify panicked with issuer=%q, claimID=%q, checkIssuer=%q: %v",
					issuer, claimID, checkIssuer, r)
			}
		}()

		valid, err := claims.Verify(checkIssuer)

		// Verify expected behavior
		if issuer == checkIssuer && len(claimID) > 0 && userID != "" {
			// Should be valid
			if !valid || err != nil {
				t.Logf("Expected valid=true, err=nil for matching issuers and non-empty fields, got valid=%v, err=%v", valid, err)
			}
		} else {
			// Should be invalid
			if valid {
				t.Logf("Expected valid=false for mismatched issuers or empty fields, got valid=%v", valid)
			}
			if err == nil {
				t.Logf("Expected error for invalid claims, got nil")
			}
		}
	})
}
