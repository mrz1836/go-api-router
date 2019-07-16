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
	Buffer          bytes.Buffer
	CacheIdentifier []string
	CacheTTL        time.Duration
	IPAddress       string
	Method          string
	NoWrite         bool
	RequestID       string
	Status          int
	URL             string
	UserAgent       string
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
	// ErrCodeUnknown unknown error code
	ErrCodeUnknown int = 600
)

// APIError is the enriched error message for API related errors
type APIError struct {
	Code            int         `json:"code" url:"code"`
	Data            interface{} `json:"data" url:"data"`
	InternalMessage string      `json:"-" url:"-"`
	IPAddress       string      `json:"ip_address" url:"ip_address"`
	Method          string      `json:"method" url:"method"`
	PublicMessage   string      `json:"message" url:"message"`
	RequestGUID     string      `json:"request_guid" url:"request_guid"`
	URL             string      `json:"url" url:"url"`
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

// LogRequest will write the request to the logs before calling the handler. The request parameters will be filtered.
func LogRequest(h httprouter.Handle) httprouter.Handle {
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
		SetupCrossOrigin(writer, req, ps)

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

// GetParams gets the params from the http request (parsed once on log request)
func GetParams(req *http.Request) url.Values {
	params := req.Context().Value("params").(url.Values)
	return params
}

// BasicAuth wraps a request for Basic Authentication (RFC 2617)
func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := req.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, req, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

// SetupCrossOrigin sets the cross-origin headers
//todo: configure cors based on init() ?
func SetupCrossOrigin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	//w.Header().Set("Access-Control-Allow-Origin", "*") // turn off if using BasicAuth
	//w.Header().Set("Access-Control-Allow-Origin", strings.TrimRight(req.Header.Get("Referer"), "/"))
	w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, LINK, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Cache-Control, Pragma, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Cookie")
	w.Header().Set("Vary", "Origin")
}

// ReturnResponse helps return a status code and message to the user
func ReturnResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(message)) // todo: catch this error - shoot to logs
}

// ReturnJSONResponse helps return a JSON response to the user
func ReturnJSONResponse(w http.ResponseWriter, json string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(json)); err != nil {
		logger.Data(2, logger.WARN, err.Error())
	}
}

//----------------------------------------------------------------------------------------------------------------------

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
