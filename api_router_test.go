package apirouter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-logger"
)

// testStruct is for testing restricted fields
type testStruct struct {
	ID              uint64 `json:"id"`
	FieldOne        string `json:"field_one"`
	FieldTwo        string `json:"field_two"`
	RestrictedField string `json:"restricted_field"`
}

var (
	// All fields that can be displayed
	testRestrictedFields = []string{
		"id",
		"field_one",
		"field_two",
	}
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

// TestRouter_Request tests a basic request
func TestRouter_Request(t *testing.T) {

	router := New()

	router.HTTPRouter.GET("/test", router.Request(indexTestJSON))

	req, _ := http.NewRequest("GET", "/test?this=that&id=1234", nil)
	rr := httptest.NewRecorder()

	router.HTTPRouter.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Wrong status %d", status)
	}
}

// TestRouter_RequestSkipPath tests a basic request
func TestRouter_RequestSkipPath(t *testing.T) {

	router := New()
	router.SkipLoggingPaths = append(router.SkipLoggingPaths, "/health")

	router.HTTPRouter.GET("/health", router.Request(indexTestJSON))

	req, _ := http.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	router.HTTPRouter.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Wrong status %d", status)
	}
}

// TestRouter_RequestNoLogging tests a basic request
func TestRouter_RequestNoLogging(t *testing.T) {

	router := New()

	router.HTTPRouter.GET("/test", router.RequestNoLogging(indexTestJSON))

	req, _ := http.NewRequest("GET", "/test?this=that&id=1234", nil)
	rr := httptest.NewRecorder()

	router.HTTPRouter.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Wrong status %d", status)
	}
}

// TestReturnResponseWithJSON tests the ReturnResponse()
// Only tests the basics, method is very simple
func TestReturnJSONEncode(t *testing.T) {

	// Create new test recorder
	req := httptest.NewRequest("GET", "/test?this=that&id=123", nil)
	w := httptest.NewRecorder()

	// Fire the index test
	indexTestReturnJSONEncode(w, req, nil)

	// Test the content type
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
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
	response := strings.TrimSpace(string(body))
	if response != `{"field_one":"this","field_two":"that","id":123}` {
		t.Fatalf("expected value: %s, got: %s", `{"field_one":"this","field_two":"that","id":123}`, response)
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

	// Test the content type
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected value: %s, got: %s", "application/json", contentType)
	}

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

	// Check the response
	response := strings.TrimSpace(string(body))
	if response != `{"message":"Welcome to this simple API example!"}` {
		t.Fatalf("expected value: %s, got: %s", `{"message":"Welcome to this simple API example!"}`, response)
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
	if !strings.Contains(contentType, "application/json") {
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
	response := strings.TrimSpace(string(body))
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

// TestPanic will test the panic feature in Request logging
func TestPanic(t *testing.T) {
	router := New()

	router.HTTPRouter.GET("/test", router.Request(indexTestPanic))

	req, _ := http.NewRequest("GET", "/test?this=that&id=1234", nil)
	rr := httptest.NewRecorder()

	router.HTTPRouter.ServeHTTP(rr, req)
}

// indexTestPanic basic request to trigger a panic
func indexTestPanic(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	panic(fmt.Errorf("error occurred"))
}

// indexTestNoJSON basic request to /
func indexTestNoJSON(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var testDataJSON = map[string]interface{}{"message": "Welcome to this simple API example!"}
	ReturnResponse(w, req, http.StatusOK, testDataJSON)
}

// indexTestNoJSON basic request to /
func indexTestJSON(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var testDataJSON = map[string]interface{}{"message": "test"}
	ReturnResponse(w, req, http.StatusCreated, testDataJSON)
}

// indexTestNoJSON basic request to /
func indexTestReturnJSONEncode(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	testFields := new(testStruct)
	testFields.ID = 123
	testFields.FieldOne = "this"
	testFields.FieldTwo = "that"

	err := ReturnJSONEncode(w, http.StatusCreated, json.NewEncoder(w), testFields, testRestrictedFields)
	if err != nil {
		logger.Data(2, logger.ERROR, err.Error())
	}
}
