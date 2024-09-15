package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"site/data/tokens"
	"site/internal"
	"time"
)

func (config *Config) OpenConnection() (*sql.DB, error) {
	connString := fmt.Sprintf(`host=%s port=%s dbname=%s user=%s password=%s sslmode=disable`,
		config.DB.DbHost, config.DB.DbPort, config.DB.DbName, config.DB.DbUser, config.DB.DbPassword,
	)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(100)
	dbConn.SetMaxIdleConns(100)
	dbConn.SetConnMaxIdleTime(time.Hour)
	return dbConn, nil
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *templatesData) {
	ts, ok := app.templatesChache[page]
	if !ok {
		app.logger.LogError.Println(errors.New("template not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.logger.LogError.Println(errors.New("template not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.LogError.Println(errors.New("template not found"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *Application) ReadJSON(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) WriteJSON(w http.ResponseWriter, r *http.Request, data any, header http.Header, status int) error {
	qs := r.URL.Query()
	print := app.ReadString(qs, "print", "")

	if print == "pretty" {
		d, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		for key, val := range header {
			w.Header()[key] = val
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(d)

		return nil
	} else {
		d, err := json.Marshal(data)
		if err != nil {
			return err
		}
		for key, val := range header {
			w.Header()[key] = val
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(d)
		return nil
	}
}

func (app *Application) ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *Application) activateUser(id int64, tokenPlainText string) (bool, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))
	tokenHashslice := tokenHash[:]

	userTokenHash, err := app.models.Tokens.GetUserToken(id)
	if err != nil {
		return false, err
	}

	user, err := app.models.Users.GetUserByID(id)
	if err != nil {
		return false, err
	}
	if user.Activated {
		return true, errors.New("user is already activated")
	}

	if bytes.Equal(tokenHashslice, userTokenHash) {
		user.Activated = true
		app.logger.LogInfo.Println("user.Activated = true")
	} else {
		return false, errors.New("invalid token")
	}

	if err := app.models.Users.UpdateUser(user); err != nil {
		return false, err
	}

	if err = app.models.Tokens.DeleteAllForUser(id, tokens.ScopeToken); err != nil {
		return false, err
	}

	return true, nil
}

func (app *Application) authenticateUser(email, password string) (int64, error) {
	user, err := app.models.Users.GetUserByEmail(email)
	if err != nil {
		return -1, err
	}
	if !user.Activated {
		return -1, internal.ErrNotActivated
	}

	match, err := user.Password.Matches(password)
	if err != nil {
		return -1, err
	}

	if !match {
		return -1, internal.ErrInvalidCredentials
	}

	return user.ID, nil
}

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")

func (app *Application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *Application) background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.LogError.Printf("%s", err)
			}
		}()
		fn()
	}()
}
