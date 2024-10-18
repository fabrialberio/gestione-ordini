package handlers

import (
	"gestione-ordini/pkg/auth"
	"gestione-ordini/pkg/database"
	"gestione-ordini/pkg/reqContext"
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

	err := reqContext.GetRequestContext(r).AuthenticationErr
	if err == auth.ErrNoCookie {
		data.ErrorMsg = html.EscapeString(r.URL.Query().Get("errormsg"))
	} else if err != nil {
		data.ErrorMsg = "Sessione scaduta"
		auth.UnsetAuthenticatedUser(w)
	} else {
	}

	reqContext.GetRequestContext(r).Templ.ExecuteTemplate(w, "login.html", data)
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	dest := ""

	var ok bool
	user, _ := reqContext.GetRequestContext(r).DB.FindUserWithUsername(username)
	if user != nil {
		ok = auth.VerifyPassword(user.PasswordHash, password)
	} else {
		ok = false
	}

	if ok {
		auth.SetAuthenticatedUser(w, user.ID, user.RoleID)
		switch user.RoleID {
		case database.RoleIDCook:
			dest = "/cook"
		case database.RoleIDManager:
			dest = "/manager"
		case database.RoleIDAdministrator:
			dest = "/admin"
		}
	} else {
		auth.UnsetAuthenticatedUser(w)
		dest = "/?errormsg=Password errata"
	}

	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	auth.UnsetAuthenticatedUser(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}