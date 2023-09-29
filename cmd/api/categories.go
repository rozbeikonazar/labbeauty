package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"cosmetcab.dp.ua/internal/data"
	"cosmetcab.dp.ua/internal/validator"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if file with that name already exists
	_, err = os.ReadFile(ImagesDir + header.Filename)
	// if file does not exist then we are creating new file in specified dir
	if err == nil {
		app.fileAlreadyExistResponse(w, r, header.Filename)
		return

	}
	// Create new file with the same name as the uploaded file
	dst, err := os.Create(ImagesDir + header.Filename)
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
		PhotoURL    string `json:"photo_url"`
	}
	input.PhotoURL = ImagesDir + header.Filename
	input.Title = r.FormValue("title")
	input.Description = r.FormValue("description")

	category := &data.Category{
		Title:       input.Title,
		Description: input.Description,
		PhotoURL:    input.PhotoURL,
	}
	v := validator.New()
	if data.ValidateCategory(category, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		// delete image that have been saved to static
		os.Remove(ImagesDir + header.Filename)
		app.dbErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/categories/%d", category.ID))

	err = app.writeJSON(w, http.StatusCreated, category, headers)

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, category, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
