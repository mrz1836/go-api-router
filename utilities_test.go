package apirouter

import (
	"context"
	"net/http"
	"testing"

	"github.com/mrz1836/go-logger"
	"github.com/mrz1836/go-parameters"
)

// TestSnakeCase test our snake case method
func TestSnakeCase(t *testing.T) {

	//Test a valid case
	word := "testCamelCase"
	result := SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "TestCamelCase"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "TEstCamelCase"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "testCamelCASE"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "testCamel!CASE"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}
}

// TestFindString test our find string method
func TestFindString(t *testing.T) {
	//needle string, haystack []string
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
	//router := New()

	//router.HTTPRouter.GET("/test", router.Request(indexTestJSON))

	req, _ := http.NewRequest("GET", "/test?this=that&id=1234", nil)

	req = req.WithContext(context.WithValue(req.Context(), "params", parameters.ParseParams(req)))

	logger.Println(req.Context().Value("params"))

	p := req.Context().Value("params").(*parameters.Params)

	logger.Println(p.GetString("this"))

	params := parameters.GetParams(req)
	if id, success := params.GetUint64Ok("id"); !success {
		t.Fatal("failed to find the param: id", success, id, params)
	} else if id == 0 {
		t.Fatal("failed to find the param: id", success, id, params)
	}
}

// TestPermitParams test permitting parameters
/*func TestPermitParams(t *testing.T) {

	// Test parsing a url
	testUrl, err := url.Parse("https://example.com/endpoint/?param1=test1&param2=test2")
	if err != nil {
		t.Fatal("error parsing url", err)
	}

	// Test parsing values from a url
	testValues := testUrl.Query()
	param1 := testValues.Get("param1")
	param2 := testValues.Get("param2")
	if len(param1) == 0 || param1 != "test1" {
		t.Fatal("missing param1")
	} else if len(param2) == 0 || param2 != "test2" {
		t.Fatal("missing param2")
	}

	// Test permit params (testing all "all lower case"
	allowedKeys := []string{"anotherParam", "PAram1"}

	// Test the allowed keys vs the values
	PermitParams(testValues, allowedKeys)

	testParam1 := testValues.Get("param1")
	testParam2 := testValues.Get("param2")
	if testParam1 != param1 {
		t.Fatal("failed, expected param1 to eq param1:", testParam1, param1)
	}

	if testParam2 == param2 {
		t.Fatal("expected this value to be empty, removed from permit:", testParam2, param2)
	}
}*/

// TestGetIPFromRequest test getting IP from req
func TestGetIPFromRequest(t *testing.T) {

}

// TestGetRequestID test getting request ID from req
func TestGetRequestID(t *testing.T) {

}

// TestGetClientIPAddress test getting client IP
func TestGetClientIPAddress(t *testing.T) {

}
