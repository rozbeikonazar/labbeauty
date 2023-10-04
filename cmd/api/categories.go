package main

import (
	"errors"
	"fmt"
	"net/http"

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
	//before uploading check if file with that name already does not exist
	//to avoid overwriting
	fileName := generateUniqueFileName(header.Filename)
	err = app.azureBlobStorage.UploadBlob(fileName, &file)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PhotoURL    string `json:"photo_url"`
	}
	input.PhotoURL = blobURL + containerName + fileName
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
		// if err occured while saving to DB, perform rollback
		// delete image that have been saved to blob storage
		delErr := app.azureBlobStorage.DeleteBlob(fileName)
		if delErr != nil {
			err = fmt.Errorf("%w; additionally, an error occured while deleting blob: %v", err, delErr)
		}
		app.dbErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/categories/%d", category.ID))

	err = app.writeJSON(w, http.StatusCreated, category, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

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

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
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
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		PhotoURL    string `json:"photo_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	category.Title = input.Title
	category.Description = input.Description
	category.PhotoURL = input.PhotoURL
	v := validator.New()
	if data.ValidateCategory(category, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Update(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, category, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
