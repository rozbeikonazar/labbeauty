package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"cosmetcab.dp.ua/internal/data"
	"cosmetcab.dp.ua/internal/validator"
)

func (app *application) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // max size 10MB
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Retrieve the file from the form data. The 'photo' key corresponds to the 'name' attribute
	// of the file input field in the form.
	file, header, err := r.FormFile("photo")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()
	// Create new file with the same name as the uploaded file
	dst, err := os.Create("./ui/static/img/" + header.Filename)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer dst.Close()
	// Copy the contents of the uploaded file into the new file
	if _, err := io.Copy(dst, file); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		PhotoURL    string `json:"photo_url"`
	}
	input.PhotoURL = "./ui/static/img" + header.Filename
	input.Title = r.FormValue("title")
	input.Description = r.FormValue("description")
	input.URL = r.FormValue("url")
	service := &data.Service{
		Title:       input.Title,
		Description: input.Description,
		URL:         input.URL,
		PhotoURL:    input.PhotoURL,
	}
	v := validator.New()
	if data.ValidateService(service, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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
