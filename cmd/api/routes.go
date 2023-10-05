package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.notAllowedResponse)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.HandlerFunc(http.MethodPost, "/categories", app.createCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/categories/:id", app.showCategoryHandler)
	router.HandlerFunc(http.MethodPatch, "/categories/:id", app.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/categories/:id", app.deleteCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	return app.recoverPanic(app.secureHeaders(router))
}
