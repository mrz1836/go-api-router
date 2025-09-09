package apirouter

import (
	"encoding/json"
	"net/http"

	"github.com/mrz1836/go-logger"
)

const (
	// ErrCodeUnknown unknown error code (example)
	ErrCodeUnknown int = 600

	// StatusCodeUnknown unknown HTTP status code (example)
	StatusCodeUnknown int = 600
)

// APIError is the enriched error message for API related errors
type APIError struct {
	Code            int         `json:"code" url:"code"`                 // Associated error code
	Data            interface{} `json:"data" url:"data"`                 // Arbitrary data that is relevant
	InternalMessage string      `json:"-" url:"-"`                       // An internal message for engineers
	IPAddress       string      `json:"ip_address" url:"ip_address"`     // Current IP of user
	Method          string      `json:"method" url:"method"`             // Method requested (IE: POST)
	PublicMessage   string      `json:"message" url:"message"`           // Public error message
	RequestGUID     string      `json:"request_guid" url:"request_guid"` // Unique Request ID for tracking
	StatusCode      int         `json:"status_code" url:"status_code"`   // Associated HTTP status code (should be in request as well)
	URL             string      `json:"url" url:"url"`                   // Requesting URL
}

// ErrorFromResponse generates a new error struct using CustomResponseWriter from LogRequest()
func ErrorFromResponse(w *APIResponseWriter, internalMessage, publicMessage string, errorCode, statusCode int, data interface{}) *APIError {
	// Log the error
	logError(statusCode, internalMessage, w.RequestID, w.IPAddress)

	// Return an error
	return &APIError{
		Code:            errorCode,
		Data:            data,
		InternalMessage: internalMessage,
		IPAddress:       w.IPAddress,
		Method:          w.Method,
		PublicMessage:   publicMessage,
		RequestGUID:     w.RequestID,
		StatusCode:      statusCode,
		URL:             w.URL,
	}
}

// ErrorFromRequest gives an error without a response writer using the request
func ErrorFromRequest(req *http.Request, internalMessage, publicMessage string, errorCode, statusCode int, data interface{}) *APIError {
	// Get values from req if available
	ip, _ := GetIPFromRequest(req)
	id, _ := GetRequestID(req)

	// Log the error
	logError(statusCode, internalMessage, id, ip)

	// Return an error
	return &APIError{
		Code:            errorCode,
		Data:            data,
		InternalMessage: internalMessage,
		IPAddress:       ip,
		Method:          req.Method,
		PublicMessage:   publicMessage,
		RequestGUID:     id,
		StatusCode:      statusCode,
		URL:             req.URL.String(),
	}
}

// logError will log the internal message and code for diagnosing
func logError(statusCode int, internalMessage, requestID, ipAddress string) {
	// Skip non-error codes
	if statusCode < http.StatusBadRequest || statusCode == http.StatusNotFound {
		return
	}

	// Start with error
	logLevel := "error"

	// Switch based on known statuses
	if statusCode == http.StatusBadRequest ||
		statusCode == http.StatusUnauthorized ||
		statusCode == http.StatusMethodNotAllowed ||
		statusCode == http.StatusLocked ||
		statusCode == http.StatusForbidden ||
		statusCode == http.StatusUnprocessableEntity {
		logLevel = "warn"
	}

	// Show the login a standard way
	logger.NoFilePrintf(LogErrorFormat, requestID, ipAddress, logLevel, internalMessage, statusCode)
}

// Error returns the string error message (only public message)
func (e *APIError) Error() string {
	return e.PublicMessage
}

// ErrorCode returns the error code
func (e *APIError) ErrorCode() int {
	return e.Code
}

// JSON returns the entire public version of the error message
func (e *APIError) JSON() (string, error) {
	m, err := json.Marshal(e)
	return string(m), err
}

// Internal returns the string error message (only internal message)
func (e *APIError) Internal() string {
	return e.InternalMessage
}
