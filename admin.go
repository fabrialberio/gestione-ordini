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

func checkRole(w http.ResponseWriter, r *http.Request, roleId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		http.Error(w, "Non autorizzato", http.StatusForbidden)
		return false
	}

	if claims.RoleID != roleId {
		http.Error(w, "Non autorizzato", http.StatusForbidden)
		return false
	}

	return true
}

func admin(w http.ResponseWriter, r *http.Request) {
	if !checkRole(w, r, database.RoleIDAdministrator) {
		return
	}

	templ.ExecuteTemplate(w, "admin.html", nil)
}

func users(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	var err error
	var data struct {
		OrderBy int
		Users   []database.User
	}

	data.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.OrderBy = database.UserOrderByID
	}

	data.Users, err = db.GetUsers(data.OrderBy)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "users.html", data)
}

func user(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	var data struct {
		User  *database.User
		IsNew bool
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		data.IsNew = true
		data.User = &database.User{}
	} else {
		user, err := db.GetUser(id)
		if err != nil {
			http.Error(w, "Errore interno", http.StatusInternalServerError)
			return
		}

		data.User = user
	}

	templ.ExecuteTemplate(w, "user.html", data)
}
