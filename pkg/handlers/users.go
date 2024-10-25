package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func GetUsersTable(w http.ResponseWriter, r *http.Request) {
	var err error
	data := components.UsersTable{
		Table: components.Table{
			TableURL: DestUsersTable,
		},
	}

	data.Table.Headings = []components.TableHeading{
		{database.OrderUserByID, "ID"},
		{database.OrderUserByRole, "Ruolo"},
		{database.OrderUserByUsername, "Username"},
		{database.OrderUserByName, "Nome"},
		{database.OrderUserBySurname, "Cognome"},
	}

	data.Table.OrderBy, err = strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		data.Table.OrderBy = database.OrderUserByID
	}
	data.Table.OrderDesc = r.URL.Query().Get("orderDesc") == "true"

	data.Users, err = appContext.Database(r).FindAllUsers(data.Table.OrderBy, data.Table.OrderDesc)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	appContext.ExecuteTemplate(w, r, "usersTable", data)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		User          database.User
		NameInput     components.Input
		SurnameInput  components.Input
		UsernameInput components.Input
		RoleSelect    components.Select
		IsNew         bool
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		data.IsNew = true
		data.User = database.User{}
	} else {
		user, err := appContext.Database(r).FindUser(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.User = user
	}

	data.NameInput = components.Input{"Nome", keyUserName, "text", data.User.Name}
	data.SurnameInput = components.Input{"Cognome", keyUserSurname, "text", data.User.Surname}
	data.UsernameInput = components.Input{"Username", keyUserUsername, "text", data.User.Username}

	roles, err := appContext.Database(r).FindAllRoles()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.RoleSelect = components.Select{"Ruolo", keyUserRoleID, data.User.RoleID, []components.SelectOption{}}
	for _, r := range roles {
		data.RoleSelect.Options = append(data.RoleSelect.Options, components.SelectOption{int(r.ID), r.Name})
	}

	appContext.ExecuteTemplate(w, r, "user.html", data)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	roleId, err := strconv.Atoi(r.FormValue(keyUserRoleID))
	if err != nil {
		HandleError(w, r, err)
		return
	}
	username := r.FormValue(keyUserUsername)
	name := r.FormValue(keyUserName)
	surname := r.FormValue(keyUserSurname)

	if isNew {
		password := r.FormValue(keyUserPassword)
		passwordHash, err := auth.HashPassword(password)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		err = appContext.Database(r).CreateUser(database.User{
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
		id, err := strconv.Atoi(r.FormValue(keyUserID))
		if err != nil {
			HandleError(w, r, err)
			return
		}

		if delete {
			err := appContext.Database(r).DeleteUser(id)
			if err != nil {
				HandleError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateUser(database.User{
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

	http.Redirect(w, r, DestUsers, http.StatusSeeOther)
}
