package apirouter

import "errors"

// ErrHeaderInvalid is when the header is missing or invalid
var ErrHeaderInvalid = errors.New("authorization header was mal-formatted or missing")

// ErrClaimsValidationFailed is when the claim's validation has failed
var ErrClaimsValidationFailed = errors.New("claims failed validation")

// ErrJWTInvalid is when the JWT payload is invalid
var ErrJWTInvalid = errors.New("jwt was invalid")

// ErrIssuerMismatch is when the issuer does not match the system issuer
var ErrIssuerMismatch = errors.New("issuer did not match")

// ErrInvalidSessionID is when the session id is invalid or missing
var ErrInvalidSessionID = errors.New("invalid session id detected")

// ErrInvalidUserID is when the user ID is invalid or missing
var ErrInvalidUserID = errors.New("invalid user id detected")

// ErrInvalidSigningMethod is when the signing method is invalid
var ErrInvalidSigningMethod = errors.New("invalid signing method")
