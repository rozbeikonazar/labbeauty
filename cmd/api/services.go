package main

import (
	"database/sql"
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
	}
	_, err = app.models.Categories.Get(input.CategoryID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	_, err = app.models.SubCategories.Get(input.SubCategoryID)
	if err != nil {
		app.notFoundResponse(w, r)
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
