package main

import (
	"gestione-ordini/database"
	"net/http"
)

func adminUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := getSessionCookie(r)
	if err != nil || claims.RoleID != database.RoleIDAdministrator {
		http.Error(w, "Non autorizzato", http.StatusForbidden)
		return
	}

	users, err := db.GetUsers()
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "users.html", users)
}
