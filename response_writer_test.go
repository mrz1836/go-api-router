package apirouter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockResponseWriter is a mock implementation of http.ResponseWriter for testing purposes.
type MockResponseWriter struct {
	mock.Mock
}

// Header returns the header map that will be sent by WriteHeader.
func (m *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write writes the data to the connection as part of an HTTP reply.
func (m *MockResponseWriter) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Called(statusCode)
}

// TestAPIResponseWriter_AddCacheIdentifier tests the AddCacheIdentifier method of the APIResponseWriter struct.
func TestAPIResponseWriter_AddCacheIdentifier(t *testing.T) {
	t.Parallel()

	t.Run("initializes and appends to nil slice", func(t *testing.T) {
		w := &APIResponseWriter{}
		require.Nil(t, w.CacheIdentifier)

		w.AddCacheIdentifier("first")
		require.Len(t, w.CacheIdentifier, 1)
		require.Equal(t, "first", w.CacheIdentifier[0])
	})

	t.Run("appends additional identifiers", func(t *testing.T) {
		w := &APIResponseWriter{
			CacheIdentifier: []string{"existing"},
		}
		w.AddCacheIdentifier("new")
		require.Len(t, w.CacheIdentifier, 2)
		require.Equal(t, []string{"existing", "new"}, w.CacheIdentifier)
	})
}

// TestAPIResponseWriter_StatusCode tests the StatusCode method of the APIResponseWriter struct.
func TestAPIResponseWriter_StatusCode(t *testing.T) {
	t.Parallel()

	t.Run("returns stored status code", func(t *testing.T) {
		w := &APIResponseWriter{
			Status: 404,
		}
		require.Equal(t, 404, w.StatusCode())
	})
}

// TestAPIResponseWriter_WriteHeader tests the WriteHeader method of the APIResponseWriter struct.
func TestAPIResponseWriter_Write(t *testing.T) {
	t.Parallel()

	t.Run("writes to ResponseWriter when NoWrite is false", func(t *testing.T) {
		mockWriter := &MockResponseWriter{}
		mockWriter.On("Write", []byte("hello")).Return(5, nil)

		w := &APIResponseWriter{
			ResponseWriter: mockWriter,
			Status:         0, // Should default to 200
			NoWrite:        false,
		}

		n, err := w.Write([]byte("hello"))
		require.NoError(t, err)
		require.Equal(t, 5, n)
		require.Equal(t, http.StatusOK, w.Status)

		mockWriter.AssertExpectations(t)
	})

	t.Run("writes to internal buffer when NoWrite is true", func(t *testing.T) {
		w := &APIResponseWriter{
			NoWrite: true,
		}

		n, err := w.Write([]byte("hello buffer"))
		require.NoError(t, err)
		require.Equal(t, 12, n)
		require.Equal(t, http.StatusOK, w.Status)
		require.Equal(t, "hello buffer", w.Buffer.String())
	})

	t.Run("preserves existing status code", func(t *testing.T) {
		mockWriter := &MockResponseWriter{}
		mockWriter.On("Write", []byte("status")).Return(6, nil)

		w := &APIResponseWriter{
			ResponseWriter: mockWriter,
			Status:         http.StatusTeapot,
		}

		n, err := w.Write([]byte("status"))
		require.NoError(t, err)
		require.Equal(t, 6, n)
		require.Equal(t, http.StatusTeapot, w.Status)

		mockWriter.AssertExpectations(t)
	})
}
