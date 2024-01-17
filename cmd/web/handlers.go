package main

import (
	"github.com/agung96tm/go-phone-test/internal/authentication"
	"net/http"
	"time"
)

func (app *application) phoneListHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "phone_list.tmpl", data)
}

func (app *application) phoneCreateHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = struct {
		PhoneNumber string `form:"phone_number"`
		Provider    string `form:"provider"`
	}{}
	app.render(w, http.StatusOK, "phone_create.tmpl", data)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     authentication.AccessTokenKey,
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Secure:   false,
	})

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func (app *application) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := app.googleOauth2.GenerateStateOauthCookie(w)
	u := app.googleOauth2.Config.AuthCodeURL(oauthState)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (app *application) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	//http.SetCookie(w, &http.Cookie{
	//	Name:     "accessToken",
	//	Value:    "testtoken",
	//	HttpOnly: false,
	//	Path:     "/",
	//	MaxAge:   3600,
	//	Secure:   false,
	//})
	//
	//acc, err := r.Cookie("accessToken")
	//if err != nil {
	//	fmt.Println("error", err)
	//	return
	//}
	//
	//fmt.Println("Cookie Value:", acc.Value)
	//
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Cookie set and retrieved successfully"))
	if valid := app.googleOauth2.StateValid(r); !valid {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	googleToken, err := app.googleOauth2.GetGoogleToken(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	token, err := app.googleOauth2.SendLoginGoogle(googleToken)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authentication.AccessTokenKey,
		Value:    token.Token,
		HttpOnly: false,
		Path:     "/",
		MaxAge:   3600,
		Secure:   false,
	})

	http.Redirect(w, r, "/phones/input", http.StatusSeeOther)
}
