/*
Package main shows examples using the API Router
*/
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

// main fires on load (go run examples.go)
func main() {
	// Load the router & middleware
	router := apirouter.New()
	port := "3000"

	// Create a middleware stack
	s := apirouter.NewStack()

	// Use a Julien middleware
	s.Use(passThrough)

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", s.Wrap(router.Request(index)))

	// Set a test method (testing converting a standard handler to a handle)
	router.HTTPRouter.GET("/test", s.Wrap(router.Request(apirouter.StandardHandlerToHandle(StdHandler()))))

	// Set the options request on slash for CrossOrigin
	router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders)

	// Logout the loading of the API
	logger.Data(2, logger.DEBUG, "starting API server...", logger.MakeParameter("port", port))
	logger.Fatalln(http.ListenAndServe(":"+port, router.HTTPRouter)) //nolint:gosec // G114
}

// index basic request to /
func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var welcomeMessageJSON = map[string]interface{}{"message": "Welcome to this simple API example!"}
	apirouter.ReturnResponse(w, req, http.StatusOK, welcomeMessageJSON)
}

// passThrough is an example middleware
func passThrough(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logger.Data(2, logger.DEBUG, "middleware method hit!")
		fn(w, r, p)
	}
}

// StdHandler is an example standard handler
func StdHandler() http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		logger.Data(2, logger.DEBUG, "standard handler hit!")
	})
}
