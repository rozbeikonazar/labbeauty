package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

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
	// Retrieve the file from the form data.
	// The 'photo' key corresponds to the 'name' attribute
	// of the file input field in the form
	file, header, err := r.FormFile("photo")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()
	fileName, err := generateUniqueImageName(header.Filename)
	if err != nil {
		app.badRequestResponse(w, r, err)
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
	var wg sync.WaitGroup
	// If validation is correct then we can upload image in a background goroutine
	app.background(func() {
		defer wg.Done()
		wg.Add(1)
		uplErr := app.azureBlobStorage.UploadBlob(fileName, &file)
		if uplErr != nil {
			// if image could not be uploaded, log error and send it to telegram
			app.logAndSendErr("image was not uploaded", header.Filename, uplErr)

		}
	})

	err = app.models.Categories.Insert(category)
	if err != nil {
		// if err occured while saving to DB, perform deletion in a background goroutine
		// of image that have been saved to blob storage
		app.background(func() {
			// wait for uploading to finish
			wg.Wait()
			delErr := app.azureBlobStorage.DeleteBlob(fileName)
			if delErr != nil {
				app.logAndSendErr("image was not deleted", header.Filename, delErr)

			}
		})
		app.dbErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/categories/%d", category.ID))

	err = app.writeJSON(w, http.StatusAccepted, envelope{"category": category}, headers)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCategoriesHanlder(w http.ResponseWriter, r *http.Request) {
	categories, err := app.models.Categories.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
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
	err = r.ParseMultipartForm(10 << 20) // max size 10MB
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Retrieve the file from the form data.
	// The 'photo' key corresponds to the 'name' attribute
	// of the file input field in the form
	file, header, err := r.FormFile("photo")
	if err == nil {
		// this means user specified the file and therefore
		// we need to upload it to the blob
		fileName, err := generateUniqueImageName(header.Filename)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
		// if file provided has supported extension, we have to delete previous blob
		// before uploading a new one
		photoURL := category.PhotoURL
		// split photoURL to get blobName
		blobName := strings.Split(photoURL, containerName)
		// delete old image in a background goroutine
		app.background(func() {
			err = app.azureBlobStorage.DeleteBlob(blobName[1])
			if err != nil {
				app.logAndSendErr("image was not deleted", photoURL, err)
			}
		})
		// upload new image in a background goroutine

		app.background(func() {
			err = app.azureBlobStorage.UploadBlob(fileName, &file)
			if err != nil {
				app.logAndSendErr("image was not uploaded", header.Filename, err)

			}
		})
		defer file.Close()
		category.PhotoURL = blobURL + containerName + fileName

	}

	title := r.FormValue("title")
	if title != "" {
		category.Title = title
	}

	description := r.FormValue("description")
	if description != "" {
		category.Description = description
	}

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

	err = app.writeJSON(w, http.StatusAccepted, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	photoURL, err := app.models.Categories.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	// split photoURL by containerName to get blobName
	blobName := strings.Split(photoURL, containerName)

	app.background(func() {
		delErr := app.azureBlobStorage.DeleteBlob(blobName[1])
		if delErr != nil {
			app.logAndSendErr("image was not deleted", photoURL, delErr)
		}
	})
	err = app.writeJSON(w, http.StatusAccepted, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
