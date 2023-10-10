package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"cosmetcab.dp.ua/internal/data"
	"cosmetcab.dp.ua/internal/validator"
)

func (app *application) createServiceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Time          sql.NullInt16 `json:"time"`
		Description   string        `json:"description"`
		Price         int           `json:"price"`
		CategoryID    int64         `json:"category_id"`
		SubCategoryID int64         `json:"subcategory_id"`
	}
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	_, err = app.models.Categories.Get(input.CategoryID)
	if err != nil {
		app.notFoundWithIDResponse(w, r, input.CategoryID)
		return
	}
	_, err = app.models.SubCategories.Get(input.SubCategoryID)
	if err != nil {
		app.notFoundWithIDResponse(w, r, input.SubCategoryID)
		return
	}
	service := &data.Service{
		Time:          input.Time,
		Description:   input.Description,
		Price:         input.Price,
		CategoryID:    input.CategoryID,
		SubCategoryID: input.SubCategoryID,
	}

	v := validator.New()

	if data.ValidateService(service, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Services.Insert(service)
	if err != nil {
		app.dbErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/services/%d", service.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"service": service}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showServiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	service, err := app.models.Services.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"service": service}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := app.models.Services.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"services": services}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateServiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	service, err := app.models.Services.Get(id)
	if err != nil {

	}

	var input struct {
		Time          *sql.NullInt16 `json:"time"`
		Description   *string        `json:"description"`
		Price         *int           `json:"price"`
		CategoryID    *int64         `json:"category_id"`
		SubCategoryID *int64         `json:"subcategory_id"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Time != nil {
		service.Time = *input.Time
	}
	if input.Description != nil {
		service.Description = *input.Description
	}
	if input.Price != nil {
		service.Price = *input.Price
	}
	if input.CategoryID != nil {
		// check if category with new id exists
		_, err := app.models.Categories.Get(*input.CategoryID)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		// if category exists then assign it to service
		service.CategoryID = *input.CategoryID

	}
	if input.SubCategoryID != nil {
		_, err := app.models.SubCategories.Get(*input.SubCategoryID)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		service.SubCategoryID = *input.SubCategoryID
	}

	v := validator.New()
	if data.ValidateService(service, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Services.Update(service)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"service": service}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	err = app.models.Services.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "service succesfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
