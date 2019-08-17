package apirouter

import "github.com/julienschmidt/httprouter"

// Stack is the interface for middleware
type Stack interface {
	/*
		Adds a middleware to the stack. MWs will be
		called in the same order that they are added,
		such that:
			Use(Request ID Middleware)
			Use(Request Timing Middleware)
		would result in the request id middleware being
		the outermost layer, called first, before the
		timing middleware.
	*/
	// Use for adding new middlewares
	Use(Middleware)

	/*
		Wraps a given handle with the current stack
		from the result of Use() calls.
	*/
	// Wrap wraps the router
	Wrap(httprouter.Handle) httprouter.Handle
}

// Middleware is the Handle implementation
type Middleware func(httprouter.Handle) httprouter.Handle

// stack internal stack type
type stack struct {
	middlewares []Middleware
}

// NewStack create a new stack
func NewStack() *stack {
	return &stack{
		middlewares: []Middleware{},
	}
}

// Use adds the middle ware to the list
func (s *stack) Use(mw Middleware) {
	s.middlewares = append(s.middlewares, mw)
}

// Wrap wraps the router
func (s *stack) Wrap(fn httprouter.Handle) httprouter.Handle {
	l := len(s.middlewares)
	if l == 0 {
		return fn
	}

	// There is at least one item in the list. Starting
	// with the last item, create the handler to be
	// returned:
	var result httprouter.Handle
	result = s.middlewares[l-1](fn)

	// Reverse through the stack for the remaining elements,
	// and wrap the result with each layer:
	for i := 0; i < (l - 1); i++ {
		result = s.middlewares[l-(2+i)](result)
	}

	return result
}
