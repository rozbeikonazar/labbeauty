package main

import (
	"fmt"
	"net/http"
	"net/url"

	"cosmetcab.dp.ua/internal/validator"
)

type FormData struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func ValidateFormData(formData *FormData, v *validator.Validator) {
	v.Check(formData.Name != "", "name", "must be provided")
	v.Check(len(formData.Name) > 1, "name", "must be at least 2 bytes long")
	v.Check(len(formData.Name) <= 100, "name", "must not be more than 100 bytes long")
	v.Check(formData.Phone != "", "phone", "must be provided")
	v.Check(v.Matches(formData.Phone, validator.PhoneRX), "phone", "must be valid phone number")
	v.Check(formData.Message != "", "message", "must be provided")
	v.Check(len(formData.Message) <= 500, "message", "must not be more than 500 bytes long")
}

func (app *application) sendToTelegramHandler(w http.ResponseWriter, r *http.Request) {
	var input FormData
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if ValidateFormData(&input, v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	formattedMessage := fmt.Sprintf("Ім'я: %s\nТелефон: %s\nПовідомлення: %s", input.Name, input.Phone, input.Message)
	// encode formatted message so it can be safely placed inside url query
	encodedMessage := url.QueryEscape(formattedMessage)

	err = app.sendToBot(encodedMessage)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "sent"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
