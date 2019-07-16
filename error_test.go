package apiMiddleware

import (
	"net/http"
	"testing"
)

// setupTest creates the foundation for each error test
func setupTest() *APIResponseWriter {

	var w http.ResponseWriter

	// Start the custom response writer
	return &APIResponseWriter{
		IPAddress:      "127.0.0.1",
		Method:         "GET",
		RequestID:      "unique-guid-per-user",
		ResponseWriter: w,
		Status:         0,
		URL:            "/this/path",
		UserAgent:      "Our User Agent",
	}
}

// TestNewError test the creation of an error TestNewError()
func TestNewError(t *testing.T) {

	w := setupTest()

	err := NewError(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	if err.InternalMessage != "internal message" {
		t.Fatalf("value expected %s, value received: %s", "internal message", err.InternalMessage)
	}

	if err.PublicMessage != "public message" {
		t.Fatalf("value expected %s, value received: %s", "public message", err.PublicMessage)
	}

	if err.Code != ErrCodeUnknown {
		t.Fatalf("value expected %d, value received: %d", ErrCodeUnknown, err.Code)
	}

	if err.Data != `{"something":"else"}` {
		t.Fatalf("value expected %s, value received: %s", `{"something":"else"}`, err.Data)
	}
}

// TestError_Error test the method Error()
func TestError_Error(t *testing.T) {

	w := setupTest()

	err := NewError(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString := err.Error()
	if errorString != `public message` {
		t.Fatal("error response is not correct", errorString)
	}
}

// TestError_JSON test the method JSON()
func TestError_JSON(t *testing.T) {

	w := setupTest()

	err := NewError(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString, _ := err.JSON()
	if errorString != `{"code":600,"data":"{\"something\":\"else\"}","ip_address":"127.0.0.1","method":"GET","message":"public message","request_guid":"unique-guid-per-user","url":"/this/path"}` {
		t.Fatal("error response is not correct", errorString)
	}
}

// TestError_Internal test the method Internal()
func TestError_Internal(t *testing.T) {

	w := setupTest()

	err := NewError(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString := err.Internal()
	if errorString != `internal message` {
		t.Fatal("error response is not correct", errorString)
	}
}

// TestError_ErrorCode test the method ErrorCode()
func TestError_ErrorCode(t *testing.T) {

	w := setupTest()

	err := NewError(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	code := err.ErrorCode()
	if code != err.Code {
		t.Fatal("error response is not correct", code)
	}
}
