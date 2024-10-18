package handlers

import (
	"gestione-ordini/pkg/appContext"
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"html"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	var data struct {
		ErrorMsg string
	}

	err := appContext.FromRequest(r).AuthenticationErr
	if err == auth.ErrNoCookie {
		data.ErrorMsg = html.EscapeString(r.URL.Query().Get("errormsg"))
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		auth.UnsetAuthenticatedUser(w)
	} else {
	}

	appContext.FromRequest(r).Templ.ExecuteTemplate(w, "login.html", data)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, _ := appContext.FromRequest(r).DB.FindUserWithUsername(username)
	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := auth.SetAuthenticatedUser(w, user, password)
	if err != nil {
		http.Redirect(w, r, "/?errormsg=Password errata", http.StatusSeeOther)
		return
	}

	var dest string
	switch user.RoleID {
	case database.RoleIDCook:
		dest = "/cook"
	case database.RoleIDManager:
		dest = "/manager"
	case database.RoleIDAdministrator:
		dest = "/admin"
	}

	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	auth.UnsetAuthenticatedUser(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
