package main

import (
	"fmt"
	"net/http"

	"cosmetcab.dp.ua/internal/data"
)

func (app *application) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PhotoURL    string `json:"photo_url"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showServiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	service := data.Service{
		ID:          id,
		Title:       "Title",
		Description: "Desc",
		URL:         "url",
		PhotoURL:    "Photo",
	}

	err = app.writeJSON(w, http.StatusOK, service, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
