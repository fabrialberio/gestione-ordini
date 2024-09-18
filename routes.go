package main

import (
	"fmt"
	"log"
	"net/http"
)

func logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h(w, r)
	}
}

func onlyPost(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		} else {
			h(w, r)
		}
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var data struct {
		IDUtente int
		IDRuolo  int
		ErrorMsg string
	}

	claims, err := getSessionCookie(r)
	if err == http.ErrNoCookie {
		data.ErrorMsg = r.URL.Query().Get("errormsg")
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		unsetSessionCookie(w)
	} else {
		data.IDUtente = claims.IDUtente
		data.IDRuolo = claims.IDRuolo
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	errorMsg := ""

	utente, _ := db.GetUtenteByUsername(username)
	ok := verifyPassword(utente, password)
	if ok {
		setSessionCookie(w, utente.ID, utente.IDRuolo)
	} else {
		unsetSessionCookie(w)
		errorMsg = "?errormsg=Password errata"
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", errorMsg), http.StatusSeeOther)
}
