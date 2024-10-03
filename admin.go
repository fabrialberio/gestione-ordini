package main

import (
	"gestione-ordini/database"
	"net/http"
	"strconv"
)

func checkPerm(w http.ResponseWriter, r *http.Request, permId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		http.Error(w, "Non autorizzato", http.StatusForbidden)
		return false
	}

	if ok, err := db.UserHasPermission(claims.UserID, permId); err != nil || !ok {
		http.Error(w, "Non autorizzato", http.StatusForbidden)
		return false
	}

	return true
}

func users(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "users.html", users)
}

func user(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		id = 0
	}

	user, err := db.GetUser(id)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "user.html", user)
}
