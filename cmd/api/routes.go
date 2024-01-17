package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	routes := httprouter.New()
	routes.HandlerFunc(http.MethodPost, "/v1/phones-auto/", app.apiPhoneAutoHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/phones/", app.apiPhoneListHandler)
	routes.HandlerFunc(http.MethodPost, "/v1/phones/", app.apiPhoneCreateHandler)
	return app.enableCORS(routes)
}
