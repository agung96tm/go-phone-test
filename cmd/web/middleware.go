package main

import (
	"context"
	"github.com/agung96tm/go-phone-test/internal/authentication"
	"net/http"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(authentication.AccessTokenKey)
		if err == nil {
			ctx := context.WithValue(r.Context(), authentication.IsAuthenticatedKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
