package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

func main() {
	// Load the router & middleware
	router := apirouter.New()
	port := "3000"

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", router.Request(index))

	// Set the options request on slash for Cors
	router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders)

	// Logout the loading of the API
	logger.Data(2, logger.DEBUG, "starting API server...", logger.MakeParameter("port", port))
	logger.Fatalln(http.ListenAndServe(":"+port, router.HTTPRouter))
}

// index basic request to /
func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	_, _ = fmt.Fprint(w, "Welcome to this simple API example!")
}
