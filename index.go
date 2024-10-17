package main

import (
	"gestione-ordini/database"
	"html"
	"net/http"
)

func HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	var data struct {
		ErrorMsg string
	}

	_, err := GetAuthenticatedUser(r)
	if err == ErrNoCookie {
		data.ErrorMsg = html.EscapeString(r.URL.Query().Get("errormsg"))
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		UnsetAuthenticatedUser(w)
	} else {
	}

	templ.ExecuteTemplate(w, "login.html", data)
}

func HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	dest := ""

	var ok bool
	user, _ := db.FindUserWithUsername(username)
	if user != nil {
		ok = verifyPassword(user.PasswordHash, password)
	} else {
		ok = false
	}

	if ok {
		SetAuthenticatedUser(w, user.ID, user.RoleID)
		switch user.RoleID {
		case database.RoleIDCook:
			dest = "/cook"
		case database.RoleIDManager:
			dest = "/manager"
		case database.RoleIDAdministrator:
			dest = "/admin"
		}
	} else {
		UnsetAuthenticatedUser(w)
		dest = "/?errormsg=Password errata"
	}

	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func HandlePostLogout(w http.ResponseWriter, r *http.Request) {
	UnsetAuthenticatedUser(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
