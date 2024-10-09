package main

import (
	"gestione-ordini/database"
	"html"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

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

	return data, nil
}

func login(w http.ResponseWriter, r *http.Request) error {
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
	return nil
}

func logout(w http.ResponseWriter, r *http.Request) error {
	unsetSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
