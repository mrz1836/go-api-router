package apirouter

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetTokenHeader will test the method SetTokenHeader()
func TestSetTokenHeader(t *testing.T) {
	t.Parallel()

	t.Run("nil writer and req, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			SetTokenHeader(nil, nil, "", 3*time.Minute)
		})
	})

	t.Run("writer, no req, no token", func(t *testing.T) {
		w := httptest.NewRecorder()
		assert.Panics(t, func() {
			SetTokenHeader(w, nil, "", 3*time.Minute)
		})
	})

	t.Run("writer, req, no token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)
		SetTokenHeader(w, req, "", 3*time.Minute)
	})

	t.Run("writer, req, valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)
		SetTokenHeader(w, req, "token", 3*time.Minute)
		assert.Equal(t, AuthorizationBearer+" token", w.Header().Get(AuthorizationHeader))
		assert.Equal(t, AuthorizationBearer+" token", req.Header.Get(AuthorizationHeader))

		var cookie *http.Cookie
		cookie, err = req.Cookie(CookieName)
		require.NoError(t, err)
		require.NotNil(t, cookie)
		assert.Equal(t, "token", cookie.Value)
	})
}

// TestGetTokenFromHeader will test the method GetTokenFromHeader()
func TestGetTokenFromHeader(t *testing.T) {
	t.Parallel()

	t.Run("nil writer, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			token := GetTokenFromHeader(nil)
			assert.Equal(t, "", token)
		})
	})

	t.Run("empty token", func(t *testing.T) {
		w := httptest.NewRecorder()
		token := GetTokenFromHeader(w)
		assert.Equal(t, "", token)
	})

	t.Run("get valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)
		SetTokenHeader(w, req, "token", 3*time.Minute)
		token := GetTokenFromHeader(w)
		assert.Equal(t, "token", token)
	})
}

// TestGetTokenFromResponse will test the method GetTokenFromResponse()
func TestGetTokenFromResponse(t *testing.T) {
	t.Parallel()

	t.Run("nil writer, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			token := GetTokenFromResponse(nil)
			assert.Equal(t, "", token)
		})
	})

	t.Run("empty token", func(t *testing.T) {
		w := httptest.NewRecorder()
		res := w.Result()
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		token := GetTokenFromResponse(res)
		assert.Equal(t, "", token)
	})

	t.Run("get valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(
			context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)
		SetTokenHeader(w, req, "token", 3*time.Minute)
		res := w.Result()
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		token := GetTokenFromResponse(res)
		assert.Equal(t, "token", token)
	})
}

// TestClearToken will test the method ClearToken()
func TestClearToken(t *testing.T) {
	t.Parallel()

	t.Run("nil writer, panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ClearToken(nil, nil)
		})
	})

	t.Run("empty token", func(_ *testing.T) {
		w := httptest.NewRecorder()
		ClearToken(w, nil)
	})

	t.Run("clear valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)
		SetTokenHeader(w, req, "token", 3*time.Minute)

		token := GetTokenFromHeader(w)
		assert.Equal(t, "token", token)

		token2 := GetTokenFromHeaderFromRequest(req)
		assert.Equal(t, "token", token2)

		var cookie *http.Cookie
		cookie, err = req.Cookie(CookieName)
		require.NotNil(t, cookie)
		require.NoError(t, err)

		token3 := cookie.Value
		assert.Equal(t, "token", token3)

		ClearToken(w, req)

		token = GetTokenFromHeader(w)
		assert.Equal(t, "", token)

		token2 = GetTokenFromHeaderFromRequest(req)
		assert.Equal(t, "", token2)

		cookie, err = req.Cookie(CookieName)
		require.NotNil(t, cookie)
		require.NoError(t, err)
		token3 = cookie.Value
		assert.Equal(t, "", token3)
	})
}

// TestClaims_CreateToken will test the method CreateToken()
func TestClaims_CreateToken(t *testing.T) {
	t.Parallel()

	sessionID, err := randomHex(32)
	require.NoError(t, err)

	var secret string
	secret, err = randomHex(16)
	require.NoError(t, err)
	assert.NotEqual(t, "", secret)

	t.Run("valid token", func(t *testing.T) {
		claims := createClaims(
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		var valid bool
		valid, err = claims.Verify("web-server-test")
		require.NoError(t, err)
		assert.True(t, valid)

		assert.False(t, claims.IsEmpty())
	})

	t.Run("user id is empty", func(t *testing.T) {
		claims := createClaims(
			"",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		assert.True(t, claims.IsEmpty())
	})

	t.Run("missing user id", func(t *testing.T) {
		claims := createClaims(
			"",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		var valid bool
		valid, err = claims.Verify("web-server-test")
		require.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("wrong issuer", func(t *testing.T) {
		claims := createClaims(
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		var valid bool
		valid, err = claims.Verify("web-server-wrong")
		require.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("missing session id", func(t *testing.T) {
		claims := createClaims(
			"123",
			"web-server-test",
			"",
			5*time.Minute,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		var valid bool
		valid, err = claims.Verify("web-server-test")
		require.Error(t, err)
		assert.False(t, valid)
	})

	t.Run("valid - empty expiration, set default", func(t *testing.T) {
		claims := createClaims(
			"123",
			"web-server-test",
			sessionID,
			0,
		)

		var token string
		token, err = claims.CreateToken(5*time.Minute, secret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		var valid bool
		valid, err = claims.Verify("web-server-test")
		require.NoError(t, err)
		assert.True(t, valid)
	})

}

// TestCreateToken will test the method CreateToken()
func TestCreateToken(t *testing.T) {
	t.Parallel()

	sessionID, err := randomHex(16)
	require.NoError(t, err)

	var secret string
	secret, err = randomHex(16)
	require.NoError(t, err)
	assert.NotEqual(t, "", secret)

	t.Run("valid token", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("missing user id", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("missing issuer", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("missing session id", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			"",
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("missing secret", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			"",
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("create token - verify", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		w := httptest.NewRecorder()

		var req *http.Request
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)

		req.Header.Add(AuthorizationHeader, AuthorizationBearer+" "+token)

		var authenticated bool
		authenticated, req, err = Check(w, req, secret, "web-server-test", 10)
		require.NoError(t, err)
		assert.NotNil(t, req)
		assert.True(t, authenticated)

		reqClaims := GetClaims(req)
		assert.Equal(t, "123", reqClaims.UserID)
		assert.Equal(t, sessionID, reqClaims.Id)
		assert.Equal(t, "web-server-test", reqClaims.Issuer)
		assert.WithinDuration(t, time.Now().UTC().Add(5*time.Minute), time.Unix(reqClaims.ExpiresAt, 0), 5*time.Second)
	})

	t.Run("verify - missing token in header", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		w := httptest.NewRecorder()

		var req *http.Request
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)

		var authenticated bool
		authenticated, req, err = Check(w, req, secret, "web-server-test", 10)
		require.Error(t, err)
		assert.Nil(t, req)
		assert.False(t, authenticated)
	})

	t.Run("verify - invalid token", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		w := httptest.NewRecorder()

		var req *http.Request
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)

		req.Header.Add(AuthorizationHeader, AuthorizationBearer+" "+token+"-invalid")

		var authenticated bool
		authenticated, req, err = Check(w, req, secret, "web-server-test", 10)
		require.Error(t, err)
		assert.Nil(t, req)
		assert.False(t, authenticated)
	})

	t.Run("verify - invalid issuer", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"123",
			"web-server-test",
			sessionID,
			5*time.Minute,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		w := httptest.NewRecorder()

		var req *http.Request
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)

		req.Header.Add(AuthorizationHeader, AuthorizationBearer+" "+token)

		var authenticated bool
		authenticated, req, err = Check(w, req, secret, "web-server-wrong", 10)
		require.Error(t, err)
		assert.Nil(t, req)
		assert.False(t, authenticated)
	})

	t.Run("verify - invalid expiration time", func(t *testing.T) {

		var token string
		token, err = CreateToken(
			secret,
			"",
			"web-server-test",
			sessionID,
			1*time.Nanosecond,
		)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		w := httptest.NewRecorder()

		var req *http.Request
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "https://domain.com", nil)
		require.NoError(t, err)
		assert.NotNil(t, req)

		req.Header.Add(AuthorizationHeader, AuthorizationBearer+" "+token)

		time.Sleep(1 * time.Second)

		var authenticated bool
		authenticated, req, err = Check(w, req, secret, "web-server-test", 10)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired by")
		assert.Nil(t, req)
		assert.False(t, authenticated)
	})
}

// TestGetClaims will test the method GetClaims()
func TestGetClaims(t *testing.T) {
	req := httptest.NewRequest(http.MethodConnect, "/", nil)
	claims := GetClaims(req)
	assert.NotNil(t, claims)
	assert.True(t, claims.IsEmpty())
}

// randomHex returns a random hex string, or an error and empty string
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
