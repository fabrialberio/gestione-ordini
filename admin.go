package main

import (
	"gestione-ordini/database"
	"net/http"
	"strconv"
)

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
		htmlError(w, http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "usersTable.html", data)
}

func usersEdit(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	var data struct {
		User  database.User
		Roles []database.Role
		IsNew bool
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		data.IsNew = true
		data.User = database.User{}
	} else {
		user, err := db.GetUser(id)
		if err != nil {
			htmlError(w, http.StatusInternalServerError)
			return
		}

		data.User = user
	}

	data.Roles, err = db.GetRoles()
	if err != nil {
		htmlError(w, http.StatusInternalServerError)
		return
	}

	templ.ExecuteTemplate(w, "user.html", data)
}

func usersApplyEdit(w http.ResponseWriter, r *http.Request) {
	if !checkPerm(w, r, database.PermIDEditUsers) {
		return
	}

	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	roleId, err := strconv.Atoi(r.FormValue("roleId"))
	if err != nil {
		htmlError(w, http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	name := r.FormValue("name")
	surname := r.FormValue("surname")

	if isNew {
		password := r.FormValue("password")
		passwordHash, err := hashPassword(password)
		if err != nil {
			htmlError(w, http.StatusInternalServerError)
			return
		}

		err = db.CreateUser(database.User{
			RoleID:       roleId,
			Username:     username,
			PasswordHash: passwordHash,
			Name:         name,
			Surname:      surname,
		})
		if err != nil {
			htmlError(w, http.StatusInternalServerError)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			htmlError(w, http.StatusInternalServerError)
			return
		}

		if delete {
			err = db.DeleteUser(id)
			if err != nil {
				htmlError(w, http.StatusInternalServerError)
				return
			}
		} else {
			err = db.UpdateUser(database.User{
				ID:       id,
				RoleID:   roleId,
				Username: username,
				Name:     name,
				Surname:  surname,
			})
			if err != nil {
				htmlError(w, http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
