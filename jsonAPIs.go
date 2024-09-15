package main

import (
	"net/http"
	"strconv"
)

type envelope map[string]any

func (app *Application) booksListJsonHandler(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Books.GetAll()
	if err != nil {
		app.WriteJSON(w, r, envelope{"error": err}, nil, http.StatusInternalServerError)
		return
	}

	app.WriteJSON(w, r, envelope{"books": books}, nil, http.StatusOK)
}

func (app *Application) bookViewJsonHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		app.WriteJSON(w, r, envelope{"error": err}, nil, http.StatusNotFound)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		app.WriteJSON(w, r, envelope{"error": err}, nil, http.StatusNotFound)
		return
	}

	app.WriteJSON(w, r, envelope{"book": book}, nil, http.StatusOK)
}
