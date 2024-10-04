package main

import (
	"gestione-ordini/database"
	"html"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID   int
		RoleID   int
		ErrorMsg string
	}

	claims, err := getSessionCookie(r)
	if err == http.ErrNoCookie {
		data.ErrorMsg = html.EscapeString(r.URL.Query().Get("errormsg"))
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		unsetSessionCookie(w)
	} else {
		data.UserID = claims.UserID
		data.RoleID = claims.RoleID
	}

	templ.ExecuteTemplate(w, "index.html", data)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	dest := ""

	var ok bool
	user, _ := db.GetUserByUsername(username)
	if user != nil {
		ok = verifyPassword(user.PasswordHash, password)
	} else {
		ok = false
	}

	if ok {
		setSessionCookie(w, user.ID, user.RoleID)
		switch user.RoleID {
		case database.RoleIDCook:
			dest = "/cook"
		case database.RoleIDManager:
			dest = "/manager"
		case database.RoleIDAdministrator:
			dest = "/admin"
		}
	} else {
		unsetSessionCookie(w)
		dest = "/?errormsg=Password errata"
	}

	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {
	unsetSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
