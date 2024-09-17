package main

import (
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
		Username string
		ErrorMsg string
	}

	claims, err := getSessionCookie(r)
	if err != nil {
		data.ErrorMsg = "Sessione scaduta"
	} else {
		data.Username = claims.Username
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	ok, _ := verifyPassword(username, password)
	if ok {
		setSessionCookie(w, username)
	} else {
		unsetSessionCookie(w)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
