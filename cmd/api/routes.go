package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	routes := httprouter.New()

	dynamic := alice.New(app.enableCORS, app.authenticate)

	routes.Handler(http.MethodPost, "/v1/social/google/", dynamic.ThenFunc(app.apiSocialGoogleHandler))
	routes.Handler(http.MethodPost, "/v1/phones-auto/", dynamic.ThenFunc(app.apiPhoneAutoHandler))

	protected := dynamic.Append(app.requireAuthentication)
	routes.Handler(http.MethodGet, "/v1/phones/", protected.ThenFunc(app.apiPhoneListHandler))
	routes.Handler(http.MethodPost, "/v1/phones/", protected.ThenFunc(app.apiPhoneCreateHandler))

	return routes
}
