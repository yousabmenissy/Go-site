package main

import (
	"context"
	"net/http"

	"github.com/justinas/nosurf"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.logger.LogError.Println(err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{HttpOnly: true, Path: "/", Secure: true})

	return csrfHandler
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := app.sessionManager.GetInt64(r.Context(), "authenticatedUserID")
		if id == 0 {
			app.logger.LogDebug.Println("user not authenticated")
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.models.Users.Exists(id)
		if err != nil {
			app.logger.LogError.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if exists {
			app.logger.LogDebug.Println("user is authenticated")
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
