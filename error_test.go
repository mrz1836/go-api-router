package apirouter

import (
	"fmt"
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

// TestErrorFromResponse test the creation of an error ErrorFromResponse()
func TestErrorFromResponse(t *testing.T) {
	t.Parallel()

	w := setupTest()

	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

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

// ExampleErrorFromResponse example using ErrorFromResponse()
func ExampleErrorFromResponse() {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	fmt.Println(err.Error())
	// Output:public message
}

// BenchmarkErrorFromResponse benchmarks the ErrorFromResponse() method
func BenchmarkErrorFromResponse(b *testing.B) {
	w := setupTest()
	for i := 0; i < b.N; i++ {
		_ = ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	}
}

// TestError_Error test the method Error()
func TestAPIError_Error(t *testing.T) {
	t.Parallel()

	w := setupTest()

	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString := err.Error()
	if errorString != `public message` {
		t.Fatal("error response is not correct", errorString)
	}
}

// ExampleAPIError_Error example using Error()
func ExampleAPIError_Error() {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	fmt.Println(err.Error())
	// Output:public message
}

// BenchmarkAPIError_Error benchmarks the Error() method
func BenchmarkAPIError_Error(b *testing.B) {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

// TestError_JSON test the method JSON()
func TestAPIError_JSON(t *testing.T) {
	t.Parallel()

	w := setupTest()

	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString, _ := err.JSON()
	if errorString != `{"code":600,"data":"{\"something\":\"else\"}","ip_address":"127.0.0.1","method":"GET","message":"public message","request_guid":"unique-guid-per-user","url":"/this/path"}` {
		t.Fatal("error response is not correct", errorString)
	}
}

// ExampleAPIError_JSON example using JSON()
func ExampleAPIError_JSON() {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	str, _ := err.JSON()
	fmt.Println(str)
	// Output:{"code":600,"data":"{\"something\":\"else\"}","ip_address":"127.0.0.1","method":"GET","message":"public message","request_guid":"unique-guid-per-user","url":"/this/path"}
}

// BenchmarkAPIError_JSON benchmarks the NewError() method
func BenchmarkAPIError_JSON(b *testing.B) {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	for i := 0; i < b.N; i++ {
		_, _ = err.JSON()
	}
}

// TestError_Internal test the method Internal()
func TestAPIError_Internal(t *testing.T) {
	t.Parallel()

	w := setupTest()

	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	errorString := err.Internal()
	if errorString != `internal message` {
		t.Fatal("error response is not correct", errorString)
	}
}

// ExampleAPIError_Internal example using Internal()
func ExampleAPIError_Internal() {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	fmt.Println(err.Internal())
	// Output:internal message
}

// BenchmarkAPIError_Internal benchmarks the Internal() method
func BenchmarkAPIError_Internal(b *testing.B) {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	for i := 0; i < b.N; i++ {
		_ = err.Internal()
	}
}

// TestError_ErrorCode test the method ErrorCode()
func TestAPIError_ErrorCode(t *testing.T) {
	t.Parallel()

	w := setupTest()

	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)

	code := err.ErrorCode()
	if code != err.Code {
		t.Fatal("error response is not correct", code)
	}
}

// ExampleAPIError_ErrorCode example using ErrorCode()
func ExampleAPIError_ErrorCode() {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	fmt.Println(err.ErrorCode())
	// Output:600
}

// BenchmarkAPIError_ErrorCode benchmarks the ErrorCode() method
func BenchmarkAPIError_ErrorCode(b *testing.B) {
	w := setupTest()
	err := ErrorFromResponse(w, "internal message", "public message", ErrCodeUnknown, `{"something":"else"}`)
	for i := 0; i < b.N; i++ {
		_ = err.ErrorCode()
	}
}
