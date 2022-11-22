package apirouter

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	defaultExpiration = 1 * time.Hour

	// AuthorizationHeader is the auth header
	AuthorizationHeader = "Authorization"

	// AuthorizationBearer is the second part of the auth header
	AuthorizationBearer = "Bearer"

	// CookieName is for the secure cookie that also has the JWT token
	CookieName = "jwt_token"
)

// Claims is our custom JWT claims
type Claims struct {
	jwt.StandardClaims        // The standard JWT fields
	UserID             string `json:"user_id"` // The user ID set on the claims
}

// CreateToken will make a token from claims
func (c Claims) CreateToken(expiration time.Duration, sessionSecret string) (string, error) {

	// Create a new token object, specifying signing method, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaims(c.UserID, c.Issuer, c.Id, expiration))

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(sessionSecret))
}

// Verify will check the claims against known verifications
func (c Claims) Verify(issuer string) (bool, error) {

	// Invalid issuer
	if c.Issuer != issuer {
		return false, ErrIssuerMismatch
	}

	// Valid Session ID
	if len(c.Id) == 0 {
		return false, ErrInvalidSessionID
	}

	// Valid User ID
	if c.UserID == "" {
		return false, ErrInvalidUserID
	}

	return true, nil
}

// IsEmpty will detect if the claims are empty or not
func (c Claims) IsEmpty() bool {
	return len(c.UserID) <= 0
}

// createClaims will make a new set of claims for JWT
func createClaims(userID, issuer, sessionID string, expiration time.Duration) Claims {
	// Set default if not set
	if expiration <= 0 {
		expiration = defaultExpiration
	}
	return Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).UTC().Unix(),
			Id:        sessionID,
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    issuer,
			NotBefore: time.Now().UTC().Unix(),
		},
		userID,
	}
}

// CreateToken will make the claims, and then make/sign the token
func CreateToken(sessionSecret, userID, issuer, sessionID string,
	expiration time.Duration) (string, error) {

	// Create a new token object, specifying signing method, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaims(userID, issuer, sessionID, expiration))

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(sessionSecret))
}

// ClearToken will remove the token from the response and request
func ClearToken(w http.ResponseWriter, req *http.Request) {

	// Remove from response
	w.Header().Del(AuthorizationHeader)

	// Create empty cookie
	cookie := &http.Cookie{
		Path:    "/",
		Name:    CookieName,
		Value:   "",
		Expires: time.Now().Add(-24 * time.Hour),
	}

	// Remove from request
	if req != nil && req.Header != nil {
		req.Header.Del(AuthorizationHeader)
		req.Header.Del("Cookie") // Remove all cookies
		req.AddCookie(cookie)    // Add the empty cookie
	}

	// Clear any cookie out
	http.SetCookie(w, cookie)
}

// Check will check if the JWT is present and valid in the request and then extend the token
func Check(w http.ResponseWriter, r *http.Request, sessionSecret, issuer string,
	sessionAge time.Duration) (authenticated bool, req *http.Request, err error) {

	var jwtToken string

	// Look for a cookie value first
	var cookie *http.Cookie
	cookie, _ = r.Cookie(CookieName)
	if cookie != nil {
		jwtToken = cookie.Value
	} else { // Get from the auth header
		authHeaderValue := r.Header.Get(AuthorizationHeader)
		authHeader := strings.Split(authHeaderValue, AuthorizationBearer+" ")
		if len(authHeader) != 2 {
			err = ErrHeaderInvalid
			return
		}
		// Set the token value
		jwtToken = authHeader[1]
	}

	// Parse the JWT token
	var token *jwt.Token
	if token, err = jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// This error NEVER occurs since we are using our own signing method
		/*if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, ErrUnexpectedSigningMethod
		}*/
		return []byte(sessionSecret), nil
	}); err != nil {
		return
	}

	// This NEVER occurs, so far all errors have been returned first so the above case would get hit
	/* else if token == nil {
		// err = fmt.Errorf("token was nil from: %s", jwtToken)
		err = ErrTokenNil
		return
	}*/

	// Check we have claims and validity of token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {

		// Now verify the claims are good
		if _, claimErr := claims.Verify(issuer); claimErr != nil {
			// err = fmt.Errorf("claims failed validation: %w", claimErr)
			err = ErrClaimsValidationFailed
			return
		}

		// This case is NEVER hit since there is always a corresponding error
		/* else if !verified {
			err = ErrClaimsValidationFailed
			return
		}*/

		// Create new token
		var newToken string
		if newToken, err = CreateToken(
			sessionSecret,
			claims.UserID,
			issuer,
			claims.Id,
			sessionAge,
		); err != nil {
			return
		}

		// Set the token in the writer (response)
		SetTokenHeader(w, r, newToken, sessionAge)

		// Add the claims to the request for future use in router actions
		req = SetCustomData(r, claims)
		authenticated = true
	} else { // Not sure how or why this error would be triggered...
		err = ErrJWTInvalid
	}

	return
}

// GetClaims will return the current claims from the request
func GetClaims(req *http.Request) Claims {
	if claims := GetCustomData(req); claims != nil {
		return *claims.(*Claims)
	}
	return Claims{}
}

// GetTokenFromHeaderFromRequest will get the token value from the request
func GetTokenFromHeaderFromRequest(req *http.Request) string {
	headerVal := req.Header.Get(AuthorizationHeader)
	if parts := strings.Split(headerVal, " "); len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// GetTokenFromHeader will get the token value from the header
func GetTokenFromHeader(w http.ResponseWriter) string {
	headerVal := w.Header().Get(AuthorizationHeader)
	if parts := strings.Split(headerVal, " "); len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// GetTokenFromResponse will get the token value from the HTTP response
func GetTokenFromResponse(res *http.Response) string {
	headerVal := res.Header.Get(AuthorizationHeader)
	if parts := strings.Split(headerVal, " "); len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// SetTokenHeader will set the authentication token on the response and set a cookie
func SetTokenHeader(w http.ResponseWriter, r *http.Request, token string, expiration time.Duration) {

	// Set on the response
	w.Header().Set(AuthorizationHeader, AuthorizationBearer+" "+token)

	// Set on the request
	r.Header.Set(AuthorizationHeader, AuthorizationBearer+" "+token)

	// Create the cookie
	cookie := &http.Cookie{
		Path:    "/",
		Name:    CookieName,
		Value:   token,
		Expires: time.Now().UTC().Add(expiration),
		// todo: secure / http only etc
	}

	// Set the cookie on the request
	r.AddCookie(cookie)

	// Set the cookie (response)
	http.SetCookie(w, cookie)
}
