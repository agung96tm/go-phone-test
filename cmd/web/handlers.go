package main

import (
	"github.com/agung96tm/go-phone-test/internal/validator"
	"net/http"
)

func (app *application) phoneListHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "phone_list.tmpl", data)
}

type PhoneForm struct {
	PhoneNumber string `form:"phone_number"`
	Provider    string `form:"provider"`
	validator.Validator
}

func (app *application) phoneCreateHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = PhoneForm{}
	app.render(w, http.StatusOK, "phone_create.tmpl", data)
}

func (app *application) phoneCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var form PhoneForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.PhoneNumber), "phone_number", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Provider), "provider", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.PhoneNumber, 10), "phone_number", "This field must greater or equal than 10 characters")
	form.CheckField(validator.MaxChars(form.PhoneNumber, 15), "phone_number", "This field must lower or equal than 15 characters")
	form.CheckField(validator.In(form.Provider, []string{"telkomsel", "xl", "indosat", "tri", "smartfreen"}...), "provider", "This field include wrong provider")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "phone_form.tmpl", data)
		return
	}
}
