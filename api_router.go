/*
Package apirouter provides lightweight HTTP middleware for CORS handling, structured logging, and standardized API response formatting.

It is designed to integrate seamlessly with Julien Schmidt's httprouter and leverages the go-logger package by MrZ for consistent logging across services.
*/
package apirouter

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-logger"
	"github.com/mrz1836/go-parameters"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// Headers for CORs and Authentication
const (
	// connectionHeader       string = "Connection"
	allowCredentialsHeader string = "Access-Control-Allow-Credentials"
	allowHeadersHeader     string = "Access-Control-Allow-Headers"
	allowMethodsHeader     string = "Access-Control-Allow-Methods"
	allowOriginHeader      string = "Access-Control-Allow-Origin"
	authenticateHeader     string = "WWW-Authenticate"
	contentTypeHeader      string = "Content-Type"
	defaultHeaders         string = "Accept, Content-Type, Content-Length, Cache-Control, Pragma, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Cookie"
	defaultMethods         string = "POST, GET, OPTIONS, PUT, DELETE, HEAD"
	exposeHeader           string = "Access-Control-Expose-Headers"
	forwardedHost          string = "x-forwarded-host"
	forwardedProtocol      string = "x-forwarded-proto"
	origin                 string = "Origin"
	varyHeaderString       string = "Vary"
)

// Log formats for the request
const (
	LogErrorFormat  string = "request_id=\"%s\" ip_address=\"%s\" type=\"%s\" internal_message=\"%s\" code=%d\n"
	LogPanicFormat  string = "request_id=\"%s\" method=\"%s\" path=\"%s\" type=\"%s\" error_message=\"%s\" stack_trace=\"%s\"\n"
	LogParamsFormat string = "request_id=\"%s\" method=\"%s\" path=\"%s\" ip_address=\"%s\" user_agent=\"%s\" params=\"%v\"\n"
	LogTimeFormat   string = "request_id=\"%s\" method=\"%s\" path=\"%s\" ip_address=\"%s\" user_agent=\"%s\" service=%dms status=%d\n"
)

// Package variables
var (
	authTokenKey  paramRequestKey = "auth_token"
	customDataKey paramRequestKey = "custom_data"
	ipAddressKey  paramRequestKey = "ip_address"
	requestIDKey  paramRequestKey = "request_id"

	// defaultFilterFields is the fields to filter from logs
	defaultFilterFields = []string{
		"api_key",
		"jwt",
		"new_password",
		"new_password_confirmation",
		"oauth",
		"oauth_token",
		"password",
		"password_check",
		"password_confirm",
		"password_confirmation",
		"social_security_number",
		"ssn",
		"token",
	}
)

// paramRequestKey for a context key
type paramRequestKey string

// Router is the configuration for the middleware service
type Router struct {
	AccessControlExposeHeaders  string               `json:"access_control_expose_headers" url:"access_control_expose_headers"`   // Allow specific headers for cors
	CrossOriginAllowCredentials bool                 `json:"cross_origin_allow_credentials" url:"cross_origin_allow_credentials"` // Allow credentials for BasicAuth()
	CrossOriginAllowHeaders     string               `json:"cross_origin_allow_headers" url:"cross_origin_allow_headers"`         // Allowed headers
	CrossOriginAllowMethods     string               `json:"cross_origin_allow_methods" url:"cross_origin_allow_methods"`         // Allowed methods
	CrossOriginAllowOrigin      string               `json:"cross_origin_allow_origin" url:"cross_origin_allow_origin"`           // Custom value for allow origin
	CrossOriginAllowOriginAll   bool                 `json:"cross_origin_allow_origin_all" url:"cross_origin_allow_origin_all"`   // Allow all origins
	CrossOriginEnabled          bool                 `json:"cross_origin_enabled" url:"cross_origin_enabled"`                     // Enable or Disable CrossOrigin
	FilterFields                []string             `json:"filter_fields" url:"filter_fields"`                                   // Filter out protected fields from logging
	HTTPRouter                  *nrhttprouter.Router `json:"-" url:"-"`                                                           // NewRelic wrapper for J Schmidt's httprouter
	Logger                      LoggerInterface      `json:"-" url:"-"`                                                           // Logger interface
	SkipLoggingPaths            []string             `json:"skip_logging_paths" url:"skip_logging_paths"`                         // Skip logging on these paths (IE: /health)
	loadedNewRelic              bool
}

// NewWithNewRelic returns a router middleware configuration with NewRelic enabled
func NewWithNewRelic(app *newrelic.Application) *Router {
	return defaultRouter(app)
}

// defaultRouter is the default settings of the Router/Config
func defaultRouter(app *newrelic.Application) (r *Router) {

	// Create a new configuration
	r = new(Router)

	// Default is cross_origin = enabled
	r.CrossOriginEnabled = true

	// The default is to allow credentials for BasicAuth()
	r.CrossOriginAllowCredentials = true

	// The default is to allow all (easier to get started)
	r.CrossOriginAllowOriginAll = true

	// Default is defaultHeaders
	r.CrossOriginAllowHeaders = defaultHeaders

	// Default is for the common request methods
	r.CrossOriginAllowMethods = defaultMethods

	// Create the router (nil if app is not set)
	r.HTTPRouter = nrhttprouter.New(app)
	r.loadedNewRelic = app != nil

	// Set the defaults
	r.setDefaults()

	// Set the filter fields to default
	r.FilterFields = defaultFilterFields

	// Set the default implementation (which can now be overridden)
	r.Logger = logger.GetImplementation()

	return
}

// setDefaults will set the router defaults
func (r *Router) setDefaults() {

	// Turn on trailing slash redirect
	r.HTTPRouter.RedirectTrailingSlash = true
	r.HTTPRouter.RedirectFixedPath = true

	// Turn on the default CORs options handler
	r.HTTPRouter.HandleOPTIONS = true
	r.HTTPRouter.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Turned cross_origin off?
		if !r.CrossOriginEnabled {
			return
		}

		// Set the header
		header := w.Header()

		// If we're using NewRelic - ignore options requests (default)
		if r.loadedNewRelic {
			txn := newrelic.FromContext(req.Context())
			txn.Ignore()
		}

		// On for all origins?
		if r.CrossOriginAllowOriginAll {

			// Normal requests use the Origin header
			originDomain := req.Header.Get(origin)
			if len(originDomain) == 0 {

				// Maybe it's behind a proxy?
				originDomain = req.Header.Get(forwardedHost)
				if len(originDomain) > 0 {
					originDomain = req.Header.Get(forwardedProtocol) + "//" + originDomain
				}
			}
			header.Set(allowOriginHeader, originDomain)
			header.Set(varyHeaderString, origin)
		} else { // Only the origin set by config
			header.Set(allowOriginHeader, r.CrossOriginAllowOrigin)
		}

		// Allow credentials (used for BasicAuth)
		if r.CrossOriginAllowCredentials {
			header.Set(allowCredentialsHeader, "true")
		}

		// Set access control
		header.Set(allowMethodsHeader, r.CrossOriginAllowMethods)
		header.Set(allowHeadersHeader, r.CrossOriginAllowHeaders)

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
}

// New returns a router middleware configuration to use for all future requests
func New() *Router {
	return defaultRouter(nil)
}

// Request will write the request to the logs before and after calling the handler
func (r *Router) Request(h httprouter.Handle) httprouter.Handle {
	return parameters.MakeHTTPRouterParsedReq(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		// Get the params from parameters.GetParams(req)
		params := GetParams(req)

		// Start the custom response writer
		// var writer *APIResponseWriter
		guid, _ := uuid.NewV4()
		writer := &APIResponseWriter{
			IPAddress:      GetClientIPAddress(req),
			Method:         req.Method,
			RequestID:      guid.String(),
			ResponseWriter: w,
			Status:         0, // future use with E-tags
			URL:            req.URL.String(),
			UserAgent:      req.UserAgent(),
		}

		// Store key information into the request that can be used by other methods
		req = SetOnRequest(req, ipAddressKey, writer.IPAddress)
		req = SetOnRequest(req, requestIDKey, writer.RequestID)

		// Set cross-origin on each request that goes through logging
		r.SetCrossOriginHeaders(writer, req, ps)

		// Set access control headers
		if len(r.AccessControlExposeHeaders) > 0 {
			w.Header().Set(exposeHeader, r.AccessControlExposeHeaders)
		}

		// Do we have paths to skip?
		// todo: this was added because some requests are confidential or "health-checks" and they can't be split apart from the router
		var skipLogging bool
		if len(r.SkipLoggingPaths) > 0 {
			for _, path := range r.SkipLoggingPaths {
				if path == req.URL.Path {
					skipLogging = true
					break
				}
			}
		}

		// Skip logging this specific request
		if !skipLogging {

			// Capture the panics and log
			defer func() {
				if err := recover(); err != nil {
					r.Logger.Printf(LogPanicFormat, writer.RequestID, writer.Method, writer.URL, "error", err.(error).Error(), strings.ReplaceAll(string(debug.Stack()), "\n", ";"))
				}
			}()

			// Start the log (timer)
			r.Logger.Printf(LogParamsFormat, writer.RequestID, writer.Method, writer.URL, writer.IPAddress, writer.UserAgent, FilterMap(params, r.FilterFields).Values)
			start := time.Now()

			// Fire the request
			h(writer, req, ps)

			// Complete the timer and final log
			elapsed := time.Since(start)
			r.Logger.Printf(LogTimeFormat, writer.RequestID, writer.Method, writer.URL, writer.IPAddress, writer.UserAgent, int64(elapsed/time.Millisecond), writer.Status)

		} else {
			// Fire the request (no logging)
			h(writer, req, ps)
		}
	})
}

// RequestNoLogging will just call the handler without any logging
// Used for API calls that do not require any logging overhead
func (r *Router) RequestNoLogging(h httprouter.Handle) httprouter.Handle {
	return parameters.MakeHTTPRouterParsedReq(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		// Start the custom response writer
		guid, _ := uuid.NewV4()
		writer := &APIResponseWriter{
			IPAddress:      GetClientIPAddress(req),
			Method:         req.Method,
			RequestID:      guid.String(),
			ResponseWriter: w,
			Status:         0, // future use with E-tags
			URL:            req.URL.String(),
			UserAgent:      req.UserAgent(),
		}

		// Store key information into the request that can be used by other methods
		req = SetOnRequest(req, ipAddressKey, writer.IPAddress)
		req = SetOnRequest(req, requestIDKey, writer.RequestID)

		// Set cross-origin on each request that goes through logging
		r.SetCrossOriginHeaders(writer, req, ps)

		// Set access control headers
		if len(r.AccessControlExposeHeaders) > 0 {
			w.Header().Set(exposeHeader, r.AccessControlExposeHeaders)
		}

		// Fire the request
		h(writer, req, ps)
	})
}

// BasicAuth wraps a request for Basic Authentication (RFC 2617)
func (r *Router) BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string, errorResponse interface{}) httprouter.Handle {

	// Return the function up the chain
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := req.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate a request to the given handle
			h(w, req, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set(authenticateHeader, "Basic realm=Restricted")
			RespondWith(w, req, http.StatusUnauthorized, errorResponse)
		}
	}
}

// SetCrossOriginHeaders sets the cross-origin headers if enabled
// todo: combine this method and the GlobalOPTIONS  http.HandlerFunc() method (@mrz had an issue combining)
func (r *Router) SetCrossOriginHeaders(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Turned cross_origin off?
	if !r.CrossOriginEnabled {
		return
	}

	// Set the header
	header := w.Header()

	// On for all origins?
	if r.CrossOriginAllowOriginAll {

		// Normal requests use the Origin header
		originDomain := req.Header.Get(origin)
		if len(originDomain) == 0 {

			// Maybe it's behind a proxy?
			originDomain = req.Header.Get(forwardedHost)
			if len(originDomain) > 0 {
				originDomain = req.Header.Get(forwardedProtocol) + "//" + originDomain
			}
		}
		header.Set(allowOriginHeader, originDomain)
		header.Set(varyHeaderString, origin)
	} else { // Only the origin set by config
		header.Set(allowOriginHeader, r.CrossOriginAllowOrigin)
	}

	// Allow credentials (used for BasicAuth)
	if r.CrossOriginAllowCredentials {
		header.Set(allowCredentialsHeader, "true")
	}

	// Set access control
	header.Set(allowMethodsHeader, r.CrossOriginAllowMethods)
	header.Set(allowHeadersHeader, r.CrossOriginAllowHeaders)

	// Adjust status code to 204 (Leaving this out, allowing customized response)
	// w.WriteHeader(http.StatusNoContent)
}
