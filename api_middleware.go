/*
Package apiMiddleware is a lightweight middleware for logging, error handling and custom response writer.
*/
package apiMiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-logger"
	"github.com/satori/go.uuid"
)

// APIResponseWriter wraps the ResponseWriter and stores the status of the request.
// It is used by the LogRequest middleware
type APIResponseWriter struct {
	http.ResponseWriter
	Buffer          bytes.Buffer  `json:"-" url:"-"`
	CacheIdentifier []string      `json:"cache_identifier" url:"cache_identifier"`
	CacheTTL        time.Duration `json:"cache_ttl" url:"cache_ttl"`
	IPAddress       string        `json:"ip_address" url:"ip_address"`
	Method          string        `json:"method" url:"method"`
	NoWrite         bool          `json:"no_write" url:"no_write"`
	RequestID       string        `json:"request_id" url:"request_id"`
	Status          int           `json:"status" url:"status"`
	URL             string        `json:"url" url:"url"`
	UserAgent       string        `json:"user_agent" url:"user_agent"`
}

// AddCacheIdentifier add cache identifier to the response writer
func (r *APIResponseWriter) AddCacheIdentifier(identifier string) {
	if r.CacheIdentifier == nil {
		r.CacheIdentifier = make([]string, 0, 2)
	}
	r.CacheIdentifier = append(r.CacheIdentifier, identifier)
}

// Header returns the http.Header that will be written to the response
func (r *APIResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// WriteHeader will write the header to the client, setting the status code
func (r *APIResponseWriter) WriteHeader(status int) {
	r.Status = status
	if !r.NoWrite {
		r.ResponseWriter.WriteHeader(status)
	}
}

// Write writes the data out to the client, if WriteHeader was not called, it will write status http.StatusOK (200)
func (r *APIResponseWriter) Write(data []byte) (int, error) {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}

	if r.NoWrite {
		return r.Buffer.Write(data)
	}

	return r.ResponseWriter.Write(data)
}

//----------------------------------------------------------------------------------------------------------------------

const (
	// ErrCodeUnknown unknown error code (example)
	ErrCodeUnknown int = 600
)

// APIError is the enriched error message for API related errors
type APIError struct {
	Code            int         `json:"code" url:"code"`                 // Associated error code
	Data            interface{} `json:"data" url:"data"`                 // Arbitrary data that is relevant
	InternalMessage string      `json:"-" url:"-"`                       // Internal message for engineers
	IPAddress       string      `json:"ip_address" url:"ip_address"`     // Current IP of user
	Method          string      `json:"method" url:"method"`             // Method requested (IE: POST)
	PublicMessage   string      `json:"message" url:"message"`           // Public error message
	RequestGUID     string      `json:"request_guid" url:"request_guid"` // Unique Request ID for tracking
	URL             string      `json:"url" url:"url"`                   // Requesting URL
}

// NewError generates a new error struct using CustomResponseWriter from LogRequest()
func NewError(w *APIResponseWriter, internalMessage string, publicMessage string, errorCode int, data interface{}) *APIError {
	return &APIError{
		Code:            errorCode,
		Data:            data,
		InternalMessage: internalMessage,
		IPAddress:       w.IPAddress,
		Method:          w.Method,
		PublicMessage:   publicMessage,
		RequestGUID:     w.RequestID,
		URL:             w.URL,
	}
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

//----------------------------------------------------------------------------------------------------------------------

// Log formats for the request
const (
	logParamsFormat = "request_id=\"%s\" method=%s path=\"%s\" ip_address=\"%s\" user_agent=\"%s\" params=%v\n"
	logTimeFormat   = "request_id=\"%s\" method=%s path=\"%s\" ip_address=\"%s\" user_agent=\"%s\" service=%dms status=%d\n"
)

// MiddlewareConfig is the configuration for the middleware service
type MiddlewareConfig struct {
	CorsEnabled          bool   `json:"cors_enabled" url:"cors_enabled"`                     // Enable or Disable Cors
	CorsAllowOriginAll   bool   `json:"cors_allow_origin_all" url:"cors_allow_origin_all"`   // Allow all origins
	CorsAllowOrigin      string `json:"cors_allow_origin" url:"cors_allow_origin"`           // Custom value for allow origin
	CorsAllowCredentials bool   `json:"cors_allow_credentials" url:"cors_allow_credentials"` // Allow credentials for BasicAuth()
	CorsAllowMethods     string `json:"cors_allow_methods" url:"cors_allow_methods"`         // Allowed methods
	CorsAllowHeaders     string `json:"cors_allow_headers" url:"cors_allow_headers"`         // Allowed headers
}

// NewMiddleware returns a middleware configuration to use for all future requests
func NewMiddleware() *MiddlewareConfig {
	config := new(MiddlewareConfig)

	// Default is to allow credentials for BasicAuth()
	config.CorsAllowCredentials = true

	// Default is for Cors to be enabled and these are common headers
	config.CorsAllowHeaders = "Accept, Content-Type, Content-Length, Cache-Control, Pragma, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Cookie"

	// Default is for the common request methods
	config.CorsAllowMethods = "POST, GET, OPTIONS, PUT, DELETE, LINK, HEAD"

	// Default is to allow all (easier to get started)
	config.CorsAllowOriginAll = true

	// Default is cors = enabled
	config.CorsEnabled = true

	// return the default configuration
	return config
}

// LogRequest will write the request to the logs before calling the handler. The request parameters will be filtered.
func (m *MiddlewareConfig) LogRequest(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		// Parse the params (once here, then store in the request)
		params := req.URL.Query()
		req = req.WithContext(context.WithValue(req.Context(), "params", params))

		// Start the custom response writer
		var writer *APIResponseWriter
		writer = &APIResponseWriter{
			IPAddress:      GetClientIPAddress(req),
			Method:         fmt.Sprintf("%s", req.Method),
			RequestID:      uuid.NewV4().String(),
			ResponseWriter: w,
			Status:         0, // future use with Etags
			URL:            fmt.Sprintf("%s", req.URL),
			UserAgent:      req.UserAgent(),
		}

		// Set cross origin on each request that goes through logging
		m.SetupCrossOrigin(writer, req, ps)

		// Start the log (timer)
		logger.Printf(logParamsFormat, writer.RequestID, writer.Method, writer.URL, writer.IPAddress, writer.UserAgent, GetParams(req))
		start := time.Now()

		// Fire the request
		h(writer, req, ps)

		// Complete the timer and final log
		elapsed := time.Since(start)
		logger.Printf(logTimeFormat, writer.RequestID, writer.Method, writer.URL, writer.IPAddress, writer.UserAgent, int64(elapsed/time.Millisecond), writer.Status)
	}
}

// BasicAuth wraps a request for Basic Authentication (RFC 2617)
func (m *MiddlewareConfig) BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string, errorMessage string) httprouter.Handle {

	// Return the function up the chain
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := req.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, req, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			m.ReturnResponse(w, http.StatusUnauthorized, errorMessage, false)
		}
	}
}

// SetupCrossOrigin sets the cross-origin headers if enabled
func (m *MiddlewareConfig) SetupCrossOrigin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Turned cors off? Just return
	if !m.CorsEnabled {
		return
	}

	// On for all origins?
	if m.CorsAllowOriginAll {
		w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Set("Vary", "Origin")
	} else { //Only the origin set by config
		w.Header().Set("Access-Control-Allow-Origin", m.CorsAllowOrigin)
	}

	// Allow credentials (used for BasicAuth)
	if m.CorsAllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Allowed methods to accept
	w.Header().Set("Access-Control-Allow-Methods", m.CorsAllowMethods)

	// Allowed headers to accept
	w.Header().Set("Access-Control-Allow-Headers", m.CorsAllowHeaders)
}

// ReturnResponse helps return a status code and message to the end user
func (m *MiddlewareConfig) ReturnResponse(w http.ResponseWriter, code int, message string, json bool) {

	// Set the header status code
	w.WriteHeader(code)

	// Set the content if JSON
	if json {
		w.Header().Set("Content-Type", "application/json")
	}

	// Write the content, log error if occurs
	if _, err := w.Write([]byte(message)); err != nil {
		logger.Data(2, logger.WARN, err.Error())
	}
}

//----------------------------------------------------------------------------------------------------------------------

// GetParams gets the params from the http request (parsed once on log request)
func GetParams(req *http.Request) url.Values {
	params := req.Context().Value("params").(url.Values)
	return params
}

// Gets the client ip address
func GetClientIPAddress(req *http.Request) string {
	//The ip address
	var ip string

	//Do we have a load balancer
	if xForward := req.Header.Get("X-Forwarded-For"); xForward != "" {
		//Set the ip as the given forwarded ip
		ip = xForward

		//Do we have more than one?
		if strings.Contains(ip, ",") {

			//Set the first ip address (from AWS)
			ip = strings.Split(ip, ",")[0]
		}
	} else {
		//Use the client address
		ip = strings.Split(req.RemoteAddr, ":")[0]

		//Remove bracket if local host
		ip = strings.Replace(ip, "[", "", 1)

		//Hack if no ip is found
		if len(ip) == 0 {
			ip = "localhost"
		}
	}

	//Return the ip address
	return ip
}

//----------------------------------------------------------------------------------------------------------------------
