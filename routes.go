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

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var errorMsg string

		username := r.FormValue("username")
		password := r.FormValue("password")
		unsetSessionCookie(w)

		ok, err := verifyPassword(username, password)
		if err != nil || !ok {
			errorMsg = "Password errata"
		} else {
			err = setSessionCookie(w, username)
			if err != nil {
				log.Printf("Error setting cookie: %v", err)
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/?msg=%s", errorMsg), http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		var data struct {
			Username string
			ErrorMsg string
		}

		claims, err := getSessionCookie(r)
		if err == http.ErrNoCookie {
			data.ErrorMsg = r.URL.Query().Get("msg")
		} else if err != nil {
			log.Printf("Error getting cookie: %v", err)
		} else {
			data.Username = claims.Username
		}

		templates.ExecuteTemplate(w, "index.html", data)
	}
}
