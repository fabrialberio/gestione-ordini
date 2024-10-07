package main

import (
	"gestione-ordini/database"
	"net/http"
	"strconv"
)

func displayUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Non autorizzato", http.StatusForbidden)
}

func displayInternalError(w http.ResponseWriter) {
	http.Error(w, "Errore interno", http.StatusInternalServerError)
}

func checkPerm(w http.ResponseWriter, r *http.Request, permId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		displayUnauthorized(w)
		return false
	}

	if ok, err := db.UserHasPerm(claims.UserID, permId); err != nil || !ok {
		displayUnauthorized(w)
		return false
	}

	return true
}

func checkRole(w http.ResponseWriter, r *http.Request, roleId int) bool {
	claims, err := getSessionCookie(r)
	if err != nil {
		displayUnauthorized(w)
		return false
	}

	if claims.RoleID != roleId {
		displayUnauthorized(w)
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

func usersPage(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	templ.ExecuteTemplate(w, "usersPage.html", nil)
}

func usersTable(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	var err error
	var data struct {
		OrderBy   int
		OrderDesc bool
		Headers   interface{}
		Users     []database.User
	}

	data.Headers = []struct {
		Index int
		Name  string
	}{
		{database.UserOrderByID, "ID"},
		{database.UserOrderByRole, "Ruolo"},
		{database.UserOrderByUsername, "Username"},
		{database.UserOrderByName, "Nome"},
		{database.UserOrderBySurname, "Cognome"},
	}

	data.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.OrderBy = database.UserOrderByID
	}
	data.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Users, err = db.GetUsers(data.OrderBy, data.OrderDesc)
	if err != nil {
		displayInternalError(w)
		return
	}

	templ.ExecuteTemplate(w, "usersTable.html", data)
}

func usersEdit(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	var data struct {
		User  *database.User
		Roles []database.Role
		IsNew bool
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		data.IsNew = true
		data.User = &database.User{}
	} else {
		user, err := db.GetUser(id)
		if err != nil {
			displayInternalError(w)
			return
		}

		data.User = user
	}

	data.Roles, err = db.GetRoles()
	if err != nil {
		displayInternalError(w)
		return
	}

	templ.ExecuteTemplate(w, "user.html", data)
}

func usersApplyEdit(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	isNew := r.FormValue("isNew") == "true"

	roleId, err := strconv.Atoi(r.FormValue("roleId"))
	if err != nil {
		displayInternalError(w)
		return
	}
	username := r.FormValue("username")
	name := r.FormValue("name")
	surname := r.FormValue("surname")

	if isNew {
		password := r.FormValue("password")
		passwordHash, err := hashPassword(password)
		if err != nil {
			displayInternalError(w)
			return
		}

		err = db.AddUser(database.User{
			RoleID:       roleId,
			Username:     username,
			PasswordHash: passwordHash,
			Name:         name,
			Surname:      surname,
		})
		if err != nil {
			displayInternalError(w)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			displayInternalError(w)
			return
		}

		err = db.UpdateUser(database.User{
			ID:       id,
			RoleID:   roleId,
			Username: username,
			Name:     name,
			Surname:  surname,
		})
		if err != nil {
			displayInternalError(w)
			return
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
