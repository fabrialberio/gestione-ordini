package main

import (
	"gestione-ordini/database"
	"net/http"
	"strconv"
)

func admin(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, checkRole(r, database.RoleIDAdministrator)
}

func adminUsers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, checkPerm(r, database.PermIDEditUsers)
}

func adminUsersTable(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkPerm(r, database.PermIDEditUsers); err != nil {
		return nil, err
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
		return nil, err
	}

	return data, nil
}

func adminUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if err := checkPerm(r, database.PermIDEditUsers); err != nil {
		return nil, err
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
			return nil, err
		}

		data.User = user
	}

	data.Roles, err = db.GetRoles()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func adminUserEdit(w http.ResponseWriter, r *http.Request) error {
	if err := checkPerm(r, database.PermIDEditUsers); err != nil {
		return err
	}

	isNew := r.FormValue("isNew") == "true"
	delete := r.Form.Has("delete")

	roleId, err := strconv.Atoi(r.FormValue("roleId"))
	if err != nil {
		return err
	}
	username := r.FormValue("username")
	name := r.FormValue("name")
	surname := r.FormValue("surname")

	if isNew {
		password := r.FormValue("password")
		passwordHash, err := hashPassword(password)
		if err != nil {
			return err
		}

		err = db.CreateUser(database.User{
			RoleID:       roleId,
			Username:     username,
			PasswordHash: passwordHash,
			Name:         name,
			Surname:      surname,
		})
		if err != nil {
			return err
		}
	} else {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			return err
		}

		if delete {
			err = db.DeleteUser(id)
			if err != nil {
				return err
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
				return err
			}
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	return nil
}
