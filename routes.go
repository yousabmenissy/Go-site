package main

import (
	"net/http"
	"site/ui"

	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()

	basic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	router.Handle("GET /", basic.ThenFunc(app.homeHandler))
	router.Handle("GET /test", basic.ThenFunc(app.testHandler))
	router.Handle("GET /books/{id}", basic.ThenFunc(app.bookViewHandler))
	router.Handle("GET /books", basic.ThenFunc(app.booksListHandler))
	router.Handle("GET /users/signup", basic.ThenFunc(app.signupFormHandler))
	router.Handle("POST /users/signup", basic.ThenFunc(app.POSTSignupHandler))
	router.Handle("GET /users/login", basic.ThenFunc(app.loginFormHandler))
	router.Handle("POST /users/login", basic.ThenFunc(app.POSTLoginHandler))

	protected := basic.Append(app.requireAuthentication)
	router.Handle("GET /books/create", protected.ThenFunc(app.bookFormHandler))
	router.Handle("POST /books/create", protected.ThenFunc(app.POSTBookHandler))
	router.Handle("POST /users/logout", protected.ThenFunc(app.POSTLogoutHandler))

	router.HandleFunc("GET /users/{id}/activated/{key}", app.activateUserHandler)

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handle("GET /static/", fileServer)

	standard := alice.New(app.recoverPanic)
	return standard.Then(router)
}
