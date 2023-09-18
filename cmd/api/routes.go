package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.notAllowedResponse)
	router.HandlerFunc(http.MethodPost, "/services", app.createServiceHandler)
	router.HandlerFunc(http.MethodGet, "/services/:id", app.showServiceHandler)
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	return app.recoverPanic(app.secureHeaders(router))
}
