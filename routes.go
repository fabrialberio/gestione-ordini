package main

import (
	"log"
	"net/http"
)

type IndexData struct {
	Username string
	ErrorMsg string
}

func logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		h(w, r)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		indexGet(w, r)
	} else if r.Method == http.MethodPost {
		indexPost(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func indexGet(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		cookie = &http.Cookie{}
	}

	claims, err := validateJWT(cookie.Value)
	if err != nil {
		claims = &UserClaims{}
	}

	templates.ExecuteTemplate(w, "index.html", IndexData{
		Username: claims.Username,
		ErrorMsg: r.URL.Query().Get("msg"),
	})
}

func indexPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := login(username, password, w)
	if err != nil {
		logout(w)
		http.Redirect(w, r, "/?msg=Password errata", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
