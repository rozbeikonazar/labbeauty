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
	// categories routes
	router.HandlerFunc(http.MethodGet, "/categories", app.listCategoriesHanlder)
	router.HandlerFunc(http.MethodPost, "/categories", app.createCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/categories/:id", app.showCategoryHandler)
	router.HandlerFunc(http.MethodPatch, "/categories/:id", app.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/categories/:id", app.deleteCategoryHandler)
	// subcategories routes
	router.HandlerFunc(http.MethodGet, "/subcategories", app.listSubCategoriesHandler)
	router.HandlerFunc(http.MethodPost, "/subcategories", app.createSubCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/subcategories/:id", app.showSubCategoryHandler)
	router.HandlerFunc(http.MethodPut, "/subcategories/:id", app.updateSubCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/subcategories/:id", app.deleteSubCategoryHandler)
	// services routes
	router.HandlerFunc(http.MethodPost, "/services", app.createServiceHandler)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	return app.recoverPanic(app.rateLimit(app.secureHeaders(router)))
}
