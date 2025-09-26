package apirouter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/matryer/respond"
)

// AllowedKeys is for allowed keys
type AllowedKeys map[string]interface{}

// ReturnResponse helps return a status code and message to the end user
// deprecated: use RespondWith instead
func ReturnResponse(w http.ResponseWriter, req *http.Request, code int, data interface{}) {
	// w.Header().Set(connectionHeader, "close")
	respond.With(w, req, code, data)
}

// ReturnJSONEncode is a mixture of ReturnResponse and JSONEncode
func ReturnJSONEncode(w http.ResponseWriter, code int, e *json.Encoder, objects interface{}, allowed []string) (err error) {
	// Set the content if JSON
	w.Header().Set(contentTypeHeader, "application/json")

	// Close the connection
	// w.Header().Set(connectionHeader, "close")

	// Set the header status code
	w.WriteHeader(code)

	// Attempt to encode the objects
	err = JSONEncode(e, objects, allowed)

	return err
}

// JSONEncodeHierarchy will execute JSONEncode for multiple nested objects
func JSONEncodeHierarchy(w io.Writer, objects, allowed interface{}) error {
	if allowed == nil {
		return json.NewEncoder(w).Encode(objects)
	}

	if slice, ok := allowed.([]string); ok {
		return JSONEncode(json.NewEncoder(w), objects, slice)
	} else if obj, found := allowed.(AllowedKeys); found {
		val := reflect.ValueOf(objects)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		data := val.Interface()
		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		numFields := t.NumField()

		fieldOutputs := make([]string, 0, numFields)

		for i := 0; i < numFields; i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = field.Name
			}
			keys, good := obj[jsonTag]
			if !good {
				continue
			}

			var buf bytes.Buffer
			buf.WriteString(`"`)
			buf.WriteString(jsonTag)
			buf.WriteString(`": `)

			fieldValue := v.Field(i)
			fieldInterface := fieldValue.Interface()
			if fieldValue.Kind() == reflect.Struct && fieldValue.CanAddr() {
				fieldInterface = fieldValue.Addr().Interface()
			}

			var sub bytes.Buffer
			err := JSONEncodeHierarchy(&sub, fieldInterface, keys)
			if err != nil {
				return err
			}
			buf.Write(sub.Bytes())

			fieldOutputs = append(fieldOutputs, buf.String())
		}

		_, _ = w.Write([]byte("{"))
		_, _ = w.Write([]byte(strings.Join(fieldOutputs, ",")))
		_, _ = w.Write([]byte("}"))
	}
	return nil
}

// JSONEncode will encode only the allowed fields of the models
func JSONEncode(e *json.Encoder, objects interface{}, allowed []string) error {
	var data []map[string]interface{}
	isMulti := false
	count := 0

	if reflect.TypeOf(objects).Kind() == reflect.Slice {
		count = reflect.ValueOf(objects).Len()
		data = make([]map[string]interface{}, count)
		isMulti = true
	}

	if isMulti {
		if count == 0 {
			return e.Encode(make([]interface{}, 0))
		}

		raw := reflect.ValueOf(objects)

		obj := jsonMap(raw.Index(0).Interface())
		toRemove := make([]string, 0)

		for k := range obj {
			if FindString(k, allowed) == -1 {
				toRemove = append(toRemove, k)
			}
		}

		for _, k := range toRemove {
			delete(obj, k)
		}

		if data != nil {
			data[0] = obj
		}

		for i := 1; i < count; i++ {
			obj = jsonMap(raw.Index(i).Interface())

			for _, k := range toRemove {
				delete(obj, k)
			}

			if data != nil {
				data[i] = obj
			}
		}

		return e.Encode(data)
	}

	obj := jsonMap(objects)
	toRemove := make([]string, 0)

	for k := range obj {
		if FindString(k, allowed) == -1 {
			toRemove = append(toRemove, k)
		}
	}

	for _, k := range toRemove {
		delete(obj, k)
	}

	return e.Encode(obj)
}

// jsonMap converts an object to a map of string interfaces
func jsonMap(obj interface{}) map[string]interface{} {
	fieldValues := make(map[string]interface{})

	var s, stringPointer reflect.Value

	// Dereference the obj if it is a pointer
	if reflect.ValueOf(obj).Kind() == reflect.Ptr {
		stringPointer = reflect.ValueOf(obj)
		s = stringPointer.Elem()
	} else {
		s = reflect.ValueOf(obj)
		// stringPointer = reflect.ValueOf(&obj)
	}

	typeOfT := s.Type()
	for i := 0; i < typeOfT.NumField(); i++ {
		structField := typeOfT.Field(i)
		fieldName := structField.Name
		if fieldName[0] != strings.ToUpper(string(fieldName[0]))[0] {
			continue
		}

		// Exclude any field starting with an underscore
		if strings.Index(fieldName, "_") == 0 {
			continue
		}
		val := s.Field(i)
		// Check for embedded types
		if structField.Anonymous {
			subFields := jsonMap(val.Interface())
			for k, v := range subFields {
				fieldValues[k] = v
			}
			continue
		}
		key := SnakeCase(fieldName)
		comps := strings.Split(key, ",")
		key = comps[0]
		fieldType := structField.Type
		if fieldType.Kind() != reflect.Ptr && val.CanAddr() {
			// fieldType = reflect.PtrTo(fieldType)
			val = val.Addr()
		}
		fieldValues[key] = val.Interface()
	}

	return fieldValues
}

// RespondWith writes a JSON response with the specified status code and data to the ResponseWriter.
// It sets the "Content-Type" header to "application/json; charset=utf-8". The data is serialized to JSON.
//
// If data is an error, it responds with a JSON object {"error": <error message>}.
// If data is nil and the status is an error (>= 400), it responds with {"error": <StatusText>, "code": <status>}.
// If the status is 204 (No Content) or 304 (Not Modified), no response body is sent.
//
// This function ensures a single response per request and is safe for use in HTTP handlers.
func RespondWith(w http.ResponseWriter, _ *http.Request, status int, data interface{}) {
	// If no content is expected, send just the status and no "body"
	if status == http.StatusNoContent || status == http.StatusNotModified {
		w.WriteHeader(status)
		return
	}

	// Convert error to a JSON error payload for better readability
	if err, ok := data.(error); ok && err != nil {
		data = map[string]interface{}{"error": err.Error()}
	}
	// Provide a default body for error status codes with no data
	if data == nil && status >= 400 {
		data = map[string]interface{}{
			"error": http.StatusText(status),
			"code":  status,
		}
	}

	// Serialize data to JSON
	responseBody, err := json.Marshal(data)
	if err != nil {
		// If serialization fails, respond with a generic error message
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"failed to encode response"}`))
		return
	}

	// Set headers and write the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
	w.WriteHeader(status)
	_, _ = w.Write(responseBody)
}
