package handlers

import (
	"errors"
	"gestione-ordini/pkg/appContext"
	"log"
	"net/http"
)

var (
	ErrItemNotFound     = errors.New("item not found")
	ErrItemNotAllowed   = errors.New("item access restricted")
	ErrItemNotDeletable = errors.New("item cannot be deleted")
	ErrItemInvalidForm  = errors.New("item form data is invalid")
	ErrLogout           = errors.New("authentication error")
	ErrDatabaseQuery    = errors.New("database query error")
)

func ShowItemNotAllowedError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	executeErrorTemplate(w, r, "Accesso non consentito", err)
}

func ShowItemNotDeletableError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	executeErrorTemplate(w, r, "Eliminazione non consentita", err)
}

func ShowItemInvalidFormError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	executeErrorTemplate(w, r, "Richiesta non valida", err)
}

func ShowDatabaseQueryError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	executeErrorTemplate(w, r, "Errore interno", err)
}

func LogoutError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	logoutRedirect(w, r, true)
}

func logError(r *http.Request, err error) {
	log.Printf("%s %s %s Error: %v", r.RemoteAddr, r.Method, r.URL, err)
}

func executeErrorTemplate(w http.ResponseWriter, r *http.Request, title string, err error) {
	data := struct {
		Title    string
		ErrorMsg string
	}{
		Title:    title,
		ErrorMsg: err.Error(),
	}

	appContext.ExecuteTemplate(w, r, "error.html", data)
}
