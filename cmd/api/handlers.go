package main

import (
	"github.com/agung96tm/go-phone-test/internal/models"
	"github.com/agung96tm/go-phone-test/internal/validator"
	"net/http"
)

func (app *application) apiPhoneListHandler(w http.ResponseWriter, r *http.Request) {
	var queryFilter struct {
		OddEven string `json:"odd_even"`
	}
	qs := r.URL.Query()
	queryFilter.OddEven = app.readString(qs, "odd_even", "")

	phones, err := app.models.Phone.GetAll(queryFilter.OddEven)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"phones": phones}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) apiPhoneCreateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		Provider    string `json:"provider"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	phone := models.Phone{
		PhoneNumber: input.PhoneNumber,
		Provider:    input.Provider,
	}

	v := validator.New()
	v.CheckField(validator.NotBlank(input.PhoneNumber), "phone_number", "This field cannot be blank")
	v.CheckField(validator.NotBlank(input.Provider), "provider", "This field cannot be blank")
	v.CheckField(validator.MinChars(input.PhoneNumber, 10), "phone_number", "This field must greater or equal than 10 characters")
	v.CheckField(validator.MaxChars(input.PhoneNumber, 15), "phone_number", "This field must lower or equal than 15 characters")
	v.CheckField(validator.In(input.Provider, []string{"telkomsel", "xl", "indosat", "tri", "smartfreen"}...), "provider", "This field include wrong provider")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.FieldErrors)
		return
	}

	err = app.models.Phone.Insert(&phone)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"phone": phone}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) apiPhoneAutoHandler(w http.ResponseWriter, r *http.Request) {
	provider, phoneNumber, err := app.models.Phone.GetRandPhoneNumber()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"provider":     provider,
		"phone_number": phoneNumber,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
