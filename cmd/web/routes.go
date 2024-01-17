package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	routes := httprouter.New()

	routes.HandlerFunc(http.MethodGet, "/phones/input", app.phoneCreateHandler)
	routes.HandlerFunc(http.MethodPost, "/phones/input", app.phoneCreatePostHandler)
	routes.HandlerFunc(http.MethodGet, "/phones/output", app.phoneListHandler)
	return routes
}
