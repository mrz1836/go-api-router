package apirouter

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

// TestUse test the use method
func TestUse(t *testing.T) {
	t.Parallel()

	s := NewStack()
	mw := func(fn httprouter.Handle) httprouter.Handle {
		return fn
	}
	c := len(s.middlewares)

	s.Use(mw)

	if len(s.middlewares) != c+1 {
		t.Error("expected Use() to increase the number of items in the InternalStack")
	}
}

// TestWrap test the wrap method
func TestWrap(t *testing.T) {
	t.Parallel()

	s := NewStack()

	var middlewareCalled bool
	mw := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			middlewareCalled = true
			fn(w, r, p)
		}
	}
	s.Use(mw)

	var handlerCalled bool
	hn := func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		handlerCalled = true
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := plainHandler(wrapped)
	handler.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("expected handler to have been called")
	}

	if !middlewareCalled {
		t.Error("expected middleware to have been called")
	}
}

// TestWrap_Ordering test wrap ordering
func TestWrap_Ordering(t *testing.T) {
	t.Parallel()

	s := NewStack()

	var firstCallAt *time.Time
	var secondCallAt *time.Time
	var thirdCallAt *time.Time
	var fourthCallAt *time.Time
	var handlerCallAt *time.Time

	first := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ts := time.Now()
			firstCallAt = &ts
			fn(w, r, p)
		}
	}

	second := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ts := time.Now()
			secondCallAt = &ts
			fn(w, r, p)
		}
	}
	third := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ts := time.Now()
			thirdCallAt = &ts
			fn(w, r, p)
		}
	}
	fourth := func(fn httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ts := time.Now()
			fourthCallAt = &ts
			fn(w, r, p)
		}
	}

	s.Use(first)
	s.Use(second)
	s.Use(third)
	s.Use(fourth)

	hn := func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		ts := time.Now()
		handlerCallAt = &ts
	}

	wrapped := s.Wrap(hn)
	req := httptest.NewRequest("GET", "/example", nil)
	w := httptest.NewRecorder()
	handler := plainHandler(wrapped)
	handler.ServeHTTP(w, req)

	if firstCallAt == nil || secondCallAt == nil || thirdCallAt == nil || fourthCallAt == nil || handlerCallAt == nil {
		t.Fatal("failed to call one or more functions")
	}

	if firstCallAt.After(*secondCallAt) || firstCallAt.After(*thirdCallAt) || firstCallAt.After(*fourthCallAt) || firstCallAt.After(*handlerCallAt) {
		t.Error("failed to call first middleware first")
	}

	if fourthCallAt.Before(*thirdCallAt) || fourthCallAt.Before(*secondCallAt) || fourthCallAt.After(*handlerCallAt) {
		t.Error("failed to call fourth middleware last before the handler")
	}

	if secondCallAt.After(*thirdCallAt) {
		t.Error("expected second middleware to come before the third")
	}
}

// TestWrap_WhenEmpty test wrap when empty
func TestWrap_WhenEmpty(t *testing.T) {
	t.Parallel()

	s := NewStack()
	hn := func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {}
	w := s.Wrap(hn)

	if reflect.ValueOf(hn).Pointer() != reflect.ValueOf(w).Pointer() {
		t.Error("expected that Wrap() would return the given function when InternalStack is empty")
	}
}

// plainHandler vanilla handler
func plainHandler(fn httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, httprouter.Params{})
	}
}

// TestStandardHandlerToHandle tests the conversion of a standard http.Handler to httprouter.Handle
func TestStandardHandlerToHandle(t *testing.T) {
	var wasCalled bool

	// Standard handler that sets a flag when called
	standardHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		wasCalled = true
		w.WriteHeader(http.StatusOK)
	})

	// Convert to httprouter.Handle
	routerHandle := StandardHandlerToHandle(standardHandler)

	// Create a test request and response
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	routerHandle(rr, req, httprouter.Params{})

	require.True(t, wasCalled, "Expected the standard handler to be called")
	require.Equal(t, http.StatusOK, rr.Code)
}

// TestStandardHandlerToMiddleware tests the conversion of a standard http.Handler to apirouter.Middleware
func TestStandardHandlerToMiddleware(t *testing.T) {
	var wasMiddlewareCalled bool

	// Standard middleware that sets a flag
	standardMiddleware := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		wasMiddlewareCalled = true
		w.WriteHeader(http.StatusAccepted)
	})

	// Convert to apirouter.Middleware
	middleware := StandardHandlerToMiddleware(standardMiddleware)

	// Final handler (wrapped by middleware)
	finalHandler := func(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		t.Fatal("Final handler should not be called directly")
	}

	// Wrap the final handler using the middleware
	wrappedHandler := middleware(finalHandler)

	// Execute the handler
	req := httptest.NewRequest(http.MethodGet, "/middleware-test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler(rr, req, httprouter.Params{})

	require.True(t, wasMiddlewareCalled, "Expected the standard middleware to be called")
	require.Equal(t, http.StatusAccepted, rr.Code)
}
