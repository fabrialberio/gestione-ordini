package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "admin.html", nil)
}

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "users.html", nil)
}

func GetAdminUsersTable(w http.ResponseWriter, r *http.Request) {
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

	data.Users, err = appContext.FromRequest(r).DB.FindAllUsers(data.OrderBy, data.OrderDesc)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "usersTable.html", data)
}

func GetAdminUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		User  database.User
		Roles []database.Role
		IsNew bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.User = database.User{}
	} else {
		user, err := appContext.FromRequest(r).DB.FindUser(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.User = user
	}

	data.Roles, err = appContext.FromRequest(r).DB.FindAllRoles()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "user.html", data)
}

func PostAdminUser(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	roleId, err := strconv.Atoi(r.FormValue("roleId"))
	if err != nil {
		HandleError(w, r, err)
		return
	}
	username := r.FormValue("username")
	name := r.FormValue("name")
	surname := r.FormValue("surname")

	if isNew {
		password := r.FormValue("password")
		passwordHash, err := auth.HashPassword(password)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		err = appContext.FromRequest(r).DB.CreateUser(database.User{
			RoleID:       roleId,
			Username:     username,
			PasswordHash: passwordHash,
			Name:         name,
			Surname:      surname,
		})
		if err != nil {
			HandleError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err := appContext.FromRequest(r).DB.DeleteUser(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.FromRequest(r).DB.UpdateUser(database.User{
				ID:       id,
				RoleID:   roleId,
				Username: username,
				Name:     name,
				Surname:  surname,
			})
			if err != nil {
				HandleError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
