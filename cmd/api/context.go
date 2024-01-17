package main

import (
	"context"
	"github.com/agung96tm/go-phone-test/internal/models"
	"net/http"
)

type ContextKey string

const userContextKey = ContextKey("user")

func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in context")
	}
	return user
}
