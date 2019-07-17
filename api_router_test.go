package apirouter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// TestNew tests the New() method
func TestNew(t *testing.T) {

	// Create a new router with default properties
	router := New()

	// Check default configuration
	if !router.CrossOriginEnabled {
		t.Fatalf("expected value: %s, got: %s", "true", "false")
	}

	// Check default configuration
	if !router.CrossOriginAllowCredentials {
		t.Fatalf("expected value: %s, got: %s", "true", "false")
	}

	// Check default configuration
	if !router.CrossOriginAllowOriginAll {
		t.Fatalf("expected value: %s, got: %s", "true", "false")
	}

	// Check default configuration
	if router.CrossOriginAllowHeaders != defaultHeaders {
		t.Fatalf("expected value: %s, got: %s", defaultHeaders, router.CrossOriginAllowHeaders)
	}

	// Check default configuration
	if router.CrossOriginAllowMethods != defaultMethods {
		t.Fatalf("expected value: %s, got: %s", defaultMethods, router.CrossOriginAllowMethods)
	}

	// Make sure we have a HTTP router
	if router.HTTPRouter == nil {
		t.Fatal("expected to have http router, got nil")
	}
}

// TestReturnResponse tests the ReturnResponse()
// Only tests the basics, method is very simple
func TestReturnResponse(t *testing.T) {

	// Create new test recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Fire the index test
	indexTestNoJSON(w, req, nil)

	// Get the result
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("got an error", err.Error())
	}

	// Test the code returned
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected value: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	// Check the content type
	header := resp.Header.Get("Content-Type")
	if header != "" {
		t.Fatalf("expected value: %s, got: %s", "", header)
	}

	// Check the response
	response := string(body)
	if response != "Welcome to this simple API example!" {
		t.Fatalf("expected value: %s, got: %s", "Welcome to this simple API example!", response)
	}
}

// TestReturnResponseWithJSON tests the ReturnResponse()
// Only tests the basics, method is very simple
func TestReturnResponse_WithJSON(t *testing.T) {

	// Create new test recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Fire the index test
	indexTestJSON(w, req, nil)

	// Test the content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected value: %s, got: %s", "application/json", contentType)
	}

	// Get the result
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("got an error", err.Error())
	}

	// Test the code returned
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected value: %d, got: %d", http.StatusCreated, resp.StatusCode)
	}

	// Check the response
	response := string(body)
	if response != `{"message":"test"}` {
		t.Fatalf("expected value: %s, got: %s", `{"message":"test"}`, response)
	}
}

// TestRouter_SetCrossOriginHeaders tests SetCrossOriginHeaders() method
func TestRouter_SetCrossOriginHeaders(t *testing.T) {
	// Create new test recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := New()

	// Fire the index test
	router.SetCrossOriginHeaders(w, req, nil)

	// Test the header
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != req.Header.Get("Origin") {
		t.Fatalf("expected value: %s, got: %s", req.Header.Get("Origin"), allowOrigin)
	}

	// Test the header
	vary := w.Header().Get("Vary")
	if vary != "Origin" {
		t.Fatalf("expected value: %s, got: %s", "Origin", vary)
	}

	// Test the header
	credentials := w.Header().Get("Access-Control-Allow-Credentials")
	if credentials != "true" {
		t.Fatalf("expected value: %s, got: %s", "true", credentials)
	}

	// Test the header
	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods != defaultMethods {
		t.Fatalf("expected value: %s, got: %s", defaultMethods, methods)
	}

	// Test the header
	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers != defaultHeaders {
		t.Fatalf("expected value: %s, got: %s", defaultHeaders, headers)
	}

	// Get the result
	resp := w.Result()
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("got an error", err.Error())
	}
}

// TestRouter_SetCrossOriginHeaders_Disabled tests SetCrossOriginHeaders() method
func TestRouter_SetCrossOriginHeaders_Disabled(t *testing.T) {
	// Create new test recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := New()
	router.CrossOriginEnabled = false

	// Fire the index test
	router.SetCrossOriginHeaders(w, req, nil)

	// Test the header
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "" {
		t.Fatalf("expected value: %s, got: %s", "", allowOrigin)
	}

	// Test the header
	vary := w.Header().Get("Vary")
	if vary == "Origin" {
		t.Fatalf("expected value: %s, got: %s", "", vary)
	}

	// Test the header
	credentials := w.Header().Get("Access-Control-Allow-Credentials")
	if credentials == "true" {
		t.Fatalf("expected value: %s, got: %s", "", credentials)
	}

	// Test the header
	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods == defaultMethods {
		t.Fatalf("expected value: %s, got: %s", "", methods)
	}

	// Test the header
	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers == defaultHeaders {
		t.Fatalf("expected value: %s, got: %s", "", headers)
	}

	// Get the result
	resp := w.Result()
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("got an error", err.Error())
	}
}

// TestRouter_SetCrossOriginHeaders_CustomOrigin tests SetCrossOriginHeaders() method
func TestRouter_SetCrossOriginHeaders_CustomOrigin(t *testing.T) {
	// Create new test recorder
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router := New()
	router.CrossOriginAllowOriginAll = false
	router.CrossOriginAllowOrigin = "testdomain.com"

	// Fire the index test
	router.SetCrossOriginHeaders(w, req, nil)

	// Test the header
	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != router.CrossOriginAllowOrigin {
		t.Fatalf("expected value: %s, got: %s", router.CrossOriginAllowOrigin, allowOrigin)
	}

	// Test the header
	vary := w.Header().Get("Vary")
	if vary == "Origin" {
		t.Fatalf("expected value: %s, got: %s", "", vary)
	}

	// Get the result
	resp := w.Result()
	_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("got an error", err.Error())
	}
}

// indexTestNoJSON basic request to /
func indexTestNoJSON(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ReturnResponse(w, http.StatusOK, "Welcome to this simple API example!", false)
}

// indexTestNoJSON basic request to /
func indexTestJSON(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ReturnResponse(w, http.StatusCreated, `{"message":"test"}`, true)
}