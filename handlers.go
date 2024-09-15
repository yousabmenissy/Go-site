package main

import (
	"errors"
	"fmt"
	"net/http"
	"site/data/books"
	"site/data/tokens"
	"site/data/users"
	"site/internal"
	"site/internal/validation"
	"strconv"
	"strings"
	"time"
)

func (app *Application) errorResponse(w http.ResponseWriter, status int, err error) {
	app.logger.LogError.Output(2, err.Error())
	w.WriteHeader(status)
}

func (app *Application) testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "handler works!")
}

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/home" {
		data := app.NewTemplateData(r)
		app.render(w, http.StatusOK, "home.html", data)
	} else {
		w.WriteHeader(404)
	}
}

func (app *Application) POSTBookHandler(w http.ResponseWriter, r *http.Request) {
	form := books.BookSubmitForm{V: validation.New()}
	err := app.ReadJSON(r, &form.Input)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}
	app.logger.LogDebug.Printf("%+v", form.Input)

	err = form.Prepare()
	app.logger.LogDebug.Printf("%+v\n", form.V.Errors)
	if err != nil {
		app.errorResponse(w, http.StatusBadRequest, err)
		return
	}
	if !form.V.Valid() {
		app.logger.LogError.Println("not valid")
		app.WriteJSON(w, r, form.V.Errors, nil, http.StatusUnprocessableEntity)
		return
	}

	err = app.models.Books.Insert(&form.Output)
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}
	app.logger.LogDebug.Printf("%d\n", form.Output.ID)

	w.WriteHeader(http.StatusOK)
}

func (app *Application) bookFormHandler(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	app.render(w, http.StatusOK, "books.Create.html", data)
}

func (app *Application) bookViewHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Host, "api.") {
		app.bookViewJsonHandler(w, r)
		return
	} else {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			app.errorResponse(w, http.StatusNotFound, err)
			return
		}
		book, err := app.models.Books.Get(id)
		if err != nil {
			app.errorResponse(w, http.StatusNotFound, err)
			return
		}

		data := app.NewTemplateData(r)
		data.Book = book
		app.render(w, http.StatusOK, "books.View.html", data)
	}
}

func (app *Application) booksListHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Host, "api.") {
		app.booksListJsonHandler(w, r)
		return
	} else {
		books, err := app.models.Books.GetAll()
		app.logger.LogDebug.Printf("%+v\n", books)
		if err != nil {
			app.errorResponse(w, http.StatusInternalServerError, err)
			return
		}
		data := app.NewTemplateData(r)
		data.Books = books
		app.render(w, http.StatusOK, "books.List.html", data)
	}
}

func (app *Application) signupFormHandler(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = users.SignUpForm{V: validation.New()}
	app.render(w, http.StatusOK, "users.signup.html", data)
}

func (app *Application) POSTSignupHandler(w http.ResponseWriter, r *http.Request) {
	form := users.SignUpForm{V: validation.New()}

	app.ReadJSON(r, &form)
	app.logger.LogDebug.Printf("%+v", form)

	if form.Validate(); !form.V.Valid() {
		app.WriteJSON(w, r, form.V.Errors, nil, http.StatusUnprocessableEntity)
	}
	user := &users.User{Name: form.Name, Email: form.Email, Activated: false}
	if err := user.Password.Set(form.Password); err != nil {
		app.logger.LogError.Println(err)
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err := app.models.Users.Insert(user)
	if err != nil {
		app.logger.LogError.Println(err)

		switch {
		case errors.Is(err, internal.ErrDuplicateEmail):
			app.WriteJSON(w, r, envelope{"email": err.Error()}, nil, http.StatusUnprocessableEntity)
		default:
			app.WriteJSON(w, r, envelope{"errors": err.Error()}, nil, http.StatusInternalServerError)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "user successfully created!")
	app.logger.LogInfo.Printf("successfully inserted user: %+v\n", user)

	token, err := app.models.Tokens.New(user.ID, time.Hour*12, tokens.ScopeToken)
	if err != nil {
		app.logger.LogError.Println(err)
		app.WriteJSON(w, r, envelope{"errors": err}, nil, http.StatusInternalServerError)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.PlainText,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.html", data)
		if err != nil {
			app.logger.LogError.Printf("failed to send activation email to %s: %s", user.Email, err)
		}
	})

	app.logger.LogDebug.Printf("%+v", token.PlainText)
	app.WriteJSON(w, r, envelope{"userID": user.ID, "created_at": user.Created_at}, nil, http.StatusOK)
}

func (app *Application) loginFormHandler(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	app.render(w, http.StatusOK, "users.login.html", data)
}

func (app *Application) POSTLoginHandler(w http.ResponseWriter, r *http.Request) {
	form := users.LoginForm{V: validation.New()}
	err := app.ReadJSON(r, &form)
	if err != nil {
		app.WriteJSON(w, r, envelope{"error": err}, nil, http.StatusBadRequest)
		return
	}
	if form.Validate(); !form.V.Valid() {
		app.WriteJSON(w, r, envelope{"errors": form.V.Errors}, nil, http.StatusUnprocessableEntity)
		return
	}
	id, err := app.authenticateUser(form.Email, form.Password)
	if err != nil {
		app.logger.LogError.Println(err)
		switch {
		case errors.Is(err, internal.ErrInvalidCredentials):
			app.WriteJSON(w, r, envelope{"all": "Email or password is incorrect"}, nil, http.StatusUnprocessableEntity)

		case errors.Is(err, internal.ErrNotActivated):
			app.WriteJSON(w, r, envelope{"all": "the account exists but was not activated"}, nil, http.StatusUnprocessableEntity)

		default:
			app.errorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorResponse(w, http.StatusInternalServerError, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.logger.LogDebug.Printf("the id in the request context: %d", app.sessionManager.GetInt64(r.Context(), "authenticatedUserID"))
	w.WriteHeader(http.StatusOK)
}

func (app *Application) POSTLogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.logger.LogError.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	idstr := r.PathValue("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		app.errorResponse(w, http.StatusNotFound, err)
		return
	}

	tokenPlainText := r.PathValue("key")
	activated, err := app.activateUser(id, tokenPlainText)
	if err != nil {
		app.errorResponse(w, http.StatusNotFound, err)
		return
	}
	if activated {
		fmt.Fprintln(w, "your account was successfully activated")
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
