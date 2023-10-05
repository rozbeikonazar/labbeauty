package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/lib/pq"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	_, file, line, _ := runtime.Caller(2)
	msg := fmt.Sprintf("error occurred in file %s, line %d\n", file, line)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "err", msg)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := envelope{"error": message}
	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) notAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Method %s is not allowed for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) dbErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// if error contains UNIQUE violations then serve an error response
	// with first value that caused it
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		detail := pgErr.Detail
		value := strings.Split(detail, "=")[1]
		message := fmt.Sprintf("A resource with the same identifier (%s)", value)
		app.errorResponse(w, r, http.StatusConflict, message)

	} else {
		app.serverErrorResponse(w, r, err)
	}
}
