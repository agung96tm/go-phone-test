package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	routes := httprouter.New()

	dynamic := alice.New(app.authenticate)

	routes.Handler(http.MethodGet, "/login", dynamic.ThenFunc(app.loginHandler))
	routes.Handler(http.MethodGet, "/logout", dynamic.ThenFunc(app.logoutHandler))
	routes.Handler(http.MethodGet, "/auth/google/login", dynamic.ThenFunc(app.oauthGoogleLogin))
	routes.Handler(http.MethodGet, "/auth/google/callback", dynamic.ThenFunc(app.oauthGoogleCallback))

	protected := dynamic.Append(app.requireAuthentication)
	routes.Handler(http.MethodGet, "/phones/input", protected.ThenFunc(app.phoneCreateHandler))
	routes.Handler(http.MethodGet, "/phones/output", protected.ThenFunc(app.phoneListHandler))

	return routes
}
