package main

import (
	"fmt"
	"html"
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
		UserID   int
		RoleID   int
		ErrorMsg string
	}

	claims, err := getSessionCookie(r)
	if err == http.ErrNoCookie {
		data.ErrorMsg = html.EscapeString(r.URL.Query().Get("errormsg"))
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		unsetSessionCookie(w)
	} else {
		data.UserID = claims.UserID
		data.RoleID = claims.RoleID
	}

	templates.ExecuteTemplate(w, "index.html", data)
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	errorMsg := ""

	var ok bool
	user, _ := db.GetUserByUsername(username)
	if user != nil {
		ok = verifyPassword(user.PasswordHash, password)
	} else {
		ok = false
	}

	if ok {
		setSessionCookie(w, user.ID, user.RoleID)
	} else {
		unsetSessionCookie(w)
		errorMsg = "?errormsg=Password errata"
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", errorMsg), http.StatusSeeOther)
}
