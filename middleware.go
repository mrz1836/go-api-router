package apirouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Stack is the interface for middleware
type Stack interface {
	/*
		Adds middleware to the InternalStack. MWs will be
		called in the same order that they are added
		such that:
			Use(Request ID Middleware)
			Use(Request Timing Middleware)
		would result in the request id middleware being
		the outermost layer, called first, before the
		timing middleware.
	*/
	// Use for adding new middlewares
	Use(m Middleware)

	/*
		Wraps a given handle with the current InternalStack
		from the result of Use() calls.
	*/
	// Wrap wraps the router
	Wrap(h httprouter.Handle) httprouter.Handle
}

// Middleware is the Handle implementation
type Middleware func(httprouter.Handle) httprouter.Handle

// InternalStack internal stack type
type InternalStack struct {
	middlewares []Middleware
}

// NewStack will create an InternalStack struct
func NewStack() *InternalStack {
	return &InternalStack{
		middlewares: []Middleware{},
	}
}

// Use adds the middleware to the list
func (s *InternalStack) Use(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

// Wrap wraps the router
func (s *InternalStack) Wrap(fn httprouter.Handle) httprouter.Handle {
	l := len(s.middlewares)
	if l == 0 {
		return fn
	}

	// There is at least one item in the list. Starting
	// with the last item, create the handler to be
	// returned:
	var result httprouter.Handle
	result = s.middlewares[l-1](fn)

	// Reverse through the InternalStack for the remaining elements and wrap the result with each layer:
	for i := 0; i < (l - 1); i++ {
		result = s.middlewares[l-(2+i)](result)
	}

	return result
}

// StandardHandlerToHandle converts a standard middleware to Julien handle version
func StandardHandlerToHandle(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		next.ServeHTTP(w, r)
	}
}

// StandardHandlerToMiddleware converts standard middleware to type Middleware
func StandardHandlerToMiddleware(next http.Handler) Middleware {
	return func(_ httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			next.ServeHTTP(w, r)
		}
	}
}
