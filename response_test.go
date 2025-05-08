package apirouter

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestJSONEncode_Basic tests the JSON encoder removes fields that are not approved
func TestJSONEncode_Basic(t *testing.T) {
	t.Parallel()

	// Set up a mock struct for testing
	type TestStruct struct {
		TestKey    string `json:"test_key"`
		TestKeyTwo string `json:"test_key_two"`
		notAllowed string
	}

	// Base model and test model
	var model = new(TestStruct)
	var modelTest = new(TestStruct)
	var allowedFields = []string{"test_key", "test_key_two"} // notice omitted: notAllowed

	// Set the testing data
	model.TestKey = "TestValue1"
	model.TestKeyTwo = "TestValue2"
	model.notAllowed = "PrivateValue"

	// Set the buffer and encoder
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	// Run the encoder
	err := JSONEncode(enc, model, allowedFields)
	if err != nil {
		t.Fatal(err)
	}

	// Now unmarshal and test
	if err = json.Unmarshal(b.Bytes(), &modelTest); err != nil {
		t.Fatal(err)
	}

	// Test for our fields and values now
	if modelTest.TestKey != "TestValue1" {
		t.Fatal("TestKey does not have the right value! Encoding failed.", modelTest.TestKey)
	} else if modelTest.TestKeyTwo != "TestValue2" {
		t.Fatal("TestKeyTwo does not have the right value! Encoding failed.", modelTest.TestKeyTwo)
	} else if modelTest.notAllowed == "PrivateValue" {
		t.Fatal("Field not removed! notAllowed does not have the right value! Encoding failed.", modelTest.notAllowed)
	}
}

// TestJsonEncode_SubStruct tests the JSON encoder removes fields that are not approved
func TestJsonEncode_SubStruct(t *testing.T) {
	t.Parallel()

	// Set up a new mock sub-struct
	type TestSubStruct struct {
		TestSubKey string `json:"test_sub_key"`
	}
	// Set up a mock struct for testing
	type TestStruct struct {
		TestKey    string        `json:"test_key"`
		TestKeyTwo TestSubStruct `json:"test_key_two"`
		NotAllowed string        `json:"not_allowed"`
	}

	// Base model and test model
	var model = new(TestStruct)
	var modelTest = new(TestStruct)
	var allowedFields = []string{"test_key", "test_key_two"} // notice omitted: notAllowed

	// Set the testing data
	model.TestKey = "TestValue1"
	model.TestKeyTwo.TestSubKey = "TestSubValue"
	model.NotAllowed = "PrivateValue"

	// Set the buffer and encoder
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	// Run the encoder
	err := JSONEncode(enc, model, allowedFields)
	if err != nil {
		t.Fatal(err)
	}

	// Now unmarshal and test
	if err = json.Unmarshal(b.Bytes(), &modelTest); err != nil {
		t.Fatal(err)
	}

	// Test for our fields and values now
	if modelTest.TestKey != "TestValue1" {
		t.Fatal("TestKey does not have the right value! Encoding failed.", modelTest.TestKey)
	} else if modelTest.TestKeyTwo.TestSubKey != "TestSubValue" {
		t.Fatal("TestKeyTwo does not have the right value! Encoding failed.", modelTest.TestKeyTwo)
	} else if modelTest.NotAllowed == "PrivateValue" {
		t.Fatal("Field not removed! notAllowed does not have the right value! Encoding failed.", modelTest.NotAllowed)
	}

}

// TestJSONEncodeHierarchy tests the JSONEncodeHierarchy function
func TestJSONEncodeHierarchy(t *testing.T) {
	t.Parallel()

	type Nested struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	type Parent struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Nested Nested `json:"nested"`
	}

	tests := []struct {
		name     string
		input    interface{}
		allowed  interface{}
		expected string
	}{
		{
			name: "flat struct with allowed fields as []string",
			input: &Parent{
				ID:   1,
				Name: "Test",
			},
			allowed:  []string{"id", "name"},
			expected: `{"id":1,"name":"Test"}`,
		},
		{
			name: "nested struct with AllowedKeys",
			input: &Parent{
				ID:   99,
				Name: "NestedName",
				Nested: Nested{
					Foo: "alpha",
					Bar: 42,
				},
			},
			allowed: AllowedKeys{
				"id":   nil,
				"name": nil,
				"nested": AllowedKeys{
					"foo": nil,
				},
			},
			expected: `{"id": 99,"name": "NestedName","nested": {"foo": "alpha"}}`,
		},
		{
			name: "partial nested allowed keys",
			input: &Parent{
				ID:   7,
				Name: "Partial",
				Nested: Nested{
					Foo: "included",
					Bar: 9000,
				},
			},
			allowed: AllowedKeys{
				"nested": AllowedKeys{
					"bar": nil,
				},
			},
			expected: `{"nested": {"bar": 9000}}`,
		},
		{
			name: "empty allowed keys yields empty object",
			input: &Parent{
				ID:     1,
				Name:   "NoFields",
				Nested: Nested{Foo: "A", Bar: 1},
			},
			allowed:  AllowedKeys{},
			expected: `{}`,
		},
		{
			name: "fallback: unknown allowed type does nothing",
			input: &Parent{
				ID:   1,
				Name: "Ignore",
			},
			allowed:  123, // unsupported
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := JSONEncodeHierarchy(&buf, tt.input, tt.allowed)
			require.NoError(t, err)

			// Compact both to normalize spacing and formatting
			var actualObj, expectedObj any
			if tt.expected != `` {
				// t.Log(string(bytes.TrimSpace(buf.Bytes())))
				require.NoError(t, json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &actualObj))
				require.NoError(t, json.Unmarshal([]byte(tt.expected), &expectedObj))
				require.Equal(t, expectedObj, actualObj)
			} else {
				require.Equal(t, tt.expected, buf.String())
			}
		})
	}
}

// TestJSONEncode tests the JSONEncode function
func TestJSONEncode(t *testing.T) {
	t.Parallel()

	type Example struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	t.Run("encodes single struct with allowed fields", func(t *testing.T) {
		input := Example{ID: 1, Name: "Alice", Email: "alice@example.com"}

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{"id", "name"})
		require.NoError(t, err)

		var output map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)

		require.InDelta(t, float64(1), output["id"], 0)
		require.Equal(t, "Alice", output["name"])
		require.NotContains(t, output, "email")
	})

	t.Run("encodes slice of structs with allowed fields", func(t *testing.T) {
		input := []Example{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
		}

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), input, []string{"id", "email"})
		require.NoError(t, err)

		var output []map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)
		require.Len(t, output, 2)

		require.InDelta(t, float64(1), output[0]["id"], 0)
		require.Equal(t, "alice@example.com", output[0]["email"])
		require.InDelta(t, float64(2), output[1]["id"], 0)
		require.Equal(t, "bob@example.com", output[1]["email"])

		require.NotContains(t, output[0], "name")
		require.NotContains(t, output[1], "name")
	})

	t.Run("encodes empty slice to empty array", func(t *testing.T) {
		var input []Example

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), input, []string{"id"})
		require.NoError(t, err)

		require.Equal(t, "[]\n", buf.String())
	})

	t.Run("excludes unexported and underscored fields", func(t *testing.T) {
		type Test struct {
			ID       int    `json:"id"`
			_name    string // underscore-prefixed
			private  string // unexported
			Exported string `json:"exported"`
		}

		input := Test{ID: 123, _name: "hidden", private: "hidden", Exported: "ok"}
		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{"id", "exported"})
		require.NoError(t, err)

		var output map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)

		require.InDelta(t, float64(123), output["id"], 0)
		require.Equal(t, "ok", output["exported"])
		require.NotContains(t, output, "_name")
		require.NotContains(t, output, "private")
	})

	t.Run("includes fields from embedded structs", func(t *testing.T) {
		type Inner struct {
			InnerField string `json:"inner_field"`
		}
		type Outer struct {
			ID int `json:"id"`
			Inner
		}

		input := Outer{ID: 1, Inner: Inner{InnerField: "nested"}}
		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{"id", "inner_field"})
		require.NoError(t, err)

		var output map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)

		require.InDelta(t, float64(1), output["id"], 0)
		require.Equal(t, "nested", output["inner_field"])
	})

	t.Run("handles pointer fields", func(t *testing.T) {
		type WithPtr struct {
			Name *string `json:"name"`
		}

		str := "pointer name"
		input := WithPtr{Name: &str}

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{"name"})
		require.NoError(t, err)

		var output map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)

		require.Equal(t, "pointer name", output["name"])
	})

	t.Run("omits disallowed fields", func(t *testing.T) {
		type Person struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		input := Person{ID: 42, Name: "John", Email: "john@example.com"}

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{"id"}) // omit name, email
		require.NoError(t, err)

		var output map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &output)
		require.NoError(t, err)

		require.InDelta(t, float64(42), output["id"], 0)
		require.NotContains(t, output, "name")
		require.NotContains(t, output, "email")
	})

	t.Run("encodes only empty object if no fields are allowed", func(t *testing.T) {
		type Foo struct {
			One string `json:"one"`
			Two string `json:"two"`
		}
		input := Foo{One: "a", Two: "b"}

		var buf bytes.Buffer
		err := JSONEncode(json.NewEncoder(&buf), &input, []string{}) // no allowed fields
		require.NoError(t, err)
		require.JSONEq(t, `{}`, buf.String())
	})

}

// TestRespondWith tests the RespondWith function
func TestRespondWith(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		data         interface{}
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Struct data success",
			status:       http.StatusOK,
			data:         map[string]string{"hello": "world"},
			expectedCode: http.StatusOK,
			expectedBody: `{"hello":"world"}`,
		},
		{
			name:         "Error data with client error status",
			status:       http.StatusBadRequest,
			data:         errors.New("invalid input"),
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid input"}`,
		},
		{
			name:         "Nil data with error status",
			status:       http.StatusNotFound,
			data:         nil,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Not Found","code":404}`,
		},
		{
			name:         "Nil data with no content status",
			status:       http.StatusNoContent,
			data:         nil,
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:         "Nil data with success status",
			status:       http.StatusOK,
			data:         nil,
			expectedCode: http.StatusOK,
			expectedBody: "null",
		},
		{
			name:         "Unencodable data triggers internal error",
			status:       http.StatusOK,
			data:         make(chan int), // not JSON-encodable
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"failed to encode response"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			RespondWith(rr, req, tt.status, tt.data)

			// Verify the status code matches expectation
			require.Equal(t, tt.expectedCode, rr.Code, "HTTP status code mismatch")

			// Verify Content-Type header (should be set for all responses with a body)
			if tt.expectedCode != http.StatusNoContent && tt.expectedCode != http.StatusNotModified {
				require.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			} else {
				require.Empty(t, rr.Header().Get("Content-Type"), "Content-Type should be empty for status %d", tt.expectedCode)
			}

			// Verify response body content
			bodyStr := rr.Body.String()
			if tt.expectedBody == "" {
				// For no-content responses, the body should be empty
				require.Zero(t, rr.Body.Len(), "Expected no body for status %d", tt.expectedCode)
			} else {
				// Use JSONEq to ignore key order or formatting differences in JSON
				require.JSONEq(t, tt.expectedBody, bodyStr, "Response body JSON mismatch")
				// If Content-Length is set, verify it matches the body length
				if cl := rr.Header().Get("Content-Length"); cl != "" {
					require.Equal(t, strconv.Itoa(len(bodyStr)), cl, "Content-Length header mismatch")
				}
			}
		})
	}
}

// TestRespondWith_Expanded tests the RespondWith function with various data types
func TestRespondWith_Expanded(t *testing.T) {
	type customStruct struct {
		Message string `json:"message"`
		Value   int    `json:"value"`
	}
	type pointerStruct struct {
		ID int `json:"id"`
	}

	tests := []struct {
		name         string
		status       int
		data         interface{}
		expectedCode int
		expectedJSON string // leave empty when the body should be empty
	}{
		{
			name:         "Map success",
			status:       http.StatusOK,
			data:         map[string]interface{}{"foo": "bar"},
			expectedCode: http.StatusOK,
			expectedJSON: `{"foo":"bar"}`,
		},
		{
			name:         "Struct success",
			status:       http.StatusCreated,
			data:         customStruct{Message: "created", Value: 42},
			expectedCode: http.StatusCreated,
			expectedJSON: `{"message":"created","value":42}`,
		},
		{
			name:         "Slice success",
			status:       http.StatusOK,
			data:         []string{"a", "b", "c"},
			expectedCode: http.StatusOK,
			expectedJSON: `["a","b","c"]`,
		},
		{
			name:         "Primitive string",
			status:       http.StatusOK,
			data:         "hello world",
			expectedCode: http.StatusOK,
			expectedJSON: `"hello world"`,
		},
		{
			name:         "Primitive int",
			status:       http.StatusOK,
			data:         123,
			expectedCode: http.StatusOK,
			expectedJSON: `123`,
		},
		{
			name:         "Pointer to struct",
			status:       http.StatusOK,
			data:         &pointerStruct{ID: 7},
			expectedCode: http.StatusOK,
			expectedJSON: `{"id":7}`,
		},
		{
			name:         "Error value with 400",
			status:       http.StatusBadRequest,
			data:         errors.New("bad input"),
			expectedCode: http.StatusBadRequest,
			expectedJSON: `{"error":"bad input"}`,
		},
		{
			name:         "Nil payload with 404",
			status:       http.StatusNotFound,
			data:         nil,
			expectedCode: http.StatusNotFound,
			expectedJSON: `{"error":"Not Found","code":404}`,
		},
		{
			name:         "Nil payload with 200 (explicit null)",
			status:       http.StatusOK,
			data:         nil,
			expectedCode: http.StatusOK,
			expectedJSON: `null`,
		},
		{
			name:         "No‑content 204",
			status:       http.StatusNoContent,
			data:         nil,
			expectedCode: http.StatusNoContent,
			expectedJSON: ``,
		},
		{
			name:         "Unencodable channel",
			status:       http.StatusOK,
			data:         make(chan int),
			expectedCode: http.StatusInternalServerError,
			expectedJSON: `{"error":"failed to encode response"}`,
		},
		{
			name:         "Error status 422 with struct payload",
			status:       http.StatusUnprocessableEntity,
			data:         map[string]string{"reason": "validation failed"},
			expectedCode: http.StatusUnprocessableEntity,
			expectedJSON: `{"reason":"validation failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// New request/recorder per sub‑test to avoid cross‑contamination
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()

			RespondWith(rr, req, tt.status, tt.data)

			// --- Status code ---
			require.Equal(t, tt.expectedCode, rr.Code, "status code")

			// --- Body / headers ---
			if tt.expectedJSON == "" {
				// Expect empty body for 204 & 304
				require.Zero(t, rr.Body.Len(), "body should be empty")
				require.Empty(t, rr.Header().Get("Content-Type"))
				return
			}

			require.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
			require.JSONEq(t, tt.expectedJSON, rr.Body.String(), "body JSON mismatch")

			// Content‑Length should match body length when set
			if cl := rr.Header().Get("Content-Length"); cl != "" {
				require.Equal(t, strconv.Itoa(len(rr.Body.Bytes())), cl, "Content-Length mismatch")
			}
		})
	}
}

// ----------   Benchmarks   ----------

// Note: ReturnResponse is assumed to be the legacy wrapper that still calls matryer/respond.With.
// If its signature differs, adjust the benchmark accordingly.

// BenchmarkRespondWith tests the performance of RespondWith
func BenchmarkRespondWith(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	payload := map[string]string{"hello": "world"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		RespondWith(rr, req, http.StatusOK, payload)
	}
}

// BenchmarkReturnResponse tests the performance of ReturnResponse
func BenchmarkReturnResponse(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	payload := map[string]string{"hello": "world"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		ReturnResponse(rr, req, http.StatusOK, payload)
	}
}
