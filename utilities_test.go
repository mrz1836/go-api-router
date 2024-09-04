package apirouter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mrz1836/go-parameters"
)

// TestSnakeCase test our snake case method
func TestSnakeCase(t *testing.T) {
	t.Parallel()

	// Create the list of tests
	var tests = []struct {
		input    string
		expected string
	}{
		{"testCamelCase", "test_camel_case"},
		{"TestCamelCase", "test_camel_case"},
		{"TEstCamelCase", "test_camel_case"},
		{"testCamelCASE", "test_camel_case"},
		{"testCamel!CASE", "test_camel_case"},
		{"testCamelAPI", "test_camel_api"},
		{"testCamelIP", "test_camel_ip"},
		{"testCamelURL", "test_camel_url"},
		{"testCamelJSON", "test_camel_json"},
	}

	// Test all
	for _, test := range tests {
		if output := SnakeCase(test.input); output != test.expected {
			t.Errorf("%s Failed: [%s] inputted and [%s] expected, received: [%s]", t.Name(), test.input, test.expected, output)
		}
	}
}

// TestFindString test our find string method
func TestFindString(t *testing.T) {
	t.Parallel()

	// needle string, haystack []string
	haystack := []string{"test", "stack"}
	needle := "stack"

	if index := FindString(needle, haystack); index == -1 {
		t.Fatal("FindString does not work correctly!")
	}

	if index := FindString("wrong", haystack); index >= 0 {
		t.Fatal("FindString does not work correctly!")
	}

}

// TestGetParams test getting params
func TestGetParams(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = req.WithContext(
		context.WithValue(req.Context(), parameters.ParamsKeyName, parameters.ParseParams(req)),
	)

	params := GetParams(req)
	if params == nil {
		t.Fatal("params should not be nil")
	}

	if id, success := params.GetUint64Ok("id"); !success {
		t.Fatal("failed to find the param: id", success, id, params)
	} else if id == 0 {
		t.Fatal("failed to find the param: id", success, id, params)
	}
}

// TestGetParams_BadKey tests a bad key on the context storage
func TestGetParams_BadKey(t *testing.T) {
	t.Parallel()

	type badParamKey string
	const BadParamKey badParamKey = "bad_key"

	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = req.WithContext(context.WithValue(req.Context(), BadParamKey, parameters.ParseParams(req)))

	params := GetParams(req)

	if params != nil {
		t.Fatal("params should be nil")
	}
}

// TestPermitParams tests the permitting params
func TestPermitParams(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234&private=data", strings.NewReader(""),
	)

	req = req.WithContext(context.WithValue(req.Context(), parameters.ParamsKeyName, parameters.ParseParams(req)))

	params := GetParams(req)
	if params == nil {
		t.Fatal("params should not be nil")
	}

	PermitParams(params, []string{"this", "id"})

	p, ok := params.GetStringOk("private")
	if ok {
		t.Fatal("parameter should not exit", p, ok)
	} else if len(p) > 0 {
		t.Fatal("parameter value should be empty")
	}
}

// TestGetIPFromRequest test getting IP from req
func TestGetIPFromRequest(t *testing.T) {
	t.Parallel()

	// Fake storing the ip address
	testIP := "127.0.0.1"
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)
	req = SetOnRequest(req, ipAddressKey, testIP)

	ip, ok := GetIPFromRequest(req)
	if !ok {
		t.Fatal("failed to get ip address", ip, ok)
	} else if ip != testIP {
		t.Fatal("ip address was not what was returned", ip, ok)
	}
}

// TestGetRequestID test getting request ID from req
func TestGetRequestID(t *testing.T) {
	t.Parallel()

	// Fake storing the request id
	testFakeID := "ern8347t88e7zrhg8eh48e7hg8e"
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = SetOnRequest(req, requestIDKey, testFakeID)

	id, ok := GetRequestID(req)
	if !ok {
		t.Fatal("failed to get request id", id, ok)
	} else if id != testFakeID {
		t.Fatal("request id was not what was returned", id, ok)
	}
}

// TestGetClientIPAddress test getting client ip address
func TestGetClientIPAddress(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, "/test?this=that&id=1234",
		strings.NewReader(""),
	)

	ip := GetClientIPAddress(req)
	if ip != "" {
		t.Fatal("expected ip to be localhost on the test, IP:", ip)
	}
}

// TestSetAuthToken test setting the auth token
func TestSetAuthToken(t *testing.T) {
	t.Parallel()

	// Fake storing the token
	testFakeToken := "ern8347t88e7zrhg8eh48e7hg8e" //nolint:gosec // this is a fake token
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = SetAuthToken(req, testFakeToken)

	token, ok := GetAuthToken(req)
	if !ok {
		t.Fatal("failed to get auth token", token, ok)
	} else if token != testFakeToken {
		t.Fatal("token was not what was returned", token, ok)
	}
}

// TestGetAuthToken test setting the auth token
func TestGetAuthToken(t *testing.T) {
	t.Parallel()

	// Test getting the token
	testFakeToken := "ern8347t88e7zrhg8eh48e7hg8e" //nolint:gosec // this is a fake token
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = SetAuthToken(req, testFakeToken)

	token, ok := GetAuthToken(req)
	if !ok {
		t.Fatal("failed to get auth token", token, ok)
	} else if token != testFakeToken {
		t.Fatal("token was not what was returned", token, ok)
	}
}

// TestSetUserData test setting the auth token
func TestSetUserData(t *testing.T) {
	t.Parallel()

	type TestThis struct {
		FieldName string
		FieldTwo  string
	}

	// Fake storing the ip address
	testFakeUserData := new(TestThis)
	testFakeUserData.FieldName = "this"
	testFakeUserData.FieldTwo = "that"
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = SetCustomData(req, testFakeUserData)

	data := GetCustomData(req)
	newData := data.(*TestThis)
	if newData.FieldTwo != testFakeUserData.FieldTwo {
		t.Fatal("failed get the correct data", newData.FieldTwo, testFakeUserData.FieldTwo)
	}
}

// TestGetUserData test setting the auth token
func TestGetUserData(t *testing.T) {
	t.Parallel()

	type TestThis struct {
		FieldName string
		FieldTwo  string
	}

	// Fake storing the ip address
	testFakeUserData := new(TestThis)
	testFakeUserData.FieldName = "this"
	testFakeUserData.FieldTwo = "that"
	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)

	req = SetCustomData(req, testFakeUserData)

	data := GetCustomData(req)

	newData := data.(*TestThis)
	if newData.FieldName != testFakeUserData.FieldName {
		t.Fatal("failed get the correct data", newData.FieldName, testFakeUserData.FieldName)
	}
}

// TestNoCache test using the NoCache on a request
func TestNoCache(t *testing.T) {
	t.Parallel()

	req, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet,
		"/test?this=that&id=1234", strings.NewReader(""),
	)
	w := httptest.NewRecorder()

	for _, v := range etagHeaders {
		req.Header.Set(v, "something")
	}

	NoCache(w, req)

	for _, v := range etagHeaders {
		if req.Header.Get(v) != "" {
			t.Fatal("header should be removed", v)
		}
	}
	for key := range noCacheHeaders {
		if w.Header().Get(key) == "" {
			t.Fatal("header should have been added", key)
		}
	}

}
