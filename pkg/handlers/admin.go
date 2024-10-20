package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/components"
	"gestione-ordini/pkg/database"
	"net/http"
	"strconv"
)

func adminSidebar(selected int) []components.SidebarDest {
	sidebar := []components.SidebarDest{
		{destAdminUsers, "fa-users", "Utenti", false},
		{destAdminProducts, "fa-box-open", "Prodotti", false},
	}
	sidebar[selected].Selected = true

	return sidebar
}

func GetAdmin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, destAdminUsers, http.StatusSeeOther)
}

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: adminSidebar(0),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminUsers.html", data)
}

func GetAdminProducts(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Sidebar []components.SidebarDest
	}{
		Sidebar: adminSidebar(1),
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminProducts.html", data)
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

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminUsersTable.html", data)
}

func GetAdminUser(w http.ResponseWriter, r *http.Request) {
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
		user, err := appContext.FromRequest(r).DB.FindUser(id)
		if err != nil {
			HandleError(w, r, err)
			return
		}

		data.User = user
	}

	data.NameInput = components.Input{"Nome", keyUserName, "text", data.User.Name}
	data.SurnameInput = components.Input{"Cognome", keyUserSurname, "text", data.User.Surname}
	data.UsernameInput = components.Input{"Username", keyUserUsername, "text", data.User.Username}

	roles, err := appContext.FromRequest(r).DB.FindAllRoles()
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data.RoleSelect = components.Select{"Ruolo", keyUserRoleID, data.User.RoleID, []components.SelectOption{}}
	for _, r := range roles {
		data.RoleSelect.Options = append(data.RoleSelect.Options, components.SelectOption{int(r.ID), r.Name})
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "adminUser.html", data)
}

func PostAdminUser(w http.ResponseWriter, r *http.Request) {
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
		id, err := strconv.Atoi(r.FormValue(keyUserID))
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

	http.Redirect(w, r, destAdminUsers, http.StatusSeeOther)
}
