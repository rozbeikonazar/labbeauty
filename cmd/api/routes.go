package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.notAllowedResponse)
	fileServer := http.FileServer(http.Dir("./ui/static"))
	authorizedChain := alice.New(app.recoverPanic, app.rateLimit, app.secureHeaders, app.checkAuth)
	stdChain := alice.New(app.recoverPanic, app.rateLimit, app.secureHeaders)
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	// categories routesstdChain(
	router.Handler(http.MethodGet, "/categories", stdChain.ThenFunc(app.listCategoriesHanlder))
	router.Handler(http.MethodPost, "/categories", authorizedChain.ThenFunc(app.createCategoryHandler))
	router.Handler(http.MethodGet, "/categories/:id", stdChain.ThenFunc(app.showCategoryHandler))
	router.Handler(http.MethodPatch, "/categories/:id", authorizedChain.ThenFunc(app.updateCategoryHandler))
	router.Handler(http.MethodDelete, "/categories/:id", authorizedChain.ThenFunc(app.deleteCategoryHandler))
	// subcategories routes
	router.Handler(http.MethodGet, "/subcategories", stdChain.ThenFunc(app.listSubCategoriesHandler))
	router.Handler(http.MethodPost, "/subcategories", authorizedChain.ThenFunc(app.createSubCategoryHandler))
	router.Handler(http.MethodGet, "/subcategories/:id", stdChain.ThenFunc(app.showSubCategoryHandler))
	router.Handler(http.MethodPut, "/subcategories/:id", authorizedChain.ThenFunc(app.updateSubCategoryHandler))
	router.Handler(http.MethodDelete, "/subcategories/:id", authorizedChain.ThenFunc(app.deleteSubCategoryHandler))
	// services routes
	router.Handler(http.MethodGet, "/services", stdChain.ThenFunc(app.listServicesHandler))
	router.Handler(http.MethodPost, "/services", authorizedChain.ThenFunc(app.createServiceHandler))
	router.Handler(http.MethodGet, "/services/:id", stdChain.ThenFunc(app.showServiceHandler))
	router.Handler(http.MethodPatch, "/services/:id", authorizedChain.ThenFunc(app.updateServiceHandler))
	router.Handler(http.MethodDelete, "/services/:id", authorizedChain.ThenFunc(app.deleteServiceHandler))

	router.Handler(http.MethodGet, "/services_with_subcategories/:id", stdChain.ThenFunc(app.listServicesWithSubcategoriesByCategory))
	// users routes
	router.Handler(http.MethodPost, "/user/register", authorizedChain.ThenFunc(app.registerUserHandler))
	router.Handler(http.MethodPost, "/user/login", stdChain.ThenFunc(app.loginHandler))
	// TODO change logout to POST
	router.Handler(http.MethodGet, "/user/logout", authorizedChain.ThenFunc(app.logoutHandler))
	router.Handler(http.MethodGet, "/healthcheck", authorizedChain.ThenFunc(app.healthcheckHandler))

	return router
}
