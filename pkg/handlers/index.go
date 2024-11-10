package handlers

import (
	"errors"
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"net/http"
)

const dateFormat = "2006-01-02"

func loginRedirect(w http.ResponseWriter, r *http.Request, roleId int) {
	var dest string
	if roleId == database.RoleIDChef {
		dest = DestChef
	} else {
		dest = DestConsole
	}

	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func logoutRedirect(w http.ResponseWriter, r *http.Request, errorMsg bool) {
	var dest string
	if errorMsg {
		dest = "/?errormsg"
	} else {
		dest = "/"
	}

	auth.UnsetAuthenticatedUser(w)
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func GetIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	var data struct {
		ErrorMsg string
	}

	user, err := appContext.AuthenticatedUser(r)
	if errors.Is(err, auth.ErrNoCookie) {
		if r.URL.Query().Has("errormsg") {
			data.ErrorMsg = "Utente o password errati"
		}
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
	} else {
		loginRedirect(w, r, user.RoleID)
		return
	}

	appContext.ExecuteTemplate(w, r, "login.html", data)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := appContext.Database(r).FindUserWithUsername(username)
	if err != nil {
		logoutRedirect(w, r, true)
		return
	}

	err = auth.SetAuthenticatedUser(w, user, password)
	if err != nil {
		logoutRedirect(w, r, true)
		return
	}

	loginRedirect(w, r, user.RoleID)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logoutRedirect(w, r, false)
}
