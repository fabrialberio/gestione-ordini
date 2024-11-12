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
	orderBy, err := strconv.Atoi(r.URL.Query().Get("orderBy"))
	if err != nil {
		orderBy = database.OrderUserByID
	}
	orderDesc := r.URL.Query().Get("orderDesc") == "true"

	users, err := appContext.Database(r).FindAllUsers(orderBy, orderDesc)
	if err != nil {
		logError(r, err)
	}

	data := components.UsersTable{
		Table: components.Table{
			TableURL:  DestUsersTable,
			OrderBy:   orderBy,
			OrderDesc: orderDesc,
			Headings: []components.TableHeading{
				{Index: database.OrderUserByID, Name: "ID"},
				{Index: database.OrderUserByRole, Name: "Ruolo"},
				{Index: database.OrderUserByUsername, Name: "Username"},
				{Index: database.OrderUserByName, Name: "Nome"},
				{Index: database.OrderUserBySurname, Name: "Cognome"},
			},
		},
		Users: users,
	}

	appContext.ExecuteTemplate(w, r, "usersTable", data)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	} else if user.RoleID != database.RoleIDAdministrator {
		ShowItemNotAllowedError(w, r, auth.ErrInvalidRole)
		return
	}

	isNew := false
	defaultUser := database.User{}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		isNew = true
	} else {
		defaultUser, err = appContext.Database(r).FindUser(id)
		if err != nil {
			ShowItemNotAllowedError(w, r, err)
			return
		}
	}

	roles, err := appContext.Database(r).FindAllRoles()
	if err != nil {
		ShowDatabaseQueryError(w, r, err)
		return
	}

	roleOptions := []components.SelectOption{}
	for _, r := range roles {
		roleOptions = append(roleOptions, components.SelectOption{Value: int(r.ID), Text: r.Name})
	}

	data := struct {
		IsNew         bool
		User          database.User
		NameInput     components.Input
		SurnameInput  components.Input
		UsernameInput components.Input
		PasswordInput components.Input
		RoleSelect    components.Select
	}{
		IsNew: isNew,
		User:  defaultUser,
		NameInput: components.Input{
			Label:        "Nome",
			Name:         keyUserName,
			DefaultValue: defaultUser.Name,
		},
		SurnameInput: components.Input{
			Label:        "Cognome",
			Name:         keyUserSurname,
			DefaultValue: defaultUser.Surname,
		},
		UsernameInput: components.Input{
			Label:        "Username",
			Name:         keyUserUsername,
			DefaultValue: defaultUser.Username,
		},
		PasswordInput: components.Input{
			Label: "Password",
			Name:  keyUserPassword,
		},
		RoleSelect: components.Select{
			Label:    "Ruolo",
			Name:     keyUserRoleID,
			Selected: defaultUser.RoleID,
			Options:  roleOptions,
		},
	}

	appContext.ExecuteTemplate(w, r, "user.html", data)
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	user, err := appContext.AuthenticatedUser(r)
	if err != nil {
		LogoutError(w, r, err)
		return
	} else if user.RoleID != database.RoleIDAdministrator {
		ShowItemNotAllowedError(w, r, auth.ErrInvalidRole)
		return
	}

	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	roleId, err := strconv.Atoi(r.FormValue(keyUserRoleID))
	if err != nil {
		ShowItemInvalidFormError(w, r, err)
		return
	}
	username := r.FormValue(keyUserUsername)
	password := r.FormValue(keyUserPassword)
	name := r.FormValue(keyUserName)
	surname := r.FormValue(keyUserSurname)
	passwordHash := ""

	if password != "" {
		passwordHash, err = auth.HashPassword(password)
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}
	}

	if isNew {
		err = appContext.Database(r).CreateUser(database.User{
			RoleID:       roleId,
			Username:     username,
			PasswordHash: passwordHash,
			Name:         name,
			Surname:      surname,
		})
		if err != nil {
			ShowDatabaseQueryError(w, r, err)
			return
		}
	} else {
		id, err := strconv.Atoi(r.FormValue(keyUserID))
		if err != nil {
			ShowItemInvalidFormError(w, r, err)
			return
		}

		if delete {
			err := appContext.Database(r).DeleteUser(id)
			if err != nil {
				ShowItemNotDeletableError(w, r, err)
				return
			}
		} else {
			err = appContext.Database(r).UpdateUser(database.User{
				ID:           id,
				RoleID:       roleId,
				Username:     username,
				PasswordHash: passwordHash,
				Name:         name,
				Surname:      surname,
			})
			if err != nil {
				ShowDatabaseQueryError(w, r, err)
				return
			}
		}
	}

	http.Redirect(w, r, DestUsers, http.StatusSeeOther)
}
