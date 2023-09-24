package main

import (
	"fmt"
	"net/http"
	"runtime"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	_, file, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("error occurred in file %s, line %d\n", file, line)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "err", msg)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := map[string]any{
		"error": message,
	}
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

func (app *application) fileAlreadyExistResponse(w http.ResponseWriter, r *http.Request, filename string) {
	message := fmt.Sprintf("file %s already exist", filename)
	app.errorResponse(w, r, http.StatusConflict, message)
}
